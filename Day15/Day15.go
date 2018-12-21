package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	//"sort"
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

type coord struct {
	x int
	y int
}

type unit struct {
	// systemic pro-elf bias
	is_goblin bool 
	alive bool
	coords coord
	has_moved bool
	ap int
	hp int
}

type mapsquare struct {
	obst bool
	//unit *unit
	unit_id int
}

type dijkstranode struct {
	dist int
	visited bool
	parent coord
}

func parseNumUnits(lines[]string)int {
	count:=0
	for i:=0; i<len(lines); i++ {
		for j:= 0; j < len(lines[i]); j++ {
			if lines[i][j] == 'G' || lines[i][j] == 'E' {
				count++;
			}
		}
	}
	return count
}

func parseInput(lines []string, grid *[][]mapsquare, units *[]unit) {
	unitCount:=0
	for i:=0; i<len(lines); i++ {
		for j:= 0; j < len(lines[i]); j++ {
			if lines[i][j] == '#' {
				(*grid)[i][j].obst = true
				(*grid)[i][j].unit_id = -1
			} else {
				(*grid)[i][j].obst = false
				(*grid)[i][j].unit_id = -1
				if lines[i][j] == 'G' {
					(*units)[unitCount] = unit{is_goblin:true, alive:true, coords:coord{x:j,y:i}, has_moved:false, ap:3, hp: 200}
					(*grid)[i][j].unit_id = unitCount
					unitCount++
				} else if lines[i][j] == 'E' {
					(*units)[unitCount] = unit{is_goblin:false, alive:true, coords:coord{x:j,y:i}, has_moved:false, ap:3, hp: 200}
					(*grid)[i][j].unit_id = unitCount
					unitCount++
				}
				
			}
		}
	}
	
}

func openEnemyRangeCoords(grid [][]mapsquare, units []unit, elf bool ) ([]coord) {
	var in_range []coord
	for i:=0; i < len(grid); i++ {
		for j:=0; j<len(grid[i]); j++ {
			// if this square is proximate to an enemy
			if isObstacle(grid, units, j, i){
				continue;
			}
			if(j > 0 && grid[i][j-1].unit_id >= 0 && units[grid[i][j-1].unit_id].is_goblin == elf) { // west
				in_range = append(in_range, coord{x:j, y:i})
			} else if(j < len(grid[i]) - 1 && grid[i][j+1].unit_id >= 0 && units[grid[i][j+1].unit_id].is_goblin == elf) { // east
				in_range = append(in_range, coord{x:j, y:i})
			} else if(i > 0 &&  grid[i-1][j].unit_id >= 0 && units[grid[i-1][j].unit_id].is_goblin == elf) { // north
				in_range = append(in_range, coord{x:j, y:i})
			} else if(i < len(grid) - 1 && grid[i+1][j].unit_id >= 0 && units[grid[i+1][j].unit_id].is_goblin == elf) { // south
				in_range = append(in_range, coord{x:j, y:i})
			}
		}
	}
	return in_range
}

func isObstacle(grid [][]mapsquare, units []unit, x int, y int)bool {
	if( x < 0 || x >= len(grid[0]) || y < 0 || y >= len(grid)) {
		return true
	} else if grid[y][x].obst || (grid[y][x].unit_id >= 0 && units[grid[y][x].unit_id].alive) {
		return true
	}
	return false
}

func getMapIntBoolKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// return min dist from start to end coords
// if min dist of -1 indicates no path
func dijkstra(grid [][]mapsquare, units []unit, start coord, end coord) (int) {
	
	if isObstacle(grid, units, start.x, start.y) || isObstacle(grid, units, end.x, end.y) {
		return -1
	}

	unvisited := make(map[int]bool)
	gridwidth := len(grid[0])
	nodegrid := make([][]dijkstranode,len(grid))
	for i:= 0; i < len(grid); i++ {
		nodegrid[i] = make([]dijkstranode,gridwidth)
		for j:=0; j<gridwidth; j++ {
			nodegrid[i][j].dist = -1 // infinity
			if !isObstacle(grid, units,j,i) {
				unvisited[(i*gridwidth) + j] = true
			}
		}
	}
	
	var currnode coord
	nodegrid[start.y][start.x].dist = 0
	currnode = start
	for {
		currdist := nodegrid[currnode.y][currnode.x].dist
		nextdist := currdist+1
		for i := 0; i < 4; i++ {
			var n coord
			if i == 0 { // north
				n = coord{x:currnode.x, y:currnode.y-1}
			} else if i == 1 { // south
					n = coord{x:currnode.x, y:currnode.y+1}
			} else if i == 2 { // east
					n = coord{x:currnode.x + 1, y:currnode.y}
			} else if i == 3 { // west
					n = coord{x:currnode.x - 1 , y:currnode.y}
			}
			if !isObstacle(grid, units, n.x, n.y) {
				if(!nodegrid[n.y][n.x].visited && (nodegrid[n.y][n.x].dist > nextdist  || nodegrid[n.y][n.x].dist < 0 )) {
					nodegrid[n.y][n.x].dist = nextdist
					nodegrid[n.y][n.x].parent = currnode
				}
			}
		}
		nodegrid[currnode.y][currnode.x].visited = true
		delete(unvisited, currnode.x + (currnode.y*gridwidth))
		
		if(currnode == end) {
			break
		} else {
			// find node in unvisited set with smallest distance
			// if smallest distance is infinity, destination is unreachable
			candidates := getMapIntBoolKeys(unvisited)
			least_x := -1
			least_y := -1
			least_dist := -1
			for _,k := range candidates {
				y:=k/gridwidth
				x:=k - (y*gridwidth)
				if(nodegrid[y][x].dist != -1 && (least_dist == -1 || nodegrid[y][x].dist < least_dist)) {
					least_dist = nodegrid[y][x].dist
					least_x = x
					least_y = y
				}
			}
			if(least_dist == -1) {
				break
			} else {
				currnode.x = least_x
				currnode.y = least_y
			}
		}	
		
	}
	
	if(currnode == end) { // if we reached the target
		return nodegrid[end.y][end.x].dist
	} else { // no path
		return -1
	}
}


func printState(grid [][]mapsquare, units []unit) {
	for i:=0; i < len(grid); i++ {
		for j:= 0; j < len(grid[i]); j++ {
			if grid[i][j].obst {
				print("#")
			} else {
				if grid[i][j].unit_id >= 0 {
					if units[grid[i][j].unit_id].is_goblin {
						print("G")
					} else {
						print("E")
					}
				} else {
					print(".")
				}
			}
		}
		fmt.Println()
	}
}

func containsEnemy(grid [][]mapsquare, units []unit, goblin bool, x int, y int) bool {
	if(x < 0 || x >= len(grid[0]) || y < 0 || y >= len(grid)) {
		return false
	}
	if(grid[y][x].unit_id >= 0 && units[grid[y][x].unit_id].is_goblin == goblin) {
		return true
	}
	return false
}

// substitution for an inline ternary conditional
func inlineCondStr(statement bool, a string, b string)string {
	if statement {
		return a
	}
	return b
}

func attackEnemy(grid [][]mapsquare, units []unit, x int, y int, attacker_id int) {
	if containsEnemy(grid,units,!units[attacker_id].is_goblin, x, y) {
		enemy_id := grid[y][x].unit_id
		units[enemy_id].hp -= units[attacker_id].ap
		if units[enemy_id].hp < 0 {
			units[enemy_id].hp = 0
		}
		if(units[enemy_id].hp == 0) {
			units[enemy_id].alive = false
			grid[y][x].unit_id = -1
		}
	}
}

