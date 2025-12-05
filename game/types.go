package game

import "github.com/steven004/goldenfox-sudoku/engine"

// GameState represents the comprehensive state object for the UI
type GameState struct {
	Board          engine.SudokuBoard `json:"board"`
	SelectedRow    int                `json:"selectedRow"`
	SelectedCol    int                `json:"selectedCol"`
	IsSelected     bool               `json:"isSelected"`
	PencilMode     bool               `json:"pencilMode"`
	Mistakes       int                `json:"mistakes"`
	EraseCount     int                `json:"eraseCount"`
	UndoCount      int                `json:"undoCount"`
	ElapsedSeconds int                `json:"elapsedSeconds"` // in seconds
	Difficulty     string             `json:"difficulty"`
	IsSolved       bool               `json:"isSolved"`
	// User Stats
	UserLevel              int     `json:"userLevel"`
	GamesPlayed            int     `json:"gamesPlayed"`
	WinRate                float64 `json:"winRate"`
	PendingGames           int     `json:"pendingGames"`
	AverageTime            string  `json:"averageTime"` // Formatted string
	CurrentDifficultyCount int     `json:"currentDifficultyCount"`
	Progress               int     `json:"progress"` // -3 to +5
	RemainingCells         int     `json:"remainingCells"`
}
