package moke

import (
	"context"

	"github.com/go_course_master/review_manager/pkg/reviewrate"
)

type Rate struct{}

func New() *Rate {
	return &Rate{}
}

func (r *Rate) Rate(ctx context.Context, review string) (rate reviewrate.RateLevel, err error) {
	// Simulate a rate calculation
	if review == "" {
		return 0, nil // No rate for empty reviews
	}

	// For simplicity, let's just return a fixed rate based on the length of the review
	level := len(review)%5 + 1 // Rate between 1 and 5
	return reviewrate.RateLevel(level), nil
}

func (r *Rate) Stop(ctx context.Context) error {
	return nil
}
