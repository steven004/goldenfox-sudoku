import { SudokuBoard } from '../types';

export const isValidMove = (board: SudokuBoard, row: number, col: number, value: number): boolean => {
    if (value === 0) return true;

    // Check Row
    for (let c = 0; c < 9; c++) {
        if (c !== col && board.cells[row][c].value === value) return false;
    }

    // Check Col
    for (let r = 0; r < 9; r++) {
        if (r !== row && board.cells[r][col].value === value) return false;
    }

    // Check Block
    const startR = Math.floor(row / 3) * 3;
    const startC = Math.floor(col / 3) * 3;
    for (let r = startR; r < startR + 3; r++) {
        for (let c = startC; c < startC + 3; c++) {
            if ((r !== row || c !== col) && board.cells[r][c].value === value) return false;
        }
    }

    return true;
};

export const getConflictingCandidates = (board: SudokuBoard, row: number, col: number): number[] => {
    const cell = board.cells[row][col];
    if (!cell.candidates) return [];

    const conflicting: number[] = [];
    Object.keys(cell.candidates).forEach(key => {
        const num = parseInt(key);
        if (cell.candidates[num]) {
            if (!isValidMove(board, row, col, num)) {
                conflicting.push(num);
            }
        }
    });

    return conflicting;
};
