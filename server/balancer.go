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
			b.arrived(w)
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	w := b.pool[0]
	picked := false
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
		if (elevator.currFloor > req.currFloor && elevator.dir == 1) || (elevator.currFloor < req.currFloor && elevator.dir == -1) {
			// If elevator is going up but passenger requested from below, we skip. And vice versa
			continue
		}
		if !picked || math.Abs(float64(w.currFloor-req.currFloor)) > math.Abs(float64(elevator.currFloor-req.currFloor)) {
			w = elevator
			picked = true
		}
	}
	if req.currFloor != w.currFloor {
		fmt.Printf(
			"Elevator %d is currently on %d, picking up passengers on floor %d going to floor %d.\n",
			w.id, w.currFloor, req.currFloor, req.destFloor)
	} else {
		fmt.Printf(
			"Elevator %d picked up passenger on floor %d, heading to floor %d.\n",
			w.id, w.currFloor, req.destFloor)
	}
	w.requests <- req
}

func (b *Balancer) arrived(w *Elevator) {
	fmt.Printf("Elevator %d has arrived at %d. Dir: %d \n", w.id, w.currFloor, w.dir)
}
