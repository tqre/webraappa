package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

func main() {
	go server()

	// Normal functionality with readers
	resp1, _ := http.Get("http://localhost:3000/test.html")
	out1a, _ := os.Create("copy1a.html")
	io.Copy(out1a, resp1.Body)
	// copy1b.html will be empty, as resp.Body is of type io.ReadCloser
	out1b, _ := os.Create("copy1b.html")
	io.Copy(out1b, resp1.Body)

	// TeeReader
	resp2, _ := http.Get("http://localhost:3000/test.html")
	out2a, _ := os.Create("copy2a.html")
	var buffer bytes.Buffer
	tee := io.TeeReader(resp2.Body, &buffer)

	io.Copy(out2a, tee)
	out2b, _ := os.Create("copy2b.html")
	io.Copy(out2b, &buffer)

}

func server() {
	http.Handle("/", http.FileServer(http.Dir("./testing")))
	http.ListenAndServe(":3000", nil)
}
