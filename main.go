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
	cmd := exec.Command("Rscript", "pred/recognition.R", "pred/voices/erwin.wav")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
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
