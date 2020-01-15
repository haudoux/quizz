//The first part is mine
//The second part is another by someone else

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var filenameFlag = flag.String("filename", "", "Set a filename for csv")
var timeoutFlag = flag.String("timeout", "", "Set a timeout for the quizz")

func init() {
	//flag.StringVar(filenameFlag, "filename", "", "Set a filename for csv")
}
func main() {
	flag.Parse()
	filename := "quizz.csv"
	timeout := 30
	var nbOfQuestions *int = new(int)
	var nbOfErrors *int = new(int)
	var nbOfValid *int = new(int)
	c1 := make(chan bool, 1)
	if strings.Compare(*filenameFlag, "") != 0 {
		filename = *filenameFlag
	}
	if strings.Compare(*timeoutFlag, "") != 0 {
		timeout, _ = strconv.Atoi(*timeoutFlag)
	}
	csvString := getFileData(filename)
	go func() {
		quizz(csvString, nbOfQuestions, nbOfErrors, nbOfValid)
		c1 <- true
	}()
	select {
	case <-c1:
		break
	case <-time.After(time.Duration(timeout) * time.Second):
		fmt.Println("Fail you are too slow !!!")
	}
	fmt.Println("Numbers of questions : " + strconv.Itoa(*nbOfQuestions))
	fmt.Println("Numbers of good answers : " + strconv.Itoa(*nbOfValid))
	fmt.Println("Numbers of wrong answers : " + strconv.Itoa(*nbOfErrors))
}

func getFileData(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func quizz(csvString string, nbOfQuestions *int, nbOfErrors *int, nbOfValid *int) {
	csv := csv.NewReader(strings.NewReader(csvString))
	*nbOfQuestions = 0
	*nbOfErrors = 0
	*nbOfValid = 0
	for {
		result, err := csv.Read()
		//Si fin fichier
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Problem # " + result[0])
		fmt.Println("Answer : ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.ToLower(strings.Trim(strings.Replace(answer, "\r\n", "", -1), " "))
		res := strings.ToLower(strings.Trim(result[1], " "))
		if strings.Compare(answer, res) != 0 {
			fmt.Println("Wrong, answer is " + result[1])
			*nbOfErrors++
		} else {
			*nbOfValid++
		}
		*nbOfQuestions++
	}
}

/*
//Second part
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	filename := flag.String("filename", "quizz.csv", "Set a filename for csv")
	timeout := flag.Int("timeout", 30, "Set a timeout for quizz")
	flag.Parse()
	problems := openFile(filename)
	quizz(problems, *timeout)
}
func openFile(filename *string) (problems []problem) {
	file, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV File: %s\n", *filename))
	}
	rd := csv.NewReader(file)
	lines, err := rd.ReadAll()
	if err != nil {
		exit("Failds to parse the CSV file")
	}
	problems = parseLines(lines)
	return problems
}
func quizz(problems []problem, timeout int) {
	correct := 0
	timer := time.NewTimer(time.Duration(timeout) * time.Second)

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
			if answer == p.answer {
				fmt.Println("Correct")
				correct++
			}
		}()

		select {
		case <-timer.C:
			break
		case answer := <-answerCh:
			if answer == p.answer {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
*/
