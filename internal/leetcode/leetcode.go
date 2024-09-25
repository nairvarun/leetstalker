package leetcode

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)


type GraphQLRequest struct {
	Query		string					`json:"query"`
	Variables	map[string]string		`json:"variables"`
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
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	return body, nil
}
