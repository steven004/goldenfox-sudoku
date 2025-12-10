package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/steven004/goldenfox-sudoku/engine"
)

// GameManager manages the lifecycle of the application and user data
type GameManager struct {
	// Infrastructure
	generator engine.PuzzleGenerator
	userData  *UserData
	mu        sync.RWMutex // Protects currentSession

	// Active Session
	currentSession *GameSession
}

// NewGameManager creates a new game manager
func NewGameManager(generator engine.PuzzleGenerator) *GameManager {
	// Load user data
	ud, err := LoadUserData("user_data.json")
	if err != nil {

		ud = NewUserData()
	}

	return &GameManager{
		generator: generator,
		userData:  ud,
	}
}

// NewGame generates and starts a new puzzle
func (gm *GameManager) NewGame(difficultyOverride string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Check for abandonment of previous game
	if gm.currentSession != nil && !gm.currentSession.currentBoard.IsSolved() && gm.currentSession.endTime.IsZero() {
		// Record Loss
		gm.userData.RecordLoss(gm.currentSession.difficulty)
		// Save
		gm.saveSessionLocked()
	}

	var targetDifficulty engine.DifficultyLevel

	if difficultyOverride != "" {
		switch difficultyOverride {
		case "Beginner":
			targetDifficulty = engine.Beginner
		case "Easy":
			targetDifficulty = engine.Easy
		case "Medium":
			targetDifficulty = engine.Medium
		case "Hard":
			targetDifficulty = engine.Hard
		case "Expert":
			targetDifficulty = engine.Expert
		case "FoxGod":
			targetDifficulty = engine.FoxGod
		default:
			targetDifficulty = gm.getDifficultyFromLevel()
		}
	} else {
		targetDifficulty = gm.getDifficultyFromLevel()
	}

	progress := gm.userData.Stats.Progress
	extraClues := CalculateDynamicClues(targetDifficulty, progress)

	puzzle, diffIndex, err := gm.generator.Generate(targetDifficulty, extraClues)
	if err != nil {
		return fmt.Errorf("failed to generate puzzle: %w", err)
	}

	// Generate Smart ID: 12 digits
	// 1st digit: Difficulty (1-6)
	// 2nd digit: Progress + 4 (Range 2-8)
	// 3-12 digits: Index (padded)

	diffValue := int(targetDifficulty) + 1 // Beginner(0) -> 1
	progValue := progress + 4              // -2 -> 2, +4 -> 8

	// Get total games started at this difficulty (Completed + Pending)
	gameIndex := gm.userData.GetGamesAtDifficulty(targetDifficulty) + 1

	id := fmt.Sprintf("%d%d%010d", diffValue, progValue, gameIndex)

	gm.currentSession = NewGameSession(puzzle, targetDifficulty, diffIndex, id)

	return nil
}

func (gm *GameManager) getDifficultyFromLevel() engine.DifficultyLevel {
	userLevel := gm.userData.Stats.Level
	if userLevel < 1 {
		return engine.Easy
	}
	if userLevel > 6 {
		return engine.FoxGod
	}
	return engine.DifficultyLevel(userLevel - 1)
}

// RestartGame resets the current session
func (gm *GameManager) RestartGame() error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	if gm.currentSession == nil {
		return fmt.Errorf("no game in progress")
	}

	// Create new session with same initial board
	gm.currentSession = NewGameSession(
		gm.currentSession.initialBoard, // this creates a clone inside NewGameSession
		gm.currentSession.difficulty,
		gm.currentSession.difficultyIndex,
		fmt.Sprintf("%d", time.Now().UnixNano()),
	)

	return nil
}

// ---- Delegates ----

func (gm *GameManager) InputNumber(row, col, val int) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if gm.currentSession == nil {
		return fmt.Errorf("no game")
	}

	if err := gm.currentSession.InputNumber(row, col, val); err != nil {
		return err
	}

	// Auto-save (Throttling could go here)
	go gm.SaveCurrentGame()
	return nil
}

func (gm *GameManager) ToggleCandidate(row, col, val int) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if gm.currentSession == nil {
		return fmt.Errorf("no game")
	}

	if err := gm.currentSession.ToggleCandidate(row, col, val); err != nil {
		return err
	}

	go gm.SaveCurrentGame()
	return nil
}

func (gm *GameManager) ClearCell(row, col int) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if gm.currentSession == nil {
		return fmt.Errorf("no game")
	}

	if err := gm.currentSession.ClearCell(row, col); err != nil {
		return err
	}

	go gm.SaveCurrentGame()
	return nil
}

func (gm *GameManager) Undo() error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if gm.currentSession == nil {
		return fmt.Errorf("no game")
	}

	if err := gm.currentSession.Undo(); err != nil {
		return err
	}

	go gm.SaveCurrentGame()
	return nil
}

// ---- Accessors ----

func (gm *GameManager) GetBoard() *engine.SudokuBoard {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	if gm.currentSession == nil {
		return nil
	}
	return gm.currentSession.currentBoard
}

func (gm *GameManager) IsSolved() bool {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	if gm.currentSession == nil {
		return false
	}
	return gm.currentSession.currentBoard.IsSolved()
}

func (gm *GameManager) GetUserLevel() int {
	if gm.userData == nil {
		return 1
	}
	gm.userData.mu.RLock()
	defer gm.userData.mu.RUnlock()
	return gm.userData.Stats.Level
}

func (gm *GameManager) GetHistory() []PuzzleRecord {
	if gm.userData == nil {
		return nil
	}
	// Return combined history
	return gm.userData.GetAllHistory()
}

