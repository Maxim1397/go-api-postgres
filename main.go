package main

import (
	"go-api-postgres/router"
	"log"
	"net/http"
)




func main() {
	r := router.Router()

	log.Println("Starting server on the port 8084...")
	//Starting server on the port 8084
	log.Fatal(http.ListenAndServe(":8084", r))
}

