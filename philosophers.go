package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	MAX_DELAY    int = 5000
	PHILOSOPHERS int = 5
)

var waitForLock sync.WaitGroup

type Fork struct {
	mutex    *sync.Mutex
	isLocked bool
}

type Philosopher struct {
	num         int
	left, right *Fork
}

func NewFork() *Fork {
	var res Fork
	res.mutex = &sync.Mutex{}
	res.isLocked = false
	return &res
}

func NewPhilosopher(num int, l, r *Fork) *Philosopher {
	res := Philosopher{num, l, r}
	return &res
}

//Philosopher methods
func (p *Philosopher) Eat() {
	dur := rand.Intn(MAX_DELAY)
	fmt.Printf("Philosoper %d will eat %d milliseconds.\n", p.num, dur)
	time.Sleep(time.Duration(dur) * time.Millisecond)
}

func (p *Philosopher) Think() {
	dur := rand.Intn(MAX_DELAY)
	fmt.Printf("Philosoper %d will think %d milliseconds.\n", p.num, dur)
	time.Sleep(time.Duration(dur) * time.Millisecond)
}

func (p *Philosopher) GetForks() bool {
	waitForLock.Wait()
	waitForLock.Add(1)
	defer waitForLock.Done()

	if !p.left.isLocked && !p.right.isLocked {
		p.left.isLocked = true
		p.right.isLocked = true
		p.left.mutex.Lock()
		p.right.mutex.Lock()
		return true
	}
	return false
}

func (p *Philosopher) ReleaseForks() {
	p.left.mutex.Unlock()
	p.right.mutex.Unlock()
	p.left.isLocked = false
	p.right.isLocked = false
}

func (p *Philosopher) Cycle() {
	for {
		for !p.GetForks() {
		}
		p.Eat()
		p.ReleaseForks()
		p.Think()
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var endlessWaitGorup sync.WaitGroup
	endlessWaitGorup.Add(1)

	//init philosophers and forks :)
	var forks [PHILOSOPHERS]*Fork
	var philosophers [PHILOSOPHERS]*Philosopher

	for idx := 0; idx < PHILOSOPHERS; idx++ {
		forks[idx] = NewFork()
	}

	for idx := 0; idx < PHILOSOPHERS; idx++ {
		philosophers[idx] = NewPhilosopher(idx, forks[idx], forks[(idx+1)%PHILOSOPHERS])
	}

	//start the eathinking
	for _, phil := range philosophers {
		fmt.Printf("%d\n", phil.num)
		go phil.Cycle()
	}

	endlessWaitGorup.Wait()
}
