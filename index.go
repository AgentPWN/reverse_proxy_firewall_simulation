package main
import (
	"fmt"
	"time"
	"math/rand"
	// "sync"
)
type routine struct{
	health bool
	id int
}

func worker(id int, channel chan int){
	for job := range channel{
		fmt.Printf("%d,%d\n",id,job)
	}
}
func main() {
	now := time.Now()
	var channel = make(chan int)
	thread := routine{true, 0}
	go worker(0, channel)
	for i := 0; time.Since(now) < 10*time.Second; i++{

		if rand.Intn(2)<1{
			go worker(i,channel)
		}
		if rand.Intn(3)<2{
			go worker(i,channel)
		}
		channel<-i

		time.Sleep(1*time.Second)
	}
	close(channel)
}