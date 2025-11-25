package engine

// PuzzleGenerator defines the interface for generating Sudoku puzzles
type PuzzleGenerator interface {
	// Generate creates a new Sudoku puzzle of the specified difficulty
	Generate(difficulty DifficultyLevel) (*SudokuBoard, error)
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
