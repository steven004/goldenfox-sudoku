import React from 'react';
import { Clock, Zap, Hash } from 'lucide-react';

interface TimerDisplayProps {
    timeElapsed: string;
    difficulty: string;
    difficultyIndex: number;
    remainingCells: number;
}

export const TimerDisplay: React.FC<TimerDisplayProps> = ({
    timeElapsed,
    difficulty,
    difficultyIndex,
    remainingCells
}) => {
    return (
        <div className="flex flex-col items-center bg-black/20 rounded-lg p-5 shadow-[inset_0_2px_4px_rgba(0,0,0,0.4),0_1px_0_rgba(255,255,255,0.1)] border border-white/5">
            <div className="flex items-center gap-2 mb-2">
                <Clock size={18} className="text-sudoku-primary" />
                <span className="text-sm text-gray-500 font-bold uppercase tracking-wider">Time</span>
            </div>
            <div className="text-5xl font-mono font-bold text-sudoku-primary drop-shadow-md leading-none mb-4">
                {timeElapsed}
            </div>

            <div className="w-full grid grid-cols-[1fr_auto_1fr] items-start px-4">
                {/* Column 1: Difficulty */}
                <div className="flex flex-col items-center gap-1">
                    <div className="flex items-center gap-1.5 text-gray-400 h-5">
                        <Zap size={14} />
                        <span className="text-xs font-bold uppercase">Diff</span>
                    </div>
                    <div className="flex flex-col items-center">
                        <span className="text-lg font-bold text-white leading-none">{difficulty}</span>
                        <span className="text-xs text-gray-400 mt-1 font-mono">({difficultyIndex.toFixed(1)})</span>
                    </div>
                </div>

                {/* Divider */}
                <div className="w-px h-12 bg-white/10 mx-2 self-center" />

                {/* Column 3: Left */}
                <div className="flex flex-col items-center gap-1">
                    <div className="flex items-center gap-1.5 text-gray-400 h-5">
                        <Hash size={14} />
                        <span className="text-xs font-bold uppercase">Left</span>
                    </div>
                    <span className="text-2xl font-bold text-white leading-none mt-1">{remainingCells}</span>
                </div>
            </div>
        </div>
    );
};
