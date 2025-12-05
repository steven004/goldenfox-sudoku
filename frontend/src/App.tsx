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

    useConfetti(gameState?.isSolved || false);

    const handleActionClick = (action: string) => {
        if (action === 'history') {
            setIsHistoryOpen(true);
        } else {
            handleGameAction(action);
        }
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
            onLoadGame={refreshState}
        />
    );
}

export default App;

