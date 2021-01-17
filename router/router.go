package router

import (
	"github.com/gorilla/mux"
	"go-api-postgres/handlers"
)

func Router() *mux.Router {
	//create routes
	router := mux.NewRouter()

	router.HandleFunc("/users", handlers.GetAllUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/users", handlers.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}
