import { useState, useEffect, useCallback } from 'react';

// Import sound files
// Note: These imports rely on Vite handling .wav files
import clickSound from '../assets/sounds/click.wav';
import popSound from '../assets/sounds/pop.wav';
import errorSound from '../assets/sounds/error.wav';
import eraseSound from '../assets/sounds/erase.wav';
import winSound from '../assets/sounds/win.wav';
import pencilSound from '../assets/sounds/pencil.wav';

type SoundType = 'click' | 'pop' | 'error' | 'erase' | 'win' | 'pencil';

export const useSound = () => {
    const [isMuted, setIsMuted] = useState<boolean>(() => {
        const saved = localStorage.getItem('sudoku-muted');
        return saved ? JSON.parse(saved) : false;
    });

    const [sounds, setSounds] = useState<Record<SoundType, HTMLAudioElement | null>>({
        click: null,
        pop: null,
        error: null,
        erase: null,
        win: null,
        pencil: null,
    });

    useEffect(() => {
        // Initialize Audio objects
        // We use a separate effect/object to avoid re-creating Audio on every render
        const audioMap = {
            click: new Audio(clickSound),
            pop: new Audio(popSound),
            error: new Audio(errorSound),
            erase: new Audio(eraseSound),
            win: new Audio(winSound),
            pencil: new Audio(pencilSound),
        };

        // Preload
        Object.values(audioMap).forEach(audio => {
            audio.load();
            audio.volume = 0.5; // Default volume
        });

        setSounds(audioMap);
    }, []);

    const toggleMute = useCallback(() => {
        setIsMuted(prev => {
            const newState = !prev;
            localStorage.setItem('sudoku-muted', JSON.stringify(newState));
            return newState;
        });
    }, []);

    const playSound = useCallback((type: SoundType) => {
        if (isMuted) return;

        const audio = sounds[type];
        if (audio) {
            // Reset time to allow rapid replay
            audio.currentTime = 0;
            audio.play().catch(e => {
                // Ignore auto-play errors (interaction requirements)
                console.debug("Audio play skipped:", e);
            });
        }
    }, [isMuted, sounds]);

    return { isMuted, toggleMute, playSound };
};
