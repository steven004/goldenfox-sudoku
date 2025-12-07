import { useState, useEffect, useCallback } from 'react';
import { GameState } from '../types';
import { GetGameState, SelectCell, InputNumber, ToggleCandidate, NewGame, ClearCell, RestartGame, Undo } from '../../wailsjs/go/main/App';

export const useGameLogic = () => {
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [timerSeconds, setTimerSeconds] = useState(0);
    const [pencilMode, setPencilMode] = useState(false);

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

    const handleCellClick = async (row: number, col: number) => {
        try {
            await SelectCell(row, col);
            refreshState();
        } catch (err) {
            console.error(err);
        }
    };

    const handleNumberClick = async (num: number, forcePencil: boolean = false) => {
        try {
            const isNote = forcePencil || pencilMode;
            if (isNote) {
                await ToggleCandidate(num);
            } else {
                await InputNumber(num);
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
                    await ClearCell();
                    refreshState();
                    break;
                case 'new':
                    await NewGame(difficulty || "Easy");
                    refreshState();
                    break;
                case 'restart':
                    await RestartGame();
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
        refreshState,
        handleCellClick,
        handleNumberClick,
        handleGameAction,
    };
};
