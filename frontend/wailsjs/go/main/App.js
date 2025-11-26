export function Greet(arg1) {
    return window['go']['main']['App']['Greet'](arg1);
}

export function GetBoard() {
    return window['go']['main']['App']['GetBoard']();
}

export function NewGame(arg1) {
    return window['go']['main']['App']['NewGame'](arg1);
}

export function SelectCell(arg1, arg2) {
    return window['go']['main']['App']['SelectCell'](arg1, arg2);
}

export function InputNumber(arg1) {
    return window['go']['main']['App']['InputNumber'](arg1);
}

export function TogglePencilMode() {
    return window['go']['main']['App']['TogglePencilMode']();
}

export function GetGameState() {
    return window['go']['main']['App']['GetGameState']();
}
