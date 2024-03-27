package main

import (
	"bufio"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	stringArray := fetchContent()
	questions := buildQuestions(stringArray)
	shuffle(questions)
	errors := buildQuestionsFromErrorFile(questions)
	shuffle(errors)
	startAsking(questions, errors)
}

func buildQuestionsFromErrorFile(questions []Question) []Question {
	errorFileContent := readFile("errors.md")
	errors := make([]Question, 0)
	for _, line := range errorFileContent {
		idx := slices.IndexFunc(questions, func(q Question) bool { return q.Query == line })
		errors = append(errors, questions[idx])
	}
	return errors
}

func startAsking(questions []Question, errors []Question) {
	for i := 0; i < len(questions); i++ {
		ask(questions[i])
		if i < len(errors) {
			ask(errors[i])
		}
	}
}

func ask(q Question) {
	fmt.Println("\n" + q.Query)
	expected := ""
	for index, answer := range q.Answer {
		fmt.Println(index+1, answer.Answer)
		if answer.Correct {
			expected += strconv.Itoa(index + 1)
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Text() == expected {
		fmt.Println("Correct")
	} else {
		fmt.Println("False", "("+expected+")")
		saveError(q.Query)
	}
}

func saveError(query string) {
	f, err := os.OpenFile("errors.md", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file, err2 := os.ReadFile("errors.md")
	if err2 != nil {
		panic(err2)
	}
	content := string(file)
	if strings.Contains(content, query) {
		return
	}
	_, err = f.WriteString(query + "\n")
	if err != nil {
		panic(err)
	}
}

func shuffle(questions []Question) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})
}

func buildQuestions(stringArray []string) []Question {
	questions := make([]Question, 0)
	for _, line := range stringArray {
		if strings.HasPrefix(line, "###") {
			q := Question{Query: line, Answer: make([]Answer, 0)}
			questions = append(questions, q)
		} else {
			if strings.HasPrefix(line, "- [ ] ") {
				addAnswer(line, false, questions)
			}
			if strings.HasPrefix(line, "- [x] ") {
				addAnswer(line, true, questions)
			}
		}
	}
	return questions
}

func fetchContent() []string {
	name := "downloaded-questions.md"
	out, _ := os.Create(name)
	defer out.Close()

	resp, _ := http.Get("https://raw.githubusercontent.com/Ditectrev/Professional-Scrum-Developer-I-PSD-I-Practice-Tests-Exams-Questions-Answers/master/README.md")
	defer resp.Body.Close()

	_, _ = io.Copy(out, resp.Body)

	return readFile(name)
}

func readFile(name string) []string {
	file, _ := os.Open(name)
	defer file.Close()

	stringArray := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) > 0 {
			stringArray = append(stringArray, text)
		}
	}
	return stringArray
}

func addAnswer(line string, correct bool, questions []Question) {
	answer := Answer{Answer: line[6:], Correct: correct}
	q := &questions[len(questions)-1]
	q.Answer = append(q.Answer, answer)
}

type Question struct {
	Id     int
	Query  string
	Answer []Answer
}

type Answer struct {
	Answer  string
	Correct bool
}
