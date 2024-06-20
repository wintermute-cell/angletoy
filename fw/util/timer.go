package util

import (
	"sync"
	"time"
)

type Timer struct {
	interval time.Duration
	lastTime time.Time
	pauseDur time.Duration
	pausedAt time.Time
	mu       sync.Mutex
	paused   bool
}

// New creates a new Timer with the given interval.
func NewTimer(interval float32) *Timer {
	return &Timer{
		interval: time.Duration(interval * float32(time.Second)),
		lastTime: time.Now(),
	}
}

// Check checks if the interval has passed since the last successful check.
// It returns true if it has, false otherwise.
func (t *Timer) Check() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.paused {
		return false
	}

	now := time.Now()
	elapsed := now.Sub(t.lastTime) - t.pauseDur

	if elapsed >= t.interval {
		t.lastTime = now
		t.pauseDur = 0
		return true
	}

	return false
}

// Pause pauses the timer.
func (t *Timer) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.paused {
		t.pausedAt = time.Now()
		t.paused = true
	}
}

// Continue continues the timer from where it was paused.
func (t *Timer) Continue() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.paused {
		t.pauseDur += time.Now().Sub(t.pausedAt)
		t.paused = false
	}
}

// ResetTime resets the timer's elapsed time to 0.
func (t *Timer) ResetTime() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.lastTime = time.Now()
	t.pauseDur = 0
	if t.paused {
		t.pausedAt = t.lastTime
	}
}
