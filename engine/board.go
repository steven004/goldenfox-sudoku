package engine

import "fmt"

// SetValue sets a value in a cell if it's not a given clue
func (b *SudokuBoard) SetValue(row, col, val int) error {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return fmt.Errorf("invalid position: row=%d, col=%d", row, col)
	}

	if val < 0 || val > 9 {
		return fmt.Errorf("invalid value: %d (must be 0-9)", val)
	}

	if b.Cells[row][col].Given {
		return fmt.Errorf("cannot modify given cell at [%d][%d]", row, col)
	}

	b.Cells[row][col].Value = val

	// Clear candidates when a value is set
	if val != 0 {
		b.Cells[row][col].Candidates = make(map[int]bool)
	}

	return nil
}

// GetValue returns the value at the specified position
func (b *SudokuBoard) GetValue(row, col int) (int, error) {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return 0, fmt.Errorf("invalid position: row=%d, col=%d", row, col)
	}
	return b.Cells[row][col].Value, nil
}

// IsValidMove checks if placing a value at a position violates Sudoku rules
func (b *SudokuBoard) IsValidMove(row, col, val int) bool {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return false
	}

	if val < 1 || val > 9 {
		return false
	}

	// Check row
	for c := 0; c < 9; c++ {
		if c != col && b.Cells[row][c].Value == val {
			return false
		}
	}

	// Check column
	for r := 0; r < 9; r++ {
		if r != row && b.Cells[r][col].Value == val {
			return false
		}
	}

	// Check 3x3 block
	blockRow := (row / 3) * 3
	blockCol := (col / 3) * 3
	for r := blockRow; r < blockRow+3; r++ {
		for c := blockCol; c < blockCol+3; c++ {
			if (r != row || c != col) && b.Cells[r][c].Value == val {
				return false
			}
		}
	}

	return true
}

// FindConflicts returns all cells that violate Sudoku rules
func (b *SudokuBoard) FindConflicts() []Coordinate {
	conflicts := make(map[Coordinate]bool)

	// Check all filled cells
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			val := b.Cells[row][col].Value
			if val == 0 {
				continue
			}

			// Check if this placement is valid
			// Temporarily clear the cell to check
			b.Cells[row][col].Value = 0
			if !b.IsValidMove(row, col, val) {
				conflicts[Coordinate{Row: row, Col: col}] = true
			}
			b.Cells[row][col].Value = val
		}
	}

	// Convert map to slice
	result := make([]Coordinate, 0, len(conflicts))
	for coord := range conflicts {
		result = append(result, coord)
	}

	return result
}

// IsSolved checks if the puzzle is completely and correctly solved
func (b *SudokuBoard) IsSolved() bool {
	// Check all cells are filled
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if b.Cells[row][col].Value == 0 {
				return false
			}
		}
	}

	// Check no conflicts
	return len(b.FindConflicts()) == 0
}

// AddCandidate adds a pencil note to a cell
func (b *SudokuBoard) AddCandidate(row, col, val int) error {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return fmt.Errorf("invalid position: row=%d, col=%d", row, col)
	}

	if val < 1 || val > 9 {
		return fmt.Errorf("invalid candidate value: %d (must be 1-9)", val)
	}

	if b.Cells[row][col].Value != 0 {
		return fmt.Errorf("cannot add candidate to filled cell at [%d][%d]", row, col)
	}

	b.Cells[row][col].Candidates[val] = true
	return nil
}

// RemoveCandidate removes a pencil note from a cell
func (b *SudokuBoard) RemoveCandidate(row, col, val int) error {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return fmt.Errorf("invalid position: row=%d, col=%d", row, col)
	}

	if val < 1 || val > 9 {
		return fmt.Errorf("invalid candidate value: %d (must be 1-9)", val)
	}

	delete(b.Cells[row][col].Candidates, val)
	return nil
}

// GetCandidates returns all pencil notes for a cell
func (b *SudokuBoard) GetCandidates(row, col int) ([]int, error) {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return nil, fmt.Errorf("invalid position: row=%d, col=%d", row, col)
	}

	candidates := make([]int, 0, len(b.Cells[row][col].Candidates))
	for val := range b.Cells[row][col].Candidates {
		candidates = append(candidates, val)
	}
	return candidates, nil
}

// RemoveCandidateFromPeers removes a candidate from all cells in the same row, column, and block
func (b *SudokuBoard) RemoveCandidateFromPeers(row, col, val int) {
	if row < 0 || row > 8 || col < 0 || col > 8 || val < 1 || val > 9 {
		return
	}

	// Remove from row
	for c := 0; c < 9; c++ {
		if c != col {
			delete(b.Cells[row][c].Candidates, val)
		}
	}

	// Remove from column
	for r := 0; r < 9; r++ {
		if r != row {
			delete(b.Cells[r][col].Candidates, val)
		}
	}

	// Remove from 3x3 block
	blockRow := (row / 3) * 3
	blockCol := (col / 3) * 3
	for r := blockRow; r < blockRow+3; r++ {
		for c := blockCol; c < blockCol+3; c++ {
			if r != row || c != col {
				delete(b.Cells[r][c].Candidates, val)
			}
		}
	}
}

// ClearCell removes the value from a cell (if it's not a given)
func (b *SudokuBoard) ClearCell(row, col int) error {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return fmt.Errorf("invalid position: row=%d, col=%d", row, col)
	}

	if b.Cells[row][col].Given {
		return fmt.Errorf("cannot clear given cell at [%d][%d]", row, col)
	}

	b.Cells[row][col].Value = 0
	return nil
}

// CountNumber counts how many times a number appears on the board
func (b *SudokuBoard) CountNumber(val int) int {
	if val < 1 || val > 9 {
		return 0
	}

	count := 0
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if b.Cells[row][col].Value == val {
				count++
			}
		}
	}
	return count
}

// Clone creates a deep copy of the board
func (b *SudokuBoard) Clone() *SudokuBoard {
	clone := NewBoard()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			clone.Cells[i][j].Value = b.Cells[i][j].Value
			clone.Cells[i][j].Given = b.Cells[i][j].Given

			// Deep copy candidates
			clone.Cells[i][j].Candidates = make(map[int]bool)
			for k, v := range b.Cells[i][j].Candidates {
				clone.Cells[i][j].Candidates[k] = v
			}
		}
	}
	return clone
}

// Reset clears all non-given cells (for restart functionality)
func (b *SudokuBoard) Reset() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if !b.Cells[i][j].Given {
				b.Cells[i][j].Value = 0
				b.Cells[i][j].Candidates = make(map[int]bool)
			}
		}
	}
}

// String returns a simple string representation of the board (81 digits)
func (b *SudokuBoard) String() string {
	var s string
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			s += fmt.Sprintf("%d", b.Cells[i][j].Value)
		}
	}
	return s
}

// ParseBoard creates a SudokuBoard from its string representation
func ParseBoard(s string) (*SudokuBoard, error) {
	if len(s) != 81 {
		return nil, fmt.Errorf("invalid board string length: %d (expected 81)", len(s))
	}

	b := NewBoard()
	for i := 0; i < 81; i++ {
		val := int(s[i] - '0')
		if val < 0 || val > 9 {
			return nil, fmt.Errorf("invalid character in board string at index %d: %c", i, s[i])
		}

		row := i / 9
		col := i % 9
		b.Cells[row][col].Value = val

		// Note: We don't know which were 'Given' from just the string of values.
		// The caller must handle setting 'Given' status if loading a fresh puzzle,
		// or we might need a separate mask if we want to preserve that perfectly.
		// For now, we assume this reconstructs the *current state*.
	}
	return b, nil
}
