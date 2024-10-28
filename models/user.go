package models

import (
	"github.com/google/uuid"
)

// import uuid
type User struct {
	User_ID     uuid.UUID `json:"u_id"`
	Username        string    `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"` // $TODO хэш
}
