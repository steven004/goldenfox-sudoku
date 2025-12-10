# Sudoku App Design Document

## 1. Introduction
This document outlines the architecture for a Sudoku application built using **Wails** (Go Backend + React/TypeScript Frontend). The goal is to ensure a responsive, modern user experience with a robust, modular backend. We adopt a strict separation of concerns where Go handles the business logic and persistence, while React handles the presentation and interactive state.

## 2. Features and Requirements

### Core Features
![Golden Fox Sudoku GUI Design](./fox_sudoku_final.png)
-   **Cell Selection & Highlighting**: Clicking an empty cell highlights its entire row, column, and 3×3 subgrid. Clicking a filled cell highlights all cells with the same number.
-   **Pencil Notes**: A toggleable pencil mode allows entry of candidate numbers. When a cell is filled, the same candidate is automatically removed from all peer cells (row, column, block).
-   **Transient Conflict Feedback**:
    -   **Invalid Moves**: Number inputs that violate Sudoku rules are rejected by the backend. The frontend displays the invalid number in **Red** momentarily (transient error) accompanied by an error sound.
    -   **Pencil Conflicts**: Existing pencil marks that conflict with placed numbers are visually highlighted in red (frontend logic).
-   **Eraser Tool**: Removes entries. If a cell has a user-filled number, it clears it. If it has notes, it clears notes. **Limited to 3 uses per game**.
-   **Undo**: Step-by-step undo of moves (placements, notes, erasures) using a history stack. **Limited to 3 uses per game**.
-   **Save/Load Game**: Auto-saves current game state (including values and notes) to disk (JSON).
-   **Sound Effects**: Audio feedback for interactions (Click, Pop, Scratch, Error, Win).

### Statistics & Progress Panel
-   **Timer**: Real-time game timer (Frontend driven for display, Backend tracked for stats).
-   **Restart Button**: Restart the current puzzle from the beginning.
-   **Action Limits**: Display remaining undo (max 3) and eraser (max 3) uses.
-   **History Tracking**: Track wins, losses, average times, and best times per difficulty.
-   **Game Level**: Display current puzzle difficulty (Beginner/Easy/Medium/Hard/Expert/FoxGod).
-   **User Level**: User's skill level based on performance.

### Constraints
-   **Backend**: Go (Wails).
-   **Frontend**: React (Vite, TypeScript, Tailwind/CSS).
-   **Target**: Desktop (macOS first, cross-platform capable).
-   **Scope**: Single-player, persistent local data.

## 3. Architecture Overview

The application follows the **Wails Architecture**:

1.  **Frontend (View)**: React Application.
    -   Renders the board state.
    -   Handles immediate user interactions (hover, click, selection).
    -   Manages transient visual states (animations, transient errors).
    -   Plays sound effects.
    -   Calls Backend methods via Wails bindings.
2.  **Backend (Controller/Model)**: Go Application.
    -   **GameManager**: Coordinator. Manages application lifecycle, user data loading/saving, and sessions.
    -   **GameSession**: Active Game Logic. Encapsulates the specific state of a playing session (Board, History, Timer, Limits).
    -   **Sudoku Engine**: Core domain logic (Board structure, Validation rules, Generator).

### Inter-module Interactions
1.  **Input Flow**: User clicks a number -> Frontend calls `InputNumber(row, col, val)` (Go).
2.  **Validation**: `GameSession` validates the move.
    -   **If Valid**: Updates Board, History, and Auto-saves. Returns `nil`.
    -   **If Invalid**: Returns `error`.
3.  **Feedback**:
    -   **Success**: Frontend calls `GetGameState` to refresh the board.
    -   **Error**: Frontend catches the error and triggers a "Transient Error State" (Red number flash + sound).

## 4. Module Details

### 4.1. Core Sudoku Engine (Model)
**Package**: `engine`

-   `SudokuBoard`: Represents the 9x9 grid.
-   `Cell`: Value, Given (bool), Candidates (map[int]bool), IsInvalid (legacy/unused).
-   **Key Methods**:
    -   `IsValidMove(row, col, val) bool`: Pure logic check.
    -   `FindConflicts() []Coordinate`: Legacy conflict finder (unused in active play).
    -   `IsSolved() bool`.
    -   `PuzzleGenerator`: Generates valid puzzles with difficulty grading.

### 4.2. Game Management (Backend)
**Package**: `game`

#### `GameManager` (Coordinator)
-   **Responsibilities**:
    -   Loads/Saves `UserData` (Stats, History).
    -   Manages the active `GameSession` pointer.
    -   **Thread Safety**: Uses `sync.RWMutex` to protect the session from concurrent access.
    -   **API**: Exposes simplified methods to Wails (`InputNumber`, `NewGame`, `GetGameState`).

#### `GameSession` (Active Logic)
-   **Responsibilities**:
    -   Holds `currentBoard` and `initialBoard`.
    -   Holds `HistoryManager` (Undo stack).
    -   Tracks `StartTime` / `EndTime` (for duration calc).
    -   Tracks `EraseCount` / `UndoCount`.
-   **Behavior**:
    -   `InputNumber`: strictly rejects invalid moves (Conflict-free state).
    -   `ToggleCandidate`: strictly rejects invalid notes.

### 4.3. User Interface (Frontend)
**Tech**: React, TypeScript, Vite.
**Path**: `frontend/src`

-   **Components**:
    -   `App`: Main layout, Headers, Sidebar.
    -   `Board`: Renders 9x9 grid.
    -   `Cell`: Individual cell rendering (Value, Candidates, Colors).
        -   **Visuals**: Handles "Same Number" highlight, "Peer" highlight, "Error" red state.
    -   `Controls`: Number pad, Tool toggles (Pen/Pencil).
-   **Hooks**:
    -   `useGameLogic`: The brain of the frontend.
        -   Maintains `gameState`.
        -   Maintains `transientError` (local temporary state for invalid inputs).
        -   Maintains `transientSound` logic.
    -   `useSound`: Manages AudioContext and sound assets.

### 4.4. Persistence Layer
**Structure**: `UserData` struct in `game/userdata.go`.
**Storage**: `user_data.json`. Defaults to OS Config setup, but supports local file override if present.

**Smart ID Format**:
12-Digits: `DPIIIIIIIIII`
- **D**: Difficulty Index + 1 (1-6)
- **P**: Progress + 4 (Interval 2-8)
- **I**: Game Index (10-digit padded sequence)
*Example*: `360000000015` (Medium, +2 Progress, 15th Game)

**JSON Structure**:
```json
{
  "stats": { ... },
  "completed_history": [
    {
      "id": "360000000015",
      "predefined": "...",
      "final_state": "...",
      "is_solved": true,
      "time_elapsed": 120000000000,
      "difficulty": 2
    }
  ],
  "pending_history": [
    {
      "id": "420000000001",
      "predefined": "...",
      "final_state": "...",
      "is_solved": false,
      "time_elapsed": 5000000000,
      "difficulty": 3
    }
  ]
}
```

## 5. Deployment
-   **Build**: `wails build`.
-   **Tags**: Use `-tags production` to suppress debug logs.
-   **Assets**: Frontend assets compiled by Vite, embedded by Wails.
