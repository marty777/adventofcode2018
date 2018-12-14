package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"strconv"
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

type orientation int

const (
	north			orientation = 0
	east			orientation = 1
	south			orientation = 2
	west			orientation = 3
)

type trackgrid byte 
const (
	none 			trackgrid = 0
	horizontal 		trackgrid = 1
	vertical 		trackgrid = 2
	intersection 	trackgrid = 3
	cornerslash		trackgrid = 4
	cornerbackslash	trackgrid = 5
)

func digits(a int) []int {
	var out []int
	s:=strconv.Itoa(a)
	for _,c := range s {
		d,_ := strconv.Atoi(string(c))
		out = append(out, d)
	}
	return out
}

func checkBoard(input []int, board []int) int {
	
	
	if len(board) - len(input) < 0 {
		return 0
	}
	
	match := true
	for i := 0; i < len(input); i++ {
		if board[i + len(board) - len(input)] != input[i] {
			match = false
			break
		}
	}
	
	if match {
		return len(board) - len(input)
	}
	return 0
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// parse input
	input := 0
	fmt.Sscanf(lines[0], "%d", &input)
	inputdigits := digits(input)
	
	var board []int
	elf1:=0
	elf2:=1
	board = append(board,3)
	board = append(board,7)
	partAComplete := false
	partBComplete := false
	for {
		sum:=board[elf1] + board[elf2]
		digits := digits(sum)
		board = append(board, digits...)
		elf1 = (board[elf1] + 1 + elf1) % len(board)
		elf2 = (board[elf2] + 1 + elf2) % len(board)
		if len(board) > input + 10 && !partAComplete{
			fmt.Print("Result A:")
			for i := input; i < input + 10; i++ {
				print(board[i])
			}
			fmt.Println();
			partAComplete = true
		}
		if(!partBComplete) {
			index := checkBoard(inputdigits, board)
			if index > 0 {
				fmt.Println("Result B:", index)
				partBComplete = true
			}
		}
		if partAComplete && partBComplete {
			break
		}
		
		if(len(board) % 10000000 == 0) {
			fmt.Println("Board size:", len(board))
		}
	}
	
	
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}