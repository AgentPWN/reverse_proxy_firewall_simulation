package main
import (
	"fmt"
	"time"
	"math/rand"
	// "sync"
)
type worker struct{
	id int
	value int
	health chan bool
}

type dispatcher struct{
	value chan int

}
func (w* worker) job (){	
	for healthy := range w.health{
		if healthy{
			fmt.Printf("worker %d: %d\n",w. id, w.value)
		}
	}
}

func (w *worker)health_controller(){
	for ;;{
		if rand.Intn(3)<1{
			w.health<-false
		} else {
			w.health <- true
		}
		time.Sleep(1*time.Second)
		w.value++
	}
}

func (d *dispatcher) dispatch_controller(w_1 worker, w_2 worker){
	value := <- d.value
	fmt.Println("Dispatcher", value)

}

func main(){
	var w_1 worker
	w_1.id = 1
	w_1.health = make(chan bool)
	w_1.value = 1
	go w_1.job()
	go w_1.health_controller()
	
	var w_2 worker
	w_2.id = 2
	w_2.health = make(chan bool)
	w_2.value = 1
	go w_2.job()
	go w_2.health_controller()
	now := time.Now()
	var d dispatcher
	go d.dispatch_controller(w_1, w_2)
	d.value = make(chan int)
	for i:= 0 ;time.Since(now) < 5*time.Second; i++{
		fmt.Println("timer:", i)
		d.value <- i
		time.Sleep(1*time.Second)
	}
	close(w_1.health)
	close(w_2.health)

}
