package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

import "golang.org/x/oauth2"
import "github.com/alecthomas/kong"
import "github.com/google/go-github/v39/github"
import "github.com/shurcooL/githubv4"

type CLI struct {
	Organization string `help:"Backup an organization's repositories." arg:"" optional:""`

	Debug bool `help:"Enable extra verbose output." short:"d"`
	Quiet bool `help:"Supress all standard output." short:"q"`

	GithubToken string `help:"PAT to access the GitHub API" group:"Backup options:" short:"t" env:"GITHUB_TOKEN" placeholder:"\"...\""`
	Path        string `help:"Write the backup to the given path" group:"Backup options:" short:"p" placeholder:"\".\""`
	Repository  string `help:"Backup a single repository by name" group:"Backup options:" short:"r" placeholder:"\"\""`

	GitBinary   string `help:"Path to external git binary" group:"Git flags:" xor:"gitbinary" placeholder:"\"git\""`
	NoGitBinary bool   `help:"Do not use an external git binary" group:"Git flags:" xor:"gitbinary"`
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "--version" {
		fmt.Println("ghb version 1")
	} else {
		ctx := context.Background()

		var cli CLI
		kong.Parse(&cli,
			kong.Description("Backup the user's GitHub repositories and associated data."))

		os.Exit(run(ctx, cli))
	}
}

func run(ctx context.Context, cli CLI) int {
	config, err := ValidateConfig(cli)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	auth := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GithubToken},
	)

	httpClient := oauth2.NewClient(ctx, auth)
	v3Client := github.NewClient(httpClient)
	v4Client := githubv4.NewClient(httpClient)

	actor, err := getCurrentActor(v4Client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	var repositories []*github.Repository

	if cli.Organization == "" && cli.Repository == "" {
		repositories, err = getUserRepositories(ctx, v3Client, actor)
	} else if cli.Organization == "" && cli.Repository != "" {
		repositories, err = getSingleRepository(ctx, v3Client, actor,
			cli.Repository)
	} else if cli.Organization != "" && cli.Repository == "" {
		repositories, err = getOrgRepositories(ctx, v3Client,
			cli.Organization)
	} else if cli.Organization != "" && cli.Repository != "" {
		repositories, err = getSingleRepository(ctx, v3Client,
			cli.Organization, cli.Repository)
	} else {
		fmt.Fprintf(os.Stderr, "given unsupported options\n")
		return 25
	}

	// this is from fetching the repositories above
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	if len(repositories) != 0 {
		err = os.Mkdir(filepath.Join(config.BackupPath,
			*repositories[0].Owner.Login), 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}
	}

	for _, r := range repositories {
		repo, err := getRepositoryInfo(ctx, v4Client, r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}

		err = writeRepositoryMetadata(config, repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}

		issues, err := getRepositoryIssues(ctx, v4Client, repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}

		prs, err := getRepositoryPrs(ctx, v4Client, repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}

		err = writeIssuesToDisk(config, repo, issues)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}

		err = writePrsToDisk(config, repo, prs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}

		err = downloadAllBodyLinks(ctx, config, repo, issues)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}
	}

	return 0
}
