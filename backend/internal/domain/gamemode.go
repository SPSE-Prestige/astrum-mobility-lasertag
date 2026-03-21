package domain

// GameModeHandler defines the strategy interface for different game modes.
// Each game mode implements its own rules for hits, kills, and win conditions.
type GameModeHandler interface {
	Mode() GameMode
	OnHit(game *Game, attacker *GamePlayer, victim *GamePlayer, damage int, isHeadshot bool) *HitResult
	OnKill(game *Game, attacker *GamePlayer, victim *GamePlayer) *KillResult
	CheckWinCondition(game *Game, players []*GamePlayer, teams []*Team) *WinResult
	CanRespawn(game *Game, player *GamePlayer) bool
}

type HitResult struct {
	DamageApplied       int
	AttackerScoreChange int
	IsKill              bool
}

type KillResult struct {
	AttackerScoreChange int
	VictimScoreChange   int
}

type WinResult struct {
	GameOver bool
	WinnerID string
	Reason   string
}
