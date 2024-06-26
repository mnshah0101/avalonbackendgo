package main

type User struct {
	ID             string   `json:"_id"`
	Email          string   `json:"email"`
	Cases          []string `json:"cases"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Organization   string   `json:"organization"`
	Password       string   `json:"password"`
	ProfilePicture string   `json:"profile_picture"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Case struct {
	ID                string `json:"_id"`
	CaseTitle         string `json:"case_title"`
	AttorneyFirstName string `json:"attorney_first_name"`
	AttorneyLastName  string `json:"attorney_last_name"`
	CaseInfo          string `json:"case_info"`
	CaseType          string `json:"case_type"`
	City              string `json:"city"`
	Date              string `json:"date"`
	JudgeName         string `json:"judge_name"`
	NumberFiles       int    `json:"number_files"`
	State             string `json:"state"`
	UserID            string `json:"user_id"`
}

type Document struct {
	ID        string  `json:"_id"`
	FileName  string  `json:"file_name"`
	CaseID    string  `json:"case"`
	Date      string  `json:"date"`
	FileURL   string  `json:"file_url"`
	Relevancy float64 `json:"relevancy"`
	Stored    bool    `json:"stored"`
}

type Chat struct {
	ID           string    `json:"_id"`
	Messages     []Message `json:"messages"`
	SelectedDocs []string  `json:"selected_docs"`
	UserID       string    `json:"user_id"`
}

type Message struct {
	Text      string `json:"text"`
	Sender    string `json:"sender"`
	Timestamp string `json:"timestamp"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Object  interface{} `json:"object"`
}
