package main

import (
	"context"
	"fmt"
	"os"
)

import "golang.org/x/oauth2"
import "github.com/google/go-github/v39/github"

func main() {
	ctx := context.Background()
	auth := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	httpClient := oauth2.NewClient(ctx, auth)
	v3Client := github.NewClient(httpClient)

	var repos []*github.Repository
	var err error

	// if no CLI args then user
	if len(os.Args[1:]) == 0 {
		// TODO: get current actor and use it instead of empty string
		repos, _, err = v3Client.Repositories.List(ctx, "mfinelli", nil)

		if err != nil {
			fmt.Println(err)
		}
	// else assume org
	} else if len(os.Args[1:]) == 1 {
		repos, _, err = v3Client.Repositories.ListByOrg(ctx, os.Args[1], nil)

		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("usage...")
	}

	// fmt.Println(repos)

	for _, repo := range repos {
		fmt.Println(github.Stringify(repo.FullName))
	}
}
