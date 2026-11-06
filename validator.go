package router

import (
	"errors"
	"strings"
)

var (
	ErrEmptyTenantID    = errors.New("tenant ID cannot be empty")
	ErrUnauthorizedRoom = errors.New("room ID is not in the allowlist")
	ErrPayloadTooLarge  = errors.New("message payload exceeds maximum token limit")
	ErrInvalidPrompt    = errors.New("prompt contains kosong or whitespace only")
)

type DispatchConfig struct {
	AllowedRooms  []string
	MaxPayloadLen int
}

// ValidateDispatch memeriksa kelayakan request pesan sebelum diteruskan ke container AI Agent.
func ValidateDispatch(tenantID, roomID, prompt string, cfg DispatchConfig) (string, error) {
	cleanTenant := strings.TrimSpace(tenantID)
	if cleanTenant == "" {
		return "", ErrEmptyTenantID
	}

	if len(cfg.AllowedRooms) == 0 {
		return "", ErrUnauthorizedRoom
	}

	allowed := false
	for _, r := range cfg.AllowedRooms {
		if r == roomID {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", ErrUnauthorizedRoom
	}

	cleanPrompt := strings.TrimSpace(prompt)
	if cleanPrompt == "" {
		return "", ErrInvalidPrompt
	}

	maxLen := cfg.MaxPayloadLen
	if maxLen <= 0 {
		maxLen = 4096 // default fallback
	}

	if len(cleanPrompt) > maxLen {
		return "", ErrPayloadTooLarge
	}

	return cleanTenant + ":" + roomID, nil
}
