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

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}
	filePath := "tmp.wav"
	log.Print(body[:10])
	err2 := ioutil.WriteFile(filePath, body, 0644)
	log.Printf("File %s", err2)
	if err2 != nil {
		log.Print(err2)
	}
	cmd := exec.Command("Rscript", "/app/pred/recognition.R", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	result := cmd.Run()
	if result != nil {
		fmt.Fprintf(w, out.String())
		log.Print(result)
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
