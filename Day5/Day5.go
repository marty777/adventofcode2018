 package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"strings"
)

func getMillis() int64 {
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

func singlePass(str string) (str2 string, reactions int) {
	reaction_count := 0
	for _, r1 := range "abcdefghijklmnopqrstuvwxyz" {
		r:=string(r1)
		search := r + strings.ToUpper(r)
		reaction_count += strings.Count(str, search)
		str = strings.Replace(str, search, "", -1)
		search = strings.ToUpper(r) + r
		reaction_count += strings.Count(str, search)
		str = strings.Replace(str, search, "", -1)
	}
	return str, reaction_count
}

func removeUnit(str string, r string)string {
	str = strings.Replace(str, r, "", -1)
	str = strings.Replace(str, strings.ToUpper(r), "", -1)
	return str
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// Part A
	reaction_count := 1
	polymer := lines[0]
	for reaction_count > 0 {
		polymer, reaction_count = singlePass(polymer)
	}
	fmt.Println("Result A:",len(polymer))
	
	// Part B
	least_length := len(lines[0])
	for _, r1 := range "abcdefghijklmnopqrstuvwxyz" {
		r:=string(r1)
		polymer = removeUnit(lines[0], r)
		reaction_count = 1
		for reaction_count > 0 {
			polymer, reaction_count = singlePass(polymer)
		}
		if(len(polymer) < least_length) {
			least_length = len(polymer)
		}
	}
	fmt.Println("Result B: ",least_length)
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}