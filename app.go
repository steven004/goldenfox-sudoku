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

	var gameManager *game.GameManager
	if err != nil {
		fmt.Printf("Warning: Failed to load embedded puzzle dataset: %v\n", err)
		fmt.Println("Falling back to simple generator")
		simpleGen := generator.NewSimpleGenerator()
		gameManager = game.NewGameManager(simpleGen)
	} else {
		fmt.Println("Successfully loaded embedded puzzle dataset")
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

	// Enforce 1:1 Aspect Ratio (macOS only)
	SetWindowAspectRatio()

	// Determine difficulty based on user level
	level := a.gameManager.GetUserLevel()
	var diff engine.DifficultyLevel

	switch {
	case level < 5:
		diff = engine.Beginner
	case level < 10:
		diff = engine.Easy
	case level < 20:
		diff = engine.Medium
	case level < 50:
		diff = engine.Hard
	default:
		diff = engine.Expert
	}

	fmt.Printf("Starting initial game for User Level %d (Difficulty: %s)\n", level, diff)

	// Start a new game using the standard logic (generates ID, etc.)
	if err := a.gameManager.NewGame(diff); err != nil {
		fmt.Printf("Error starting initial game: %v\n", err)
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

// LoadGame loads a specific game from history
func (a *App) LoadGame(id string) error {
	return a.gameManager.LoadGame(id)
}

// SaveGame saves the current game progress
func (a *App) SaveGame() error {
	return a.gameManager.SaveCurrentGame()
}

// TogglePencilMode toggles the pencil mode
func (a *App) TogglePencilMode() bool {
	a.gameManager.TogglePencilMode()
	return a.gameManager.IsPencilMode()
}

// GetGameState returns a comprehensive state object for the UI
func (a *App) GetGameState() game.GameState {
	return a.gameManager.GetGameState()
}

// ClearCell clears the selected cell
func (a *App) ClearCell() error {
	row, col, selected := a.gameManager.GetSelectedCell()
	if !selected {
		return fmt.Errorf("no cell selected")
	}
	return a.gameManager.ClearCell(row, col)
}

// Undo reverts the last move
func (a *App) Undo() error {
	return a.gameManager.Undo()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetHistory returns the user's puzzle history
func (a *App) GetHistory() []game.PuzzleRecord {
	return a.gameManager.GetHistory()
}
