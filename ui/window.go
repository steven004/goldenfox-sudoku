package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/steven004/goldenfox-sudoku/engine"
	"github.com/steven004/goldenfox-sudoku/game"
)

// MainWindow represents the main application window
type MainWindow struct {
	app         fyne.App
	window      fyne.Window
	gameManager *game.GameManager

	// UI Components
	cells       [9][9]*CellWidget
	numberBtns  [9]*widget.Button
	pencilBtn   *widget.Button
	statusLabel *widget.Label

	// Containers
	boardContainer *fyne.Container
}

// NewMainWindow creates a new main window instance
func NewMainWindow(app fyne.App, generator engine.PuzzleGenerator) *MainWindow {
	w := app.NewWindow("Golden Fox Sudoku")

	mw := &MainWindow{
		app:         app,
		window:      w,
		gameManager: game.NewGameManager(generator),
	}

	mw.buildUI()
	mw.startNewGame(engine.Easy) // Start with Easy by default

	w.Resize(fyne.NewSize(900, 700))
	w.CenterOnScreen()

	return mw
}

// ShowAndRun shows the window and runs the application
func (mw *MainWindow) ShowAndRun() {
	mw.window.ShowAndRun()
}

func (mw *MainWindow) buildUI() {
	// 1. Create the Sudoku Board
	mw.createBoard()

	// 2. Create Control Panel (Right side)
	controls := mw.createControls()

	// 3. Create Number Input Bar (Bottom)
	inputBar := mw.createInputBar()

	// 4. Main Layout
	// Use a Border layout: Input bar at bottom, Controls at right, Board in center
	content := container.NewBorder(
		nil,                                    // Top
		inputBar,                               // Bottom
		nil,                                    // Left
		controls,                               // Right
		container.NewCenter(mw.boardContainer), // Center (Board)
	)

	// Add a background container with padding
	finalLayout := container.NewStack(
		canvas.NewRectangle(SoftWhite), // Background color
		container.NewPadded(content),
	)

	mw.window.SetContent(finalLayout)
}

func (mw *MainWindow) createBoard() {
	// Create 9x9 grid of cells
	// We use a Grid layout with 9 columns
	grid := container.NewGridWithColumns(9)

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			cell := NewCellWidget(r, c, mw)
			mw.cells[r][c] = cell
			grid.Add(cell)
		}
	}

	// Wrap grid in a container that preserves aspect ratio and adds spacing
	// For a true Sudoku look, we'd ideally want thicker borders between 3x3 blocks
	// For now, we'll use a simple grid and rely on cell styling
	mw.boardContainer = container.NewPadded(grid)
}

func (mw *MainWindow) createInputBar() *fyne.Container {
	// Horizontal bar with buttons 1-9 and Eraser
	buttons := make([]fyne.CanvasObject, 0, 10)

	// 1-9 Buttons
	for i := 1; i <= 9; i++ {
		val := i
		btn := widget.NewButton(fmt.Sprintf("%d", val), func() {
			mw.handleNumberInput(val)
		})
		mw.numberBtns[i-1] = btn
		buttons = append(buttons, btn)
	}

	// Eraser Button
	eraserBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		mw.handleEraser()
	})
	buttons = append(buttons, eraserBtn)

	return container.NewHBox(container.NewCenter(container.NewHBox(buttons...)))
}

func (mw *MainWindow) createControls() *fyne.Container {
	// Vertical column of controls

	// Title
	title := widget.NewLabel("Golden Fox\nSudoku")
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// Status / Difficulty
	mw.statusLabel = widget.NewLabel("Difficulty: Easy")
	mw.statusLabel.Alignment = fyne.TextAlignCenter

	// Pencil Mode Toggle
	mw.pencilBtn = widget.NewButtonWithIcon("Pencil Mode", theme.DocumentCreateIcon(), func() {
		mw.togglePencilMode()
	})

	// New Game Button
	newGameBtn := widget.NewButtonWithIcon("New Game", theme.ViewRefreshIcon(), func() {
		mw.showNewGameMenu()
	})

	// Restart Button
	restartBtn := widget.NewButtonWithIcon("Restart", theme.HistoryIcon(), func() {
		mw.restartGame()
	})

	// Spacer to push content to top/bottom if needed
	return container.NewVBox(
		title,
		widget.NewSeparator(),
		mw.statusLabel,
		layout.NewSpacer(),
		mw.pencilBtn,
		widget.NewSeparator(),
		newGameBtn,
		restartBtn,
	)
}

// --- Game Actions ---

func (mw *MainWindow) startNewGame(diff engine.DifficultyLevel) {
	err := mw.gameManager.NewGame(diff)
	if err != nil {
		fmt.Println("Error starting game:", err)
		return
	}
	mw.statusLabel.SetText(fmt.Sprintf("Difficulty: %s", diff))
	mw.refreshBoard()
}

func (mw *MainWindow) restartGame() {
	mw.gameManager.RestartGame()
	mw.refreshBoard()
}

func (mw *MainWindow) selectCell(row, col int) {
	mw.gameManager.SelectCell(row, col)
	mw.refreshBoard() // Refresh to show selection highlight
}

func (mw *MainWindow) handleNumberInput(val int) {
	row, col, selected := mw.gameManager.GetSelectedCell()
	if !selected {
		return
	}

	err := mw.gameManager.InputNumber(row, col, val)
	if err != nil {
		// Could show error in status bar or toast
		fmt.Println("Input error:", err)
	}

	mw.refreshBoard()

	if mw.gameManager.IsSolved() {
		mw.showWinDialog()
	}
}

func (mw *MainWindow) handleEraser() {
	row, col, selected := mw.gameManager.GetSelectedCell()
	if !selected {
		return
	}

	mw.gameManager.ClearCell(row, col)
	mw.refreshBoard()
}

func (mw *MainWindow) togglePencilMode() {
	mw.gameManager.TogglePencilMode()
	if mw.gameManager.IsPencilMode() {
		mw.pencilBtn.Importance = widget.HighImportance
	} else {
		mw.pencilBtn.Importance = widget.MediumImportance
	}
	mw.pencilBtn.Refresh()
}

func (mw *MainWindow) refreshBoard() {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			mw.cells[r][c].Update()
		}
	}
}

func (mw *MainWindow) showNewGameMenu() {
	// Simple popup menu for difficulty selection
	// In a real app, this might be a modal dialog
	menu := fyne.NewMenu("New Game",
		fyne.NewMenuItem("Beginner", func() { mw.startNewGame(engine.Beginner) }),
		fyne.NewMenuItem("Easy", func() { mw.startNewGame(engine.Easy) }),
		fyne.NewMenuItem("Medium", func() { mw.startNewGame(engine.Medium) }),
		fyne.NewMenuItem("Hard", func() { mw.startNewGame(engine.Hard) }),
		fyne.NewMenuItem("Expert", func() { mw.startNewGame(engine.Expert) }),
	)

	popUp := widget.NewPopUpMenu(menu, mw.window.Canvas())
	popUp.ShowAtPosition(fyne.NewPos(100, 100)) // Position arbitrarily for now
}

func (mw *MainWindow) showWinDialog() {
	dialog := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabelWithStyle("Congratulations!\nYou Solved the Puzzle!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewButton("Play Again", func() {
				mw.startNewGame(mw.gameManager.GetDifficulty())
				mw.window.Canvas().Overlays().Top().Hide()
			}),
		),
		mw.window.Canvas(),
	)
	dialog.Show()
}
