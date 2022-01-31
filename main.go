package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Problem struct {
	Question string
	Answer   int
}

func main() {
	//TODO Filename should be customized with a flag
	f, err := os.Open("./problems.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	problemSet, err := createProblemSet(f)
	if err != nil {
		log.Fatal(err)
	}
	timer := time.NewTimer(time.Duration(30) * time.Second)

	correctAnswers,totalQuestions := startQuiz(problemSet, timer)
	fmt.Printf("Got %d out of %d correct", correctAnswers, totalQuestions)
}

/*
 Read the CSV and convert each record into a problem set with answers
*/
func createProblemSet(f *os.File) ([]Problem, error) {
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var problemSet []Problem

	for _, record := range records {
		answerInt, _ := strconv.Atoi(record[1])
		data := Problem{
			Question: record[0],
			Answer:   answerInt,
		}
		problemSet = append(problemSet, data)

	}

	return problemSet, err
}

/*
 Main logic for quiz game: Starts, stop and returns state of the game.
*/
func startQuiz(ps []Problem, timer *time.Timer) (int, int) {
	correctAnswers := 0
	totalQuestions := len(ps)
	fmt.Println("Starting Quiz...")

	problemLoop:
		for idx,problem := range ps {
			fmt.Printf("Question %d: %s\n", idx+1, problem.Question)
			answerCh := make(chan int)
			go func() {
				var answer int
				fmt.Scanln(&answer)
				answerCh <- answer
			}()
			select {
			case <- timer.C:
				fmt.Println()
				break problemLoop
			case answer := <-answerCh:
				if answer == problem.Answer {
					correctAnswers++
				}
			}
		}
	return correctAnswers,totalQuestions
}