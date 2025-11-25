package engine

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	if board == nil {
		t.Fatal("NewBoard returned nil")
	}

	// Verify all cells are initialized
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board.Cells[i][j].Value != 0 {
				t.Errorf("Cell [%d][%d] should be empty, got %d", i, j, board.Cells[i][j].Value)
			}
			if board.Cells[i][j].Given {
				t.Errorf("Cell [%d][%d] should not be marked as given", i, j)
			}
			if board.Cells[i][j].Candidates == nil {
				t.Errorf("Cell [%d][%d] candidates map is nil", i, j)
			}
		}
	}
}

func TestSetValue(t *testing.T) {
	board := NewBoard()

	// Valid set
	err := board.SetValue(0, 0, 5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if board.Cells[0][0].Value != 5 {
		t.Errorf("Expected value 5, got %d", board.Cells[0][0].Value)
	}

	// Invalid position
	err = board.SetValue(-1, 0, 5)
	if err == nil {
		t.Error("Expected error for invalid row, got nil")
	}

	err = board.SetValue(0, 10, 5)
	if err == nil {
		t.Error("Expected error for invalid column, got nil")
	}

	// Invalid value
	err = board.SetValue(0, 1, 10)
	if err == nil {
		t.Error("Expected error for invalid value, got nil")
	}

	// Cannot modify given cell
	board.Cells[1][1].Given = true
	err = board.SetValue(1, 1, 5)
	if err == nil {
		t.Error("Expected error for modifying given cell, got nil")
	}
}

func TestGetValue(t *testing.T) {
	board := NewBoard()
	board.Cells[3][4].Value = 7

	val, err := board.GetValue(3, 4)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if val != 7 {
		t.Errorf("Expected 7, got %d", val)
	}

	// Invalid position
	_, err = board.GetValue(-1, 0)
	if err == nil {
		t.Error("Expected error for invalid position, got nil")
	}
}

func TestIsValidMove(t *testing.T) {
	board := NewBoard()

	// Empty board - all moves should be valid
	if !board.IsValidMove(0, 0, 5) {
		t.Error("Move should be valid on empty board")
	}

	// Place a value and test conflicts
	board.Cells[0][0].Value = 5

	// Row conflict
	if board.IsValidMove(0, 5, 5) {
		t.Error("Should detect row conflict")
	}

	// Column conflict
	if board.IsValidMove(5, 0, 5) {
		t.Error("Should detect column conflict")
	}

	// Block conflict
	if board.IsValidMove(1, 1, 5) {
		t.Error("Should detect block conflict")
	}

	// Valid move
	if !board.IsValidMove(0, 5, 3) {
		t.Error("Valid move should be allowed")
	}

	// Invalid value
	if board.IsValidMove(0, 0, 0) {
		t.Error("Should reject value 0")
	}
	if board.IsValidMove(0, 0, 10) {
		t.Error("Should reject value 10")
	}
}

func TestFindConflicts(t *testing.T) {
	board := NewBoard()

	// No conflicts on empty board
	conflicts := board.FindConflicts()
	if len(conflicts) != 0 {
		t.Errorf("Expected no conflicts, got %d", len(conflicts))
	}

	// Create row conflict
	board.Cells[0][0].Value = 5
	board.Cells[0][5].Value = 5

	conflicts = board.FindConflicts()
	if len(conflicts) == 0 {
		t.Error("Should detect row conflict")
	}

	// Create column conflict
	board2 := NewBoard()
	board2.Cells[0][0].Value = 3
	board2.Cells[5][0].Value = 3

	conflicts = board2.FindConflicts()
	if len(conflicts) == 0 {
		t.Error("Should detect column conflict")
	}

	// Create block conflict
	board3 := NewBoard()
	board3.Cells[0][0].Value = 7
	board3.Cells[2][2].Value = 7

	conflicts = board3.FindConflicts()
	if len(conflicts) == 0 {
		t.Error("Should detect block conflict")
	}
}

func TestIsSolved(t *testing.T) {
	board := NewBoard()

	// Empty board is not solved
	if board.IsSolved() {
		t.Error("Empty board should not be solved")
	}

	// Partially filled board is not solved
	board.Cells[0][0].Value = 5
	if board.IsSolved() {
		t.Error("Partially filled board should not be solved")
	}

	// Create a valid complete board (simplified test)
	// Fill with a valid pattern
	validBoard := createValidSolvedBoard()
	if !validBoard.IsSolved() {
		t.Error("Valid complete board should be solved")
	}

	// Board with conflicts is not solved
	conflictBoard := createValidSolvedBoard()
	conflictBoard.Cells[0][0].Value = conflictBoard.Cells[0][1].Value
	if conflictBoard.IsSolved() {
		t.Error("Board with conflicts should not be solved")
	}
}

func TestAddCandidate(t *testing.T) {
	board := NewBoard()

	// Add valid candidate
	err := board.AddCandidate(0, 0, 5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !board.Cells[0][0].Candidates[5] {
		t.Error("Candidate 5 should be added")
	}

	// Add multiple candidates
	board.AddCandidate(0, 0, 3)
	board.AddCandidate(0, 0, 7)
	if len(board.Cells[0][0].Candidates) != 3 {
		t.Errorf("Expected 3 candidates, got %d", len(board.Cells[0][0].Candidates))
	}

	// Cannot add candidate to filled cell
	board.Cells[1][1].Value = 9
	err = board.AddCandidate(1, 1, 5)
	if err == nil {
		t.Error("Expected error for adding candidate to filled cell")
	}

	// Invalid position
	err = board.AddCandidate(-1, 0, 5)
	if err == nil {
		t.Error("Expected error for invalid position")
	}

	// Invalid value
	err = board.AddCandidate(0, 0, 10)
	if err == nil {
		t.Error("Expected error for invalid candidate value")
	}
}

func TestRemoveCandidate(t *testing.T) {
	board := NewBoard()

	// Add and remove candidate
	board.AddCandidate(0, 0, 5)
	board.AddCandidate(0, 0, 3)

	err := board.RemoveCandidate(0, 0, 5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if board.Cells[0][0].Candidates[5] {
		t.Error("Candidate 5 should be removed")
	}
	if !board.Cells[0][0].Candidates[3] {
		t.Error("Candidate 3 should still exist")
	}

	// Remove non-existent candidate (should not error)
	err = board.RemoveCandidate(0, 0, 9)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetCandidates(t *testing.T) {
	board := NewBoard()

	// Empty candidates
	candidates, err := board.GetCandidates(0, 0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(candidates) != 0 {
		t.Errorf("Expected 0 candidates, got %d", len(candidates))
	}

	// Add candidates
	board.AddCandidate(0, 0, 5)
	board.AddCandidate(0, 0, 3)
	board.AddCandidate(0, 0, 7)

	candidates, err = board.GetCandidates(0, 0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(candidates) != 3 {
		t.Errorf("Expected 3 candidates, got %d", len(candidates))
	}
}

func TestRemoveCandidateFromPeers(t *testing.T) {
	board := NewBoard()

	// Add candidate 5 to multiple cells
	for i := 0; i < 9; i++ {
		board.AddCandidate(0, i, 5) // Same row
		board.AddCandidate(i, 0, 5) // Same column
	}
	board.AddCandidate(1, 1, 5) // Same block

	// Remove from peers of [0][0]
	board.RemoveCandidateFromPeers(0, 0, 5)

	// Check row
	for c := 1; c < 9; c++ {
		if board.Cells[0][c].Candidates[5] {
			t.Errorf("Candidate 5 should be removed from row at [0][%d]", c)
		}
	}

	// Check column
	for r := 1; r < 9; r++ {
		if board.Cells[r][0].Candidates[5] {
			t.Errorf("Candidate 5 should be removed from column at [%d][0]", r)
		}
	}

	// Check block
	if board.Cells[1][1].Candidates[5] {
		t.Error("Candidate 5 should be removed from block at [1][1]")
	}

	// Original cell should keep candidate
	if !board.Cells[0][0].Candidates[5] {
		t.Error("Original cell [0][0] should keep candidate 5")
	}
}

func TestClearCell(t *testing.T) {
	board := NewBoard()

	// Set and clear a value
	board.SetValue(0, 0, 5)
	err := board.ClearCell(0, 0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if board.Cells[0][0].Value != 0 {
		t.Errorf("Cell should be cleared, got %d", board.Cells[0][0].Value)
	}

	// Cannot clear given cell
	board.Cells[1][1].Given = true
	board.Cells[1][1].Value = 5
	err = board.ClearCell(1, 1)
	if err == nil {
		t.Error("Expected error for clearing given cell")
	}
}

func TestCountNumber(t *testing.T) {
	board := NewBoard()

	// Empty board
	count := board.CountNumber(5)
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	// Add some 5s
	board.Cells[0][0].Value = 5
	board.Cells[1][1].Value = 5
	board.Cells[2][2].Value = 5

	count = board.CountNumber(5)
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	// Count different number
	count = board.CountNumber(7)
	if count != 0 {
		t.Errorf("Expected count 0 for number 7, got %d", count)
	}

	// Invalid number
	count = board.CountNumber(10)
	if count != 0 {
		t.Errorf("Expected count 0 for invalid number, got %d", count)
	}
}

func TestClone(t *testing.T) {
	board := NewBoard()
	board.Cells[0][0].Value = 5
	board.Cells[0][0].Given = true
	board.AddCandidate(1, 1, 3)
	board.AddCandidate(1, 1, 7)

	clone := board.Clone()

	// Verify clone has same values
	if clone.Cells[0][0].Value != 5 {
		t.Error("Clone should have same value")
	}
	if !clone.Cells[0][0].Given {
		t.Error("Clone should have same Given flag")
	}
	if !clone.Cells[1][1].Candidates[3] || !clone.Cells[1][1].Candidates[7] {
		t.Error("Clone should have same candidates")
	}

	// Verify it's a deep copy
	board.Cells[0][0].Value = 9
	if clone.Cells[0][0].Value == 9 {
		t.Error("Clone should be independent of original")
	}

	board.AddCandidate(1, 1, 5)
	if clone.Cells[1][1].Candidates[5] {
		t.Error("Clone candidates should be independent")
	}
}

func TestReset(t *testing.T) {
	board := NewBoard()

	// Set some given cells
	board.Cells[0][0].Value = 5
	board.Cells[0][0].Given = true

	// Set some user cells
	board.SetValue(1, 1, 3)
	board.AddCandidate(2, 2, 7)

	// Reset
	board.Reset()

	// Given cells should remain
	if board.Cells[0][0].Value != 5 {
		t.Error("Given cell should not be cleared")
	}
	if !board.Cells[0][0].Given {
		t.Error("Given flag should remain")
	}

	// User cells should be cleared
	if board.Cells[1][1].Value != 0 {
		t.Error("User cell should be cleared")
	}
	if len(board.Cells[2][2].Candidates) != 0 {
		t.Error("User candidates should be cleared")
	}
}

// Helper function to create a valid solved board for testing
func createValidSolvedBoard() *SudokuBoard {
	board := NewBoard()
	// A valid Sudoku solution
	solution := [][]int{
		{5, 3, 4, 6, 7, 8, 9, 1, 2},
		{6, 7, 2, 1, 9, 5, 3, 4, 8},
		{1, 9, 8, 3, 4, 2, 5, 6, 7},
		{8, 5, 9, 7, 6, 1, 4, 2, 3},
		{4, 2, 6, 8, 5, 3, 7, 9, 1},
		{7, 1, 3, 9, 2, 4, 8, 5, 6},
		{9, 6, 1, 5, 3, 7, 2, 8, 4},
		{2, 8, 7, 4, 1, 9, 6, 3, 5},
		{3, 4, 5, 2, 8, 6, 1, 7, 9},
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			board.Cells[i][j].Value = solution[i][j]
		}
	}

	return board
}
