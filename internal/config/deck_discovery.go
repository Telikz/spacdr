package config

import (
	"os"
	"path/filepath"
	"strings"
)

type DeckInfo struct {
	Name         string
	Category     string
	FullPath     string
	RelativePath string
}

type CategoryDecks struct {
	Category string
	Decks    []DeckInfo
}

func DiscoverDecks() ([]CategoryDecks, error) {
	categoryMap := make(map[string][]DeckInfo)

	err := filepath.Walk(spacdrDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".json") && path != spacdrDir {
			rel, err := filepath.Rel(spacdrDir, path)
			if err != nil {
				return err
			}

			parts := strings.Split(rel, string(filepath.Separator))
			var category string
			var deckName string

			if len(parts) == 1 {
				category = ""
				deckName = strings.TrimSuffix(parts[0], ".json")
			} else {
				category = parts[0]
				deckName = strings.TrimSuffix(parts[len(parts)-1], ".json")
			}

			deckRef := strings.TrimSuffix(rel, ".json")

			deckInfo := DeckInfo{
				Name:         deckName,
				Category:     category,
				FullPath:     path,
				RelativePath: deckRef,
			}

			categoryMap[category] = append(categoryMap[category], deckInfo)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var result []CategoryDecks
	if decks, ok := categoryMap[""]; ok {
		result = append(result, CategoryDecks{
			Category: "Uncategorized",
			Decks:    decks,
		})
	}

	for category := range categoryMap {
		if category != "" {
			result = append(result, CategoryDecks{
				Category: category,
				Decks:    categoryMap[category],
			})
		}
	}

	return result, nil
}
