# Curated Sudoku Dataset

## Overview
This dataset contains **5,000 carefully selected Sudoku puzzles** extracted from the Kaggle "3 million Sudoku puzzles with ratings" dataset.

## File Location
`/Users/xin-force/.gemini/antigravity/projects/sudoku/Data/sudoku_curated_5000.csv`

## Distribution
- **Beginner**: 1,000 puzzles (difficulty 0.0 - 0.9)
- **Easy**: 1,000 puzzles (difficulty 1.0 - 2.2)
- **Medium**: 1,000 puzzles (difficulty 2.3 - 3.9)
- **Hard**: 1,000 puzzles (difficulty 4.0 - 5.5)
- **Expert**: 1,000 puzzles (difficulty 5.6 - 8.5)

**Total**: 5,001 rows (including header)

## CSV Format
```csv
puzzle,solution,clues,difficulty,level
```

### Columns:
1. **puzzle**: 81-character string (0 = empty cell, 1-9 = given clue)
2. **solution**: 81-character string (complete solved puzzle)
3. **clues**: Number of given clues (19-31)
4. **difficulty**: Original difficulty rating from Kaggle dataset (0.0-8.5)
5. **level**: Our categorized difficulty level (Beginner/Easy/Medium/Hard/Expert)

### Example Row:
```csv
......95....64..7......7.1..38.15..25..87..6...72.....7...5...9.5.....2.3.94.....,471382956985641273263597418638915742592874361147263895724158639856739124319426587,26,0.0,Beginner
```

## Usage in Go

### Parsing a Puzzle String
```go
func parsePuzzleString(puzzleStr string) SudokuBoard {
    board := NewBoard()
    for i, char := range puzzleStr {
        row := i / 9
        col := i % 9
        if char >= '1' && char <= '9' {
            val := int(char - '0')
            board.SetValue(row, col, val)
            board.Cells[row][col].Given = true
        }
    }
    return board
}
```

### Loading from CSV
```go
func loadPuzzles(filename string, level DifficultyLevel) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    reader.Read() // Skip header
    
    var puzzles []string
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
        
        // record[4] is the level column
        if record[4] == level.String() {
            puzzles = append(puzzles, record[0])
        }
    }
    
    return puzzles, nil
}
```

## Extraction Details
- **Source**: Kaggle 3 million Sudoku puzzles dataset
- **Extraction Date**: 2025-11-25
- **Method**: Random selection from each difficulty range
- **Seed**: 42 (for reproducibility)
- **Script**: `Temp/extract_puzzles.py`

## License
The original Kaggle dataset is licensed under **CC0: Public Domain**, so this curated subset is also free to use without restrictions.
