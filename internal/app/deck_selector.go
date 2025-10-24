package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/telikz/spacdr/internal/config"
)

type DeckItem struct {
	name       string
	category   string
	path       string
	isCategory bool
}

type DeckSelectorModel struct {
	categoryDecks []config.CategoryDecks
	allItems      []DeckItem
	filteredIdx   []int
	selectedIdx   int
	width         int
	height        int
	scrollOffset  int
	searchMode    bool
	searchQuery   string
	confirmed     bool
}

func NewDeckSelectorModel(categoryDecks []config.CategoryDecks) *DeckSelectorModel {
	m := &DeckSelectorModel{
		categoryDecks: categoryDecks,
		selectedIdx:   0,
		scrollOffset:  0,
		confirmed:     false,
	}
	m.buildItems()
	return m
}

func (m *DeckSelectorModel) buildItems() {
	m.allItems = []DeckItem{}

	for _, cd := range m.categoryDecks {
		m.allItems = append(m.allItems, DeckItem{
			name:       cd.Category,
			isCategory: true,
		})

		for _, deck := range cd.Decks {
			m.allItems = append(m.allItems, DeckItem{
				name:     deck.Name,
				category: cd.Category,
				path:     deck.RelativePath,
			})
		}
	}

	m.updateFilter()
}

func (m *DeckSelectorModel) updateFilter() {
	m.filteredIdx = []int{}

	if m.searchQuery == "" {
		for i := range m.allItems {
			m.filteredIdx = append(m.filteredIdx, i)
		}
	} else {
		query := strings.ToLower(m.searchQuery)
		for i, item := range m.allItems {
			if strings.Contains(strings.ToLower(item.name), query) {
				m.filteredIdx = append(m.filteredIdx, i)
			}
		}
	}

	if m.selectedIdx >= len(m.filteredIdx) {
		m.selectedIdx = len(m.filteredIdx) - 1
		if m.selectedIdx < 0 {
			m.selectedIdx = 0
		}
	}

	m.scrollOffset = 0
	m.ensureVisible()
}

func (m *DeckSelectorModel) Init() tea.Cmd {
	return nil
}

func (m *DeckSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.searchMode {
			switch msg.String() {
			case "esc":
				m.searchMode = false
				m.searchQuery = ""
				m.updateFilter()
			case "enter":
				m.searchMode = false
				if m.selectedIdx < len(m.filteredIdx) {
					realIdx := m.filteredIdx[m.selectedIdx]
					if realIdx < len(m.allItems) && !m.allItems[realIdx].isCategory {
						m.confirmed = true
						return m, tea.Quit
					}
				}
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.updateFilter()
				}
			default:
				if len(msg.String()) == 1 && msg.String()[0] >= 32 && msg.String()[0] < 127 {
					m.searchQuery += msg.String()
					m.updateFilter()
				}
			}
		} else {
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "/":
				m.searchMode = true
				m.searchQuery = ""
				m.updateFilter()
			case "l", "enter":
				if m.selectedIdx < len(m.filteredIdx) {
					realIdx := m.filteredIdx[m.selectedIdx]
					if realIdx < len(m.allItems) && !m.allItems[realIdx].isCategory {
						m.confirmed = true
						return m, tea.Quit
					}
				}
			case "j", "down":
				m.selectedIdx++
				if m.selectedIdx >= len(m.filteredIdx) {
					m.selectedIdx = len(m.filteredIdx) - 1
				}
				m.ensureVisible()
			case "k", "up":
				m.selectedIdx--
				if m.selectedIdx < 0 {
					m.selectedIdx = 0
				}
				m.ensureVisible()
			}
		}
	}
	return m, nil
}

func (m *DeckSelectorModel) ensureVisible() {
	viewportHeight := m.height - 10
	if viewportHeight < 2 {
		viewportHeight = 2
	}

	filteredCount := len(m.filteredIdx)
	if viewportHeight > filteredCount {
		viewportHeight = filteredCount
	}

	if m.selectedIdx < m.scrollOffset {
		m.scrollOffset = m.selectedIdx
	}
	if m.selectedIdx >= m.scrollOffset+viewportHeight {
		m.scrollOffset = m.selectedIdx - viewportHeight + 1
	}
}

