package handlers

import (
	"context"
	"encoding/json" // package to encode and decode the json into struct 
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go-api-postgres/models" // models package with User model
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  //  used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

//connection with postgres db
func createConnection() *pgxpool.Pool {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	pool, err := pgxpool.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}

	//fmt.Println("Successfully connected!")
	// return the connection
	return pool
}

// CreateUser create a user in the postgres db
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode request body.  %v", err)
	}

	// call insert user function and pass the user to the table
	insertID := insertUser(user)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// GetAllUser will return all the users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	// get all users from  db
	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all users. %v", err)
	}

	// send all users as response
	json.NewEncoder(w).Encode(users)
}

// UpdateUser update user's detail in the postgres db
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	// get id from the request params, key is "id" and convert it to string
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var user models.User

	// decode the json request to user
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call update user to update the user
	updatedRows := updateUser(int64(id), user)

	// format the message string
	msg := fmt.Sprintf("User updated successfully. Total rows/record affected %v", updatedRows)

	// format of the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// DeleteUser delete user's detail in the postgres db
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// get the id from the request params, key is "id" and convert int to string
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	// call the deleteUser, convert the int to int64
	deletedRows := deleteUser(int64(id))

	// format the message string
	msg := fmt.Sprintf("User deleted successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
// insert one user to the DB
func insertUser(user models.User) int64 {

	// create connection with postgres db
	db := createConnection()


	defer db.Close()

	// create the insert sql query
	sqlStatement := `INSERT INTO users (name, lastname, age, birthdate) VALUES ($1, $2, $3, $4) RETURNING id`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(context.Background(),sqlStatement, user.Name, user.Lastname, user.Age, user.Birthdate).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	//fmt.Printf("Inserted a single record %v", id)

	return id
}

// get all users from the DB
func getAllUsers() ([]models.User, error) {
	// create connection with postgres db
	db := createConnection()

	defer db.Close()

	var users []models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users`

	// execute the sql statement
	rows, err := db.Query(context.Background(),sqlStatement)


	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var user models.User

		// unmarshal the row object to user
		err = rows.Scan(&user.ID, &user.Name, &user.Lastname, &user.Age, &user.Birthdate)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		users = append(users, user)

	}

	return users, err
}

// update user in the DB
func updateUser(id int64, user models.User) int64 {

	// create connection with postgres db
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE users SET name=$2, lastname=$3, age=$4, birthdate=$5 WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(context.Background(),sqlStatement, id, user.Name, user.Lastname, user.Age, user.Birthdate)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected := res.RowsAffected()


	//fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete user from the DB
func deleteUser(id int64) int64 {

	// create connection with postgres db
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM users WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(context.Background(),sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected := res.RowsAffected()


	//fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
