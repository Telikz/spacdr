package service

import (
	"sort"
	"spacdr/internal/domain"
	"spacdr/internal/repo"
	"time"
)

type DeckService interface {
	LoadDeck(filePath string) (*domain.Deck, error)
	SaveDeck(filePath string, deck *domain.Deck) error
	RateCard(deck *domain.Deck, cardIndex int, score int) error
	SortCardsByScore(deck *domain.Deck)
	NextCard(deck *domain.Deck, current int) int
	PreviousCard(current int) int
	AdjustCardScoresByReviewDate(deck *domain.Deck)
}

type DeckServiceImpl struct {
	repo repo.DeckRepository
}

func NewDeckService(repo repo.DeckRepository) DeckService {
	return &DeckServiceImpl{repo: repo}
}

func (s *DeckServiceImpl) LoadDeck(filePath string) (*domain.Deck, error) {
	return s.repo.Load(filePath)
}

func (s *DeckServiceImpl) SaveDeck(filePath string, deck *domain.Deck) error {
	return s.repo.Save(filePath, deck)
}

func (s *DeckServiceImpl) RateCard(deck *domain.Deck, cardIndex int, score int) error {
	if cardIndex < 0 || cardIndex >= len(deck.Cards) {
		return nil
	}

	deck.Cards[cardIndex].Score = score
	deck.Cards[cardIndex].LastReview = time.Now()
	return nil
}

func (s *DeckServiceImpl) SortCardsByScore(deck *domain.Deck) {
	sort.Slice(deck.Cards, func(i, j int) bool {
		if deck.Cards[i].Score != deck.Cards[j].Score {
			return deck.Cards[i].Score < deck.Cards[j].Score
		}
		return deck.Cards[i].LastReview.Before(deck.Cards[j].LastReview)
	})
}

func (s *DeckServiceImpl) NextCard(deck *domain.Deck, current int) int {
	if current < len(deck.Cards)-1 {
		return current + 1
	}
	return current
}

func (s *DeckServiceImpl) PreviousCard(current int) int {
	if current > 0 {
		return current - 1
	}
	return current
}

func (s *DeckServiceImpl) AdjustCardScoresByReviewDate(deck *domain.Deck) {
	now := time.Now()
	oneWeek := 7 * 24 * time.Hour
	oneMonth := 30 * 24 * time.Hour

	for i := range deck.Cards {
		if deck.Cards[i].LastReview.IsZero() {
			continue
		}

		timeSinceReview := now.Sub(deck.Cards[i].LastReview)

		if timeSinceReview > oneMonth {
			deck.Cards[i].Score = 1
		} else if timeSinceReview > oneWeek {
			if deck.Cards[i].Score > 1 {
				deck.Cards[i].Score--
			}
		}
	}
}
