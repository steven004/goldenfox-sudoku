import React, { memo } from 'react';
import { Gamepad2, Percent, Hourglass, Timer, Layers } from 'lucide-react';

interface StatListProps {
    gamesPlayed: number;
    userLevel: number;
    currentDifficultyCount: number;
    winRate: number;
    pendingGames: number;
    averageTime: string;
}

const StatListComponent: React.FC<StatListProps> = ({
    gamesPlayed,
    userLevel,
    currentDifficultyCount,
    winRate,
    pendingGames,
    averageTime
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
        <div className="flex flex-col gap-1 justify-between flex-grow">
            <StatRow icon={Gamepad2} label="Total Games" value={gamesPlayed} />
            <StatRow icon={Layers} label={`Games (Lv.${userLevel})`} value={currentDifficultyCount} />
            <StatRow icon={Percent} label="Win Rate" value={`${winRate.toFixed(1)}%`} />
            <StatRow icon={Hourglass} label="Pending" value={pendingGames} />
            <StatRow icon={Timer} label="Avg Time" value={averageTime} />
        </div>
    );
};

export const StatList = memo(StatListComponent);
