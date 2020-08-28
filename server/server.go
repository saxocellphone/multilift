package server

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	currFloor int
	destFloor int
}

func request(work chan Request, currFloor int, destFloor int) {
	c := make(chan int)
	work <- Request{currFloor, destFloor}
	<-c
}

func GetRequests(work chan Request) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		inputArr := strings.Split(input, " ")
		currFloor, destFloor := 0, 0
		if len(inputArr) == 1 {
			destFloor, _ = strconv.Atoi(inputArr[0])
		} else if len(inputArr) == 2 {
			currFloor, _ = strconv.Atoi(inputArr[0])
			destFloor, _ = strconv.Atoi(inputArr[1])
		} else {
			fmt.Println("Error: input error")
			continue
		}
		go request(work, currFloor, destFloor)
	}

}