func unitAttack(grid [][]mapsquare, units []unit, unit_id int) {
	// check if we're currently proximate to an enemy
	if( !containsEnemy(grid, units, !units[unit_id].is_goblin, units[unit_id].coords.x, units[unit_id].coords.y-1 )  && 
		!containsEnemy(grid, units, !units[unit_id].is_goblin, units[unit_id].coords.x-1, units[unit_id].coords.y )  && 
		!containsEnemy(grid, units, !units[unit_id].is_goblin, units[unit_id].coords.x+1, units[unit_id].coords.y )  &&
		!containsEnemy(grid, units, !units[unit_id].is_goblin, units[unit_id].coords.x, units[unit_id].coords.y+1 )  ) {
		return
	}
	
	// attack reading order unit with lowest HP
	min_hp := 200
	temp_coord := units[unit_id].coords
	if(containsEnemy(grid, units, !units[unit_id].is_goblin, temp_coord.x, temp_coord.y-1 ) && units[grid[temp_coord.y-1][temp_coord.x].unit_id].hp < min_hp) {
		min_hp = units[grid[temp_coord.y-1][temp_coord.x].unit_id].hp
	}
	if(containsEnemy(grid, units, !units[unit_id].is_goblin, temp_coord.x-1, temp_coord.y ) && units[grid[temp_coord.y][temp_coord.x-1].unit_id].hp < min_hp) {
		min_hp = units[grid[temp_coord.y][temp_coord.x-1].unit_id].hp
	}
	if(containsEnemy(grid, units, !units[unit_id].is_goblin, temp_coord.x+1, temp_coord.y ) && units[grid[temp_coord.y][temp_coord.x+1].unit_id].hp < min_hp) {
		min_hp = units[grid[temp_coord.y][temp_coord.x+1].unit_id].hp
	}
	if(containsEnemy(grid, units, !units[unit_id].is_goblin, temp_coord.x, temp_coord.y+1 ) && units[grid[temp_coord.y+1][temp_coord.x].unit_id].hp < min_hp) {
		min_hp = units[grid[temp_coord.y+1][temp_coord.x].unit_id].hp
	}
	
	for i:=0; i < 4; i++ {
		var enemy_coords coord
		switch i {
			case 0:
				enemy_coords = coord{x:temp_coord.x, y:temp_coord.y-1}
			case 1:
				enemy_coords = coord{x:temp_coord.x-1, y:temp_coord.y}
			case 2:
				enemy_coords = coord{x:temp_coord.x+1, y:temp_coord.y}
			case 3:
				enemy_coords = coord{x:temp_coord.x, y:temp_coord.y+1}
		}
		if(containsEnemy(grid, units, !units[unit_id].is_goblin, enemy_coords.x, enemy_coords.y) && units[grid[enemy_coords.y][enemy_coords.x].unit_id].hp == min_hp) {
			attackEnemy(grid,units,enemy_coords.x, enemy_coords.y,unit_id)
			break
		}
	}
	
}

func unitMove(grid [][]mapsquare, units []unit, unit_id int) {
	
	// check if we're currently proximate to an enemy
	if( containsEnemy(grid, units, !units[unit_id].is_goblin, units[unit_id].coords.x, units[unit_id].coords.y-1 )  || 
		containsEnemy(grid, units, !units[unit_id].is_goblin, units[unit_id].coords.x-1, units[unit_id].coords.y )  || 
		containsEnemy(grid, units, !units[unit_id].is_goblin, units[unit_id].coords.x+1, units[unit_id].coords.y )  ||
		containsEnemy(grid, units, !units[unit_id].is_goblin, units[unit_id].coords.x, units[unit_id].coords.y+1 )  ) {
		return
	}
	
	enemy_open := openEnemyRangeCoords(grid, units, !units[unit_id].is_goblin)
	
	open_dist := make([]int, 4*len(enemy_open))
	for i:= 0; i < len(enemy_open); i++ {
		// north
		open_dist[i] = dijkstra(grid, units, coord{x:units[unit_id].coords.x, y:units[unit_id].coords.y-1}, enemy_open[i])
		// west
		open_dist[i + len(enemy_open)] = dijkstra(grid, units, coord{x:units[unit_id].coords.x-1, y:units[unit_id].coords.y}, enemy_open[i])
		// east
		open_dist[i + 2*len(enemy_open)] = dijkstra(grid, units, coord{x:units[unit_id].coords.x+1, y:units[unit_id].coords.y}, enemy_open[i])
		// south
		open_dist[i + 3*len(enemy_open)] = dijkstra(grid, units, coord{x:units[unit_id].coords.x, y:units[unit_id].coords.y+1}, enemy_open[i])
	}
	
	
	// find min dist, then take first instance of it for goal
	min_dist := -1
	for i:= 0; i < 4*len(enemy_open); i++ {
		if min_dist == -1 || (open_dist[i] < min_dist && open_dist[i] != -1) {
			min_dist = open_dist[i]
		}
	}
	// if -1, no open moves
	if(min_dist == -1) {
		return;
	}
	
	min_i := -1;
	for i:=0; i < 4*len(enemy_open); i++ {
		if(open_dist[i] == min_dist) {
			min_i = i;
			break;
		}
	}
	
	// move
	start_coord := units[unit_id].coords
	var end_coord coord
	
	if(min_i < len(enemy_open)) { // move north		
		end_coord = coord{x:start_coord.x, y:start_coord.y - 1}
		//fmt.Println("North")
	} else if(min_i < 2*len(enemy_open)) { // move west
		end_coord = coord{x:start_coord.x-1, y:start_coord.y}
		//fmt.Println("West")
	} else if(min_i < 3*len(enemy_open)) { // move east
		end_coord = coord{x:start_coord.x+1, y:start_coord.y}
		//fmt.Println("East")
	} else if(min_i < 4*len(enemy_open)) { // move south
		end_coord = coord{x:start_coord.x, y:start_coord.y+1}
		//fmt.Println("South")
	} else {
		//fmt.Println("Uh oh")
	}
	
	// update position
	units[unit_id].coords = end_coord
	grid[start_coord.y][start_coord.x].unit_id = -1
	grid[end_coord.y][end_coord.x].unit_id = unit_id	
}

