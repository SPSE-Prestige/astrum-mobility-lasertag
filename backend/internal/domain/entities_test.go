package domain

import (
	"errors"
	"testing"
)

func TestGameSettings_Validate(t *testing.T) {
	valid := DefaultGameSettings()

	tests := []struct {
		name    string
		mod     func(GameSettings) GameSettings
		wantErr bool
	}{
		{"valid defaults", func(s GameSettings) GameSettings { return s }, false},
		{"MaxPlayers too low (1)", func(s GameSettings) GameSettings { s.MaxPlayers = 1; return s }, true},
		{"MaxPlayers too high (101)", func(s GameSettings) GameSettings { s.MaxPlayers = 101; return s }, true},
		{"MaxPlayers boundary low (2)", func(s GameSettings) GameSettings { s.MaxPlayers = 2; return s }, false},
		{"MaxPlayers boundary high (100)", func(s GameSettings) GameSettings { s.MaxPlayers = 100; return s }, false},
		{"negative RespawnDelay", func(s GameSettings) GameSettings { s.RespawnDelay = -1; return s }, true},
		{"RespawnDelay too high (301)", func(s GameSettings) GameSettings { s.RespawnDelay = 301; return s }, true},
		{"negative GameDuration", func(s GameSettings) GameSettings { s.GameDuration = -1; return s }, true},
		{"GameDuration too high (7201)", func(s GameSettings) GameSettings { s.GameDuration = 7201; return s }, true},
		{"negative ScorePerKill", func(s GameSettings) GameSettings { s.ScorePerKill = -1; return s }, true},
		{"ScorePerKill too high (10001)", func(s GameSettings) GameSettings { s.ScorePerKill = 10001; return s }, true},
		{"ScorePerKill zero (valid)", func(s GameSettings) GameSettings { s.ScorePerKill = 0; return s }, false},
		{"GameDuration zero unlimited (valid)", func(s GameSettings) GameSettings { s.GameDuration = 0; return s }, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := tc.mod(valid)
			err := s.Validate()
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !errors.Is(err, ErrValidation) {
					t.Errorf("expected error wrapping ErrValidation, got: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestDefaultGameSettings(t *testing.T) {
	s := DefaultGameSettings()

	if s.RespawnDelay != 5 {
		t.Errorf("RespawnDelay = %d, want 5", s.RespawnDelay)
	}
	if s.GameDuration != 300 {
		t.Errorf("GameDuration = %d, want 300", s.GameDuration)
	}
	if s.FriendlyFire != false {
		t.Error("FriendlyFire = true, want false")
	}
	if s.MaxPlayers != 20 {
		t.Errorf("MaxPlayers = %d, want 20", s.MaxPlayers)
	}
	if s.ScorePerKill != 100 {
		t.Errorf("ScorePerKill = %d, want 100", s.ScorePerKill)
	}
}
