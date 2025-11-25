package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// CellWidget represents a single Sudoku cell
type CellWidget struct {
	widget.BaseWidget
	row, col   int
	mainWindow *MainWindow
	background *canvas.Rectangle
	valueText  *canvas.Text
	tapHandler func()
}

func NewCellWidget(row, col int, mainWindow *MainWindow) *CellWidget {
	cell := &CellWidget{
		row:        row,
		col:        col,
		mainWindow: mainWindow,
	}

	cell.ExtendBaseWidget(cell)
	return cell
}

func (c *CellWidget) CreateRenderer() fyne.WidgetRenderer {
	c.background = canvas.NewRectangle(LightGray)
	c.valueText = canvas.NewText("", Charcoal)
	c.valueText.TextSize = 24
	c.valueText.Alignment = fyne.TextAlignCenter

	return &cellRenderer{
		cell:       c,
		background: c.background,
		valueText:  c.valueText,
		objects:    []fyne.CanvasObject{c.background, c.valueText},
	}
}

func (c *CellWidget) Tapped(_ *fyne.PointEvent) {
	c.mainWindow.selectCell(c.row, c.col)
}

func (c *CellWidget) Update() {
	board := c.mainWindow.gameManager.GetBoard()
	if board == nil {
		return
	}

	val := board.Cells[c.row][c.col].Value
	if val == 0 {
		c.valueText.Text = ""
	} else {
		c.valueText.Text = fmt.Sprintf("%d", val)

		// Color based on whether it's given or user input
		if c.mainWindow.gameManager.IsCellGiven(c.row, c.col) {
			c.valueText.Color = GivenColor
		} else {
			c.valueText.Color = UserColor
		}
	}

	// Update background based on selection/highlighting
	c.updateBackground()

	c.Refresh()
}

func (c *CellWidget) updateBackground() {
	selectedRow, selectedCol, hasSelection := c.mainWindow.gameManager.GetSelectedCell()

	if !hasSelection {
		c.background.FillColor = LightGray
		return
	}

	// Check if this is the selected cell
	if c.row == selectedRow && c.col == selectedCol {
		c.background.FillColor = SelectionHighlight
		return
	}

	// Check if in same row, column, or block as selected cell
	sameRow := c.row == selectedRow
	sameCol := c.col == selectedCol
	sameBlock := (c.row/3 == selectedRow/3) && (c.col/3 == selectedCol/3)

	if sameRow || sameCol || sameBlock {
		c.background.FillColor = PeerHighlight
		return
	}

	c.background.FillColor = LightGray
}

type cellRenderer struct {
	cell       *CellWidget
	background *canvas.Rectangle
	valueText  *canvas.Text
	objects    []fyne.CanvasObject
}

func (r *cellRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.valueText.Resize(size)
	r.valueText.Move(fyne.NewPos(0, 0))
}

func (r *cellRenderer) MinSize() fyne.Size {
	return fyne.NewSize(50, 50)
}

func (r *cellRenderer) Refresh() {
	r.background.Refresh()
	r.valueText.Refresh()
}

func (r *cellRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *cellRenderer) Destroy() {}
