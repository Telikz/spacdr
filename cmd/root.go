package cmd

import (
	"github.com/telikz/spacdr/internal/app"
	"github.com/telikz/spacdr/internal/config"

	"github.com/spf13/cobra"
)

var deckPath string

var RootCmd = &cobra.Command{
	Use:   "spacdr",
	Short: "A flashcard CLI application",
	Long:  "spacdr is a command-line flashcard application for studying",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.InitializeConfig()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.StartStudySession(deckPath)
	},
}

func init() {
	RootCmd.Flags().StringVar(&deckPath, "deck", "", "path to deck file (relative to .spacdr, e.g. 'spanish/vocabulary' or 'deck.json'). If empty, an interactive menu will be shown")
}
