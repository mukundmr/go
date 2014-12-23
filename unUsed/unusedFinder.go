package main

/*
   find used properties from a properties file in all files contaned under a source dir
*/

import (
  "bufio"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
  "sync"
  "sync/atomic"
)

func main() {
  var mutex sync.Mutex;

  var propertyList atomic.Value
  
  // create closure functions
  read := func(key string) bool {
    map1 := propertyList.Load().(map[string]bool)
    return map1[key]
  }
  
  update := func(key string, val bool) {
    mutex.Lock()
    defer mutex.Unlock()
    map1 := propertyList.Load().(map[string]bool)
    map2 := make(map[string]bool)
    for k,v := range map1 {
      map2[k]=v
    }
    map2[key] = val
    propertyList.Store(map2)
  }

  /* set the values as per environment */
  propertyFile := ""
  searchDir := ""

  propertyList.Store(make(map[string]bool)) 
  fileHandle_en, _ := os.Open(propertyFile)
  defer fileHandle_en.Close()
  fileScanner_en := bufio.NewScanner(fileHandle_en)
  for fileScanner_en.Scan() {
    if (len(fileScanner_en.Text()) > 0) && !strings.HasPrefix(fileScanner_en.Text(), "#") {
      update(strings.Split(fileScanner_en.Text(), "=")[0], true)
    }
  }
  
  // build file list

  fileList := []string{}
  
  err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
    if !f.IsDir() {
      fileList = append(fileList, path)
    }
    return nil
  })
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  num_file := len(fileList)
  fmt.Println("number of files", num_file)
  
  // create a wait group to ensure all files are processed
  var wg sync.WaitGroup
  wg.Add(num_file)
  map1 := propertyList.Load().(map[string]bool)
  var keyList = make([]string , len(map1))
  ctr := 0
  for key, _ := range map1 {
    keyList[ctr] = key
    ctr++
  }
  for _, file := range fileList {
    go func(file string) {
        data, _ := ioutil.ReadFile(file)
        for key := range keyList {
          if read(keyList[key]) {
            if strings.Contains(string(data), keyList[key]) {
              update(keyList[key], false)
            }
          }
        }
        wg.Done()
    }(file)
  }
  wg.Wait();
  map1 = propertyList.Load().(map[string]bool)
  for key, _ := range map1 {
    if map1[key] {
      fmt.Println(key)
    }
  }
  // end of main
}



