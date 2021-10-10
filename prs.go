package main

import (
	"context"
	"fmt"
	// "os"
	// "path/filepath"
	// "strconv"
)

import "github.com/shurcooL/githubv4"

// import "gopkg.in/yaml.v2"

func getRepositoryPrs(ctx context.Context, client *githubv4.Client, repo repository) ([]pr, error) {
	var q struct {
		Repository struct {
			PullRequests struct {
				Nodes    []apiPr
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage githubv4.Boolean
				}
			} `graphql:"pullRequests(first: 100, after: $prCursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	vars := map[string]interface{}{
		"owner":    githubv4.String(repo.Owner),
		"name":     githubv4.String(repo.Name),
		"prCursor": (*githubv4.String)(nil),
	}

	var prs []pr
	var apiPrs []apiPr

	for {
		err := client.Query(ctx, &q, vars)
		if err != nil {
			return []pr{}, err
		}

		apiPrs = append(apiPrs, q.Repository.PullRequests.Nodes...)
		if !q.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}
		vars["prCursor"] = githubv4.NewString(q.Repository.PullRequests.PageInfo.EndCursor)
	}

	fmt.Println(apiPrs)

	return prs, nil
}
