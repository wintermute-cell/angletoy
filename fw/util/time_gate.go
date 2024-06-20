package util

import (
	"sync"
	"time"
)

type TimeGate struct {
	duration  time.Duration
	startTime time.Time
	active    bool
	mu        sync.Mutex
}

// NewTimeGate creates a new TimeGate with the given duration (in seconds).
func NewTimeGate(duration float32) *TimeGate {
	return &TimeGate{
		duration: time.Duration(duration * float32(time.Second)),
	}
}

// Start activates the TimeGate.
func (g *TimeGate) Start() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.startTime = time.Now()
	g.active = true
}

// Check checks if the TimeGate is active and if the time since it was last started is less than the specified duration.
func (g *TimeGate) Check() bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.active {
		return false
	}

	elapsed := time.Now().Sub(g.startTime)

	if elapsed >= g.duration {
		g.active = false
		return false
	}

	return true
}

// Reset deactivates the TimeGate without waiting for its duration to elapse.
func (g *TimeGate) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.active = false
}
