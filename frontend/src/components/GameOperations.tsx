import React, { memo } from 'react';
import { Eraser, RotateCcw } from 'lucide-react';

interface GameOperationsProps {
    eraseCount: number;
    undoCount: number;
}

const GameOperationsComponent: React.FC<GameOperationsProps> = ({ eraseCount, undoCount }) => {
    return (
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
    );
};

export const GameOperations = memo(GameOperationsComponent);
