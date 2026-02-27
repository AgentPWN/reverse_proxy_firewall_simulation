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

func (fw *firewallData) firewall() {
	for data := range fw.requestData {
		fmt.Println(data.addr, data.lastVisit, data.port, data.requestCount)
	}
}

func loadPage(filename string) []byte {
	html, _ := os.ReadFile("frontend/" + filename + ".html")
	return html
}

func (fw *firewallData) handler(w http.ResponseWriter, r *http.Request) {
	w.Write(loadPage("index"))
	requestAddr := strings.Split(r.RemoteAddr, ":")
	fmt.Println(requestAddr)
	requestIp := strings.Join(requestAddr[:len(requestAddr)-1], ":")
	requestPort := requestAddr[len(requestAddr)-1]
	var req request
	if time.Since(req.lastVisit) > 1*time.Minute {
		req.addr = requestIp
		req.lastVisit = time.Now()
		req.requestCount = 0
		req.port = requestPort
	} else {
		req.requestCount++
	}
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
