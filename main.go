package main

// todo:
// github actions to build and release

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/nairvarun/leetstalker/internal/color"
	"github.com/nairvarun/leetstalker/internal/config"
)

// {"data":{"matchedUser":{"profile":{"ranking":198312},"submitStatsGlobal":{"acSubmissionNum":[{"difficulty":"All","count":324},{"difficulty":"Easy","count":133},{"difficulty":"Medium","count":168},{"difficulty":"Hard","count":23}]}}}}
type UserData struct {
	Data struct {
		MatchedUser struct {
			Profile struct {
				Ranking int
			}
			SubmitStatsGlobal struct {
				AcSubmissionNum []struct {
					Difficulty string
					Count int
				}
			}
		}
	}
}

type chanResult struct {
	username string
	userData UserData
}

func query(username string, ch chan<-chanResult, wg *sync.WaitGroup) {
	defer wg.Done()
	graphqlQuery := map[string]string {
		"query": `
			query userPublicProfileAndProblemsSolved($username: String!) {
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
		`,
		"variables": fmt.Sprintf(`{
            "username": "%s"
        }`, username),
	}

	graphqlQueryJSON, err := json.Marshal(graphqlQuery)
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", "https://leetcode.com/graphql/", bytes.NewBuffer(graphqlQueryJSON))
	if err != nil {
		panic(err)
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	userData := new(UserData)
	if err = json.NewDecoder(response.Body).Decode(&userData); err != nil {
		panic(err)
	}

	ch <- chanResult{
		username: username,
		userData: (*userData),
	}
}

func main() {
	ch := make(chan chanResult)
	var wg sync.WaitGroup

	configuration, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	unames := configuration.Usernames
	for _, uname := range unames {
		wg.Add(1)
		go query(uname, ch, &wg)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	results := make([]chanResult, 0)
	for res := range ch {
		results = append(results, res)
	}
	rankCmp := func (a, b chanResult) int {
		return cmp.Compare(a.userData.Data.MatchedUser.Profile.Ranking, b.userData.Data.MatchedUser.Profile.Ranking)
	}
	slices.SortFunc(results, rankCmp)

	for _, result := range results {
		fmt.Printf("%s (%s) [%s: %s + %s + %s]\n",
			color.Format(color.BOLD,	result.username),
			color.Format(color.ITALIC, 	strconv.Itoa(result.userData.Data.MatchedUser.Profile.Ranking)),
			color.Format(color.NONE, 	strconv.Itoa(result.userData.Data.MatchedUser.SubmitStatsGlobal.AcSubmissionNum[0].Count)),
			color.Format(color.GREEN,	strconv.Itoa(result.userData.Data.MatchedUser.SubmitStatsGlobal.AcSubmissionNum[1].Count)),
			color.Format(color.YELLOW, 	strconv.Itoa(result.userData.Data.MatchedUser.SubmitStatsGlobal.AcSubmissionNum[2].Count)),
			color.Format(color.RED, 	strconv.Itoa(result.userData.Data.MatchedUser.SubmitStatsGlobal.AcSubmissionNum[3].Count)),
		)
	}
}
