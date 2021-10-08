package main

import (
	"time"
)

import "github.com/google/go-github/v39/github"

// massaged "repository" for marshaling into yaml
type repository struct {
	Name               string `yaml:"repository"`
	Description        string `yaml:"description,omitempty"`
	HomepageURL        string `yaml:"homepage,omitempty"`
	CreatedAt          string `yaml:"created"`
	IsArchived         bool   `yaml:"archived"`
	IsPrivate          bool   `yaml:"private"`
	IsTemplate         bool   `yaml:"-"`
	TemplateRepository string `yaml:"template,omitempty"`
	SshURL             string `yaml:"ssh"`
	DiskUsage          int    `yaml:"size"` // KB
}

// "issue" as returned from the graphql api
type apiIssue struct {
	Author struct {
		Login string
	}

	Editor struct {
		Login string
	}

	Number       int
	Title        string
	Body         string
	CreatedAt    string
	ClosedAt     string
	LastEditedAt string
	IsPinned     bool
	State        string

	Assignees struct {
		Nodes []struct {
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

// massaged "issue" for marshaling into yaml
type issue struct {
	Number       int       `yaml:"number"`
	Title        string    `yaml:"title"`
	Body         string    `yaml:"-"`
	Author       string    `yaml:"author"`
	Editor       string    `yaml:"editor,omitempty"`
	State        string    `yaml:"state"`
	CreatedAt    string    `yaml:"created"`
	ClosedAt     string    `yaml:"closed,omitempty"`
	LastEditedAt string    `yaml:"edited,omitempty"`
	IsPinned     bool      `yaml:"pinned"`
	Assignees    []string  `yaml:"assignees,omitempty"`
	Labels       []string  `yaml:"labels,omitempty"`
	Milestone    string    `yaml:"milestone,omitempty"`
	Comments     []comment `yaml:"-"`
}

// issue "comment" as returned by the graphql endpoint
type apiComment struct {
	DatabaseId int
	Author     struct {
		Login string
	}

	Editor struct {
		Login string
	}

	Body         string
	CreatedAt    string
	LastEditedAt string
}

type comment struct {
	Body         string `yaml:"-"`
	DatabaseId   int    `yaml:"id"`
	Author       string `yaml:"author"`
	Editor       string `yaml:"editor,omitempty"`
	CreatedAt    string `yaml:"created"`
	LastEditedAt string `yaml:"edited,omitempty"`
}

func convertApiIssueToIssue(input apiIssue, comments []apiComment) issue {
	output := issue{
		Number:       input.Number,
		Title:        input.Title,
		Body:         input.Body,
		Author:       input.Author.Login,
		Editor:       input.Editor.Login,
		State:        input.State,
		CreatedAt:    input.CreatedAt,
		ClosedAt:     input.ClosedAt,
		LastEditedAt: input.LastEditedAt,
		IsPinned:     input.IsPinned,
		Milestone:    input.Milestone.Title,
	}

	for _, assignee := range input.Assignees.Nodes {
		output.Assignees = append(output.Assignees, assignee.Login)
	}

	for _, label := range input.Labels.Nodes {
		output.Labels = append(output.Labels, label.Name)
	}

	for _, c := range comments {
		output.Comments = append(output.Comments, comment{
			Body:         c.Body,
			DatabaseId:   c.DatabaseId,
			Author:       c.Author.Login,
			Editor:       c.Editor.Login,
			CreatedAt:    c.CreatedAt,
			LastEditedAt: c.LastEditedAt,
		})
	}

	return output
}

func convertGithubRepositoryToRepository(repo *github.Repository) repository {
	r := repository{
		Name:        *repo.FullName,
		Description: *repo.Description,
		HomepageURL: *repo.Homepage,
		CreatedAt:   repo.CreatedAt.Format(time.RFC3339),
		IsArchived:  *repo.Archived,
		IsPrivate:   *repo.Private,
		IsTemplate:  *repo.IsTemplate,
		// TemplateRepository: *repo.TemplateRepository.FullName,
		SshURL:    *repo.SSHURL,
		DiskUsage: *repo.Size,
	}

	if repo.TemplateRepository != nil {
		r.TemplateRepository = *repo.TemplateRepository.FullName
	}

	return r
}
