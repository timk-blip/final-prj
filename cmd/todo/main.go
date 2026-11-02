package main

import (
	"final-prj/pkg/api"
	"final-prj/pkg/db"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	dataBase, err := db.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer dataBase.Close()
	fmt.Println("start db -", dbFile)
	r := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("web"))
	r.Handle("/*", fileServer)
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	fmt.Println("сервер запущен: http://127.0.0.1:" + port)

	api.Init(r)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
