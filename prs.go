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

	for _, iPr := range apiPrs {
		comments, err := getPrComments(ctx, client, repo, iPr.Number)
		if err != nil {
			return []pr{}, err
		}

		prs = append(prs, convertApiPrToPr(iPr, comments))
	}

	return prs, nil
}

func getPrComments(ctx context.Context, client *githubv4.Client, repo repository, prNo int) ([]apiComment, error) {

	var q struct {
		Repository struct {
			PullRequest struct {
				Comments struct {
					Nodes    []apiComment
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage githubv4.Boolean
					}
				} `graphql:"comments(first: 100, after: $commentCursor)"`
			} `graphql:"pullRequest(number: $pr)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	vars := map[string]interface{}{
		"owner":         githubv4.String(repo.Owner),
		"name":          githubv4.String(repo.Name),
		"pr":            githubv4.Int(prNo),
		"commentCursor": (*githubv4.String)(nil),
	}

	var comments []apiComment

	for {
		err := client.Query(ctx, &q, vars)
		if err != nil {
			return []apiComment{}, err
		}

		comments = append(comments, q.Repository.PullRequest.Comments.Nodes...)
		if !q.Repository.PullRequest.Comments.PageInfo.HasNextPage {
			break
		}
		vars["commentCursor"] = githubv4.NewString(q.Repository.PullRequest.Comments.PageInfo.EndCursor)
	}

	return comments, nil

}

func writePrsToDisk(config Config, repo repository, prs []pr) error {
	if len(prs) == 0 {
		return nil
	}

	basepath := filepath.Join(config.BackupPath, repo.Owner, repo.Name, "pulls")

	err := os.Mkdir(basepath, 0755)
	if err != nil {
		return err
	}

	for _, i := range prs {
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

		err = os.WriteFile(filepath.Join(p, "pr.yml"), []byte(meta), 0644)
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
