package main

import (
	"fmt"
	"log"

	"github.com/steven004/goldenfox-sudoku/engine"
	"github.com/steven004/goldenfox-sudoku/generator"
)

func main() {
	// Create a new generator
	gen, err := generator.NewPreloadedGenerator(generator.GetDefaultDataPath())
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	fmt.Println("Golden Fox Sudoku - Generator Demo")
	fmt.Println("===================================\n")

	// Generate and display a puzzle for each difficulty level
	for level := engine.Beginner; level <= engine.Expert; level++ {
		board, err := gen.Generate(level)
		if err != nil {
			log.Fatalf("Failed to generate %s puzzle: %v", level.String(), err)
		}

		fmt.Printf("%s Puzzle:\n", level.String())
		printBoard(board)
		fmt.Println()
	}
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
