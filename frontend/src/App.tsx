import { useState, useEffect } from 'react';
import { useGameLogic } from './hooks/useGameLogic';
import { useConfetti } from './hooks/useConfetti';
import { useKeyboard } from './hooks/useKeyboard';
import { useSound } from './hooks/useSound';
import { GameLayout } from './components/GameLayout';

function App() {
    const { playSound, isMuted, toggleMute } = useSound();

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

    const [isHistoryOpen, setIsHistoryOpen] = useState(false);
    const [isHelpOpen, setIsHelpOpen] = useState(false);
    const [isNewGameOpen, setIsNewGameOpen] = useState(false);

    useConfetti(gameState?.isSolved || false);

    // Play Win Sound
    useEffect(() => {
        if (gameState?.isSolved) {
            playSound('win');
        }
    }, [gameState?.isSolved, playSound]);

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

    const handleNewGame = (difficulty: string) => {
        // ... (existing logic commented out in previous turn, actually not needed if using handleActionClick 'new') ...
    };

    // Keyboard Controls
    useKeyboard({
        gameState,
        selection,
        onCellMove: handleCellClick,
        onNumberInput: handleNumberClick,
        onAction: handleActionClick
    });

    return (
        <GameLayout
            gameState={gameState}
            timerSeconds={timerSeconds}
            pencilMode={pencilMode}
            selection={selection}
            onCellClick={handleCellClick}
            onNumberClick={handleNumberClick}
            onActionClick={handleActionClick}
            isHistoryOpen={isHistoryOpen}
            setIsHistoryOpen={setIsHistoryOpen}
            isHelpOpen={isHelpOpen}
            setIsHelpOpen={setIsHelpOpen}
            isNewGameOpen={isNewGameOpen}
            setIsNewGameOpen={setIsNewGameOpen}
            onLoadGame={refreshState}
            onNewGame={(diff) => handleGameAction('new', diff)}
            isMuted={isMuted}
            onToggleMute={toggleMute}
        />
    );
}

export default App;

