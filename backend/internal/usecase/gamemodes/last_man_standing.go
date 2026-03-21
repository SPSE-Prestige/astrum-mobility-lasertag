package gamemodes

import "github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"

// LastManStanding — no respawns, last player/team standing wins.
type LastManStanding struct{}

func (l *LastManStanding) Mode() domain.GameMode { return domain.GameModeLastManStanding }

func (l *LastManStanding) OnHit(game *domain.Game, attacker, victim *domain.GamePlayer, damage int, isHeadshot bool) *domain.HitResult {
	scoring := game.Config.Scoring
	actualDamage := damage
	if isHeadshot {
		actualDamage = int(float64(damage) * scoring.HeadshotMultiplier)
	}
	if actualDamage > victim.HP {
		actualDamage = victim.HP
	}

	sameTeam := attacker.TeamID != nil && victim.TeamID != nil && *attacker.TeamID == *victim.TeamID
	if sameTeam && !game.Config.Player.FriendlyFire {
		return &domain.HitResult{DamageApplied: 0, AttackerScoreChange: 0, IsKill: false}
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

func (l *LastManStanding) OnKill(game *domain.Game, attacker, victim *domain.GamePlayer) *domain.KillResult {
	return &domain.KillResult{AttackerScoreChange: 0, VictimScoreChange: 0}
}

func (l *LastManStanding) CheckWinCondition(game *domain.Game, players []*domain.GamePlayer, teams []*domain.Team) *domain.WinResult {
	if game.Config.TeamCount > 0 {
		return l.checkTeamWin(players, teams)
	}
	alive := 0
	var lastAlive *domain.GamePlayer
	for _, p := range players {
		if p.IsAlive {
			alive++
			lastAlive = p
		}
	}
	if alive <= 1 && lastAlive != nil {
		return &domain.WinResult{GameOver: true, WinnerID: lastAlive.ID, Reason: "last player standing"}
	}
	return &domain.WinResult{GameOver: false}
}

func (l *LastManStanding) checkTeamWin(players []*domain.GamePlayer, teams []*domain.Team) *domain.WinResult {
	teamAlive := make(map[string]int)
	for _, p := range players {
		if p.IsAlive && p.TeamID != nil {
			teamAlive[*p.TeamID]++
		}
	}
	aliveTeams := 0
	var lastTeamID string
	for tid, count := range teamAlive {
		if count > 0 {
			aliveTeams++
			lastTeamID = tid
		}
	}
	if aliveTeams <= 1 && lastTeamID != "" {
		return &domain.WinResult{GameOver: true, WinnerID: lastTeamID, Reason: "last team standing"}
	}
	return &domain.WinResult{GameOver: false}
}

func (l *LastManStanding) CanRespawn(_ *domain.Game, _ *domain.GamePlayer) bool {
	return false // no respawns in LMS
}
