package service

import (
	"path/filepath"
	"github.com/telikz/spacdr/internal/domain"
	"github.com/telikz/spacdr/internal/repo"
	"testing"
	"time"
)

func createTestDeck() *domain.Deck {
	return &domain.Deck{
		Name: "Test Deck",
		Cards: []domain.Card{
			{Front: "Q1", Back: "A1", Score: 0, LastReview: time.Time{}},
			{Front: "Q2", Back: "A2", Score: 2, LastReview: time.Now().Add(-24 * time.Hour)},
			{Front: "Q3", Back: "A3", Score: 1, LastReview: time.Now().Add(-48 * time.Hour)},
		},
	}
}

func TestDeckServiceRateCard(t *testing.T) {
	svc := NewDeckService(repo.NewFileDeckRepository())
	deck := createTestDeck()

	err := svc.RateCard(deck, 0, 5)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if deck.Cards[0].Score != 5 {
		t.Errorf("Expected score 5, got %d", deck.Cards[0].Score)
	}
	if deck.Cards[0].LastReview.IsZero() {
		t.Error("Expected LastReview to be updated, got zero time")
	}
}

func TestDeckServiceRateCardInvalidIndex(t *testing.T) {
	svc := NewDeckService(repo.NewFileDeckRepository())
	deck := createTestDeck()

	err := svc.RateCard(deck, 10, 5)
	if err != nil {
		t.Errorf("Expected nil error for out of bounds, got %v", err)
	}

	err = svc.RateCard(deck, -1, 5)
	if err != nil {
		t.Errorf("Expected nil error for negative index, got %v", err)
	}
}

func TestDeckServiceSortCardsByScore(t *testing.T) {
	svc := NewDeckService(repo.NewFileDeckRepository())
	deck := createTestDeck()

	svc.SortCardsByScore(deck)

	if deck.Cards[0].Score != 0 {
		t.Errorf("Expected first card to have score 0, got %d", deck.Cards[0].Score)
	}
	if deck.Cards[1].Score != 1 {
		t.Errorf("Expected second card to have score 1, got %d", deck.Cards[1].Score)
	}
	if deck.Cards[2].Score != 2 {
		t.Errorf("Expected third card to have score 2, got %d", deck.Cards[2].Score)
	}
}

func TestDeckServiceSortCardsByScoreThenByDate(t *testing.T) {
	svc := NewDeckService(repo.NewFileDeckRepository())
	deck := &domain.Deck{
		Name: "Test Deck",
		Cards: []domain.Card{
			{Front: "Q1", Back: "A1", Score: 1, LastReview: time.Now()},
			{Front: "Q2", Back: "A2", Score: 1, LastReview: time.Now().Add(-24 * time.Hour)},
			{Front: "Q3", Back: "A3", Score: 2, LastReview: time.Now()},
		},
	}

	svc.SortCardsByScore(deck)

	if deck.Cards[0].Front != "Q2" {
		t.Errorf("Expected Q2 first (older review), got %s", deck.Cards[0].Front)
	}
	if deck.Cards[1].Front != "Q1" {
		t.Errorf("Expected Q1 second, got %s", deck.Cards[1].Front)
	}
	if deck.Cards[2].Front != "Q3" {
		t.Errorf("Expected Q3 last (highest score), got %s", deck.Cards[2].Front)
	}
}

func TestDeckServiceNextCard(t *testing.T) {
	svc := NewDeckService(repo.NewFileDeckRepository())
	deck := createTestDeck()

	next := svc.NextCard(deck, 0)
	if next != 1 {
		t.Errorf("Expected next to be 1, got %d", next)
	}

	next = svc.NextCard(deck, 2)
	if next != 2 {
		t.Errorf("Expected next to be 2 (at end), got %d", next)
	}
}

func TestDeckServicePreviousCard(t *testing.T) {
	svc := NewDeckService(repo.NewFileDeckRepository())

	prev := svc.PreviousCard(1)
	if prev != 0 {
		t.Errorf("Expected previous to be 0, got %d", prev)
	}

	prev = svc.PreviousCard(0)
	if prev != 0 {
		t.Errorf("Expected previous to be 0 (at start), got %d", prev)
	}
}

func TestDeckServiceLoadAndSaveDeck(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "service_test.json")

	repo := repo.NewFileDeckRepository()
	svc := NewDeckService(repo)

	deck := createTestDeck()
	err := svc.SaveDeck(filePath, deck)
	if err != nil {
		t.Fatalf("Failed to save deck: %v", err)
	}

	loaded, err := svc.LoadDeck(filePath)
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

func TestDeckServiceRateMultipleCards(t *testing.T) {
	svc := NewDeckService(repo.NewFileDeckRepository())
	deck := createTestDeck()

	err := svc.RateCard(deck, 0, 5)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	err = svc.RateCard(deck, 1, 3)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	err = svc.RateCard(deck, 2, 4)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	if deck.Cards[0].Score != 5 {
		t.Errorf("Card 0: expected 5, got %d", deck.Cards[0].Score)
	}
	if deck.Cards[1].Score != 3 {
		t.Errorf("Card 1: expected 3, got %d", deck.Cards[1].Score)
	}
	if deck.Cards[2].Score != 4 {
		t.Errorf("Card 2: expected 4, got %d", deck.Cards[2].Score)
	}
}
