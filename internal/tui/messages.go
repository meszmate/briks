package tui

import "time"

// TickMsg is sent on every game gravity tick.
type TickMsg struct {
	Time time.Time
}

// ScreenMsg triggers a screen re-render.
type ScreenMsg struct{}

// DASTickMsg handles delayed auto-shift repeats.
type DASTickMsg struct {
	Time time.Time
}

// RainbowTickMsg advances the rainbow theme animation.
type RainbowTickMsg struct {
	Time time.Time
}

// LockTickMsg checks lock delay expiration.
type LockTickMsg struct {
	Time time.Time
}
