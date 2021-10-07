package main

import (
	"context"
)

import "github.com/google/go-github/v39/github"

func getUserRepositories(ctx context.Context, client *github.Client, owner string) ([]*github.Repository, error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var repositories []*github.Repository

	for {
		repos, resp, err := client.Repositories.List(ctx, owner, opt)
		if err != nil {
			return []*github.Repository{}, err
		}

		repositories = append(repositories, repos...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return repositories, nil
}

func getOrgRepositories(ctx context.Context, client *github.Client, org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var repositories []*github.Repository

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			return []*github.Repository{}, err
		}

		repositories = append(repositories, repos...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return repositories, nil
}

func getSingleRepository(ctx context.Context, client *github.Client, owner, repo string) ([]*github.Repository, error) {
	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return []*github.Repository{}, err
	}

	return []*github.Repository{repository}, nil
}
