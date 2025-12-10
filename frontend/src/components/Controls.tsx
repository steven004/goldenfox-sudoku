import React from 'react';
import { Edit3, Eraser, RotateCcw, RefreshCw, Plus, History as HistoryIcon } from 'lucide-react';

interface ControlsProps {
    onNumberClick: (num: number) => void;
    onActionClick: (action: string) => void;
    pencilMode: boolean;
    selectedNumber?: number;
    completionCounts: Record<number, number>;
}

export const Controls: React.FC<ControlsProps> = ({ onNumberClick, onActionClick, pencilMode, selectedNumber, completionCounts }) => {
    return (
        <div className="flex flex-col gap-6 w-full">
            {/* Row 1: Number Pad (Glass Orbs) */}
            <div className="flex justify-between items-center px-4">
                {[1, 2, 3, 4, 5, 6, 7, 8, 9].map(num => {
                    const isSelected = selectedNumber === num;
                    const isCompleted = (completionCounts[num] || 0) >= 9;

                    // Determine styling based on state priority: Selected > Completed > Normal
                    let buttonStyle = '';
                    let inlineStyle = {};

                    if (isSelected) {
                        // 1. Selected State (Glowing Sun Stone)
                        // Bright radial gradient, white border, dark text, strong outer glow
                        buttonStyle = 'text-[#1a1a1a] border-2 border-white shadow-[0_0_20px_#FF9F43,inset_0_2px_4px_rgba(255,255,255,0.5)] scale-110 z-10';
                        inlineStyle = { background: 'radial-gradient(circle at 30% 30%, #FFD28F, #FF9F43, #D68D38)' };
                    } else if (isCompleted) {
                        // 2. Completed State (Dim Stone)
                        // Flat, dark, low contrast
                        buttonStyle = 'text-[#666] bg-[#1a1a1a] border border-[#333] shadow-[inset_0_2px_4px_rgba(0,0,0,0.5)] cursor-default';
                    } else {
                        // 3. Normal State (Amber Glass Stone)
                        // Convex look, deep amber, glossy highlight
                        buttonStyle = 'text-sudoku-primary-light border border-sudoku-primary-dark/50 shadow-[inset_0_2px_4px_rgba(255,255,255,0.2),inset_0_-2px_4px_rgba(0,0,0,0.4),0_4px_8px_rgba(0,0,0,0.4)] hover:scale-105 active:scale-95 transition-transform';
                        inlineStyle = { background: 'radial-gradient(circle at 30% 30%, rgba(214,141,56,0.3), rgba(214,141,56,0.1))' };
                    }

                    return (
                        <button
                            key={num}
                            className={`
                                relative flex items-center justify-center
                                w-14 h-14 rounded-full transition-all duration-200
                                text-2xl font-bold
                                ${buttonStyle}
                            `}
                            style={inlineStyle}
                            onClick={() => onNumberClick(num)}
                            disabled={isCompleted && !isSelected}
                        >
                            {num}
                        </button>
                    );
                })}
            </div>

            {/* Row 2: Tool Actions (Tactile Buttons) */}
            <div className="flex justify-between gap-2 px-2">
                {[
                    { id: 'pencil', label: 'Pencil', icon: Edit3, active: pencilMode },
                    { id: 'eraser', label: 'Eraser', icon: Eraser },
                    { id: 'undo', label: 'Undo', icon: RotateCcw },
                    { id: 'restart', label: 'Restart', icon: RefreshCw },
                    { id: 'new', label: 'New', icon: Plus },
                    { id: 'history', label: 'History', icon: HistoryIcon },
                ].map(action => {
                    const Icon = action.icon;
                    const isActive = action.active;

                    return (
                        <button
                            key={action.id}
                            className={`
                                flex flex-col items-center justify-center
                                py-2 px-3 rounded-xl border transition-all duration-100
                                min-w-[4.5rem] relative overflow-hidden
                                ${isActive
                                    ? 'bg-gradient-to-b from-sudoku-primary-dark to-sudoku-primary-darker border-sudoku-primary-darker text-[#121418] shadow-[inset_0_2px_4px_rgba(0,0,0,0.3)] translate-y-[1px]'
                                    : 'bg-gradient-to-b from-sudoku-panel to-sudoku-panel-dark border-[#3E4552] text-sudoku-primary-dark shadow-[inset_0_1px_0_rgba(255,255,255,0.1),0_4px_0_#151820,0_5px_10px_rgba(0,0,0,0.5)] hover:translate-y-[1px] hover:shadow-[inset_0_1px_0_rgba(255,255,255,0.1),0_3px_0_#151820,0_4px_8px_rgba(0,0,0,0.5)] active:translate-y-[4px] active:shadow-[inset_0_2px_4px_rgba(0,0,0,0.4)]'
                                }
                            `}
                            onClick={() => onActionClick(action.id)}
                        >
                            <Icon size={24} strokeWidth={2} className="drop-shadow-sm" />
                            <span className="text-[10px] font-bold mt-1 uppercase tracking-wider drop-shadow-sm">{action.label}</span>
                        </button>
                    );
                })}
            </div>
        </div>
    );
};
