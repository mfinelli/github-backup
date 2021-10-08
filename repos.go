package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

import "github.com/google/go-github/v39/github"
import "github.com/shurcooL/githubv4"
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

// the v3 API doesn't return template repository information when doing a
// repository list (compared to getting a single repository)
// we eventually want to move everything to the v4 graphql endpoint anyway so
// for now we're just going to make an extra api call per repository to get
// all of the information that we want
func getRepositoryInfo(ctx context.Context, client *githubv4.Client, repo *github.Repository) (repository, error) {
	var q struct {
		Repository struct {
			Description string
			CreatedAt   string
			DiskUsage   int // KB
			HomepageURL string
			IsArchived  bool
			IsPrivate   bool
			IsTemplate  bool
			SshURL      string

			TemplateRepository struct {
				NameWithOwner string
			}
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	vars := map[string]interface{}{
		"owner": githubv4.String(*repo.Owner.Login),
		"name":  githubv4.String(*repo.Name),
	}

	err := client.Query(ctx, &q, vars)
	if err != nil {
		return repository{}, err
	}

	r := repository{
		Owner:              *repo.Owner.Login,
		Name:               *repo.Name,
		FullName:           *repo.FullName,
		Description:        q.Repository.Description,
		HomepageURL:        q.Repository.HomepageURL,
		CreatedAt:          q.Repository.CreatedAt,
		IsArchived:         q.Repository.IsArchived,
		IsPrivate:          q.Repository.IsPrivate,
		IsTemplate:         q.Repository.IsTemplate,
		TemplateRepository: q.Repository.TemplateRepository.NameWithOwner,
		SshURL:             q.Repository.SshURL,
		DiskUsage:          q.Repository.DiskUsage,
	}

	return r, nil
}

func writeRepositoryMetadata(config Config, repo repository) error {
	basepath := filepath.Join(config.BackupPath, repo.Owner, repo.Name)

	err := os.Mkdir(basepath, 0755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&repo)
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
