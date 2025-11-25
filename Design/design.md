# Sudoku App Design Document

## 1. Introduction
This document outlines a modular architecture for a Sudoku application written in Go with a Fyne-based GUI. The goal is to ensure the solution is highly modular, extensible, and maintainable, allowing future features to be added with minimal refactoring. We adopt an MVC-inspired separation of concerns, loosely coupling the graphical interface from the core puzzle logic.

## 2. Features and Requirements

### Core Features
![Golden Fox Sudoku GUI Design](./fox_sudoku_final.png)
- **Cell Selection & Highlighting**: Clicking an empty cell highlights its entire row, column, and 3×3 subgrid. Clicking a filled cell highlights all cells with the same number.
- **Pencil Notes**: A toggleable pencil mode allows entry of candidate numbers. When a cell is filled, the same candidate is automatically removed from all peer cells (row, column, block).
- **Number Completion Indicator**: When all 9 instances of a number are placed on the board, that number button in the horizontal input bar is grayed out (disabled).
- **Eraser Tool**: Removes entries. If a cell has a user-filled number, it clears it. If it has notes, it clears notes. **Limited to 3 uses per game**.
- **Undo**: Step-by-step undo of moves (placements, notes, erasures) using a history stack. **Limited to 3 uses per game**.
- **Real-Time Validation**: Immediate highlighting of rule conflicts (duplicates in row/col/block).
- **Save/Load Game**: Persist current game state (including values and notes) to disk.

### Statistics & Progress Panel
- **Timer**: Real-time game timer showing elapsed time.
- **Restart Button**: Restart the current puzzle from the beginning.
- **Action Limits**: Display remaining undo (max 3) and eraser (max 3) uses.
- **Mistake Counter**: Track number of mistakes (wrong placements, conflicts created).
- **Game Level**: Display current puzzle difficulty (Beginner/Easy/Medium/Hard/Expert).
- **User Level**: User's skill level based on performance.
- **Games Played**: Total number of games completed by the user.
- **Average Time**: Average completion time for puzzles of the current difficulty level.

### Constraints
- **Language**: Go.
- **GUI Toolkit**: Fyne (targeting macOS/cross-platform).
- **Scope**: Single-player. No user login/AI solver initially, but architecture must support future addition.

## 3. Architecture Overview

The application is organized into distinct layers corresponding to Model-View-Controller (MVC):

1.  **Core Sudoku Engine (Model)**: The puzzle logic layer. Encapsulates board state, rules, and algorithms. Knows nothing of the GUI.
2.  **Game State Management (Controller)**: Orchestrates gameplay. Mediates between UI events and the core engine. Handles input, updates puzzle state, manages pencil mode, history (undo), and persistence.
3.  **User Interface (View)**: Fyne-based GUI. Renders the board and controls. Captures interactions and delegates intent to the controller. Observes game state changes to update visuals.

### Inter-module Interactions
1.  **User Input Flow**: UI calls Controller methods (e.g., `InputNumber`, `Undo`). The UI does not modify puzzle state directly.
2.  **State Update**: Controller updates Engine. Engine enforces rules and returns results/conflicts. Controller handles side-effects (e.g., clearing peer notes).
3.  **UI Refresh (Observer)**: Controller/Model notifies UI of changes. The UI subscribes to state change events to redraw.
4.  **Persistence**: Controller uses a Persistence module to perform file I/O, isolating file formats from the rest of the app.

## 4. Module Details

### 4.1. Core Sudoku Engine (Model)
**Package**: `sudokuengine`

This module is the heart of the logic. It can be used independently of the GUI (e.g., for testing or a CLI version).

**Data Model**:
-   `SudokuBoard`: Represents the 9x9 grid.
-   `Cell`:
    -   `Value`: int (0 for empty, 1-9 for filled).
    -   `Given`: bool (true if original clue, immutable).
    -   `Candidates`: map[int]bool or bitmask (for pencil notes).

**Key Methods**:
-   `IsValidMove(row, col, val) bool`: Checks if placing `val` violates Sudoku rules.
-   `FindConflicts() []Coordinate`: Scans board for rule violations (useful for red highlighting).
-   `IsSolved() bool`: Checks if puzzle is complete and valid.
-   `SetValue(row, col, val) error`: Places a number. Rejects if cell is `Given`.
-   `AddCandidate(row, col, val)` / `RemoveCandidate(...)`: Manages pencil notes.
-   `ClearCell(row, col)`: Removes user value or notes.

**Interfaces (for extensibility)**:
-   `PuzzleGenerator`:
    ```go
    type PuzzleGenerator interface {
        Generate(difficulty DifficultyLevel) SudokuBoard
    }
    ```
    Allows swapping generation algorithms (random, backtracking, external service).

-   `SudokuSolver`:
    ```go
    type SudokuSolver interface {
        Solve(board SudokuBoard) (solution SudokuBoard, solvable bool)
        Hint(board SudokuBoard) MoveHint
        AnalyzeDifficulty(board SudokuBoard) DifficultyRating
    }
    ```
    Supports future AI hints or difficulty analysis.

