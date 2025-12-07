import { useState, useEffect, useCallback } from 'react';
import { GameState } from '../types';
import { GetGameState, InputNumber, ToggleCandidate, NewGame, ClearCell, RestartGame, Undo } from '../../wailsjs/go/main/App';

export const useGameLogic = () => {
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [timerSeconds, setTimerSeconds] = useState(0);
    const [pencilMode, setPencilMode] = useState(false);
    const [selection, setSelection] = useState({ row: 4, col: 4 });

    const refreshState = useCallback(() => {
        GetGameState().then((state: any) => {
            setGameState(state as GameState);
        }).catch(err => console.error("Failed to get game state:", err));
    }, []);

    // Sync local timer with backend state
    useEffect(() => {
        if (gameState) {
            setTimerSeconds(gameState.elapsedSeconds);
        }
    }, [gameState?.elapsedSeconds]);

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
        // No backend call needed
    };

    const handleNumberClick = async (num: number, forcePencil: boolean = false) => {
        try {
            const { row, col } = selection;
            if (row === -1 || col === -1) return;

            const isNote = forcePencil || pencilMode;
            if (isNote) {
                await ToggleCandidate(row, col, num);
            } else {
                await InputNumber(row, col, num);
            }
            refreshState();
            setTimeout(refreshState, 50);
        } catch (err) {
            console.error(err);
        }
    };

    const handleGameAction = async (action: string, difficulty?: string) => {
        if (!gameState) return;
        try {
            switch (action) {
                case 'pencil':
                    setPencilMode(prev => !prev);
                    break;
                case 'eraser':
                    if (gameState.isSolved) return;
                    if (selection.row !== -1 && selection.col !== -1) {
                        await ClearCell(selection.row, selection.col);
                        refreshState();
                    }
                    break;
                case 'new':
                    await NewGame(difficulty || "Easy");
                    setSelection({ row: 4, col: 4 });
                    refreshState();
                    break;
                case 'restart':
                    await RestartGame();
                    setSelection({ row: 4, col: 4 });
                    refreshState();
                    break;
                case 'undo':
                    if (gameState.isSolved) return;
                    await Undo();
                    refreshState();
                    break;
            }
        } catch (err) {
            console.error(err);
        }
    };

    return {
        gameState,
        timerSeconds,
        pencilMode,
        selection,
        refreshState,
        handleCellClick,
        handleNumberClick,
        handleGameAction,
    };
};
