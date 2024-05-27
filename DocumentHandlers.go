package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetDocumentsByCaseHandler(w http.ResponseWriter, r *http.Request) {
	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get documents",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// unmarshal the request body
	var getDocsByCaseIdRequest struct {
		CaseID string `json:"case_id"`
	}

	if err := json.Unmarshal(body, &getDocsByCaseIdRequest); err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		response := ErrorResponse{
			Message: "Failed to read request body",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// get case object
	var return_case Case
	return_case, err = GetCaseFromId(getDocsByCaseIdRequest.CaseID)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get case",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// check if case exists
	if return_case.ID == "" {
		response := ErrorResponse{
			Message: "Case not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// get documents by case id
	var documents []Document
	documents, err = GetDocumentsByCaseId(getDocsByCaseIdRequest.CaseID)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get documents",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SuccessResponse{
		Message: "Documents retrieved successfully",
		Status:  http.StatusOK,
		Object:  documents,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func GetDocumentByIDHandler(w http.ResponseWriter, r *http.Request) {
	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get document",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// unmarshal the request body
	var getDocByIdRequest struct {
		ID string `json:"_id"`
	}

	if err := json.Unmarshal(body, &getDocByIdRequest); err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		response := ErrorResponse{
			Message: "Failed to read request body",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// get document by id
	var document Document
	document, err = GetDocumentById(getDocByIdRequest.ID)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get document",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// check if document exists
	if document.ID == "" {
		response := ErrorResponse{
			Message: "Document not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// return document
	response := SuccessResponse{
		Message: "Document retrieved successfully",
		Status:  http.StatusOK,
		Object:  document,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func DeleteDocumentByIDHandler(w http.ResponseWriter, r *http.Request) {
	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to delete document",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// unmarshal the request body
	var deleteDocByIdRequest struct {
		ID string `json:"_id"`
	}

	if err := json.Unmarshal(body, &deleteDocByIdRequest); err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		response := ErrorResponse{
			Message: "Failed to read request body",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// get document by id
	var document Document
	document, err = GetDocumentById(deleteDocByIdRequest.ID)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get document",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// check if document exists
	if document.ID == "" {
		response := ErrorResponse{
			Message: "Document not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// delete document by id
	err = DeleteDocumentById(deleteDocByIdRequest.ID)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to delete document",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SuccessResponse{
		Message: "Document deleted successfully",
		Status:  http.StatusOK,
		Object:  document,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func DeleteDocumentsByCaseHandler(w http.ResponseWriter, r *http.Request) {
	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to delete documents",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// unmarshal the request body
	var deleteDocsByCaseRequest struct {
		CaseID string `json:"case_id"`
	}

	if err := json.Unmarshal(body, &deleteDocsByCaseRequest); err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		response := ErrorResponse{
			Message: "Failed to read request body",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// get case from caseid
	var return_case Case
	return_case, err = GetCaseFromId(deleteDocsByCaseRequest.CaseID)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get case",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// check if case exists
	if return_case.ID == "" {
		response := ErrorResponse{
			Message: "Case not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// delete documents by case id
	var documents []Document
	documents, err = DeleteDocumentsByCaseId(deleteDocsByCaseRequest.CaseID)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to delete documents",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SuccessResponse{
		Message: "Documents deleted successfully",
		Status:  http.StatusOK,
		Object:  documents,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
