package main

import (
	"os"
	"path/filepath"
)

// TODO: we either need to fail on a directory that's not empty or skip
//       directories that already exist
func setupDirectories(owner string, repos []string) error {
	subs := []string{"issues", "pulls"}

	err := os.Mkdir(owner, 0755)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		err := os.Mkdir(filepath.Join(owner, repo), 0755)
		if err != nil {
			return err
		}

		for _, sub := range subs {
			err = os.Mkdir(filepath.Join(owner, repo, sub), 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
