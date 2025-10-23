package app

import (
	"cram/internal/domain"
	"cram/internal/service"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UIModel struct {
	deck     *domain.Deck
	current  int
	flipped  bool
	quitting bool
	filePath string
	width    int
	height   int
	err      string
	svc      service.DeckService
}

func NewUIModel(deck *domain.Deck, filePath string, svc service.DeckService) *UIModel {
	return &UIModel{
		deck:     deck,
		current:  0,
		flipped:  false,
		quitting: false,
		filePath: filePath,
		svc:      svc,
	}
}

func (m *UIModel) Init() tea.Cmd {
	return nil
}

func (m *UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "h", "l":
			m.flipped = !m.flipped
		case "j":
			m.current = m.svc.NextCard(m.deck, m.current)
			m.flipped = false
		case "k":
			m.current = m.svc.PreviousCard(m.current)
			m.flipped = false

		case "1", "2", "3", "4", "5":
			score := int(msg.String()[0] - '0')
			err := m.svc.RateCard(m.deck, m.current, score)
			if err != nil {
				return nil, nil
			}
			err = m.svc.SaveDeck(m.filePath, m.deck)
			if err != nil {
				return nil, nil
			}
			m.current = m.svc.NextCard(m.deck, m.current)
			m.flipped = false
		}
	}
	return m, nil
}

func (m *UIModel) View() string {
	if len(m.deck.Cards) == 0 {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
		return errorStyle.Render("No cards in deck")
	}

	card := m.deck.Cards[m.current]
	content := card.Front
	if m.flipped {
		content = card.Back
	}

	availableWidth := m.width
	if availableWidth < 20 {
		availableWidth = 20
	}

	cardPadding := availableWidth / 4
	cardWidth := availableWidth - cardPadding

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Width(cardWidth).
		Align(lipgloss.Center)

	contentStyle := lipgloss.NewStyle()

	progress := fmt.Sprintf("(%d/%d)", m.current+1, len(m.deck.Cards))
	scoreStr := ""
	if card.Score > 0 {
		scoreStyle := lipgloss.NewStyle()
		scoreStr = " " + scoreStyle.Render(fmt.Sprintf(" - %d/5", card.Score))
	}

	header := headerStyle.Render(fmt.Sprintf("%s  %s%s", m.deck.Name, progress, scoreStr))
	cardBox := cardStyle.Render(contentStyle.Render(strings.TrimSpace(content)))

	help := "Rate: [1] [2] [3] [4] [5] \n" + "[H/L] flip  |  [J/K] navigate |  [Q] quit"

	helpStyle := lipgloss.NewStyle().
		Italic(true).
		Align(lipgloss.Center).
		Width(m.width)

	helpText := helpStyle.Render(help)
	helpHeight := lipgloss.Height(helpText)
	verticalSpace := m.height - lipgloss.Height(cardBox) - helpHeight - 4
	if verticalSpace < 0 {
		verticalSpace = 0
	}

	topSpacer := strings.Repeat("\n", verticalSpace/2)
	bottomSpacer := strings.Repeat("\n", verticalSpace/2)

	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)

	centeredCard := containerStyle.Render(cardBox)

	middle := lipgloss.JoinVertical(lipgloss.Center,
		topSpacer,
		centeredCard,
		bottomSpacer,
	)

	mainContent := lipgloss.JoinVertical(lipgloss.Top,
		header,
		middle,
	)

	helpSpacing := lipgloss.NewStyle().
		Height(m.height - lipgloss.Height(mainContent) - helpHeight)

	return lipgloss.JoinVertical(lipgloss.Top,
		mainContent,
		helpSpacing.Render(""),
		helpText,
	)
}
