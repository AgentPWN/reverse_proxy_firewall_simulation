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
func (w* worker) job (){
	
	for healthy := range w.health{
		if healthy{
			fmt.Printf("%d\n",w.value)
		}
		// fmt.Printf("%b\n",healthy)

		
	}
}
func main(){
	var w worker
	w.id = 1
	w.health = make(chan bool)
	w.value = 1
	now := time.Now()
	go w.job()
	// fmt.Printf("hellow world")
	for i:= 0 ;time.Since(now) < 30*time.Second; i++{
		if rand.Intn(3)<1{
			w.health<-false
		} else {
			w.health <- true
		}
		time.Sleep(1*time.Second)
		w.value++
	}
	// time.Sleep(10*time.Second)
	close(w.health)
}
