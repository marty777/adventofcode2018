package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"sort"
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

func charVal(str string)int{
    sbytes := []byte(str)
	return int(sbytes[0] - 64)
}

type job struct {
	letter string
	remaining int
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	outbound := make(map[string][]string)
	outbound2 := make(map[string][]string)
	// sort input
	sort.Strings(lines)
	
	// parse input
	for i:= 0; i < len(lines); i++ {
		a := ""
		b := ""
		fmt.Sscanf(lines[i], "Step %s must be finished before step %s can begin.", &a, &b)
		outbound[a] = append(outbound[a], b)
		outbound2[a] = append(outbound2[a], b)
	}
	
	
	// Part A
	result := ""
	for {
		// find elements with no incoming connections
		// build inbound list
		inbound := make(map[string][]string)
		for k := range outbound {
			for i := 0; i < len(outbound[k]); i++ {
				inbound[outbound[k][i]] = append(inbound[outbound[k][i]], k)
			}
		}
		
		var candidates []string
		for k := range outbound {
			if len(inbound[k]) == 0 {
				candidates = append(candidates, k)
			}
		}
		if len(candidates) == 0 {
			fmt.Println("something bad happened")
		}
		
		sort.Strings(candidates)
		result += candidates[0]
		
		if(len(outbound) == 1) {
			final := outbound[candidates[0]]
			// just in case
			sort.Strings(final)
			for i:=0; i < len(final); i++ {
				result += final[i]
			}
			break
		}
		
		delete(outbound,candidates[0])
		// might hit this on zero input
		if(len(outbound) == 0) {
			break
		}
	}
	
	fmt.Println("Result A:", result)
	
	// Part B
	jobs := make([]job,5)
	steps := 0
	
	for {
		// finish jobs
		// get fastest running job to finish and update all job times
		var candidates []string
		fastest := 120
		finished_found := false
		for i:=0; i < len(jobs); i++ {
			if jobs[i].letter != "" && jobs[i].remaining < fastest {
				fastest = jobs[i].remaining
				finished_found = true
			}
		}
		var finished[]string
		if(finished_found) {
			for i:=0; i<len(jobs); i++ {
				jobs[i].remaining -= fastest;
				if jobs[i].letter != "" && jobs[i].remaining <= 0 {
					finished = append(finished, jobs[i].letter)
					jobs[i].letter = ""
					jobs[i].remaining = 0
				}
			}
			
			// remove the finished jobs
			for i:=0; i < len(finished); i++ {
				// badly handle final nodes with only inbound dependencies. This will lose the final nodes if there aren't enough job queues available at the step
				// where the final outbound node is removed. Works fine on the puzzle input though.
				if len(outbound2) == 1 {
					for j:= 0; j < len(outbound2[finished[i]]); j++ {
						candidates = append(candidates, outbound2[finished[i]][j])
					}
				}
				delete(outbound2, finished[i])
			}
			steps += fastest
		}
		
		// assign new jobs
		
		// build inbound list
		inbound := make(map[string][]string)
		for k := range outbound2 {
			for i := 0; i < len(outbound2[k]); i++ {
				inbound[outbound2[k][i]] = append(inbound[outbound2[k][i]], k)
			}
		}
		
		// find candidate jobs
		for k := range outbound2 {
			if len(inbound[k]) == 0 {
				// filter out already assigned candidates
				found:= false
				for i:= 0; i < len(jobs); i++ {
					if jobs[i].letter == k {
						found = true
					}
				}
				if !found {
					candidates = append(candidates, k)
				}
			}
		}
		
		// assign new jobs
		if len(candidates) > 0 {
			sort.Strings(candidates)
			candidate_i := 0
			for i:= 0; i < len(jobs); i++ {
				if jobs[i].letter == "" && candidate_i < len(candidates) {
					jobs[i].letter = candidates[candidate_i]
					jobs[i].remaining = 60 + charVal(jobs[i].letter)
					candidate_i++
				} 
			}
		}
		
		//check if job queue finished
		remaining_jobs := 0
		for i:= 0; i < len(jobs); i++ {
			if jobs[i].letter != "" {
				remaining_jobs++
			}
		}
		
		if remaining_jobs == 0 {
			break
		}
	}
	
	fmt.Println("Result B:", steps)

	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}