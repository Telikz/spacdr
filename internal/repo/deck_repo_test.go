package repo

import (
	"os"
	"path/filepath"
	"spacdr/internal/domain"
	"testing"
	"time"
)

func TestFileDeckRepositoryLoad(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_deck.json")

	deck := &domain.Deck{
		Name: "Test Deck",
		Cards: []domain.Card{
			{Front: "Q1", Back: "A1", Score: 1, LastReview: time.Now()},
			{Front: "Q2", Back: "A2", Score: 3, LastReview: time.Now()},
		},
	}

	repo := NewFileDeckRepository()
	err := repo.Save(filePath, deck)
	if err != nil {
		t.Fatalf("Failed to save deck: %v", err)
	}

	loaded, err := repo.Load(filePath)
	if err != nil {
		t.Fatalf("Failed to load deck: %v", err)
	}

	if loaded.Name != deck.Name {
		t.Errorf("Expected name %s, got %s", deck.Name, loaded.Name)
	}
	if len(loaded.Cards) != len(deck.Cards) {
		t.Errorf("Expected %d cards, got %d", len(deck.Cards), len(loaded.Cards))
	}
}

func TestFileDeckRepositorySave(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_deck.json")

	deck := &domain.Deck{
		Name: "New Deck",
		Cards: []domain.Card{
			{Front: "Question", Back: "Answer", Score: 5, LastReview: time.Now()},
		},
	}

	repo := NewFileDeckRepository()
	err := repo.Save(filePath, deck)
	if err != nil {
		t.Fatalf("Failed to save deck: %v", err)
	}

	_, err = os.Stat(filePath)
	if err != nil {
		t.Fatalf("File was not created: %v", err)
	}
}

func TestFileDeckRepositoryLoadNotFound(t *testing.T) {
	repo := NewFileDeckRepository()
	_, err := repo.Load("/nonexistent/path/deck.json")
	if err == nil {
		t.Fatal("Expected error for nonexistent file, got nil")
	}
}

func TestFileDeckRepositoryRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "roundtrip.json")

	original := &domain.Deck{
		Name: "Round Trip Deck",
		Cards: []domain.Card{
			{Front: "Q1", Back: "A1", Score: 2, LastReview: time.Now()},
			{Front: "Q2", Back: "A2", Score: 4, LastReview: time.Now()},
			{Front: "Q3", Back: "A3", Score: 1, LastReview: time.Now()},
		},
	}

	repo := NewFileDeckRepository()
	err := repo.Save(filePath, original)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	loaded, err := repo.Load(filePath)
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if loaded.Name != original.Name {
		t.Errorf("Name mismatch: %s vs %s", loaded.Name, original.Name)
	}
	if len(loaded.Cards) != len(original.Cards) {
		t.Errorf("Card count mismatch: %d vs %d", len(loaded.Cards), len(original.Cards))
	}

	for i, card := range loaded.Cards {
		if card.Front != original.Cards[i].Front || card.Back != original.Cards[i].Back {
			t.Errorf("Card %d content mismatch", i)
		}
		if card.Score != original.Cards[i].Score {
			t.Errorf("Card %d score mismatch: %d vs %d", i, card.Score, original.Cards[i].Score)
		}
	}
}
