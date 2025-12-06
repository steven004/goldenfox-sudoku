import React from 'react';
import {
    Clock,
    Zap,
    Hash,
    Eraser,
    RotateCcw,
    Crown,
    Target,
    Gamepad2,
    Percent,
    Hourglass,
    Timer,
    Layers
} from 'lucide-react';

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


    const StatRow = ({ icon: Icon, label, value, color = "text-white" }: { icon: any, label: string, value: string | number, color?: string }) => (
        <div className="flex items-center justify-between py-3 border-b border-white/5 last:border-0">
            <div className="flex items-center gap-3 text-gray-400">
                <Icon size={18} />
                <span className="text-sm font-bold uppercase tracking-wider">{label}</span>
            </div>
            <span className={`text-base font-bold ${color}`}>{value}</span>
        </div>
    );

    return (
        <div className="flex flex-col gap-4 bg-gradient-to-b from-[#323846] to-[#1E222D] p-5 rounded-xl w-full backdrop-blur-sm border-4 border-[#FF9F43] shadow-[inset_0_1px_1px_rgba(255,255,255,0.3),inset_0_-1px_1px_rgba(0,0,0,0.5),0_10px_30px_rgba(0,0,0,0.6),0_0_0_2px_rgba(255,159,67,0.3)] relative overflow-hidden h-full">
            {/* Glossy sheen effect */}
            <div className="absolute top-0 left-0 w-full h-1/3 bg-gradient-to-b from-white/5 to-transparent pointer-events-none" />

            {/* SECTION 1: Current Game Info */}
            <div className="flex flex-col items-center bg-black/20 rounded-lg p-5 shadow-[inset_0_2px_4px_rgba(0,0,0,0.4),0_1px_0_rgba(255,255,255,0.1)] border border-white/5">
                <div className="flex items-center gap-2 mb-2">
                    <Clock size={18} className="text-fox-orange" />
                    <span className="text-sm text-gray-500 font-bold uppercase tracking-wider">Time</span>
                </div>
                <div className="text-5xl font-mono font-bold text-fox-orange drop-shadow-md leading-none mb-4">
                    {timeElapsed}
                </div>

                <div className="w-full flex justify-between items-center px-4">
                    <div className="flex flex-col items-center gap-1">
                        <div className="flex items-center gap-1.5 text-gray-400">
                            <Zap size={14} />
                            <span className="text-xs font-bold uppercase">Diff</span>
                        </div>
                        <span className="text-lg font-bold text-white">
                            {difficulty} <span className="text-sm text-gray-500">({difficultyIndex.toFixed(1)})</span>
                        </span>
                    </div>
                    <div className="w-px h-10 bg-white/10" />
                    <div className="flex flex-col items-center gap-1">
                        <div className="flex items-center gap-1.5 text-gray-400">
                            <Hash size={14} />
                            <span className="text-xs font-bold uppercase">Left</span>
                        </div>
                        <span className="text-lg font-bold text-white">{remainingCells}</span>
                    </div>
                </div>
            </div>

            {/* SECTION 2: Operations (Erase & Undo) */}
            <div className="flex flex-col gap-4 bg-black/20 rounded-lg p-5 shadow-[inset_0_2px_4px_rgba(0,0,0,0.4),0_1px_0_rgba(255,255,255,0.1)] border border-white/5">
                <div className="flex justify-between items-center border-b border-white/5 pb-3">
                    <div className="flex items-center gap-2 text-gray-400">
                        <Eraser size={18} />
                        <span className="text-sm font-bold uppercase tracking-wider">Erase</span>
                    </div>
                    <div className="flex gap-2">
                        {[1, 2, 3].map(i => (
                            <div key={i} className={`w-4 h-4 rounded-full ${i <= (3 - eraseCount) ? 'bg-fox-orange shadow-[0_0_8px_#FF9F43]' : 'bg-gray-700/50'}`} />
                        ))}
                    </div>
                </div>

                <div className="flex justify-between items-center">
                    <div className="flex items-center gap-2 text-gray-400">
                        <RotateCcw size={18} />
                        <span className="text-sm font-bold uppercase tracking-wider">Undo</span>
                    </div>
                    <div className="flex gap-2">
                        {[1, 2, 3].map(i => (
                            <div key={i} className={`w-4 h-4 rounded-full ${i <= (3 - undoCount) ? 'bg-blue-400 shadow-[0_0_8px_#60A5FA]' : 'bg-gray-700/50'}`} />
                        ))}
                    </div>
                </div>
            </div>

            {/* SECTION 3: User Stats */}
            <div className="flex flex-col bg-black/20 rounded-lg p-5 shadow-[inset_0_2px_4px_rgba(0,0,0,0.4),0_1px_0_rgba(255,255,255,0.1)] border border-white/5 flex-grow overflow-y-auto">

                {/* Level Progress Visual */}
                <div className="mb-6">
                    <div className="flex justify-between text-xs font-bold uppercase tracking-wider text-gray-500 mb-2">
                        <span className={progress < 0 ? "text-red-400" : ""}>
                            {userLevel > 1 ? `Lv.${userLevel - 1}` : "Min"}
                        </span>
                        <span className="text-fox-orange">Lv.{userLevel}</span>
                        <span className={progress > 0 ? "text-green-400" : ""}>
                            {userLevel < 6 ? `Lv.${userLevel + 1}` : "Max"}
                        </span>
                    </div>

                    <div className="relative h-8 flex items-center justify-center">
                        {/* Track Line */}
                        <div className="absolute w-full h-1 bg-gray-700 rounded-full" />

                        {/* Nodes */}
                        <div className="absolute w-full flex justify-between px-1">
                            {/* Loss Zone (-3 to -1) */}
                            {[-3, -2, -1].map((step) => (
                                <div
                                    key={step}
                                    className={`w-2 h-2 rounded-full transition-all duration-300 ${progress <= step ? 'bg-red-500 shadow-[0_0_8px_#EF4444] scale-125' : 'bg-gray-600'
                                        }`}
                                />
                            ))}

                            {/* Neutral Zone (0) */}
                            <div className={`w-3 h-3 rounded-full transition-all duration-300 z-10 ${progress === 0 ? 'bg-white shadow-[0_0_10px_white] scale-125' : 'bg-gray-500'
                                }`} />

                            {/* Win Zone (+1 to +5) */}
                            {[1, 2, 3, 4, 5].map((step) => (
                                <div
                                    key={step}
                                    className={`w-2 h-2 rounded-full transition-all duration-300 ${progress >= step ? 'bg-green-500 shadow-[0_0_8px_#22C55E] scale-125' : 'bg-gray-600'
                                        }`}
                                />
                            ))}
                        </div>

                        {/* Current Position Indicator (Fox Head / Marker) */}
                        {/* We can animate this later, for now the highlighted nodes show position well enough */}
                    </div>

                    <div className="text-center mt-1 text-xs text-gray-400 font-medium">
                        {progress === 0 && "Ready to start"}
                        {progress > 0 && `${5 - progress} wins to Level Up!`}
                        {progress < 0 && `${3 + progress} losses to Demotion`}
                    </div>
                </div>

                <div className="flex flex-col gap-1 justify-between flex-grow">
                    <StatRow icon={Gamepad2} label="Total Games" value={gamesPlayed} />
                    <StatRow icon={Layers} label={`Games (Lv.${userLevel})`} value={currentDifficultyCount} />
                    <StatRow icon={Percent} label="Win Rate" value={`${winRate.toFixed(1)}%`} />
                    <StatRow icon={Hourglass} label="Pending" value={pendingGames} />
                    <StatRow icon={Timer} label="Avg Time" value={averageTime} />
                </div>
            </div>
        </div>
    );
};
