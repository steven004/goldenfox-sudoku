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

    // Computed Game State (Base + Transient Overlay)
    const displayState = useMemo(() => {
        if (!gameState) return null;
        if (!transientError) return gameState;

        // Clone deeply enough to modify the specific cell
        // We do a shallow clone of board/cells structure for performance, assuming immutable updates elsewhere
        // But here we need to be careful not to mutate the original state if it's reused.
        // For safety, let's clone the row at least.
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

    const handleCellClick = (row: number, col: number) => {
        setSelection({ row, col });
        onSound?.('click');
    };

    const handleNumberClick = async (num: number, forcePencil: boolean = false) => {
        try {
            const { row, col } = selection;
            if (row === -1 || col === -1) return;

            const isNote = forcePencil || pencilMode;
            if (isNote) {
                await ToggleCandidate(row, col, num);
                onSound?.('pencil');
            } else {
                await InputNumber(row, col, num);
                onSound?.('pop');
            }
            refreshState();
            setTimeout(refreshState, 50);
        } catch (err: any) {
            console.error(err);
            // Check if error is due to conflict (backend returns "conflict" string or similar in error message usually)
            // But we can assume any error during input/toggle on a valid cell is likely a rule violation or immutable cell.
            // For better UX, let's show the transient error if it matches expected conflict patterns or just generic "invalid".

            // To show the red number, we need the number we TRIED to input.
            // "num" is available here.

            // We trigger transient error for visual feedback
            const { row, col } = selection;
            if (!pencilMode && !forcePencil) {
                // Only show big red number for Value input. 
                // For Notes, maybe we just play error sound? User asked "logic... is the same".
                // But notes are small. Showing a big red number for a failed note might be confusing?
                // Wait, if addNote fails, it means the NOTE conflicts with a NUMBER. 
                // Displaying the big red number implies "This NUMBER is invalid here", which is true!
                // So yes, showing the transient error is correct for both cases if the value itself is the problem.
                setTransientError({ row, col, value: num });
            } else {
                // If it was a note attempt that failed, should we show the big number?
                // If I try to pencil "5" and 5 is already in the row...
                // Seeing a big red "5" fade out conveys "5 is not allowed here". That works.
                setTransientError({ row, col, value: num });
            }

            onSound?.('error');
        }
    };

    const handleGameAction = async (action: string, difficulty?: string) => {
        if (!gameState) return;
        try {
            switch (action) {
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
                    refreshState();
                    onSound?.('click');
                    break;
                case 'restart':
                    await RestartGame();
                    setSelection({ row: 4, col: 4 });
                    setTimerSeconds(0);
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
        selection,
        refreshState,
        handleCellClick,
        handleNumberClick,
        handleGameAction,
    };
};
