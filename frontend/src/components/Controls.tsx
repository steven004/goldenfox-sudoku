import React from 'react';
import { Edit3, Eraser, RotateCcw, RefreshCw, Plus, Save, Download } from 'lucide-react';

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
                        // 1. Selected State (Burning Sun)
                        // Solid Radial Gradient, White Border, Dark Text, Strong Glow
                        buttonStyle = 'text-[#1a1a1a] border border-white shadow-[0_0_15px_#FF9F43] scale-110 z-10';
                        inlineStyle = { background: 'radial-gradient(circle at 30% 30%, #FFD28F, #FF9F43, #D68D38)' };
                    } else if (isCompleted) {
                        // 2. Completed State (Burnt Out)
                        // Transparent bg, Orange border (circle), dark grey text
                        buttonStyle = 'text-[#444] bg-transparent border border-[rgba(214,141,56,0.5)] cursor-default';
                    } else {
                        // 3. Normal State (Amber Glass)
                        // Subtle Orange tint bg, visible Orange ring, Pale Gold text
                        buttonStyle = 'text-[#FFD28F] bg-[rgba(214,141,56,0.15)] border border-[rgba(214,141,56,0.5)] hover:bg-[rgba(214,141,56,0.25)] hover:text-white shadow-[0_2px_5px_rgba(0,0,0,0.3)]';
                    }

                    return (
                        <button
                            key={num}
                            className={`
                                relative flex items-center justify-center
                                w-14 h-14 rounded-full transition-all duration-300
                                text-2xl font-bold
                                ${buttonStyle}
                            `}
                            style={inlineStyle}
                            onClick={() => onNumberClick(num)}
                            disabled={isCompleted && !isSelected} // Optional: disable clicking if completed? User didn't specify, but implies "ignores it". Keeping clickable for now unless specified otherwise, but visual style implies disabled.
                        >
                            {num}
                        </button>
                    );
                })}
            </div>

            {/* Row 2: Tool Actions (Outline Capsules) */}
            <div className="flex justify-between gap-2 px-2">
                {[
                    { id: 'pencil', label: 'Pencil', icon: Edit3, active: pencilMode },
                    { id: 'eraser', label: 'Eraser', icon: Eraser },
                    { id: 'undo', label: 'Undo', icon: RotateCcw },
                    { id: 'restart', label: 'Restart', icon: RefreshCw },
                    { id: 'new', label: 'New', icon: Plus },
                    { id: 'save', label: 'Save', icon: Save },
                    { id: 'load', label: 'Load', icon: Download },
                ].map(action => {
                    const Icon = action.icon;
                    const isActive = action.active;

                    return (
                        <button
                            key={action.id}
                            className={`
                                flex flex-col items-center justify-center
                                py-2 px-3 rounded-2xl border transition-all duration-300
                                min-w-[4.5rem]
                                ${isActive
                                    ? 'bg-[#D68D38] border-[#D68D38] text-[#121418]'
                                    : 'bg-transparent border-[#D68D38] text-[#D68D38] hover:bg-[#D68D38]/10'
                                }
                            `}
                            onClick={() => onActionClick(action.id)}
                        >
                            <Icon size={24} strokeWidth={2} />
                            <span className="text-[10px] font-bold mt-1 uppercase tracking-wider">{action.label}</span>
                        </button>
                    );
                })}
            </div>
        </div>
    );
};
