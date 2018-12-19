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
	inst [4]uint64
}

func parseInput(lines[] string)(int, []instruction) {
	var ip int
	fmt.Sscanf(lines[0], "#ip %d", &ip)
	var program []instruction
	var instlist = []string{"addr","addi","mulr","muli","banr","bani","borr","bori","setr","seti","gtir","gtri","gtrr","eqir","eqri","eqrr"}
	fmt.Println(instlist)
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
		program = append(program, instruction{inst: [4]uint64{uint64(inst_i), uint64(a), uint64(b), uint64(c)}})
	}
	return ip, program
}

func execute(inst instruction, state *[6]uint64) {
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
	var state [6]uint64
	
	var index uint64
	fmt.Println(state)
	for {
		state[ip] = index
		execute(program[index], &state)
		index = state[ip]
		index++
		if(index >= uint64(len(program)) || index < 0) {
			break
		}
	}
	fmt.Println("Result A:", state[0]);
	
	state = [6]uint64{1,0,0,0,0,0}
	fmt.Println(state)
	index = 0
	for {
		//fmt.Println(state)
		state[ip] = index		
		execute(program[index], &state)
		index = state[ip]
		index++
		if(index >= uint64(len(program)) || index < 0) {
			break
		}
	}
	fmt.Println("Result B:", state[0]);
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}