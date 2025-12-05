import { useEffect } from 'react';
import confetti from 'canvas-confetti';

export const useConfetti = (shouldFire: boolean) => {
    useEffect(() => {
        if (shouldFire) {
            const duration = 2000;
            const end = Date.now() + duration;

            const frame = () => {
                confetti({
                    particleCount: 2,
                    angle: 60,
                    spread: 55,
                    origin: { x: 0 },
                    colors: ['#D68D38', '#FFD28F', '#FFFFFF']
                });
                confetti({
                    particleCount: 2,
                    angle: 120,
                    spread: 55,
                    origin: { x: 1 },
                    colors: ['#D68D38', '#FFD28F', '#FFFFFF']
                });

                if (Date.now() < end) {
                    requestAnimationFrame(frame);
                }
            };
            frame();
        }
    }, [shouldFire]);
};
