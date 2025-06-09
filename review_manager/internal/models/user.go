package models

import "time"

type User struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
}
