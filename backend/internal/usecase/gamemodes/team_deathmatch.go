package gamemodes

import "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"

// TeamDeathmatch — two or more teams, team with highest score wins.
type TeamDeathmatch struct{}

func (t *TeamDeathmatch) Mode() domain.GameMode { return domain.GameModeTeamDeathmatch }

func (t *TeamDeathmatch) OnHit(game *domain.Game, attacker, victim *domain.GamePlayer, damage int, isHeadshot bool) *domain.HitResult {
	scoring := game.Config.Scoring
	actualDamage := damage
	if isHeadshot {
		actualDamage = int(float64(damage) * scoring.HeadshotMultiplier)
	}
	if actualDamage > victim.HP {
		actualDamage = victim.HP
	}

	// Friendly fire check
	sameTeam := attacker.TeamID != nil && victim.TeamID != nil && *attacker.TeamID == *victim.TeamID
	if sameTeam && !game.Config.Player.FriendlyFire {
		return &domain.HitResult{DamageApplied: 0, AttackerScoreChange: 0, IsKill: false}
	}

	isKill := victim.HP-actualDamage <= 0
	scoreChange := scoring.PointsPerHit
	if sameTeam {
		scoreChange = -scoring.TeamkillPenalty
	}
	if isKill {
		if sameTeam {
			scoreChange -= scoring.TeamkillPenalty
		} else {
			scoreChange += scoring.PointsPerKill
		}
	}
	return &domain.HitResult{
		DamageApplied:       actualDamage,
		AttackerScoreChange: scoreChange,
		IsKill:              isKill,
	}
}

func (t *TeamDeathmatch) OnKill(game *domain.Game, attacker, victim *domain.GamePlayer) *domain.KillResult {
	return &domain.KillResult{AttackerScoreChange: 0, VictimScoreChange: 0}
}

func (t *TeamDeathmatch) CheckWinCondition(game *domain.Game, players []*domain.GamePlayer, teams []*domain.Team) *domain.WinResult {
	// Team deathmatch ends by timer; no automatic win condition during play
	return &domain.WinResult{GameOver: false}
}

func (t *TeamDeathmatch) CanRespawn(game *domain.Game, player *domain.GamePlayer) bool {
	if player.LivesRemaining == 0 {
		return false
	}
	return true
}
