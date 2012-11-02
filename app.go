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

func statusImage(r *http.Request) string {
	if !strings.HasSuffix(r.URL.Path, ".png") {
		return UnknownImage
	}
	jobName := r.URL.Path[1 : len(r.URL.Path)-4]

	resp, err := http.Get("http://ci.moozement.net/job/" + jobName + "/api/json")
	if err != nil {
		log.Print("Cannot fetch data")
		return UnknownImage
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return UnknownImage
	}

	var job Job
	err = json.Unmarshal(body, &job)
	if err != nil {
		log.Print("Unmarshal Error")
		return UnknownImage
	}

	switch {
	case job.Color == BlueColor:
		return PassingImage
	case job.Color == RedColor:
		return FailingImage
	}

	return UnknownImage
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	image := statusImage(r)
	http.ServeFile(w, r, image)
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
