package main

/*
   find used properties from a propeties file in all files contaned under a source dir
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
	type PropertyMap map[string]bool
	var mutex sync.Mutex
	var propertyList atomic.Value
	propertyList.Store(make(PropertyMap))

	// create closure functions
	read := func(key string) bool {
		map1 := propertyList.Load().(PropertyMap)
		return map1[key]
	}
	update := func(key string, val bool) {
		mutex.Lock()
		defer mutex.Unlock()
		map1 := propertyList.Load().(PropertyMap)
		map2 := make(PropertyMap)
		for k, v := range map1 {
			map2[k] = v
		}
		map2[key] = val
		propertyList.Store(map2)
	}

	/* set the values as per environment */
	propertyFile := ""
	searchDir := ""

	fileHandle_en, _ := os.Open(propertyFile)
	fileScanner_en := bufio.NewScanner(fileHandle_en)
	var keyStr string
	for fileScanner_en.Scan() {
		if (len(fileScanner_en.Text()) > 0) && !strings.HasPrefix(fileScanner_en.Text(), "#") {
			keyStr = strings.Split(fileScanner_en.Text(), "=")[0]
			update(keyStr, true)
		}
	}
	fileHandle_en.Close() // can defer, but lets not waste file handles

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
	fmt.Println("number of files", num_file, ".  number of properties", len(propertyList.Load().(PropertyMap)))

	// create a wait group to ensure all files are processed
	var wg sync.WaitGroup
	wg.Add(num_file)
	map1 := propertyList.Load().(PropertyMap)
	var keyList = make([]string, len(map1)) //using a slice & appending makes it slower
	ctr := 0
	for key, _ := range map1 {
		keyList[ctr] = key
		ctr++
	}
	// lets limit number of files being handled at any given time
	var cblock = make(chan int, 100)
	for _, file := range fileList {
		go func(file string) {
			cblock <- 1
			data, _ := ioutil.ReadFile(file)
			strdata := string(data)
			for key := range keyList {
				if read(keyList[key]) {
					if strings.Contains(strdata, keyList[key]) {
						update(keyList[key], false)
					}
				}
			}
			wg.Done()
			<-cblock
		}(file)
	}
	wg.Wait()
	map1 = propertyList.Load().(PropertyMap)
	for key, val := range map1 {
		if val {
			fmt.Println(key)
		}
	}
	// end of main
}
