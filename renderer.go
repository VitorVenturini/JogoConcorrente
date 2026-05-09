package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

func RunRenderer(ctx context.Context, wg *sync.WaitGroup, renderCh <-chan GameSnapshot) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case snapshot := <-renderCh:
			drawSnapshot(snapshot)
		}
	}
}

func drawSnapshot(snapshot GameSnapshot) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Arena Concorrente")
	fmt.Println(renderArena(snapshot))
	fmt.Printf("HP: %d  Tick: %d\n", snapshot.HUD.PlayerLife, snapshot.HUD.Tick)
	fmt.Printf("Mensagem: %s\n", snapshot.HUD.Message)
	fmt.Println("Comandos: up, down, left, right, attack, quit")
}

func renderArena(snapshot GameSnapshot) string {
	grid := make([][]rune, snapshot.Arena.Height)
	for y := 0; y < snapshot.Arena.Height; y++ {
		grid[y] = make([]rune, snapshot.Arena.Width)
		for x := 0; x < snapshot.Arena.Width; x++ {
			grid[y][x] = '.'
		}
	}

	playerPos := snapshot.Player.Position
	if isInsideArena(snapshot.Arena, playerPos) {
		grid[playerPos.Y][playerPos.X] = snapshot.Player.Symbol
	}

	for _, enemy := range snapshot.Enemies {
		if isInsideArena(snapshot.Arena, enemy.Position) {
			grid[enemy.Position.Y][enemy.Position.X] = enemy.Symbol
		}
	}

	lines := make([]string, 0, snapshot.Arena.Height+2)
	border := "+" + strings.Repeat("-", snapshot.Arena.Width) + "+"
	lines = append(lines, border)
	for _, row := range grid {
		lines = append(lines, "|"+string(row)+"|")
	}
	lines = append(lines, border)

	return strings.Join(lines, "\n")
}

func isInsideArena(arena Arena, position Position) bool {
	return position.X >= 0 &&
		position.X < arena.Width &&
		position.Y >= 0 &&
		position.Y < arena.Height
}
