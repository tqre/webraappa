package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

func main() {
	go server()
	resp, _ := http.Get("http://localhost:3000/test.html")
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(doc)

}

func server() {
	http.Handle("/", http.FileServer(http.Dir("./testing")))
	http.ListenAndServe(":3000", nil)
}
