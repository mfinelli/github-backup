package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

import "gopkg.in/yaml.v2"

func writeIssueToFiles(owner, repo string, i issue) error {
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

	if len(i.Comments) > 0 {
		cPath := filepath.Join(p, "comments")
		err := os.Mkdir(cPath, 0755)
		if err != nil {
			return err
		}

		err = writeCommentsToFiles(cPath, i.Comments)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeCommentsToFiles(commentsPath string, comments []comment) error {
	for i, c := range comments {
		index := i + 1 // don't do 0-based directories

		p := filepath.Join(commentsPath, strconv.Itoa(index))
		err := os.Mkdir(p, 0755)
		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(p, "body.md"), []byte(c.Body),
			0644)
		if err != nil {
			return err
		}

		data, err := yaml.Marshal(&c)
		if err != nil {
			return err
		}

		meta := fmt.Sprintf("---\n%s\n", string(data))
		err = os.WriteFile(filepath.Join(p, "comment.yml"),
			[]byte(meta), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
