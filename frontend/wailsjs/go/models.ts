export namespace engine {
	
	export class Cell {
	    value: number;
	    given: boolean;
	    candidates: {[key: number]: boolean};
	    isInvalid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Cell(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	        this.given = source["given"];
	        this.candidates = source["candidates"];
	        this.isInvalid = source["isInvalid"];
	    }
	}
	export class SudokuBoard {
	    cells: Cell[][];
	
	    static createFrom(source: any = {}) {
	        return new SudokuBoard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cells = this.convertValues(source["cells"], Cell);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace game {
	
	export class GameState {
	    board: engine.SudokuBoard;
	    mistakes: number;
	    eraseCount: number;
	    undoCount: number;
	    elapsedSeconds: number;
	    difficulty: string;
	    difficultyIndex: number;
	    isSolved: boolean;
	    userLevel: number;
	    gamesPlayed: number;
	    winRate: number;
	    pendingGames: number;
	    averageTime: string;
	    currentDifficultyCount: number;
	    progress: number;
	    remainingCells: number;
	
	    static createFrom(source: any = {}) {
	        return new GameState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.board = this.convertValues(source["board"], engine.SudokuBoard);
	        this.mistakes = source["mistakes"];
	        this.eraseCount = source["eraseCount"];
	        this.undoCount = source["undoCount"];
	        this.elapsedSeconds = source["elapsedSeconds"];
	        this.difficulty = source["difficulty"];
	        this.difficultyIndex = source["difficultyIndex"];
	        this.isSolved = source["isSolved"];
	        this.userLevel = source["userLevel"];
	        this.gamesPlayed = source["gamesPlayed"];
	        this.winRate = source["winRate"];
	        this.pendingGames = source["pendingGames"];
	        this.averageTime = source["averageTime"];
	        this.currentDifficultyCount = source["currentDifficultyCount"];
	        this.progress = source["progress"];
	        this.remainingCells = source["remainingCells"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PuzzleRecord {
	    id: string;
	    predefined: string;
	    final_state: string;
	    is_solved: boolean;
	    time_elapsed: number;
	    // Go type: time
	    played_at: any;
	    difficulty: number;
	    difficulty_index: number;
	    mistakes: number;
	
	    static createFrom(source: any = {}) {
	        return new PuzzleRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.predefined = source["predefined"];
	        this.final_state = source["final_state"];
	        this.is_solved = source["is_solved"];
	        this.time_elapsed = source["time_elapsed"];
	        this.played_at = this.convertValues(source["played_at"], null);
	        this.difficulty = source["difficulty"];
	        this.difficulty_index = source["difficulty_index"];
	        this.mistakes = source["mistakes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

