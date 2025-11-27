# Implementation Plan: Golden Fox Sudoku

## Goal
Implement a modular, extensible Sudoku application in Go with a Fyne-based GUI following the approved design. The application will feature a clean MVC architecture with core game logic, state management, and a rich user interface including statistics tracking, limited undo/eraser actions, and real-time validation.

## User Review Required

> [!IMPORTANT]
> **Limited Action Counters**: Undo and Eraser are limited to 3 uses per game. This is a significant gameplay constraint that affects user experience. Once implemented, changing this limit would require modifying the GameManager logic.

> [!IMPORTANT]
> **User Statistics Persistence**: The design includes user-level statistics (games played, average time, user level). We need to decide on the storage mechanism:
> - Option 1: Store in a local JSON file (e.g., `~/.goldenfox/user_stats.json`)
> - Option 2: Store in the same file as game saves (combined approach)
> - Option 3: Defer this feature to a later phase
> 
> **Recommendation**: Option 1 for simplicity and separation of concerns.

## Proposed Changes

### Core Engine Package

#### [NEW] [engine/types.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/engine/types.go)
Define core data structures:
- `Cell` struct with `Value`, `Given`, and `Candidates` fields
- `SudokuBoard` struct containing 9x9 array of `Cell`
- `Coordinate` struct for row/col pairs
- `DifficultyLevel` enum (Beginner, Easy, Medium, Hard, Expert)

#### [NEW] [engine/board.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/engine/board.go)
Implement board operations:
- `NewBoard()` - Initialize empty board
- `SetValue(row, col, val)` - Place number with validation
- `GetValue(row, col)` - Retrieve cell value
- `IsValidMove(row, col, val)` - Check Sudoku rules
- `FindConflicts()` - Identify rule violations
- `IsSolved()` - Check completion
- `AddCandidate(row, col, val)` - Add pencil note
- `RemoveCandidate(row, col, val)` - Remove pencil note
- `RemoveCandidateFromPeers(row, col, val)` - Auto-clear notes from row/col/block
- `ClearCell(row, col)` - Reset cell
- `CountNumber(val)` - Count instances of a number on the board

#### [NEW] [engine/interfaces.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/engine/interfaces.go)
Define interfaces for extensibility (implementations will be in `generator` package):
- `PuzzleGenerator` interface:
  ```go
  type PuzzleGenerator interface {
      Generate(difficulty DifficultyLevel) SudokuBoard
  }
  ```
- `SudokuSolver` interface:
  ```go
  type SudokuSolver interface {
      Solve(board SudokuBoard) (solution SudokuBoard, solvable bool)
      Hint(board SudokuBoard) MoveHint
      AnalyzeDifficulty(board SudokuBoard) DifficultyRating
  }
  ```

---

### Generator Package (Separate Module)

> **Note**: This package **implements** the interfaces defined in `engine/interfaces.go`. The engine package only defines the contracts, this package fulfills them.

#### [NEW] [generator/generator.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/generator/generator.go)
Implement puzzle generation (implements `engine.PuzzleGenerator` interface):
- `PuzzleGenerator` interface with `Generate(difficulty)` method
- `SimplePuzzleGenerator` struct implementing basic backtracking generation
- Algorithm: Generate complete valid board, then remove cells based on difficulty
- Difficulty levels:
  - **Beginner**: 45-50 clues (very easy, good for learning)
  - **Easy**: 40-44 clues (simple logic only)
  - **Medium**: 35-39 clues (requires some techniques)
  - **Hard**: 30-34 clues (advanced techniques needed)
  - **Expert**: 25-29 clues (very challenging)

#### [NEW] [generator/solver.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/generator/solver.go)
Implement solving logic:
- `SudokuSolver` interface with `Solve()`, `Hint()`, `AnalyzeDifficulty()` methods
- `BacktrackingSolver` struct implementing basic backtracking algorithm
- Used to verify generated puzzles have unique solutions
- Hint and difficulty analysis can be stubbed for future implementation

