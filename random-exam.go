package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type question struct {
	q   string
	ans string
}

var questionList []question
var score int = 0
var done int = 0
var timeUsed int = 0

func main() {
	// declear flags
	fileFlag := flag.String("file", "./question.csv", "file path")
	timeoutFlag := flag.Int("timeout", 180, "timeout in second")
	countFlag := flag.Int("count", 5, "question count")
	flag.Parse()

	fmt.Println("Time out: ", *timeoutFlag)

	// open the file
	file, err := os.Open(*fileFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read the file line by line then append to questionList
	scanner := bufio.NewScanner(file)

	var line []string
	for scanner.Scan() {
		line = strings.Split(scanner.Text(), ",")
		questionList = append(questionList, question{
			q:   line[0],
			ans: line[1],
		})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// shuffle the questionList
	questionList = shuffle(questionList)
	questionList = questionList[:*countFlag]

	endTimeout := make(chan bool)
	endFinish := make(chan bool)

	go timer(endTimeout, *timeoutFlag)
	go test(endFinish, questionList)

	endTestChecker(endTimeout, endFinish, *countFlag, *timeoutFlag)
}

// functions

func shuffle(list []question) []question {
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result []question
	for len(list) != 0 {
		randIndex := randomSource.Intn(len(list))
		result = append(result, list[randIndex])
		list = append(list[:randIndex], list[randIndex+1:]...)
	}
	return result
}

func endTestChecker(endTimeout, endFinish <-chan bool, numberOfQuestion int, timeLimit int) {
	for {
		select {
		case <-endTimeout:
			fmt.Printf("\n\t--END OF THE TEST--\ntime limit has been reached.\n")
			displayResult(score, timeUsed, numberOfQuestion, timeLimit, done)
			return
		case <-endFinish:
			fmt.Printf("\n\t--END OF THE TEST--\nyou have finished the test.\n")
			displayResult(score, timeUsed, numberOfQuestion, timeLimit, done)
			return
		}
	}
}

func displayResult(score int, timeUsed int, numberOfQuestion int, timeLimit int, done int) {
	fmt.Printf("\n\t--RESULT--\n")
	fmt.Println("Done: ", done, "of", numberOfQuestion)
	fmt.Println("Score: ", score, "/", numberOfQuestion)
	fmt.Println("Time used: ", timeUsed, " of ", timeLimit, "seconds")
	fmt.Println("----------")
}

// goroutines
func timer(endTimeout chan<- bool, timeout int) {
	for t := 0; t < timeout; t++ {
		time.Sleep(1 * time.Second)
		timeUsed++
	}
	close(endTimeout)
}

func test(endFinish chan<- bool, questionList []question) {
	reader := bufio.NewReader(os.Stdin)
	for i, v := range questionList {
		fmt.Println("Question ", i+1)
		fmt.Print(v.q, "= ? : ")
		answer, _ := reader.ReadString('\n')
		answer = strings.Replace(answer, "\n", "", -1)

		if answer == v.ans {
			score++
			fmt.Printf("\nCorrect !\n\n")
		} else {
			fmt.Printf("\nWrong !\n\n")
		}
		done++
	}
	close(endFinish)
}
