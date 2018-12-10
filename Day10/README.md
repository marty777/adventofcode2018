# Advent of Code 2018 - Day 10

## Usage
```
Day10 input.txt
```

## Notes
The sample data from the program description is included (as sample.txt) as well as my copy of the full size puzzle input.

My initial solution advanced the particle system until all particles fell within a pre-defined bound, then printed the result to output. Advancing through and examining these gave me the result string. Since I know that the aligned particles produce a pattern with letters 10 units high, this implementation waits until it sees all the particles within a 20-high boundary and prints out one advance of the system (and one step count for Part B) each time the Enter key is pressed to allow the user to verify the correct alignment. The hardcoded boundary may need to be tweaked for alternate inputs.