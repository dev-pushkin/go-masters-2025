package rate

import "context"

type Queue interface {
	Add(ctx context.Context, reviewID uint64) error
}
