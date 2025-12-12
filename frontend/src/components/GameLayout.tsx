import React from 'react';
import { CircleHelp, Volume2, VolumeX } from 'lucide-react';
import bgImage from '../assets/images/background.jpg';
import { Board } from './Board';
import { Controls } from './Controls';
import { StatsPanel } from './StatsPanel';
import { HistoryModal } from './HistoryModal';
import { HelpModal } from './HelpModal';
import { NewGameModal } from './NewGameModal';
import { useLayoutStyles } from '../hooks/useLayoutStyles';
import { useGame } from '../context/GameContext';

export const GameLayout: React.FC = () => {
    const layoutStyles = useLayoutStyles();
    const {
        gameState,
        timerSeconds,
        selection,
        pencilMode,
        isMuted,
        toggleMute,
        handleCellClick,
        handleNumberClick,
        handleActionClick,
        handleNewGame,
        refreshState,
        isHistoryOpen,
        setIsHistoryOpen,
        isHelpOpen,
        setIsHelpOpen,
        isNewGameOpen,
        setIsNewGameOpen,
        completionCounts
    } = useGame();

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
                ...layoutStyles
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
                        className="absolute right-8 top-1/2 -translate-y-1/2 p-2 text-sudoku-primary-dark hover:text-sudoku-primary-light transition-colors hover:scale-110 active:scale-95"
                        title="How to Play"
                    >
                        <CircleHelp size={32} strokeWidth={2} />
                    </button>

                    {/* Mute Button - Left of Help */}
                    <button
                        onClick={toggleMute}
                        className="absolute right-24 top-1/2 -translate-y-1/2 p-2 text-sudoku-primary-dark hover:text-sudoku-primary-light transition-colors hover:scale-110 active:scale-95"
                        title={isMuted ? "Unmute" : "Mute"}
                    >
                        {isMuted ? <VolumeX size={32} strokeWidth={2} /> : <Volume2 size={32} strokeWidth={2} />}
                    </button>

                    {gameState?.isSolved && (
                        <div className="text-6xl font-black text-transparent bg-clip-text bg-gradient-to-b from-sudoku-primary-light to-sudoku-primary-dark drop-shadow-[0_2px_10px_rgba(214,141,56,0.5)] animate-in fade-in zoom-in duration-500 tracking-widest">
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
                                selectedRow={selection.row}
                                selectedCol={selection.col}
                                onCellClick={handleCellClick}
                                isPencilMode={pencilMode}
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
                                pencilMode={pencilMode}
                                selectedNumber={
                                    selection.row !== -1 && selection.col !== -1
                                        ? gameState.board.cells[selection.row][selection.col].value
                                        : undefined
                                }
                                completionCounts={completionCounts}
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
                            difficultyIndex={gameState.difficultyIndex}
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
                onLoadGame={refreshState}
                currentUserLevel={gameState.userLevel || 1}
                currentUserProgress={gameState.progress || 0}
            />

            <HelpModal
                isOpen={isHelpOpen}
                onClose={() => setIsHelpOpen(false)}
            />

            <NewGameModal
                isOpen={isNewGameOpen}
                onClose={() => setIsNewGameOpen(false)}
                onSelectDifficulty={handleNewGame}
                userLevel={gameState.userLevel}
            />

        </div>
    );
};
