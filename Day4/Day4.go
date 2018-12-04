 package main

import (
	"fmt"
	"bufio"
	"time"
	"os"
	"strconv"
	"strings"
	"sort"
)

func makeTimestamp() int64 {
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

func getMapStringIntKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func getMapIntIntKeys(m map[int]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func getTimeStampStr(str string) string {
	index1 := 1
	index2 := strings.Index(str, "]")
	return str[index1:index2]
}

func getGuardID(str string) int {
	index1 := strings.Index(str, "Guard #")
	if index1 != -1  {
		index2 := strings.Index(str, " begins")
		guardID, _ := strconv.Atoi(str[index1 + 7: index2])
		return guardID
	} else {
		return -1
	}
}

func getMinute(str string) int {
	minute,_ := strconv.Atoi(str[14:16])
	return minute
}

func getTimeStamp(str string) int {
	month, _ := strconv.Atoi(str[5:7])
	day,_ := strconv.Atoi(str[8:10])
	hour,_ := strconv.Atoi(str[11:13])
	minute,_ := strconv.Atoi(str[14:16])
	return minute + (hour*60) + ((day + month*32) * 3600);
}

func main() {
	starttime := makeTimestamp()
	
	lines, err := readLines(os.Args[1])
	check(err)
	fmt.Println(len(lines), "lines found in input file")
	
	m := make(map[int]string)
	
	// Read lines with associated timestamp index
	for i:= 0; i < len(lines); i++ {
		timeStamp := getTimeStamp(getTimeStampStr(lines[i]))
		m[timeStamp] = lines[i];
	}
	
	// sort timestamp indexes
	var keys []int
	guardTimes :=[][60]int {}
    for k := range m {
        keys = append(keys, k)
    }
    sort.Ints(keys)
	
	guard_index := 0
	guard_indexes := make(map[int]int)
	key_index := 0
	
	// Parse sorted lines into sleep data per guard
	for {
		guardID := getGuardID(m[keys[key_index]])
		value, okay := guard_indexes[guardID]
		curr_guard_index := 0
		if !okay {
			guard_indexes[guardID] = guard_index
			var row [60]int
			for i := 0; i < 60; i++ {
				row[i] = 0
			}
			guardTimes = append(guardTimes, row)
			curr_guard_index = guard_index
			guard_index++
		} else {
			curr_guard_index = value
		}
		key_index++;
		
		sleep_start_index := 0
		sleep_end_index := 0
		for {
			line := m[keys[key_index]]
			if getGuardID(m[keys[key_index]]) > -1 {
				break
			}
			if strings.Index(line, "falls asleep") > -1 {
				sleep_start_index = getMinute(getTimeStampStr(line))
			} else if strings.Index(line, "wakes up") > -1 {
				sleep_end_index = getMinute(getTimeStampStr(line))
				for i := sleep_start_index; i < sleep_end_index; i++ {
					guardTimes[curr_guard_index][i]++
				}
			}
			key_index++
			if key_index >= len(keys) {
				break
			}
		}
		
		if key_index >= len(keys) {
			break
		}
		
		if getGuardID(m[keys[key_index]]) > -1 {
			continue
		}
	}
	
	// Result A
	sleepiest := 0
	sleepiest_count := 0
	guardList := getMapIntIntKeys(guard_indexes)
	
	// find sleepiest guard
	for i := 0; i < len(guardList); i++ {
		sleep_count := 0
		guard_id := guardList[i]
		guard_times_index := guard_indexes[guard_id]
		for j := 0; j < 60; j++  {
			sleep_count += guardTimes[guard_times_index][j]
		}
		if sleep_count > sleepiest_count {
			sleepiest = guard_id
			sleepiest_count = sleep_count
		}
	}
	
	// find sleepiest minute of sleepiest guard
	sleepiest_guard_index := guard_indexes[sleepiest]
	sleepiest_minute := 0
	sleepiest_minute_amount := 0
	for j := 0; j < 60; j++  {
		sleep := guardTimes[sleepiest_guard_index][j]
		if sleep > sleepiest_minute_amount {
			sleepiest_minute = j
			sleepiest_minute_amount = sleep
		}
	}
	
	fmt.Println("Result A: ", sleepiest * sleepiest_minute)
	
	// Result B
	sleepiest = 0
	sleepiest_count = 0
	sleepiest_minute = 0
	for i := 0; i < len(guardList); i++ {
		guard_id := guardList[i]
		guard_times_index := guard_indexes[guard_id]
		for j := 0; j < 60; j++  {
			sleep_count := guardTimes[guard_times_index][j]
			if sleep_count > sleepiest_count {
				sleepiest = guard_id
				sleepiest_count = sleep_count
				sleepiest_minute = j
			}
		}
		
	}
	
	fmt.Println("Result B: ", sleepiest * sleepiest_minute)
	
	endtime := makeTimestamp()
	elapsed := endtime - starttime
	fmt.Println("Elapsed time (milliseconds):", elapsed)
}