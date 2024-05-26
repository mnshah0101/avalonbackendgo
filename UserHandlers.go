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

	if user.ID == "" {
		log.Printf("User not found: %s", loginUser.Email)
		response := ErrorResponse{
			Message: "User not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to read request body",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var getUserRequest struct {
		ID string `json:"_id"`
	}

	if err := json.Unmarshal(body, &getUserRequest); err != nil {
		response := ErrorResponse{
			Message: "Failed to read request body",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err = getUserFromId(getUserRequest.ID)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get user",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	//check if user exists
	if user.ID == "" {
		response := ErrorResponse{
			Message: "User not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// return user
	response := SuccessResponse{
		Message: "User retrieved successfully",
		Status:  http.StatusOK,
		Object:  user,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

	var deleteUserRequest struct {
		ID string `json:"_id"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(body, &deleteUserRequest); err != nil {
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

	log.Printf("Request body unmarshalled successfully: %+v", deleteUserRequest)

	var user User

	user, err = deleteUserFromId(deleteUserRequest.ID, deleteUserRequest.Email)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		response := ErrorResponse{
			Message: "Failed to delete user",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	//check if user
	if user.ID == "" {
		log.Printf("User not found: %s", deleteUserRequest.Email)
		response := ErrorResponse{
			Message: "User not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("User deleted successfully: %s", deleteUserRequest.Email)

	response := SuccessResponse{
		Message: "User deleted successfully",
		Status:  http.StatusOK,
		Object:  user,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	//parse request body into user object
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

	//check to see if all fields are filled
	if user.ID == "" || user.Email == "" || user.FirstName == "" || user.LastName == "" || user.Organization == "" || user.Password == "" || user.ProfilePicture == "" {
		log.Printf("Error: All fields must be filled")
		response := ErrorResponse{
			Message: "All fields must be filled",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var return_user User

	//update user
	return_user, err = updateUser(user)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		response := ErrorResponse{
			Message: "Failed to update user",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	//check if user exists:
	if return_user.ID == "" {
		log.Printf("User not found: %s", user.Email)
		response := ErrorResponse{
			Message: "User not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	//return success
	response := SuccessResponse{
		Message: "User updated successfully",
		Status:  http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}