package models

type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	PasswordHash  string `json:"password_hash"`
	PlainPassword string `json:"plain_password"`
	Role          string `json:"role"`
}
