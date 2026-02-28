package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type request struct {
	addr         string
	port         string
	lastVisit    time.Time
	requestCount int
}

type firewallData struct {
	requestData chan request
}

type requestStruct struct { //try and get rid of this struct, this is causing some redundancy of data
	port         string
	lastVisit    time.Time
	requestCount int
}

var state = make(map[string]requestStruct)

// state = make(map[string]requestStruct)

func (fw *firewallData) firewall() {
	for data := range fw.requestData {
		// fmt.Println(data.addr, data.lastVisit, data.port, data.requestCount)
		if v, ok := state[data.addr]; !ok {
			state[data.addr] = requestStruct{
				data.port,
				data.lastVisit,
				data.requestCount,
			}
		} else {
			if time.Since(v.lastVisit) > 1*time.Minute {
				v.lastVisit = time.Now()
				v.requestCount = 0
			}
			v.requestCount ++
			state[data.addr] = v
		}
		fmt.Println(state[data.addr])
	}
}

func loadPage(filename string) []byte {
	html, _ := os.ReadFile("frontend/" + filename + ".html")
	return html
}

func (fw *firewallData) handler(w http.ResponseWriter, r *http.Request) {
	w.Write(loadPage("index"))
	requestAddr := strings.Split(r.RemoteAddr, ":")
	// fmt.Println(requestAddr)
	requestIp := strings.Join(requestAddr[:len(requestAddr)-1], ":")
	requestPort := requestAddr[len(requestAddr)-1]
	var req request
	req.addr = requestIp
	req.lastVisit = time.Now()
	req.requestCount = 0
	req.port = requestPort
	fw.requestData <- req
}

func (fw *firewallData) webserver() {
	http.HandleFunc("/", fw.handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	var fw firewallData
	fw.requestData = make(chan request)
	go fw.firewall()
	fw.webserver()
	close(fw.requestData)
}
