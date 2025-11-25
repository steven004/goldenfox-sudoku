# Golden Fox Sudoku

A modular, extensible Sudoku application written in Go with a Fyne-based GUI.

## Features (Planned)
- 🦊 Beautiful GUI with fox-themed design
- 🎮 5 difficulty levels (Beginner, Easy, Medium, Hard, Expert)
- ✏️ Pencil notes with auto-clearing
- ↩️ Limited undo/eraser (3 uses per game)
- 📊 Statistics tracking and user progress
- 💾 Save/Load game functionality
- ⏱️ Timer and mistake counter

## Current Status
✅ **Puzzle Generator Module** - Complete with 5,000 curated puzzles

## Project Structure
```
goldenfox-sudoku/
├── engine/          # Core Sudoku logic and interfaces
├── generator/       # Puzzle generation (pre-loaded puzzles)
├── Data/            # 5,000 curated puzzles dataset
├── Design/          # Design documents and GUI mockups
├── examples/        # Example programs
└── go.mod           # Go module definition
```

## Installation
```bash
git clone https://github.com/steven004/goldenfox-sudoku.git
cd goldenfox-sudoku
go mod tidy
```

## Usage

### Run Generator Demo
```bash
go run examples/generator_demo.go
```

### Run Tests
```bash
go test ./...
```

## Module Name
```
github.com/steven004/goldenfox-sudoku
```

## Development

### Implemented
- ✅ Engine types and interfaces
- ✅ Puzzle generator with 5,000 curated puzzles
- ✅ Comprehensive unit tests

### In Progress
- 🚧 Core engine (board operations, validation)
- 🚧 Game manager (state management, undo)
- 🚧 GUI with Fyne

## License
MIT License (or your preferred license)

## Author
Steven (steven004)
