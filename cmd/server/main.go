package main

import (
	"log"
	"golang/internal/route"
	"net/http"

	"golang/internal/database"
)

func main() {

	mongoURI := "mongodb://localhost:27017/GoDB"
	database.ConnectMongoDB(mongoURI,"GoDB")

	

	port := ":8080"
	r := route.NewRouter()
	
	log.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(port,r))

}
