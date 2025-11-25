# Puzzle Generator Module

## Overview
The `generator` package provides a puzzle generation system for Golden Fox Sudoku using pre-loaded puzzles from a curated dataset.

## Implementation
The generator uses the **PreloadedGenerator** which loads 5,000 curated puzzles from a CSV file and provides random selection by difficulty level.

## Usage

### Basic Usage
```go
import (
    "github.com/xin-force/goldenfox-sudoku/engine"
    "github.com/xin-force/goldenfox-sudoku/generator"
)

// Create a new generator
gen, err := generator.NewPreloadedGenerator(generator.GetDefaultDataPath())
if err != nil {
    log.Fatal(err)
}

// Generate a puzzle
board, err := gen.Generate(engine.Medium)
if err != nil {
    log.Fatal(err)
}
```

### Custom Data Path
```go
gen, err := generator.NewPreloadedGenerator("/path/to/puzzles.csv")
```

## Features
- ✅ Loads 1,000 puzzles per difficulty level (5,000 total)
- ✅ Random selection ensures variety
- ✅ Fast generation (no computation, just lookup)
- ✅ Implements `engine.PuzzleGenerator` interface
- ✅ Comprehensive error handling
- ✅ Full unit test coverage

## Difficulty Levels
- **Beginner**: 1,000 puzzles (difficulty 0.0-0.9)
- **Easy**: 1,000 puzzles (difficulty 1.0-2.2)
- **Medium**: 1,000 puzzles (difficulty 2.3-3.9)
- **Hard**: 1,000 puzzles (difficulty 4.0-5.5)
- **Expert**: 1,000 puzzles (difficulty 5.6-8.5)

## Testing
Run the unit tests:
```bash
go test ./generator -v
```

All tests pass:
- ✅ Generator initialization
- ✅ Puzzle loading from CSV
- ✅ Generation for all difficulty levels
- ✅ Randomness verification
- ✅ Puzzle string parsing
- ✅ Error handling

## Example
See `examples/generator_demo.go` for a complete example:
```bash
go run examples/generator_demo.go
```

## Future Enhancements
The modular design allows easy addition of:
- Algorithmic puzzle generation (backtracking)
- Online puzzle API integration
- Custom difficulty algorithms
- Puzzle validation and rating

Simply implement the `engine.PuzzleGenerator` interface with a new generator type.
