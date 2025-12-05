import React from 'react';
import { CircleHelp } from 'lucide-react';
import bgImage from '../assets/images/background.png';
import { Board } from './Board';
import { Controls } from './Controls';
import { StatsPanel } from './StatsPanel';
import { HistoryModal } from './HistoryModal';
import { HelpModal } from './HelpModal';
import { GameState } from '../types';
import LayoutConfig from '../config.json';

interface GameLayoutProps {
    gameState: GameState | null;
    timerSeconds: number;
    onCellClick: (row: number, col: number) => void;
    onNumberClick: (num: number) => void;
    onActionClick: (action: string) => void;
    isHistoryOpen: boolean;
    setIsHistoryOpen: (isOpen: boolean) => void;
    isHelpOpen: boolean;
    setIsHelpOpen: (isOpen: boolean) => void;
    onLoadGame: () => void;
}

export const GameLayout: React.FC<GameLayoutProps> = ({
    gameState,
    timerSeconds,
    onCellClick,
    onNumberClick,
    onActionClick,
    isHistoryOpen,
    setIsHistoryOpen,
    isHelpOpen,
    setIsHelpOpen,
    onLoadGame
}) => {
    if (!gameState) {
        return <div className="flex items-center justify-center h-screen text-white">Loading Golden Fox Sudoku...</div>;
    }

    const formatTime = (totalSeconds: number) => {
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;
        return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    };

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
                                onCellClick={onCellClick}
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
                                onNumberClick={onNumberClick}
                                onActionClick={onActionClick}
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
                            timeElapsed={formatTime(timerSeconds)}
                            difficulty={gameState.difficulty}
                            userLevel={gameState.userLevel}
                            gamesPlayed={gameState.gamesPlayed}
                            averageTime={gameState.averageTime}
                            winRate={gameState.winRate}
                            pendingGames={gameState.pendingGames}
                            currentDifficultyCount={gameState.currentDifficultyCount}
                            progress={gameState.progress}
                            remainingCells={gameState.remainingCells}
                        />
                    </div>
                </div>
            </div>

            <HistoryModal
                isOpen={isHistoryOpen}
                onClose={() => setIsHistoryOpen(false)}
                onLoadGame={onLoadGame}
            />

            <HelpModal
                isOpen={isHelpOpen}
                onClose={() => setIsHelpOpen(false)}
            />

        </div>
    );
};
