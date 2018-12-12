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

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// parse input
	puzzleID := 0
	fmt.Sscanf(lines[0], "%d", &puzzleID)
	
	var grids[300][300][300]int
	highscore3 := 0
	highscore := 0
	size:=0
	x3:=0
	y3:=0
	x:=0
	y:=0
	// fill grid
	for i := 0; i < 300; i++ {
			for j:=0; j < 300; j++ {
				power := (i + 11)*(j+1)
				power += puzzleID
				power *= (i + 11)			
				power = power / 100
				power = power % 10			
				power -= 5
				grids[0][i][j] = power
				if(power > highscore) {
					highscore = power
					x = i+1;
					y = j+1
				}
		}
	}
	for k:= 1; k < 300; k++ {
		for i:= 0; i < 300- k ; i++ {
			for j:=0; j < 300 - k; j++ {
				power := grids[k-1][i][j]
				//lower edge
				for i1 := 0; i1 < k+1; i1++ {
					power += grids[0][i+i1][j+k]
				}
				// right edge
				for j1:= 0; j1 < k; j1++ {
					power += grids[0][i+k][j+j1]
				}
				grids[k][i][j] = power
				if(power > highscore) {
					highscore = power
					x = i+1
					y = j+1
					size = k+1
				}
				if(k == 2 && power > highscore3) {
					highscore3 = power
					x3 = i+1
					y3 = j+1
				}
			}
		}
		if(k == 2) {
			fmt.Printf("Result A: %d,%d\n", x3, y3)
		}
	}
	
	fmt.Printf("Result B: %d,%d,%d\n", x,y, size)
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}