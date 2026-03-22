package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/SPSE-Prestige/aimtec2026-lasertag/backend/internal/usecase"
)

// Client wraps the Paho MQTT client for device communication.
type Client struct {
	client    pahomqtt.Client
	deviceUC  *usecase.DeviceUseCase
	hitUC     *usecase.HitUseCase
	gameUC    *usecase.GameUseCase
	broadcast func(msg any) // WebSocket broadcast callback
}

func NewClient(
	brokerURL string,
	deviceUC *usecase.DeviceUseCase,
	hitUC *usecase.HitUseCase,
	gameUC *usecase.GameUseCase,
	broadcast func(msg any),
) *Client {
	c := &Client{
		deviceUC:  deviceUC,
		hitUC:     hitUC,
		gameUC:    gameUC,
		broadcast: broadcast,
	}

	opts := pahomqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("lasertag-backend").
		SetAutoReconnect(true).
		SetOnConnectHandler(func(_ pahomqtt.Client) {
			log.Println("[MQTT] connected to broker")
			c.subscribe()
		}).
		SetConnectionLostHandler(func(_ pahomqtt.Client, err error) {
			log.Printf("[MQTT] connection lost: %v", err)
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

func (c *Client) subscribe() {
	subs := map[string]byte{
		"device/+/register":  0,
		"device/+/heartbeat": 0,
		"device/+/event":     0,
	}
	token := c.client.SubscribeMultiple(subs, c.handleMessage)
	token.Wait()
	if token.Error() != nil {
		log.Printf("[MQTT] subscribe error: %v", token.Error())
	}
	log.Println("[MQTT] subscribed to device topics")
}

func (c *Client) handleMessage(_ pahomqtt.Client, msg pahomqtt.Message) {
	topic := msg.Topic()
	parts := strings.Split(topic, "/")
	if len(parts) < 3 {
		return
	}
	deviceID := parts[1]
	action := parts[2]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch action {
	case "register":
		c.handleRegister(ctx, deviceID)
	case "heartbeat":
		c.handleHeartbeat(ctx, deviceID)
	case "event":
		c.handleEvent(ctx, deviceID, msg.Payload())
	}
}

func (c *Client) handleRegister(ctx context.Context, deviceID string) {
	device, err := c.deviceUC.Register(ctx, deviceID)
	if err != nil {
		log.Printf("[MQTT] register error for %s: %v", deviceID, err)
		return
	}
	log.Printf("[MQTT] device registered: %s", deviceID)
	c.broadcast(map[string]any{
		"type":   "device_registered",
		"device": device,
	})
}

func (c *Client) handleHeartbeat(ctx context.Context, deviceID string) {
	if err := c.deviceUC.Heartbeat(ctx, deviceID); err != nil {
		log.Printf("[MQTT] heartbeat error for %s: %v", deviceID, err)
	}
}

type hitEvent struct {
	GameID   string `json:"game_id"`
	VictimID string `json:"victim_id"` // victim device_id
}

func (c *Client) handleEvent(ctx context.Context, attackerDeviceID string, payload []byte) {
	var evt hitEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		log.Printf("[MQTT] invalid event payload from %s: %v", attackerDeviceID, err)
		return
	}

	result, err := c.hitUC.ProcessHit(ctx, evt.GameID, attackerDeviceID, evt.VictimID)
	if err != nil {
		log.Printf("[MQTT] hit rejected (%s -> %s): %v", attackerDeviceID, evt.VictimID, err)
		return
	}

	log.Printf("[MQTT] kill: %s -> %s", attackerDeviceID, evt.VictimID)

	// Notify attacker device of confirmed kill (include score+kills for CAN→MP135)
	c.SendCommand(attackerDeviceID, map[string]any{
		"action":    "kill_confirmed",
		"victim_id": evt.VictimID,
		"score":     result.AttackerScore,
		"kills":     result.AttackerKills,
	})

	// Notify victim device to die
	c.SendCommand(evt.VictimID, map[string]any{
		"action": "die",
	})

	// Schedule respawn after delay
	game, err := c.gameUC.GetGame(ctx, evt.GameID)
	if err == nil && game.Settings.RespawnDelay > 0 {
		go func() {
			time.Sleep(time.Duration(game.Settings.RespawnDelay) * time.Second)
			rctx, rcancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer rcancel()
			if err := c.hitUC.Respawn(rctx, evt.GameID, evt.VictimID); err != nil {
				log.Printf("[MQTT] respawn error for %s: %v", evt.VictimID, err)
				return
			}
			c.SendCommand(evt.VictimID, map[string]any{
				"action": "respawn",
			})
			// Broadcast respawn to dashboard
			c.broadcast(map[string]any{
				"type":      "respawn",
				"game_id":   evt.GameID,
				"device_id": evt.VictimID,
			})
		}()
	}

	// Broadcast kill to dashboard
	c.broadcast(map[string]any{
		"type":    "kill",
		"game_id": evt.GameID,
		"result":  result,
	})
}

// SendCommand publishes a command to a specific device.
func (c *Client) SendCommand(deviceID string, command any) {
	data, err := json.Marshal(command)
	if err != nil {
		return
	}
	topic := fmt.Sprintf("device/%s/command", deviceID)
	c.client.Publish(topic, 1, false, data)
}

// PublishGameStart notifies all players' devices that the game has started.
func (c *Client) PublishGameStart(deviceIDs []string, gameID string) {
	for _, did := range deviceIDs {
		c.SendCommand(did, map[string]any{
			"action":  "game_start",
			"game_id": gameID,
		})
	}
}

// PublishGameEnd notifies all players' devices that the game has ended.
func (c *Client) PublishGameEnd(deviceIDs []string) {
	for _, did := range deviceIDs {
		c.SendCommand(did, map[string]any{
			"action":  "game_end",
		})
	}
}
