import { useMemo } from 'react';
import LayoutConfig from '../config.json';

export const useLayoutStyles = () => {
    return useMemo(() => {
        return {
            '--app-side-length': `${LayoutConfig.APP_SIDE_LENGTH}px`,
            '--scale-base': LayoutConfig.SCALE_BASE,
            '--p-margin-top': LayoutConfig.MARGIN_TOP,
            '--p-margin-left': LayoutConfig.MARGIN_LEFT,
            '--p-header-height': LayoutConfig.HEADER_HEIGHT,
            '--p-board-size': LayoutConfig.BOARD_SIZE,
            '--p-info-bar-width': LayoutConfig.INFO_BAR_WIDTH,
            '--p-info-bar-height': LayoutConfig.INFO_BAR_HEIGHT,
            '--p-control-bar-height': LayoutConfig.CONTROL_BAR_HEIGHT,
            '--p-control-bar-width': LayoutConfig.CONTROL_BAR_WIDTH,
            '--p-gap-board-info': LayoutConfig.GAP_BOARD_INFO,
            '--p-gap-board-control': LayoutConfig.GAP_BOARD_CONTROL,
        } as React.CSSProperties;
    }, []);
};
