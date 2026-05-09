package main

import (
	"bufio"
	"context"
	"os"
	"strings"
	"sync"
)

func RunInputReader(ctx context.Context, wg *sync.WaitGroup, inputCh chan<- PlayerCommand) {
	defer wg.Done()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
		}

		raw := strings.TrimSpace(strings.ToLower(scanner.Text()))
		cmd, ok := parseCommand(raw)
		if !ok {
			cmd = PlayerCommand{Type: ""}
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

func parseCommand(raw string) (PlayerCommand, bool) {
	switch raw {
	case string(CommandMoveUp):
		return PlayerCommand{Type: CommandMoveUp}, true
	case string(CommandMoveDown):
		return PlayerCommand{Type: CommandMoveDown}, true
	case string(CommandMoveLeft):
		return PlayerCommand{Type: CommandMoveLeft}, true
	case string(CommandMoveRight):
		return PlayerCommand{Type: CommandMoveRight}, true
	case string(CommandAttack):
		return PlayerCommand{Type: CommandAttack}, true
	case string(CommandQuit):
		return PlayerCommand{Type: CommandQuit}, true
	default:
		return PlayerCommand{}, false
	}
}
