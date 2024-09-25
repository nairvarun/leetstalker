package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Configuration struct {
	Users []string `json:"users"`
}

func createConfigFile(configDir, configFile string) error {
    // Create config directory
    if err := os.MkdirAll(configDir, 0755); err != nil && !os.IsExist(err) {
        return err
    }

    // Create config file
    if _, err := os.Create(configFile); err != nil {
        return err
    }

	// Write sample data to config file
	viper.Set("users", []string {"larryNY"})
	if err := viper.WriteConfigAs(filepath.Join(configDir, configFile)); err != nil {
		return err
	}

    return nil
}


func LoadConfiguration() (*Configuration, error) {
	var config *Configuration

	// Get user $HOME
	homeDir, err := os.UserHomeDir();
	if err != nil {
		return nil, err
	}

	// Set config file path
	configDir := filepath.Join(homeDir, ".config", "leetstalker")
	configFile := "config"
	configType := "yaml"

	viper.AddConfigPath(configDir)
	viper.SetConfigName(configFile)
	viper.SetConfigType(configType)


	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// Create file if not exists
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found. Creating empty config file at ~/.config/leetstalker/")
			if err := createConfigFile(configDir, fmt.Sprintf("%s.%s", configFile, configType)); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}