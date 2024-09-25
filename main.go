package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"leetstalker/internal/config"
	"leetstalker/internal/leetcode"
)

type UserData struct {
	Username	string
	Data		[]byte
	Error		error
}

func worker(username string, results chan<- UserData) {
	data, err := leetcode.FetchUserData(username)
	results <- UserData{
		Username:	username,
		Data:		data,
		Error:		err,
	}
}

func main() {
	config, err := config.LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}

	results := make(chan UserData, len(config.Users))
	var wg sync.WaitGroup

	for _, username := range config.Users {
		wg.Add(1)
		go func (u string) {
			defer wg.Done()
			worker(u, results)
		}(username)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Error != nil {
			log.Fatal(result.Error)
			continue
		}

		var jsonResult interface{}
		if err := json.Unmarshal(result.Data, &jsonResult); err != nil {
			log.Fatal(err)
			continue
		}

		prettyJSON, err := json.MarshalIndent(jsonResult, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON for %s: %v\n", result.Username, err)
			continue
		}

		fmt.Printf("Data for %s:\n%s\n\n", result.Username, string(prettyJSON))
	}
}