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
	for {
		fmt.Println("Enter your username:")
		fmt.Scanf("%s", &u)
		resp, er := http.Post("http://localhost:8080/add", "Test", bytes.NewBuffer([]byte(u)))
		if er != nil {
			fmt.Println(er)
		}
		*username = u
		mes, er := ioutil.ReadAll(resp.Body)
		if er != nil {
			fmt.Println(er)
		}
		if string(mes) != "Existing Username" {
			*username = u
			break
		}
		fmt.Println("This username already exits")
	}
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
	fmt.Println("Type ESC to exit from chat")
	for {
		reader := bufio.NewReader(os.Stdin)
		mes, _ = reader.ReadString('\n')
		mes = mes[:len(mes)-1]
		if mes == "ESC" {
			resp, err := http.Post("http://localhost:8080/delete", "Text", bytes.NewBuffer([]byte(username)))
			if err != nil {
				fmt.Println(err)
			}
			mes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(mes)
			os.Exit(0)
		}
		_, err := http.Post("http://localhost:8080/", "Text", bytes.NewBuffer([]byte(username+":"+mes)))
		if err != nil {
			fmt.Println(err)
		}
	}
}
