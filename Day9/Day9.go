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

func insertAtIndex(arr []int, val int, index int)[]int {
	temp := make([]int, len(arr))
	copy(temp, arr)
	temp = append(temp[:index], val)
	temp = append(temp, arr[index:]...) 
	//arr1 = append(arr1, arr[index-1:]...)
	return temp
}

func deleteAtIndex(arr []int, index int)([]int) {
	temp := make([]int, len(arr))
	copy(temp, arr)
	temp = append(temp[:index], arr[index+1:]...)
	return temp
}

type listnode struct {
	next *listnode
	last *listnode
	val int
}

func insert(node *listnode, curVal int) *listnode {
	currNode := node
	//fmt.Println("Insert: ", curVal, currNode.val)
	for i:= 0; i < 1; i++ {
		currNode = currNode.next
	}
	var newNode listnode
	nextNode := currNode.next
	newNode.val = curVal
	newNode.last = currNode
	newNode.next = nextNode
	nextNode.last = &newNode
	currNode.next = &newNode
	
	return &newNode
}

func at23(node *listnode, curVal int)(int, *listnode) {
	currNode := node
	for i:= 0; i < 7; i++ {
		currNode = currNode.last
	}
	removalVal := currNode.val
	lastNode := currNode.last
	nextNode := currNode.next
	lastNode.next = nextNode
	nextNode.last = lastNode
	return (removalVal + curVal), nextNode
}

func printlist(node *listnode) {
	startVal := node.val
	tempNode := node
	for {
		fmt.Print(tempNode.val)
		fmt.Print(" ")
		tempNode = tempNode.next
		if tempNode.val == startVal {
			break
		}
	}
	fmt.Println()
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")

	// parse input
	numPlayers := 0
	lastMarble := 0
	fmt.Sscanf(lines[0], "%d players; last marble is worth %d points", &numPlayers, &lastMarble)
	
	
	// Part A - initial lazy approach, using an integer slice to represent a circular linked list
	// Issue is with very poor performance on inserts and deletes as list size grows. See Part B.
	scores := make([]int, numPlayers)
	var circle []int
	scoreIndex := 0
	currentMarble := 0
	for i:= 0; i <= lastMarble; i++ {
		if len(circle) == 0 {
			circle = append(circle, i)
			currentMarble = 0
			scores[scoreIndex] += i
		} else {
			if i % 23 == 0 {
				scores[scoreIndex] += i
				removalIndex := (currentMarble - 7) % len(circle)
				// Go modulo is weird for negative numbers
				if removalIndex < 0 {
					removalIndex += len(circle)
				}
				scores[scoreIndex]+=circle[removalIndex]
				circle = deleteAtIndex(circle, removalIndex)
				currentMarble = removalIndex
			} else {
				insertIndex := (currentMarble + 2) % len(circle)
				circle = insertAtIndex(circle, i, insertIndex)
				currentMarble = insertIndex
			}
		}
		scoreIndex = (scoreIndex + 1) % len(scores)
	}
	
	highScore:=0
	for i:=0;i<len(scores);i++ {
		if scores[i] > highScore {
			highScore = scores[i]
		}		
	}
	
	fmt.Println("Result A:", highScore)
	
	
	// Part B - Due to the increase in the number of elements to store, using an int slice has a performance impact on inserts and deletes.
	// A proper linked list is used instead for this portion.
	lastMarble *= 100
	var node0 listnode
	var node1 listnode
	var node *listnode
	node0.val = 0
	node1.val = 1
	node0.last = &node1
	node0.next = &node1
	node1.last = &node0
	node1.next = &node0
	scores2 := make([]int, numPlayers)
	scoreIndex = 0
	node = &node1
	for  i:= 2; i <= lastMarble; i++ {
		if( i % 23 == 0) {
			score, newNode := at23(node, i)
			scores2[scoreIndex] += score
			node = newNode
		} else {
			newNode := insert(node, i)
			node = newNode
		}
		scoreIndex = (scoreIndex + 1) % len(scores2)
	}
	
	highScore=0
	for i:=0;i<len(scores2);i++ {
		if scores2[i] > highScore {
			highScore = scores2[i]
		}
		
	}
	fmt.Println("Result B:", highScore)
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}