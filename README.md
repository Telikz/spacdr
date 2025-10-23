# cram

A terminal-based flashcard application for studying and memorizing. Built with Go and featuring an interactive TUI (Text User Interface) for smooth learning sessions.

## Features

- **Interactive Study Sessions** - Flip through flashcards with keyboard shortcuts
- **Smart Card Ordering** - Cards are automatically sorted by score, prioritizing cards you struggle with
- **Score Tracking** - Rate each card from 1-5 and track review history
- **Spaced Repetition** - Card scores adjust based on review dates to optimize retention
- **Persistent Storage** - Decks are saved as JSON files for easy sharing and version control

## Installation

### Prerequisites

- Go 1.25.3 or higher

### Build

```bash
make build
```

This creates an executable called `cram` in the current directory.


This loads the default `deck.json` from the current directory.

### Keyboard Controls

**During Card View:**
- `h` / `l` - Flip card (show front/back)
- `j` / `k` - Navigate to next/previous card
- `r` - Rate the current card
- `q` / `Ctrl+C` - Quit

**During Rating:**
- `1` - 1/5 (hard)
- `2` - 2/5
- `3` - 3/5
- `4` - 4/5
- `5` - 5/5 (easy)
- `q` / `Ctrl+C` - Quit

## Deck Format

Decks are stored as JSON files with the following structure:

```json
{
  "name": "Spanish Vocabulary",
  "cards": [
    {
      "front": "¿Cómo estás?",
      "back": "How are you?",
      "score": 0,
      "last_review": "0001-01-01T00:00:00Z"
    },
    {
      "front": "Hola",
      "back": "Hello",
      "score": 5,
      "last_review": "2025-10-20T14:30:00Z"
    }
  ]
}
```

### Fields

- `name` - Deck name (displayed in header)
- `cards` - Array of card objects
  - `front` - Question/prompt side of the card
  - `back` - Answer side of the card
  - `score` - Current rating (0-5, starts at 0)
  - `last_review` - ISO 8601 timestamp of last review

## Development


### Available Commands

```bash
make build      # Build the binary
make run        # Build and run the application
make test       # Run unit tests
make lint       # Run golangci-lint
make clean      # Remove binary
make help       # Show available targets
```

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```
## License

MIT
