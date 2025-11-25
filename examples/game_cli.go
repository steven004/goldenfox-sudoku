package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/steven004/goldenfox-sudoku/engine"
	"github.com/steven004/goldenfox-sudoku/game"
	"github.com/steven004/goldenfox-sudoku/generator"
)

func main() {
	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║   Golden Fox Sudoku - CLI Demo    ║")
	fmt.Println("╚════════════════════════════════════╝")
	fmt.Println()

	// Create generator
	gen, err := generator.NewPreloadedGenerator(generator.GetDefaultDataPath())
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	// Create game manager
	gm := game.NewGameManager(gen)

	// Start a new game
	fmt.Println("Starting a new Medium difficulty puzzle...")
	if err := gm.NewGame(engine.Medium); err != nil {
		log.Fatalf("Failed to start game: %v", err)
	}

	fmt.Println()
	printBoard(gm)
	fmt.Println()

	// Interactive game loop
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Commands:")
	fmt.Println("  place <row> <col> <num>  - Place a number (e.g., 'place 0 0 5')")
	fmt.Println("  pencil <row> <col> <num> - Add pencil note")
	fmt.Println("  clear <row> <col>        - Clear a cell")
	fmt.Println("  restart                  - Restart the puzzle")
	fmt.Println("  quit                     - Exit the game")
	fmt.Println()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "place":
			if len(parts) != 4 {
				fmt.Println("Usage: place <row> <col> <num>")
				continue
			}
			row, _ := strconv.Atoi(parts[1])
			col, _ := strconv.Atoi(parts[2])
			num, _ := strconv.Atoi(parts[3])

			gm.SetPencilMode(false)
			if err := gm.InputNumber(row, col, num); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Placed %d at [%d][%d]\n", num, row, col)
				printBoard(gm)

				// Check if solved
				if gm.IsSolved() {
					fmt.Println("\n🎉 Congratulations! You solved the puzzle! 🎉")
					return
				}

				// Check for conflicts
				conflicts := gm.FindConflicts()
				if len(conflicts) > 0 {
					fmt.Printf("⚠️  Warning: %d conflict(s) detected\n", len(conflicts))
				}
			}

		case "pencil":
			if len(parts) != 4 {
				fmt.Println("Usage: pencil <row> <col> <num>")
				continue
			}
			row, _ := strconv.Atoi(parts[1])
			col, _ := strconv.Atoi(parts[2])
			num, _ := strconv.Atoi(parts[3])

			gm.SetPencilMode(true)
			if err := gm.InputNumber(row, col, num); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Added pencil note %d at [%d][%d]\n", num, row, col)
			}

		case "clear":
			if len(parts) != 3 {
				fmt.Println("Usage: clear <row> <col>")
				continue
			}
			row, _ := strconv.Atoi(parts[1])
			col, _ := strconv.Atoi(parts[2])

			if err := gm.ClearCell(row, col); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Cleared cell [%d][%d]\n", row, col)
				printBoard(gm)
			}

		case "restart":
			if err := gm.RestartGame(); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Puzzle restarted!")
				printBoard(gm)
			}

		case "quit":
			fmt.Println("Thanks for playing Golden Fox Sudoku!")
			return

		default:
			fmt.Println("Unknown command. Try: place, pencil, clear, restart, or quit")
		}
	}
}

func printBoard(gm *game.GameManager) {
	board := gm.GetBoard()

	fmt.Println("    0 1 2   3 4 5   6 7 8")
	fmt.Println("  ╔═══════╦═══════╦═══════╗")

	for i := 0; i < 9; i++ {
		if i == 3 || i == 6 {
			fmt.Println("  ╠═══════╬═══════╬═══════╣")
		}

		fmt.Printf("%d ║ ", i)

		for j := 0; j < 9; j++ {
			if j == 3 || j == 6 {
				fmt.Print("║ ")
			}

			val := board.Cells[i][j].Value
			if val == 0 {
				fmt.Print(". ")
			} else {
				if gm.IsCellGiven(i, j) {
					// Given cells in bold (using ANSI codes)
					fmt.Printf("\033[1m%d\033[0m ", val)
				} else {
					fmt.Printf("%d ", val)
				}
			}
		}
		fmt.Println("║")
	}

	fmt.Println("  ╚═══════╩═══════╩═══════╝")

	// Show number completion status
	fmt.Print("Complete: ")
	for num := 1; num <= 9; num++ {
		if gm.IsNumberComplete(num) {
			fmt.Printf("%d ", num)
		}
	}
	fmt.Println()
}
