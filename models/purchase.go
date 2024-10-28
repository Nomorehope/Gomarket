package models

import (
	"github.com/google/uuid"
)

// import uuid
type Purchase struct {
	Product_ID  uuid.UUID `json:"pur_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
}
