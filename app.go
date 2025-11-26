package main

import (
	"context"
	"fmt"

	"github.com/steven004/goldenfox-sudoku/engine"
	"github.com/steven004/goldenfox-sudoku/game"
	"github.com/steven004/goldenfox-sudoku/generator"
)

// App struct
type App struct {
	ctx         context.Context
	gameManager *game.GameManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Try to initialize the preloaded generator with the dataset
	dataPath := generator.GetDefaultDataPath()
	gen, err := generator.NewPreloadedGenerator(dataPath)

	var gameManager *game.GameManager
	if err != nil {
		fmt.Printf("Warning: Failed to load puzzle dataset from %s: %v\n", dataPath, err)
		fmt.Println("Falling back to simple generator")
		simpleGen := generator.NewSimpleGenerator()
		gameManager = game.NewGameManager(simpleGen)
	} else {
		fmt.Printf("Successfully loaded puzzle dataset from %s\n", dataPath)
		gameManager = game.NewGameManager(gen)
	}

	return &App{
		gameManager: gameManager,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Start a default game
	a.gameManager.NewGame(engine.Beginner)
}

// --- Exposed Methods ---

// GetBoard returns the current board state
func (a *App) GetBoard() engine.SudokuBoard {
	if board := a.gameManager.GetBoard(); board != nil {
		return *board
	}
	return engine.SudokuBoard{}
}

// NewGame starts a new game with the given difficulty
func (a *App) NewGame(difficultyStr string) error {
	var diff engine.DifficultyLevel
	switch difficultyStr {
	case "Beginner":
		diff = engine.Beginner
	case "Easy":
		diff = engine.Easy
	case "Medium":
		diff = engine.Medium
	case "Hard":
		diff = engine.Hard
	case "Expert":
		diff = engine.Expert
	case "FoxGod":
		diff = engine.FoxGod
	default:
		return fmt.Errorf("invalid difficulty: %s", difficultyStr)
	}
	return a.gameManager.NewGame(diff)
}

// SelectCell selects a cell at the given row and column
func (a *App) SelectCell(row, col int) {
	a.gameManager.SelectCell(row, col)
}

// InputNumber inputs a number into the selected cell
func (a *App) InputNumber(val int) error {
	row, col, selected := a.gameManager.GetSelectedCell()
	if !selected {
		return fmt.Errorf("no cell selected")
	}
	return a.gameManager.InputNumber(row, col, val)
}

// TogglePencilMode toggles the pencil mode
func (a *App) TogglePencilMode() bool {
	a.gameManager.TogglePencilMode()
	return a.gameManager.IsPencilMode()
}

// GetGameState returns a comprehensive state object for the UI
type GameState struct {
	Board       engine.SudokuBoard `json:"board"`
	SelectedRow int                `json:"selectedRow"`
	SelectedCol int                `json:"selectedCol"`
	IsSelected  bool               `json:"isSelected"`
	PencilMode  bool               `json:"pencilMode"`
	Mistakes    int                `json:"mistakes"`
	TimeElapsed string             `json:"timeElapsed"` // formatted string
	Difficulty  string             `json:"difficulty"`
}

func (a *App) GetGameState() GameState {
	row, col, selected := a.gameManager.GetSelectedCell()
	board := engine.SudokuBoard{}
	if b := a.gameManager.GetBoard(); b != nil {
		board = *b
	}
	return GameState{
		Board:       board,
		SelectedRow: row,
		SelectedCol: col,
		IsSelected:  selected,
		PencilMode:  a.gameManager.IsPencilMode(),
		Mistakes:    a.gameManager.GetMistakes(),
		TimeElapsed: a.gameManager.GetElapsedTime(), // Assuming this returns string or duration
		Difficulty:  a.gameManager.GetDifficulty().String(),
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
