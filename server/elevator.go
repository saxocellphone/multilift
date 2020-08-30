package server

import (
	"fmt"
	"sort"
	"time"
)

const nElevator = 5
const ElevatorSpeed = 2

type Elevator struct {
	id          int          // id of the elevator
	i           int          // For heap implementation
	currFloor   int          // The current floor the elevator is on
	dir         int          // 1: going up, -1: going down, 0: stationary
	requests    chan Request // Channel for the elevator to get requests from the balancer
	pending     []int        // The order of the floors the elevator needs to visit
	operational bool
}

func contains(target int, slice []int) bool {
	for _, n := range slice {
		if n == target {
			return true
		}
	}
	return false
}

func (w *Elevator) move() {
	// move one floor, atomic operation cannot be interrupted
	time.Sleep(time.Duration(ElevatorSpeed) * time.Second) // Simulates moving
	w.currFloor += w.dir
}

func (w *Elevator) startCycle(curFloor, destFloor int) {
	// When a new request comes in and the elevator is stationary,
	// move to the requested floor with no interruption
	w.dir = 1
	if curFloor != w.currFloor {
		fmt.Printf("New move cycle detected for elevator %d. Moving to floor %d.\n", w.id, curFloor)
		delta := curFloor - w.currFloor
		if delta < 0 {
			delta = delta * -1
		}
		time.Sleep(time.Duration(delta*ElevatorSpeed) * time.Second) // Simulates moving
		w.currFloor = curFloor
		if destFloor < w.currFloor {
			w.dir = -1
		}
	} else {
		fmt.Printf("New move cycle detected for elevator %d. Already in position. Moving to floor %d.\n", w.id, destFloor)
		delta := destFloor - w.currFloor
		if delta < 0 {
			w.dir = -1
		}
	}
}

func (w *Elevator) operate(dropoff chan *Elevator) {
	for w.operational {
		// Infinite loop to keep looking for passengers
		for {
			select {
			case req := <-w.requests:
				if w.dir != 0 && !contains(req.currFloor, w.pending) {
					w.pending = append(w.pending, req.currFloor)
				}
				if !contains(req.destFloor, w.pending) {
					w.pending = append(w.pending, req.destFloor)
				}
				if w.dir == 0 {
					w.startCycle(req.currFloor, req.destFloor)
				}
				sort.SliceStable(w.pending, func(i int, j int) bool {
					return w.dir*w.pending[i] < w.dir*w.pending[j]
				})
			default:
				// Keeping here to make the request non-blocking
			}
			if len(w.pending) == 0 {
				break
			}
			fmt.Println(w.id, w.dir, w.currFloor, w.pending)
			w.move()
			if len(w.pending) == 0 || w.currFloor == w.pending[0] {
				if len(w.pending) > 0 {
					w.pending = w.pending[1:]
				}
				if len(w.pending) == 0 {
					w.dir = 0
				}
				dropoff <- w
			}
		}
	}
}
