package main

import (
	"context"
	"fmt"
)

import "github.com/shurcooL/githubv4"

func getIssuesAndCommentsForRepository(client *githubv4.Client, repo, owner string) error {
	var q struct {
		Repository struct {
			Description string
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	vars := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name": githubv4.String(repo),
	}

	err := client.Query(context.Background(), &q, vars)

	if err != nil {
		return err
	}

	fmt.Println(q.Repository.Description)

	return nil
}
