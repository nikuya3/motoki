package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"
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

	fmt.Fprintf(w, out.String())
}

func handleRate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	handleError(err)
	voiceId := strings.Split(string(body), " ")[0]
	correct := strings.Split(string(body), " ")[1]
	db, err2 := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	handleError(err2)
	stmt, err3 := db.Prepare("update predictions set correct=$1 where voice_id_fk=$2")
	handleError(err3)
	res, err4 := stmt.Exec(voiceId, correct[1])
	handleError(err4)
	log.Print(res)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./www/dist")))
	http.HandleFunc("/recognize", handler)
	http.HandleFunc("/rate", handleRate)
	fmt.Printf(os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
