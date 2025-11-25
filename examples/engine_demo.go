package main

import (
	"fmt"
	"log"

	"github.com/steven004/goldenfox-sudoku/engine"
	"github.com/steven004/goldenfox-sudoku/generator"
)

func main() {
	fmt.Println("Golden Fox Sudoku - Engine Demo")
	fmt.Println("================================\n")

	// Create generator
	gen, err := generator.NewPreloadedGenerator(generator.GetDefaultDataPath())
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	// Generate a Medium puzzle
	puzzle, err := gen.Generate(engine.Medium)
	if err != nil {
		log.Fatalf("Failed to generate puzzle: %v", err)
	}

	fmt.Println("Generated Medium Puzzle:")
	printBoard(puzzle)
	fmt.Println()

	// Demonstrate board operations
	fmt.Println("Testing Board Operations:")
	fmt.Println("-------------------------")

	// 1. Try a valid move
	fmt.Println("\n1. Testing valid move at [0][0]...")
	if val, _ := puzzle.GetValue(0, 0); val == 0 {
		if puzzle.IsValidMove(0, 0, 5) {
			puzzle.SetValue(0, 0, 5)
			fmt.Println("   ✓ Placed 5 at [0][0]")
		}
	}

	// 2. Add pencil notes
	fmt.Println("\n2. Adding pencil notes to [1][1]...")
	if val, _ := puzzle.GetValue(1, 1); val == 0 {
		puzzle.AddCandidate(1, 1, 3)
		puzzle.AddCandidate(1, 1, 7)
		puzzle.AddCandidate(1, 1, 9)
		candidates, _ := puzzle.GetCandidates(1, 1)
		fmt.Printf("   ✓ Added candidates: %v\n", candidates)
	}

	// 3. Count numbers
	fmt.Println("\n3. Counting numbers on the board...")
	for num := 1; num <= 9; num++ {
		count := puzzle.CountNumber(num)
		fmt.Printf("   Number %d appears %d times", num, count)
		if count == 9 {
			fmt.Print(" (COMPLETE)")
		}
		fmt.Println()
	}

	// 4. Check for conflicts
	fmt.Println("\n4. Checking for conflicts...")
	conflicts := puzzle.FindConflicts()
	if len(conflicts) == 0 {
		fmt.Println("   ✓ No conflicts found")
	} else {
		fmt.Printf("   ✗ Found %d conflicts\n", len(conflicts))
	}

	// 5. Test clone
	fmt.Println("\n5. Testing board clone...")
	clone := puzzle.Clone()
	fmt.Println("   ✓ Board cloned successfully")
	clone.SetValue(2, 2, 8)
	origVal, _ := puzzle.GetValue(2, 2)
	cloneVal, _ := clone.GetValue(2, 2)
	fmt.Printf("   Original [2][2]: %d, Clone [2][2]: %d\n", origVal, cloneVal)
	if origVal != cloneVal {
		fmt.Println("   ✓ Clone is independent")
	}

	// 6. Test reset
	fmt.Println("\n6. Testing board reset...")
	testBoard := puzzle.Clone()
	testBoard.SetValue(3, 3, 9)
	fmt.Println("   Before reset:")
	printBoardRow(testBoard, 3)
	testBoard.Reset()
	fmt.Println("   After reset:")
	printBoardRow(testBoard, 3)
	fmt.Println("   ✓ Non-given cells cleared")

	// 7. Check if solved
	fmt.Println("\n7. Checking if puzzle is solved...")
	if puzzle.IsSolved() {
		fmt.Println("   ✓ Puzzle is solved!")
	} else {
		fmt.Println("   ○ Puzzle is not yet solved")
	}

	fmt.Println("\n================================")
	fmt.Println("Engine Demo Complete!")
}

func printBoard(board *engine.SudokuBoard) {
	for i := 0; i < 9; i++ {
		if i%3 == 0 && i != 0 {
			fmt.Println("------+-------+------")
		}
		for j := 0; j < 9; j++ {
			if j%3 == 0 && j != 0 {
				fmt.Print("| ")
			}
			val := board.Cells[i][j].Value
			if val == 0 {
				fmt.Print(". ")
			} else {
				fmt.Printf("%d ", val)
			}
		}
		fmt.Println()
	}
}

func printBoardRow(board *engine.SudokuBoard, row int) {
	fmt.Print("   ")
	for j := 0; j < 9; j++ {
		val := board.Cells[row][j].Value
		if val == 0 {
			fmt.Print(". ")
		} else {
			fmt.Printf("%d ", val)
		}
	}
	fmt.Println()
}
