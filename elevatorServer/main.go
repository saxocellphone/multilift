package server

import (
	"container/heap"
	"fmt"
	"strconv"
	"time"
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

type Pool []*Elevator

func (p Pool) Len() int { return len(p) }

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p *Pool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].i = i
	a[j].i = j
}

func (p *Pool) Push(x interface{}) {
	a := *p
	n := len(a)
	a = a[0 : n+1]
	w := x.(*Elevator)
	a[n] = w
	w.i = n
	*p = a
}

func (p *Pool) Pop() interface{} {
	a := *p
	*p = a[0 : len(a)-1]
	w := a[len(a)-1]
	w.i = -1 // for safety
	return w
}

type Balancer struct {
	pool Pool
	done chan *Elevator
	i    int
}

func NewBalancer() *Balancer {
	done := make(chan *Elevator, nElevator)
	b := &Balancer{make(Pool, 0, nElevator), done, 0}
	for i := 0; i < nElevator; i++ {
		w := &Elevator{id: i, requests: make(chan Request)}
		heap.Push(&b.pool, w)
		go w.operate(b.done)
	}
	return b
}

func (b *Balancer) balance(work chan Request) {
	for {
		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.done:
			b.completed(w)
		}
		// b.print()
	}
}

func (b *Balancer) print() {

}

func (b *Balancer) dispatch(req Request) {
	w := heap.Pop(&b.pool).(*Elevator)
	w.requests <- req
	w.pending++
	fmt.Printf("Elevator %d is currently on %d, and is picking up passenger from floor %d heading to floor %d. \n", w.id, w.currFloor, req.currFlorr, req.destFloor)
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Elevator) {
	w.pending--
	fmt.Printf("Elevator %d has arrived at destination. It has no further requests. \n", w.id)
	heap.Remove(&b.pool, w.i)
	heap.Push(&b.pool, w)
}

func getRequests(work chan Request) {
	for {
		input := ""
		fmt.Scanln(&input)
		floor, _ := strconv.Atoi(input)
		go request(work, 0, floor) // For now, assume all requests are comming from floor 0
	}
}

func main() {
	work := make(chan Request)
	go getRequests(work)
	NewBalancer().balance(work)
}
