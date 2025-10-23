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
	rating   bool
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
		rating:   false,
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
		if m.rating {
			switch msg.String() {
			case "1", "2", "3", "4", "5":
				score := int(msg.String()[0] - '0')
				m.svc.RateCard(m.deck, m.current, score)
				err := m.svc.SaveDeck(m.filePath, m.deck)
				if err != nil {
					return nil, nil
				}
				m.current = m.svc.NextCard(m.deck, m.current)
				m.rating = false
				m.flipped = false
			case "q", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
		} else {
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
			case "r":
				m.rating = true
			}
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
	if availableWidth < 40 {
		availableWidth = 40
	}

	cardPadding := 4
	cardWidth := availableWidth - cardPadding

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		MarginBottom(1).
		Align(lipgloss.Center).
		Width(availableWidth)

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(2).
		MarginTop(1).
		MarginBottom(2).
		Width(cardWidth).
		Align(lipgloss.Center)

	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	progress := fmt.Sprintf("(%d/%d)", m.current+1, len(m.deck.Cards))
	scoreStr := ""
	if card.Score > 0 {
		scoreStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
		scoreStr = " " + scoreStyle.Render(fmt.Sprintf("⭐ %d/5", card.Score))
	}

	header := headerStyle.Render(fmt.Sprintf("%s  %s%s", m.deck.Name, progress, scoreStr))
	cardBox := cardStyle.Render(contentStyle.Render(strings.TrimSpace(content)))

	var help string
	if m.rating {
		help = "Rate: [1] [2] [3] [4] [5] | [Q] quit"
	} else {
		help = "[H/L] flip  |  [J/K] navigate  |  [R] rate  |  [Q] quit"
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Italic(true).
		MarginTop(1).
		Align(lipgloss.Center).
		Width(availableWidth)

	verticalSpace := m.height - lipgloss.Height(header) - lipgloss.Height(cardBox) - lipgloss.Height(helpStyle.Render(help)) - 4
	if verticalSpace < 0 {
		verticalSpace = 0
	}

	spacer := strings.Repeat("\n", verticalSpace/2)

	return lipgloss.JoinVertical(lipgloss.Center,
		spacer,
		header,
		cardBox,
		helpStyle.Render(help),
	)
}
