package cron

import (
	"context"
	"time"
)

type (
	Option interface {
		apply(*settings)
	}

	IntervalOption time.Duration
	LeadOption     time.Duration

	settings struct {
		interval time.Duration
		lead     time.Duration
	}
)

func NoExitJob(fn func(ctx context.Context)) Job {
	return func(ctx context.Context) bool {
		fn(ctx)
		return true
	}
}

func PlainJob(fn func()) Job {
	return func(_ context.Context) bool {
		fn()
		return true
	}
}

func WithInterval(interval time.Duration) IntervalOption {
	return IntervalOption(interval)
}

func (t IntervalOption) apply(s *settings) {
	s.interval = time.Duration(t)
}

func WithLead(interval time.Duration) LeadOption {
	return LeadOption(interval)
}

func (t LeadOption) apply(s *settings) {
	s.lead = time.Duration(t)
}
