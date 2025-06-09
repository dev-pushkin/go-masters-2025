package reviewrate

import "context"

type RateLevel int

const (
	RateLevel_1 RateLevel = iota + 1
	RateLevel_2
	RateLevel_3
	RateLevel_4
	RateLevel_5
)

type ReviewRate interface {
	Rate(ctx context.Context, review string) (rate RateLevel, err error)
}
