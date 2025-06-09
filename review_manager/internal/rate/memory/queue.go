package memory

import (
	"context"
	"log/slog"
	"sync"

	"github.com/go_course_master/review_manager/internal/db"
	"github.com/go_course_master/review_manager/internal/models"
	"github.com/go_course_master/review_manager/pkg/reviewrate"
)

type Queue struct {
	in chan uint64
	wg sync.WaitGroup

	r  reviewrate.ReviewRate
	db db.RateDB

	log *slog.Logger
}

func NewQueue(log *slog.Logger, r reviewrate.ReviewRate, db db.RateDB) *Queue {
	return &Queue{
		in:  make(chan uint64),
		r:   r,
		db:  db,
		log: log,
	}
}

func (q *Queue) Add(_ context.Context, reviewID uint64) error {
	q.log.Debug("Adding review to queue", "reviewID", reviewID)
	q.in <- reviewID
	return nil
}

func (q *Queue) Run() {
	q.log.Info("Starting rate queue workers")
	go func() {
		for id := range q.in {
			q.log.Debug("Received review ID from queue", "reviewID", id)
			q.wg.Add(1)
			go func(id uint64) {
				defer q.wg.Done()
				q.handle(id)
			}(id)
		}
	}()
}

func (q *Queue) Stop(ctx context.Context) {
	q.log.Info("Stopping rate queue workers")
	close(q.in)
	q.wg.Wait()
}

func (q *Queue) handle(reviewID uint64) {
	q.log.Debug("handling review", "reviewID", reviewID)

	review, err := q.db.GetReviewByID(context.Background(), reviewID)
	if err != nil {
		q.log.Error("get review by ID", "reviewID", reviewID, "error", err)
		return
	}

	rate, err := q.r.Rate(context.Background(), review.Text)
	switch {
	case err != nil:
		q.log.Error("calculate review rate", "reviewID", reviewID, "error", err)
		err = q.db.UpdateReviewRate(context.Background(), reviewID, models.ReviewStatusRateError, 0)
		if err != nil {
			q.log.Error("update review rate", "reviewID", reviewID, "error", err)
		}
		return

	default:
		q.log.Debug("review rate calculated", "reviewID", reviewID, "rate", rate)
		err = q.db.UpdateReviewRate(context.Background(), reviewID, models.ReviewStatusSuccess, int(rate))
		if err != nil {
			q.log.Error("update review rate", "reviewID", reviewID, "error", err)
			return
		}
	}
}
