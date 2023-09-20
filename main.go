package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type problem struct {
	question string
	answer   string
}

func main() {
	file, err := os.Open("problems.csv")

	if err != nil {
		panic(err)
	}
	defer file.Close()

	problems := getProblems(file)
	var correct, wrong int

	for _, v := range problems {
		var input string

		fmt.Printf("What is %s? ", v.question)
		fmt.Scan(&input)

		if v.answer == input {
			correct++
		} else {
			wrong++
		}
	}
	fmt.Println("\nTotal number of questions: ", len(problems))
	fmt.Printf("You answered %v questions correctly and %v questions wrongly.\n", correct, wrong)
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
