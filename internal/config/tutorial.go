package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/telikz/spacdr/internal/domain"
)

func CreateTutorialDeck() error {
	tutorialDeck := &domain.Deck{
		Name: "Tutorial Deck",
		Cards: []domain.Card{
			{
				Front:      "What is a flashcard?",
				Back:       "A flashcard is a study tool with a question on the front and answer on the back used for learning and memorization.",
				Score:      0,
				LastReview: time.Time{},
			},
			{
				Front:      "How do you rate a card in spacdr?",
				Back:       "Press 'r' to rate a card, then select a score from 1-5 (1=hard, 5=easy).",
				Score:      0,
				LastReview: time.Time{},
			},
			{
				Front:      "What keyboard shortcut flips a card?",
				Back:       "Press 'h' or 'l' to flip between the front and back of a card.",
				Score:      0,
				LastReview: time.Time{},
			},
			{
				Front:      "How do you navigate between cards?",
				Back:       "Press 'j' to go to the next card or 'k' to go to the previous card.",
				Score:      0,
				LastReview: time.Time{},
			},
			{
				Front:      "What does spaced repetition do?",
				Back:       "Spaced repetition optimizes learning by reviewing cards at increasing intervals, focusing on cards you struggle with.",
				Score:      0,
				LastReview: time.Time{},
			},
		},
	}

	tutorialPath := filepath.Join(spacdrDir, "tutorial.json")

	if _, err := os.Stat(tutorialPath); err == nil {
		return nil
	}

	data, err := json.MarshalIndent(tutorialDeck, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tutorialPath, data, 0644)
}
