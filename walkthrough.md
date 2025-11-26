# Golden Fox Sudoku - GUI Verification Walkthrough

**Goal**: Verify that the Fyne-based GUI works correctly and matches the design requirements.

## 1. Launch Application
Run the following command in your terminal:
```bash
go run cmd/sudoku/main.go
```
- [ ] **Verify**: Application window opens.
- [ ] **Verify**: Title is "Golden Fox Sudoku".
- [ ] **Verify**: Window is centered and sized appropriately (900x700).

## 2. Visual Inspection
- [ ] **Theme**: Colors should be "Fox Orange" (primary), Charcoal (text), and Soft White (background).
- [ ] **Layout**:
    - **Center**: 9x9 Sudoku Grid.
    - **Bottom**: Number bar (1-9) and Eraser button.
    - **Right**: Control panel (Title, Difficulty, Pencil Mode, New Game, Restart).
- [ ] **Board**:
    - Given numbers (clues) should be visible (darker text).
    - Empty cells should be blank.

## 3. Interaction Testing

### Cell Selection
- [ ] **Action**: Click on any empty cell.
- [ ] **Verify**: Cell background turns Orange/Gold (Selected).
- [ ] **Verify**: Related row, column, and block are highlighted (lighter orange).

### Number Input
- [ ] **Action**: Select an empty cell. Click a number button (e.g., "5") in the bottom bar.
- [ ] **Verify**: The number appears in the cell (Fox Orange color).
- [ ] **Action**: Select a cell with a *Given* number (dark text). Try to change it.
- [ ] **Verify**: The number does NOT change.

### Pencil Mode
- [ ] **Action**: Click "Pencil Mode" button in the right panel.
- [ ] **Verify**: Button highlights or indicates active state.
- [ ] **Action**: Select an empty cell. Click numbers "1", "2", "3".
- [ ] **Verify**: Numbers appear as small candidates (Note: Current implementation might just list them or overwrite; verify behavior).
    - *Note*: In Phase 1, `InputNumber` in pencil mode adds candidates. The `CellWidget` should display them.
    - *Check*: Does `CellWidget` render candidates?
    - *Correction*: Looking at `ui/cell.go`, `Update()` only shows `c.valueText.Text = fmt.Sprintf("%d", val)`. It does **not** yet visualize candidates (pencil notes).
    - **Expected Result**: Pencil notes might not be visible yet in the UI, even if logic handles them. This is a known limitation for Phase 1 GUI.

### Eraser
- [ ] **Action**: Select a user-filled cell. Click the "Eraser" icon button.
- [ ] **Verify**: The cell becomes empty.

### Game Controls
- [ ] **Restart**:
    - **Action**: Fill a few cells. Click "Restart".
    - **Verify**: Board resets to original state (user inputs cleared).
- [ ] **New Game**:
    - **Action**: Click "New Game". Select "Medium".
    - **Verify**: Board refreshes with a new puzzle. Difficulty label updates.

## 4. Win Condition (Optional)
- [ ] **Action**: Solve the puzzle (or use a debug command if available).
- [ ] **Verify**: "Congratulations" dialog appears.

## Known Issues / Future Work (Phase 2)
- **Pencil Marks Visualization**: The current `CellWidget` might not display multiple small numbers for candidates yet.
- **Timer**: Not implemented yet.
- **Undo/Redo**: Not implemented yet.
