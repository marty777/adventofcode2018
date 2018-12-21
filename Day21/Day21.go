package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
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

type instruction struct {
	inst [4]int
}


func getMapIntIntKeys(m map[int]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func parseInput(lines[] string)(int, []instruction) {
	var ip int
	fmt.Sscanf(lines[0], "#ip %d", &ip)
	var program []instruction
	var instlist = []string{"addr","addi","mulr","muli","banr","bani","borr","bori","setr","seti","gtir","gtri","gtrr","eqir","eqri","eqrr"}
	for i:=1; i<len(lines);i++ {
	    var inst string
		var a,b,c int
		fmt.Sscanf(lines[i], "%s %d %d %d",&inst, &a, &b, &c)
		inst_i := -1
		for j:=0; j<len(instlist);j++ {
			if inst == instlist[j] {
				inst_i = j
				break
			}
		}
		program = append(program, instruction{inst: [4]int{int(inst_i), int(a), int(b), int(c)}})
	}
	return ip, program
}

func execute(inst instruction, state *[6]int) {
	switch inst.inst[0] {
		case 0: // addr
			state[inst.inst[3]] = state[inst.inst[1]] + state[inst.inst[2]]
		case 1: // addi
			state[inst.inst[3]] = state[inst.inst[1]] + inst.inst[2]
		case 2: // mulr
			state[inst.inst[3]] = state[inst.inst[1]] * state[inst.inst[2]]
		case 3: // muli
			state[inst.inst[3]] = state[inst.inst[1]] * inst.inst[2]
		case 4: // banr
			state[inst.inst[3]] = state[inst.inst[1]] & state[inst.inst[2]]
		case 5: // bani
			state[inst.inst[3]] = state[inst.inst[1]] & inst.inst[2]
		case 6: // borr
			state[inst.inst[3]] = state[inst.inst[1]] | state[inst.inst[2]]
		case 7: // bori
			state[inst.inst[3]] = state[inst.inst[1]] | inst.inst[2]
		case 8: // setr
			state[inst.inst[3]] = state[inst.inst[1]]
		case 9: // seti
			state[inst.inst[3]] = inst.inst[1]
		case 10: // gtir
			if inst.inst[1] > state[inst.inst[2]] {
				state[inst.inst[3]] = 1
			} else {
				state[inst.inst[3]] = 0
			}
		case 11: //gtri
			if state[inst.inst[1]] > inst.inst[2] {
				state[inst.inst[3]] = 1
			} else {
				state[inst.inst[3]] = 0
			}
		case 12: //gtrr
			if state[inst.inst[1]] > state[inst.inst[2]] {
				state[inst.inst[3]] = 1
			} else {
				state[inst.inst[3]] = 0
			}
		case 13: // eqir
			if inst.inst[1] == state[inst.inst[2]] {
				state[inst.inst[3]] = 1
			} else {
				state[inst.inst[3]] = 0
			}
		case 14: //eqri
			if state[inst.inst[1]] == inst.inst[2] {
				state[inst.inst[3]] = 1
			} else {
				state[inst.inst[3]] = 0
			}
		case 15: //eqrr
			if state[inst.inst[1]] == state[inst.inst[2]] {
				state[inst.inst[3]] = 1
			} else {
				state[inst.inst[3]] = 0
			}
	}
	
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	ip, program:=parseInput(lines)
	var state [6]int
	state[0] = 0
	var index int
	steps:=0
	for {
		
		state[ip] = index
		
		execute(program[index], &state)
		if state[ip] == 28 {
			break;
		}
		index = state[ip]
		index++
		steps++
		if(index >= int(len(program)) || index < 0) {
			break
		}
	}
	resultA := 0
	testReg := 0
	if program[28].inst[1] == 0 {
		resultA = state[program[28].inst[2]]
		testReg = program[28].inst[2]
	} else {
		resultA = state[program[28].inst[1]]
		testReg = program[28].inst[1]
	}
	fmt.Println(program[28])
	fmt.Println(state)
	fmt.Println("Result A:", resultA);
	
	// Part B - Odd wording for the solution, but this is the final value we hit on the comparison register to r0 at the end
	// of the program before we see a repeat value.
	state = [6]int{0,0,0,0,0,0}
	index = 0
	steps = 0
	valmap := make(map[int]int)
	last_val := 0
	for {
		
		state[ip] = index
		execute(program[index], &state)
		if(state[ip] == 28) {
			valmap[state[testReg]]++
			if(valmap[state[testReg]] > 1) {
				break;
			}
			last_val = state[testReg]
		}
		index = state[ip]+1
		steps++
	}
	
	fmt.Println("Result B:", last_val)
	
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}