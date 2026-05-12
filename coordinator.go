package main

import (
	"context"
	"strconv"
	"sync"
)

func RunGameCoordinator(
	ctx context.Context,
	wg *sync.WaitGroup,
	initialState GameState,
	inputCh <-chan PlayerCommand,
	enemyActionCh <-chan EnemyAction,
	tickCh <-chan Tick,
	renderCh chan GameSnapshot,
	cancel context.CancelFunc,
	enemyChannels []chan GameSnapshot,
) {
	defer wg.Done()

	state := initialState
	updateAllSnapshots(&state, renderCh, enemyChannels)

	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-inputCh:
			if state.GameOver || state.Victory {
				if cmd.Type == CommandQuit {
					cancel()
					return
				}
				continue
			}
			handlePlayerCommand(&state, cmd, cancel)
			updateAllSnapshots(&state, renderCh, enemyChannels)
			if state.ShouldQuit {
				cancel()
				return
			}
		case action := <-enemyActionCh:
			if state.GameOver || state.Victory {
				continue
			}
			handleEnemyAction(&state, action)
			updateAllSnapshots(&state, renderCh, enemyChannels)

		case tick := <-tickCh:
			if state.GameOver || state.Victory {
				continue
			}
			handleTick(&state, tick)
		}
	}
}

func updateAllSnapshots(state *GameState, r chan GameSnapshot, enemyChannels []chan GameSnapshot) {
	snap := state.Snapshot()
	sendSnapshot(r, snap)
	for _, ch := range enemyChannels {
		sendSnapshot(ch, snap)
	}
}

func sendSnapshot(ch chan GameSnapshot, snapshot GameSnapshot) {
	select {
	case ch <- snapshot:
	default:
		<-ch
		ch <- snapshot
	}
}

func handlePlayerCommand(state *GameState, cmd PlayerCommand, cancel context.CancelFunc) {
	switch cmd.Type {
	case CommandMoveUp:
		movePlayer(state, 0, -1)
		state.HUD.Message = "Jogador moveu para cima"
	case CommandMoveDown:
		movePlayer(state, 0, 1)
		state.HUD.Message = "Jogador moveu para baixo"
	case CommandMoveLeft:
		movePlayer(state, -1, 0)
		state.HUD.Message = "Jogador moveu para a esquerda"
	case CommandMoveRight:
		movePlayer(state, 1, 0)
		state.HUD.Message = "Jogador moveu para a direita"
	case CommandAttack:
		attackEnemies(state)
	case CommandQuit:
		state.ShouldQuit = true
		state.GameOver = true
		state.HUD.Message = "Encerrando jogo"
		cancel()
	default:
		state.HUD.Message = "Comando ignorado"
	}

	state.HUD.PlayerLife = state.Player.Health
}

func movePlayer(state *GameState, deltaX, deltaY int) {
	next := Position{
		X: state.Player.Position.X + deltaX,
		Y: state.Player.Position.Y + deltaY,
	}

	if next.X < 0 || next.X >= state.Arena.Width || next.Y < 0 || next.Y >= state.Arena.Height {
		state.HUD.Message = "Movimento bloqueado pela borda da arena"
		return
	}

	state.Player.Position = next
}

func attackEnemies(state *GameState) {
	px := state.Player.Position.X
	py := state.Player.Position.Y
	adjacent := []Position{
		{X: px, Y: py - 1},
		{X: px, Y: py + 1},
		{X: px - 1, Y: py},
		{X: px + 1, Y: py},
	}

	hit := false
	var remaining []Enemy
	for _, e := range state.Enemies {
		damaged := false
		for _, pos := range adjacent {
			if e.Position == pos {
				e.Health--
				damaged = true
				hit = true
				break
			}
		}
		if !damaged || e.Health > 0 {
			remaining = append(remaining, e)
		}
	}
	state.Enemies = remaining

	switch {
	case !hit:
		state.HUD.Message = "Ataque no vazio"
	case len(state.Enemies) == 0:
		state.Victory = true
		state.HUD.Message = "Vitoria! Todos os inimigos foram derrotados!"
	default:
		state.HUD.Message = "Ataque executado!"
	}
}

func handleEnemyAction(state *GameState, action EnemyAction) {
	if action.Type == EnemyActionWait {
		return
	}
	if action.Type != EnemyActionMove {
		return
	}

	idx := -1
	for i, e := range state.Enemies {
		if e.ID == action.EnemyID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return
	}

	next := action.Target
	if next.X < 0 || next.X >= state.Arena.Width || next.Y < 0 || next.Y >= state.Arena.Height {
		return
	}

	state.Enemies[idx].Position = next

	if next == state.Player.Position {
		state.Player.Health--
		state.HUD.PlayerLife = state.Player.Health
		state.HUD.Message = action.EnemyID + " atacou o jogador! HP: " + strconv.Itoa(state.Player.Health)
		if state.Player.Health <= 0 {
			state.GameOver = true
			state.HUD.Message = "Game Over! Voce foi derrotado."
		}
	} else {
		state.HUD.Message = action.EnemyID + " moveu"
	}
}

func handleTick(state *GameState, tick Tick) {
	state.HUD.Tick = tick.Number
}
