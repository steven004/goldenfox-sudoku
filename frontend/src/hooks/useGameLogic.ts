import { useState, useEffect, useCallback } from 'react';
import { GameState } from '../types';
import { GetGameState, SelectCell, InputNumber, TogglePencilMode, NewGame, ClearCell, RestartGame, Undo } from '../../wailsjs/go/main/App';

export const useGameLogic = () => {
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [timerSeconds, setTimerSeconds] = useState(0);

    const refreshState = useCallback(() => {
        GetGameState().then((state: any) => {
            setGameState(state as GameState);
        }).catch(err => console.error("Failed to get game state:", err));
    }, []);

    // Sync local timer with backend state when it updates
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

    const handleNumberClick = async (num: number) => {
        try {
            await InputNumber(num);
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
                    await TogglePencilMode();
                    break;
                case 'eraser':
                    if (gameState.isSolved) return;
                    await ClearCell();
                    break;
                case 'new':
                    await NewGame(difficulty || "Easy");
                    break;
                case 'restart':
                    await RestartGame();
                    break;
                case 'undo':
                    if (gameState.isSolved) return;
                    await Undo();
                    break;
            }
            refreshState();
        } catch (err) {
            console.error(err);
        }
    };

    return {
        gameState,
        timerSeconds,
        refreshState,
        handleCellClick,
        handleNumberClick,
        handleGameAction
    };
};
