package game

import "github.com/steven004/goldenfox-sudoku/engine"

// GameState represents the comprehensive state object for the UI
type GameState struct {
	Board       engine.SudokuBoard `json:"board"`
	SelectedRow int                `json:"selectedRow"`
	SelectedCol int                `json:"selectedCol"`
	IsSelected  bool               `json:"isSelected"`
	PencilMode  bool               `json:"pencilMode"`
	Mistakes    int                `json:"mistakes"`
	TimeElapsed string             `json:"timeElapsed"` // formatted string
	Difficulty  string             `json:"difficulty"`
	IsSolved    bool               `json:"isSolved"`
}
