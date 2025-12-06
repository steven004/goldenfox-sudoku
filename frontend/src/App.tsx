import { useState } from 'react';
import { useGameLogic } from './hooks/useGameLogic';
import { useConfetti } from './hooks/useConfetti';
import { GameLayout } from './components/GameLayout';

function App() {
    const {
        gameState,
        timerSeconds,
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
        // We route this through handleGameAction ('new') but now we need to pass diff.
        // wait, handleGameAction in useGameLogic uses StartNewGame() which matches standard behavior.
        // We need to call the SPECIFIC NewGame(diff) backend logic.
        // Let's modify useGameLogic or just call StartNewGame(diff) if we update it?
        // Actually, handleGameAction('new') was just calling the internal logic.
        // Let's call the backend directly or via a new hook method.
        // StartNewGame is exposed by Wails as main.NewGame(diff).
        // Let's check useGameLogic.
        // For now, I will assume useGameLogic needs an update OR I can call the window.go.main.App.NewGame directly here.
        // Better: Update useGameLogic to support difficulty.
    };

    return (
        <GameLayout
            gameState={gameState}
            timerSeconds={timerSeconds}
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

