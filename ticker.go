package main

import (
	"context"
	"sync"
	"time"
)

func RunTicker(ctx context.Context, wg *sync.WaitGroup, tickCh chan<- Tick) {
	defer wg.Done()
	defer logPanic("ticker")
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			count++
			select {
			case tickCh <- Tick{Number: count}:
			default:
			}
		}
	}
}
