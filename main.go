package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Nao foi possivel abrir debug.log: %v", err)
	} else {
		log.SetOutput(logFile)
		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
		defer func() {
			log.Println("Encerrando jogo")
			_ = logFile.Close()
		}()
	}
	log.Println("Iniciando jogo")

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

	//criacao de canais dinamicamente para quantos inimigos existirem
	var enemyChannels []chan GameSnapshot
	for i := 0; i < len(state.Enemies); i++ {
		enemyChannels = append(enemyChannels, make(chan GameSnapshot, 1))
	}

	var wg sync.WaitGroup
	// 4 goroutines base + 1 goroutine por inimigo
	wg.Add(4 + len(state.Enemies))

	go RunGameCoordinator(ctx, &wg, state, inputCh, enemyActionCh, tickCh, renderCh, cancel, enemyChannels)
	go RunRenderer(ctx, &wg, renderCh, screen)
	go RunInputReader(ctx, &wg, inputCh, screen)
	go RunTicker(ctx, &wg, tickCh)

	// Lancando as 6 entidades autonomas com tempos ligeiramente diferentes
	go RunEnemy(ctx, &wg, "enemy-a", "chase", 400*time.Millisecond, enemyChannels[0], enemyActionCh)
	go RunEnemy(ctx, &wg, "enemy-b", "chase", 600*time.Millisecond, enemyChannels[1], enemyActionCh)
	go RunEnemy(ctx, &wg, "enemy-c", "chase", 450*time.Millisecond, enemyChannels[2], enemyActionCh)
	go RunEnemy(ctx, &wg, "enemy-d", "chase", 700*time.Millisecond, enemyChannels[3], enemyActionCh)
	go RunEnemy(ctx, &wg, "enemy-e", "chase", 550*time.Millisecond, enemyChannels[4], enemyActionCh)
	go RunEnemy(ctx, &wg, "enemy-f", "chase", 650*time.Millisecond, enemyChannels[5], enemyActionCh)

	<-ctx.Done()
	wg.Wait()
}
