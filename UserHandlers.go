package main

import (
	"encoding/json"
	"io"
	"log"

	"net/http"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request method: %s, Request URL: %s", r.Method, r.URL)

	var user User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		response := ErrorResponse{
			Message: "Failed to read request body",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("Request body read successfully: %s", body)

	if err := json.Unmarshal(body, &user); err != nil {
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

	log.Printf("Request body unmarshalled successfully: %+v", user)

	err = createUser(user)

	log.Printf("User Created")

	if err != nil {
		log.Printf("Error creating user: %v", err)
		response := ErrorResponse{
			Message: "Failed to create user",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("User created successfully: %+v", user)

	response := SuccessResponse{
		Message: "User created successfully",
		Status:  http.StatusCreated,
		Object:  user,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func AuthorizeUserHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		response := ErrorResponse{
			Message: "Failed to read request body",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var loginUser LoginUser

	if err := json.Unmarshal(body, &loginUser); err != nil {
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

	var user User

	log.Printf("Request body unmarshalled successfully: %+v", loginUser)

	user, err = getUserFromEmail(loginUser.Email)

	log.Printf("User: %+v", user)

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

	if loginUser.Password != user.Password {
		log.Printf("Passwords do not match")
		response := ErrorResponse{
			Message: "Passwords do not match",
			Status:  http.StatusUnauthorized,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SuccessResponse{
		Message: "User authorized successfully",
		Status:  http.StatusOK,
		Object:  user,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
