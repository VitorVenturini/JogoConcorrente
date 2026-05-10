package main

import (
	"context"
	"sync"

	"github.com/gdamore/tcell/v2"
)

func RunInputReader(ctx context.Context, wg *sync.WaitGroup, inputCh chan<- PlayerCommand, screen tcell.Screen) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// screen.PollEvent() aguarda uma tecla ser pressionada instantaneamente
		ev := screen.PollEvent()
		if ev == nil {
			return // Ocorre durante o shutdown quando a tela é finalizada
		}

		var cmd PlayerCommand
		valid := false

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				cmd = PlayerCommand{Type: CommandMoveUp}
				valid = true
			case tcell.KeyDown:
				cmd = PlayerCommand{Type: CommandMoveDown}
				valid = true
			case tcell.KeyLeft:
				cmd = PlayerCommand{Type: CommandMoveLeft}
				valid = true
			case tcell.KeyRight:
				cmd = PlayerCommand{Type: CommandMoveRight}
				valid = true
			case tcell.KeyEscape: // Botão Esc para sair
				cmd = PlayerCommand{Type: CommandQuit}
				valid = true
			case tcell.KeyRune:
				// Tecla de Espaço para atacar
				if ev.Rune() == ' ' {
					cmd = PlayerCommand{Type: CommandAttack}
					valid = true
				}
			}
		}

		if !valid {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case inputCh <- cmd:
		}

		if cmd.Type == CommandQuit {
			return
		}
	}
}
