package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// 1. Inicialização do tcell
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Erro ao criar tela: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Erro ao inicializar tela: %v", err)
	}
	// Isso garante que o terminal volte ao normal (sem bugar) quando o jogo fechar
	defer screen.Fini()

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

	// 2. Passamos o "screen" para quem precisa (Renderer e InputReader)
	go RunGameCoordinator(ctx, &wg, state, inputCh, enemyActionCh, tickCh, renderCh, cancel, enemyAStateCh, enemyBStateCh)
	go RunRenderer(ctx, &wg, renderCh, screen)
	go RunInputReader(ctx, &wg, inputCh, screen)
	go RunTicker(ctx, &wg, tickCh)

	// Como a movimentação agora é instantânea com o tcell, podemos deixar os inimigos rápidos de novo!
	go RunEnemy(ctx, &wg, "enemy-a", "chase", 800*time.Millisecond, enemyAStateCh, enemyActionCh)
	go RunEnemy(ctx, &wg, "enemy-b", "patrol", 1200*time.Millisecond, enemyBStateCh, enemyActionCh)

	<-ctx.Done()
	wg.Wait()
}
