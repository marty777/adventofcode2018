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

type testelement struct {
	start [4]int
	inst [4]int
	after [4]int
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

func copy4(dest *[4]int, src [4]int) {
	for i:=0;i<4;i++ {
		dest[i] = src[i]
	}
}

func cmp4(a [4]int, b [4]int) bool {
	for i:=0; i<4; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func addr(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = input[a.inst[1]] + input[a.inst[2]]
	return cmp4(input,a.after)
}

func addi(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = input[a.inst[1]] + a.inst[2]
	return cmp4(input,a.after)
}

func mulr(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = input[a.inst[1]] * input[a.inst[2]]
	return cmp4(input,a.after)
}

func muli(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = input[a.inst[1]] * a.inst[2]
	return cmp4(input,a.after)
}

func banr(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = input[a.inst[1]] & input[a.inst[2]]
	return cmp4(input,a.after)
}

func bani(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = input[a.inst[1]] & a.inst[2]
	return cmp4(input,a.after)
}

func borr(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = input[a.inst[1]] | input[a.inst[2]]
	return cmp4(input,a.after)
}

func bori(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = input[a.inst[1]] | a.inst[2]
	//fmt.Println(a.start, input, a.after, input[a.inst[1]] | a.inst[2], input[a.inst[1]], a.inst[2])
	return cmp4(input,a.after)
}

func setr(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = input[a.inst[1]]
	return cmp4(input,a.after)
}

func seti(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	input[a.inst[3]] = a.inst[1]
	return cmp4(input,a.after)
}

func gtir(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	if a.inst[1] > input[a.inst[2]] {
		input[a.inst[3]] = 1
	} else {
		input[a.inst[3]] = 0
	}
	return cmp4(input,a.after)
}

func gtri(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	if input[a.inst[1]] > a.inst[2] {
		input[a.inst[3]] = 1
	} else {
		input[a.inst[3]] = 0
	}
	return cmp4(input,a.after)
}

func gtrr(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	if input[a.inst[1]] > input[a.inst[2]] {
		input[a.inst[3]] = 1
	} else {
		input[a.inst[3]] = 0
	}
	return cmp4(input,a.after)
}

func eqir(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	if a.inst[1] == input[a.inst[2]] {
		input[a.inst[3]] = 1
	} else {
		input[a.inst[3]] = 0
	}
	return cmp4(input,a.after)
}

func eqri(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	if input[a.inst[1]] == a.inst[2] {
		input[a.inst[3]] = 1
	} else {
		input[a.inst[3]] = 0
	}
	return cmp4(input,a.after)
}

func eqrr(a testelement)bool {
	var input [4]int
	copy4(&input,a.start)
	if input[a.inst[1]] == input[a.inst[2]] {
		input[a.inst[3]] = 1
	} else {
		input[a.inst[3]] = 0
	}
	return cmp4(input,a.after)
}

func mask(power int)int {
	out := 65535
	a := 1 << uint(power)
	return out ^ a
}


func partA(elements []testelement)int {
	
	partAcount:=0
	for i:= 0; i < len(elements); i++ {
		samplecount := 0
		// addr
		if addr(elements[i]) {
			samplecount++
		}
		if addi(elements[i]) {
			samplecount++
		}
		if mulr(elements[i]) {
			samplecount++
		}
		if muli(elements[i]) {
			samplecount++
		}
		if banr(elements[i]) {
			samplecount++
		}
		if bani(elements[i]) {
			samplecount++
		}
		if borr(elements[i]) {
			samplecount++
		}
		if bori(elements[i]) {
			samplecount++
		}
		if setr(elements[i]) {
			samplecount++
		}
		if seti(elements[i]) {
			samplecount++
		}
		if gtir(elements[i]) {
			samplecount++
		}
		if gtri(elements[i]) {
			samplecount++
		}
		if gtrr(elements[i]) {
			samplecount++
		}
		if eqir(elements[i]) {
			samplecount++
		}
		if eqri(elements[i]) {
			samplecount++
		}
		if eqrr(elements[i]) {
			samplecount++
		}
		
		if(samplecount >= 3) {
			partAcount++
		}
	}
	
	return partAcount
}

func countbits(a int) int {
	bits := a
	count := 0
	for bits > 0 {
		if bits % 2 == 1 {
			count++
		}
		bits = bits >> 1
	}
	return count
}

//the index of the highest set bit
func highbitindex(a int) int {
	index := 0
	for a > 0 {
		a = a >> 1
		index++
	}
	return index-1
}

func bitindexoff(a int, index int) int {
	mask := 1 << uint(index)
	mask = ^mask
	return a & mask
}

func execute(inst instruction, state *[4]int) {
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

func partB(elements []testelement, program []instruction)int {
	results := make(map[int]int)
	for i := 0; i < 16; i++  {
		results[i] = 65535 // first 16 bits on
	}
	
	// remove all cases an opcode fails
	for i:= 0; i < len(elements); i++ {
		opcode := elements[i].inst[0]
		if !addr(elements[i]) {
			results[opcode] &= mask(0)
		}
		if !addi(elements[i]) {
			results[opcode] &= mask(1)
		}
		if !mulr(elements[i]) {
			results[opcode] &= mask(2)
		}
		if !muli(elements[i]) {
			results[opcode] &= mask(3)
		}
		if !banr(elements[i]) {
			results[opcode] &= mask(4)
		}
		if !bani(elements[i]) {
			results[opcode] &= mask(5)
		}
		if !borr(elements[i]) {
			results[opcode] &= mask(6)
		}
		if !bori(elements[i]) {
			results[opcode] &= mask(7)
		}
		if !setr(elements[i]) {
			results[opcode] &= mask(8)
		}
		if !seti(elements[i]) {
			results[opcode] &= mask(9)
		}
		if !gtir(elements[i]) {
			results[opcode] &= mask(10)
		}
		if !gtri(elements[i]) {
			results[opcode] &= mask(11)
		}
		if !gtrr(elements[i]) {
			results[opcode] &= mask(12)
		}
		if !eqir(elements[i]) {
			results[opcode] &= mask(13)
		}
		if !eqri(elements[i]) {
			results[opcode] &= mask(14)
		}
		if !eqrr(elements[i]) {
			results[opcode] &= mask(15)
		}
		
	}
	
	opcodemap := make(map[int]int)
	
	// narrow down which opcode is witch by ruling out certainties (ie, opcodes that only have one possible option)
	for {
		// find element with 1 bit
		min_bits := 16
		min_op := -1
		for i:=0;i<16;i++  {
			bits := countbits(results[i])
			if bits == 1 {
				min_bits = 1
				min_op = i
				break
			}
		}
		if min_bits != 1 {
			// we're either finished or in trouble
			break
		}
		bitindex := highbitindex(results[min_op])
		opcodemap[min_op] = bitindex
		for i:=0;i<16;i++  {
			results[i] = bitindexoff(results[i],bitindex)
		}
		
	}
	
	var state [4]int // initializes to zeroes
	for i:=0; i < len(program); i++ {
		statement:=program[i]
		statement.inst[0] = opcodemap[statement.inst[0]]
		execute(statement,&state)
	}
	return state[0]
}


func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	// parse input
	var tests []testelement
	var program []instruction
	
	i:=0
	for i<len(lines) {
		var a1,b1,c1,d1 int
		var a2,b2,c2,d2 int
		var a3,b3,c3,d3 int
		if len(lines[i+0]) == 0 && len(lines[i+1]) == 0 {
			break
		}
		fmt.Sscanf(lines[0+i], "Before: [%d, %d, %d, %d]", &a1, &b1, &c1, &d1)
		fmt.Sscanf(lines[1+i], "%d %d %d %d]", &a2, &b2, &c2, &d2)
		fmt.Sscanf(lines[2+i], "After: [%d, %d, %d, %d]", &a3, &b3, &c3, &d3)
		var element testelement
		element.start[0] = a1
		element.start[1] = b1
		element.start[2] = c1
		element.start[3] = d1
		element.inst[0] = a2
		element.inst[1] = b2
		element.inst[2] = c2
		element.inst[3] = d2
		element.after[0] = a3
		element.after[1] = b3
		element.after[2] = c3
		element.after[3] = d3
		tests = append(tests, element)
		i+=4
	}
	
	for i < len(lines) {
		if len(lines[i]) == 0 {
			i++
			continue
		}
		var a,b,c,d int
		fmt.Sscanf(lines[i], "%d %d %d %d", &a,&b,&c,&d)
		var inst instruction
		inst.inst[0] = a
		inst.inst[1] = b
		inst.inst[2] = c
		inst.inst[3] = d
		program = append(program, inst)
		i++
	}
	
	
	resultA := partA(tests)
	fmt.Println("Result A:", resultA)
	resultB := partB(tests, program)
	fmt.Println("Result B:", resultB)
	
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}