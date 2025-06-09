package db

import (
	"context"

	"github.com/go_course_master/review_manager/internal/models"
)

type DB interface {
	GetUserByEmail(ctx context.Context, email string) (user models.User, err error)
	GetObjectByID(ctx context.Context, id uint64) (object models.Object, err error)
	CreateUser(ctx context.Context, user models.User) (uint64, error)
	GetReviews(ctx context.Context, filter models.ReviewFilter) ([]models.Review, error)
	CreateReview(ctx context.Context, review models.Review) (uint64, error)
	GetAverageRateByObjectID(ctx context.Context, objectID uint64, status models.ReviewStatus) (float64, error)
}

type RateDB interface {
	UpdateReviewRate(ctx context.Context, reviewID uint64, status models.ReviewStatus, rate int) error
	GetReviewByID(ctx context.Context, reviewID uint64) (models.Review, error)
}

type Migrator interface {
	MogrationUp(ctx context.Context) error
}
