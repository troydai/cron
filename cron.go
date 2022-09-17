package cron

import (
	"context"
	"fmt"
	"time"
)

type (
	Job func(context.Context) bool

	TermCh <-chan struct{}
)

func Start(ctx context.Context, job Job, opt ...Option) (TermCh, error) {
	s := &settings{}
	for _, o := range opt {
		o.apply(s)
	}

	if s.interval == 0 {
		return nil, fmt.Errorf("invalid interval: it can't be zero")
	}

	term := make(chan struct{})

	go func() {
		defer close(term)

		wait(ctx, s.lead)

		repeat(ctx, s.interval, job)
	}()

	return term, nil
}

func wait(ctx context.Context, period time.Duration) {
	if period == 0 {
		return
	}

	timer := time.NewTimer(period)

	select {
	case <-ctx.Done():
		stopAndDrain(timer)
	case <-timer.C:
	}
}

func repeat(ctx context.Context, internval time.Duration, job Job) {
	jobCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(internval)

	for {
		select {
		case <-jobCtx.Done():
			return
		case <-ticker.C:
			go func() {
				if cont := job(jobCtx); !cont {
					cancel()
				}
			}()
		}
	}
}

func stopAndDrain(t *time.Timer) {
	fmt.Println("stop and drain")
	if t == nil {
		return
	}

	if !t.Stop() {
		<-t.C
	}
}
