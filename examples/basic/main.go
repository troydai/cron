package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/troydai/cron"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	term, err := cron.Start(ctx, cron.NoExitJob(routine), cron.WithLead(3*time.Second), cron.WithInterval(time.Second))
	if err != nil {
		log.Fatalf("fail to start cron: %s", err.Error())
	}

	fmt.Println("Cron starts in 3 seconds. Its interval is 1 second.")

	select {
	case <-term:
		fmt.Println("Cron is terminated by itself.")
	case <-time.After(10 * time.Second):
		fmt.Println("Cron is signeld to terminate.")
		cancel()
	}

	select {
	case <-term:
		fmt.Println("Cron is terminated by itself.")
	case <-time.After(time.Second):
		fmt.Println("Cron did not terminate itself.")
	}
}

func routine(_ context.Context) {
	fmt.Printf("[%v] Routine output\n", time.Now())
}
