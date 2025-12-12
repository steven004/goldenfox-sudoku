package engine

// PuzzleGenerator defines the interface for generating Sudoku puzzles
type PuzzleGenerator interface {
	// Generate creates a new Sudoku puzzle of the specified difficulty with extra clues
	// seed: A string used for deterministic generation (e.g. "daily-2025-12-12"). If empty, behavior is implementation specific (usually random).
	Generate(difficulty DifficultyLevel, extraClues int, seed string) (*SudokuBoard, float64, error)
}

// SudokuSolver defines the interface for solving Sudoku puzzles
type SudokuSolver interface {
	// Solve attempts to solve the given puzzle
	Solve(board *SudokuBoard) (*SudokuBoard, bool)

	// Hint provides a suggested move for the current board state
	Hint(board *SudokuBoard) (*Coordinate, int, error)

	// AnalyzeDifficulty estimates the difficulty of a puzzle
	AnalyzeDifficulty(board *SudokuBoard) (float64, error)
}
