package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"math"
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

type bot struct {
	px,py,pz int
	r int
	//blx,bhx,bly,bhy,blz,bhz int
}

type coordrange struct {
	min_x,min_y,min_z,max_x,max_y,max_z int
}

func abs(a int)int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a int, b int)int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int)int {
	if a > b {
		return a
	}
	return b
}

func intersection(a coordrange, b coordrange)(int, coordrange) {
	overlap := -1
	var ret coordrange
	ret.min_x = max(a.min_x,b.min_x)
	ret.max_x = min(a.max_x,b.max_x)
	ret.min_y = max(a.min_y,b.min_y)
	ret.max_y = min(a.max_y,b.max_y)
	ret.min_z = max(a.min_z,b.min_z)
	ret.max_z = min(a.max_z,b.max_z)
	
	if(ret.min_x > ret.max_x || ret.min_y > ret.max_y || ret.min_z > ret.max_z) {
		return 0, coordrange{min_x:0,max_x:0,min_y:0,max_y:0,min_z:0,max_z:0}
	}
	
	overlap = (ret.max_x - ret.min_x + 1) * (ret.max_y - ret.min_y + 1) * (ret.max_z - ret.min_z + 1)
	return overlap, ret
}



func getMapIntBoolKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func boundingBox(a bot, bounds *[6]int ) {
	if a.px - a.r < bounds[0] {
		bounds[0] = a.px - a.r
	}
	if a.px + a.r > bounds[1] {
		bounds[1] = a.px + a.r
	}
	if a.py - a.r < bounds[2] {
		bounds[2] = a.py - a.r
	}
	if a.py + a.r > bounds[3] {
		bounds[3] = a.py + a.r
	}
	if a.pz - a.r < bounds[4] {
		bounds[4] = a.pz - a.r
	}
	if a.pz + a.r > bounds[5] {
		bounds[5] = a.pz + a.r
	}
}


func getIndex(x int, y int, z int, bounds [6]int )int {
	x_range := bounds[1] - bounds[0]
	//y_range := bounds[3] - bounds[2]
	z_range := bounds[5] - bounds[4]
	index := ((x - bounds[0]) + ((y - bounds[2])*x_range))*z_range + (z - bounds[4])
	return index
}

func dist(x1 int, y1 int, z1 int, x2 int, y2 int, z2 int)int {
	return abs(x2 - x1) + abs(y2 - y1) + abs(z2 - z1)
}

func addCoords(a bot, bounds [6]int, coords map[int]bool) {
	for i:=a.px - a.r; i <= a.px+a.r; i++ {
		for j:= a.py - a.r; j <= a.py + a.r; j++ {
			for k:= a.pz - a.r; k <= a.pz + a.r; k++ {
				if dist(i,j,k,a.px,a.py,a.pz) <= a.r {
					coords[getIndex(i,j,k,bounds)] = true
				}
			}
		}
	}
}

func inRange2(a bot, x int, y int, z int)bool {
	dist := dist(a.px,a.py,a.pz,x,y,z)
	if(dist <= a.r) {
		return true
	}
	return false
}

func inRange(a bot, b bot)bool {
	dist := dist(a.px, a.py, a.pz, b.px, b.py, b.pz)
	if(dist <= a.r) {
		return true
	}
	return false
}


func pointIntersects(a bot, x int, y int, z int)bool {
		dist := abs(x - a.px) + abs(y - a.py) + abs(z - a.pz)
		if(dist <= a.r) {
			return true
		}
		return false
}
func botIntersects(a bot, b bot)bool {
	dist:= abs(a.px - b.px) + abs(a.py - b.py) + abs(a.pz - b.pz)
	if abs(a.r - b.r) >= dist {
		return true
	}
	return false
}

func botToCoordRange(a bot)coordrange {
	crange := coordrange{min_x:a.px - a.r, max_x:a.px + a.r, min_y:a.py - a.r, max_y:a.py + a.r,min_z:a.pz - a.r, max_z:a.pz + a.r }
	return crange
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	
	bots := make([]bot, len(lines))
	for i:= 0; i < len(lines); i++ {
		var px,py,pz int
		var r int
		fmt.Sscanf(lines[i], "pos=<%d,%d,%d>, r=%d", &px, &py, &pz, &r)
		bots[i] = bot{px:px, py:py, pz:pz, r:r}
		fmt.Println(bots[i])
	}
	
	large_range := 0
	large_bot := 0
	for i:= 0; i < len(bots); i++ {
		if bots[i].r > large_range {
			large_range=bots[i].r
			large_bot = i
		}
	}
	
	fmt.Println(large_range, large_bot)
	
	range_count := 0
	for i := 0 ; i < len(bots); i++ {
		if inRange(bots[large_bot], bots[i]) {
			range_count++
		}
	}
	
	fmt.Println("Result A:", range_count)
	
	
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}