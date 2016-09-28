package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"net"
	"net/http"
	"net/url"
	"time"
)

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

type Minion struct {
	mastersrv         string
	register_endpoint string
	srv               string
	client            *http.Client
}

func (m *Minion) isDown(url string) bool {
	_, err := m.client.Get(url)
	if err != nil {
		return true
	}
	return false
}

func (m *Minion) Register() {
	form_value := url.Values{}
	form_value.Add("endpoint", m.srv)
	_, err := m.client.PostForm(m.register_endpoint, form_value)
	if err == nil {
		color.Green("Registered with master server.\n")
	} else {
		color.Red("Failed to register: %v", err)
	}

}

func (m *Minion) Health(w http.ResponseWriter, r *http.Request) {
	// color.Blue("heart beat")
	fmt.Fprintf(w, "alive")
}

func (m *Minion) Task(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "POST" {
		url := r.FormValue("url")
		isdown := m.isDown(url)
		color.Green("Recived task: url: %s isdown %v", url, isdown)
		fmt.Fprintf(w, "%v", isdown)
	} else {
		fmt.Fprintf(w, "Endpoint Doc:\n Route: /isdown\nPOST:Takes in form value \"url\"\nDesc.: Writes true if site is down.\n")
	}
}

func NewMinion(mastersrv string, srv int) *Minion {
	minion := &Minion{}
	minion.client = &http.Client{
		Timeout: time.Second * 5,
	}
	minion.mastersrv = mastersrv
	minion.srv = fmt.Sprintf("%s:%d", GetLocalIP(), srv)
	println(minion.srv)
	minion.register_endpoint = fmt.Sprintf("http://%s/register", mastersrv)
	return minion
}

func main() {
	mastersrv := flag.String("master", "localhost:8080", "addr. for master server")
	minionport := flag.Int("minport", 7070, "addr. for minion server")
	flag.Parse()

	minion := NewMinion(*mastersrv, *minionport)
	minion.Register()

	http.HandleFunc("/isdown", minion.Task)
	http.HandleFunc("/health", minion.Health)

	color.Green("Starting server on %d.", *minionport)
	http.ListenAndServe(fmt.Sprintf(":%d", *minionport), nil)
}
