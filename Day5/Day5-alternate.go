 package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"strings"
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

func onePass(str1 string) string {
	if len(str1) == 0 {
		return str1
	}
	str := []rune(str1)
	mainIndex := 0
	for {
		index1 := mainIndex
		index2 := mainIndex+1
		for {
			for index1 >= 0 && str[index1] == ' '  {
				index1--
			}
			for index2 <= len(str) - 1 && str[index2] == ' '  {
				index2++
			}
			
			if index1 < 0 || index2 >= len(str) {
				break
			}
			
			if (byte(str[index1]) - byte(str[index2])) == 32 || (byte(str[index2]) - byte(str[index1])) == 32{
				str[index1] = 32
				str[index2] = 32
				if index1 == 0 || index2 == len(str) -1 {
					break
				}
				
			} else {
				break
			}
			
		}
		
		// empty string
		if index1 <= 0 && index2 >= len(str) - 1 {
			break;
		}
		
		if mainIndex == len(str) - 2 {
			break
		}
		mainIndex++ 
	}
	
	output := ""
	for i := 0; i < len(str); i++ {
		if str[i] != ' ' {
			output += string(str[i])
		}
	}
	return output
}

func removeUnit(str string, r string)string {
	str = strings.Replace(str, r, "", -1)
	str = strings.Replace(str, strings.ToUpper(r), "", -1)
	return str
}

func main() {
	starttime := getMillis()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	// Part A
	polymer := onePass(lines[0])
	fmt.Println("Result A:",len(polymer))
	
	// Part B
	
	least_length := len(lines[0])
	for _, r1 := range "abcdefghijklmnopqrstuvwxyz" {
		r:=string(r1)
		polymer = removeUnit(lines[0], r)
		polymer = onePass(polymer)
		if(len(polymer) < least_length) {
			least_length = len(polymer)
		}
	}
	fmt.Println("Result B: ",least_length)
	
	
	endtime := getMillis()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}