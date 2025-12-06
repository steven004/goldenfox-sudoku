# Golden Fox Sudoku 🦊

A premium, modern Sudoku experience built with Go (Wails) and React.

![Golden Fox Sudoku](Design/fox_sudoku_final.png)

## Features

- **✨ Premium Design:** Sleek dark mode interface with "Golden Fox" aesthetics.
- **🎮 5000+ Puzzles:** Curated dataset ranging from Beginner to Expert.
- **🏆 Leveling System:** Progress from Level 1 to Level 6 (Fox God).
- **⚔️ Challenger Mode:** (New in V2.0) Choose your difficulty manually and earn bonus rank progress for taking on tougher challenges!
- **📝 Smart Tools:** Pencil marks, auto-erase, undo/redo, and conflict highlighting.
- **💾 Auto-Save:** Never lose your progress; games save automatically.
- **📊 Rich Statistics:** Track your win rate, see difficulty-specific stats, and monitor your rank progress.
- **⚡ Optimized Performance:** Silky smooth 60fps UI with React component optimization.

## Installation

### macOS
1.  Download the latest `.dmg` from the [Releases](https://github.com/steven004/goldenfox-sudoku/releases) page.
2.  Open the disk image and drag the app to your **Applications** folder.

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Running Locally
```bash
# Install dependencies
cd frontend && npm install && cd ..

# Run in development mode
wails dev
```

### Building
```bash
# Build app bundle
wails build

# Create DMG installer (macOS)
brew install create-dmg
create-dmg --volname "GoldenFox Sudoku" ... (see workflow)
```

## License

MIT License