---

### Game Manager Package

#### [NEW] [game/manager.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/game/manager.go)
Implement game state controller:
- `GameManager` struct with:
  - Current `SudokuBoard`
  - Initial board (for restart)
  - History stack (`[]GameStateMemento`)
  - Session state (selected cell, pencil mode, eraser active)
  - Action counters (undo remaining, eraser remaining)
  - Mistake counter
  - Timer (start time, elapsed)
- Methods:
  - `NewGame(difficulty)` - Generate and start new puzzle
  - `RestartGame()` - Reset to initial state
  - `SelectCell(row, col)`
  - `InputNumber(row, col, val)` - Handle number input based on mode
  - `EraseCell(row, col)` - Clear cell (if eraser uses remain)
  - `Undo()` - Revert last move (if undo uses remain)
  - `TogglePencilMode(on)`
  - `IsNumberComplete(val)` - Check if all 9 instances placed
  - `GetElapsedTime()` - Calculate time since game start
  - `IncrementMistakes()` - Track errors

#### [NEW] [game/memento.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/game/memento.go)
Implement state snapshot:
- `GameStateMemento` struct containing deep copy of board state
- `CreateMemento()` - Snapshot current state
- `RestoreMemento(memento)` - Restore from snapshot

---

### Persistence Package

#### [NEW] [persistence/storage.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/persistence/storage.go)
Implement save/load functionality:
- `Persistence` interface with `SaveGame()` and `LoadGame()` methods
- `JSONPersistence` struct implementing file-based JSON storage
- `GameState` struct for serialization (initial board, current board, notes, metadata)
- Save format matches design specification

### User Data Persistence
- Create `UserData` struct to hold stats and history.
- Implement JSON saving/loading.
- Update `GameManager` to track and save progress.

### History Viewer & Game Loading
- **Backend**:
    - Add `GetHistory` to `GameManager`.
    - Add `LoadGame(id)` to `GameManager`.
    - Add `ParseBoard(string)` to `engine` to reconstruct boards.
- **Frontend**:
    - Create `HistoryModal` component.
    - Implement tabs for "Uncompleted" and "Finished".
    - Add "Load" button to controls.

### Save Functionality
- **Backend**:
    - Add `currentGameID` to `GameManager`.
    - Update `SaveCurrentGame` to handle updates vs new records.
    - Expose `SaveGame` in `App.go`.
- **Frontend**:
    - Connect "Save" button in `App.tsx`.
    - Ensure `Board` component respects `given` flag for styling (already implemented, verify).

#### [NEW] [persistence/user_stats.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/persistence/user_stats.go)
Implement user statistics tracking:
- `UserStats` struct with games played, average times per difficulty, user level
- `LoadUserStats()` and `SaveUserStats()` methods
- Calculate user level based on performance metrics
- Store in `~/.goldenfox/user_stats.json`

---

### GUI Package

#### [NEW] [ui/main.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/ui/main.go)
Main application entry:
- Initialize Fyne application
- Create main window with layout
- Set up board, controls, and statistics panel
- Wire up event handlers

#### [NEW] [ui/cell_widget.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/ui/cell_widget.go)
Custom cell widget:
- Render cell with value or pencil notes
- Handle highlighting (selected, peers, same number, conflicts)
- Handle tap events to select cell
- Support keyboard input when focused

#### [NEW] [ui/board_view.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/ui/board_view.go)
Board rendering:
- 9x9 grid of `CellWidget`s
- Draw block boundaries (thicker lines every 3 cells)
- Update cells when game state changes

#### [NEW] [ui/controls.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/ui/controls.go)
Control panel UI:
- Horizontal number input bar (1-9 buttons)
- Gray out completed numbers using `IsNumberComplete()`
- Pencil toggle button
- Eraser button with remaining uses display (X/3)
- Undo button with remaining uses display (X/3)
- Restart, New Game, Save, Load buttons