func (m *DeckSelectorModel) View() string {
	if len(m.allItems) == 0 {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
		return errorStyle.Render("✗ No decks found in ~/.spacdr/")
	}

	availableWidth := m.width
	if availableWidth < 20 {
		availableWidth = 20
	}

	listPadding := availableWidth / 4
	listWidth := availableWidth - listPadding

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width)

	title := titleStyle.Render("Select a Deck")

	var listContent strings.Builder

	viewportHeight := m.height - 10
	if viewportHeight < 2 {
		viewportHeight = 2
	}

	if len(m.filteredIdx) == 0 {
		noResultsStyle := lipgloss.NewStyle().
			Italic(true)
		listContent.WriteString(noResultsStyle.Render("  No results"))
	} else {
		endIdx := m.scrollOffset + viewportHeight
		if endIdx > len(m.filteredIdx) {
			endIdx = len(m.filteredIdx)
		}

		for i := m.scrollOffset; i < endIdx; i++ {
			realIdx := m.filteredIdx[i]
			item := m.allItems[realIdx]
			isSelected := i == m.selectedIdx

			if item.isCategory {
				catStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("33")).
					Bold(true)
				if isSelected {
					catStyle = catStyle.
						Foreground(lipgloss.Color("220"))
				}
				listContent.WriteString(catStyle.Render(item.name) + "\n")
			} else {
				deckStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("250"))

				prefix := "  "
				if isSelected {
					deckStyle = deckStyle.
						Bold(true)
					prefix = "▶ "
				}

				deckName := item.name
				maxLen := listWidth - len(prefix) - 4
				if len(deckName) > maxLen && maxLen > 3 {
					deckName = deckName[:maxLen-3] + "..."
				}

				listContent.WriteString(deckStyle.Render(prefix+deckName) + "\n")
			}
		}
	}

	listStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 4).
		Width(listWidth)

	listBox := listStyle.Render(listContent.String())

	containerStyle := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)

	centeredList := containerStyle.Render(listBox)

	helpStyle := lipgloss.NewStyle().
		Italic(true).
		Align(lipgloss.Center).
		Width(m.width)

	var help string
	if m.searchMode {
		help = helpStyle.Render("Type to search • Esc to exit search")
	} else {
		help = helpStyle.Render("↑↓ JK | L select | / search  | Q quit")
	}

	helpHeight := lipgloss.Height(help)
	listBoxHeight := lipgloss.Height(listBox)
	verticalSpace := m.height - listBoxHeight - helpHeight - 4
	if verticalSpace < 0 {
		verticalSpace = 0
	}

	topSpacer := strings.Repeat("\n", verticalSpace/2)
	bottomSpacer := strings.Repeat("\n", verticalSpace/2)

	middle := lipgloss.JoinVertical(lipgloss.Center,
		topSpacer,
		centeredList,
		bottomSpacer,
	)

	headerParts := []string{title}
	if m.searchMode {
		searchStyle := lipgloss.NewStyle().
			Bold(true).
			Align(lipgloss.Center).
			Width(listWidth)
		cursor := "│"
		searchBar := searchStyle.Render("Search: " + m.searchQuery + cursor)
		headerParts = append(headerParts, "", searchBar)
	}

	mainContent := lipgloss.JoinVertical(lipgloss.Top, append(headerParts, middle)...)

	helpSpacing := lipgloss.NewStyle().
		Height(m.height - lipgloss.Height(mainContent) - helpHeight)

	return lipgloss.JoinVertical(lipgloss.Top,
		mainContent,
		helpSpacing.Render(""),
		help,
	)
}

func (m *DeckSelectorModel) GetSelectedDeck() string {
	if !m.confirmed {
		return ""
	}
	if m.selectedIdx < len(m.filteredIdx) {
		realIdx := m.filteredIdx[m.selectedIdx]
		if realIdx < len(m.allItems) && !m.allItems[realIdx].isCategory {
			return m.allItems[realIdx].path
		}
	}
	return ""
}
