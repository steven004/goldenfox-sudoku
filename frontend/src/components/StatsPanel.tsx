import React from 'react';

interface StatsPanelProps {
    mistakes: number;
    timeElapsed: string;
    difficulty: string;
}

export const StatsPanel: React.FC<StatsPanelProps> = ({ mistakes, timeElapsed, difficulty }) => {
    return (
        <div className="flex flex-col gap-6 bg-[#1E222D]/85 p-6 rounded-xl w-full backdrop-blur-sm border border-[#FF9F43]/50 shadow-[0_4px_20px_rgba(0,0,0,0.3),0_0_15px_rgba(255,159,67,0.2)]">
            {/* Clock */}
            <div className="flex flex-col items-center">
                <div className="text-4xl font-bold text-fox-orange font-mono tracking-wider">
                    {timeElapsed}
                </div>
                <div className="text-xs text-gray-500 mt-1">TIME</div>
            </div>

            <div className="h-px bg-gray-700 w-full" />

            {/* Mistakes */}
            <div className="flex flex-col items-center">
                <div className="text-sm text-gray-400 mb-2">Mistakes (0/3)</div>
                <div className="flex gap-2">
                    {[1, 2, 3].map(i => (
                        <span key={i} className={`text-2xl font-bold ${i <= mistakes ? 'text-red-500' : 'text-gray-700'}`}>
                            X
                        </span>
                    ))}
                </div>
            </div>

            <div className="h-px bg-gray-700 w-full" />

            {/* Stats Grid */}
            <div className="flex flex-col gap-4">
                <div className="flex flex-col">
                    <span className="text-xs text-gray-500">Difficulty</span>
                    <span className="text-lg font-bold text-white">{difficulty}</span>
                </div>
                <div className="flex flex-col">
                    <span className="text-xs text-gray-500">User Level</span>
                    <span className="text-lg font-bold text-white">Intermediate</span>
                </div>
                <div className="flex flex-col">
                    <span className="text-xs text-gray-500">Games Played</span>
                    <span className="text-lg font-bold text-white">45</span>
                </div>
                <div className="flex flex-col">
                    <span className="text-xs text-gray-500">Avg Time</span>
                    <span className="text-lg font-bold text-white">08:23</span>
                </div>
            </div>
        </div>
    );
};
