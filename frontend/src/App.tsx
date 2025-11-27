import { useState, useEffect, useCallback } from 'react';
import confetti from 'canvas-confetti';
import logo from './assets/images/fox-sudoku.png';
import bgImage from './assets/images/background.png';

import { Board } from './components/Board';
import { Controls } from './components/Controls';
import { StatsPanel } from './components/StatsPanel';
import { HistoryModal } from './components/HistoryModal';
import { GameState } from './types';
import { GetGameState, SelectCell, InputNumber, TogglePencilMode, NewGame, ClearCell, SaveGame } from '../wailsjs/go/main/App';

function App() {
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [isHistoryOpen, setIsHistoryOpen] = useState(false);

    const refreshState = useCallback(() => {
        GetGameState().then((state: any) => {
            setGameState(state as GameState);
        }).catch(err => console.error("Failed to get game state:", err));
    }, []);

    useEffect(() => {
        // Initial load
        refreshState();

        // Set up a poller for timer (every second)
        const interval = setInterval(() => {
            // In a real app, we might just increment local timer and sync occasionally
            // For now, let's sync from backend to get the "GetElapsedTime" string
            refreshState();
        }, 1000);

        return () => clearInterval(interval);
    }, [refreshState]);

    // Trigger fireworks when solved
    useEffect(() => {
        if (gameState?.isSolved) {
            const duration = 2000;
            const end = Date.now() + duration;

            const frame = () => {
                confetti({
                    particleCount: 2,
                    angle: 60,
                    spread: 55,
                    origin: { x: 0 },
                    colors: ['#D68D38', '#FFD28F', '#FFFFFF']
                });
                confetti({
                    particleCount: 2,
                    angle: 120,
                    spread: 55,
                    origin: { x: 1 },
                    colors: ['#D68D38', '#FFD28F', '#FFFFFF']
                });

                if (Date.now() < end) {
                    requestAnimationFrame(frame);
                }
            };
            frame();
        }
    }, [gameState?.isSolved]);

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
        } catch (err) {
            console.error(err);
        }
    };

    const handleActionClick = async (action: string) => {
        try {
            switch (action) {
                case 'pencil':
                    await TogglePencilMode();
                    break;
                case 'eraser': // Match the ID from Controls.tsx
                    if (gameState?.isSolved) return;
                    await ClearCell();
                    break;
                case 'new':
                    await NewGame("Easy"); // Default to Easy for now
                    break;
                case 'restart':
                    // TODO: Implement Restart in backend binding if not already
                    await NewGame("Easy"); // Temporary
                    break;
                case 'history':
                    setIsHistoryOpen(true);
                    break;
                case 'save':
                    await SaveGame();
                    // Optional: Show toast
                    break;
                // Other actions...
            }
            refreshState();
        } catch (err) {
            console.error(err);
        }
    };

    if (!gameState) {
        return <div className="flex items-center justify-center h-screen text-white">Loading Golden Fox Sudoku...</div>;
    }

    return (
        <div
            className="flex flex-col h-screen w-screen overflow-hidden relative bg-cover bg-center bg-no-repeat"
            style={{ backgroundImage: `url(${bgImage})` }}
        >
            {/* 1. Top Zone: The Header */}
            <header className="h-20 flex items-center justify-center gap-4 shrink-0 pt-4">
                <img src={logo} alt="Golden Fox Logo" className="h-12 w-auto object-contain drop-shadow-lg" />
                <h1 className={`text-3xl font-bold tracking-wider drop-shadow-md font-display transition-all duration-500 ${gameState.isSolved ? 'text-[#FFD28F] scale-110' : 'text-white'}`}>
                    {gameState.isSolved ? '✨ VICTORY! ✨' : 'Golden Fox Sudoku'}
                </h1>
            </header>

            {/* 2. Middle Zone: The Main Stage */}
            <main className="flex-1 flex flex-row w-full max-w-[1600px] mx-auto p-4 gap-8 min-h-0">

                {/* Left Column (Game Area) - Flex 3 */}
                <div className="flex-[3] flex flex-col items-center justify-center gap-6 h-full">
                    {/* Board Area - Square, Centered */}
                    <div className="aspect-square w-full max-h-[65vh] flex items-center justify-center">
                        <Board
                            board={gameState.board}
                            selectedRow={gameState.selectedRow}
                            selectedCol={gameState.selectedCol}
                            onCellClick={handleCellClick}
                        />
                    </div>
                </div>

                {/* Right Column (Info Panel) - Flex 1 */}
                <div className="flex-[1] flex flex-col justify-center h-full max-w-sm">
                    <StatsPanel
                        mistakes={gameState.mistakes}
                        timeElapsed={gameState.timeElapsed}
                        difficulty={gameState.difficulty}
                    />
                </div>
            </main>

            {/* 3. Bottom Zone: Controls */}
            <div className="w-full max-w-4xl mx-auto pb-8 px-4">
                <Controls
                    onNumberClick={handleNumberClick}
                    onActionClick={handleActionClick}
                    pencilMode={gameState.pencilMode}
                    selectedNumber={
                        gameState.selectedRow !== -1 && gameState.selectedCol !== -1
                            ? gameState.board.cells[gameState.selectedRow][gameState.selectedCol].value
                            : undefined
                    }
                    completionCounts={(() => {
                        const counts: Record<number, number> = {};
                        gameState.board.cells.forEach(row => {
                            row.forEach(cell => {
                                if (cell.value !== 0) {
                                    counts[cell.value] = (counts[cell.value] || 0) + 1;
                                }
                            });
                        });
                        return counts;
                    })()}
                />
            </div>

            <HistoryModal
                isOpen={isHistoryOpen}
                onClose={() => setIsHistoryOpen(false)}
                onLoadGame={() => {
                    refreshState();
                    // Maybe show a toast?
                }}
            />
        </div>
    );
}

export default App;
