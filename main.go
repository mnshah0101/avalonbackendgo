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
	router.HandleFunc("POST /getCase", GetCaseByIDHandler)
	router.HandleFunc("POST /getUserCases", GetCaseByUserHandler)
	router.HandleFunc("POST /deleteCaseById", DeleteCaseByIDHandler)
	router.HandleFunc("POST /deleteUserCases", DeleteCasesByUserHandler)

	// Document Routes
	router.HandleFunc("POST /getCaseDocuments", GetDocumentsByCaseHandler)
	router.HandleFunc("POST /getDocumentById", GetDocumentByIDHandler)
	router.HandleFunc("POST /deleteDocumentById", DeleteDocumentByIDHandler)
	router.HandleFunc("POST /deleteCaseDocuments", DeleteDocumentsByCaseHandler)

	// endpoints to be done
	// upload document
	// delete all docs by user

	// Chat Routes
	router.HandleFunc("POST /getCaseChat", GetChatByCaseIDHandler) // the case_id and chat id are the same

	//endpoints to be done
	// delete by chat id
	// delete by user id
	// MAKE SURE CHAT IS BEING CREATED WHEN CASE IS CREATED

	log.Println("Server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", router))

}
