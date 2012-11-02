package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	UnknownImage = "public/unknown.png"
	FailingImage = "public/failing.png"
	PassingImage = "public/passing.png"
	FaviconImage = "public/favicon.ico"

	BlueColor = "blue"
	RedColor  = "red"
)

type Job struct {
	Color string
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".png") {
		http.ServeFile(w, r, UnknownImage)
		return
	}
	jobName := r.URL.Path[1 : len(r.URL.Path)-4]
	log.Print(jobName)

	resp, err := http.Get("http://ci.moozement.net/job/" + jobName + "/api/json")
	if err != nil {
		http.ServeFile(w, r, UnknownImage)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var job Job
	err = json.Unmarshal(body, &job)
	if err != nil {
		log.Print("bad job")
		http.ServeFile(w, r, UnknownImage)
		return
	}
	log.Print(job.Color)

	switch {
	case job.Color == BlueColor:
		http.ServeFile(w, r, PassingImage)
	case job.Color == RedColor:
		http.ServeFile(w, r, FailingImage)
	default:
		http.ServeFile(w, r, UnknownImage)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, FaviconImage)
}

func main() {
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", statusHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListAndServer: ", err)
	}
}
