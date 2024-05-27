package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func GetChatByCaseIDHandler(w http.ResponseWriter, r *http.Request) {
	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get chat",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// unmarshal the request body
	var getChatByCaseIdRequest struct {
		CaseID string `json:"case_id"`
	}

	if err := json.Unmarshal(body, &getChatByCaseIdRequest); err != nil {
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

	// get chat object
	var return_chat Chat
	return_chat, err = GetChatFromCaseId(getChatByCaseIdRequest.CaseID)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get chat",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// check if chat exists
	if return_chat.ID == "" {
		response := ErrorResponse{
			Message: "Chat not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(return_chat)
}

func AddMessageToChatHandler(w http.ResponseWriter, r *http.Request) {
	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := ErrorResponse{
			Message: "Failed to add message",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// unmarshal the request body
	var addMessageToChatRequest struct {
		CaseID  string  `json:"case_id"`
		Message Message `json:"message"`
	}

	if err := json.Unmarshal(body, &addMessageToChatRequest); err != nil {
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

	//make sure fields are not empty
	if addMessageToChatRequest.CaseID == "" || addMessageToChatRequest.Message.Text == "" || addMessageToChatRequest.Message.Sender == "" || addMessageToChatRequest.Message.Timestamp == "" {
		response := ErrorResponse{
			Message: "Fields cannot be empty",
			Status:  http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// get chat object
	var return_chat Chat
	return_chat, err = GetChatFromCaseId(addMessageToChatRequest.CaseID)

	if err != nil {
		response := ErrorResponse{
			Message: "Failed to get chat",
			Status:  http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("Chat returned")

	// check if chat exists
	if return_chat.ID == "" {
		response := ErrorResponse{
			Message: "Chat not found",
			Status:  http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("Chat exists")

	// add message to chat
	err = AddMessageToChat(addMessageToChatRequest.CaseID, addMessageToChatRequest.Message.Text, addMessageToChatRequest.Message.Sender, addMessageToChatRequest.Message.Timestamp)

	log.Printf("Message added")

	if err != nil {
		response := ErrorResponse{
			Message: "Failed to add message",
			Status:  http.StatusInternalServerError,
		}

		// log error
		log.Printf("Error adding message: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "Message added successfully",
		Status:  http.StatusOK,
	})
}
