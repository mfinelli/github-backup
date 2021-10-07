package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Debug         bool
	Quiet         bool
	GithubToken   string
	BackupPath    string
	GitBinaryPath string
	InternalGit   bool
}

func ValidateConfig(cli CLI) (Config, error) {
	c := Config{
		Debug: cli.Debug,
		Quiet: cli.Quiet,
	}

	if cli.GithubToken == "" {
		return Config{}, errors.New("No GitHub token provided")
	}

	c.GithubToken = cli.GithubToken

	if cli.Path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return Config{}, err
		}

		c.BackupPath = cwd
	} else {
		p, err := filepath.Abs(cli.Path)
		if err != nil {
			return Config{}, err
		}

		info, err := os.Stat(p)
		if os.IsNotExist(err) {
			return Config{}, errors.New(fmt.Sprintf(
				"%s does not exist", p))
		}

		if !info.IsDir() {
			return Config{}, errors.New(fmt.Sprintf(
				"%s is not a directory", p))
		}

		c.BackupPath = p
	}

	fp, err := os.Open(c.BackupPath)
	if err != nil {
		return Config{}, err
	}
	defer fp.Close()

	_, err = fp.Readdirnames(1)
	if err == nil {
		return Config{}, errors.New(fmt.Sprintf(
			"%s is not empty, refusing to clobber", c.BackupPath))
	} else if err != io.EOF {
		return Config{}, err
	} // else the directory is empty!

	if cli.NoGitBinary {
		c.GitBinaryPath = ""
		c.InternalGit = true
	} else {
		c.InternalGit = false

		if cli.GitBinary == "" {
			c.GitBinaryPath = "git" // use whatever git in $PATH
		} else {
			gp, err := filepath.Abs(cli.GitBinary)
			if err != nil {
				return Config{}, err
			}

			c.GitBinaryPath = gp
			_, err = os.Stat(c.GitBinaryPath)
			if os.IsNotExist(err) {
				return Config{}, errors.New(fmt.Sprintf(
					"%s does not exist", c.GitBinaryPath))
			}
		}

		// test to make sure that we actually have a git
		cmd := exec.Command(c.GitBinaryPath, "version")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return Config{}, err
		}

		if !strings.Contains(string(out), "git version") {
			return Config{}, errors.New(fmt.Sprintf(
				"%s doesn't appear to be a git",
				c.GitBinaryPath))
		}
	}

	return c, nil
}
