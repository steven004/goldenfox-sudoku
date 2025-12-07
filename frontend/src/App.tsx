import { useState } from 'react';
import { useGameLogic } from './hooks/useGameLogic';
import { useConfetti } from './hooks/useConfetti';
import { useKeyboard } from './hooks/useKeyboard';
import { GameLayout } from './components/GameLayout';

function App() {
    const {
        gameState,
        timerSeconds,
        pencilMode,
        refreshState,
        handleCellClick,
        handleNumberClick,
        handleGameAction
    } = useGameLogic();

    const [isHistoryOpen, setIsHistoryOpen] = useState(false);
    const [isHelpOpen, setIsHelpOpen] = useState(false);
    const [isNewGameOpen, setIsNewGameOpen] = useState(false);

    useConfetti(gameState?.isSolved || false);

    const handleActionClick = (action: string) => {
        if (action === 'history') {
            setIsHistoryOpen(true);
        } else if (action === 'new') {
            setIsNewGameOpen(true);
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
        onCellMove: handleCellClick,
        onNumberInput: handleNumberClick,
        onAction: handleActionClick
    });

    return (
        <GameLayout
            gameState={gameState}
            timerSeconds={timerSeconds}
            pencilMode={pencilMode}
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
        />
    );
}

export default App;

