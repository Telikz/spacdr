package app

import (
	"fmt"
	"github.com/telikz/spacdr/internal/repo"
	"github.com/telikz/spacdr/internal/service"

	tea "github.com/charmbracelet/bubbletea"
)

func StartStudySession(deckPath string) error {
	repo := repo.NewFileDeckRepository()
	svc := service.NewDeckService(repo)

	deck, err := svc.LoadDeck(deckPath)
	if err != nil {
		return fmt.Errorf("error loading deck from %s: %w", deckPath, err)
	}

	svc.SortCardsByScore(deck)

	uiModel := NewUIModel(deck, deckPath, svc)
	p := tea.NewProgram(uiModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
