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
	clay bool
	wet bool
	dripping bool
}

func inputBounds(lines []string)(int,int,int,int) {
	min_x := -1
	min_y := 0
	max_x := 0
	max_y := 0
	for i:=0; i < len(lines); i++ {
		var v1, v2, v3 int
		if(lines[i][0] == 'x') {
			fmt.Sscanf(lines[i], "x=%d, y=%d..%d", &v1, &v2, &v3)
			if(v1 < min_x || min_x == -1) {
				min_x = v1
			}
			if(v1 > max_x || max_x == -1) {
				max_x = v1
			}
			
			if(v2 < min_y || min_y == -1) {
				min_y = v2
			}
			if(v2 > max_y || max_y == -1) {
				max_y = v2
			}
			if(v3 < min_y || min_y == -1) {
				min_y = v3
			}
			if(v3 > max_y || max_y == -1) {
				max_y = v3
			}
			
		} else {
			fmt.Sscanf(lines[i], "y=%d, x=%d..%d", &v1, &v2, &v3)
			if(v1 < min_y || min_y == -1) {
				min_y = v1
			}
			if(v1 > max_y || max_y == -1) {
				max_y = v1
			}
			
			if(v2 < min_x || min_x == -1) {
				min_x = v2
			}
			if(v2 > max_x || max_x == -1) {
				max_x = v2
			}
			if(v3 < min_x || min_x == -1) {
				min_x = v3
			}
			if(v3 > max_x || max_x == -1) {
				max_x = v3
			}
		}
	}
	
	return min_x, min_y, max_x, max_y
}

func parseInput(lines []string, grid [][]gridSquare, x_offset int, y_offset int) {
	for i:= 0; i < len(lines); i++ {
		var v1, v2, v3 int
		
		if(lines[i][0] == 'x') {
			fmt.Sscanf(lines[i], "x=%d, y=%d..%d", &v1, &v2, &v3)
			for j:=v2-y_offset; j <= v3 - y_offset; j++ {
				grid[v1 - x_offset][j].clay = true
			}
		} else {
			fmt.Sscanf(lines[i], "y=%d, x=%d..%d", &v1, &v2, &v3)
			for j:=v2-x_offset; j <= v3 - x_offset; j++ {
				grid[j][v1-y_offset].clay = true
			}
		}
	}
}

func printState(grid [][]gridSquare) {
	for j:=0;j<len(grid[0]);j++ {
		for i:=0;i<len(grid);i++ {
			if grid[i][j].clay {
				fmt.Print("#")
			} else if grid[i][j].wet {
				fmt.Print("~")
			} else if grid[i][j].dripping {
				fmt.Print("|")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func advance(grid [][]gridSquare, start_depth int, actual_depth int) (int,int) {
	// spread into containers
	waterSquares := 0
	stillSquares := 0
	for j:=0;j<len(grid[0])-1;j++ {
		for i:=0; i < len(grid)-1; i++ {
			if grid[i][j].dripping {
				if(grid[i][j+1].clay || grid[i][j+1].wet) {
					// get a layer wet
					container_min_x := i
					container_max_x := i
					falling_layer := false
					for k:=i; k > 0; k-- {
						if grid[k-1][j].clay {
							container_min_x = k
							break
						} else if grid[k-1][j+1].dripping { 
							container_min_x = -1
							break
						} else if !grid[k-1][j+1].clay && !grid[k-1][j+1].wet {
							container_min_x = k
							grid[k-1][j].dripping = true
							falling_layer = true
							break
						}
					}
					for k:=i; k < len(grid) -1; k++ {
						if grid[k+1][j].clay {
							container_max_x = k
							break
						} else if grid[k+1][j+1].dripping { 
							container_max_x = -1
							break
						} else if(!grid[k+1][j+1].clay && !grid[k+1][j+1].wet) {
							container_max_x = k
							grid[k+1][j].dripping = true
							falling_layer = true
							break
						}
					}
					
					if(container_min_x == -1 && container_max_x == -1) {
						continue
					}
					if container_min_x == -1 {
						falling_layer = true
						container_min_x = i
					}
					if container_max_x == -1 {
						falling_layer = true
						container_max_x = i
					}					
					
					
					for k:=container_min_x; k <= container_max_x; k++ {
						if(falling_layer) {
							grid[k][j].dripping = true
						} else {
							grid[k][j].dripping = false
							grid[k][j].wet = true
						}
					}
				} else {
					// fall
					grid[i][j+1].dripping = true
				}
			}
		}
		
		if(j >= start_depth && j <= actual_depth) { 
			for i:=0; i < len(grid)-1; i++ {
				if grid[i][j].wet || grid[i][j].dripping {
					waterSquares++
				}
				if grid[i][j].wet {
					stillSquares++
				}
				
			}
		}
	}
	return waterSquares, stillSquares
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	min_x, min_y, max_x, max_y := inputBounds(lines)
	min_x -= 5
	max_x += 5
	max_y++
	
	grid := make([][]gridSquare, max_x - min_x + 1) 
	for i:=0; i < max_x - min_x + 1; i++ {
		grid[i] = make([]gridSquare, max_y - min_y + 1)
	}
	parseInput(lines, grid, min_x, min_y)
	grid[500 - min_x][0].dripping = true
	lastCount := 0
	currCount := 0
	stillCount := 0
	steadyCount := 0
	for steadyCount < 10 {
		lastCount = currCount
		currCount, stillCount = advance(grid,min_y, max_y)
		if(currCount == lastCount) {
			steadyCount++
		} else {
			steadyCount = 0
		}
	}
	
	//printState(grid)
	fmt.Println("Result A:", currCount - 1)
	fmt.Println("Result B:", stillCount)
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}