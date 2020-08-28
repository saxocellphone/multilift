package server

import (
	"fmt"
	"strconv"
)

type Request struct {
	currFlorr int
	destFloor int
}

func request(work chan Request, currFloor int, destFloor int) {
	c := make(chan int)
	work <- Request{currFloor, destFloor}
	<-c
}

func GetRequests(work chan Request) {
	for {
		input := ""
		fmt.Scanln(&input)
		floor, _ := strconv.Atoi(input)
		go request(work, 0, floor) // For now, assume all requests are comming from floor 0
	}
}
