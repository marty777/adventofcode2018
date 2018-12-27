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

type treenode struct {
	x int
	y int
	north, east, south, west int
	depth int
	index int
}

// somewhat fragile stack implementation borrowed from
// https://stackoverflow.com/questions/28541609/looking-for-reasonable-stack-implementation-in-golang
type stack []int

func (s stack) Push(v int) stack {
    return append(s, v)
}

func (s stack) Pop() (stack, int) {
	l := len(s)
    return  s[:l-1], s[l-1]
}

func (s stack) Peek() int {
	return s[len(s) - 1]
}

func (s stack) Print() {
	for _,v := range s {
		fmt.Print(" ", v)
	}
	fmt.Println()
}

func addNode(nodes []treenode, parentId int, dir int) ([]treenode, int) {
	newNode := treenode{}
	retNodeId := -1	
	if(dir == 0)  {// north
		if nodes[parentId].north == -1 {
			newNode.south = parentId
			newNode.north = -1
			newNode.east = -1
			newNode.west = -1
			newNode.depth = nodes[parentId].depth + 1
			newNode.x = nodes[parentId].x
			newNode.y = nodes[parentId].y - 1
			nodes = append(nodes, newNode)
			nodes[len(nodes) - 1].index = len(nodes) - 1
			nodes[parentId].north = len(nodes) - 1
			retNodeId = nodes[parentId].north
		} else {
			retNodeId = nodes[parentId].north
		}
	} else if dir == 1 {// east 
		if nodes[parentId].east == -1 {
			newNode.west = parentId
			newNode.north = -1
			newNode.east = -1
			newNode.south = -1
			newNode.depth = nodes[parentId].depth + 1
			newNode.x = nodes[parentId].x + 1
			newNode.y = nodes[parentId].y
			nodes = append(nodes, newNode)
			nodes[len(nodes) - 1].index = len(nodes) - 1
			nodes[parentId].east = len(nodes) - 1
			retNodeId = nodes[parentId].east
		} else {
			retNodeId = nodes[parentId].east
		}
	} else if dir == 2 {// south 
		if nodes[parentId].south == -1 {
			newNode.north = parentId
			newNode.south = -1
			newNode.east = -1
			newNode.west = -1
			newNode.depth = nodes[parentId].depth + 1
			newNode.x = nodes[parentId].x
			newNode.y = nodes[parentId].y + 1
			nodes = append(nodes, newNode)
			nodes[len(nodes) - 1].index = len(nodes) - 1
			nodes[parentId].south = len(nodes) - 1
			retNodeId = nodes[parentId].south
		} else {
			retNodeId = nodes[parentId].south
		}
	} else if dir == 3 {// west 
		if nodes[parentId].west == -1 {
			newNode.east = parentId
			newNode.south = -1
			newNode.north = -1
			newNode.west = -1
			newNode.depth = nodes[parentId].depth + 1
			newNode.x = nodes[parentId].x - 1
			newNode.y = nodes[parentId].y
			nodes = append(nodes, newNode)
			nodes[len(nodes) - 1].index = len(nodes) - 1
			nodes[parentId].west = len(nodes) - 1
			retNodeId = nodes[parentId].west
		} else {
			retNodeId = nodes[parentId].west
		}
	}
	return nodes, retNodeId
}

func parseInput(input string) {
	in := input[1:len(input)-1]
	
	nodes := make([]treenode,0)
	s := make(stack,0)
	
	// root node
	nodes = append(nodes, treenode{x:0,y:0,north:-1,east:-1,south:-1,west:-1,depth:0,index:0})
	
	index := 0
	nodeIndex := 0
	for index < len(in) {
		if in[index] == 'N' {
			nodes, nodeIndex = addNode(nodes, nodeIndex, 0)
		} else if in[index] == 'E' {
			nodes, nodeIndex = addNode(nodes, nodeIndex, 1)
		} else if in[index] == 'S' {
			nodes, nodeIndex = addNode(nodes, nodeIndex, 2)
		} else if in[index] == 'W' {
			nodes, nodeIndex = addNode(nodes, nodeIndex, 3)
		} else if in[index] == '(' {
			s = s.Push(nodeIndex)
		} else if in[index] == '|' {
			nodeIndex = s.Peek()
		} else if in[index] == ')' {
			s, nodeIndex = s.Pop()
		}
		index++
	}
	
	max_depth := 0
	depth_count := 0
	for i:=0; i < len(nodes); i++ {
		if nodes[i].depth > max_depth {
			max_depth = nodes[i].depth
		}
		if(nodes[i].depth >= 1000) {
			depth_count++
		}
	}
	fmt.Println("Part A:", max_depth)
	fmt.Println("Part B:", depth_count)
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// parse input
	parseInput(lines[0])
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}