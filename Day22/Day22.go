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

type gridsquare struct {
	geo int
	erosion int
	region_type int
}

type coord struct {
	x int
	y int
	tool int
}

type dijkstranode struct {
	x int
	y int
	tool int
	dist int
	visited bool
	parent coord
}



// go modulo operator can produce negative numbers
func modulo(a int, b int)int {
	val := a % b
	if(val < 0) {
		val += b
	}
	return val
}


func getMapIntBoolKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// 0 = torch
// 1 = climbing gear
// 2 = neither
func isObstacle(grid [][][]gridsquare, x int, y int, tool int)bool {
	// torch - rocky or narrow
	// climbing gear - rocky or wet
	// neither - wet or narrow
	if x < 0 || y < 0 || x >= len(grid) || y >= len(grid[0]) {
		return true
	}
	if(tool == 0 && (grid[x][y][tool].region_type == 0 || grid[x][y][tool].region_type == 2)) {
		return false
	} else if(tool == 1 && (grid[x][y][tool].region_type == 0 || grid[x][y][tool].region_type == 1)) {
		return false
	} else if(tool == 2 && (grid[x][y][tool].region_type == 1 || grid[x][y][tool].region_type == 2)) {
		return false
	}
	return true
}

// return min dist from start to end coords
// if min dist of -1 indicates no path
func dijkstra(grid [][][]gridsquare, start coord, end coord) (int) {
	
	unvisited := make(map[int]bool)
	gridwidth := len(grid)
	//gridheight := len(grid[0])
	nodegrid := make([][][]dijkstranode,len(grid))
	for i:= 0; i < len(grid); i++ {
		nodegrid[i] = make([][]dijkstranode,len(grid[i]))
		for j:=0; j<len(grid[i]); j++ {
			nodegrid[i][j] = make([]dijkstranode, 3)
			for k:=0; k<3; k++ {
				nodegrid[i][j][k].dist = -1 // infinity
				nodegrid[i][j][k].x = i
				nodegrid[i][j][k].y = j
				nodegrid[i][j][k].tool = k
				if !isObstacle(grid,i,j,k) {
					unvisited[((j*gridwidth) + i)*3 + k] = true
				}
			}
		}
	}
	
	var currnode coord
	nodegrid[start.x][start.y][start.tool].dist = 0
	currnode = start
	steps := 0
	for {
		currdist := nodegrid[currnode.x][currnode.y][currnode.tool].dist
		nextdist := currdist+1
		
		for i := 0; i < 6; i++ {
			var n coord
			if i == 0 { // north
				n = coord{x:currnode.x, y:currnode.y-1, tool:currnode.tool}
			} else if i == 1 { // south
					n = coord{x:currnode.x, y:currnode.y+1, tool:currnode.tool}
			} else if i == 2 { // east
					n = coord{x:currnode.x + 1, y:currnode.y, tool:currnode.tool}
			} else if i == 3 { // west
					n = coord{x:currnode.x - 1 , y:currnode.y, tool:currnode.tool}
			} else if i == 4 { // change tool 1
					n = coord{x:currnode.x , y:currnode.y, tool:(currnode.tool+1) % 3}
			} else if i == 5 { // change tool 2
					n = coord{x:currnode.x , y:currnode.y, tool:(currnode.tool+2) % 3}
			}
			if !isObstacle(grid, n.x, n.y, n.tool) {
				nextdist2 := nextdist
				if(i == 4 || i == 5) {
					nextdist2+=6
				}
				if(!nodegrid[n.x][n.y][n.tool].visited && (nodegrid[n.x][n.y][n.tool].dist > nextdist2  || nodegrid[n.x][n.y][n.tool].dist < 0 )) {
					nodegrid[n.x][n.y][n.tool].dist = nextdist2
					nodegrid[n.x][n.y][n.tool].parent = currnode
				}
			}
		}
		nodegrid[currnode.x][currnode.y][currnode.tool].visited = true
		delete(unvisited, (currnode.x + (currnode.y*gridwidth))*3 + currnode.tool)
		
		if(currnode == end) {
			break
		} else {
			// find node in unvisited set with smallest distance
			// if smallest distance is infinity, destination is unreachable
			candidates := getMapIntBoolKeys(unvisited)
			least_x := -1
			least_y := -1
			least_tool := -1
			least_dist := -1
			for _,k := range candidates {
				tool:=k % 3
				k = (k - tool)/3
				y:=k/gridwidth
				x:=k - (y*gridwidth)
				if(nodegrid[x][y][tool].dist != -1 && (least_dist == -1 || nodegrid[x][y][tool].dist < least_dist)) {
					least_dist = nodegrid[x][y][tool].dist
					least_x = x
					least_y = y
					least_tool = tool
				}
			}
			if(least_dist == -1) {
				break
			} else {
				currnode.x = least_x
				currnode.y = least_y
				currnode.tool = least_tool
			}
		}	
		steps++
	}
	
	if(currnode == end) { // if we reached the target
		temp := currnode
		temp_node := nodegrid[temp.x][temp.y][temp.tool]
		for {
			temp = temp_node.parent
			temp_node = nodegrid[temp.x][temp.y][temp.tool]
			if(temp.x == start.x && temp.y == start.y && temp.tool == start.tool) {
				break
			}
		}
		
	
		return nodegrid[end.x][end.y][end.tool].dist
	} else { // no path
		return -1
	}
}


