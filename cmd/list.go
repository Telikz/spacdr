package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/telikz/spacdr/internal/config"
)

var ListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all available decks",
	Long:  "List all available decks in the .spacdr directory organized by category",
	RunE: func(cmd *cobra.Command, args []string) error {
		categoryDecks, err := config.DiscoverDecks()
		if err != nil {
			return err
		}

		if len(categoryDecks) == 0 {
			fmt.Println("No decks found in .spacdr/")
			return nil
		}

		fmt.Println("Available Decks:")
		fmt.Println()

		for _, cd := range categoryDecks {
			fmt.Printf("📂 %s\n", cd.Category)
			for _, deck := range cd.Decks {
				fmt.Printf("   └─ %s (%s)\n", deck.Name, deck.RelativePath)
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(ListCmd)
}
