import React, { useEffect, useState } from 'react';
import { X, Crown, Shield, Star, Award, Zap, Ghost } from 'lucide-react';

interface NewGameModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSelectDifficulty: (difficulty: string) => void;
    userLevel: number;
}

export const NewGameModal: React.FC<NewGameModalProps> = ({
    isOpen,
    onClose,
    onSelectDifficulty,
    userLevel
}) => {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        if (isOpen) {
            setIsVisible(true);
        } else {
            const timer = setTimeout(() => setIsVisible(false), 300);
            return () => clearTimeout(timer);
        }
    }, [isOpen]);

    if (!isVisible) return null;

    const difficulties = [
        { id: 'Beginner', level: 1, label: 'Beginner', icon: Shield, color: 'text-green-400', borderColor: 'border-green-400/50', bgHover: 'hover:bg-green-400/10' },
        { id: 'Easy', level: 2, label: 'Easy', icon: Star, color: 'text-blue-400', borderColor: 'border-blue-400/50', bgHover: 'hover:bg-blue-400/10' },
        { id: 'Medium', level: 3, label: 'Medium', icon: Award, color: 'text-yellow-400', borderColor: 'border-yellow-400/50', bgHover: 'hover:bg-yellow-400/10' },
        { id: 'Hard', level: 4, label: 'Hard', icon: Zap, color: 'text-orange-400', borderColor: 'border-orange-400/50', bgHover: 'hover:bg-orange-400/10' },
        { id: 'Expert', level: 5, label: 'Expert', icon: Crown, color: 'text-red-400', borderColor: 'border-red-400/50', bgHover: 'hover:bg-red-400/10' },
        { id: 'FoxGod', level: 6, label: 'Fox God', icon: Ghost, color: 'text-purple-400', borderColor: 'border-purple-400/50', bgHover: 'hover:bg-purple-400/10' },
    ];

    return (
        <div className={`fixed inset-0 z-50 flex items-center justify-center p-4 transition-opacity duration-300 ${isOpen ? 'opacity-100' : 'opacity-0'}`}>
            <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
            <div className={`relative bg-[#1E222D] border-2 border-[#D68D38] rounded-xl w-full max-w-lg overflow-hidden shadow-[0_0_50px_rgba(0,0,0,0.5)] transform transition-all duration-300 ${isOpen ? 'scale-100 translate-y-0' : 'scale-95 translate-y-4'}`}>

                {/* Header */}
                <div className="flex items-center justify-between p-5 border-b border-white/10 bg-black/20">
                    <h2 className="text-xl font-bold text-white flex items-center gap-2">
                        <Crown className="text-fox-orange" size={24} />
                        Select Challenge
                    </h2>
                    <button onClick={onClose} className="text-gray-400 hover:text-white transition-colors">
                        <X size={24} />
                    </button>
                </div>

                {/* Content */}
                <div className="p-6 grid grid-cols-2 gap-4">
                    {difficulties.map((diff) => {
                        const isLocked = diff.id === 'FoxGod' && userLevel < 5;
                        const isRecommended = diff.level === userLevel;
                        const Icon = diff.icon;

                        return (
                            <button
                                key={diff.id}
                                disabled={isLocked}
                                onClick={() => {
                                    onSelectDifficulty(diff.id);
                                    onClose();
                                }}
                                className={`
                                    relative flex flex-col items-center justify-center p-4 rounded-lg border-2 transition-all duration-200
                                    ${isLocked
                                        ? 'border-gray-700 bg-gray-800/50 opacity-50 cursor-not-allowed grayscale'
                                        : isRecommended
                                            ? `${diff.borderColor} bg-white/5 shadow-[0_0_15px_rgba(255,159,67,0.3)] scale-105 ring-2 ring-offset-2 ring-offset-[#1E222D] ring-[#D68D38] z-10`
                                            : `${diff.borderColor} bg-transparent ${diff.bgHover} hover:scale-105 active:scale-95 shadow-[0_4px_0_rgba(0,0,0,0.2)]`
                                    }
                                `}
                            >
                                {isRecommended && (
                                    <div className="absolute -top-3 bg-[#D68D38] text-[#1E222D] text-[10px] font-bold px-2 py-0.5 rounded-full shadow-lg flex items-center gap-1">
                                        <Award size={10} /> YOUR RANK
                                    </div>
                                )}
                                <Icon size={32} className={`mb-2 ${isLocked ? 'text-gray-500' : diff.color}`} />
                                <span className={`text-lg font-bold ${isLocked ? 'text-gray-500' : 'text-white'}`}>{diff.label}</span>
                                {isLocked && <span className="text-xs text-red-500 mt-1 uppercase font-bold tracking-wider">Locked</span>}
                            </button>
                        );
                    })}
                </div>
            </div>
        </div>
    );
};
