package server

import (
	"math"
	"sort"
	"time"
)

const nElevator = 2
const ElevatorSpeed = 1

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

func (w *Elevator) operate(dropoff chan *Elevator) {
	for w.operational {
		// Infinite loop to keep looking for passengers
		for {
			select {
			case req := <-w.requests:
				if req.currFloor != w.currFloor && !contains(req.currFloor, w.pending) {
					w.pending = append(w.pending, req.currFloor)
				}
				if !contains(req.destFloor, w.pending) {
					w.pending = append(w.pending, req.destFloor)
				}
			default:
				// Keeping here to make the request non-blocking
			}

			if len(w.pending) == 0 {
				break
			}
			sort.SliceStable(w.pending, func(i int, j int) bool {
				return math.Abs(float64(w.pending[i]-w.currFloor)) < math.Abs(float64(w.pending[j]-w.currFloor))
			})
			w.dir = 1
			if w.currFloor > w.pending[0] {
				w.dir = -1
			}
			w.move()
			if len(w.pending) == 0 || w.currFloor == w.pending[0] {
				if len(w.pending) > 0 {
					w.pending = w.pending[1:]
				}
				dropoff <- w
			}
		}
	}
}
