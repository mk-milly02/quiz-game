package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var filename = flag.String("source", "problems.csv", "the csv file that contains the problems")

type problem struct {
	question string
	answer   string
}

func main() {
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

	correct, wrong := askQuestions(problems)
	fmt.Println("\nTotal number of questions: ", len(problems))
	fmt.Printf("You answered %v questions correctly and %v questions wrongly.\n", correct, wrong)
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

func askQuestions(questions []problem) (int, int) {
	var correct, wrong int
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
	return correct, wrong
}
