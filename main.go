package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
	cmd := exec.Command("Rscript", "/app/pred/recognition.R", "/app/pred/voices/erwin.wav")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "Internal server error")
	}
	fmt.Fprintf(w, out.String())
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./www/dist")))
	http.HandleFunc("/recognize", handler)
	fmt.Printf(os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
