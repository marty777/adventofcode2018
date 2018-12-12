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

func evalState(state []bool, zero_index int)int {
	count:=0
	for i := 0; i < len(state); i++ {
		if(state[i]) {
			count += (i - zero_index)
		}
	}
	return count
}

func readPattern(str string) int {
	pattern := 0
	var i uint
	for i = 0; i < uint(len(str)); i++ {
		if str[i] == '#' {
			pattern |= 1 << i
		}
	}
	return pattern
}

func getPattern(plants []bool, index int)int {
	pattern := 0
	var i uint
	ind := uint(index)
	for i =0; i < 5; i++ {
		if plants[ind + i] {
			pattern |= 1 << i
		}
	}
	return pattern
}

func advanceState(source []bool, dest []bool, rules map[int]bool) {
	for i := 0; i < len(source) - 5; i++ {
		pattern := getPattern(source, i)
		if(rules[pattern]) {
			dest[i+2] = true
		} else {
			dest[i+2] = false
		}
	}
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// parse input
	initial_str := lines[0][len("initial input: "):]
	
	// finite padding bound for the region to evaluate with our cellular automata
	padding:=5000
	state := make([]bool, len(initial_str) +(2*padding))
	state2 := make([]bool, len(initial_str) + (2*padding))
	// read the initial state into the middle of our cells
	for i,_ := range initial_str {
		if initial_str[i] == '.' {
			state[i+padding] = false
		} else if initial_str[i] == '#' {
			state[i+padding] = true
		} 
	}
	
	rules := make(map[int]bool)
	for i:=2; i < len(lines); i++ {
		slice := strings.Split(lines[i], " => ")
		pattern := readPattern(slice[0])
		if slice[1] == "#" {
			rules[pattern] = true
		} else {
			rules[pattern] = false
		}
	}
	
	// Part A
	for i:=0; i<20; i++ {
		if(i %2 == 0) {
			advanceState(state, state2, rules)
		} else {
			advanceState(state2, state, rules)
		}
	}
	
	fmt.Println("Result A:", evalState(state, padding))
	
	// Part B
	// the trick is that the automata will settle into a preditable advancing pattern that 
	// differs by a constant evaluated amount from state to state. Once we advance to that 
	// steady state, the evaluated amount at 50000000000 steps is simple to calculate
	diff:=0
	last_diff := 0
	i:=19
	for {
		last_diff = diff
		if(i %2 == 0) {
			advanceState(state, state2, rules)
			diff = evalState(state2, padding) - evalState(state, padding)
		} else {
			advanceState(state2, state, rules)
			diff = evalState(state, padding) - evalState(state2, padding)
		}
		if(last_diff == diff) {
			break
		}
		i++
	}
	
	term := 50000000000 - (i)
	count := 0
	if i % 2 == 0 {
		count = evalState(state,padding)
	} else {
		count = evalState(state2, padding)
	}
	
	fmt.Println("Result B:", count + (diff * term))
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}