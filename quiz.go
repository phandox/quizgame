package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Question struct {
	question string
	answer   string
}

type Score struct {
	correct int
	total   int
}

func (s *Score) Update(answer bool) {
	if answer {
		s.correct++
	}
	s.total++
}

func LoadQuestions(filePath string) ([]Question, error) {
	var questions []Question
	f, err := os.Open(filePath)
	if err != nil {
		return []Question{}, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records {
		questions = append(questions, Question{record[0], record[1]})
	}
	return questions, nil
}

func AskQuestion(q Question, r io.Reader) bool {
	reader := bufio.NewReader(r)
	fmt.Printf("%s: ", q.question)
	usrAnswer, _ := reader.ReadString('\n')
	usrAnswer = strings.TrimSuffix(usrAnswer, "\n")
	return usrAnswer == q.answer
}

func main() {
	questions, err := LoadQuestions("problems.csv")
	if err != nil {
		log.Fatal("Can't load the questions.")
	}
	var answer bool
	score := Score{}
	for _, q := range questions {
		answer = AskQuestion(q, os.Stdin)
		score.Update(answer)
	}
	fmt.Println("Quiz is over! Correct answers:", score.correct, " Total questions:", score.total)
}
