package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/domain"
)

// Client wraps the Paho MQTT client for device communication.
// Depends on domain port interfaces, not concrete use cases.
type Client struct {
	client    pahomqtt.Client
	deviceUC  domain.DeviceUseCasePort
	hitUC     domain.HitUseCasePort
	gameUC    domain.GameUseCasePort
	broadcast domain.WSBroadcaster
}

func NewClient(
	brokerURL string,
	deviceUC domain.DeviceUseCasePort,
	hitUC domain.HitUseCasePort,
	gameUC domain.GameUseCasePort,
	broadcast domain.WSBroadcaster,
) *Client {
	c := &Client{
		deviceUC:  deviceUC,
		hitUC:     hitUC,
		gameUC:    gameUC,
		broadcast: broadcast,
	}

	// Unique client ID per instance to allow multiple replicas
	clientID := fmt.Sprintf("lasertag-backend-%s", uuid.New().String()[:8])

	opts := pahomqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetOnConnectHandler(func(_ pahomqtt.Client) {
			slog.Info("mqtt connected to broker", "client_id", clientID)
			c.subscribe()
		}).
		SetConnectionLostHandler(func(_ pahomqtt.Client, err error) {
			slog.Error("mqtt connection lost", "error", err)
		})

	c.client = pahomqtt.NewClient(opts)
	return c
}

func (c *Client) Connect() error {
	token := c.client.Connect()
	token.Wait()
	return token.Error()
}

func (c *Client) Disconnect() {
	c.client.Disconnect(1000)
}

// Ping checks if MQTT is connected (used by health check).
func (c *Client) Ping() bool {
	return c.client.IsConnected()
}

func (c *Client) subscribe() {
	subs := map[string]byte{
		"device/+/register":  1,
		"device/+/heartbeat": 0,
		"device/+/event":     1,
	}
	token := c.client.SubscribeMultiple(subs, c.handleMessage)
	token.Wait()
	if token.Error() != nil {
		slog.Error("mqtt subscribe error", "error", token.Error())
		return
	}
	slog.Info("mqtt subscribed to device topics")
}

func (c *Client) handleMessage(_ pahomqtt.Client, msg pahomqtt.Message) {
	topic := msg.Topic()
	parts := strings.Split(topic, "/")
	if len(parts) < 3 {
		slog.Warn("mqtt invalid topic format", "topic", topic)
		return
	}
	deviceID := parts[1]
	action := parts[2]

	if deviceID == "" {
		slog.Warn("mqtt empty device ID in topic", "topic", topic)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch action {
	case "register":
		c.handleRegister(ctx, deviceID)
	case "heartbeat":
		c.handleHeartbeat(ctx, deviceID)
	case "event":
		c.handleEvent(ctx, deviceID, msg.Payload())
	default:
		slog.Warn("mqtt unknown action", "action", action, "device_id", deviceID)
	}
}

func (c *Client) handleRegister(ctx context.Context, deviceID string) {
	device, err := c.deviceUC.Register(ctx, deviceID)
	if err != nil {
		slog.Error("mqtt register failed", "device_id", deviceID, "error", err)
		return
	}
	slog.Info("mqtt device registered", "device_id", deviceID)

	// Check if device was in an active game and send state back
	c.tryReconnect(ctx, deviceID)

	c.broadcast.Broadcast(map[string]any{
		"type":   "device_registered",
		"device": device,
	})
}

func (c *Client) handleHeartbeat(ctx context.Context, deviceID string) {
	if err := c.deviceUC.Heartbeat(ctx, deviceID); err != nil {
		slog.Error("mqtt heartbeat failed", "device_id", deviceID, "error", err)
	}
}

// tryReconnect checks if a device has an active game session and re-sends game state.
func (c *Client) tryReconnect(ctx context.Context, deviceID string) {
	info, err := c.deviceUC.Reconnect(ctx, deviceID)
	if err != nil {
		slog.Error("mqtt reconnect check failed", "device_id", deviceID, "error", err)
		return
	}
	if info == nil {
		return // device not in any active game
	}

	slog.Info("mqtt sending game state to reconnected device",
		"device_id", deviceID,
		"game_id", info.Game.ID,
		"player_id", info.Player.ID,
	)

	c.PublishGameState(deviceID, info)

	c.broadcast.Broadcast(map[string]any{
		"type":      "device_reconnected",
		"game_id":   info.Game.ID,
		"device_id": deviceID,
		"player_id": info.Player.ID,
	})
}

type hitEvent struct {
	GameID   string `json:"game_id"`
	VictimID string `json:"victim_id"`
}

// deviceEvent is the generic envelope for all device event payloads.
type deviceEvent struct {
	Action   string `json:"action"`
	GameID   string `json:"game_id"`
	VictimID string `json:"victim_id"`
}

func (c *Client) handleEvent(ctx context.Context, attackerDeviceID string, payload []byte) {
	var evt deviceEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		slog.Error("mqtt invalid event payload", "device_id", attackerDeviceID, "error", err)
		return
	}

	switch evt.Action {
	case "weapon_shoot":
		c.handleShot(ctx, attackerDeviceID, &evt)
	default:
		// Legacy hit event (game_id + victim_id, no action field)
		c.handleHit(ctx, attackerDeviceID, &evt)
	}
}

