package server

import "time"

const nElevator = 10
const ElevatorSpeed = 1

type Elevator struct {
	id        int
	i         int
	currFloor int
	requests  chan Request
	pending   int
}

func (w *Elevator) operate(done chan *Elevator) {
	for {
		// Infinite loop to keep looking for passengers
		req := <-w.requests
		delta := 1
		if w.currFloor > req.destFloor {
			delta = -1
		}
		for w.currFloor != req.destFloor {
			w.currFloor += delta
			time.Sleep(time.Duration(ElevatorSpeed) * time.Second)
		}
		done <- w
	}
}