func generate(grid [][][]gridsquare, depth int, width int, height int, target_x int, target_y int)int {
	risk := 0
	for i := 0; i < len(grid); i++ {
		grid[i][0][0].geo = 16807*i
		grid[i][0][0].erosion = modulo(grid[i][0][0].geo + depth,20183)
		grid[i][0][0].region_type = modulo(grid[i][0][0].erosion,3)
		if i <= target_x {
			risk += grid[i][0][0].region_type
		}
		grid[i][0][1] = grid[i][0][0]
		grid[i][0][2] = grid[i][0][0]
	}
	
	for i := 1; i < len(grid[0]); i++ {
		grid[0][i][0].geo = 48271*i
		grid[0][i][0].erosion = modulo(grid[0][i][0].geo + depth,20183)
		grid[0][i][0].region_type = modulo(grid[0][i][0].erosion,3)
		if i <= target_y  {
			risk += grid[0][i][0].region_type
		}
		grid[0][i][1] = grid[0][i][0]
		grid[0][i][2] = grid[0][i][0]
	}
	
	for i := 1; i < len(grid); i++ {
		for j:= 1; j < len(grid[i]); j++ {
			if(i == target_x && j == target_y) {
				grid[i][j][0].geo = 0
			} else {
				grid[i][j][0].geo = grid[i-1][j][0].erosion * grid[i][j-1][0].erosion
			}
			grid[i][j][0].erosion = modulo(grid[i][j][0].geo + depth,20183)
			grid[i][j][0].region_type = modulo(grid[i][j][0].erosion,3)
			if (i <= target_x && j <= target_y) {
				risk += grid[i][j][0].region_type
			}
			grid[i][j][1] = grid[i][j][0]
			grid[i][j][2] = grid[i][j][0]
		}
	}
	return risk
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	var depth, target_x, target_y int
	fmt.Sscanf(lines[0], "depth: %d", &depth)
	fmt.Sscanf(lines[1], "target: %d,%d", &target_x, &target_y)
	
	// these bounds worked for my specific input on part 2. 
	width := target_x * 7
	height := target_y * 2
	
	grid :=make([][][]gridsquare, width)
	for i:=0; i < width; i++ {
		grid[i] = make([][]gridsquare, height)
		for j:=0; j < height; j++ {
			grid[i][j] = make([]gridsquare,3)
		}
	}
		
	risk := generate(grid, depth, width, height, target_x, target_y)
	
	fmt.Println("Result A:", risk)
	
	dist_torch := dijkstra(grid, coord{x:0,y:0,tool:0}, coord{x:target_x,y:target_y, tool:0})
	
	fmt.Println("Result B:", dist_torch)
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}