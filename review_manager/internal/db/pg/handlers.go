package pg

import (
	"context"
	"fmt"

	"github.com/go_course_master/review_manager/internal/models"
)

// DB interface

func (p *PostgresDB) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := p.p.QueryRow(ctx, `
		SELECT id, created_at, email
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.CreatedAt, &user.Email)
	return user, err
}

func (p *PostgresDB) GetObjectByID(ctx context.Context, id uint64) (models.Object, error) {
	var obj models.Object
	err := p.p.QueryRow(ctx, `
		SELECT id, created_at, name
		FROM objects
		WHERE id = $1
	`, id).Scan(&obj.ID, &obj.CreatedAt, &obj.Name)
	return obj, err
}

func (p *PostgresDB) CreateUser(ctx context.Context, user models.User) (uint64, error) {
	var id uint64
	err := p.p.QueryRow(ctx, `
		INSERT INTO users (created_at, email)
		VALUES ( $2, $3)
		RETURNING id
	`, user.CreatedAt, user.Email).Scan(&id)
	return id, err
}

func (p *PostgresDB) GetReviews(ctx context.Context, filter models.ReviewFilter) ([]models.Review, error) {
	query := `SELECT id, created_at, status, text, rate_level, user_id, object_id FROM reviews WHERE true`
	args := []any{}
	i := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", i)
		args = append(args, *filter.UserID)
		i++
	}
	if filter.ObjectID != nil {
		query += fmt.Sprintf(" AND object_id = $%d", i)
		args = append(args, *filter.ObjectID)
		i++
	}
	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", i)
		args = append(args, *filter.Status)
		i++
	}

	rows, err := p.p.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Review
	for rows.Next() {
		var r models.Review
		if err := rows.Scan(&r.ID, &r.CreatedAt, &r.Status, &r.Text, &r.RateLevel, &r.UserID, &r.ObjectID); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

func (p *PostgresDB) CreateReview(ctx context.Context, review models.Review) (uint64, error) {
	var id uint64
	err := p.p.QueryRow(ctx, `
		INSERT INTO reviews (created_at, status, text, rate_level, user_id, object_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, review.CreatedAt, review.Status, review.Text, review.RateLevel, review.UserID, review.ObjectID).Scan(&id)
	return id, err
}

// RateDB interface

func (p *PostgresDB) UpdateReviewRate(ctx context.Context, reviewID uint64, status models.ReviewStatus, rate int) error {
	cmd, err := p.p.Exec(ctx, `
		UPDATE reviews
		SET status = $1, rate_level = $2
		WHERE id = $3
	`, status, rate, reviewID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("no review found with id %d", reviewID)
	}
	return nil
}

func (p *PostgresDB) GetReviewByID(ctx context.Context, reviewID uint64) (models.Review, error) {
	var r models.Review
	err := p.p.QueryRow(ctx, `
		SELECT id, created_at, status, text, rate_level, user_id, object_id
		FROM reviews
		WHERE id = $1
	`, reviewID).Scan(&r.ID, &r.CreatedAt, &r.Status, &r.Text, &r.RateLevel, &r.UserID, &r.ObjectID)
	return r, err
}

func (p *PostgresDB) GetAverageRateByObjectID(ctx context.Context, objectID uint64, status models.ReviewStatus) (float64, error) {
	var avgRate float64
	err := p.p.QueryRow(ctx, `
		SELECT AVG(rate_level)
		FROM reviews
		WHERE object_id = $1 AND status = $2
	`, objectID, status).Scan(&avgRate)
	if err != nil {
		return 0, err
	}
	return avgRate, nil
}
