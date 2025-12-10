import React from 'react';
import { Cell } from '../types';

interface CellProps {
    cell: Cell;
    row: number;
    col: number;
    isSelected: boolean;
    isPeer: boolean; // Same row/col/block
    isSameValue: boolean; // Same value as selected
    conflictingCandidates?: number[];
    onClick: (row: number, col: number) => void;
    isPencilMode?: boolean;
}

// --- Style Helper Functions ---

const getBackgroundColor = (
    cell: Cell,
    isSelected: boolean,
    isPencilMode: boolean,
    isSameValue: boolean,
    isPeer: boolean
): string => {
    if (isSelected) {
        if (isPencilMode) return 'bg-gradient-to-br from-sudoku-teal to-sudoku-teal-light text-white';
        // Selection Gradient
        return 'bg-gradient-to-br from-sudoku-primary to-[#EE5A24] text-white';
    }
    if (isSameValue && cell.value !== 0) return 'bg-sudoku-highlight'; // Same-Value Highlight
    if (isPeer) return 'bg-sudoku-peer'; // Peer Highlight
    return cell.given ? 'bg-sudoku-cell-given' : 'bg-sudoku-cell-bg'; // Base Color
};

const getTextColor = (cell: Cell, isSelected: boolean): string => {
    if (cell.isInvalid && !cell.given) {
        return isSelected ? 'text-red-200' : 'text-red-500';
    }
    const fontStyle = cell.given ? 'font-bold' : 'font-light italic';
    const baseColor = isSelected ? 'text-white' : 'text-sudoku-text';
    return `${baseColor} ${fontStyle}`;
};

const getBorderClasses = (row: number, col: number): string => {
    let classes = '';
    // Right Border
    if (col !== 8) {
        classes += (col + 1) % 3 === 0
            ? ' border-r-[2px] border-sudoku-primary-dark'
            : ' border-r border-sudoku-grid';
    }
    // Bottom Border
    if (row !== 8) {
        classes += (row + 1) % 3 === 0
            ? ' border-b-[2px] border-sudoku-primary-dark'
            : ' border-b border-sudoku-grid';
    }
    return classes;
};

// --- Component ---

export const CellComponent: React.FC<CellProps> = ({
    cell,
    row,
    col,
    isSelected,
    isPeer,
    isSameValue,
    conflictingCandidates = [],
    onClick,
    isPencilMode = false
}) => {
    const bgColor = getBackgroundColor(cell, isSelected, isPencilMode, isSameValue, isPeer);
    const textColor = getTextColor(cell, isSelected);
    const borderClasses = getBorderClasses(row, col);

    return (
        <div
            className={`
                w-full h-full flex items-center justify-center cursor-pointer select-none
                text-3xl transition-all duration-100 relative
                ${bgColor} ${textColor} ${borderClasses}
            `}
            onClick={() => onClick(row, col)}
        >
            {cell.value !== 0 ? (
                <span className="text-[2.7rem] leading-none">{cell.value}</span>
            ) : (
                // Render candidates (Pencil Marks)
                <div className="grid grid-cols-3 grid-rows-3 w-full h-full pointer-events-none">
                    {[1, 2, 3, 4, 5, 6, 7, 8, 9].map(num => {
                        const isConflict = conflictingCandidates.includes(num);
                        return (
                            <div key={num}
                                className={`flex items-center justify-center text-[15px] leading-none font-medium ${isConflict ? 'text-red-500 font-bold' : 'text-sudoku-text-secondary'}`}
                            >
                                {cell.candidates && cell.candidates[num] ? num : ''}
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
};
