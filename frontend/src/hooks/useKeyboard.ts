import { useEffect } from 'react';
import { GameState } from '../types';

interface KeyboardProps {
    gameState: GameState | null;
    onCellMove: (row: number, col: number) => void;
    onNumberInput: (num: number, forcePencil?: boolean) => void;
    onAction: (action: string) => void;
}

export const useKeyboard = ({ gameState, onCellMove, onNumberInput, onAction }: KeyboardProps) => {
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (!gameState || gameState.isSolved) return;

            // Prevent default scrolling for arrows/space
            if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', ' '].includes(e.key)) {
                e.preventDefault();
            }

            // --- 1. Selection Navigation (Arrows) ---
            if (e.key.startsWith('Arrow')) {
                let { selectedRow, selectedCol } = gameState;

                // If nothing selected, start at center (4,4)
                if (selectedRow === -1 || selectedCol === -1) {
                    onCellMove(4, 4);
                    return;
                }

                switch (e.key) {
                    case 'ArrowUp':
                        if (selectedRow > 0) onCellMove(selectedRow - 1, selectedCol);
                        break;
                    case 'ArrowDown':
                        if (selectedRow < 8) onCellMove(selectedRow + 1, selectedCol);
                        break;
                    case 'ArrowLeft':
                        if (selectedCol > 0) onCellMove(selectedRow, selectedCol - 1);
                        break;
                    case 'ArrowRight':
                        if (selectedCol < 8) onCellMove(selectedRow, selectedCol + 1);
                        break;
                }
                return;
            }

            // --- 2. Number Input (1-9) ---
            const num = parseInt(e.key);
            if (!isNaN(num) && num >= 1 && num <= 9) {
                // Shift + Number -> Force Pencil Note
                onNumberInput(num, e.shiftKey);
                return;
            }

            // --- 3. Tools ---
            switch (e.key.toLowerCase()) {
                case 'backspace':
                case 'delete':
                    onAction('eraser');
                    break;
                case 'p': // Pencil Toggle
                case 'n': // Notes Toggle
                    onAction('pencil');
                    break;
                case 'u': // Undo
                    onAction('undo');
                    break;
                case 'z': // Undo (Ctrl+Z / Cmd+Z)
                    if (e.metaKey || e.ctrlKey) {
                        e.preventDefault(); // Stop browser undo
                        onAction('undo');
                    }
                    break;
                case 'n':
                    // Cmd+N -> New Game
                    if (e.metaKey || e.ctrlKey) {
                        e.preventDefault();
                        onAction('new');
                    }
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [gameState, onCellMove, onNumberInput, onAction]);
};
