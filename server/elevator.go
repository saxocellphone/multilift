package server

import (
	"fmt"
	"time"
)

const nElevator = 1
const ElevatorSpeed = 1

type Elevator struct {
	id          int
	i           int // For heap implementation
	currFloor   int
	requests    chan Request
	pending     []int
	operational bool
}

func (w *Elevator) move(isUp bool) {
	// move one floor, atomic operation cannot be interrupted
	time.Sleep(time.Duration(ElevatorSpeed) * time.Second) // Simulates moving
	dir := 1
	if !isUp {
		dir = -1
	}
	w.currFloor += dir
}

func (w *Elevator) operate(done chan *Elevator) {
	for w.operational {
		// Infinite loop to keep looking for passengers
		req := <-w.requests
		fmt.Printf(
			"Elevator %d is currently on %d, and is picking up passenger from floor %d.\n",
			w.id, w.currFloor, req.currFloor)
		isUp := true
		if w.currFloor > req.currFloor {
			isUp = false
		}
		for w.currFloor != req.currFloor {
			w.move(isUp)
		}
		fmt.Printf(
			"Elevator %d has reached passenger from floor %d. Now heading to floor %d.\n",
			w.id, req.currFloor, req.destFloor)
		// TODO: This is not right, an elevator can pick up other passengers
		// when dropping off existing passenger
		isUp = true
		if w.currFloor > req.destFloor {
			isUp = false
		}
		for w.currFloor != req.destFloor {
			w.move(isUp)
		}
		done <- w
	}
}
