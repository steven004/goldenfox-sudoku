import React, { createContext, useContext, ReactNode, useState, useEffect, useMemo } from 'react';
import { GameState } from '../types';
import { useGameLogic } from '../hooks/useGameLogic';
import { useSound } from '../hooks/useSound';
import { useKeyboard } from '../hooks/useKeyboard';
import { useConfetti } from '../hooks/useConfetti';

interface GameContextProps {
    gameState: GameState | null;
    timerSeconds: number;
    pencilMode: boolean;
    selection: { row: number; col: number };
    isMuted: boolean;

    // Actions
    handleCellClick: (row: number, col: number) => void;
    handleNumberClick: (num: number) => void;
    handleActionClick: (action: string) => void;
    toggleMute: () => void;
    handleNewGame: (difficulty: string) => void;
    refreshState: () => void;

    // UI States (Modals)
    isHistoryOpen: boolean;
    setIsHistoryOpen: (isOpen: boolean) => void;
    isHelpOpen: boolean;
    setIsHelpOpen: (isOpen: boolean) => void;
    isNewGameOpen: boolean;
    setIsNewGameOpen: (isOpen: boolean) => void;

    // Derived Data
    completionCounts: Record<number, number>;
}

const GameContext = createContext<GameContextProps | undefined>(undefined);

export const GameProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    // 1. Sound Logic
    const { playSound, isMuted, toggleMute } = useSound();

    // 2. Game Core Logic
    const {
        gameState,
        timerSeconds,
        pencilMode,
        selection,
        refreshState,
        handleCellClick,
        handleNumberClick,
        handleGameAction
    } = useGameLogic(playSound);

    // 3. UI Local State (Modals)
    const [isHistoryOpen, setIsHistoryOpen] = useState(false);
    const [isHelpOpen, setIsHelpOpen] = useState(false);
    const [isNewGameOpen, setIsNewGameOpen] = useState(false);

    // 4. Effects (Confetti & Win Sound)
    useConfetti(gameState?.isSolved || false);

    useEffect(() => {
        if (gameState?.isSolved) {
            playSound('win');
        }
    }, [gameState?.isSolved, playSound]);

    // 5. Action Handlers (Bridge)
    const handleActionClick = (action: string) => {
        if (action === 'history') {
            setIsHistoryOpen(true);
            playSound('click');
        } else if (action === 'new') {
            setIsNewGameOpen(true);
            playSound('click');
        } else {
            handleGameAction(action);
        }
    };

    const handleNewGame = (diff: string) => {
        handleGameAction('new', diff);
    };

    // 6. Keyboard Integration
    useKeyboard({
        gameState,
        selection,
        onCellMove: handleCellClick,
        onNumberInput: handleNumberClick,
        onAction: handleActionClick
    });

    // 7. Derived Data (Optimized)
    const completionCounts = useMemo(() => {
        if (!gameState?.board?.cells) return {};
        const counts: Record<number, number> = {};
        gameState.board.cells.forEach(row => {
            row.forEach(cell => {
                if (cell.value !== 0) {
                    counts[cell.value] = (counts[cell.value] || 0) + 1;
                }
            });
        });
        return counts;
    }, [gameState?.board]);

    const value: GameContextProps = {
        gameState,
        timerSeconds,
        pencilMode,
        selection,
        isMuted,
        handleCellClick,
        handleNumberClick,
        handleActionClick,
        toggleMute,
        handleNewGame,
        refreshState,
        isHistoryOpen,
        setIsHistoryOpen,
        isHelpOpen,
        setIsHelpOpen,
        isNewGameOpen,
        setIsNewGameOpen,
        completionCounts
    };

    return <GameContext.Provider value={value}>{children}</GameContext.Provider>;
};

export const useGame = () => {
    const context = useContext(GameContext);
    if (!context) {
        throw new Error('useGame must be used within a GameProvider');
    }
    return context;
};
