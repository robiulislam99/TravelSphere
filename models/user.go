package models

// User represents a logged-in session user.
type User struct {
    Username  string `json:"username"`
    FirstName string `json:"first_name"`
}