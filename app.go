package main

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/steven004/goldenfox-sudoku/engine"
	"github.com/steven004/goldenfox-sudoku/game"
	"github.com/steven004/goldenfox-sudoku/generator"
)

//go:embed Data/sudoku_curated_5000.csv
var puzzleData []byte

// App struct
type App struct {
	ctx         context.Context
	gameManager *game.GameManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Initialize the preloaded generator with the embedded dataset
	gen, err := generator.NewPreloadedGenerator(puzzleData)
	if err != nil {
		// This should theoretically never happen if the embed is correct
		panic(fmt.Sprintf("Critical Error: Failed to load embedded puzzle dataset: %v", err))
	}

	fmt.Println("Successfully loaded embedded puzzle dataset")
	gameManager := game.NewGameManager(gen)

	return &App{
		gameManager: gameManager,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Enforce 1:1 Aspect Ratio (macOS only)
	SetWindowAspectRatio()

	// Check legacy difficulty determination logic - we can simplify this now
	// Since NewGame("") handles the user level default internally.

	// Try to load the last played game (if any)
	lastGameID := a.gameManager.GetLastGameID()
	if lastGameID != "" {
		fmt.Printf("Resuming last game: %s\n", lastGameID)
		if err := a.gameManager.LoadGame(lastGameID); err == nil {
			return // Successfully resumed
		} else {
			fmt.Printf("Failed to resume last game: %v\n", err)
		}
	}

	// If loading failed or no last game, start a new one (using default/user level)
	if err := a.gameManager.NewGame(""); err != nil {
		fmt.Printf("Startup Warning: Failed to create new game: %v\n", err)
	}
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
	// difficultyStr can be "Hard" or empty "" (to use user level)
	return a.gameManager.NewGame(difficultyStr)
}

// RestartGame restarts the current game to its initial state
func (a *App) RestartGame() error {
	return a.gameManager.RestartGame()
}

// InputNumber inputs a number into the specified cell (Value)
func (a *App) InputNumber(row, col, val int) error {
	return a.gameManager.InputNumber(row, col, val)
}

// ToggleCandidate toggles a candidate note in the specified cell
func (a *App) ToggleCandidate(row, col, val int) error {
	return a.gameManager.ToggleCandidate(row, col, val)
}

// LoadGame loads a specific game from history
func (a *App) LoadGame(id string) error {
	return a.gameManager.LoadGame(id)
}

// SaveGame saves the current game progress
func (a *App) SaveGame() error {
	return a.gameManager.SaveCurrentGame()
}

// GetGameState returns a comprehensive state object for the UI
func (a *App) GetGameState() game.GameState {
	return a.gameManager.GetGameState()
}

// ClearCell clears the specified cell
func (a *App) ClearCell(row, col int) error {
	return a.gameManager.ClearCell(row, col)
}

// Undo reverts the last move
func (a *App) Undo() error {
	return a.gameManager.Undo()
}

// GetHistory returns the user's puzzle history
func (a *App) GetHistory() []game.PuzzleRecord {
	return a.gameManager.GetHistory()
}
