# Core Engine Module

## Overview
The `engine` package provides the core Sudoku game logic including board operations, validation, and rule enforcement.

## Components

### Data Structures
- **`Cell`** - Represents a single cell with value, given flag, and candidates
- **`SudokuBoard`** - 9x9 grid of cells
- **`Coordinate`** - Row/column position
- **`DifficultyLevel`** - Enum for puzzle difficulty

### Board Operations

#### Basic Operations
```go
board := engine.NewBoard()

// Set/Get values
board.SetValue(row, col, val)
val, err := board.GetValue(row, col)

// Clear a cell
board.ClearCell(row, col)
```

#### Validation
```go
// Check if a move is valid
isValid := board.IsValidMove(row, col, val)

// Find all conflicts on the board
conflicts := board.FindConflicts()

// Check if puzzle is solved
solved := board.IsSolved()
```

#### Pencil Notes (Candidates)
```go
// Add candidate
board.AddCandidate(row, col, val)

// Remove candidate
board.RemoveCandidate(row, col, val)

// Get all candidates for a cell
candidates, err := board.GetCandidates(row, col)

// Auto-remove candidates from peers
board.RemoveCandidateFromPeers(row, col, val)
```

#### Utility Methods
```go
// Count occurrences of a number
count := board.CountNumber(val)

// Clone the board (deep copy)
clone := board.Clone()

// Reset to initial state (clear non-given cells)
board.Reset()
```

## Features Implemented
- ✅ Complete board operations
- ✅ Sudoku rule validation (row, column, block)
- ✅ Conflict detection
- ✅ Pencil notes with auto-clearing
- ✅ Given cell protection
- ✅ Board cloning for undo
- ✅ Reset functionality
- ✅ Number counting for completion indicator
- ✅ Comprehensive error handling
- ✅ Full unit test coverage (14 test functions)

## Testing
Run the unit tests:
```bash
go test ./engine -v
```

All tests pass:
- ✅ TestNewBoard
- ✅ TestSetValue
- ✅ TestGetValue
- ✅ TestIsValidMove
- ✅ TestFindConflicts
- ✅ TestIsSolved
- ✅ TestAddCandidate
- ✅ TestRemoveCandidate
- ✅ TestGetCandidates
- ✅ TestRemoveCandidateFromPeers
- ✅ TestClearCell
- ✅ TestCountNumber
- ✅ TestClone
- ✅ TestReset

## Example
See `examples/engine_demo.go` for a complete demonstration:
```bash
go run examples/engine_demo.go
```

## Integration
The engine is designed to work seamlessly with:
- **Generator** - Provides puzzles with Given flags set
- **Game Manager** - Uses board operations for gameplay
- **GUI** - Displays board state and handles user input

## Design Principles
1. **Immutability of Given Cells** - Cannot modify or clear given clues
2. **Validation First** - `IsValidMove()` checks before `SetValue()`
3. **Auto-clearing** - Setting a value clears candidates automatically
4. **Deep Copying** - `Clone()` creates independent copies for undo
5. **Error Handling** - All operations return errors for invalid inputs
