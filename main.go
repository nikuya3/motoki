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
)

func handleError(e error) {
	if e != nil {
		log.Print(e)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	handleError(err)
	filePath := "http://freewavesamples.com/files/Alesis-Sanctuary-QCard-AcoustcBas-C2.wav"
	log.Print(body[:10])
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

func saveVoice(data []byte) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	handleError(err)
	db.Exec("CREATE TABLE IF NOT EXISTS voices (id bigint unsigned primary key, data bytea not null)")
	db.Exec("INSERT INTO voices VALUES (?)", data)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./www/dist")))
	http.HandleFunc("/recognize", handler)
	fmt.Printf(os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
