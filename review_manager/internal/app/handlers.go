package app

import (
	"context"
	"errors"
	"strconv"

	"github.com/go_course_master/review_manager/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (a *App) ReviewAdd(ctx context.Context, r models.ReviewAddRequest) error {
	tr := otel.Tracer("rate-review-app")
	ctx, span := tr.Start(ctx, "ReviewAdd")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.email", r.UserEmail),
		attribute.String("object.id", strconv.Itoa(int(r.ObjectID))),
	)

	// 1. Найти пользователя по email
	user, err := a.db.GetUserByEmail(ctx, r.UserEmail)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "user not found")
		return errors.Join(ErrUserNotFound, err)
	}

	object, err := a.db.GetObjectByID(ctx, r.ObjectID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "object not found")
		return errors.Join(ErrObjectNotFound, err)
	}

	newReview := models.Review{
		Status:   models.ReviewStatusCreated,
		UserID:   user.ID,
		ObjectID: object.ID,
		Text:     r.Text,
	}

	// 2. Добавить новый отзыв в базу данных
	newReview.ID, err = a.db.CreateReview(ctx, newReview)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create review")
		return err
	}

	// 3. Положить в очередь на обработку отзыва для расчета рейтинга
	a.rateQueue.Add(ctx, newReview.ID)

	span.SetStatus(codes.Ok, "review added successfully")
	span.SetAttributes(attribute.String("review.id", strconv.Itoa(int(r.ObjectID))))

	return nil
}

func (a *App) ReviewList(ctx context.Context, r models.ReviewFilter) ([]models.Review, error) {
	// This is a stub implementation. Replace with actual logic.
	return a.db.GetReviews(ctx, r)
}

func (a *App) ObjectRate(ctx context.Context, objectID uint64) (float64, error) {
	object, err := a.db.GetObjectByID(ctx, objectID)
	if err != nil {
		return 0, errors.Join(ErrObjectNotFound, err)
	}

	return a.db.GetAverageRateByObjectID(ctx, object.ID, models.ReviewStatusSuccess)
}
