package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func handleError(e error) {
	if e != nil {
		log.Print(e)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	handleError(err)
	filePath := string(body)
	log.Print(string(body))
	cmd := exec.Command("Rscript", "/app/pred/recognition.R", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	result := cmd.Run()
	if result != nil {
		fmt.Fprintf(w, out.String())
		log.Print(result)
		fmt.Fprintf(w, "Internal server error")
	}

	switch out.String() {
	case "1":
		fmt.Fprintf(w, "Female")

	case "2":
		fmt.Fprintf(w, "Male")
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./www/dist")))
	http.HandleFunc("/recognize", handler)
	fmt.Printf(os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
