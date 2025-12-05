package game

import "time"

// GameTimer handles the timekeeping for a game session
type GameTimer struct {
	startTime  time.Time
	endTime    time.Time
	pausedTime time.Duration
}

// NewGameTimer creates a new timer instance
func NewGameTimer() *GameTimer {
	return &GameTimer{
		startTime: time.Now(),
	}
}

// Reset resets the timer to the current time
func (gt *GameTimer) Reset() {
	gt.startTime = time.Now()
	gt.endTime = time.Time{}
	gt.pausedTime = 0
}

// Stop marks the timer as stopped (game finished)
func (gt *GameTimer) Stop() {
	gt.endTime = time.Now()
}

// IsRunning returns true if the timer has started but not stopped
func (gt *GameTimer) IsRunning() bool {
	return !gt.startTime.IsZero() && gt.endTime.IsZero()
}

// GetElapsedDuration returns the total elapsed time
func (gt *GameTimer) GetElapsedDuration() time.Duration {
	if !gt.endTime.IsZero() {
		return gt.endTime.Sub(gt.startTime)
	}
	return time.Since(gt.startTime)
}

// SetStartTime allows manually setting the start time (for loading games)
func (gt *GameTimer) SetStartTime(t time.Time) {
	gt.startTime = t
}

// SetEndTime allows manually setting the end time (for loading games)
func (gt *GameTimer) SetEndTime(t time.Time) {
	gt.endTime = t
}
