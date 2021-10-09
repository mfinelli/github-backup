package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

import "github.com/shurcooL/githubv4"
import "gopkg.in/yaml.v2"

func getRepositoryIssues(ctx context.Context, client *githubv4.Client, repo repository) ([]issue, error) {
	var q struct {
		Repository struct {
			Issues struct {
				Nodes    []apiIssue
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage githubv4.Boolean
				}
			} `graphql:"issues(first: 100, after: $issueCursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	vars := map[string]interface{}{
		"owner":       githubv4.String(repo.Owner),
		"name":        githubv4.String(repo.Name),
		"issueCursor": (*githubv4.String)(nil),
	}

	var issues []issue
	var apiIssues []apiIssue

	for {
		err := client.Query(ctx, &q, vars)
		if err != nil {
			return []issue{}, err
		}

		apiIssues = append(apiIssues, q.Repository.Issues.Nodes...)
		if !q.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		vars["issueCursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)
	}

	for _, iIssue := range apiIssues {
		comments, err := getIssueComments(ctx, client, repo, iIssue.Number)
		if err != nil {
			return []issue{}, err
		}

		issues = append(issues, convertApiIssueToIssue(iIssue, comments))
	}

	return issues, nil
}

func getIssueComments(ctx context.Context, client *githubv4.Client, repo repository, issueNo int) ([]apiComment, error) {
	var q struct {
		Repository struct {
			Issue struct {
				Comments struct {
					Nodes    []apiComment
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage githubv4.Boolean
					}
				} `graphql:"comments(first: 100, after: $commentCursor)"`
			} `graphql:"issue(number: $issue)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	vars := map[string]interface{}{
		"owner":         githubv4.String(repo.Owner),
		"name":          githubv4.String(repo.Name),
		"issue":         githubv4.Int(issueNo),
		"commentCursor": (*githubv4.String)(nil),
	}

	var comments []apiComment

	for {
		err := client.Query(ctx, &q, vars)
		if err != nil {
			return []apiComment{}, err
		}

		comments = append(comments, q.Repository.Issue.Comments.Nodes...)
		if !q.Repository.Issue.Comments.PageInfo.HasNextPage {
			break
		}
		vars["commentCursor"] = githubv4.NewString(q.Repository.Issue.Comments.PageInfo.EndCursor)
	}

	return comments, nil

}

func writeIssuesToDisk(config Config, repo repository, issues []issue) error {
	if len(issues) == 0 {
		return nil
	}

	basepath := filepath.Join(config.BackupPath, repo.Owner, repo.Name, "issues")

	err := os.Mkdir(basepath, 0755)
	if err != nil {
		return err
	}

	for _, i := range issues {
		p := filepath.Join(basepath, strconv.Itoa(i.Number))
		err = os.Mkdir(p, 0755)
		if err != nil {
			return err
		}

		data, err := yaml.Marshal(&i)
		if err != nil {
			return err
		}

		meta := fmt.Sprintf("---\n%s\n", string(data))

		err = os.WriteFile(filepath.Join(p, "body.md"), []byte(i.Body), 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(p, "issue.yml"), []byte(meta), 0644)
		if err != nil {
			return err
		}

		if len(i.Comments) > 0 {
			cp := filepath.Join(p, "comments")
			err = os.Mkdir(cp, 0755)
			if err != nil {
				return err
			}

			for _, c := range i.Comments {
				cip := filepath.Join(cp, strconv.Itoa(c.Number))

				err = os.Mkdir(cip, 0755)
				if err != nil {
					return err
				}

				cdata, err := yaml.Marshal(&c)
				if err != nil {
					return err
				}

				cmeta := fmt.Sprintf("---\n%s\n", string(cdata))

				err = os.WriteFile(filepath.Join(cip, "body.md"), []byte(c.Body), 0644)
				if err != nil {
					return err
				}

				err = os.WriteFile(filepath.Join(cip, "comment.yml"), []byte(cmeta), 0644)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
