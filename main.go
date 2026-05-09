package main

import (
	"context"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	state := NewInitialGameState()

	inputCh := make(chan PlayerCommand)
	enemyActionCh := make(chan EnemyAction)
	tickCh := make(chan Tick)
	renderCh := make(chan GameSnapshot, 1)
	enemyAStateCh := make(chan GameSnapshot, 1)
	enemyBStateCh := make(chan GameSnapshot, 1)

	var wg sync.WaitGroup
	wg.Add(6)

	go RunGameCoordinator(ctx, &wg, state, inputCh, enemyActionCh, tickCh, renderCh, cancel, enemyAStateCh, enemyBStateCh)
	go RunRenderer(ctx, &wg, renderCh)
	go RunInputReader(ctx, &wg, inputCh)
	go RunTicker(ctx, &wg, tickCh)
	go RunEnemy(ctx, &wg, "enemy-a", "chase", enemyAStateCh, enemyActionCh, 800*time.Millisecond)
	go RunEnemy(ctx, &wg, "enemy-b", "patrol", enemyBStateCh, enemyActionCh, 1200*time.Millisecond)

	<-ctx.Done()
	wg.Wait()
}
