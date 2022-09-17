package cron_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/troydai/cron"
)

func TestCron(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		repeat      int
		interval    time.Duration
		timeout     time.Duration
		assertCount func(t *testing.T, count int)
	}{
		{
			name:     "run once",
			interval: 50 * time.Millisecond,
			repeat:   1,
			timeout:  time.Second,
			assertCount: func(t *testing.T, count int) {
				assert.Equal(t, 1, count)
			},
		},
		{
			name:     "run twice",
			interval: 50 * time.Millisecond,
			repeat:   2,
			timeout:  time.Second,
			assertCount: func(t *testing.T, count int) {
				assert.Equal(t, 2, count)
			},
		},
		{
			name:     "run till timeout",
			interval: 50 * time.Millisecond,
			repeat:   9999,
			timeout:  time.Second,
			assertCount: func(t *testing.T, count int) {
				assert.GreaterOrEqual(t, count, 19)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			var counter int
			job := func(ctx context.Context) bool {
				if counter >= tc.repeat {
					return false
				}

				counter++
				return true
			}

			term, err := cron.Start(ctx, job, cron.WithInterval(tc.interval))
			assert.NoError(t, err)

			<-term
			tc.assertCount(t, counter)
		})
	}
}
