package main

// "issue" as returned from the graphql api
type apiIssue struct {
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
	LastEditedAt string
	IsPinned bool
	State string

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
	Number int `yaml:"number"`
	Title string `yaml:"title"`
	Body string `yaml:"-"`
	Author string `yaml:"author"`
	Editor string `yaml:"editor,omitempty"`
	State string `yaml:"state"`
	CreatedAt string `yaml:"created"`
	ClosedAt string `yaml:"closed,omitempty"`
	LastEditedAt string `yaml:"edited,omitempty"`
	IsPinned bool `yaml:"pinned"`
	Assignees []string `yaml:"assignees,omitempty"`
	Labels []string `yaml:"labels,omitempty"`
	Milestone string `yaml:"milestone,omitempty"`
}

func convertApiIssueToIssue(input apiIssue) issue {
	output := issue{
		Number: input.Number,
		Title: input.Title,
		Author: input.Author.Login,
		Editor: input.Editor.Login,
		State: input.State,
		CreatedAt: input.CreatedAt,
		ClosedAt: input.ClosedAt,
		LastEditedAt: input.LastEditedAt,
		IsPinned: input.IsPinned,
		Milestone: input.Milestone.Title,
	}

	for _, assignee := range input.Assignees.Nodes {
		output.Assignees = append(output.Assignees, assignee.Login)
	}

	for _, label := range input.Labels.Nodes {
		output.Labels = append(output.Labels, label.Name)
	}

	return output
}
