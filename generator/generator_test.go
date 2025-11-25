package generator

import (
	"path/filepath"
	"testing"

	"github.com/steven004/goldenfox-sudoku/engine"
)

func TestNewPreloadedGenerator(t *testing.T) {
	dataPath := filepath.Join("..", "Data", "sudoku_curated_5000.csv")

	gen, err := NewPreloadedGenerator(dataPath)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Verify all difficulty levels have puzzles
	for level := engine.Beginner; level <= engine.Expert; level++ {
		count := gen.GetPuzzleCount(level)
		if count == 0 {
			t.Errorf("No puzzles loaded for difficulty %s", level.String())
		}
		t.Logf("Loaded %d puzzles for difficulty %s", count, level.String())
	}
}

func TestNewPreloadedGenerator_InvalidPath(t *testing.T) {
	_, err := NewPreloadedGenerator("nonexistent.csv")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestGenerate(t *testing.T) {
	dataPath := filepath.Join("..", "Data", "sudoku_curated_5000.csv")
	gen, err := NewPreloadedGenerator(dataPath)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Test generating puzzles for each difficulty level
	for level := engine.Beginner; level <= engine.Expert; level++ {
		t.Run(level.String(), func(t *testing.T) {
			board, err := gen.Generate(level)
			if err != nil {
				t.Fatalf("Failed to generate puzzle: %v", err)
			}

			if board == nil {
				t.Fatal("Generated board is nil")
			}

			// Verify the board has some given clues
			givenCount := 0
			for i := 0; i < 9; i++ {
				for j := 0; j < 9; j++ {
					if board.Cells[i][j].Given {
						givenCount++
						// Verify given cells have valid values
						if board.Cells[i][j].Value < 1 || board.Cells[i][j].Value > 9 {
							t.Errorf("Invalid value in given cell [%d][%d]: %d",
								i, j, board.Cells[i][j].Value)
						}
					}
				}
			}

			if givenCount == 0 {
				t.Error("Generated puzzle has no given clues")
			}

			t.Logf("Generated %s puzzle with %d clues", level.String(), givenCount)
		})
	}
}

func TestGenerate_Randomness(t *testing.T) {
	dataPath := filepath.Join("..", "Data", "sudoku_curated_5000.csv")
	gen, err := NewPreloadedGenerator(dataPath)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Generate multiple puzzles and verify they're different
	const numPuzzles = 10
	puzzles := make([]string, numPuzzles)

	for i := 0; i < numPuzzles; i++ {
		board, err := gen.Generate(engine.Medium)
		if err != nil {
			t.Fatalf("Failed to generate puzzle %d: %v", i, err)
		}

		// Convert board to string for comparison
		puzzles[i] = boardToString(board)
	}

	// Check that at least some puzzles are different
	allSame := true
	for i := 1; i < numPuzzles; i++ {
		if puzzles[i] != puzzles[0] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("All generated puzzles are identical - randomness may not be working")
	}
}

func TestParsePuzzleString(t *testing.T) {
	tests := []struct {
		name        string
		puzzleStr   string
		expectError bool
		givenCount  int
	}{
		{
			name:        "Valid puzzle with dots",
			puzzleStr:   "......95....64..7......7.1..38.15..25..87..6...72.....7...5...9.5.....2.3.94.....",
			expectError: false,
			givenCount:  26,
		},
		{
			name:        "Valid puzzle with zeros",
			puzzleStr:   "000000950000640070000000701003801500250087006000720000007000500090500000203094000",
			expectError: false,
			givenCount:  26,
		},
		{
			name:        "Invalid length",
			puzzleStr:   "123456789",
			expectError: true,
		},
		{
			name:        "Invalid character",
			puzzleStr:   "......95....64..7......7.1..38.15..25..87..6...72.....7...5...9.5.....2.3.94...X.",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board, err := parsePuzzleString(tt.puzzleStr)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if board == nil {
				t.Fatal("Board is nil")
			}

			// Count given clues
			givenCount := 0
			for i := 0; i < 9; i++ {
				for j := 0; j < 9; j++ {
					if board.Cells[i][j].Given {
						givenCount++
					}
				}
			}

			if givenCount != tt.givenCount {
				t.Errorf("Expected %d given clues, got %d", tt.givenCount, givenCount)
			}
		})
	}
}

func TestGetDefaultDataPath(t *testing.T) {
	path := GetDefaultDataPath()
	expected := filepath.Join("Data", "sudoku_curated_5000.csv")

	if path != expected {
		t.Errorf("Expected default path %s, got %s", expected, path)
	}
}

// Helper function to convert board to string for comparison
func boardToString(board *engine.SudokuBoard) string {
	result := make([]byte, 81)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			val := board.Cells[i][j].Value
			if val == 0 {
				result[i*9+j] = '.'
			} else {
				result[i*9+j] = byte('0' + val)
			}
		}
	}
	return string(result)
}
