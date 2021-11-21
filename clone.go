package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

func cloneRepository(ctx context.Context, config Config, repo repository) error {
	start := time.Now()

	if !config.Quiet {
		fmt.Printf("cloning repository %s\n", repo.Name)
	}

	if config.InternalGit {
		return errors.New(fmt.Sprintf("TODO"))
	} else {
		cmd := exec.Command(config.GitBinaryPath, "clone", "--mirror",
			repo.SshURL, fmt.Sprintf("%s/%s/repository",
				repo.Owner, repo.Name))

		if config.Debug {
			fmt.Println(cmd)
		}

		_, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
	}

	if !config.Quiet {
		fmt.Printf("cloned %s in %v\n", repo.Name,
			time.Now().Sub(start))
	}

	return nil
}
