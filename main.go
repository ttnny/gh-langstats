package gh_langstats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
)

func GetLangStats(username string) {
	langStats := getLangStats(username)

	// Prepare JSON to respond to the request
	jsonStats, err := json.Marshal(langStats)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonStats))
}

// Get GitHub language statistics
func getLangStats(username string) map[string]int {
	ctx := context.Background()

	// Get GitHub token from "token.txt"
	token, err := ioutil.ReadFile("token.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Uncomment these to create a GitHub authenticated client with your token
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: string(token)}, )
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Create a GitHub non-authenticated client
	// client := github.NewClient(nil)

	// Get a list of 100 most recent pushed/updated repos from GitHub account
	listOptions := github.ListOptions{Page: 1, PerPage: 100}
	opt := &github.RepositoryListOptions{ListOptions: listOptions, Sort: "pushed"}
	repos, _, err := client.Repositories.List(ctx, username, opt)

	// Address API rate limit and other errors
	if err != nil {
		fmt.Printf("Error(s): %v\n", err)
	}

	// Convert the list of repos to type string slice
	var list []string
	for _, repo := range repos {
		list = append(list, *repo.Name)
	}

	// Get a sum of languages in all repos
	langStats := make(map[string]int)
	for _, repo := range list {
		lang, _, err := client.Repositories.ListLanguages(ctx, username, repo)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range lang {
			if value, found := langStats[k]; found { // if exists, add up value
				langStats[k] = v + value
			} else {
				langStats[k] = v
			}
		}
	}

	return langStats
}