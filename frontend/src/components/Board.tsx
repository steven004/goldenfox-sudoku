import React from 'react';
import { SudokuBoard } from '../types';
import { CellComponent } from './Cell';

interface BoardProps {
    board: SudokuBoard;
    selectedRow: number;
    selectedCol: number;
    onCellClick: (row: number, col: number) => void;
}

export const Board: React.FC<BoardProps> = ({ board, selectedRow, selectedCol, onCellClick }) => {
    if (!board || !board.cells) {
        return <div className="text-white">Loading board...</div>;
    }

    // Helper to check if a cell is a peer of the selected cell
    const isPeer = (r: number, c: number) => {
        if (selectedRow === -1 || selectedCol === -1) return false;
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
        <div className="aspect-square h-full w-auto bg-[#FDF6E3] rounded-2xl overflow-hidden border-4 border-[#D68D38] shadow-[0_10px_30px_rgba(0,0,0,0.5)]">
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
                            onClick={onCellClick}
                        />
                    ))
                ))}
            </div>
        </div>
    );
};
