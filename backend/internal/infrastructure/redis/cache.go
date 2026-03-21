package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

type GameCacheImpl struct {
	client *redis.Client
}

func NewGameCache(client *redis.Client) *GameCacheImpl {
	return &GameCacheImpl{client: client}
}

func playerKey(gameID, playerID string) string {
	return fmt.Sprintf("game:%s:player:%s", gameID, playerID)
}

func gameKey(gameID string) string {
	return fmt.Sprintf("game:%s:state", gameID)
}

func playersSetKey(gameID string) string {
	return fmt.Sprintf("game:%s:players", gameID)
}

func teamScoreKey(gameID string) string {
	return fmt.Sprintf("game:%s:team_scores", gameID)
}

func (c *GameCacheImpl) SetPlayerState(ctx context.Context, gameID string, state *domain.PlayerLiveState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	pipe := c.client.Pipeline()
	pipe.Set(ctx, playerKey(gameID, state.PlayerID), data, 0)
	pipe.SAdd(ctx, playersSetKey(gameID), state.PlayerID)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *GameCacheImpl) GetPlayerState(ctx context.Context, gameID, playerID string) (*domain.PlayerLiveState, error) {
	data, err := c.client.Get(ctx, playerKey(gameID, playerID)).Bytes()
	if err != nil {
		return nil, err
	}
	var state domain.PlayerLiveState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (c *GameCacheImpl) GetAllPlayerStates(ctx context.Context, gameID string) ([]*domain.PlayerLiveState, error) {
	playerIDs, err := c.client.SMembers(ctx, playersSetKey(gameID)).Result()
	if err != nil {
		return nil, err
	}
	states := make([]*domain.PlayerLiveState, 0, len(playerIDs))
	for _, pid := range playerIDs {
		state, err := c.GetPlayerState(ctx, gameID, pid)
		if err != nil {
			continue
		}
		states = append(states, state)
	}
	return states, nil
}

func (c *GameCacheImpl) SetGameState(ctx context.Context, state *domain.GameLiveState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, gameKey(state.GameID), data, 0).Err()
}

func (c *GameCacheImpl) GetGameState(ctx context.Context, gameID string) (*domain.GameLiveState, error) {
	data, err := c.client.Get(ctx, gameKey(gameID)).Bytes()
	if err != nil {
		return nil, err
	}
	var state domain.GameLiveState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (c *GameCacheImpl) IncrTeamScore(ctx context.Context, gameID, teamID string, delta int) error {
	return c.client.HIncrBy(ctx, teamScoreKey(gameID), teamID, int64(delta)).Err()
}

func (c *GameCacheImpl) SetTimeRemaining(ctx context.Context, gameID string, seconds int) error {
	key := gameKey(gameID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	var state domain.GameLiveState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	state.TimeRemainingS = seconds

	// Also get team scores from hash
	scores, _ := c.client.HGetAll(ctx, teamScoreKey(gameID)).Result()
	if state.TeamScores == nil {
		state.TeamScores = make(map[string]int)
	}
	for k, v := range scores {
		val, _ := strconv.Atoi(v)
		state.TeamScores[k] = val
	}

	return c.SetGameState(ctx, &state)
}

func (c *GameCacheImpl) DeleteGameState(ctx context.Context, gameID string) error {
	playerIDs, _ := c.client.SMembers(ctx, playersSetKey(gameID)).Result()
	pipe := c.client.Pipeline()
	pipe.Del(ctx, gameKey(gameID))
	pipe.Del(ctx, playersSetKey(gameID))
	pipe.Del(ctx, teamScoreKey(gameID))
	for _, pid := range playerIDs {
		pipe.Del(ctx, playerKey(gameID, pid))
	}
	_, err := pipe.Exec(ctx)
	return err
}
