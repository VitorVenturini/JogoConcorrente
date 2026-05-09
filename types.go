package main

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

type HUD struct {
	Tick       int
	Message    string
	PlayerLife int
}

type GameState struct {
	Arena      Arena
	Player     Player
	Enemies    []Enemy
	HUD        HUD
	GameOver   bool
	Victory    bool
	ShouldQuit bool
}

type GameSnapshot struct {
	Arena    Arena
	Player   Player
	Enemies  []Enemy
	HUD      HUD
	GameOver bool
	Victory  bool
}

func NewInitialGameState() GameState {
	player := Player{
		Name:     "Player",
		Position: Position{X: 5, Y: 5},
		Health:   10,
		Symbol:   '@',
	}

	enemies := []Enemy{
		{
			ID:       "enemy-a",
			Name:     "Enemy A",
			Position: Position{X: 1, Y: 1},
			Health:   3,
			Symbol:   'A',
			Behavior: "chase",
		},
		{
			ID:       "enemy-b",
			Name:     "Enemy B",
			Position: Position{X: 8, Y: 8},
			Health:   3,
			Symbol:   'B',
			Behavior: "patrol",
		},
	}

	return GameState{
		Arena: Arena{
			Width:  10,
			Height: 10,
		},
		Player:  player,
		Enemies: enemies,
		HUD: HUD{
			Tick:       0,
			Message:    "Digite up, down, left, right, attack ou quit",
			PlayerLife: player.Health,
		},
	}
}

func (s GameState) Snapshot() GameSnapshot {
	enemiesCopy := make([]Enemy, len(s.Enemies))
	copy(enemiesCopy, s.Enemies)

	return GameSnapshot{
		Arena:    s.Arena,
		Player:   s.Player,
		Enemies:  enemiesCopy,
		HUD:      s.HUD,
		GameOver: s.GameOver,
		Victory:  s.Victory,
	}
}
