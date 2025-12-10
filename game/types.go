package game

import "github.com/steven004/goldenfox-sudoku/engine"

// GameState represents the comprehensive state object for the UI
type GameState struct {
	Board engine.SudokuBoard `json:"board"`

	Mistakes        int     `json:"mistakes"`
	EraseCount      int     `json:"eraseCount"`
	UndoCount       int     `json:"undoCount"`
	Difficulty      string  `json:"difficulty"`
	DifficultyIndex float64 `json:"difficultyIndex"`
	IsSolved        bool    `json:"isSolved"`
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
