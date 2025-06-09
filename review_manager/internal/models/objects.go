package models

import "time"

type Object struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
}
