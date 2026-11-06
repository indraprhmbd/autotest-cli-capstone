package router

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateDispatch(t *testing.T) {
	defaultCfg := DispatchConfig{
		AllowedRooms:  []string{"room-1", "room-2"},
		MaxPayloadLen: 100,
	}

	tests := []struct {
		name      string
		tenantID  string
		roomID    string
		prompt    string
		cfg       DispatchConfig
		wantRoute string
		wantErr   error
	}{
		// ── Happy Path ────────────────────────────────────────────────────────
		{
			name:      "happy path: valid inputs with custom MaxPayloadLen",
			tenantID:  "tenant-abc",
			roomID:    "room-1",
			prompt:    "Hello, agent!",
			cfg:       defaultCfg,
			wantRoute: "tenant-abc:room-1",
			wantErr:   nil,
		},
		{
			name:     "happy path: second allowed room",
			tenantID: "tenant-xyz",
			roomID:   "room-2",
			prompt:   "What is the status?",
			cfg:      defaultCfg,
			wantRoute: "tenant-xyz:room-2",
			wantErr:  nil,
		},
		{
			name:     "happy path: MaxPayloadLen <= 0 falls back to 4096 default",
			tenantID: "tenant-fallback",
			roomID:   "room-1",
			prompt:   "A short prompt.",
			cfg: DispatchConfig{
				AllowedRooms:  []string{"room-1"},
				MaxPayloadLen: 0,
			},
			wantRoute: "tenant-fallback:room-1",
			wantErr:   nil,
		},
		{
			name:     "happy path: negative MaxPayloadLen also falls back to 4096 default",
			tenantID: "tenant-neg",
			roomID:   "room-1",
			prompt:   "Another short prompt.",
			cfg: DispatchConfig{
				AllowedRooms:  []string{"room-1"},
				MaxPayloadLen: -50,
			},
			wantRoute: "tenant-neg:room-1",
			wantErr:   nil,
		},
		{
			name:     "happy path: prompt at exact MaxPayloadLen boundary",
			tenantID: "tenant-boundary",
			roomID:   "room-1",
			prompt:   strings.Repeat("x", 100),
			cfg:      defaultCfg,
			wantRoute: "tenant-boundary:room-1",
			wantErr:  nil,
		},
		{
			name:     "happy path: tenant ID with surrounding whitespace is trimmed",
			tenantID: "  tenant-spaces  ",
			roomID:   "room-1",
			prompt:   "Valid prompt.",
			cfg:      defaultCfg,
			wantRoute: "tenant-spaces:room-1",
			wantErr:  nil,
		},

		// ── Boundary Conditions ───────────────────────────────────────────────
		{
			name:     "boundary: prompt at maxLen + 1 returns ErrPayloadTooLarge",
			tenantID: "tenant-abc",
			roomID:   "room-1",
			prompt:   strings.Repeat("y", 101),
			cfg:      defaultCfg,
			wantRoute: "",
			wantErr:  ErrPayloadTooLarge,
		},
		{
			name:     "boundary: prompt at default fallback max (4096) is accepted",
			tenantID: "tenant-abc",
			roomID:   "room-1",
			prompt:   strings.Repeat("a", 4096),
			cfg: DispatchConfig{
				AllowedRooms:  []string{"room-1"},
				MaxPayloadLen: 0,
			},
			wantRoute: "tenant-abc:room-1",
			wantErr:   nil,
		},
		{
			name:     "boundary: prompt at default fallback max + 1 returns ErrPayloadTooLarge",
			tenantID: "tenant-abc",
			roomID:   "room-1",
			prompt:   strings.Repeat("a", 4097),
			cfg: DispatchConfig{
				AllowedRooms:  []string{"room-1"},
				MaxPayloadLen: 0,
			},
			wantRoute: "",
			wantErr:   ErrPayloadTooLarge,
		},
		{
			name:     "boundary: prompt is single character (minimal valid)",
			tenantID: "t",
			roomID:   "room-1",
			prompt:   "x",
			cfg:      defaultCfg,
			wantRoute: "t:room-1",
			wantErr:  nil,
		},

		// ── Negative & Security Cases ─────────────────────────────────────────
		{
			name:     "negative: empty tenant ID returns ErrEmptyTenantID",
			tenantID: "",
			roomID:   "room-1",
			prompt:   "Valid prompt.",
			cfg:      defaultCfg,
			wantRoute: "",
			wantErr:  ErrEmptyTenantID,
		},
		{
			name:     "negative: whitespace-only tenant ID returns ErrEmptyTenantID",
			tenantID: "   ",
			roomID:   "room-1",
			prompt:   "Valid prompt.",
			cfg:      defaultCfg,
			wantRoute: "",
			wantErr:  ErrEmptyTenantID,
		},
		{
			name:     "negative: tab-only tenant ID returns ErrEmptyTenantID",
			tenantID: "\t\n",
			roomID:   "room-1",
			prompt:   "Valid prompt.",
			cfg:      defaultCfg,
			wantRoute: "",
			wantErr:  ErrEmptyTenantID,
		},
		{
			name:     "negative: empty AllowedRooms slice returns ErrUnauthorizedRoom",
			tenantID: "tenant-abc",
			roomID:   "room-1",
			prompt:   "Valid prompt.",
			cfg: DispatchConfig{
				AllowedRooms:  []string{},
				MaxPayloadLen: 100,
			},
			wantRoute: "",
			wantErr:  ErrUnauthorizedRoom,
		},
		{
			name:     "negative: nil AllowedRooms slice returns ErrUnauthorizedRoom",
			tenantID: "tenant-abc",
			roomID:   "room-1",
			prompt:   "Valid prompt.",
			cfg: DispatchConfig{
				AllowedRooms:  nil,
				MaxPayloadLen: 100,
			},
			wantRoute: "",
			wantErr:  ErrUnauthorizedRoom,
		},
		{
			name:     "negative: roomID not in allowlist returns ErrUnauthorizedRoom",
			tenantID: "tenant-abc",
			roomID:   "room-99",
			prompt:   "Valid prompt.",
			cfg:      defaultCfg,
			wantRoute: "",
			wantErr:  ErrUnauthorizedRoom,
		},
		{
			name:     "negative: empty roomID not in allowlist returns ErrUnauthorizedRoom",
			tenantID: "tenant-abc",
			roomID:   "",
			prompt:   "Valid prompt.",
			cfg:      defaultCfg,
			wantRoute: "",
			wantErr:  ErrUnauthorizedRoom,
		},
		{
			name:     "negative: empty prompt returns ErrInvalidPrompt",
			tenantID: "tenant-abc",
			roomID:   "room-1",
			prompt:   "",
			cfg:      defaultCfg,
			wantRoute: "",
			wantErr:  ErrInvalidPrompt,
		},
		{
			name:     "negative: whitespace-only prompt returns ErrInvalidPrompt",
			tenantID: "tenant-abc",
			roomID:   "room-1",
			prompt:   "   ",
			cfg:      defaultCfg,
			wantRoute: "",
			wantErr:  ErrInvalidPrompt,
		},
		{
			name:     "negative: tab-and-newline-only prompt returns ErrInvalidPrompt",
			tenantID: "tenant-abc",
			roomID:   "room-1",
			prompt:   "\t\n\r",
			cfg:      defaultCfg,
			wantRoute: "",
			wantErr:  ErrInvalidPrompt,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotRoute, gotErr := ValidateDispatch(tc.tenantID, tc.roomID, tc.prompt, tc.cfg)

			if gotRoute != tc.wantRoute {
				t.Errorf("route: got %q, want %q", gotRoute, tc.wantRoute)
			}

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("error: got %v, want %v", gotErr, tc.wantErr)
			}
		})
	}
}
