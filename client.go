package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func user(username *string) {
	var u string
	fmt.Println("Enter your username:")
	fmt.Scanf("%s", &u)
	_, er := http.Post("http://localhost:8080/add", "Test", bytes.NewBuffer([]byte(u)))
	if er != nil {
		fmt.Println(er)
	}
	*username = u
}

func get(username string) {
	for {
		resp, er := http.Get("http://localhost:8080/" + username)
		if er != nil {
			log.Fatalln(er)
		}
		mes, er := ioutil.ReadAll(resp.Body)
		if er != nil {
			log.Fatalln(er)
		}
		if string(mes) != "" {
			fmt.Println(string(mes))
		}
		time.Sleep(time.Second)
	}
}

func main() {
	var username string
	var mes string
	user(&username)
	go get(username)
	for {
		reader := bufio.NewReader(os.Stdin)
		mes, _ = reader.ReadString('\n')
		resp, _ := http.Post("http://localhost:8080", "Text", bytes.NewBuffer([]byte(username+":"+mes)))
		_ = resp
	}
}
