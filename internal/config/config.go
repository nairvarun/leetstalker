package config

import (
	"os"
	"path"
	"strings"
)

type config struct {
	Usernames []string
}

func createConfigFile(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			return
		}
	}()

	return nil
}

func LoadConfig() (c config, err error) {
	homeDir, err := os.UserHomeDir()
    if err != nil {
        return
    }

	configDir := path.Join(homeDir, ".config", "leetstalker")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return
	}

	configFile := path.Join(configDir, "config.txt")

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		err = createConfigFile(configFile)
		if err != nil {
			return
		}
	} else if err != nil {
		return
	}

	if _, err = os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			err = createConfigFile(configFile)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	configData, err := os.ReadFile(configFile)
	if err != nil {
		return
	}

	usernames := strings.Split(strings.Trim(string(configData), "\n"), "\n")
	c.Usernames = usernames
	return c, nil
}
