package main

// massaged "repository" for marshaling into yaml
type repository struct {
	Owner              string `yaml:"-"`
	Name               string `yaml:"-"`
	FullName           string `yaml:"repository"`
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

// "pr" as returned from the graphql api
type apiPr struct {
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
	IsDraft      bool
	State        string

	MergedAt string
	MergedBy struct {
		Login string
	}

	BaseRefName    string
	BaseRepository struct {
		NameWithOwner string
	}
	HeadRefName    string
	HeadRepository struct {
		NameWithOwner string
	}

	Assignees struct {
		Nodes []struct {
			Login string
		}
	} `graphql:"assignees(first: 100)"`

	LatestReviews struct {
		Nodes []struct {
			Author struct {
				Login string
			}

			State       string
			SubmittedAt string
		}
	} `graphql:"latestReviews(first: 100)"`

	Labels struct {
		Nodes []struct {
			Name string
		}
	} `graphql:"labels(first: 100)"`

	Milestone struct {
		Title string
	}
}

// massaged "pr" for marshaling into yaml
type pr struct {
	Number       int           `yaml:"number"`
	Title        string        `yaml:"title"`
	Body         string        `yaml:"-"`
	Author       string        `yaml:"author"`
	Editor       string        `yaml:"editor,omitempty"`
	State        string        `yaml:"state"`
	CreatedAt    string        `yaml:"created"`
	ClosedAt     string        `yaml:"closed,omitempty"`
	LastEditedAt string        `yaml:"edited,omitempty"`
	Merged       prMerged      `yaml:"merged"`
	IsDraft      bool          `yaml:"draft"`
	PullRequest  prPullRequest `yaml:"pull_request"`
	Reviews      []prReview    `yaml:"reviews"`
	Assignees    []string      `yaml:"assignees,omitempty"`
	Labels       []string      `yaml:"labels,omitempty"`
	Milestone    string        `yaml:"milestone,omitempty"`
	Comments     []comment     `yaml:"-"`
}

type prMerged struct {
	MergedOn string `yaml:"on"`
	MergedBy string `yaml:"by"`
}

type prPullRequest struct {
	Target     string `yaml:"target"`
	Source     string `yaml:"source"`
	Repository string `yaml:"repository"`
}

type prReview struct {
	Author string `yaml:"author"`
	Review string `yaml:"review"`
	Date   string `yaml:"date"`
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
	Number       int    `yaml:"-"`
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

	for i, c := range comments {
		output.Comments = append(output.Comments, comment{
			Number:       i + 1, // don't do zero-based comments
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

func convertApiPrToPr(input apiPr, comments []apiComment) pr {
	output := pr{
		Number:       input.Number,
		Title:        input.Title,
		Body:         input.Body,
		Author:       input.Author.Login,
		Editor:       input.Editor.Login,
		State:        input.State,
		CreatedAt:    input.CreatedAt,
		ClosedAt:     input.ClosedAt,
		LastEditedAt: input.LastEditedAt,
		IsDraft:      input.IsDraft,
		Merged: prMerged{
			MergedOn: input.MergedAt,
			MergedBy: input.MergedBy.Login,
		},
		PullRequest: prPullRequest{
			Target:     input.BaseRefName,
			Source:     input.HeadRefName,
			Repository: input.HeadRepository.NameWithOwner,
		},
		Milestone: input.Milestone.Title,
	}

	for _, assignee := range input.Assignees.Nodes {
		output.Assignees = append(output.Assignees, assignee.Login)
	}

	for _, label := range input.Labels.Nodes {
		output.Labels = append(output.Labels, label.Name)
	}

	for _, review := range input.LatestReviews.Nodes {
		output.Reviews = append(output.Reviews, prReview{
			Author: review.Author.Login,
			Review: review.State,
			Date:   review.SubmittedAt,
		})
	}

	for i, c := range comments {
		output.Comments = append(output.Comments, comment{
			Number:       i + 1, // don't do zero-based comments
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
