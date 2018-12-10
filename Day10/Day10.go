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

type particle struct {
	posx int
	posy int
	velx int
	vely int
}

func advance(particles []particle, timestep int) {
	for i:= 0; i < len(particles); i++ {
		delx := particles[i].velx * timestep
		dely := particles[i].vely * timestep
		particles[i].posx += delx
		particles[i].posy += dely
	}
}

func bounds(particles []particle)(int,int,int,int) {
	min_x := 100000
	max_x := -100000
	min_y := 100000
	max_y := -100000
	for i:= 0; i < len(particles); i++ {
		if particles[i].posx < min_x {
			min_x = particles[i].posx
		}
		if particles[i].posy < min_y {
			min_y = particles[i].posy
		}
		if particles[i].posx > max_x {
			max_x = particles[i].posx
		}
		if particles[i].posy > max_y {
			max_y = particles[i].posy
		}
	}
	return  min_x, min_y, max_x, max_y
}

func printParticles(particles []particle) {
	min_x, min_y, max_x, max_y := bounds(particles)
	//fmt.Println(min_x, min_y, max_x, max_y)
	grid := make([][]int, max_x - min_x + 1)
	for i:=0; i < max_x - min_x + 1; i++ {
		grid[i] = make([]int, max_y - min_y + 1)
	}
	//fmt.Println(max_x - min_x + 1, max_y - min_y + 1)
	for i := 0; i < len(particles); i++ {
		//fmt.Println(particles[i].posx, particles[i].posx - min_x, particles[i].posy, particles[i].posy - min_y)
		grid[particles[i].posx - min_x][particles[i].posy - min_y] = 1
	}
	for j := 0; j < max_y - min_y + 1; j++ {
		for i:=0; i < max_x - min_x + 1; i++ {
			if grid[i][j] == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
	
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")

	
	reader := bufio.NewReader(os.Stdin)
	
	// parse input
	particles := make([]particle, len(lines))
	for i:= 0; i<len(lines); i++ {
		var x int
		var y int
		var velx int
		var vely int
		fmt.Sscanf(lines[i], "position=<%d, %d> velocity=<%d,  %d>", &x, &y, &velx, &vely)
		particles[i] = particle{posx:x, posy:y, velx: velx, vely:vely}
	}
		
	// From examining the results, the particles should spell out a word 10 positions tall in my puzzle input
	max_width := len(particles)/2
	max_height := 15
	advanceCount:= 1
	// assuming we start with widely diverged particles, narrow briefly to a small area, and then widely diverge again
	withinBounds := false
	for {
		advance(particles,1)
		min_x, min_y, max_x, max_y := bounds(particles)
		if(max_x - min_x <= max_width && max_y - min_y <= max_height) {
			withinBounds = true
			fmt.Println("Result A (tentative):")
			printParticles(particles)
			fmt.Println("Result B (tentative):", advanceCount)
			fmt.Println("Hit Enter to advance")
			_,_ = reader.ReadString('\n')
		} else {
			// if we entered the boundary area and have now left it
			if withinBounds {
				break
			}
		}
		advanceCount++
	}
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}