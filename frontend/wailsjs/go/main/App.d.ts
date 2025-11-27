import { engine } from '../models';
import { game } from '../models';
import { main } from '../models';

export function Greet(arg1: string): Promise<string>;

export function ClearCell(): Promise<void>;

export function GetBoard(): Promise<engine.SudokuBoard>;

export function NewGame(arg1: string): Promise<void>;

export function SelectCell(arg1: number, arg2: number): Promise<void>;

export function InputNumber(arg1: number): Promise<void>;

export function TogglePencilMode(): Promise<boolean>;

export function GetGameState(): Promise<any>; // Using any for GameState for simplicity or define interface

export function GetHistory(): Promise<Array<game.PuzzleRecord>>;

export function LoadGame(arg1: string): Promise<void>;

export function SaveGame(): Promise<void>;
