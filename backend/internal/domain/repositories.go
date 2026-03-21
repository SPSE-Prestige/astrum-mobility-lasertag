package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type GameRepository interface {
	Create(ctx context.Context, game *Game) error
	GetByID(ctx context.Context, id string) (*Game, error)
	Update(ctx context.Context, game *Game) error
	List(ctx context.Context) ([]*Game, error)
}

type GamePlayerRepository interface {
	Create(ctx context.Context, player *GamePlayer) error
	GetByID(ctx context.Context, id string) (*GamePlayer, error)
	GetByDeviceID(ctx context.Context, gameID, deviceID string) (*GamePlayer, error)
	GetByGunID(ctx context.Context, gameID, gunID string) (*GamePlayer, error)
	GetByGameID(ctx context.Context, gameID string) ([]*GamePlayer, error)
	Update(ctx context.Context, player *GamePlayer) error
	CountByGameID(ctx context.Context, gameID string) (int, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team *Team) error
	GetByID(ctx context.Context, id string) (*Team, error)
	GetByGameID(ctx context.Context, gameID string) ([]*Team, error)
}

type WeaponRepository interface {
	Create(ctx context.Context, weapon *Weapon) error
	GetByID(ctx context.Context, id string) (*Weapon, error)
	List(ctx context.Context) ([]*Weapon, error)
}

type GameEventRepository interface {
	Create(ctx context.Context, event *GameEvent) error
	GetByGameID(ctx context.Context, gameID string) ([]*GameEvent, error)
}

type AdminSessionRepository interface {
	Create(ctx context.Context, session *AdminSession) error
	GetByToken(ctx context.Context, token string) (*AdminSession, error)
	DeleteByUserID(ctx context.Context, userID string) error
}

type GameCache interface {
	SetPlayerState(ctx context.Context, gameID string, state *PlayerLiveState) error
	GetPlayerState(ctx context.Context, gameID, playerID string) (*PlayerLiveState, error)
	GetAllPlayerStates(ctx context.Context, gameID string) ([]*PlayerLiveState, error)
	SetGameState(ctx context.Context, state *GameLiveState) error
	GetGameState(ctx context.Context, gameID string) (*GameLiveState, error)
	IncrTeamScore(ctx context.Context, gameID, teamID string, delta int) error
	SetTimeRemaining(ctx context.Context, gameID string, seconds int) error
	DeleteGameState(ctx context.Context, gameID string) error
}

type EventBus interface {
	Publish(gameID string, event WSMessage)
	Subscribe(gameID string) chan WSMessage
	Unsubscribe(gameID string, ch chan WSMessage)
}

type WSMessage struct {
	Type    string      `json:"type"`
	GameID  string      `json:"game_id"`
	Payload interface{} `json:"payload"`
}
