 package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"strconv"
	"strings"
)

func makeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func check (e error) {
	if(e != nil) {
		panic(e);
	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func getMapStringIntKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func parseFabricString(str string) (int, int, int, int, int)  {
	index1 := 1
	index2 := strings.Index(str, " @ ")
	index3 := strings.Index(str, ",")
	index4 := strings.Index(str, ": ")
	index5 := strings.Index(str, "x")
	
	id, err := strconv.Atoi(str[index1:index2])
	x1, err := strconv.Atoi(str[index2 + 3:index3])
	y1, err := strconv.Atoi(str[index3+1:index4])
	width, err := strconv.Atoi(str[index4+2:index5])
	height, err := strconv.Atoi(str[index5+1:])
	
	check(err)
	return  id, x1, y1, width, height
}

func main() {
	starttime := makeTimestamp()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	m := make(map[string]int)
	
	for i:= 0; i < len(lines); i++ {
		_, x1, y1, width, height := parseFabricString(lines[i])
		for j:=x1; j <= x1+width - 1; j++ {
			for k:=y1; k <= y1+height - 1; k++ {
				index := strconv.Itoa(j) + "x" + strconv.Itoa(k)
				m[index]++
			}
		}
	}
	// Part A
	overlapCount:= 0
	keys:=getMapStringIntKeys(m)
	for i:=0; i < len(keys); i++ {
		if m[keys[i]] >= 2  {
			overlapCount++
		}
	}
	fmt.Println("Result A:", overlapCount)
	
	// Part B
	for i:= 0; i < len(lines); i++ {
		noOverlap:= true
		id, x1, y1, width, height := parseFabricString(lines[i])
		for j:=x1; j <= x1+width - 1; j++ {
			for k:=y1; k <= y1+height - 1; k++ {
				index := strconv.Itoa(j) + "x" + strconv.Itoa(k)
				if m[index] > 1 {
					noOverlap = false;
					break
				}
			}
			if(!noOverlap) {
				break
			}
		}
		
		if(noOverlap) {
			fmt.Println("Result B:", id)
			break
		}
	}
	
	endtime := makeTimestamp()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}