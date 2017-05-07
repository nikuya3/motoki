package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test %s", r.URL.Path[1:])
}

func main() {
	fmt.Printf("Hello web")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
