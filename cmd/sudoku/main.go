package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"github.com/steven004/goldenfox-sudoku/generator"
	"github.com/steven004/goldenfox-sudoku/ui"
)

func main() {
	// Create Fyne application
	myApp := app.NewWithID("com.steven004.goldenfox-sudoku")
	myApp.Settings().SetTheme(&ui.GoldenFoxTheme{})

	// Create generator
	gen, err := generator.NewPreloadedGenerator(generator.GetDefaultDataPath())
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	// Create and show main window
	mainWindow := ui.NewMainWindow(myApp, gen)
	mainWindow.ShowAndRun()
}