### 4.2. Game State Management (Controller)
**Component**: `GameManager`

**Responsibilities**:
-   Holds the current `SudokuBoard`.
-   Manages Session State: Selected cell coordinates, `PencilMode` (bool), `EraserActive` (bool).
-   Manages History: A stack of `GameStateMemento` objects for Undo.

**Controller API (called by UI)**:
-   `SelectCell(row, col)`: Updates selected cell state.
-   `InputNumber(row, col, val)`:
    -   If `EraserActive`: Calls `EraseCell` (if eraser uses remaining).
    -   If `PencilMode`: Calls `AddCandidate`.
    -   Else: Calls `SetValue`. Handles auto-removal of peer candidates.
-   `IsNumberComplete(val) bool`: Returns true if all 9 instances of `val` are placed on the board.
-   `EraseCell(row, col)`: Clears value or notes. Decrements eraser uses counter.
-   `TogglePencilMode(on)`: Updates mode flag.
-   `Undo()`: Reverts last move. Decrements undo uses counter.
-   `RestartGame()`: Resets the current puzzle to its initial state.
-   `NewGame(difficulty)`: Generates new puzzle using `PuzzleGenerator`. Resets timer and counters.
-   `SaveGame(path)` / `LoadGame(path)`: Invokes Persistence layer.

**Undo System (Memento Pattern)**:
-   **State Snapshot**: `GameStateMemento` contains a deep copy of the grid values and notes.
-   **Process**: Before any state change, push current snapshot to `history` stack. On Undo, pop last snapshot and restore `SudokuBoard` state.
-   **Efficiency**: Storing full 9x9 snapshots is memory-cheap and robust against complex state changes (like auto-clearing notes).

### 4.3. GUI (View)
**Toolkit**: Fyne

**Components**:
-   **Main Window**: Grid layout for board, control panel for buttons.
-   **Board**: 9x9 grid of `CellWidget`s.
    -   **CellWidget**: Custom widget.
        -   Visuals: Background color (highlight/conflict), Main text (Value), Mini-grid text (Notes).
        -   Input: Handles `Tapped` events to call `SelectCell`.
-   **Controls** (bottom area):
    -   **Number Input Bar**: A horizontal row of 9 buttons (1-9) for easy number entry. Buttons are grayed out when all 9 instances are placed.
    -   Toggle Buttons: Pencil, Eraser (shows remaining uses: X/3).
    -   Action Buttons: Undo (shows remaining uses: X/3), Restart, Save, Load, New Game.
-   **Statistics Panel** (right side):
    -   **Timer**: Live elapsed time display (MM:SS).
    -   **Restart Button**: Reset current puzzle.
    -   **Action Counters**: Undo remaining (X/3), Eraser remaining (X/3).
    -   **Mistakes**: Count of errors made.
    -   **Game Level**: Current difficulty (Beginner/Easy/Medium/Hard/Expert).
    -   **User Level**: Skill rating (Beginner/Intermediate/Advanced/Expert).
    -   **Games Played**: Total completed games.
    -   **Average Time**: Average completion time for current difficulty.

**Event Handling & Highlighting**:
-   **Selection**: When `SelectCell` is called, UI computes highlights:
    -   Selected cell: Distinct border/color.
    -   Peers (Row/Col/Block): Light highlight.
    -   Same Number: Highlight all cells with value == selected value.
-   **Conflicts**: After moves, if Engine reports conflicts, mark those cells (e.g., red text/bg).
-   **Observer**: UI subscribes to `GameManager` events (e.g., `OnBoardUpdate`). When triggered, UI redraws affected cells or whole board.

### 4.4. Persistence Layer
**Component**: `Persistence` / `StorageManager`

**Format**: JSON. Human-readable and easy to debug.

**JSON Structure Example**:
```json
{
  "initial": [[5,0,0,...], ...],  // To identify Givens
  "current": [[5,3,0,...], ...],  // Current values
  "notes": {
    "0,1": [2,4],
    "8,8": [1,9]
  },
  "difficulty": "Medium"
}
```

**Interface**:
```go
type Persistence interface {
    SaveGame(state GameState, path string) error
    LoadGame(path string) (GameState, error)
}
```
This abstraction allows future support for cloud saves or databases without changing the Game Manager.

## 5. Future Extensions
The modular design specifically enables:
-   **Analytics**: Add an observer to `GameManager` to log events (moves, wins) without touching UI code.
-   **Multiplayer**: Implement a `NetworkGameController` that sends/receives moves. The UI remains unchanged.
-   **AI/Hints**: Implement the `SudokuSolver` interface to provide "Next Move" hints.
-   **Difficulty Analysis**: Use `SudokuSolver.AnalyzeDifficulty` to rate generated puzzles.
-   **Internationalization**: UI layer can load strings from resource files; Engine/Controller are language-agnostic.
