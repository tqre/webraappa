package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

func main() {
	go server()
	resp, _ := http.Get("http://localhost:3000/test.html")
	tok := html.NewTokenizer(resp.Body)
	for {
		tokenType := tok.Next()
		//currentToken := tok.Token()
		if tokenType == html.ErrorToken {
			return
		}

		tag, _ := tok.TagName()

		switch string(tag) {
		case "noscript":
			if tokenType == html.EndTagToken {
				continue
			}
			//frag := html.NewTokenizerFragment(resp.Body, "noscript")
			fmt.Println("<noscript> internals:")
			tok.NextIsNotRawText()
			tok.Next() // If there is text here, this fails
			tag, _ = tok.TagName()
			fmt.Println(string(tag))
			tok.Next()
			tag, _ = tok.TagName()
			fmt.Println(string(tag))

			//tokenType = tok.Next()
			//fmt.Println(tokenType)
		}

		//fmt.Print(currentToken.Type)
		//fmt.Print(currentToken.DataAtom)
		//fmt.Print(currentToken.Data)
		//fmt.Print(currentToken.Attr)
		//fmt.Print(currentToken)
		//fmt.Print(currentToken.Attr)
	}
}

func server() {
	http.Handle("/", http.FileServer(http.Dir("./testing")))
	http.ListenAndServe(":3000", nil)
}
