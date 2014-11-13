package main

import ( 
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
)



func calc(i int, c int, values [] int) (x, y int) {
	var val1 int
	for itr := 0; itr < c - 1; itr++ {
		val1 = values[itr]
		for itr2 := itr + 1; itr2 < c; itr2++ {
			if i - val1 == values[itr2] {
				x = itr + 1
				y = itr2 + 1
				return
			}
		}
	}
	x , y = -1, -1
	return
		 
}
/*
func readCase(fileHandle *File) (i int, c int, values [] int) 
{
	
}
*/
func main() {
	
	fileHandle, _ := os.Open("storedata.txt")
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)
	
	fileScanner.Scan() 
	
	numCases, _ := strconv.ParseInt(fileScanner.Text(), 10, strconv.IntSize)
	
	var c, i int 
	var itr int64
	
	for itr = 1 ; itr <= numCases; itr++ {
		fileScanner.Scan()
		c, _ = strconv.Atoi(fileScanner.Text())
		fileScanner.Scan()
		i, _ = strconv.Atoi(fileScanner.Text())
        fileScanner.Scan()
		arrStrings := strings.Split(fileScanner.Text(), " ")
		data := make ([] int, len(arrStrings))
		for idx, str := range arrStrings {
			data[idx], _ = strconv.Atoi(str)
		}
		num1, num2 := calc(c, i, data)
		fmt.Println("Case #",itr,": ", num1, num2)
	}
}

/*

Input 

Output 
3
100
3
5 75 25
200
7
150 24 79 50 88 345 3
8
8
2 1 9 4 4 56 90 3
Case #1: 2 3
Case #2: 1 4
Case #3: 4 5
*/
