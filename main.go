package main

import (
	"context"
	"fmt"
	"os"
)

import "golang.org/x/oauth2"
import "github.com/google/go-github/v39/github"
import "github.com/shurcooL/githubv4"

func main() {
	ctx := context.Background()
	auth := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	httpClient := oauth2.NewClient(ctx, auth)
	v3Client := github.NewClient(httpClient)
	v4Client := githubv4.NewClient(httpClient)

	var repos []*github.Repository
	var owner string

	actor, err := getCurrentActor(v4Client)

	if err != nil {
		fmt.Println(err)
		// TODO: exit
	}

	// if no CLI args then user
	if len(os.Args[1:]) == 0 {
		owner = actor

		opt := &github.RepositoryListOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			// TODO: get current actor and use it instead of empty string
			rrepos, resp, err := v3Client.Repositories.List(ctx, actor, opt)

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
		owner = os.Args[1]

		opt := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}

		for {
			rrepos, resp, err := v3Client.Repositories.ListByOrg(ctx, owner, opt)

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

	repoNames := []string{}

	for _, repo := range repos {
		repoNames = append(repoNames, *repo.Name)
	}

	// fmt.Println(repoNames)
	err = setupDirectories(owner, repoNames)
	if err != nil {
		fmt.Println(err)
		// TODO: exit
	}

	for i, repo := range repos {
		fmt.Println(github.Stringify(repo.FullName))

		err := getIssuesAndCommentsForRepository(v4Client, *repo.Name, owner)

		if err != nil {
			fmt.Println(err)
			// TODO: exit
		}

		if i == 1 {
			break
		}
	}
}
