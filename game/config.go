package game

import "github.com/steven004/goldenfox-sudoku/engine"

// DifficultyConfig holds configuration for each difficulty level
type DifficultyConfig struct {
	ExtraClues int
}

// GameConfig holds the global game configuration
var GameConfig = map[engine.DifficultyLevel]DifficultyConfig{
	engine.Beginner: {ExtraClues: 15},
	engine.Easy:     {ExtraClues: 10},
	engine.Medium:   {ExtraClues: 6},
	engine.Hard:     {ExtraClues: 3},
	engine.Expert:   {ExtraClues: 2},
	engine.FoxGod:   {ExtraClues: 0},
}

// GetExtraClues returns the number of extra clues for a given difficulty
func GetExtraClues(difficulty engine.DifficultyLevel) int {
	if config, ok := GameConfig[difficulty]; ok {
		return config.ExtraClues
	}
	return 0 // Default to 0 if unknown
}

// CalculateDynamicClues adjusts the extra clues based on user progress
func CalculateDynamicClues(difficulty engine.DifficultyLevel, progress int) int {
	baseClues := GetExtraClues(difficulty)

	// Fox God: Always 0
	if difficulty == engine.FoxGod {
		return 0
	}

	// Beginner: No penalty for loss (progress < 0 treated as 0)
	if difficulty == engine.Beginner && progress < 0 {
		progress = 0
	}

	var finalClues int

	// Logic Split
	if difficulty <= engine.Medium {
		// Group 1: Beginner, Easy, Medium
		// Formula: Base - Progress
		finalClues = baseClues - progress
	} else {
		// Group 2: Hard, Expert
		// Formula: Base - (Progress / 2)
		// Note: Integer division handles rounding towards zero
		finalClues = baseClues - (progress / 2)
	}

	// Safety Bounds
	if finalClues < 0 {
		finalClues = 0
	}
	// Optional: Cap max clues to avoid it becoming too easy?
	// For now, we trust the base + negative progress won't be absurd.

	return finalClues
}
