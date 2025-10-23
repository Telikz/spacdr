package cmd

import (
	"cram/internal/app"

	"github.com/spf13/cobra"
)

var deckPath string

var RootCmd = &cobra.Command{
	Use:   "cram",
	Short: "A flashcard CLI application",
	Long:  "cram is a command-line flashcard application for studying",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.StartStudySession(deckPath)
	},
}

func init() {
	RootCmd.Flags().StringVar(&deckPath, "deck", "deck.json", "path to deck file")
}
