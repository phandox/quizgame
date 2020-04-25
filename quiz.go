package main

import (
	"bufio"
	"encoding/csv"
	"flag"
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

type UserArgs struct {
	questionfile string
}

func (s *Score) Update(answer bool) {
	if answer {
		s.correct++
	}
	s.total++
}
func Flags(args []string) (UserArgs, error) {
	ua := UserArgs{}
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.StringVar(&ua.questionfile, "questions", "problems.csv", "Path to CSV file with questions in format of 'question,answer'")
	err := fs.Parse(args)
	if err != nil {
		return UserArgs{}, err
	}
	return ua, nil
}

func LoadQuestions(filePath string) ([]Question, error) {
	questions := []Question{}
	f, err := os.Open(filePath)
	if err != nil {
		return []Question{}, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	// Allow only format 'question,answer'
	r.FieldsPerRecord = 2

	records, err := r.ReadAll()
	if err != nil {
		return []Question{}, err
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
	usrArgs, err := Flags(os.Args[1:])
	if err != nil {
		log.Fatalf("Problem parsing user arguments: %v", err)
	}
	questions, err := LoadQuestions(usrArgs.questionfile)
	if err != nil {
		log.Fatalf("Can't load the questions. Error: %v", err)
	}
	var answer bool
	score := Score{}
	for _, q := range questions {
		answer = AskQuestion(q, os.Stdin)
		score.Update(answer)
	}
	fmt.Println("Quiz is over! Correct answers:", score.correct, " Total questions:", score.total)
}
