package main

import (
	"context"
)

import "github.com/shurcooL/githubv4"

func getCurrentActor(client *githubv4.Client) (string, error) {
	var q struct {
		Viewer struct {
			Login string
		}
	}

	err := client.Query(context.Background(), &q, nil)

	if err != nil {
		return "", err
	}

	return q.Viewer.Login, nil
}
