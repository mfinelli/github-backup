package main

import (
	"context"
	"fmt"
)

import "github.com/shurcooL/githubv4"

func getIssuesAndCommentsForRepository(client *githubv4.Client, repo, owner string) error {
	type issue struct {
		Author struct {
			Login string
		}

		Editor struct {
			Login string
		}

		Number int
		Title string
		Body string
		CreatedAt string
		ClosedAt string
		IsPinned bool
		State string

		Assignees struct {
			Nodes []struct{
				// Name string
				Login string
			}
		} `graphql:"assignees(first: 100)"`

		Labels struct {
			Nodes []struct {
				Name string
			}
		} `graphql:"labels(first: 100)"`

		Milestone struct {
			Title string
		}
	}

	var q struct {
		Repository struct {
			Description string
			CreatedAt string
			DiskUsage int // KB
			HomepageURL string
			IsArchived bool
			IsPrivate bool
			IsTemplate bool
			SshURL string

			TemplateRepository struct {
				NameWithOwner string
			}

			Issues struct {
				Nodes []issue
				// Nodes []struct{
				// 	Body string
				// }

				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage githubv4.Boolean
				}
			} `graphql:"issues(first: 100, after: $issuesCursor)"`

			// TODO: rate limiting

			// TODO: branch protection rules
			// TODO: deploy keys (just the ssh bits)
			// TODO: has wiki enabled
			// TODO: use has issues enabled
			// TODO: issues
			// TODO: pull requests
			// TODO: projects
			// TODO: discussions
			// TODO: releases
			// TODO: release artifacts
			// TODO: packages
			// TODO: milestones

			// Issue struct {
			// }
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	vars := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name": githubv4.String(repo),
		"issuesCursor": (*githubv4.String)(nil),
	}

	// err := client.Query(context.Background(), &q, vars)

	// if err != nil {
	// 	return err
	// }

	// fmt.Println(q.Repository.Description)
	// fmt.Println(q.Repository)

	var allIssues []issue
	for {
		err := client.Query(context.Background(), &q, vars)
		if err != nil {
			return err
		}
		allIssues = append(allIssues, q.Repository.Issues.Nodes...)
		if !q.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		vars["issuesCursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)
	}

	fmt.Println(allIssues)

	return nil
}
