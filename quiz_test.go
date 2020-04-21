package main

import (
	"strings"
	"testing"
)

type QuestionUsrAnswer struct {
	q  Question
	ua string
}

func assertError(t *testing.T, got error) {
	if got == nil {
		t.Errorf("Expected error to occur")
	}
}

func assertNoError(t *testing.T, got error) {
	if got != nil {
		t.Errorf("Error was not expected")
	}
}

func TestLoadQuestions(t *testing.T) {
	t.Run("load with existing file path", func(t *testing.T) {
		csvFilePath := "problems.csv"
		// this is from problems.csv and needs to be refactored.
		want := []Question{
			{"5+5", "10"},
			{"1+1", "2"},
			{"8+3", "11"},
			{"1+2", "3"},
			{"8+6", "14"},
			{"3+1", "4"},
			{"1+4", "5"},
			{"5+1", "6"},
			{"2+3", "5"},
			{"3+3", "6"},
			{"2+4", "6"},
			{"5+2", "7"},
		}
		got, err := LoadQuestions(csvFilePath)
		if err != nil {
			t.Errorf("Error occured but wasn't expected: %v", err)
		}
		if len(got) != len(want) {
			t.Errorf("Number of elements in loaded %v - %d is not same as expected in %v - %d", got, len(got), want, len(want))
		}
	})
	t.Run("try loading non-existent file", func(t *testing.T) {
		invalidFile := ""
		_, got := LoadQuestions(invalidFile)
		assertError(t, got)
	})
}

func TestUpdateScore(t *testing.T) {
	t.Run("correct answer", func(t *testing.T) {
		score := Score{}
		score.Update(true)
		want := Score{correct: 1, total: 1}

		if score != want {
			t.Errorf("Expected score to be %v but got %v", want, score)
		}
	})
	t.Run("incorrect answer", func(t *testing.T) {
		score := Score{}
		score.Update(false)
		want := Score{correct: 0, total: 1}

		if score != want {
			t.Errorf("Expected score to be %v but got %v", want, score)
		}
	})
	t.Run("correct and incorrect - multi question", func(t *testing.T) {
		score := Score{}
		score.Update(false)
		score.Update(true)
		want := Score{correct: 1, total: 2}

		if score != want {
			t.Errorf("Expected score to be %v but got %v", want, score)
		}
	})
}

func TestAskQuestion(t *testing.T) {
	t.Run("correct answer", func(t *testing.T) {
		question := Question{question: "1+1", answer: "2"}
		stringReader := strings.NewReader("2\n")
		got := AskQuestion(question, stringReader)
		want := true

		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("incorrect answer", func(t *testing.T) {
		question := Question{question: "1+1", answer: "2"}
		stringReader := strings.NewReader("3\n")
		got := AskQuestion(question, stringReader)
		want := false

		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("count score all correct", func(t *testing.T) {
		tc := []QuestionUsrAnswer{
			{Question{"1+1", "2"}, "2\n"},
			{Question{"2+2", "4"}, "4\n"},
		}
		var got int
		want := 2
		for _, test := range tc {
			if AskQuestion(test.q, strings.NewReader(test.ua)) {
				got++
			}
		}
		if got != want {
			t.Errorf("Expected number of correct answer %d is not the same as what was given %d", want, got)
		}
	})

}

func TestFlags(t *testing.T) {
	usrArgsAssert := func(t *testing.T, got, want UserArgs) {
		t.Helper()
		if got != want {
			t.Errorf("I got %v but %v was expected.", got, want)
		}
	}
	t.Run("custom questions file", func(t *testing.T) {
		usrArgs := []string{"-questions", "abc.csv"}
		got, err := Flags(usrArgs)
		want := UserArgs{questionfile: "abc.csv"}

		assertNoError(t, err)
		usrArgsAssert(t, got, want)
	})
	t.Run("default question file", func(t *testing.T) {
		usrArgs := []string{}
		got, err := Flags(usrArgs)
		want := UserArgs{questionfile: "problems.csv"}

		assertNoError(t, err)
		usrArgsAssert(t, got, want)
	})
	t.Run("non-existing flag", func(t *testing.T) {
		usrArgs := []string{"-invalid", "abc.csv"}
		_, err := Flags(usrArgs)

		assertError(t, err)

	})
}
