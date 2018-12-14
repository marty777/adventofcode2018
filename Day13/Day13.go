package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"sort"
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

type car struct {
	x int
	y int
	dir orientation
	turn int
	crashed bool
}

func parseInput(lines []string) ([][]trackgrid, []car) {
	width := len(lines[0])
	height := len(lines)
	
	track := make([][]trackgrid, height)
	var cars []car
	for j,line:= range lines {
		track[j] = make([]trackgrid, width)
		for i,c := range line {
			switch c {
				case ' ' :
					track[j][i] = none
				case '-' :
					track[j][i] = horizontal
				case '|' :
					track[j][i] = vertical
				case '+' :
					track[j][i] = intersection
				case '/' :
					track[j][i] = cornerslash
				case '\\' :
					track[j][i] = cornerbackslash
				case '<': 
					track[j][i] = horizontal
					cars = append(cars, car{x:i, y:j, dir:west, turn:0, crashed:false})
				case '>': 
					track[j][i] = horizontal
					cars = append(cars, car{x:i, y:j, dir:east, turn:0, crashed:false})
				case '^': 
					track[j][i] = vertical
					cars = append(cars, car{x:i, y:j, dir:north, turn:0, crashed:false})
				case 'v': 
					track[j][i] = vertical
					cars = append(cars, car{x:i, y:j, dir:south, turn:0, crashed:false})
				default:
					fmt.Print("Unexpected input", c)
			}
		}
	}
	return track, cars
}

func carOrder(cars []car, width int) []int {
	carScores := make([]int, len(cars))
	carScores2 := make([]int, len(cars))
	for i,car := range cars {
		score := car.x + (width*car.y)
		carScores[i] = score
		carScores2[i] = score
	}
	sort.Ints(carScores)
	
	carOrd := make([]int, len(cars))
	for i:=0; i < len(cars); i++ {
		score := carScores[i]
		index := 0
		for j:= 0; j < len(cars); j++ {
			if carScores2[j] == score {
				index = j
				break
			}
		}
		carOrd[i] = index
	}
	
	return carOrd
	
}

func printState(track [][]trackgrid, cars []car) {
	trackStr := " -|+/\\"
	carStr := "^>v<"
	for j := 0; j < len(track); j++ {
		for i:=0; i < len(track[j]); i++ {
			car_found := false
			car_dir := north
			for k:=0; k<len(cars); k++ {
				
				if cars[k].x == i && cars[k].y == j {
					car_found = true
					car_dir = cars[k].dir
				}
				
				
			}
			if car_found {
				fmt.Print(string(carStr[car_dir]))
			} else {
				fmt.Print(string(trackStr[track[j][i]:track[j][i] + 1]))
			}
		}
		fmt.Println()
	}
}

func checkCrash(cars []car, width int) (bool,int,int){
	order := carOrder(cars, width)
	fmt.Println(order, width)
	crashFound := false
	x:=0
	y:=0
	for i:= 0; i < len(order); i++ {
		for j:= i+1; j < len(order); j++ {
			if cars[order[i]].x == cars[order[j]].x && cars[order[i]].y == cars[order[j]].y {
				crashFound = true
				x = cars[order[i]].x
				y = cars[order[i]].y
			}
		}
	}
	return crashFound, x, y
}

func checkCrashSingle(cars []car, index int) bool {
	for i:= 0; i < len(cars); i++ {
		if i == index {
			continue
		}
		if cars[i].x == cars[index].x && cars[i].y == cars[index].y {
			cars[index].crashed = true
			cars[i].crashed = true
			return true
		}
	}
	return false
}

func advance(track [][]trackgrid, cars []car)(bool,int,int) {
	order := carOrder(cars, len(track[0]))
	crashFound := false
	x := 0
	y := 0
	for i:= 0; i < len(order); i++ {
		c := order[i]
		switch track[cars[c].y][cars[c].x] {
			case horizontal:
				if cars[c].dir == east {
					cars[c].x++
				} else if cars[c].dir == west {
					cars[c].x--
				} else {
					fmt.Println("error")
				}
			case vertical:
				if cars[c].dir == north {
					cars[c].y--
				} else if cars[c].dir == south {
					cars[c].y++
				} else {
					fmt.Println("error")
				}
			case cornerslash: 
				if cars[c].dir == north {
					cars[c].dir = east
					cars[c].x++
				} else if cars[c].dir == south {
					cars[c].dir = west
					cars[c].x--					
				} else if cars[c].dir == east {
					cars[c].dir = north
					cars[c].y--
				} else if cars[c].dir == west {
					cars[c].dir = south
					cars[c].y++
				} 
			case cornerbackslash: 
				if cars[c].dir == north {
					cars[c].dir = west
					cars[c].x--
				} else if cars[c].dir == south {
					cars[c].dir = east
					cars[c].x++					
				} else if cars[c].dir == east {
					cars[c].dir = south
					cars[c].y++
				} else if cars[c].dir == west {
					cars[c].dir = north
					cars[c].y--
				} 
			case intersection:
				if cars[c].turn == 0 {
					cars[c].dir = (cars[c].dir + 3) % 4
				} else if cars[c].turn == 1 {
					cars[c].dir = cars[c].dir
				} else {
					cars[c].dir = (cars[c].dir + 1) % 4
				}
				if  cars[c].dir == north {
					cars[c].y--
				} else if cars[c].dir == south {
					cars[c].y++
				} else if cars[c].dir == east {
					cars[c].x++
				} else {
					cars[c].x--
				}
 				
				cars[c].turn = (cars[c].turn + 1) % 3
			default :
				fmt.Println("Error")
		}
		if !crashFound {
			crashFound = checkCrashSingle(cars, c)
			if crashFound {
				x = cars[c].x
				y = cars[c].y
			}
		}
		checkCrashSingle(cars, c)
	}
	
	if crashFound {
		return true, x, y
	}
	return false, 0, 0
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// parse input
	grid, cars := parseInput(lines)
	printState(grid,cars)
	// Part A
	steps:=0
	for {
		crash,x,y := advance(grid, cars)
		if(crash) {
			fmt.Printf("Result A: %d,%d\n", x, y)
			break
		}
		steps++
	}
	// delete crashed
	for i:=0; i<len(cars); i++ {
		if cars[i].crashed {
			cars = append(cars[:i], cars[i+1:]...)
			i--
		}
	}
	
	for len(cars) > 1 {
		advance(grid, cars)
		// delete crashed
		for i:=0; i<len(cars); i++ {
			
			if cars[i].crashed {
				cars = append(cars[:i], cars[i+1:]...)
				i--
			}
		}
		steps++
	}
	
	fmt.Printf("Result B: %d,%d\n", cars[0].x, cars[0].y)
	
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}