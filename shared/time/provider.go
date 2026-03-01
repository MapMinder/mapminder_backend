package timeProvider

import "time"

type RealTimeProvider interface {
	Now() time.Time
}

type realTimeProvider struct{}

func NewTimeProvider() RealTimeProvider {
	return &realTimeProvider{}
}

func (t *realTimeProvider) Now() time.Time {
	return time.Now()
}