#### [NEW] [ui/stats_panel.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/ui/stats_panel.go)
Statistics panel UI:
- Timer display (MM:SS format, updates every second)
- Action counters (Undo X/3, Eraser X/3)
- Mistakes counter
- Game level display
- User level display
- Games played count
- Average time for current difficulty

#### [NEW] [ui/theme.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/ui/theme.go)
Custom theme and colors:
- Define Golden Fox color scheme (fox orange #FF8C00, charcoal, whites)
- Cell highlight colors (selection, peers, same number, conflicts)
- Button styles and spacing

---

### Main Package

#### [NEW] [main.go](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/main.go)
Application entry point:
- Initialize game manager
- Launch Fyne UI
- Set up application lifecycle (save on exit, etc.)

#### [NEW] [go.mod](file:///Users/xin-force/.gemini/antigravity/projects/sudoku/go.mod)
Go module definition:
- Module name: `github.com/xin-force/goldenfox-sudoku` (or appropriate path)
- Dependencies: `fyne.io/fyne/v2`

## Verification Plan

### Automated Tests

#### Engine Tests
```bash
# Run from project root
go test ./engine -v
```
Tests to implement:
- `TestSetValue_ValidMove` - Verify valid number placement
- `TestSetValue_InvalidMove` - Verify rejection of invalid moves
- `TestIsValidMove` - Test rule validation (row, col, block)
- `TestFindConflicts` - Verify conflict detection
- `TestIsSolved` - Test completion detection
- `TestCandidates` - Test pencil note operations
- `TestRemoveCandidateFromPeers` - Verify auto-clearing of notes
- `TestCountNumber` - Verify number counting

#### Generator Tests
```bash
go test ./generator -v
```
Tests to implement:
- `TestGenerate_ValidPuzzle` - Verify generated puzzles are valid
- `TestGenerate_UniqueSolution` - Verify puzzles have unique solutions
- `TestSolver_SolvablePuzzle` - Test solver on known puzzles

#### Game Manager Tests
```bash
go test ./game -v
```
Tests to implement:
- `TestUndo_RestoresState` - Verify undo functionality
- `TestActionLimits` - Verify 3-use limits for undo/eraser
- `TestRestartGame` - Verify restart resets to initial state
- `TestMistakeCounter` - Verify mistake tracking
- `TestIsNumberComplete` - Verify completion detection

#### Persistence Tests
```bash
go test ./persistence -v
```
Tests to implement:
- `TestSaveLoad_RoundTrip` - Verify save/load preserves state
- `TestUserStats_Persistence` - Verify user stats save/load

### Manual Verification

#### GUI Testing
1. **Run the application**:
   ```bash
   go run main.go
   ```

2. **Test Core Gameplay**:
   - Start a new game (Easy/Medium/Hard)
   - Select cells and verify highlighting (row/col/block)
   - Enter numbers using keyboard and number bar
   - Verify number completion graying
   - Toggle pencil mode and add notes
   - Verify auto-removal of notes from peers
   - Create conflicts and verify red highlighting
   - Use eraser (verify 3-use limit)
   - Use undo (verify 3-use limit)
   - Restart game and verify reset

3. **Test Statistics Panel**:
   - Verify timer starts and updates
   - Verify mistake counter increments on conflicts
   - Complete a game and verify stats update (games played, average time)
   - Verify user level calculation

4. **Test Persistence**:
   - Save a game in progress
   - Close application
   - Reopen and load saved game
   - Verify all state is preserved (board, notes, timer, counters)

5. **Test UI/UX**:
   - Verify layout matches design (large board, horizontal number bar, right stats panel)
   - Verify color scheme (fox orange, charcoal)
   - Test window resizing
   - Verify all buttons are functional and labeled correctly

#### Platform Testing
- Test on macOS (primary target)
- Optionally test cross-platform build (Linux/Windows) if time permits
