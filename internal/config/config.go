package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	spacdrDir string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	spacdrDir = filepath.Join(home, ".spacdr")
}

func GetSpacdrDir() string {
	return spacdrDir
}

func InitializeConfig() error {
	if _, err := os.Stat(spacdrDir); os.IsNotExist(err) {
		if err := os.MkdirAll(spacdrDir, 0755); err != nil {
			return err
		}
		if err := CreateTutorialDeck(); err != nil {
			return err
		}
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(spacdrDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}

func GetDeckPath(deckRef string) string {
	return filepath.Join(spacdrDir, deckRef+".json")
}
