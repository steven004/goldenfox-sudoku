package generator

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// PreloadedGenerator implements engine.PuzzleGenerator using pre-loaded puzzles
type PreloadedGenerator struct {
	puzzles map[engine.DifficultyLevel][]string
	rand    *rand.Rand
}

// NewPreloadedGenerator creates a new generator that loads puzzles from a CSV file
func NewPreloadedGenerator(dataPath string) (*PreloadedGenerator, error) {
	gen := &PreloadedGenerator{
		puzzles: make(map[engine.DifficultyLevel][]string),
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	if err := gen.loadPuzzles(dataPath); err != nil {
		return nil, fmt.Errorf("failed to load puzzles: %w", err)
	}

	return gen, nil
}

// loadPuzzles reads the CSV file and organizes puzzles by difficulty
func (g *PreloadedGenerator) loadPuzzles(dataPath string) error {
	file, err := os.Open(dataPath)
	if err != nil {
		return fmt.Errorf("failed to open data file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Read all puzzles
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read record: %w", err)
		}

		// CSV format: puzzle,solution,clues,difficulty,level
		if len(record) != 5 {
			continue
		}

		puzzleStr := record[0]
		levelStr := record[4]

		// Map level string to DifficultyLevel
		var level engine.DifficultyLevel
		switch levelStr {
		case "Beginner":
			level = engine.Beginner
		case "Easy":
			level = engine.Easy
		case "Medium":
			level = engine.Medium
		case "Hard":
			level = engine.Hard
		case "Expert":
			level = engine.Expert
		default:
			continue // Skip unknown levels
		}

		g.puzzles[level] = append(g.puzzles[level], puzzleStr)
	}

	// Verify we have puzzles for all levels
	for level := engine.Beginner; level <= engine.Expert; level++ {
		count := len(g.puzzles[level])
		if count == 0 {
			return fmt.Errorf("no puzzles found for difficulty level: %s", level.String())
		}
	}

	return nil
}

// Generate returns a random puzzle of the specified difficulty
func (g *PreloadedGenerator) Generate(difficulty engine.DifficultyLevel) (*engine.SudokuBoard, error) {
	puzzles, ok := g.puzzles[difficulty]
	if !ok || len(puzzles) == 0 {
		return nil, fmt.Errorf("no puzzles available for difficulty: %s", difficulty.String())
	}

	// Select a random puzzle
	puzzleStr := puzzles[g.rand.Intn(len(puzzles))]

	// Parse the puzzle string into a SudokuBoard
	board, err := parsePuzzleString(puzzleStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse puzzle: %w", err)
	}

	return board, nil
}

// parsePuzzleString converts an 81-character string into a SudokuBoard
// Format: '0' = empty cell, '1'-'9' = given clue
func parsePuzzleString(puzzleStr string) (*engine.SudokuBoard, error) {
	if len(puzzleStr) != 81 {
		return nil, fmt.Errorf("invalid puzzle string length: expected 81, got %d", len(puzzleStr))
	}

	board := engine.NewBoard()

	for i, char := range puzzleStr {
		row := i / 9
		col := i % 9

		if char == '.' || char == '0' {
			// Empty cell
			continue
		}

		if char >= '1' && char <= '9' {
			val := int(char - '0')
			board.Cells[row][col].Value = val
			board.Cells[row][col].Given = true
		} else {
			return nil, fmt.Errorf("invalid character in puzzle string at position %d: %c", i, char)
		}
	}

	return board, nil
}

// GetPuzzleCount returns the number of puzzles available for a given difficulty
func (g *PreloadedGenerator) GetPuzzleCount(difficulty engine.DifficultyLevel) int {
	return len(g.puzzles[difficulty])
}

// GetDefaultDataPath returns the default path to the puzzle data file
func GetDefaultDataPath() string {
	return filepath.Join("Data", "sudoku_curated_5000.csv")
}
