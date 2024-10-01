package leetcode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)


type GraphQLRequest struct {
	Query		string					`json:"query"`
	Variables	map[string]string		`json:"variables"`
}

type UserData struct {
	Username	string
	Data		[]byte
	Error		error
}


const (
    url = "https://leetcode.com/graphql/"

    query = `
		query userData($username: String!) {
			userContestRanking(username: $username) {
				rating
				globalRanking
				topPercentage
			}
			matchedUser(username: $username) {
				profile {
					ranking
				}
				submitStatsGlobal {
					acSubmissionNum {
						difficulty
						count
					}
				}
			}
		}
	`
)


func FetchUserData(username string) ([]byte, error) {
	variables := map[string]string{
		"username": username,
	}

	requestData := GraphQLRequest{
		Query: 		query,
		Variables: 	variables,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.133 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}; defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return body, nil
}

func worker(username string, results chan<- UserData) {
	data, err := FetchUserData(username)
	results <- UserData{
		Username:	username,
		Data:		data,
		Error:		err,
	}
}

func FetchMultiple(users []string) {
	var wg sync.WaitGroup
	results := make(chan UserData, len(users))
	for _, username := range users {
		wg.Add(1)
		go func (u string) {
			defer wg.Done()
			worker(u, results)
		}(username)

	}

	go func () {
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