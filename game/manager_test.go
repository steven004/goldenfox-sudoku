package game

import (
	"path/filepath"
	"testing"

	"github.com/steven004/goldenfox-sudoku/engine"
	"github.com/steven004/goldenfox-sudoku/generator"
)

func createTestGameManager(t *testing.T) *GameManager {
	dataPath := filepath.Join("..", "Data", "sudoku_curated_5000.csv")
	gen, err := generator.NewPreloadedGenerator(dataPath)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	return NewGameManager(gen)
}

func TestNewGameManager(t *testing.T) {
	gm := createTestGameManager(t)

	if gm == nil {
		t.Fatal("NewGameManager returned nil")
	}

	if gm.generator == nil {
		t.Error("Generator should not be nil")
	}

	// Should have no selection initially
	_, _, hasSelection := gm.GetSelectedCell()
	if hasSelection {
		t.Error("Should have no selection initially")
	}

	// Should not be in pencil mode
	if gm.IsPencilMode() {
		t.Error("Should not be in pencil mode initially")
	}
}

func TestNewGame(t *testing.T) {
	gm := createTestGameManager(t)

	// Start a new game
	err := gm.NewGame(engine.Medium)
	if err != nil {
		t.Fatalf("Failed to start new game: %v", err)
	}

	// Should have a board
	board := gm.GetBoard()
	if board == nil {
		t.Fatal("Board should not be nil after NewGame")
	}

	// Should have difficulty set
	if gm.GetDifficulty() != engine.Medium {
		t.Errorf("Expected difficulty Medium, got %s", gm.GetDifficulty().String())
	}

	// Should have some given cells
	givenCount := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if gm.IsCellGiven(i, j) {
				givenCount++
			}
		}
	}
	if givenCount == 0 {
		t.Error("Should have some given cells")
	}
	t.Logf("Generated puzzle with %d given cells", givenCount)

	// Should not be solved initially
	if gm.IsSolved() {
		t.Error("Puzzle should not be solved initially")
	}
}

func TestRestartGame(t *testing.T) {
	gm := createTestGameManager(t)

	// Try restart without a game
	err := gm.RestartGame()
	if err == nil {
		t.Error("Expected error when restarting without a game")
	}

	// Start a game
	gm.NewGame(engine.Easy)

	// Make some moves
	gm.SelectCell(0, 0)
	if !gm.IsCellGiven(0, 0) {
		gm.InputNumber(0, 0, 5)
	}
	gm.TogglePencilMode()

	// Restart
	err = gm.RestartGame()
	if err != nil {
		t.Errorf("Unexpected error on restart: %v", err)
	}

	// Selection should be cleared
	_, _, hasSelection := gm.GetSelectedCell()
	if hasSelection {
		t.Error("Selection should be cleared after restart")
	}

	// Pencil mode should be off
	if gm.IsPencilMode() {
		t.Error("Pencil mode should be off after restart")
	}

	// Board should be reset to initial state
	board := gm.GetBoard()
	if board == nil {
		t.Fatal("Board should not be nil after restart")
	}
}

func TestSelectCell(t *testing.T) {
	gm := createTestGameManager(t)
	gm.NewGame(engine.Beginner)

	// Valid selection
	err := gm.SelectCell(3, 4)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	row, col, hasSelection := gm.GetSelectedCell()
	if !hasSelection {
		t.Error("Should have selection")
	}
	if row != 3 || col != 4 {
		t.Errorf("Expected selection [3][4], got [%d][%d]", row, col)
	}

	// Invalid selection
	err = gm.SelectCell(-1, 0)
	if err == nil {
		t.Error("Expected error for invalid row")
	}

	err = gm.SelectCell(0, 10)
	if err == nil {
		t.Error("Expected error for invalid column")
	}

	// Clear selection
	gm.ClearSelection()
	_, _, hasSelection = gm.GetSelectedCell()
	if hasSelection {
		t.Error("Selection should be cleared")
	}
}

func TestPencilMode(t *testing.T) {
	gm := createTestGameManager(t)

	// Initially off
	if gm.IsPencilMode() {
		t.Error("Pencil mode should be off initially")
	}

	// Toggle on
	gm.TogglePencilMode()
	if !gm.IsPencilMode() {
		t.Error("Pencil mode should be on after toggle")
	}

	// Toggle off
	gm.TogglePencilMode()
	if gm.IsPencilMode() {
		t.Error("Pencil mode should be off after second toggle")
	}

	// Set explicitly
	gm.SetPencilMode(true)
	if !gm.IsPencilMode() {
		t.Error("Pencil mode should be on after SetPencilMode(true)")
	}

	gm.SetPencilMode(false)
	if gm.IsPencilMode() {
		t.Error("Pencil mode should be off after SetPencilMode(false)")
	}
}

