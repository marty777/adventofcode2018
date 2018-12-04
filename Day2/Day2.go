 package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
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

func main() {
	starttime := makeTimestamp()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	
	// Part A
	twoCount := 0
	threeCount := 0
	for i:= 0; i < len(lines); i++ {
		line := lines[i]
		var m map[string]int
		m = make(map[string]int)
		
		for j := 0; j < len(line); j++ {
			m[string(line[j])]++
		}
		keys := getMapStringIntKeys(m)
		hasTwo := false
		hasThree := false
		for j := 0; j < len(keys); j++ {
			if m[keys[j]] == 2 {
				hasTwo = true
			}
			if m[keys[j]] == 3 {
				hasThree = true
			}
		}
		if hasTwo {
			twoCount++
		}
		
		if hasThree {
			threeCount++
		}
	}
	
	fmt.Println( "Result A:", twoCount, "x", threeCount, "=" , twoCount * threeCount)
	
	// Part B
	for i:= 0; i < len(lines) - 1; i++ {
		for j := i+1; j < len(lines); j++ {
			diffCount := 0
			for k:=0; k < len(lines[i]); k++ {
				if lines[i][k] != lines[j][k] {
					diffCount++
				}
			}
			if(diffCount == 1) {
				fmt.Print("Result B: ")
				for  k:=0; k < len(lines[i]); k++ {
					if lines[i][k] == lines[j][k]  {
						fmt.Print(string(lines[i][k]))
					}
					
				}
				fmt.Println();
			}
		}
	}
	
	endtime := makeTimestamp()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}