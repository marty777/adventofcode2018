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


type bot struct {
	px,py,pz int
	r int
}

type coord struct {
	x,y,z int
}

func abs(a int)int {
	if a < 0 {
		return -a
	}
	return a
}

func dist(x1 int, y1 int, z1 int, x2 int, y2 int, z2 int)int {
	return abs(x2 - x1) + abs(y2 - y1) + abs(z2 - z1)
}

func inRange2(a bot, c coord)bool {
	dist := dist(a.px, a.py, a.pz, c.x, c.y, c.z)
	if(dist <= a.r) {
		return true
	}
	return false
}

func countInRange(bots []bot, c coord) int {
	count := 0;
	for i := 0; i < len(bots); i++ {
		if inRange2(bots[i], c) {
			count++
		}
	}
	return count
}

func inRange(a bot, b bot)bool {
	dist := dist(a.px, a.py, a.pz, b.px, b.py, b.pz)
	if(dist <= a.r) {
		return true
	}
	return false
}


func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// Parse input
	bots := make([]bot, len(lines))
	for i:= 0; i < len(lines); i++ {
		var px,py,pz int
		var r int
		fmt.Sscanf(lines[i], "pos=<%d,%d,%d>, r=%d", &px, &py, &pz, &r)
		bots[i] = bot{px:px, py:py, pz:pz, r:r}
	}
	
	// Part A
	large_range := 0
	large_bot := 0
	for i:= 0; i < len(bots); i++ {
		if bots[i].r > large_range {
			large_range=bots[i].r
			large_bot = i
		}
	}
	
	range_count := 0
	for i := 0 ; i < len(bots); i++ {
		if inRange(bots[large_bot], bots[i]) {
			range_count++
		}
	}
	
	fmt.Println("Result A:", range_count)
	
	// Part B
	// Supposition: Optimal points are found at or near on of the corners of the octahedron bounding one of the bots
	var corners []coord
	points := [][]int{{-1,0,0},{1,0,0},{0,-1,0},{0,1,0},{0,0,-1},{0,0,1}}
	for i:= 0; i < len(bots); i++ {
		for j:= 0; j < 6; j++ {
			corners = append(corners, coord{x:bots[i].px + (bots[i].r * points[j][0]),y:bots[i].py + (bots[i].r * points[j][1]),z:bots[i].pz + (bots[i].r * points[j][2])})
		}
	}
	
	
	// search a 17x17x17 volume around each corner
	var best_coord coord
	best_count := 0
	min_dist := 1 << 32
	
	for i:=0; i < len(corners); i++ {
		c := corners[i]
		count := 0
		for  j:=0; j<len(bots); j++ {
			if inRange2(bots[j], c) {
				count++
			}
		}
		if count > best_count || (count == best_count && dist(c.x,c.y,c.z,0,0,0) < min_dist) {
			best_count = count
			min_dist = dist(c.x,c.y,c.z,0,0,0)
			best_coord = c
		}
	}
	
	// hunt around for neighboring improvements until none found
	best_coord = coord{x:27240491,y:44370529, z:54887618}
	best_count = 882
	min_dist = 126498638
	bound := 8
	updatecount :=0
	for {
		updated := false
		c := best_coord
		for j:=c.x - bound; j <= c.x + bound; j++ {
			for k:=c.y - bound; k <= c.y + bound; k++ {
				for m:=c.z - bound; m <= c.z + bound; m++ {
					p := coord{j,k,m}
					count := 0
					for n:= 0; n < len(bots); n++ {
						if inRange2(bots[n], p) {
							count++
						}
					}
					dist := dist(0,0,0,p.x,p.y,p.z)
					if count > best_count || (count == best_count && dist < min_dist) {
						best_coord = p
						min_dist = dist
						best_count = count
						updated = true
						updatecount++
						//if updatecount % 1000 == 0 {
						//	fmt.Println("New minimum", best_coord, "count", best_count, "dist", min_dist)
						//}
					}
				}
			}
		}
		if !updated {
			break
		}
	}
	
	fmt.Println("Result B:", min_dist)
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}