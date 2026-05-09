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
	stateCh <-chan GameSnapshot,
	actionCh chan<- EnemyAction,
	interval time.Duration,
) {
	defer wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var snap GameSnapshot
	hasSnap := false
	patrolDir := 1

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
