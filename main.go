package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

var (
	UsersTable     = "AvalonUsers"
	CasesTable     = "AvalonCases"
	DocumentsTable = "AvalonDocuments"
	ChatsTable     = "AvalonChats"
	RegionName     = "us-east-1"
	Bucket         = "avaloncasesbucket"
)

func init() {
	dynamo = InitDynamoDBTClient()
	s3Client = InitS3Client()
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
	router.HandleFunc("POST /uploadDocument", UploadDocumentHandler)
	router.HandleFunc("POST /uploadDocuments", UploadDocumentsHandler)
	router.HandleFunc("POST /getCaseDocuments", GetDocumentsByCaseHandler)
	router.HandleFunc("POST /getDocumentById", GetDocumentByIDHandler)
	router.HandleFunc("POST /deleteDocumentById", DeleteDocumentByIDHandler)
	router.HandleFunc("POST /deleteCaseDocuments", DeleteDocumentsByCaseHandler)
	router.HandleFunc("POST /createDocuments", CreateDocumentsHandler)

	router.HandleFunc("POST /getCaseChat", GetChatByCaseIDHandler) // the case_id and chat id are the same
	router.HandleFunc("POST /addMessage", AddMessageToChatHandler)

	log.Println("Server started on :8080")
	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))

}
