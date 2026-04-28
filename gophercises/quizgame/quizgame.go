package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// 1. Define flags
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 5, "Time limit for the quiz in seconds")
	flag.Parse()

	// 2. Open the CSV file
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	defer file.Close()

	// 3. Read the CSV file
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	// 4. Parse the lines into problems
	problems := parseLines(lines)

	// 5. Run the quiz
	correct := 0
	ch := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(*timeLimit) * time.Second)
		ch <- struct{}{}
	}()

	go func() {
		for i, p := range problems {
			fmt.Printf("Problem #%d: %s = ", i+1, p.q)
			var answer string
			// Using fmt.Scanf to read a single word/number answer.
			// It will stop at whitespace.
			fmt.Scanf("%s\n", &answer)
			if strings.TrimSpace(answer) == p.a {
				correct++
			}
		}
		ch <- struct{}{}
	}()

	defer close(ch)

	<-ch
	// 6. Print the results
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
