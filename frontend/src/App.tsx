import { useState, useEffect, useCallback } from 'react';
import confetti from 'canvas-confetti';
import { CircleHelp } from 'lucide-react';
import bgImage from './assets/images/background.png';

import { Board } from './components/Board';
import { Controls } from './components/Controls';
import { StatsPanel } from './components/StatsPanel';
import { HistoryModal } from './components/HistoryModal';
import { HelpModal } from './components/HelpModal';
import { GameState } from './types';
import LayoutConfig from './config.json';
import { GetGameState, SelectCell, InputNumber, TogglePencilMode, NewGame, ClearCell } from '../wailsjs/go/main/App';

function App() {
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [isHistoryOpen, setIsHistoryOpen] = useState(false);
    const [isHelpOpen, setIsHelpOpen] = useState(false);

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
            // Immediate refresh
            refreshState();
            // Follow-up refresh to ensure state propagation (fixes lag)
            setTimeout(refreshState, 50);
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
            className="flex flex-col h-screen w-screen overflow-hidden relative bg-center bg-no-repeat items-center justify-center"
            style={{
                backgroundImage: `url(${bgImage})`,
                backgroundSize: '100% 100%',
                // @ts-ignore
                '--app-side-length': `${LayoutConfig.APP_SIDE_LENGTH}px`,
                '--scale-base': LayoutConfig.SCALE_BASE,
                '--p-margin-top': LayoutConfig.MARGIN_TOP,
                '--p-margin-left': LayoutConfig.MARGIN_LEFT,
                '--p-header-height': LayoutConfig.HEADER_HEIGHT,
                '--p-board-size': LayoutConfig.BOARD_SIZE,
                '--p-info-bar-width': LayoutConfig.INFO_BAR_WIDTH,
                '--p-info-bar-height': LayoutConfig.INFO_BAR_HEIGHT,
                '--p-control-bar-height': LayoutConfig.CONTROL_BAR_HEIGHT,
                '--p-control-bar-width': LayoutConfig.CONTROL_BAR_WIDTH,
                '--p-gap-board-info': LayoutConfig.GAP_BOARD_INFO,
                '--p-gap-board-control': LayoutConfig.GAP_BOARD_CONTROL,
            }}
        >
            {/* Main Container with Global Margins */}
            <div
                className="flex flex-col w-full h-full"
                style={{
                    paddingTop: 'calc(var(--app-side-length) * (var(--p-margin-top) / var(--scale-base)))',
                    paddingLeft: 'calc(var(--app-side-length) * (var(--p-margin-left) / var(--scale-base)))',
                }}
            >
                {/* 1. Header Zone */}
                <div
                    style={{
                        height: 'calc(var(--app-side-length) * (var(--p-header-height) / var(--scale-base)))',
                        width: '100%'
                    }}
                    className="shrink-0 flex items-center justify-center relative"
                >
                    {/* Help Button */}
                    <button
                        onClick={() => setIsHelpOpen(true)}
                        className="absolute right-8 top-1/2 -translate-y-1/2 p-2 text-[#D68D38] hover:text-[#FFD28F] transition-colors hover:scale-110 active:scale-95"
                        title="How to Play"
                    >
                        <CircleHelp size={32} strokeWidth={2} />
                    </button>

                    {gameState?.isSolved && (
                        <div className="text-6xl font-black text-transparent bg-clip-text bg-gradient-to-b from-[#FFD28F] to-[#D68D38] drop-shadow-[0_2px_10px_rgba(214,141,56,0.5)] animate-in fade-in zoom-in duration-500 tracking-widest">
                            VICTORY!
                        </div>
                    )}
                </div>

                {/* 2. Main Body (Flex Row) */}
                <div className="flex flex-row w-full h-full">
                    {/* Left Column: Board + Gap + Controls */}
                    <div
                        className="flex flex-col"
                        style={{
                            width: 'calc(var(--app-side-length) * (var(--p-board-size) / var(--scale-base)))',
                            marginRight: 'calc(var(--app-side-length) * (var(--p-gap-board-info) / var(--scale-base)))'
                        }}
                    >
                        {/* Board */}
                        <div
                            style={{
                                width: 'calc(var(--app-side-length) * (var(--p-board-size) / var(--scale-base)))',
                                height: 'calc(var(--app-side-length) * (var(--p-board-size) / var(--scale-base)))',
                                marginBottom: 'calc(var(--app-side-length) * (var(--p-gap-board-control) / var(--scale-base)))'
                            }}
                        >
                            <Board
                                board={gameState.board}
                                selectedRow={gameState.selectedRow}
                                selectedCol={gameState.selectedCol}
                                onCellClick={handleCellClick}
                            />
                        </div>

                        {/* Controls */}
                        <div
                            style={{
                                width: 'calc(var(--app-side-length) * (var(--p-control-bar-width) / var(--scale-base)))',
                                height: 'calc(var(--app-side-length) * (var(--p-control-bar-height) / var(--scale-base)))'
                            }}
                        >
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
                    </div>

                    {/* Right Column: Info Bar */}
                    <div
                        className="flex flex-col"
                        style={{
                            width: 'calc(var(--app-side-length) * (var(--p-info-bar-width) / var(--scale-base)))',
                            height: 'calc(var(--app-side-length) * (var(--p-info-bar-height) / var(--scale-base)))'
                        }}
                    >
                        <StatsPanel
                            eraseCount={gameState.eraseCount}
                            undoCount={gameState.undoCount}
                            timeElapsed={gameState.timeElapsed}
                            difficulty={gameState.difficulty}
                            userLevel={gameState.userLevel}
                            gamesPlayed={gameState.gamesPlayed}
                            averageTime={gameState.averageTime}
                            winRate={gameState.winRate}
                            pendingGames={gameState.pendingGames}
                            currentDifficultyCount={gameState.currentDifficultyCount}
                            winsForNextLevel={gameState.winsForNextLevel}
                            remainingCells={gameState.remainingCells}
                        />
                    </div>
                </div>
            </div>

            <HistoryModal
                isOpen={isHistoryOpen}
                onClose={() => setIsHistoryOpen(false)}
                onLoadGame={() => {
                    refreshState();
                }}
            />

            <HelpModal
                isOpen={isHelpOpen}
                onClose={() => setIsHelpOpen(false)}
            />

        </div>
    );
}

export default App;
