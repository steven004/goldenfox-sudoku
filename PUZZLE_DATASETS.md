# Sudoku Puzzle Dataset Options

## Recommended: Kaggle 3 Million Sudoku Puzzles

**Source**: [Kaggle - 3 million Sudoku puzzles with ratings](https://www.kaggle.com/datasets/radcliffe/3-million-sudoku-puzzles-with-ratings)

### Details:
- **Size**: 3 million puzzles with solutions
- **Format**: CSV file
- **License**: CC0 Public Domain (free to use)
- **Difficulty Ratings**: 0 to 8.5 (continuous scale)
- **Clues**: 19-31 clues per puzzle (most have 23-26)

### Difficulty Distribution:
- **43%** have difficulty 0 (solvable by simple scanning)
- Ratings go up to 8.5 (very challenging)

### Mapping to Our 5 Levels:
```
Beginner: difficulty 0.0 - 1.0  (simple scanning)
Easy:     difficulty 1.1 - 2.5  (basic techniques)
Medium:   difficulty 2.6 - 4.5  (intermediate)
Hard:     difficulty 4.6 - 6.5  (advanced)
Expert:   difficulty 6.6 - 8.5  (very challenging)
```

### CSV Format:
Each row contains:
- Puzzle string (81 characters, 0 = empty cell)
- Solution string (81 characters)
- Number of clues
- Difficulty rating

Example:
```
004300209005009001070060043006002087190007400050083000600000105003508690042910300,864371259325849761971265843436192587198657432257483916689734125713528694542916378,27,1.2
```

## How to Use

### Option 1: Download and Embed (Recommended for MVP)
1. Download the dataset from Kaggle (requires free account)
2. Extract a subset of puzzles (e.g., 1000 puzzles per difficulty level = 5000 total)
3. Embed them in the Go code as a constant array or separate data file
4. Randomly select from the appropriate difficulty level

### Option 2: Full Dataset with File Loading
1. Download the full dataset
2. Store in `~/.goldenfox/puzzles.csv`
3. Load puzzles on demand from the file
4. More flexible but requires file I/O

## Implementation Approach for MVP

### Step 1: Create a Small Curated Set
Extract ~5000 puzzles (1000 per difficulty) and embed them in Go:

```go
// generator/puzzles_data.go
package generator

var beginnerPuzzles = []string{
    "004300209005009001070060043...",
    "600120384008459072000006005...",
    // ... 998 more
}

var easyPuzzles = []string{
    // ... 1000 puzzles
}

// ... medium, hard, expert
```

### Step 2: Simple Generator Implementation
```go
// generator/preloaded_generator.go
type PreloadedGenerator struct {
    rand *rand.Rand
}

func (g *PreloadedGenerator) Generate(difficulty DifficultyLevel) SudokuBoard {
    var puzzles []string
    switch difficulty {
    case Beginner:
        puzzles = beginnerPuzzles
    case Easy:
        puzzles = easyPuzzles
    // ... etc
    }
    
    // Pick random puzzle
    puzzle := puzzles[g.rand.Intn(len(puzzles))]
    
    // Convert string to SudokuBoard
    return parsePuzzleString(puzzle)
}
```

## Alternative: Smaller Datasets

If you want something simpler to start:

### GitHub JSON Collection
- **Source**: [morcefaster/sudoku.json](https://gist.github.com/morcefaster/sudoku.json)
- Smaller collection, already in JSON format
- Good for quick prototyping

### Generate Your Own Small Set
- Use an online generator to create 50-100 puzzles per level
- Manually curate and test them
- Ensures quality but time-consuming

## Recommendation

**For MVP**: Use Option 1 with the Kaggle dataset
- Download the CSV
- Extract 1000 puzzles per difficulty level (5000 total)
- Embed in Go code
- This gives you plenty of variety without file I/O complexity
- Can upgrade to full dataset or real generator later

**File size**: ~5000 puzzles × 81 chars ≈ 400KB of data (very manageable)
