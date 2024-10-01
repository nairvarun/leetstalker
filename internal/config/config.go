package config

import (
	"embed"
	"errors"

	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// load config
// create config
// add user to config

//go:embed config.yaml
var f embed.FS

type Configuration struct {
	Users []string `json:"users"`
}

func getConfigPath() (string, error) {
	if homeDir, err := os.UserHomeDir(); err != nil {
		return "", err
	} else {
		return filepath.Join(homeDir, ".config", "leetstalker", "config.yaml"), nil
	}
}

func createConfigFile() error {
	// create config directory
	config, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(config), 0755); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}

	// create config file
	destination, err := os.Create(filepath.Base(config))
	if err != nil {
		return err
	}; defer destination.Close()

	// open embedded source file
	source, err := f.Open("config.yaml")
	if err != nil {
		return err
	}; defer source.Close()

	// copy embedded file to destination
	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return nil
}

func InitConfig(configuration *Configuration) error {
	// get config
	config, err := getConfigPath()
	if err != nil {
		return err
	}

	// set config in viper
	viper.SetConfigFile(config)

	// read config
	if err := viper.ReadInConfig(); err != nil {
		// Create file if not exists
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Creating empty config file (~/.config/leetstalker/config.yaml)")
			if err := createConfigFile(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		return err
	}

	return nil
}