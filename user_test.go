package main

import (
	"github.com/gorilla/mux" //used for routes
	"github.com/steinfletcher/apitest"  //used for api tests
	"go-api-postgres/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestGetAllUsers(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/users", handlers.GetAllUsers)
	ts := httptest.NewServer(r)
	defer ts.Close()
	apitest.New().
		Handler(r).
		Get("/users").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestDeleteUser(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", handlers.DeleteUser)
	ts := httptest.NewServer(r)
	defer ts.Close()
	apitest.New().
		Handler(r).
		Delete("/users/11").
		Expect(t).
		Status(http.StatusOK).
		End()
}

//func TestUpdateUser(t *testing.T) {
//	r := mux.NewRouter()
//	r.HandleFunc("/users/{id}", handlers.UpdateUser)
//	ts := httptest.NewServer(r)
//	defer ts.Close()
//	apitest.New().
//		Handler(r).
//		Put("/users/11").
//		Expect(t).
//		Body(`{"name": "Ivanaaa", "lastname": "Ivanov", "age":20, "birthdate":"12-12-2020"}`).
//		Status(http.StatusOK).
//		End()
//}
//func TestCreateUser(t *testing.T) {
//	r := mux.NewRouter()
//	r.HandleFunc("/users", handlers.CreateUser)
//	ts := httptest.NewServer(r)
//	defer ts.Close()
//	apitest.New().
//		Handler(r).
//		Post("/users").
//		Expect(t).
//		Status(http.StatusOK).
//		Body(`{"name":"Ivan","lastname":"Ivanov","age":20,"birthdate":"12-12-2020"}`).
//		End()
//}
