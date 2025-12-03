import React from 'react';
import { X } from 'lucide-react';

interface HelpModalProps {
    isOpen: boolean;
    onClose: () => void;
}

export const HelpModal: React.FC<HelpModalProps> = ({ isOpen, onClose }) => {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
            <div
                className="relative w-[600px] max-h-[80vh] overflow-y-auto bg-[#1E272E] border-2 border-[#D68D38] rounded-xl shadow-[0_0_30px_rgba(214,141,56,0.2)] p-8 text-[#FDF6E3]"
                onClick={(e) => e.stopPropagation()}
            >
                {/* Close Button */}
                <button
                    onClick={onClose}
                    className="absolute top-4 right-4 p-2 text-[#B2BEC3] hover:text-[#D68D38] transition-colors"
                >
                    <X size={24} />
                </button>

                {/* Header */}
                <div className="text-center mb-8">
                    <h2 className="text-3xl font-bold bg-gradient-to-r from-[#FFD28F] to-[#D68D38] bg-clip-text text-transparent mb-2">
                        How to Play
                    </h2>
                    <div className="h-0.5 w-24 bg-[#D68D38] mx-auto rounded-full opacity-50"></div>
                </div>

                {/* Content */}
                <div className="space-y-8">
                    {/* Rules Section */}
                    <section>
                        <h3 className="text-xl font-semibold text-[#D68D38] mb-4 flex items-center gap-2">
                            📜 Sudoku Rules
                        </h3>
                        <ul className="space-y-3 text-[#B2BEC3] leading-relaxed list-disc pl-5">
                            <li>Fill the grid so that every <strong>row</strong>, <strong>column</strong>, and <strong>3x3 box</strong> contains the digits 1 through 9.</li>
                            <li>Each number can only appear <strong>once</strong> in each row, column, and box.</li>
                            <li>The game starts with some cells already filled (clues). You cannot change these numbers.</li>
                        </ul>
                    </section>

                    {/* Controls Section */}
                    <section>
                        <h3 className="text-xl font-semibold text-[#D68D38] mb-4 flex items-center gap-2">
                            🎮 Controls & Tools
                        </h3>
                        <div className="grid gap-4 text-[#B2BEC3]">
                            <div className="flex items-start gap-3 bg-[#2D3436] p-3 rounded-lg border border-[#636E72]/30">
                                <div className="p-2 bg-[#D68D38]/10 rounded text-[#D68D38]">✏️</div>
                                <div>
                                    <strong className="text-[#FDF6E3] block mb-1">Pencil Mode</strong>
                                    Toggle this to add small notes (candidates) to a cell. Useful for tracking possibilities!
                                </div>
                            </div>
                            <div className="flex items-start gap-3 bg-[#2D3436] p-3 rounded-lg border border-[#636E72]/30">
                                <div className="p-2 bg-[#D68D38]/10 rounded text-[#D68D38]">↩️</div>
                                <div>
                                    <strong className="text-[#FDF6E3] block mb-1">Undo</strong>
                                    Made a mistake? You can undo your last 3 moves.
                                </div>
                            </div>
                            <div className="flex items-start gap-3 bg-[#2D3436] p-3 rounded-lg border border-[#636E72]/30">
                                <div className="p-2 bg-[#D68D38]/10 rounded text-[#D68D38]">🧹</div>
                                <div>
                                    <strong className="text-[#FDF6E3] block mb-1">Erase</strong>
                                    Clear a cell's content. You have 3 erase chances per game!
                                </div>
                            </div>
                        </div>
                    </section>

                    {/* Tips Section */}
                    <section>
                        <h3 className="text-xl font-semibold text-[#D68D38] mb-4 flex items-center gap-2">
                            🦊 Fox Wisdom
                        </h3>
                        <p className="text-[#B2BEC3] italic border-l-4 border-[#D68D38] pl-4 py-1 bg-[#D68D38]/5 rounded-r">
                            "Focus on rows, columns, or boxes that are almost full. Use Pencil Mode to mark potential numbers when you're unsure."
                        </p>
                    </section>
                </div>

                {/* Footer */}
                <div className="mt-8 pt-6 border-t border-[#636E72]/30 text-center">
                    <button
                        onClick={onClose}
                        className="px-8 py-2 bg-[#D68D38] hover:bg-[#E17055] text-[#1E272E] font-bold rounded-full transition-all transform hover:scale-105 active:scale-95 shadow-lg shadow-[#D68D38]/20"
                    >
                        Got it!
                    </button>
                </div>
            </div>
        </div>
    );
};
