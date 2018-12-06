package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"unsafe"
)

var mydata = make([]string, 0)

var member = make(map[string]int)

var ch = make(chan string)

func emptyChan() {
	for {
		<-ch
	}
}

func conString(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	b := buf.Bytes()
	user := *(*string)(unsafe.Pointer(&b))
	return user
}

func newUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := conString(r.Body)
		ch <- username
		member[username] = 0
	}
}

func getMes(w http.ResponseWriter, r *http.Request) {
	usr := r.URL.Path
	usr = strings.TrimPrefix(usr, "/")
	ch <- usr
	lastInd := member[usr]
	var allMes string
	mesArray := mydata[lastInd:]
	for i, v := range mesArray {
		allMes += v
		if i != (len(mesArray) - 1) {
			allMes += "\n"
		}
	}
	member[usr] = len(mydata)
	w.Write([]byte(allMes))
}

func postMes(w http.ResponseWriter, r *http.Request) {
	message := conString(r.Body)
	var username string
	i := strings.Index(message, ":")
	username = message[:i]
	ch <- username
	mydata = append(mydata, message)
	member[username] = len(mydata)
}

func final(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		postMes(w, r)
	} else if r.Method == "GET" {
		getMes(w, r)
	}
}

func main() {
	go emptyChan()
	http.HandleFunc("/", final)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
	http.HandleFunc("/add", newUser)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