func (c *Client) handleShot(ctx context.Context, deviceID string, evt *deviceEvent) {
	if evt.GameID == "" {
		slog.Warn("mqtt shot event missing game_id", "device_id", deviceID)
		return
	}

	if err := c.hitUC.RecordShot(ctx, evt.GameID, deviceID); err != nil {
		slog.Warn("mqtt shot rejected", "device_id", deviceID, "error", err)
		return
	}

	slog.Debug("mqtt shot recorded", "device_id", deviceID, "game_id", evt.GameID)
}

func (c *Client) handleHit(ctx context.Context, attackerDeviceID string, evt *deviceEvent) {
	if evt.GameID == "" || evt.VictimID == "" {
		slog.Warn("mqtt event missing required fields", "device_id", attackerDeviceID)
		return
	}

	result, err := c.hitUC.ProcessHit(ctx, evt.GameID, attackerDeviceID, evt.VictimID)
	if err != nil {
		slog.Warn("mqtt hit rejected", "attacker", attackerDeviceID, "victim", evt.VictimID, "error", err)
		return
	}

	slog.Info("mqtt kill processed", "attacker", attackerDeviceID, "victim", evt.VictimID)

	c.SendCommand(attackerDeviceID, map[string]any{
		"action":       "kill_confirmed",
		"victim_id":    evt.VictimID,
		"score":        result.AttackerScore,
		"kills":        result.AttackerKills,
		"kill_streak":  result.KillStreak,
		"weapon_level": result.WeaponLevel,
	})

	// Send weapon upgrade command if upgrade occurred
	if result.WeaponUpgraded {
		c.SendCommand(attackerDeviceID, map[string]any{
			"action":       "weapon_upgrade",
			"weapon_level": result.WeaponLevel,
			"kill_streak":  result.KillStreak,
		})
		slog.Info("mqtt weapon upgrade", "attacker", attackerDeviceID, "level", result.WeaponLevel)
	}

	c.SendCommand(evt.VictimID, map[string]any{
		"action": "die",
	})

	// Schedule respawn
	game, err := c.gameUC.GetGame(ctx, evt.GameID)
	if err == nil && game.Settings.RespawnDelay > 0 {
		go func() {
			timer := time.NewTimer(time.Duration(game.Settings.RespawnDelay) * time.Second)
			defer timer.Stop()

			<-timer.C
			rctx, rcancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer rcancel()

			if err := c.hitUC.Respawn(rctx, evt.GameID, evt.VictimID); err != nil {
				slog.Error("mqtt respawn failed", "device_id", evt.VictimID, "error", err)
				return
			}
			c.SendCommand(evt.VictimID, map[string]any{"action": "respawn"})
			c.broadcast.Broadcast(map[string]any{
				"type":      "respawn",
				"game_id":   evt.GameID,
				"device_id": evt.VictimID,
			})
		}()
	}

	broadcastPayload := map[string]any{
		"type":    "kill",
		"game_id": evt.GameID,
		"result":  result,
	}
	if result.WeaponUpgraded {
		broadcastPayload["weapon_upgrade"] = map[string]any{
			"player_id":    result.AttackerID,
			"weapon_level": result.WeaponLevel,
			"kill_streak":  result.KillStreak,
		}
	}
	c.broadcast.Broadcast(broadcastPayload)
}

// SendCommand publishes a command to a specific device.
func (c *Client) SendCommand(deviceID string, command any) {
	data, err := json.Marshal(command)
	if err != nil {
		slog.Error("mqtt marshal command error", "device_id", deviceID, "error", err)
		return
	}
	topic := fmt.Sprintf("device/%s/command", deviceID)
	token := c.client.Publish(topic, 1, false, data)
	token.Wait()
	if token.Error() != nil {
		slog.Error("mqtt publish error", "topic", topic, "error", token.Error())
	}
}

// PublishGameStart notifies all players' devices that the game has started.
func (c *Client) PublishGameStart(players []domain.Player, game domain.Game) {
	for _, did := range players {
		c.SendCommand(did.DeviceID, map[string]any{
			"action":           "game_start",
			"game_id":          game.ID,
			"team_id":          did.TeamID,
			"is_frendlyfire":   game.Settings.FriendlyFire, // false nebo true
			"type_of_the_game": game.Settings.TypeOfGame,   // 0 - deathmatch, 1 - team deathmatch
		})
	}
}

// PublishGameEnd notifies all players' devices that the game has ended.
func (c *Client) PublishGameEnd(deviceIDs []string) {
	for _, did := range deviceIDs {
		c.SendCommand(did, map[string]any{
			"action": "game_end",
		})
	}
}

// PublishGameState sends full game state to a reconnecting device.
func (c *Client) PublishGameState(deviceID string, info *domain.ReconnectInfo) {
	c.SendCommand(deviceID, map[string]any{
		"action":           "game_rejoin",
		"game_id":          info.Game.ID,
		"is_alive":         info.Player.IsAlive,
		"kills":            info.Player.Kills,
		"deaths":           info.Player.Deaths,
		"score":            info.Player.Score,
		"weapon_level":     info.Player.WeaponLevel,
		"kill_streak":      info.Player.KillStreak,
		"remaining_time":   info.RemainingTime,
		"team_id":          info.Player.TeamID,
		"is_frendlyfire":   info.Game.Settings.FriendlyFire, // false nebo true
		"type_of_the_game": info.Game.Settings.TypeOfGame,   // 0 - deathmatch, 1 - team deathmatch
	})
}
