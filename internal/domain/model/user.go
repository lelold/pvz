package model

type User struct {
	ID       string `json:"id" format:"uuid"`
	Email    string `json:"email" format:"email"`
	Role     string `json:"role" enum:"employee,moderator"`
	Password string
}
