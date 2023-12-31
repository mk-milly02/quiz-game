package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	filename := flag.String("source", "problems.csv", "the csv file that contains the problems")
	timelimt := flag.Duration("limit", 30*time.Second, "time limit in seconds")
	shuffle := flag.Bool("shuffle", false, "changes the order in which the questions are asked")
	flag.Parse()

	file, closer, err := getProblemsFile(*filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer closer()

	problems := getProblems(file)
	if len(problems) == 0 {
		fmt.Println("No questions available.")
		os.Exit(1)
	}

	if *shuffle {
		problems = shuffleQuestions(problems)
	}

	fmt.Println("Welcome to the QUIZ GAME...")

	// Create a channel to wait for Enter key press
	enterkeyPressed := make(chan bool)

	// Listen for a signal (e.g., Enter key press) in a goroutine
	go func() {
		var input string
		fmt.Print("Press enter to start the quiz:")
		fmt.Scanln(&input)
		enterkeyPressed <- true
		close(enterkeyPressed)
	}()

	// Wait for Enter key press
	<-enterkeyPressed

	askQuestions(problems, timelimt)
}

func getProblemsFile(name string) (*os.File, func(), error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	return file, func() { file.Close() }, nil
}

func getProblems(r io.Reader) (problems []problem) {
	cr := csv.NewReader(r)
	for {
		field, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		problems = append(problems, problem{field[0], strings.TrimSpace(field[1])})
	}
	return problems
}

func askQuestions(questions []problem, duration *time.Duration) {
	// Create a channel to receive a done signal
	done := make(chan bool)

	// Start the timer
	timer := time.NewTimer(*duration)

	var correct, wrong int

	// Ask questions
	go func() {
		for _, v := range questions {
			var input string
			fmt.Printf("What is %s? ", v.question)
			fmt.Scan(&input)

			if v.answer == input {
				correct++
			} else {
				wrong++
			}
		}
		done <- true
		close(done)
	}()

	// Wait for either the timer to expire or user to complete
	select {
	case <-timer.C:
		fmt.Println("\nTime elapsed.")
		fmt.Println("Total number of questions: ", len(questions))
		fmt.Printf("You answered %v questions correctly and %v questions wrongly.\n", correct, len(questions)-correct)
	case <-done:
		fmt.Println("\nCompleted successfully.")
		fmt.Println("Total number of questions: ", len(questions))
		fmt.Printf("You answered %v questions correctly and %v questions wrongly.\n", correct, wrong)
	}
}

func shuffleQuestions(questions []problem) []problem {
	size := len(questions)
	for i := size - 1; i > 0; i-- {
		j := rand.Intn(len(questions))
		tmp := questions[i]

		questions[i].question = questions[j].question
		questions[i].answer = questions[j].answer

		questions[j].question = tmp.question
		questions[j].answer = tmp.answer
	}
	return questions
}
