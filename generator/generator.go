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

// PuzzleData holds the puzzle string and its solution
type PuzzleData struct {
	Puzzle   string
	Solution string
}

// PreloadedGenerator implements engine.PuzzleGenerator using pre-loaded puzzles
type PreloadedGenerator struct {
	puzzles map[engine.DifficultyLevel][]PuzzleData
	rand    *rand.Rand
}

// NewPreloadedGenerator creates a new generator that loads puzzles from a CSV file
func NewPreloadedGenerator(dataPath string) (*PreloadedGenerator, error) {
	gen := &PreloadedGenerator{
		puzzles: make(map[engine.DifficultyLevel][]PuzzleData),
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

		data := PuzzleData{
			Puzzle:   puzzleStr,
			Solution: solutionStr,
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
func (g *PreloadedGenerator) Generate(difficulty engine.DifficultyLevel) (*engine.SudokuBoard, error) {
	// Determine the pool to use and extra clues to add
	var poolLevel engine.DifficultyLevel
	var extraClues int

	switch difficulty {
	case engine.Beginner:
		poolLevel = engine.Beginner
		extraClues = 10
	case engine.Easy:
		poolLevel = engine.Easy
		extraClues = 8
	case engine.Medium:
		poolLevel = engine.Medium
		extraClues = 6
	case engine.Hard:
		poolLevel = engine.Hard
		extraClues = 4
	case engine.Expert:
		poolLevel = engine.Expert
		extraClues = 2
	case engine.FoxGod:
		poolLevel = engine.Expert // FoxGod uses Expert pool with no help
		extraClues = 0
	default:
		return nil, fmt.Errorf("unknown difficulty: %s", difficulty.String())
	}

	puzzles, ok := g.puzzles[poolLevel]
	if !ok || len(puzzles) == 0 {
		return nil, fmt.Errorf("no puzzles available for difficulty: %s", poolLevel.String())
	}

	// Select a random puzzle
	index := g.rand.Intn(len(puzzles))
	puzzleData := puzzles[index]

	fmt.Printf("Generator: Requesting %s difficulty (Pool: %s, Extra Clues: %d)\n", difficulty.String(), poolLevel.String(), extraClues)
	fmt.Printf("Generator: Selected puzzle index %d\n", index)

	// Parse the puzzle string into a SudokuBoard
	board, err := parsePuzzleString(puzzleData.Puzzle)
	if err != nil {
		return nil, fmt.Errorf("failed to parse puzzle: %w", err)
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

	return board, nil
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

// GetDefaultDataPath returns the default path to the puzzle data file
func GetDefaultDataPath() string {
	return filepath.Join("Data", "sudoku_curated_5000.csv")
}
