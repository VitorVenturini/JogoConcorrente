package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
)

func RunRenderer(ctx context.Context, wg *sync.WaitGroup, renderCh <-chan GameSnapshot, screen tcell.Screen) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case snapshot := <-renderCh:
			drawSnapshot(snapshot, screen)
		}
	}
}

func drawSnapshot(snapshot GameSnapshot, screen tcell.Screen) {
	screen.Clear() // Limpa o buffer do tcell

	// Estilos visuais opcionais
	defStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	playerStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	enemyStyle := tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	borderStyle := tcell.StyleDefault.Foreground(tcell.ColorGray)

	offsetX := 2
	offsetY := 2

	// Desenhar o cenário
	for y := 0; y < snapshot.Arena.Height; y++ {
		for x := 0; x < snapshot.Arena.Width; x++ {
			screen.SetContent(offsetX+x, offsetY+y, '.', nil, defStyle)
		}
	}

	// Desenhar bordas
	for x := 0; x < snapshot.Arena.Width; x++ {
		screen.SetContent(offsetX+x, offsetY-1, '-', nil, borderStyle)
		screen.SetContent(offsetX+x, offsetY+snapshot.Arena.Height, '-', nil, borderStyle)
	}
	for y := 0; y < snapshot.Arena.Height; y++ {
		screen.SetContent(offsetX-1, offsetY+y, '|', nil, borderStyle)
		screen.SetContent(offsetX+snapshot.Arena.Width, offsetY+y, '|', nil, borderStyle)
	}

	// Desenhar o jogador
	pPos := snapshot.Player.Position
	screen.SetContent(offsetX+pPos.X, offsetY+pPos.Y, snapshot.Player.Symbol, nil, playerStyle)

	// Desenhar os inimigos
	for _, enemy := range snapshot.Enemies {
		ePos := enemy.Position
		screen.SetContent(offsetX+ePos.X, offsetY+ePos.Y, enemy.Symbol, nil, enemyStyle)
	}

	// Desenhar a HUD abaixo da arena
	hudY := offsetY + snapshot.Arena.Height + 1
	drawString(screen, offsetX, hudY, fmt.Sprintf("HP: %d  Tick: %d", snapshot.HUD.PlayerLife, snapshot.HUD.Tick), defStyle)
	drawString(screen, offsetX, hudY+1, fmt.Sprintf("Mensagem: %s", snapshot.HUD.Message), tcell.StyleDefault.Foreground(tcell.ColorYellow))
	drawString(screen, offsetX, hudY+3, "Comandos: Setas (mover), Espaço (atacar), Esc (sair)", borderStyle)

	screen.Show() // Atualiza a tela de fato
}

// Função auxiliar para imprimir strings longas no tcell
func drawString(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}
