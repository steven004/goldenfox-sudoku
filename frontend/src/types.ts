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
    eraseCount: number;
    undoCount: number;
    timeElapsed: string;
    difficulty: string;
    isSolved: boolean;
    userLevel: number;
    gamesPlayed: number;
    winRate: number;
    pendingGames: number;
    averageTime: string;
    currentDifficultyCount: number;
    winsForNextLevel: number;
    remainingCells: number;
}
