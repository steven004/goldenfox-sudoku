package generator

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// PuzzleData holds the puzzle string and its solution
type PuzzleData struct {
	Puzzle          string
	Solution        string
	DifficultyIndex float64
}

// PreloadedGenerator implements engine.PuzzleGenerator using pre-loaded puzzles
type PreloadedGenerator struct {
	puzzles map[engine.DifficultyLevel][]PuzzleData
	rand    *rand.Rand
}

// NewPreloadedGenerator creates a new generator that loads puzzles from CSV data
func NewPreloadedGenerator(data []byte) (*PreloadedGenerator, error) {
	gen := &PreloadedGenerator{
		puzzles: make(map[engine.DifficultyLevel][]PuzzleData),
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	if err := gen.loadPuzzles(data); err != nil {
		return nil, fmt.Errorf("failed to load puzzles: %w", err)
	}

	return gen, nil
}

// loadPuzzles reads the CSV data and organizes puzzles by difficulty
func (g *PreloadedGenerator) loadPuzzles(data []byte) error {
	reader := csv.NewReader(bytes.NewReader(data))

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
		solutionStr := record[1]
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

		// Parse difficulty index from the 'difficulty' column (index 3)
		// Assuming the CSV format is: puzzle,solution,clues,difficulty_index,level_name
		// Wait, the user said "level" column is the float index (e.g. "1.2")?
		// Let's re-read the CSV format comment: "puzzle,solution,clues,difficulty,level"
		// Usually 'difficulty' is the float (e.g. 1.2) and 'level' is the name (e.g. Easy).
		// Let's check the record[3] content.

		var diffIndex float64
		if val, err := strconv.ParseFloat(record[3], 64); err == nil {
			diffIndex = val
		}

		data := PuzzleData{
			Puzzle:          puzzleStr,
			Solution:        solutionStr,
			DifficultyIndex: diffIndex,
		}
		g.puzzles[level] = append(g.puzzles[level], data)
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
func (g *PreloadedGenerator) Generate(difficulty engine.DifficultyLevel, extraClues int) (*engine.SudokuBoard, float64, error) {
	// Determine the pool to use
	var poolLevel engine.DifficultyLevel

	switch difficulty {
	case engine.Beginner:
		poolLevel = engine.Beginner
	case engine.Easy:
		poolLevel = engine.Easy
	case engine.Medium:
		poolLevel = engine.Medium
	case engine.Hard:
		poolLevel = engine.Hard
	case engine.Expert:
		poolLevel = engine.Expert
	case engine.FoxGod:
		poolLevel = engine.Expert // FoxGod uses Expert pool
	default:
		return nil, 0, fmt.Errorf("unknown difficulty: %s", difficulty.String())
	}

	puzzles, ok := g.puzzles[poolLevel]
	if !ok || len(puzzles) == 0 {
		return nil, 0, fmt.Errorf("no puzzles available for difficulty: %s", poolLevel.String())
	}

	// Select a random puzzle
	index := g.rand.Intn(len(puzzles))
	puzzleData := puzzles[index]

	fmt.Printf("Generator: Requesting %s difficulty (Pool: %s, Extra Clues: %d)\n", difficulty.String(), poolLevel.String(), extraClues)
	fmt.Printf("Generator: Selected puzzle index %d (Diff: %.2f)\n", index, puzzleData.DifficultyIndex)

	// Parse the puzzle string into a SudokuBoard
	board, err := parsePuzzleString(puzzleData.Puzzle)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse puzzle: %w", err)
	}

	// Add extra clues if needed
	if extraClues > 0 {
		if err := g.addExtraClues(board, puzzleData.Solution, extraClues); err != nil {
			fmt.Printf("Warning: Failed to add extra clues: %v\n", err)
		}
	}

	// Count clues for verification
	clueCount := 0
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board.Cells[i][j].Given {
				clueCount++
			}
		}
	}
	fmt.Printf("Generator: Puzzle has %d clues (Original + Extra)\n", clueCount)

	return board, puzzleData.DifficultyIndex, nil
}

// addExtraClues reveals N random empty cells using the solution
func (g *PreloadedGenerator) addExtraClues(board *engine.SudokuBoard, solutionStr string, count int) error {
	if len(solutionStr) != 81 {
		return fmt.Errorf("invalid solution string length")
	}

	// Find all empty cells
	type coord struct{ r, c int }
	var emptyCells []coord

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if board.Cells[r][c].Value == 0 {
				emptyCells = append(emptyCells, coord{r, c})
			}
		}
	}

	// Shuffle empty cells
	g.rand.Shuffle(len(emptyCells), func(i, j int) {
		emptyCells[i], emptyCells[j] = emptyCells[j], emptyCells[i]
	})

	// Reveal the first N cells
	revealed := 0
	for _, cell := range emptyCells {
		if revealed >= count {
			break
		}

		// Get value from solution string
		idx := cell.r*9 + cell.c
		valChar := solutionStr[idx]
		if valChar < '1' || valChar > '9' {
			continue // Skip invalid chars in solution
		}
		val := int(valChar - '0')

		// Set value on board
		board.Cells[cell.r][cell.c].Value = val
		board.Cells[cell.r][cell.c].Given = true
		revealed++
	}

	return nil
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
