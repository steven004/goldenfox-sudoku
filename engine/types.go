package engine

// DifficultyLevel represents the difficulty of a Sudoku puzzle
type DifficultyLevel int

const (
	Beginner DifficultyLevel = iota
	Easy
	Medium
	Hard
	Expert
	FoxGod
)

// String returns the string representation of the difficulty level
func (d DifficultyLevel) String() string {
	switch d {
	case Beginner:
		return "Beginner"
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	case Expert:
		return "Expert"
	case FoxGod:
		return "FoxGod"
	default:
		return "Unknown"
	}
}

// Cell represents a single cell in the Sudoku grid
type Cell struct {
	Value      int          `json:"value"`      // 0 for empty, 1-9 for filled
	Given      bool         `json:"given"`      // true if this is an original clue
	Candidates map[int]bool `json:"candidates"` // pencil notes (candidate numbers)
	IsInvalid  bool         `json:"isInvalid"`  // True if the value conflicts with another cell
}

// NewCell creates a new empty cell
func NewCell() Cell {
	return Cell{
		Value:      0,
		Given:      false,
		Candidates: make(map[int]bool),
	}
}

// Coordinate represents a position on the board
type Coordinate struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// SudokuBoard represents a 9x9 Sudoku grid
type SudokuBoard struct {
	Cells [9][9]Cell `json:"cells"`
}

// NewBoard creates a new empty Sudoku board
func NewBoard() *SudokuBoard {
	board := &SudokuBoard{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			board.Cells[i][j] = NewCell()
		}
	}
	return board
}
