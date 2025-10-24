package app

import (
	"fmt"

	"github.com/telikz/spacdr/internal/config"
	"github.com/telikz/spacdr/internal/repo"
	"github.com/telikz/spacdr/internal/service"

	tea "github.com/charmbracelet/bubbletea"
)

func StartStudySession(deckPath string) error {
	repo := repo.NewFileDeckRepository()
	svc := service.NewDeckService(repo)

	for {
		if deckPath == "" {
			selectedPath, err := selectDeckInteractively()
			if err != nil {
				return err
			}
			if selectedPath == "" {
				return nil
			}
			deckPath = selectedPath
		}

		fullPath := config.GetDeckPath(deckPath)

		deck, err := svc.LoadDeck(fullPath)
		if err != nil {
			return fmt.Errorf("error loading deck from %s: %w", fullPath, err)
		}

		svc.SortCardsByScore(deck)

		uiModel := NewUIModel(deck, fullPath, svc)
		p := tea.NewProgram(uiModel, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}

		if !uiModel.goBack {
			break
		}

		deckPath = ""
	}
	return nil
}

func selectDeckInteractively() (string, error) {
	categoryDecks, err := config.DiscoverDecks()
	if err != nil {
		return "", err
	}

	if len(categoryDecks) == 0 {
		return "", fmt.Errorf("no decks found in .spacdr/")
	}

	selector := NewDeckSelectorModel(categoryDecks)
	p := tea.NewProgram(selector, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return "", err
	}

	selectedDeck := selector.GetSelectedDeck()
	if selectedDeck == "" {
		return "", nil
	}

	return selectedDeck, nil
}
