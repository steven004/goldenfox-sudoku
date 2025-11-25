# Game Manager Module - Phase 1

## Overview
The `game` package provides game state management and coordinates gameplay between the generator, engine, and (future) GUI.

**Phase 1 Status**: ✅ Core gameplay features implemented

## Features Implemented

### Game Lifecycle
- ✅ `NewGame(difficulty)` - Generate and start a new puzzle
- ✅ `RestartGame()` - Reset to initial puzzle state
- ✅ `IsSolved()` - Check if puzzle is complete

### Cell Selection
- ✅ `SelectCell(row, col)` - Track selected cell
- ✅ `GetSelectedCell()` - Get current selection
- ✅ `ClearSelection()` - Clear selection

### Input Modes
- ✅ `TogglePencilMode()` - Switch between number/pencil mode
- ✅ `IsPencilMode()` - Check current mode
- ✅ `SetPencilMode(enabled)` - Set mode explicitly

### Number Input
- ✅ `InputNumber(row, col, val)` - Place number or pencil note
  - Respects pencil mode
  - Auto-clears peer candidates
  - Protects given cells
- ✅ `ClearCell(row, col)` - Remove value/candidates

### Board Access
- ✅ `GetBoard()` - Access current board state
- ✅ `GetCellValue(row, col)` - Get cell value
- ✅ `GetCellCandidates(row, col)` - Get pencil notes
- ✅ `IsCellGiven(row, col)` - Check if cell is given

### Validation
- ✅ `FindConflicts()` - Detect rule violations
- ✅ `IsNumberComplete(val)` - Check if all 9 instances placed

### Metadata
- ✅ `GetDifficulty()` - Get current puzzle difficulty

## Phase 2 Features (Coming Later)
- ⏳ Undo system with limited uses (3 times)
- ⏳ Eraser with limited uses (3 times)
- ⏳ Mistake counter
- ⏳ Timer functionality
- ⏳ Statistics tracking
- ⏳ Save/Load game state

## Usage

### Basic Game Flow
```go
// Create game manager
gen, _ := generator.NewPreloadedGenerator(generator.GetDefaultDataPath())
gm := game.NewGameManager(gen)

// Start a new game
gm.NewGame(engine.Medium)

// Select a cell
gm.SelectCell(0, 0)

// Place a number
gm.InputNumber(0, 0, 5)

// Add pencil notes
gm.SetPencilMode(true)
gm.InputNumber(1, 1, 3)
gm.InputNumber(1, 1, 7)

// Check if solved
if gm.IsSolved() {
    fmt.Println("Puzzle solved!")
}
```

### Restart Game
```go
// Reset to initial state
gm.RestartGame()
```

### Check Conflicts
```go
conflicts := gm.FindConflicts()
if len(conflicts) > 0 {
    fmt.Printf("Found %d conflicts\n", len(conflicts))
}
```

## Testing
Run the unit tests:
```bash
go test ./game -v
```

All tests pass:
- ✅ TestNewGameManager
- ✅ TestNewGame
- ✅ TestRestartGame
- ✅ TestSelectCell
- ✅ TestPencilMode
- ✅ TestInputNumber
- ✅ TestInputNumberPencilMode
- ✅ TestClearCell
- ✅ TestIsNumberComplete
- ✅ TestFindConflicts
- ✅ TestIsSolved
- ✅ TestGetCellValue
- ✅ TestIsCellGiven

## Interactive Demo
Play Sudoku in the terminal:
```bash
go run examples/game_cli.go
```

Commands:
- `place <row> <col> <num>` - Place a number
- `pencil <row> <col> <num>` - Add pencil note
- `clear <row> <col>` - Clear a cell
- `restart` - Restart the puzzle
- `quit` - Exit

## Integration
The Game Manager integrates with:
- **Generator** - Creates new puzzles
- **Engine** - Validates moves and manages board state
- **GUI** (future) - Provides game logic for UI

## Design Principles
1. **Separation of Concerns** - Game logic separate from UI
2. **Immutability** - Given cells are protected
3. **Mode-based Input** - Pencil mode affects InputNumber behavior
4. **Validation** - Uses engine for rule enforcement
5. **Extensibility** - Ready for Phase 2 features
