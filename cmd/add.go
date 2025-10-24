package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/telikz/spacdr/internal/config"
)

var addCategory string

var AddCmd = &cobra.Command{
	Use:   "add <deck-file>",
	Short: "Add a flashcard deck to the .spacdr directory",
	Long:  "Add a flashcard deck JSON file to the .spacdr directory, optionally organized by category",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourcePath := args[0]

		if _, err := os.Stat(sourcePath); err != nil {
			return fmt.Errorf("deck file not found: %s", sourcePath)
		}

		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			return fmt.Errorf("error opening deck file: %w", err)
		}
		defer sourceFile.Close()

		fileName := filepath.Base(sourcePath)
		var destPath string

		spacdrDir := config.GetSpacdrDir()
		if addCategory != "" {
			categoryDir := filepath.Join(spacdrDir, addCategory)
			if err := os.MkdirAll(categoryDir, 0755); err != nil {
				return fmt.Errorf("error creating category directory: %w", err)
			}
			destPath = filepath.Join(categoryDir, fileName)
		} else {
			destPath = filepath.Join(spacdrDir, fileName)
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("error creating destination file: %w", err)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, sourceFile); err != nil {
			return fmt.Errorf("error copying deck file: %w", err)
		}

		relativePath := filepath.Base(destPath)
		if addCategory != "" {
			relativePath = filepath.Join(addCategory, relativePath)
		}

		fmt.Printf("✓ Deck added to %s\n", relativePath)
		return nil
	},
}

func init() {
	AddCmd.Flags().StringVar(&addCategory, "category", "", "category to organize the deck (optional)")
	RootCmd.AddCommand(AddCmd)
}
