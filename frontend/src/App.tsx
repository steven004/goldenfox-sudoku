import { GameProvider } from './context/GameContext';
import { GameLayout } from './components/GameLayout';

function App() {
    return (
        <GameProvider>
            <GameLayout />
        </GameProvider>
    );
}

export default App;

