package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

import "github.com/google/go-github/v39/github"
import "gopkg.in/yaml.v2"

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

func writeRepositoryMetadata(config Config, repo *github.Repository) error {
	basepath := filepath.Join(config.BackupPath, *repo.Owner.Login,
		*repo.Name)

	err := os.Mkdir(basepath, 0755)
	if err != nil {
		return err
	}

	r := convertGithubRepositoryToRepository(repo)

	data, err := yaml.Marshal(&r)
	if err != nil {
		return nil
	}

	meta := fmt.Sprintf("---\n%s\n", string(data))
	err = os.WriteFile(filepath.Join(basepath, "repository.yml"),
		[]byte(meta), 0644)
	if err != nil {
		return err
	}

	return nil
}
