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

	// if no CLI args then user
	if len(os.Args[1:]) == 0 {
		opt := &github.RepositoryListOptions{
			// TODO: the limit is 100, use that
			ListOptions: github.ListOptions{PerPage: 10},
		}

		for {
			// TODO: get current actor and use it instead of empty string
			rrepos, resp, err := v3Client.Repositories.List(ctx, "mfinelli", opt)

			if err != nil {
				fmt.Println(err)
				// TODO: exit
			}

			repos = append(repos, rrepos...)

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	// else assume org
	} else if len(os.Args[1:]) == 1 {
		opt := &github.RepositoryListByOrgOptions{
			// TODO: the limit is 100, use that
			ListOptions: github.ListOptions{PerPage: 5},
		}

		for {
			rrepos, resp, err := v3Client.Repositories.ListByOrg(ctx, os.Args[1], opt)

			if err != nil {
				fmt.Println(err)
				// TODO: exit
			}

			repos = append(repos, rrepos...)

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	} else {
		fmt.Println("usage...")
	}

	// fmt.Println(repos)

	for _, repo := range repos {
		fmt.Println(github.Stringify(repo.FullName))
	}
}
