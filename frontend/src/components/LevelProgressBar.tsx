import React, { memo } from 'react';

interface LevelProgressBarProps {
    userLevel: number;
    progress: number;
}

const LevelProgressBarComponent: React.FC<LevelProgressBarProps> = ({ userLevel, progress }) => {
    return (
        <div className="mb-6">
            <div className="flex justify-between text-xs font-bold uppercase tracking-wider text-gray-500 mb-2">
                <span className={progress < 0 ? "text-red-400" : ""}>
                    {userLevel > 1 ? `Lv.${userLevel - 1}` : "Min"}
                </span>
                <span className="text-sudoku-primary">Lv.{userLevel}</span>
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
            </div>

            <div className="text-center mt-1 text-xs text-gray-400 font-medium">
                {progress === 0 && "Ready to start"}
                {progress > 0 && `${5 - progress} wins to Level Up!`}
                {progress < 0 && `${3 + progress} losses to Demotion`}
            </div>
        </div>
    );
};

export const LevelProgressBar = memo(LevelProgressBarComponent);
