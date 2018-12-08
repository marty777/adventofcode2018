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

func readInts(str string)[]int {
	intStrs := strings.Split(str, " ")
	var ints[]int
	for i:=0;i<len(intStrs);i++ {
		a,_ := strconv.Atoi(intStrs[i])
		ints = append(ints, a)
	}
	return ints
}

type treeNode struct {
	parent *treeNode 
	children []*treeNode
	metadata []int
}

// returns end index of node structure in input array
func parseTree(node *treeNode, ints []int, start_index int)int {
	childCount := ints[start_index]
	metaDataCount := ints[start_index+1]
	index := start_index+2
	for i:= 0; i < childCount; i++ {
		var childNode treeNode
		childNode.parent = node
		node.children = append(node.children, &childNode)
		index = parseTree(&childNode, ints, index)
	}
	for i:= 0; i < metaDataCount; i++ {
		node.metadata = append(node.metadata, ints[index])
		index++
	}
	return index
}

func traverseTree(node *treeNode, depth int)int {
	sum := 0
	for i := 0; i < len(node.metadata); i++ {
		sum += node.metadata[i]
	}
	for i:= 0; i < len(node.children); i++ {
		sum += traverseTree(node.children[i], depth+1)
	}
	return sum
}

func traverseTree2(node *treeNode)int {
	sum := 0
	if len(node.children) == 0 {
		for i := 0; i < len(node.metadata); i++ {
			sum += node.metadata[i]
		}
		return sum
	}
	for i := 0; i < len(node.metadata); i++ {
		if node.metadata[i] > 0 && node.metadata[i] <= len(node.children) {
			sum += traverseTree2(node.children[node.metadata[i] - 1])
		}
	}
	return sum
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")

	// parse input
	ints := readInts(lines[0])
	
	// Part A
	var root treeNode
	parseTree(&root, ints, 0)
	sum := traverseTree(&root, 0)

	fmt.Println("Result A:", sum)
	
	// Part B
	sum2 := traverseTree2(&root)
	
	fmt.Println("Result B:", sum2)
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}