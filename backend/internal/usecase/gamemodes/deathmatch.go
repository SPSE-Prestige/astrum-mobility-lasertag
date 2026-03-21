package gamemodes

import "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"

// Deathmatch — free-for-all, no teams, last score or highest score wins.
type Deathmatch struct{}

func (d *Deathmatch) Mode() domain.GameMode { return domain.GameModeDeathmatch }

func (d *Deathmatch) OnHit(game *domain.Game, attacker, victim *domain.GamePlayer, damage int, isHeadshot bool) *domain.HitResult {
	scoring := game.Config.Scoring
	actualDamage := damage
	if isHeadshot {
		actualDamage = int(float64(damage) * scoring.HeadshotMultiplier)
	}
	if actualDamage > victim.HP {
		actualDamage = victim.HP
	}
	isKill := victim.HP-actualDamage <= 0
	scoreChange := scoring.PointsPerHit
	if isKill {
		scoreChange += scoring.PointsPerKill
	}
	return &domain.HitResult{
		DamageApplied:       actualDamage,
		AttackerScoreChange: scoreChange,
		IsKill:              isKill,
	}
}

func (d *Deathmatch) OnKill(game *domain.Game, attacker, victim *domain.GamePlayer) *domain.KillResult {
	return &domain.KillResult{
		AttackerScoreChange: 0,
		VictimScoreChange:   0,
	}
}

func (d *Deathmatch) CheckWinCondition(game *domain.Game, players []*domain.GamePlayer, teams []*domain.Team) *domain.WinResult {
	alive := 0
	var lastAlive *domain.GamePlayer
	for _, p := range players {
		if p.IsAlive || p.LivesRemaining != 0 {
			alive++
			lastAlive = p
		}
	}
	if alive <= 1 && lastAlive != nil {
		return &domain.WinResult{GameOver: true, WinnerID: lastAlive.ID, Reason: "last player standing"}
	}
	return &domain.WinResult{GameOver: false}
}

func (d *Deathmatch) CanRespawn(game *domain.Game, player *domain.GamePlayer) bool {
	if player.LivesRemaining == 0 {
		return false
	}
	return true
}
