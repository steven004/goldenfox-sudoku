import React from 'react';
import { SudokuBoard } from '../types';
import { CellComponent } from './Cell';
import { getConflictingCandidates } from '../utils/validation';

interface BoardProps {
    board: SudokuBoard;
    selectedRow: number;
    selectedCol: number;
    onCellClick: (row: number, col: number) => void;
    isPencilMode: boolean;
}

export const Board: React.FC<BoardProps> = ({ board, selectedRow, selectedCol, onCellClick, isPencilMode }) => {
    if (!board || !board.cells) {
        return <div className="text-white">Loading board...</div>;
    }

    // Helper to check if a cell is a peer of the selected cell
    const isPeer = (r: number, c: number) => {
        if (selectedRow === -1 || selectedCol === -1) return false;

        // Only highlight peers if the selected cell is BLANK
        const selectedCell = board.cells[selectedRow][selectedCol];
        if (selectedCell.value !== 0) return false;

        if (r === selectedRow) return true;
        if (c === selectedCol) return true;
        // Check block
        const blockRow = Math.floor(r / 3);
        const blockCol = Math.floor(c / 3);
        const selBlockRow = Math.floor(selectedRow / 3);
        const selBlockCol = Math.floor(selectedCol / 3);
        return blockRow === selBlockRow && blockCol === selBlockCol;
    };

    // Helper to check if cell has same value as selected cell
    const isSameValue = (r: number, c: number) => {
        if (selectedRow === -1 || selectedCol === -1) return false;
        const selectedVal = board.cells[selectedRow][selectedCol].value;
        if (selectedVal === 0) return false;
        return board.cells[r][c].value === selectedVal;
    };

    return (
        <div className="aspect-square h-full w-auto bg-sudoku-board rounded-2xl overflow-hidden border-8 border-sudoku-primary-dark shadow-[inset_3px_3px_6px_rgba(255,255,255,0.9),inset_-3px_-3px_6px_rgba(0,0,0,0.4),0_15px_35px_rgba(0,0,0,0.6),0_0_0_4px_rgba(214,141,56,0.4)]">
            <div className="grid grid-cols-9 grid-rows-9 w-full h-full">
                {board.cells.map((row, r) => (
                    row.map((cell, c) => (
                        <CellComponent
                            key={`${r}-${c}`}
                            cell={cell}
                            row={r}
                            col={c}
                            isSelected={r === selectedRow && c === selectedCol}
                            isPeer={isPeer(r, c)}
                            isSameValue={isSameValue(r, c)}
                            conflictingCandidates={cell.value === 0 ? getConflictingCandidates(board, r, c) : undefined}
                            onClick={onCellClick}
                            isPencilMode={isPencilMode}
                        />
                    ))
                ))}
            </div>
        </div>
    );
};
