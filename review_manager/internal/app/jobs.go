package app

import (
	"context"

	"github.com/go_course_master/review_manager/internal/models"
)

func (a *App) startJobs() error {

	err := a.cron.AddFunc(a.config.Jobs.RateReviewInterval, a.retryErrorRateReviewJob)
	if err != nil {
		a.log.Error("failed to add job for retrying error rate reviews", "error", err)
	}

	a.cron.Start()

	return nil
}

func (a *App) retryErrorRateReviewJob() {
	status := models.ReviewStatusRateError
	jobLogger := a.log.WithGroup("retryErrorRateReviewJob")
	reviews, err := a.db.GetReviews(context.Background(), models.ReviewFilter{
		Status: &status})
	if err != nil {
		jobLogger.Error("get reviews", "error", err)
		return
	}

	jobLogger.Info("found reviews to retry", "count", len(reviews))
	for _, review := range reviews {
		jobLogger.Info("retrying review", "review_id", review.ID)
		err := a.rateQueue.Add(context.Background(), review.ID)
		if err != nil {
			jobLogger.Error("failed to add review to queue", "review_id", review.ID, "error", err)
			continue
		}
		jobLogger.Info("review added to queue", "review_id", review.ID)
	}

}
