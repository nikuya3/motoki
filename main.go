package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func handleError(e error) {
	if e != nil {
		log.Print(e)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	handleError(err)
	filePath := strings.Replace(string(body), "https", "http", 1)
	log.Print(string(body))
	cmd := exec.Command("Rscript", "--verbose", "/app/pred/recognition.R", filePath)
	var out bytes.Buffer
	var errout bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errout
	result := cmd.Run()
	if result != nil {
		log.Print(out.String())
		log.Print(errout.String())
		log.Print(result)
		fmt.Fprintf(w, "Internal server error")
	}

	output := strings.Split(out.String(), " ")
	switch output[0] {
	case "1":
		fmt.Fprintf(w, "Female")

	case "2":
		fmt.Fprintf(w, "Male")
	}

	fmt.Fprintf(w, output[1])
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./www/dist")))
	http.HandleFunc("/recognize", handler)
	fmt.Printf(os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
