import React from 'react';
import { TimerDisplay } from './TimerDisplay';
import { GameOperations } from './GameOperations';
import { LevelProgressBar } from './LevelProgressBar';
import { StatList } from './StatList';

interface StatsPanelProps {
    eraseCount: number;
    undoCount: number;
    timeElapsed: string;
    difficulty: string;
    difficultyIndex: number;
    userLevel: number;
    gamesPlayed: number;
    averageTime: string;
    winRate: number;
    pendingGames: number;
    currentDifficultyCount: number;
    progress: number;
    remainingCells: number;
}

export const StatsPanel: React.FC<StatsPanelProps> = ({
    eraseCount,
    undoCount,
    timeElapsed,
    difficulty,
    difficultyIndex,
    userLevel,
    gamesPlayed,
    averageTime,
    winRate,
    pendingGames,
    currentDifficultyCount,
    progress,
    remainingCells
}) => {
    return (
        <div className="flex flex-col gap-4 bg-gradient-to-b from-sudoku-panel to-sudoku-panel-dark p-5 rounded-xl w-full backdrop-blur-sm border-4 border-sudoku-primary shadow-[inset_0_1px_1px_rgba(255,255,255,0.3),inset_0_-1px_1px_rgba(0,0,0,0.5),0_10px_30px_rgba(0,0,0,0.6),0_0_0_2px_rgba(255,159,67,0.3)] relative overflow-hidden h-full">
            {/* Glossy sheen effect */}
            <div className="absolute top-0 left-0 w-full h-1/3 bg-gradient-to-b from-white/5 to-transparent pointer-events-none" />

            {/* SECTION 1: Current Game Info (Timer updates every second) */}
            <TimerDisplay
                timeElapsed={timeElapsed}
                difficulty={difficulty}
                difficultyIndex={difficultyIndex}
                remainingCells={remainingCells}
            />

            {/* SECTION 2: Operations (Memoized) */}
            <GameOperations
                eraseCount={eraseCount}
                undoCount={undoCount}
            />

            {/* SECTION 3: User Stats (Memoized) */}
            <div className="flex flex-col bg-black/20 rounded-lg p-5 shadow-[inset_0_2px_4px_rgba(0,0,0,0.4),0_1px_0_rgba(255,255,255,0.1)] border border-white/5 flex-grow overflow-y-auto">
                <LevelProgressBar
                    userLevel={userLevel}
                    progress={progress}
                />

                <StatList
                    gamesPlayed={gamesPlayed}
                    userLevel={userLevel}
                    currentDifficultyCount={currentDifficultyCount}
                    winRate={winRate}
                    pendingGames={pendingGames}
                    averageTime={averageTime}
                />
            </div>
        </div>
    );
};
