package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

import "gopkg.in/yaml.v2"

func writeIssueToFile(owner, repo string, i issue) error {
	p := filepath.Join(owner, repo, "issues", strconv.Itoa(i.Number))
	err := os.Mkdir(p, 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(p, "body.md"), []byte(i.Body), 0644)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&i)
	if err != nil {
		return err
	}

	meta := fmt.Sprintf("---\n%s\n", string(data))
	err = os.WriteFile(filepath.Join(p, "issue.yml"), []byte(meta), 0644)
	if err != nil {
		return err
	}

	return nil
}
