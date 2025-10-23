package repo

import (
	"encoding/json"
	"os"
	"github.com/Telikz/spacdr/internal/domain"
)

type DeckRepository interface {
	Load(filePath string) (*domain.Deck, error)
	Save(filePath string, deck *domain.Deck) error
}

type FileDeckRepository struct{}

func NewFileDeckRepository() DeckRepository {
	return &FileDeckRepository{}
}

func (r *FileDeckRepository) Load(filePath string) (*domain.Deck, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var deck domain.Deck
	err = json.Unmarshal(data, &deck)
	if err != nil {
		return nil, err
	}

	return &deck, nil
}

func (r *FileDeckRepository) Save(filePath string, deck *domain.Deck) error {
	data, err := json.MarshalIndent(deck, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
