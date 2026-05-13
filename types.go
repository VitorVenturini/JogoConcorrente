package main

import (
	"math/rand"
	"time"
)

type CommandType string

const (
	CommandMoveUp    CommandType = "up"
	CommandMoveDown  CommandType = "down"
	CommandMoveLeft  CommandType = "left"
	CommandMoveRight CommandType = "right"
	CommandAttack    CommandType = "attack"
	CommandQuit      CommandType = "quit"
)

type EnemyActionType string

const (
	EnemyActionMove   EnemyActionType = "move"
	EnemyActionAttack EnemyActionType = "attack"
	EnemyActionWait   EnemyActionType = "wait"
)

type EffectKind string

const (
	EffectAttack EffectKind = "attack"
	EffectHit    EffectKind = "hit"
)

type Position struct {
	X int
	Y int
}

type Arena struct {
	Width  int
	Height int
}

type PlayerCommand struct {
	Type CommandType
}

type EnemyAction struct {
	EnemyID string
	Type    EnemyActionType
	Target  Position
}

type Tick struct {
	Number int
}

type Player struct {
	Name     string
	Position Position
	Health   int
	Symbol   rune
}

type Enemy struct {
	ID       string
	Name     string
	Position Position
	Health   int
	Symbol   rune
	Behavior string
}

type Effect struct {
	Pos   Position
	Kind  EffectKind
	Ticks int
}

type HUD struct {
	Tick       int
	Message    string
	PlayerLife int
}

type GameState struct {
	Arena      Arena
	Player     Player
	Enemies    []Enemy
	Obstacles  []Position
	Effects    []Effect
	HUD        HUD
	GameOver   bool
	Victory    bool
	ShouldQuit bool
}

type GameSnapshot struct {
	Arena     Arena
	Player    Player
	Enemies   []Enemy
	Obstacles []Position
	Effects   []Effect
	HUD       HUD
	GameOver  bool
	Victory   bool
}

func NewInitialGameState() GameState {
	arena := Arena{
		Width:  25,
		Height: 25,
	}

	player := Player{
		Name:     "Player",
		Position: Position{X: 15, Y: 15},
		Health:   10,
		Symbol:   '🙂',
	}

	enemies := []Enemy{
		{ID: "enemy-a", Name: "Enemy A", Position: Position{}, Health: 3, Symbol: '👾', Behavior: "chase"},
		{ID: "enemy-b", Name: "Enemy B", Position: Position{}, Health: 3, Symbol: '👹', Behavior: "chase"},
		{ID: "enemy-c", Name: "Enemy C", Position: Position{}, Health: 3, Symbol: '👻', Behavior: "chase"},
		{ID: "enemy-d", Name: "Enemy D", Position: Position{}, Health: 3, Symbol: '🤖', Behavior: "chase"},
		{ID: "enemy-e", Name: "Enemy E", Position: Position{}, Health: 3, Symbol: '🐍', Behavior: "chase"},
		{ID: "enemy-f", Name: "Enemy F", Position: Position{}, Health: 3, Symbol: '🦇', Behavior: "chase"},
	}

	obstacles := randomObstacles(arena, player.Position, len(enemies))

	randomizeEnemyPositions(enemies, arena, player.Position, obstacles)

	return GameState{
		Arena:     arena,
		Player:    player,
		Enemies:   enemies,
		Obstacles: obstacles,
		Effects:   nil,
		HUD: HUD{
			Tick:       0,
			Message:    "O jogo comecou!! Sobreviva. Use as Setas do teclado",
			PlayerLife: player.Health,
		},
	}
}

func randomizeEnemyPositions(enemies []Enemy, arena Arena, playerPos Position, obstacles []Position) {
	maxCells := arena.Width * arena.Height
	if len(enemies)+1 > maxCells {
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	used := map[Position]bool{
		playerPos: true,
	}
	for _, pos := range obstacles {
		used[pos] = true
	}

	for i := range enemies {
		for {
			pos := Position{X: rng.Intn(arena.Width), Y: rng.Intn(arena.Height)}
			if !used[pos] {
				used[pos] = true
				enemies[i].Position = pos
				break
			}
		}
	}
}

func randomObstacles(arena Arena, playerPos Position, enemyCount int) []Position {
	area := arena.Width * arena.Height
	maxCount := area - (1 + enemyCount)
	if maxCount <= 0 {
		return nil
	}

	count := area / 20
	if count < 20 {
		count = 20
	}
	if count > maxCount {
		count = maxCount
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	used := map[Position]bool{
		playerPos: true,
	}
	obstacles := make([]Position, 0, count)

	for len(obstacles) < count {
		pos := Position{X: rng.Intn(arena.Width), Y: rng.Intn(arena.Height)}
		if used[pos] {
			continue
		}
		used[pos] = true
		obstacles = append(obstacles, pos)
	}

	return obstacles
}

func (s GameState) Snapshot() GameSnapshot {
	enemiesCopy := make([]Enemy, len(s.Enemies))
	copy(enemiesCopy, s.Enemies)
	obstaclesCopy := make([]Position, len(s.Obstacles))
	copy(obstaclesCopy, s.Obstacles)
	effectsCopy := make([]Effect, len(s.Effects))
	copy(effectsCopy, s.Effects)

	return GameSnapshot{
		Arena:     s.Arena,
		Player:    s.Player,
		Enemies:   enemiesCopy,
		Obstacles: obstaclesCopy,
		Effects:   effectsCopy,
		HUD:       s.HUD,
		GameOver:  s.GameOver,
		Victory:   s.Victory,
	}
}
