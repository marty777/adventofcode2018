# Advent of Code 2018 - Day 9

## Usage
```
Day9 input.txt
```

## Notes
The sample data from the program description is included (as sample.txt).

Part A uses my initial approach, an int slice to stand in for a circular linked list. The performance of inserts and deletes as the list grows makes this unusable in Part B, where a proper linked list is used instead.