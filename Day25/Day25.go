package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"strings"
	"strconv"
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


type star [4]int

type constellation struct {
	stars []star
}

func abs(a int)int {
	if(a < 0) {
		return -a
	}
	return a
}

func dist(a star, b star)int {
	return abs(a[0] - b[0]) + abs(a[1] - b[1]) + abs(a[2] - b[2]) + abs(a[3] - b[3])
}

func (c *constellation) inRange(s star, r int)bool {
	for i:=0; i < len(c.stars); i++ {
		if dist(s, c.stars[i]) <= r {
			return true
		}
	}
	return false
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	var sky []star
	// Parse input
	for i:=0; i < len(lines); i++ {
		coords := strings.Split(lines[i], ",")
		var newStar star
		for j:=0; j<len(coords);j++ {
			newStar[j],_ = strconv.Atoi(coords[j])
		}
		sky = append(sky, newStar)
	}
	
	connected := make([]bool, len(sky))
	var constellations []constellation
	for {
		start := 0
		for i:=0; i < len(connected); i++ {
			if !connected[i] {
				start = i
				break
			}
		}
		var newConstellation constellation
		newConstellation.stars = append(newConstellation.stars, sky[start])
		for {
			additions := 0
			for i:=start; i < len(sky); i++ {
				if !connected[i] && newConstellation.inRange(sky[i], 3) {
					newConstellation.stars = append(newConstellation.stars, sky[i])
					connected[i] = true
					additions++
				}
			}
			if(additions == 0) {
				break
			}
		}
		constellations = append(constellations, newConstellation)
		
		unconnected := 0
		for i:=0; i < len(connected); i++ {
			if !connected[i] {
				unconnected++
			}
		}
		if(unconnected == 0) {
			break
		}
	}
	
	fmt.Println("Result A:", len(constellations));
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}