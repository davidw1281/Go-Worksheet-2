package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

// question struct stores a single question and its corresponding answer.
type question struct {
	q, a string
}

type score int

// check handles a potential error.
// It stops execution of the program ("panics") if an error has happened.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// questions reads in questions and corresponding answers from a CSV file into a slice of question structs.
func questions() []question {
	f, err := os.Open("quiz-questions.csv")
	check(err)
	reader := csv.NewReader(f)
	table, err := reader.ReadAll()
	check(err)
	var questions []question
	for _, row := range table {
		questions = append(questions, question{q: row[0], a: row[1]})
	}
	return questions
}

// ask asks a question and returns an updated score depending on the answer.
func ask(s chan score, complete chan score, questions []question) {
	for _, question := range questions {
		fmt.Println(question.q)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Enter answer: ")
		scanner.Scan()
		text := scanner.Text()
		if strings.Compare(text, question.a) == 0 {
			fmt.Println("Correct!")
			s <- 1
		} else {
			fmt.Println("Incorrect :-(")
		}
	}
	close(s)
	for mark := range s {
		complete <- mark
	}
	close(complete)
}

func main() {
	qs := questions()
	s := make(chan score, len(qs))
	complete := make(chan score, len(qs))
	timer := time.After(5 * time.Second)
	go ask(s, complete, qs)
	select {
	case <-timer:
		fmt.Println("\nTimeout")
		close(s)
		var finalScore score = 0
		for mark := range s {
			finalScore += mark
		}
		fmt.Println("Final score", finalScore)
	case finalScore := <-complete:
		for mark := range complete {
			finalScore += mark
		}
		fmt.Println("Final score", finalScore)
	}
}
