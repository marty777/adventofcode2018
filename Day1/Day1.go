 package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"strconv"
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

func main() {
	starttime := makeTimestamp()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	sum := 0;
	for i:= 0; i < len(lines); i++ {
		val, err := strconv.Atoi(lines[i])
		check(err)
		sum += val
	}
	fmt.Println("Result A:", sum)
	
	var m map[int]int
	m = make(map[int]int)
	
	sum = 0
	for i:= 0; i < len(lines); i++ {
		val, err := strconv.Atoi(lines[i])
		check(err)
		sum += val
		
		m[sum]++
		if(m[sum] > 1) {
			fmt.Println("Result B:", sum)
			break
		}
		if i >= len(lines) - 1 {
			i = -1;
		}
	}
	
	
	endtime := makeTimestamp()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}