// GetGameState returns the transfer object
func (gm *GameManager) GetGameState() GameState {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	var board engine.SudokuBoard
	var mistakes, erase, undo, remaining int
	var diff string
	var diffIndex float64
	var isSolved bool

	if gm.currentSession != nil {
		board = *gm.currentSession.currentBoard.Clone()
		mistakes = 0
		erase = gm.currentSession.eraseCount
		undo = gm.currentSession.undoCount
		diff = gm.currentSession.difficulty.String()
		diffIndex = gm.currentSession.difficultyIndex
		isSolved = gm.currentSession.currentBoard.IsSolved()

		// Count remaining
		filled := 0
		for r := 0; r < 9; r++ {
			for c := 0; c < 9; c++ {
				if gm.currentSession.currentBoard.Cells[r][c].Value != 0 {
					filled++
				}
			}
		}
		remaining = 81 - filled
	} else {
		remaining = 81
	}

	// Stats
	level := 1
	gamesPlayed := 0
	winRate := 0.0
	pending := 0
	avgTimeStr := "--:--"
	currDiffCount := 0
	progress := 0

	if gm.userData != nil {
		gm.userData.mu.RLock()
		level = gm.userData.Stats.Level
		// Games Played = Completed Games
		gamesPlayed = len(gm.userData.CompletedHistory)
		winRate = gm.userData.GetWinRate()
		pending = gm.userData.GetPendingGamesCount()
		if gm.currentSession != nil {
			currDiffCount = gm.userData.GetGamesAtDifficulty(gm.currentSession.difficulty)
			if avg, ok := gm.userData.Stats.AverageTimes[gm.currentSession.difficulty]; ok && avg > 0 {
				minutes := int(avg) / 60
				seconds := int(avg) % 60
				avgTimeStr = fmt.Sprintf("%02d:%02d", minutes, seconds)
			}
		}
		progress = gm.userData.Stats.Progress
		gm.userData.mu.RUnlock()
	}

	return GameState{
		Board:                  board,
		Mistakes:               mistakes,
		EraseCount:             erase,
		UndoCount:              undo,
		Difficulty:             diff,
		DifficultyIndex:        diffIndex,
		IsSolved:               isSolved,
		UserLevel:              level,
		GamesPlayed:            gamesPlayed,
		AverageTime:            avgTimeStr,
		WinRate:                winRate,
		PendingGames:           pending,
		CurrentDifficultyCount: currDiffCount,
		Progress:               progress,
		RemainingCells:         remaining,
	}
}

// ---- Persistence ----

// SaveCurrentGame (External) locks and calls internal
func (gm *GameManager) SaveCurrentGame() error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	return gm.saveSessionLocked()
}

func (gm *GameManager) saveSessionLocked() error {
	if gm.currentSession == nil {
		return nil
	}

	elapsed := gm.currentSession.GetElapsedDuration()

	record := PuzzleRecord{
		ID:              gm.currentSession.gameID,
		Predefined:      gm.currentSession.initialBoard.String(),
		FinalState:      gm.currentSession.currentBoard.String(),
		IsSolved:        gm.currentSession.currentBoard.IsSolved(),
		TimeElapsed:     elapsed,
		PlayedAt:        time.Now(),
		Difficulty:      gm.currentSession.difficulty,
		DifficultyIndex: gm.currentSession.difficultyIndex,
		Mistakes:        0,
	}

	gm.userData.UpsertPuzzleRecord(record)
	return gm.userData.Save("user_data.json")
}

func (gm *GameManager) LoadGame(id string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Find record
	// Find record in Pending or Completed
	var record *PuzzleRecord
	gm.userData.mu.RLock()

	// Check Pending first (most likely for Resume)
	for i := range gm.userData.PendingHistory {
		if gm.userData.PendingHistory[i].ID == id {
			record = &gm.userData.PendingHistory[i]
			break
		}
	}

	// Check Completed if not found
	if record == nil {
		for i := range gm.userData.CompletedHistory {
			if gm.userData.CompletedHistory[i].ID == id {
				record = &gm.userData.CompletedHistory[i]
				break
			}
		}
	}
	gm.userData.mu.RUnlock()

	if record == nil {
		return fmt.Errorf("record not found: %s", id)
	}

	initial, err := engine.ParseBoard(record.Predefined)
	if err != nil {
		return err
	}
	current, err := engine.ParseBoard(record.FinalState)
	if err != nil {
		return err
	}

	// Restore Given
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			val, _ := initial.GetValue(r, c)
			if val != 0 {
				initial.Cells[r][c].Given = true
				current.Cells[r][c].Given = true
			}
		}
	}

	// Create Session
	session := &GameSession{
		initialBoard:    initial,
		currentBoard:    current,
		difficulty:      record.Difficulty,
		difficultyIndex: record.DifficultyIndex,
		gameID:          record.ID,
		history:         NewHistoryManager(),
		startTime:       time.Now().Add(-record.TimeElapsed),
		eraseCount:      0, // Note: We don't verify these on reload currently, could be added to record later
		undoCount:       0,
	}

	if record.IsSolved {
		session.endTime = time.Now()
	}

	gm.currentSession = session
	return nil
}

func (gm *GameManager) GetLastGameID() string {
	if gm.userData == nil {
		return ""
	}
	gm.userData.mu.RLock()
	defer gm.userData.mu.RUnlock()

	// Check Pending History for the resume candidate
	if len(gm.userData.PendingHistory) > 0 {
		// Return the most recent one (last in list)
		return gm.userData.PendingHistory[len(gm.userData.PendingHistory)-1].ID
	}

	return ""
}
