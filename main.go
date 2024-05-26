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
	router.HandleFunc("POST /getUser", GetUserHandler)
	router.HandleFunc("POST /deleteUser", DeleteUserHandler)
	router.HandleFunc("POST /updateUser", UpdateUserHandler)
	

	// Case Routes
	router.HandleFunc("POST /createCase", CreateCaseHandler)
	router.HandleFunc("POST /c", CaseFromIdHandler)

	

	log.Println("Server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))

}
