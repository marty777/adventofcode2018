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
	
	var grid [300][300]int
	// fill grid
	for i := 0; i < 300; i++ {
		for j:=0; j < 300; j++ {
			power := (i + 11)*(j+1)
			power += puzzleID
			power *= (i + 11)			
			power = power / 100
			power = power % 10			
			power -= 5
			grid[i][j] = power
		}
	}
	
	// Part A
	max_power:=0
	x := 0
	y := 0
	for i := 0; i < 300 - 3; i++ {
		for j := 0; j < 300 - 3; j++ {
			power := grid[i][j] + grid[i+1][j] + grid[i+2][j] + grid[i][j+1] + grid[i+1][j+1] + grid[i+2][j+1] + grid[i][j+2] + grid[i+1][j+2] + grid[i+2][j+2] 
			if(power >= max_power) {
				max_power = power
				x = i
				y = j
			}			
		}
	}
	
	fmt.Printf("Result A: %d,%d\n", x+1, y+1)
	
	//Part B
	size := 1
	max_power = 0
	x = 0
	y = 0
	for k := 1; k <= 300; k++ {
		if k == 1 {
			fmt.Println("0% complete...")
		}
		if k == 75 {
			fmt.Println("25% complete...")
		}	
		if k == 150 {
			fmt.Println("50% complete...")
		}
		if k == 225 {
			fmt.Println("75% complete...")
		}			
		for i := 0; i < 300 - k; i++ {
			for j := 0; j < 300 - k; j++ {
				power := 0
				for i1 := 0; i1 < k; i1++ {
					for j1:=0;j1<k;j1++ {
						power += grid[i+i1][j+j1]
					}
				}
				if(power > max_power) {
					max_power = power
					x = i
					y = j
					size = k
				}			
			}
		}
	}
	fmt.Printf("Result B: %d,%d,%d\n", x+1,y+1, size)
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}