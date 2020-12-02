package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	handler "github.com/mtanzim/event-driven-bookstore/bookstore-server/handler"
	persister "github.com/mtanzim/event-driven-bookstore/bookstore-server/persister"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	db, disconnectDb := persister.NewMongo()
	defer disconnectDb()
	bookHandler := handler.NewBookHandler(db)
	r := mux.NewRouter()
	port := os.Getenv("REST_PORT")
	r.HandleFunc("/api/books", bookHandler.GetBooks).Methods(http.MethodGet)
	r.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
