package main
import (
	"fmt"
	// "time"
	"sync"
)

func display(s string) {
    for i := 0; i < 3; i++ {
        fmt.Println(s)
    }
}
func worker(id int){
	fmt.Printf("%d",id)
}
func main() {
	var wg sync.WaitGroup
	for i:=1; i<5; i++{
		wg.Add(1)
		go func(){
			defer wg.Done()
			worker(1)
		}()
		wg.Add(1)
		go func(){
			defer wg.Done()
			worker(2)
		}()
		
	}
	wg.Wait()

}