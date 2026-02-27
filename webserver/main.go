package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type worker struct {
	id        int
	value     chan string
	health    bool
	mu        sync.Mutex
	last_used time.Time
}

type server struct {
	d *dispatcher
}
type dispatcher struct {
	value chan string
}

type response struct {
	Textbox string `json:"textbox"`
}

func loadPage(filename string) []byte {
	html, _ := os.ReadFile("frontend/" + filename + ".html")
	return html
}

func (s *server) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write(loadPage("index"))
	} else if r.Method == "POST" {
		var res response
		json.NewDecoder(r.Body).Decode(&res)
		s.d.value <- res.Textbox

	}
}

func startServer(d *dispatcher) {
	s := server{
		d,
	}
	http.HandleFunc("/", s.Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (w *worker) updateHealthValue(current_status bool) {
	w.mu.Lock()
	w.health = current_status
	w.mu.Unlock()
}

func (w *worker) readHealthValue() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.health
}
func (w *worker) job() {
	for value := range w.value {
		fmt.Printf("worker %d: %s\n", w.id, value)
		w.last_used = time.Now()
	}
}

func (w *worker) health_controller() {
	for {
		if rand.Intn(3) < 1 {
			w.updateHealthValue(false)
			fmt.Printf("Server %d is down\n", w.id)
		} else {
			w.updateHealthValue(true)
			fmt.Printf("Server %d is up\n", w.id)
		}
		time.Sleep(1 * time.Second)
	}
}

func w_1_used_last(status map[string]*worker) bool {
	return (time.Since(status["w_1"].last_used) < time.Since(status["w_2"].last_used))
}

func isHealthy(status map[string]*worker, worker string) bool {
	return status[worker].readHealthValue()
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
	w_1.value = make(chan string)
	w_1.last_used = time.Now()
	var w_2 worker
	w_2.id = 2
	w_2.health = true
	w_2.value = make(chan string)
	w_2.last_used = time.Now()
	now := time.Now()
	var d dispatcher
	go d.dispatch_controller(&w_1, &w_2)
	d.value = make(chan string)
	go startServer(&d)
	for i := 0; time.Since(now) < 30*time.Second; i++ {
		fmt.Println("timer:", i)
		// d.value <- i
		time.Sleep(1 * time.Second)
	}
	close(w_1.value)
	close(w_2.value)
	close(d.value)

}
