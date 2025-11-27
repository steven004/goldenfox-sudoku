export namespace game {

    export class PuzzleRecord {
        id: string;
        predefined: string;
        final_state: string;
        is_solved: boolean;
        time_elapsed: number;
        played_at: string;
        difficulty: number;
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
            this.played_at = source["played_at"];
            this.difficulty = source["difficulty"];
            this.mistakes = source["mistakes"];
        }
    }

}
