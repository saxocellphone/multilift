package server

import (
	"fmt"
	"math"
)

type Pool []*Elevator

type Balancer struct {
	pool    Pool
	dropoff chan *Elevator // Channel to communicate from elevator back to balancer
	i       int
}

func NewBalancer() *Balancer {
	done := make(chan *Elevator, nElevator)
	b := &Balancer{make(Pool, 0, nElevator), done, 0}
	for i := 0; i < nElevator; i++ {
		w := &Elevator{id: i, requests: make(chan Request), operational: true}
		b.pool = append(b.pool, w)
		go w.operate(b.dropoff)
	}
	return b
}

func (b *Balancer) Balance(work chan Request) {
	for {
		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.dropoff:
			b.completed(w)
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	w := b.pool[0]
	reqDir := 1 // The direction of the request
	if req.currFloor > req.destFloor {
		reqDir = -1
	}
	for _, elevator := range b.pool {
		if elevator.dir != reqDir && elevator.dir != 0 {
			// We only want to find elevators going in the same direction as the request
			// Or elevator that is stationary
			continue
		}
		if math.Abs(float64(w.currFloor-req.currFloor)) > math.Abs(float64(elevator.currFloor-req.currFloor)) {
			w = elevator
		}
	}
	if req.currFloor != w.currFloor {
		fmt.Printf(
			"Elevator %d is currently on %d, picking up passengers on floor %d going to floor %d.\n",
			w.id, w.currFloor, req.currFloor, req.destFloor)
	} else {
		fmt.Printf(
			"Elevator %d is currently on %d, heading to floor %d.\n",
			w.id, w.currFloor, req.destFloor)
	}
	w.requests <- req
}

func (b *Balancer) completed(w *Elevator) {
	w.dir = 0
	fmt.Printf("Elevator %d has arrived at %d. \n", w.id, w.currFloor)
}
