package models

import "time"

type Review struct {
	ID        uint64       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	Status    ReviewStatus `json:"status"`
	Text      string       `json:"text"`
	RateLevel int          `json:"rate_level"`
	UserID    uint64       `json:"user_id"`
	ObjectID  uint64       `json:"object_id"`
}

type ReviewStatus string

const (
	ReviewStatusCreated   ReviewStatus = "CREATED"
	ReviewStatusRateError ReviewStatus = "RATE_ERROR"
	ReviewStatusSuccess   ReviewStatus = "SUCCESS"
)

type ReviewAddRequest struct {
	ObjectID  uint64 `json:"objectID"`
	UserEmail string `json:"userEmail"`
	Text      string `json:"text"`
}

type ReviewFilter struct {
	UserID    *uint64
	ObjectID  *uint64
	RateLevel *int
	Status    *ReviewStatus
}
