import { useState, useEffect } from 'react';
import { GetHistory, LoadGame } from '../../wailsjs/go/main/App';
import { game } from '../../wailsjs/go/models';

interface HistoryModalProps {
    isOpen: boolean;
    onClose: () => void;
    onLoadGame: () => void; // Callback to refresh parent state
}

export const HistoryModal = ({ isOpen, onClose, onLoadGame }: HistoryModalProps) => {
    const [history, setHistory] = useState<game.PuzzleRecord[]>([]);
    const [activeTab, setActiveTab] = useState<'uncompleted' | 'finished'>('uncompleted');

    useEffect(() => {
        if (isOpen) {
            try {
                GetHistory().then((records: game.PuzzleRecord[]) => {
                    if (!Array.isArray(records)) {
                        console.warn("GetHistory returned non-array:", records);
                        setHistory([]);
                        return;
                    }
                    // Sort by date descending
                    const sorted = records.sort((a: game.PuzzleRecord, b: game.PuzzleRecord) => {
                        const dateA = new Date(a.played_at).getTime();
                        const dateB = new Date(b.played_at).getTime();
                        return (isNaN(dateB) ? 0 : dateB) - (isNaN(dateA) ? 0 : dateA);
                    });
                    setHistory(sorted);
                }).catch(err => {
                    console.error("Failed to fetch history:", err);
                    setHistory([]);
                });
            } catch (err) {
                console.error("Critical error calling GetHistory (Backend might need restart):", err);
                setHistory([]);
            }
        }
    }, [isOpen]);

    if (!isOpen) return null;

    const filteredHistory = history.filter(record =>
        activeTab === 'finished' ? record.is_solved : !record.is_solved
    );

    const handleLoad = async (id: string) => {
        try {
            await LoadGame(id);
            onLoadGame();
            onClose();
        } catch (err) {
            console.error("Failed to load game:", err);
            alert("Failed to load game. See console for details.");
        }
    };

    const formatTime = (ns: number) => {
        const seconds = Math.floor(ns / 1000000000);
        const m = Math.floor(seconds / 60);
        const s = seconds % 60;
        return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
    };

    const formatDate = (dateStr: string) => {
        const d = new Date(dateStr);
        return isNaN(d.getTime()) ? 'Unknown Date' : d.toLocaleString();
    };

    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
            onClick={onClose} // Close on backdrop click
        >
            <div
                className="bg-[#2A2A2A] w-full max-w-2xl rounded-xl shadow-2xl border border-[#D68D38]/30 flex flex-col max-h-[80vh]"
                onClick={e => e.stopPropagation()} // Prevent close on modal click
            >
                {/* Header */}
                <div className="flex items-center justify-between p-6 border-b border-white/10">
                    <h2 className="text-2xl font-bold text-white font-display tracking-wide">Game History</h2>
                    <button onClick={onClose} className="text-white/50 hover:text-white transition-colors p-2">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                {/* Tabs */}
                <div className="flex border-b border-white/10">
                    <button
                        className={`flex-1 py-4 text-sm font-bold uppercase tracking-wider transition-colors ${activeTab === 'uncompleted' ? 'bg-[#D68D38] text-[#1a1a1a]' : 'text-white/70 hover:bg-white/5'}`}
                        onClick={() => setActiveTab('uncompleted')}
                    >
                        Uncompleted
                    </button>
                    <button
                        className={`flex-1 py-4 text-sm font-bold uppercase tracking-wider transition-colors ${activeTab === 'finished' ? 'bg-[#D68D38] text-[#1a1a1a]' : 'text-white/70 hover:bg-white/5'}`}
                        onClick={() => setActiveTab('finished')}
                    >
                        Finished
                    </button>
                </div>

                {/* List */}
                <div className="flex-1 overflow-y-auto p-4 space-y-3">
                    {filteredHistory.length === 0 ? (
                        <div className="text-center text-white/30 py-12 italic">
                            No {activeTab} games found.
                        </div>
                    ) : (
                        filteredHistory.map((record) => (
                            <div key={record.id} className="bg-white/5 rounded-lg p-4 flex items-center justify-between hover:bg-white/10 transition-colors border border-white/5">
                                <div>
                                    <div className="flex items-center gap-3 mb-1">
                                        <span className={`text-xs font-bold px-2 py-0.5 rounded ${record.difficulty === 0 ? 'bg-green-500/20 text-green-400' :
                                            record.difficulty === 1 ? 'bg-yellow-500/20 text-yellow-400' :
                                                'bg-red-500/20 text-red-400'
                                            }`}>
                                            {record.difficulty === 0 ? 'EASY' : record.difficulty === 1 ? 'MEDIUM' : 'HARD'}
                                        </span>
                                        <span className="text-white/50 text-xs">{formatDate(record.played_at as any)}</span>
                                    </div>
                                    <div className="text-white font-mono text-lg">
                                        {formatTime(record.time_elapsed)}
                                    </div>
                                </div>
                                <button
                                    onClick={() => handleLoad(record.id)}
                                    className="bg-[#D68D38] text-[#1a1a1a] px-4 py-2 rounded font-bold text-sm hover:bg-[#FFD28F] transition-colors shadow-lg"
                                >
                                    LOAD
                                </button>
                            </div>
                        ))
                    )}
                </div>

                {/* Footer Close Button */}
                <div className="p-4 border-t border-white/10 flex justify-end">
                    <button
                        onClick={onClose}
                        className="text-white/70 hover:text-white px-4 py-2 rounded hover:bg-white/10 transition-colors text-sm font-bold"
                    >
                        CLOSE
                    </button>
                </div>
            </div>
        </div>
    );
};
