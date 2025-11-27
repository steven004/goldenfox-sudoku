export interface Cell {
    value: number;
    given: boolean;
    candidates: Record<number, boolean>;
    isInvalid?: boolean;
}

export interface SudokuBoard {
    cells: Cell[][];
}

export interface GameState {
    board: SudokuBoard;
    selectedRow: number;
    selectedCol: number;
    isSelected: boolean;
    pencilMode: boolean;
    mistakes: number;
    timeElapsed: string;
    difficulty: string;
    isSolved: boolean;
}
