package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// SimplePuzzleGenerator generates puzzles using basic backtracking
type SimplePuzzleGenerator struct {
	rand *rand.Rand
}

// NewSimpleGenerator creates a new simple generator
func NewSimpleGenerator() *SimplePuzzleGenerator {
	return &SimplePuzzleGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate creates a new puzzle with the specified difficulty
// For now, it returns a hardcoded valid puzzle for testing purposes
// or generates a full board and removes cells (if implemented)
func (g *SimplePuzzleGenerator) Generate(difficulty engine.DifficultyLevel) (*engine.SudokuBoard, error) {
	// TODO: Implement actual generation logic
	// For now, return a simple hardcoded puzzle or a valid empty board with some clues

	board := engine.NewBoard()

	// Fill diagonal blocks (independent) to ensure validity before solving
	// g.fillDiagonal(board) <--- REMOVED: This conflicts with hardcoded clues

	// Define a pool of Easy puzzles
	// Format: [row, col, value]
	puzzlePool := [][][3]int{
		// Puzzle 1
		{
			{0, 0, 5}, {0, 1, 3}, {0, 4, 7},
			{1, 0, 6}, {1, 3, 1}, {1, 4, 9}, {1, 5, 5},
			{2, 1, 9}, {2, 2, 8}, {2, 7, 6},
			{3, 0, 8}, {3, 4, 6}, {3, 8, 3},
			{4, 0, 4}, {4, 3, 8}, {4, 5, 3}, {4, 8, 1},
			{5, 0, 7}, {5, 4, 2}, {5, 8, 6},
			{6, 1, 6}, {6, 7, 2}, {6, 8, 8},
			{7, 3, 4}, {7, 4, 1}, {7, 5, 9}, {7, 8, 5},
			{8, 4, 8}, {8, 7, 7}, {8, 8, 9},
		},
		// Puzzle 2
		{
			{0, 3, 2}, {0, 4, 6}, {0, 6, 7}, {0, 8, 1},
			{1, 0, 6}, {1, 1, 8}, {1, 4, 7}, {1, 7, 9},
			{2, 0, 1}, {2, 1, 9}, {2, 5, 4}, {2, 6, 5},
			{3, 0, 8}, {3, 1, 2}, {3, 3, 1}, {3, 7, 4},
			{4, 2, 4}, {4, 3, 6}, {4, 5, 2}, {4, 6, 9},
			{5, 1, 5}, {5, 5, 3}, {5, 7, 2}, {5, 8, 8},
			{6, 2, 9}, {6, 3, 3}, {6, 7, 7}, {6, 8, 4},
			{7, 1, 4}, {7, 4, 5}, {7, 7, 3}, {7, 8, 6},
			{8, 0, 7}, {8, 2, 3}, {8, 4, 1}, {8, 5, 8},
		},
		// Puzzle 3
		{
			{0, 0, 1}, {0, 3, 4}, {0, 4, 8}, {0, 5, 9}, {0, 8, 6},
			{1, 0, 7}, {1, 1, 3}, {1, 5, 2}, {1, 8, 4},
			{2, 2, 9}, {2, 3, 7}, {2, 4, 1}, {2, 6, 8},
			{3, 0, 5}, {3, 4, 7}, {3, 5, 3}, {3, 7, 9},
			{4, 0, 8}, {4, 3, 2}, {4, 5, 1}, {4, 8, 3},
			{5, 1, 9}, {5, 3, 5}, {5, 4, 6}, {5, 8, 7},
			{6, 2, 4}, {6, 4, 2}, {6, 5, 7}, {6, 6, 3},
			{7, 0, 3}, {7, 3, 8}, {7, 7, 6}, {7, 8, 9},
			{8, 0, 2}, {8, 3, 9}, {8, 4, 5}, {8, 5, 4}, {8, 8, 1},
		},
	}

	// Randomly select a puzzle from the pool
	puzzleIndex := g.rand.Intn(len(puzzlePool))
	selectedClues := puzzlePool[puzzleIndex]

	fmt.Printf("Generating puzzle: Selected puzzle #%d from pool\n", puzzleIndex+1)

	for _, clue := range selectedClues {
		r, c, v := clue[0], clue[1], clue[2]
		board.Cells[r][c].Value = v
		board.Cells[r][c].Given = true
	}

	return board, nil
}
