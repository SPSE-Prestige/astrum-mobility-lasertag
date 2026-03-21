package domain

import "context"

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, s *Session) error
	GetByToken(ctx context.Context, token string) (*Session, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
}

type DeviceRepository interface {
	Upsert(ctx context.Context, d *Device) error
	GetByDeviceID(ctx context.Context, deviceID string) (*Device, error)
	ListAll(ctx context.Context) ([]Device, error)
	ListByStatus(ctx context.Context, status DeviceStatus) ([]Device, error)
	UpdateStatus(ctx context.Context, deviceID string, status DeviceStatus) error
	UpdateLastSeen(ctx context.Context, deviceID string) error
}

type GameRepository interface {
	Create(ctx context.Context, g *Game) error
	GetByID(ctx context.Context, id string) (*Game, error)
	GetByCode(ctx context.Context, code string) (*Game, error)
	Update(ctx context.Context, g *Game) error
	ListAll(ctx context.Context) ([]Game, error)
	ListByStatus(ctx context.Context, status GameStatus) ([]Game, error)
}

type TeamRepository interface {
	Create(ctx context.Context, t *Team) error
	GetByID(ctx context.Context, id string) (*Team, error)
	ListByGame(ctx context.Context, gameID string) ([]Team, error)
	Delete(ctx context.Context, id string) error
}

type PlayerRepository interface {
	Create(ctx context.Context, p *Player) error
	GetByID(ctx context.Context, id string) (*Player, error)
	GetByGameAndDevice(ctx context.Context, gameID, deviceID string) (*Player, error)
	ListByGame(ctx context.Context, gameID string) ([]Player, error)
	ListByTeam(ctx context.Context, teamID string) ([]Player, error)
	Update(ctx context.Context, p *Player) error
	Delete(ctx context.Context, id string) error
}

type EventRepository interface {
	Create(ctx context.Context, e *GameEvent) error
	ListByGame(ctx context.Context, gameID string) ([]GameEvent, error)
}
