 package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
)

type int_tuple struct {
	x int
	y int
}

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

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	var points = []int_tuple{}
	// parse input
	for i:= 0; i < len(lines); i++ {
		x := 0
		y := 0
		fmt.Sscanf(lines[i], "%d, %d", &x, &y)
		tuple := int_tuple{x:x, y:y}
		points = append(points, tuple)
	}
	
	// Part A
	// go around the edges of a suitable large boundary. Anything that touches the edge can be discounted as infinite. Remaining largest area is the result	
	// Bounds have been guestimated based on sample data. They should really be programatically determined.
	const xmin = -500
	const ymin = -500
	const xmax = 1000
	const ymax = 1000
	
	// 50 points clusted between 0,0 and ~400,400 in sample. Sum of distances < 10000 means not much more than 200 pixels distant from the avg position of the input points, so these bounds
	// are probably overkill.
	const xmin2 = -1000
	const ymin2 = -1000
	const xmax2 = 1000
	const ymax2 = 1000
	
	var grid[xmax - xmin][ymax - ymin]int
	// initalized to 0
	for i := 0; i < xmax-xmin; i++ {
		for j := 0; j < ymax-ymin; j++ {
			grid[i][j] = -1
		}
	}
	
	for i := 0; i < xmax-xmin; i++ {
		for j := 0; j < ymax-ymin; j++ {
			var distances []int
			for k := 0; k < len(points); k++ {
				curr_x := xmin + i
				curr_y := ymin + j
				dist := Abs(points[k].x - curr_x) + Abs(points[k].y - curr_y)
				distances = append(distances, dist)
			}
			min_dist:= (xmax - xmin)*(ymax - ymin)
			for k := 0; k < len(distances); k++ {
				if distances[k] < min_dist  {
					min_dist = distances[k]
				}
			}
			min_dist_count := 0
			min_dist_k := -1
			for k := 0; k < len(distances); k++ {
				if distances[k] == min_dist  {
					min_dist_k = k
					min_dist_count++
				}
			}
			if min_dist_count > 1  {
				grid[i][j] = -1
			} else {
				grid[i][j] = min_dist_k
			}
		}
	}
	
	infinite_areas := make(map[int]bool)
	for i := 0; i < xmax-xmin; i++ {
		if grid[i][0] != -1 {
			infinite_areas[grid[i][0]] = true
		}
		if grid[i][ymax-ymin-1] != -1 {
			infinite_areas[grid[i][ymax-ymin-1]] = true
		}
	}
	for j := 0; j < ymax-ymin; j++ {
		if grid[0][j] != -1 {
			infinite_areas[grid[0][j]] = true
		}
		if grid[xmax-xmin-1][j] != -1 {
			infinite_areas[grid[xmax-xmin-1][j]] = true
		}
	}
	
	counts := make(map[int]int)
	for i := 0; i < xmax-xmin; i++ {
		for j := 0; j < ymax-ymin; j++ {
			if(grid[i][j] == -1  || infinite_areas[grid[i][j]]) {
				continue
			}
			counts[grid[i][j]]++
		}
	}
	
	max_points := 0
	for k := 0; k < len(points); k++ {
		if counts[k] > 0 {
			if counts[k] > max_points {
				max_points = counts[k]
			}
		}
	}
	
	fmt.Println("Result A: ", max_points)
	
	area_count := 0
	for i := 0; i < xmax2-xmin2; i++ {
			for j := 0; j < ymax2-ymin2; j++ {
				dist_sum := 0
				for k := 0; k < len(points); k++ {
					curr_x := xmin2 + i
					curr_y := ymin2 + j
					dist_sum += Abs(points[k].x - curr_x) + Abs(points[k].y - curr_y)
				}
				if(dist_sum < 10000) {
					area_count++
				}
			}		
	}
	
	fmt.Println("Result B:", area_count)

	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}