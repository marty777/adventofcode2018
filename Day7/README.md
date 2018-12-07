# Advent of Code 2018 - Day 7

## Usage
```
Day7 input.txt
```

## Notes
The sample data from the program description is included (as sample.txt) as well as my copy of the full size puzzle input.

The solution to Part B is sloppy in tracking remaining jobs available to be queued. Any jobs with only inbound dependancies will be lost if there aren't enough worker slots to accept them when their final preceding job nodes are deleted. It works fine for my puzzle input though.