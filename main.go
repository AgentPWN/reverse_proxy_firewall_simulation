package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type worker struct {
	id        int
	value     chan int
	health    bool
	mu        sync.Mutex
	last_used time.Time
}

type dispatcher struct {
	value chan int
}

func (w *worker) job() {
	for value := range w.value {
		fmt.Printf("worker %d: %d\n", w.id, value)
		w.last_used = time.Now()
	}
}

func (w *worker) health_controller() {
	for {
		if rand.Intn(3) < 1 {
			w.mu.Lock()
			w.health = false
			w.mu.Unlock()
			fmt.Printf("Server %d is down\n", w.id)
		} else {
			w.mu.Lock()
			w.health = true
			w.mu.Unlock()
			fmt.Printf("Server %d is up\n", w.id)
		}
		time.Sleep(1 * time.Second)
	}
}

func w_1_used_last(status map[string]*worker) bool {
	return (time.Since(status["w_1"].last_used) < time.Since(status["w_2"].last_used))
}

func isHealthy(status map[string]*worker, worker string) bool {
	return status[worker].health
}
func (d *dispatcher) dispatch_controller(w_1 *worker, w_2 *worker) {
	var status = make(map[string]*worker)
	status["w_1"] = w_1
	status["w_2"] = w_2
	go w_1.job()
	go w_2.job()
	go w_1.health_controller()
	go w_2.health_controller()

	for value := range d.value {
		fmt.Println(w_1.last_used)
		fmt.Println(w_2.last_used)

		switch {
		case isHealthy(status, "w_1") && isHealthy(status, "w_2"):
			if w_1_used_last(status) {
				w_2.value <- value
			} else {
				w_1.value <- value
			}
		case isHealthy(status, "w_1") && !isHealthy(status, "w_2"):
			w_1.value <- value
		case !isHealthy(status, "w_1") && isHealthy(status, "w_2"):
			w_2.value <- value
		case !isHealthy(status, "w_1") && !isHealthy(status, "w_2"):
			fmt.Println("Neither server is up, please try again in a bit")
		}
	}
}

func main() {
	var w_1 worker
	w_1.id = 1
	w_1.health = true
	w_1.value = make(chan int)
	w_1.last_used = time.Now()

	var w_2 worker
	w_2.id = 2
	w_2.health = true
	w_2.value = make(chan int)
	w_2.last_used = time.Now()
	now := time.Now()
	var d dispatcher
	go d.dispatch_controller(&w_1, &w_2)
	d.value = make(chan int)
	for i := 0; time.Since(now) < 30*time.Second; i++ {
		fmt.Println("timer:", i)
		d.value <- i
		time.Sleep(1 * time.Second)
	}
	close(w_1.value)
	close(w_2.value)
	close(d.value)

}
