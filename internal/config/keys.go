package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const keysFile = "keys.json"

// Action represents a game action that can be rebound.
type Action string

const (
	ActionMoveLeft  Action = "move_left"
	ActionMoveRight Action = "move_right"
	ActionSoftDrop  Action = "soft_drop"
	ActionHardDrop  Action = "hard_drop"
	ActionRotateCW  Action = "rotate_cw"
	ActionRotateCCW Action = "rotate_ccw"
	ActionHold      Action = "hold"
	ActionPause     Action = "pause"
)

var AllActions = []Action{
	ActionMoveLeft, ActionMoveRight, ActionSoftDrop, ActionHardDrop,
	ActionRotateCW, ActionRotateCCW, ActionHold, ActionPause,
}

// ActionLabel returns a human-readable label for an action.
func ActionLabel(a Action) string {
	switch a {
	case ActionMoveLeft:
		return "Move Left"
	case ActionMoveRight:
		return "Move Right"
	case ActionSoftDrop:
		return "Soft Drop"
	case ActionHardDrop:
		return "Hard Drop"
	case ActionRotateCW:
		return "Rotate CW"
	case ActionRotateCCW:
		return "Rotate CCW"
	case ActionHold:
		return "Hold"
	case ActionPause:
		return "Pause"
	default:
		return string(a)
	}
}

// KeyBindings maps actions to their bound keys.
type KeyBindings struct {
	Bindings map[Action][]string `json:"bindings"`
}

// DefaultKeyBindings returns the default key bindings.
func DefaultKeyBindings() *KeyBindings {
	return &KeyBindings{
		Bindings: map[Action][]string{
			ActionMoveLeft:  {"left", "a"},
			ActionMoveRight: {"right", "d"},
			ActionSoftDrop:  {"down", "s"},
			ActionHardDrop:  {"up", " "},
			ActionRotateCW:  {"w", "e"},
			ActionRotateCCW: {"q", "z"},
			ActionHold:      {"c", "shift+c"},
			ActionPause:     {"p", "esc"},
		},
	}
}

// MatchAction finds the action for a given key string.
func (kb *KeyBindings) MatchAction(key string) (Action, bool) {
	for action, keys := range kb.Bindings {
		for _, k := range keys {
			if k == key {
				return action, true
			}
		}
	}
	return "", false
}

// SetBinding sets the keys for an action.
func (kb *KeyBindings) SetBinding(action Action, keys []string) {
	kb.Bindings[action] = keys
}

func keysPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, keysFile), nil
}

// LoadKeyBindings reads key bindings from disk.
func LoadKeyBindings() *KeyBindings {
	path, err := keysPath()
	if err != nil {
		return DefaultKeyBindings()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultKeyBindings()
	}

	kb := DefaultKeyBindings()
	if err := json.Unmarshal(data, kb); err != nil {
		return DefaultKeyBindings()
	}

	return kb
}

// Save writes key bindings to disk.
func (kb *KeyBindings) Save() error {
	path, err := keysPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(kb, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
