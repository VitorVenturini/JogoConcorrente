package main

import (
	"context"
	"log"
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
	defer logPanic("coordinator")

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
		return
	default:
	}

	select {
	case <-ch:
	default:
	}

	select {
	case ch <- snapshot:
	default:
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
	if isObstacleAt(state, next) {
		state.HUD.Message = "Movimento bloqueado por obstaculo"
		return
	}
	if isEnemyAt(state, next) {
		state.HUD.Message = "Movimento bloqueado por inimigo"
		return
	}

	state.Player.Position = next
}

func attackEnemies(state *GameState) {
	px := state.Player.Position.X
	py := state.Player.Position.Y

	targetIndex := -1
	for i, e := range state.Enemies {
		dx := absInt(e.Position.X - px)
		dy := absInt(e.Position.Y - py)

		// Considera qualquer casa adjacente (8 direcoes), exceto a propria casa do jogador.
		if dx <= 1 && dy <= 1 && !(dx == 0 && dy == 0) {
			targetIndex = i
			break
		}
	}

	if targetIndex >= 0 {
		state.Enemies[targetIndex].Health--
		if state.Enemies[targetIndex].Health <= 0 {
			state.Enemies = append(state.Enemies[:targetIndex], state.Enemies[targetIndex+1:]...)
		}
	}

	switch {
	case targetIndex < 0:
		state.HUD.Message = "Ataque no vazio"
	case len(state.Enemies) == 0:
		state.Victory = true
		state.HUD.Message = "Vitoria! Todos os inimigos foram derrotados!"
	default:
		state.HUD.Message = "Ataque executado!"
	}
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func isEnemyAt(state *GameState, pos Position) bool {
	for _, e := range state.Enemies {
		if e.Position == pos {
			return true
		}
	}
	return false
}

func isObstacleAt(state *GameState, pos Position) bool {
	for _, o := range state.Obstacles {
		if o == pos {
			return true
		}
	}
	return false
}

func isEnemyAtExcept(state *GameState, pos Position, enemyID string) bool {
	for _, e := range state.Enemies {
		if e.ID != enemyID && e.Position == pos {
			return true
		}
	}
	return false
}

func findFreeAdjacentCell(state *GameState, origin Position, enemyID string) (Position, bool) {
	candidates := []Position{
		{X: origin.X, Y: origin.Y - 1},
		{X: origin.X, Y: origin.Y + 1},
		{X: origin.X - 1, Y: origin.Y},
		{X: origin.X + 1, Y: origin.Y},
	}

	for _, pos := range candidates {
		if pos.X < 0 || pos.X >= state.Arena.Width || pos.Y < 0 || pos.Y >= state.Arena.Height {
			continue
		}
		if isObstacleAt(state, pos) {
			continue
		}
		if pos == state.Player.Position {
			continue
		}
		if isEnemyAtExcept(state, pos, enemyID) {
			continue
		}
		return pos, true
	}

	return Position{}, false
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

	current := state.Enemies[idx].Position
	if isEnemyAtExcept(state, current, action.EnemyID) {
		if pos, ok := findFreeAdjacentCell(state, current, action.EnemyID); ok {
			state.Enemies[idx].Position = pos
			state.HUD.Message = action.EnemyID + " se afastou"
		} else {
			log.Printf("enemy %s overlap at %v; no free adjacent", action.EnemyID, current)
		}
		return
	}

	next := action.Target
	if next.X < 0 || next.X >= state.Arena.Width || next.Y < 0 || next.Y >= state.Arena.Height {
		log.Printf("enemy %s target out of bounds: %v", action.EnemyID, next)
		return
	}
	if isObstacleAt(state, next) {
		if pos, ok := findFreeAdjacentCell(state, current, action.EnemyID); ok {
			state.Enemies[idx].Position = pos
			state.HUD.Message = action.EnemyID + " desviou"
		} else {
			log.Printf("enemy %s blocked by obstacle at %v", action.EnemyID, next)
		}
		return
	}
	if next == state.Player.Position {
		state.Player.Health--
		state.HUD.PlayerLife = state.Player.Health
		state.HUD.Message = action.EnemyID + " atacou o jogador! HP: " + strconv.Itoa(state.Player.Health)
		if state.Player.Health <= 0 {
			state.GameOver = true
			state.HUD.Message = "Game Over! Voce foi derrotado."
		}
		return
	}
	if next == current {
		return
	}
	if isEnemyAtExcept(state, next, action.EnemyID) {
		if pos, ok := findFreeAdjacentCell(state, current, action.EnemyID); ok {
			state.Enemies[idx].Position = pos
			state.HUD.Message = action.EnemyID + " desviou"
		} else {
			log.Printf("enemy %s blocked at %v; target occupied", action.EnemyID, next)
		}
		return //se ja tem alguem naquela casa, ele aborta o movimento(colisao de inimigos)
	}

	//se a casa estiver livre, atualiza a posicao
	state.Enemies[idx].Position = next
	state.HUD.Message = action.EnemyID + " moveu"
}

func handleTick(state *GameState, tick Tick) {
	state.HUD.Tick = tick.Number
}
