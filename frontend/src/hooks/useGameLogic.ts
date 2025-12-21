import { useState, useEffect, useCallback, useMemo } from 'react';
import { GameState } from '../types';
import { GetGameState, InputNumber, ToggleCandidate, NewGame, ClearCell, RestartGame, Undo } from '../../wailsjs/go/main/App';
import { isValidMove } from '../utils/validation';

export const useGameLogic = (onSound?: (type: 'click' | 'pop' | 'error' | 'erase' | 'win' | 'pencil') => void) => {
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [timerSeconds, setTimerSeconds] = useState(0);
    const [pencilMode, setPencilMode] = useState(false);
    const [selection, setSelection] = useState({ row: 4, col: 4 });

    // Transient Error State (row, col, value)
    const [transientError, setTransientError] = useState<{ row: number, col: number, value: number } | null>(null);

    // Clear transient error after 500ms
    useEffect(() => {
        if (transientError) {
            const timer = setTimeout(() => {
                setTransientError(null);
            }, 500);
            return () => clearTimeout(timer);
        }
    }, [transientError]);

    const [fastMode, setFastMode] = useState(false);
    const [highlightedNumber, setHighlightedNumber] = useState<number | null>(null);

    // Computed Game State (Base + Transient Overlay)
    const displayState = useMemo(() => {
        if (!gameState) return null;
        if (!transientError) return gameState;
        // ... (existing clone logic kept implicit for brevity in tool call, but I must preserve it or assume context handles it. 
        // Wait, replace tool replaces the block. I need to be careful not to delete lines outside my target if I don't include them.
        // Actually, I can insert fastMode state declaration separately.
        // Let's split this into safer chunks.

        const newCells = gameState.board.cells.map((row, r) =>
            r === transientError.row
                ? row.map((cell, c) =>
                    c === transientError.col
                        ? { ...cell, value: transientError.value, isInvalid: true }
                        : cell
                )
                : row
        );

        return {
            ...gameState,
            board: {
                ...gameState.board,
                cells: newCells
            }
        };
    }, [gameState, transientError]);


    const refreshState = useCallback(() => {
        GetGameState().then((state: any) => {
            setGameState(state as GameState);
        }).catch(err => console.error("Failed to get game state:", err));
    }, []);

    // Sync local timer with backend state - REMOVED (Frontend independent timer)


    // Local Timer Ticking
    useEffect(() => {
        if (!gameState || gameState.isSolved) return;

        const interval = setInterval(() => {
            setTimerSeconds(prev => prev + 1);
        }, 1000);

        return () => clearInterval(interval);
    }, [gameState?.isSolved]);

    // Initial load
    useEffect(() => {
        refreshState();
    }, [refreshState]);

    const handleCellClick = async (row: number, col: number) => {
        setSelection({ row, col });
        onSound?.('click');

        if (fastMode && highlightedNumber !== null) {
            // FAST MODE: Paint the cell with the active tool (number)
            try {
                // If it's the SAME cell we just clicked to select (unlikely if painting, but possible), does input work? Yes.
                const isNote = pencilMode;
                if (isNote) {
                    await ToggleCandidate(row, col, highlightedNumber);
                    onSound?.('pencil');
                } else {
                    await InputNumber(row, col, highlightedNumber);
                    onSound?.('pop');
                }
                refreshState();
                // Delay refresh to allow backend to process if needed, similar to handleNumberClick
                setTimeout(refreshState, 50);
            } catch (err) {
                // Reuse error handling logic? Or simple fallback.
                // Ideally we want the transient error here too.
                setTransientError({ row, col, value: highlightedNumber });
                onSound?.('error');
            }
            return; // Skip standard highlight logic in Fast Mode painting
        }

        // Standard Highlight Logic (Consistent Cell-First)
        if (gameState && gameState.board.cells[row][col]) {
            const cellVal = gameState.board.cells[row][col].value;
            if (cellVal !== 0) {
                // Always update highlight if clicking a filled cell
                setHighlightedNumber(cellVal);
            } else {
                // If clicking a blank cell, ALWAYS clear highlight (Standard Peer View)
                setHighlightedNumber(null);
            }
        }
    };

    const handleNumberClick = async (num: number, forcePencil: boolean = false) => {
        if (fastMode) {
            // FAST MODE: Select Tool
            setHighlightedNumber(num);
            return;
        }

        try {
            const { row, col } = selection;
            if (row === -1 || col === -1) return;

            const isNote = forcePencil || pencilMode;
            if (isNote) {
                await ToggleCandidate(row, col, num);
                onSound?.('pencil');
                // Don't change highlight for notes, keep context
            } else {
                await InputNumber(row, col, num);
                onSound?.('pop');
                // When inputting a number, implicitly switch highlight to that number
                setHighlightedNumber(num);
            }
            refreshState();
            setTimeout(refreshState, 50);
        } catch (err: any) {
            console.error(err);
            const { row, col } = selection;
            // Show transient error regardless of mode/type if input failed
            setTransientError({ row, col, value: num });
            onSound?.('error');
        }
    };

    const handleGameAction = async (action: string, difficulty?: string) => {
        if (!gameState) return;
        try {
            switch (action) {
                case 'fastMode':
                    setFastMode(prev => !prev);
                    // If turning OFF, maybe clear highlight?
                    // If turning ON, maybe keep current selection?
                    // Let's keep it simple.
                    if (!fastMode) { // Turning ON
                        setHighlightedNumber(null); // Reset tool
                    }
                    onSound?.('click');
                    break;
                case 'pencil':
                    setPencilMode(prev => !prev);
                    onSound?.('click');
                    break;
                case 'eraser':
                    if (gameState.isSolved) return;
                    if (selection.row !== -1 && selection.col !== -1) {
                        await ClearCell(selection.row, selection.col);
                        refreshState();
                        onSound?.('erase');
                    } else {
                        onSound?.('error');
                    }
                    break;
                case 'new':
                    await NewGame(difficulty || "Easy");
                    setSelection({ row: 4, col: 4 });
                    setTimerSeconds(0);
                    setHighlightedNumber(null); // Reset highlight
                    refreshState();
                    onSound?.('click');
                    break;
                case 'restart':
                    await RestartGame();
                    setSelection({ row: 4, col: 4 });
                    setTimerSeconds(0);
                    setHighlightedNumber(null); // Reset highlight
                    refreshState();
                    onSound?.('erase');
                    break;
                case 'undo':
                    if (gameState.isSolved) return;
                    await Undo();
                    refreshState();
                    onSound?.('click');
                    break;
            }
        } catch (err) {
            console.error(err);
            onSound?.('error');
        }
    };

    return {
        gameState: displayState, // Return the computed state
        timerSeconds,
        pencilMode,
        fastMode, // Return fastMode state
        selection,
        highlightedNumber, // Return the sticky highlight state
        refreshState,
        handleCellClick,
        handleNumberClick,
        handleGameAction,
    };
};