func TestInputNumber(t *testing.T) {
	gm := createTestGameManager(t)
	gm.NewGame(engine.Beginner)

	// Find an empty cell
	var emptyRow, emptyCol int
	found := false
	for i := 0; i < 9 && !found; i++ {
		for j := 0; j < 9 && !found; j++ {
			if !gm.IsCellGiven(i, j) {
				emptyRow, emptyCol = i, j
				found = true
			}
		}
	}

	if !found {
		t.Fatal("Could not find empty cell for testing")
	}

	// Place a number
	err := gm.InputNumber(emptyRow, emptyCol, 5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	val, _ := gm.GetCellValue(emptyRow, emptyCol)
	if val != 5 {
		t.Errorf("Expected value 5, got %d", val)
	}

	// Try to modify a given cell
	var givenRow, givenCol int
	found = false
	for i := 0; i < 9 && !found; i++ {
		for j := 0; j < 9 && !found; j++ {
			if gm.IsCellGiven(i, j) {
				givenRow, givenCol = i, j
				found = true
			}
		}
	}

	if found {
		err = gm.InputNumber(givenRow, givenCol, 9)
		if err == nil {
			t.Error("Expected error when modifying given cell")
		}
	}

	// Invalid value
	err = gm.InputNumber(0, 0, 10)
	if err == nil {
		t.Error("Expected error for invalid value")
	}
}

func TestInputNumberPencilMode(t *testing.T) {
	gm := createTestGameManager(t)
	gm.NewGame(engine.Easy)

	// Find an empty cell
	var emptyRow, emptyCol int
	found := false
	for i := 0; i < 9 && !found; i++ {
		for j := 0; j < 9 && !found; j++ {
			if !gm.IsCellGiven(i, j) {
				emptyRow, emptyCol = i, j
				found = true
			}
		}
	}

	if !found {
		t.Fatal("Could not find empty cell for testing")
	}

	// Enable pencil mode
	gm.SetPencilMode(true)

	// Add candidates
	gm.InputNumber(emptyRow, emptyCol, 3)
	gm.InputNumber(emptyRow, emptyCol, 7)
	gm.InputNumber(emptyRow, emptyCol, 9)

	candidates, err := gm.GetCellCandidates(emptyRow, emptyCol)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(candidates) != 3 {
		t.Errorf("Expected 3 candidates, got %d", len(candidates))
	}
}

func TestClearCell(t *testing.T) {
	gm := createTestGameManager(t)
	gm.NewGame(engine.Medium)

	// Find an empty cell
	var emptyRow, emptyCol int
	found := false
	for i := 0; i < 9 && !found; i++ {
		for j := 0; j < 9 && !found; j++ {
			if !gm.IsCellGiven(i, j) {
				emptyRow, emptyCol = i, j
				found = true
			}
		}
	}

	if !found {
		t.Fatal("Could not find empty cell for testing")
	}

	// Place a number
	gm.InputNumber(emptyRow, emptyCol, 5)

	// Clear it
	err := gm.ClearCell(emptyRow, emptyCol)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	val, _ := gm.GetCellValue(emptyRow, emptyCol)
	if val != 0 {
		t.Errorf("Expected value 0 after clear, got %d", val)
	}
}

func TestIsNumberComplete(t *testing.T) {
	gm := createTestGameManager(t)
	gm.NewGame(engine.Hard)

	// Check initial state
	for num := 1; num <= 9; num++ {
		complete := gm.IsNumberComplete(num)
		// Most numbers should not be complete initially
		t.Logf("Number %d complete: %v", num, complete)
	}
}

func TestFindConflicts(t *testing.T) {
	gm := createTestGameManager(t)
	gm.NewGame(engine.Beginner)

	// Initially should have no conflicts
	conflicts := gm.FindConflicts()
	if len(conflicts) != 0 {
		t.Errorf("Expected no conflicts initially, got %d", len(conflicts))
	}

	// Create a conflict by placing duplicate in same row
	// Find two empty cells in the same row
	var row, col1, col2 int
	found := false
	for i := 0; i < 9 && !found; i++ {
		emptyCount := 0
		emptyCols := []int{}
		for j := 0; j < 9; j++ {
			if !gm.IsCellGiven(i, j) {
				emptyCols = append(emptyCols, j)
				emptyCount++
				if emptyCount >= 2 {
					row = i
					col1 = emptyCols[0]
					col2 = emptyCols[1]
					found = true
					break
				}
			}
		}
	}

	if found {
		// Place same number in both cells
		gm.InputNumber(row, col1, 5)
		gm.InputNumber(row, col2, 5)

		conflicts = gm.FindConflicts()
		if len(conflicts) == 0 {
			t.Error("Expected conflicts after placing duplicates")
		}
		t.Logf("Found %d conflicts", len(conflicts))
	}
}

func TestIsSolved(t *testing.T) {
	gm := createTestGameManager(t)
	gm.NewGame(engine.Expert)

	// Should not be solved initially
	if gm.IsSolved() {
		t.Error("Puzzle should not be solved initially")
	}

	// Note: We don't test actual solving here as that would require
	// implementing a solver or manually solving a puzzle
}

func TestGetCellValue(t *testing.T) {
	gm := createTestGameManager(t)

	// No game in progress
	_, err := gm.GetCellValue(0, 0)
	if err == nil {
		t.Error("Expected error when no game in progress")
	}

	// Start game
	gm.NewGame(engine.Medium)

	// Get value
	val, err := gm.GetCellValue(0, 0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	t.Logf("Cell [0][0] value: %d", val)
}

func TestIsCellGiven(t *testing.T) {
	gm := createTestGameManager(t)
	gm.NewGame(engine.Easy)

	// Count given cells
	givenCount := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if gm.IsCellGiven(i, j) {
				givenCount++
			}
		}
	}

	if givenCount == 0 {
		t.Error("Should have some given cells")
	}
	t.Logf("Found %d given cells", givenCount)
}
