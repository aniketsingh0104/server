package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var mydata = make([]string, 0)

var member = make(map[string]int)

var ch = make(chan string, 1)

func conString(body io.ReadCloser) string {
	mes, er := ioutil.ReadAll(body)
	if er != nil {
		fmt.Println(er)
	}
	user := string(mes)
	return user
}

func newUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := conString(r.Body)
		if _, ok := member[username]; ok == true {
			w.Write([]byte("Existing Username"))
			return
		}
		ch <- "mem"
		member[username] = 0
		<-ch
	}
}

func getMes(w http.ResponseWriter, r *http.Request) {
	usr := r.URL.Path
	usr = strings.TrimPrefix(usr, "/")
	ch <- "mem"
	lastInd := member[usr]
	var allMes string
	mesArray := mydata[lastInd:]
	for i, v := range mesArray {
		j := strings.Index(v, ":")
		username := v[:j]
		if username != usr {
			allMes += v
			if i != (len(mesArray) - 1) {
				allMes += "\n"
			}
		}
	}
	member[usr] = len(mydata)
	<-ch
	w.Write([]byte(allMes))
}

func postMes(w http.ResponseWriter, r *http.Request) {
	message := conString(r.Body)
	ch <- "mem"
	mydata = append(mydata, message)
	<-ch
}

func final(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		postMes(w, r)
	} else if r.Method == "GET" {
		getMes(w, r)
	}
}

func main() {
	http.HandleFunc("/", final)
	http.HandleFunc("/add", newUser)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