func step(grid [][]mapsquare, units []unit) {
	
	for i:=0; i < len(units); i++ {
		units[i].has_moved = false
	}
	// find in-range spaces for each unit type
	
	// for each unit in read order 
	for i:=0; i < len(grid); i++ {
		for j:=0; j<len(grid[0]); j++ {
			if grid[i][j].unit_id >= 0 {
				 if(units[grid[i][j].unit_id].has_moved) {
					continue
				 }
				unit_id := grid[i][j].unit_id
				unitMove(grid, units, unit_id)
				unitAttack(grid, units, unit_id)
				//printState(grid,units)
				units[unit_id].has_moved = true
			}
		}
	}
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// parse input
	grid := make([][]mapsquare, len(lines))
	for i:= 0; i < len(lines); i++ {
		grid[i] = make([]mapsquare, len(lines[i]))
	}
	numUnits := parseNumUnits(lines)
	units := make([]unit, numUnits)
	parseInput(lines, &grid, &units)
	printState(grid, units)

	rounds := 0
	elves_win:=false;
	for {
		step(grid, units)
		//printState(grid,units)
		rounds++
		// check if finished
		elf_count := 0
		goblin_count := 0
		for i:=0; i < len(units); i++ {
			if (units[i].is_goblin && units[i].alive) {
				goblin_count++
			} else if (!units[i].is_goblin && units[i].alive) {
				elf_count++
			}
		}
		
		if(goblin_count == 0) {
			elves_win = true
			break
		} else if elf_count == 0 {
			break
		}
	}
	
	remaining_hp := 0
	for i:=0; i < len(units); i++ {
		if(units[i].is_goblin != elves_win) {
			remaining_hp += units[i].hp
		}
	}
	
	fmt.Printf("Result A: %d x %d = %d\n", rounds-1, remaining_hp, (rounds-1) * remaining_hp )
	
	// reset state
	elf_power := 17 // i'll spare you a lengthy binary search, which is what I did by hand. This is my arrived-at value.
	grid = make([][]mapsquare, len(lines))
	for i:= 0; i < len(lines); i++ {
		grid[i] = make([]mapsquare, len(lines[i]))
	}
	numUnits = parseNumUnits(lines)
	units = make([]unit, numUnits)
	parseInput(lines, &grid, &units)
	
	for i:=0; i < len(units); i++ {
		if !units[i].is_goblin  {
			units[i].ap = elf_power;
		}
	}
	
	rounds = 0
	elves_win=false;
	for {
		step(grid, units)
		rounds++
		// check if finished
		elf_count := 0
		goblin_count := 0
		for i:=0; i < len(units); i++ {
			if (units[i].is_goblin && units[i].alive) {
				goblin_count++
			} else if (!units[i].is_goblin && units[i].alive) {
				elf_count++
			}
		}
		
		if goblin_count == 0  {
			elves_win = true
			break
		} else if elf_count == 0 {
			break
		}
	}
	
	if(elves_win) {
		fmt.Println("Elf Win")
		dead_elves := 0
		for i:=0; i < len(units); i++ {
			if(!units[i].is_goblin && !units[i].alive) {
				dead_elves++
			}
		}
		if dead_elves == 0 {
			fmt.Println("Success")
		} else {
			fmt.Println("Losses:", dead_elves)
		}
	} else {
		fmt.Println("Goblin Win")
	}
	
	remaining_hp = 0
	for i:=0; i < len(units); i++ {
		if(units[i].is_goblin != elves_win) {
			remaining_hp += units[i].hp
		}
	}
	
	fmt.Printf("Result B: %d x %d = %d\n", rounds-1, remaining_hp, (rounds-1) * remaining_hp )
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}