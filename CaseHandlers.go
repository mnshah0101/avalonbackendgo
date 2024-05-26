package main

import (
	"encoding/json"
	"io"
	"log"

	"net/http"
)

func CreateCaseHandler(w http.ResponseWriter, r *http.Request) {
	var myCase Case
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to create case",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := json.Unmarshal(body, &myCase); err != nil {
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

	var return_case Case
	return_case, err = CreateCase(myCase)
	if err != nil {
		log.Printf("Error creating case: %v", err)
		response := ErrorResponse{
			Message: "Failed to create case",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// return case
	response := SuccessResponse{
		Message: "Case created successfully",
		Status:  http.StatusCreated,
		Object:    return_case,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func CaseFromIdHandler(w http.ResponseWriter, r *http.Request){
	var myCase Case
	body, err := io.ReadAll(r.Body)
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

	if err := json.Unmarshal(body, &myCase); err != nil {
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

	var return_case Case
	return_case, err = GetCaseFromId(myCase.ID)
	if err != nil {
		log.Printf("Error getting case: %v", err)
		response := ErrorResponse{
			Message: "Failed to get case",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SuccessResponse{
		Message: "Case retrieved successfully",
		Status:  http.StatusOK,
		Object:    return_case,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}