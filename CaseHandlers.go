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

	log.Printf("Case: %+v", myCase)

	// check if no fields are blank
	if myCase.CaseTitle == "" || myCase.AttorneyFirstName == "" || myCase.AttorneyLastName == "" || myCase.BucketName == "" || myCase.CaseInfo == "" || myCase.CaseType == "" || myCase.City == "" || myCase.Date == "" || myCase.JudgeName == "" || myCase.SeedDoc == "" || myCase.SeedText == "" || myCase.State == "" || myCase.UserID == "" {
		response := ErrorResponse{
			Message: "All fields must be filled out",
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
		Object:  return_case,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetCaseByIDHandler(w http.ResponseWriter, r *http.Request) {
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

	var getCaseByIDRequest struct {
		ID string `json:"_id"`
	}

	if err := json.Unmarshal(body, &getCaseByIDRequest); err != nil {
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
	return_case, err = GetCaseFromId(getCaseByIDRequest.ID)
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

	// return case
	response := SuccessResponse{
		Message: "Case retrieved successfully",
		Status:  http.StatusOK,
		Object:  return_case,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func GetCaseByUserHandler(w http.ResponseWriter, r *http.Request) {
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

	var getCaseByUserRequest struct {
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal(body, &getCaseByUserRequest); err != nil {
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

	var return_user User
	return_user, err = getUserFromId(getCaseByUserRequest.UserID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		response := ErrorResponse{
			Message: "Failed to get user",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if return_user.ID == "" {
		response := ErrorResponse{
			Message: "User not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	var return_cases []Case
	return_cases, err = GetCasesByUserId(getCaseByUserRequest.UserID)
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
		Object:  return_cases,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
