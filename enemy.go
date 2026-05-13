package main

import (
	"context"
	"sync"
	"time"
)

func RunEnemy(
	ctx context.Context,
	wg *sync.WaitGroup,
	enemyID, behavior string,
	speed time.Duration, //velocidade individual dos inimigos
	stateCh <-chan GameSnapshot,
	actionCh chan<- EnemyAction,
) {
	defer wg.Done()
	defer logPanic("enemy:" + enemyID)

	var snap GameSnapshot
	hasSnap := false
	patrolDir := 1
	ticker := time.NewTicker(speed)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case s := <-stateCh:
			snap = s
			hasSnap = true
		case <-ticker.C:
			if !hasSnap || snap.GameOver || snap.Victory {
				continue
			}
			var action EnemyAction
			switch behavior {
			case "chase":
				action = chaseAction(enemyID, snap)
			case "patrol":
				action, patrolDir = patrolAction(enemyID, snap, patrolDir)
			default:
				action = EnemyAction{EnemyID: enemyID, Type: EnemyActionWait}
			}
			select {
			case actionCh <- action:
			case <-ctx.Done():
				return
			}
		}
	}
}

func chaseAction(enemyID string, snap GameSnapshot) EnemyAction {
	if next, ok := findNextStepToPlayer(enemyID, snap); ok {
		return EnemyAction{EnemyID: enemyID, Type: EnemyActionMove, Target: next}
	}

	return greedyChaseAction(enemyID, snap)
}

func greedyChaseAction(enemyID string, snap GameSnapshot) EnemyAction {
	var enemy Enemy
	for _, e := range snap.Enemies {
		if e.ID == enemyID {
			enemy = e
			break
		}
	}

	target := enemy.Position
	p := snap.Player.Position

	if p.X > enemy.Position.X {
		target.X++
	} else if p.X < enemy.Position.X {
		target.X--
	} else if p.Y > enemy.Position.Y {
		target.Y++
	} else if p.Y < enemy.Position.Y {
		target.Y--
	}

	return EnemyAction{EnemyID: enemyID, Type: EnemyActionMove, Target: target}
}

func findNextStepToPlayer(enemyID string, snap GameSnapshot) (Position, bool) {
	width := snap.Arena.Width
	height := snap.Arena.Height
	if width <= 0 || height <= 0 {
		return Position{}, false
	}

	var enemy Enemy
	found := false
	for _, e := range snap.Enemies {
		if e.ID == enemyID {
			enemy = e
			found = true
			break
		}
	}
	if !found {
		return Position{}, false
	}

	start := enemy.Position
	goal := snap.Player.Position
	if start == goal {
		return start, false
	}

	blocked := make([]bool, width*height)
	for _, o := range snap.Obstacles {
		if o.X < 0 || o.X >= width || o.Y < 0 || o.Y >= height {
			continue
		}
		blocked[o.Y*width+o.X] = true
	}
	for _, e := range snap.Enemies {
		if e.ID == enemyID {
			continue
		}
		if e.Position == goal {
			continue
		}
		idx := e.Position.Y*width + e.Position.X
		if idx >= 0 && idx < len(blocked) {
			blocked[idx] = true
		}
	}

	startIdx := start.Y*width + start.X
	goalIdx := goal.Y*width + goal.X
	if startIdx < 0 || startIdx >= len(blocked) || goalIdx < 0 || goalIdx >= len(blocked) {
		return Position{}, false
	}
	if blocked[startIdx] {
		return Position{}, false
	}

	queue := make([]int, 0, width*height)
	visited := make([]bool, width*height)
	parent := make([]int, width*height)
	for i := range parent {
		parent[i] = -1
	}

	queue = append(queue, startIdx)
	visited[startIdx] = true

	for q := 0; q < len(queue); q++ {
		cur := queue[q]
		if cur == goalIdx {
			break
		}

		x := cur % width
		y := cur / width

		neighbors := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		for _, d := range neighbors {
			nx := x + d[0]
			ny := y + d[1]
			if nx < 0 || nx >= width || ny < 0 || ny >= height {
				continue
			}
			ni := ny*width + nx
			if visited[ni] || blocked[ni] {
				continue
			}
			visited[ni] = true
			parent[ni] = cur
			queue = append(queue, ni)
		}
	}

	if !visited[goalIdx] {
		return Position{}, false
	}

	cur := goalIdx
	prev := parent[cur]
	for prev != -1 && prev != startIdx {
		cur = prev
		prev = parent[cur]
	}

	if prev == -1 {
		return Position{}, false
	}

	nx := cur % width
	ny := cur / width
	return Position{X: nx, Y: ny}, true
}

func patrolAction(enemyID string, snap GameSnapshot, dir int) (EnemyAction, int) {
	var enemy Enemy
	for _, e := range snap.Enemies {
		if e.ID == enemyID {
			enemy = e
			break
		}
	}

	nextX := enemy.Position.X + dir
	if nextX < 0 || nextX >= snap.Arena.Width {
		dir = -dir
		nextX = enemy.Position.X + dir
	}

	target := Position{X: nextX, Y: enemy.Position.Y}
	return EnemyAction{EnemyID: enemyID, Type: EnemyActionMove, Target: target}, dir
}
