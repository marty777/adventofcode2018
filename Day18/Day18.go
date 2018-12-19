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

type gridSquare struct {
	state int
}

func printState(grid [][]gridSquare) {
	for j:=0;j<len(grid[0]);j++ {
		for i:=0;i<len(grid);i++ {
			if grid[i][j].state == 0{
				fmt.Print(".")
			} else if grid[i][j].state == 1 {
				fmt.Print("|")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

func parseInput(lines []string, grid [][]gridSquare) {
	for j:=0; j < len(lines); j++ {
		for i:=0;i<len(lines[j]);i++ {
			if(lines[j][i] == '.') {
				grid[i][j].state = 0
			} else if(lines[j][i] == '|') {
				grid[i][j].state = 1
			} else {
				grid[i][j].state = 2
			}
		}
	}
}

func neighbortest(grid [][]gridSquare, x int, y int, teststate int, min int) bool {
	x1:=0
	y1:=0
	count := 0
	for i:= 0; i < 8; i++ {
		switch i {
			case 0:
				x1 = x - 1
				y1 = y - 1
			case 1:
				x1 = x
				y1 = y - 1
			case 2:
				x1 = x + 1
				y1 = y - 1
			case 3:
				x1 = x + 1
				y1 = y
			case 4:
				x1 = x + 1
				y1 = y + 1
			case 5:
				x1 = x 
				y1 = y + 1
			case 6:
				x1 = x - 1
				y1 = y + 1
			case 7:
				x1 = x - 1
				y1 = y
		}
		
		if x1 >= 0 && x1 < len(grid) && y1 >= 0 && y1 < len(grid[0]) {
			if grid[x1][y1].state == teststate {
				count++
			}
		}
	}
	if(count >= min) {
		return true
	}
	return false
}

func advance(src [][]gridSquare, dest [][]gridSquare) {
	for i:=0; i < len(src); i++ {
		for j:=0; j < len(src[0]); j++ {
			dest[i][j] = src[i][j]
			if src[i][j].state == 0 && neighbortest(src, i, j, 1, 3) {
				dest[i][j].state = 1
			} else if src[i][j].state == 1 && neighbortest(src, i, j, 2, 3) {
				dest[i][j].state = 2
			} else if src[i][j].state == 2 {
				if neighbortest(src, i, j, 2, 1) && neighbortest(src, i, j, 1, 1) {
					dest[i][j].state = 2
				} else {
					dest[i][j].state = 0
				}
			}
		}
	}
}

func value(src [][]gridSquare)(int,int) {
	woodCount:=0
	lumberCount:=0
	for i:=0; i < len(src); i++ {
		for j:=0; j < len(src[0]); j++ {
			if src[i][j].state == 1 {
				woodCount++
			} else if src[i][j].state == 2 {
				lumberCount++
			}
		}
	}
	return woodCount, lumberCount
}

// returns the period or -1
func period(trace []int)int {
	// we're looking for 4 points with identical values that are equidistant from each other in the trace
	// If those are found, check for several previous matching points with the same period
		
		
	
	framelen := 1
	// we want to see 4 of the pattern at equal distances
	pattern_count := 0;
	index1 := len(trace) - framelen
	index2 := -1
	index3 := -1
	index4 := -1
	found4 := false
	for i:=index1-framelen; i >= 0; i-- {
		found := true
		for k:=0; k < framelen; k++ {
			if trace[index1+k] != trace[i+k] {
				found = false
				break
			}
		}
		if(found) {
			pattern_count++
			if pattern_count == 1 {
				index2 = i
			} else if pattern_count == 2 {
				index3 = i
			} else if pattern_count == 3 {
				index4 = i
			} else {
				found4 = true
				break
			}
		}
	}
	if(found4 && index1-index2 == index2 - index3 && index2-index3 == index3 - index4) {
		// check for previous matches with same period
		still_matching := true
		checkback := 4
		if(index4-checkback < 0) {
			return -1
		}
		for i:=1; i <= checkback; i++  {
			if trace[index1 - i] != trace[index2 - i] || trace[index2 - i] != trace[index3 - i] || trace[index3 - i] != trace[index4-i] {
				still_matching = false
				break
			}
		}
		if(still_matching) {
			return index2-index3
		}
	}
	return -1
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	grid1 := make([][]gridSquare, len(lines[0]))
	grid2 := make([][]gridSquare, len(lines[0]))
	for i := 0; i < len(lines[0]); i++ {
		grid1[i] = make([]gridSquare, len(lines))
		grid2[i] = make([]gridSquare, len(lines))
	}
	
	parseInput(lines, grid1)
	
	var trace []int
	//printState(grid1)
	for i := 0; i < 10; i++ {
		a := 0
		b := 0
		if i % 2 == 0 {
			advance(grid1, grid2)
			a,b = value(grid2)
		} else {
			advance(grid2, grid1)
			a,b  = value(grid1)
		}
		trace = append(trace, a*b)
	}	
	woodCount, lumberCount := value(grid1)
	fmt.Printf("Result A: %d x %d = %d\n", woodCount, lumberCount, woodCount * lumberCount)
	for i := 10; i < 1000000000; i++ {
		a:=0
		b:=0
		if i % 2 == 0 {
			advance(grid1, grid2)
			a,b = value(grid2)
			
		} else {
			advance(grid2, grid1)
			a,b  = value(grid1)
		}
		trace = append(trace, a*b)
		period := period(trace)
		if(period > 1) {
			hops_back := 1 + ((1000000000 - i)/period)
			index := 1000000000 - (hops_back * period)
			fmt.Println("Result B:", trace[index-1])
			break
		} else if  period == 1 {
			fmt.Println("Result B:", trace[len(trace) - 1])
			break
		}
	}	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}