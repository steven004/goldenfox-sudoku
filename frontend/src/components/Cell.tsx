import React from 'react';
import { Cell } from '../types';

interface CellProps {
    cell: Cell;
    row: number;
    col: number;
    isSelected: boolean;
    isPeer: boolean; // Same row/col/block
    isSameValue: boolean; // Same value as selected
    onClick: (row: number, col: number) => void;
}

export const CellComponent: React.FC<CellProps> = ({
    cell,
    row,
    col,
    isSelected,
    isPeer,
    isSameValue,
    onClick
}) => {
    // Base background color (Cream/Off-white for the board)
    let bgColor = 'bg-[#FDF6E3]';

    // Highlighting logic
    if (isSelected) {
        bgColor = 'bg-gradient-to-br from-[#FF9F43] to-[#EE5A24] text-white'; // Selected Highlight
    } else if (isSameValue && cell.value !== 0) {
        bgColor = 'bg-[#FFEAA7]'; // Same-Number Highlight
    } else if (isPeer) {
        // Peer highlight (Row/Col/Block)
        bgColor = 'bg-[#FFE0B2]';
    }

    // Text color (Dark Charcoal for numbers, White for selected)
    // Differentiate font weight/style: Bold for Given, Light + Italic for User Input
    // Conflict Check: Red if invalid
    const fontStyle = cell.given ? 'font-bold' : 'font-light italic';
    let colorClass = isSelected ? 'text-white' : 'text-[#2D3436]';

    if (cell.isInvalid && !cell.given) {
        colorClass = isSelected ? 'text-red-200' : 'text-red-500';
    }

    const textColor = `${colorClass} ${fontStyle}`;

    // Grid Lines Logic (Strict Specs)
    let borderClasses = '';

    // Right Border
    if (col !== 8) {
        if ((col + 1) % 3 === 0) {
            borderClasses += ' border-r-[2px] border-[#D68D38]'; // Thick Divider
        } else {
            borderClasses += ' border-r border-[#B2BEC3]'; // Thin Grid Line
        }
    }

    // Bottom Border
    if (row !== 8) {
        if ((row + 1) % 3 === 0) {
            borderClasses += ' border-b-[2px] border-[#D68D38]'; // Thick Divider
        } else {
            borderClasses += ' border-b border-[#B2BEC3]'; // Thin Grid Line
        }
    }

    return (
        <div
            className={`
                w-full h-full flex items-center justify-center cursor-pointer select-none
                text-3xl transition-all duration-100 relative
                ${bgColor} ${textColor}
                ${borderClasses}
            `}
            onClick={() => onClick(row, col)}
        >
            {cell.value !== 0 ? (
                <span className="text-[2.7rem] leading-none">{cell.value}</span>
            ) : (
                // Render candidates (Pencil Marks)
                <div className="grid grid-cols-3 grid-rows-3 w-full h-full pointer-events-none">
                    {[1, 2, 3, 4, 5, 6, 7, 8, 9].map(num => (
                        <div key={num} className="flex items-center justify-center text-[10px] text-[#636e72] leading-none font-medium">
                            {cell.candidates && cell.candidates[num] ? num : ''}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};
