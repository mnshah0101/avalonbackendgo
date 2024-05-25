package main

import (
	"log"

	"net/http"
)

var (
	UsersTable     = "AvalonUsers"
	CasesTable     = "AvalonCases"
	DocumentsTable = "AvalonDocuments"
	ChatsTable     = "AvalonChats"
	RegionName     = "us-east-1"
)

func init() {
	dynamo = InitDynamoDBTClient()
}

func main() {

	router := http.NewServeMux()

	// User Routes
	router.HandleFunc("POST /createUser", CreateUserHandler)

	router.HandleFunc("POST /login", AuthorizeUserHandler)

	// Case Routes

	// Document Routes

	// Chat Routes

	log.Println("Server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))

}
