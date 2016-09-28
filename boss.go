package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"time"
)

type Boss struct {
	endpoints []string
	client    *http.Client
}

/* Register an endpoint */
func (b *Boss) Register(w http.ResponseWriter, r *http.Request) {
	endpoint := r.FormValue("endpoint")
	b.endpoints = append(b.endpoints, endpoint)
	color.Green("minion added: addr: %s", endpoint)
}

/* Checks if an endpoint is alive but polling every 1sec. */
func (b *Boss) HealthCheck() {
	ticker := time.NewTicker(time.Second)
	for _ = range ticker.C {
		for index, endpoint := range b.endpoints {
			_, err := b.client.Get("http://" + endpoint + "/health")
			if err != nil {
				color.Red("minion dead: addr: %s", endpoint)
				b.endpoints = append(b.endpoints[:index], b.endpoints[index+1:]...)
			}
		}
	}
}

/* Encode the list of endpoints to json for the client to use. */
func (b *Boss) List(w http.ResponseWriter, r *http.Request) {
	list, err := json.Marshal(b.endpoints)
	if err != nil {
		color.Red("JSON marshal error.\n")
		panic(err)
	} else {
		fmt.Fprintf(w, string(list))
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

/* Constructor */
func NewBoss() *Boss {
	boss := &Boss{}
	boss.client = &http.Client{
		Timeout: time.Second * 10,
	}
	boss.endpoints = []string{}
	return boss
}

func main() {
	mastersrv := flag.String("master", "localhost:8080", "addr. for master server")
	flag.Parse()

	boss := NewBoss()
	go boss.HealthCheck()

	http.HandleFunc("/register", boss.Register)
	http.HandleFunc("/list", boss.List)

	home := http.FileServer(http.Dir("./client"))
	http.Handle("/", home)

	color.Green("Starting server on %s.", *mastersrv)
	http.ListenAndServe(*mastersrv, nil)
}
