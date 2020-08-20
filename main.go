package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fasmide/go-mess/api"
	"github.com/fasmide/go-mess/database"
	_ "github.com/mattn/go-adodb"
)

func main() {
	path := "FestoMES.accdb"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	log.Printf("opening %s", path)

	db := &database.Database{Path: path}
	err := db.Connect()
	if err != nil {
		log.Fatalf("unable to connect: %s", err)
	}

	api := api.New(db)
	http.Handle("/", api)

	log.Printf("listening on :8000")
	log.Fatal(
		http.ListenAndServe(":8000", nil),
	)
}
