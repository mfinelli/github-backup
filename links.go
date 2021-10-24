package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

import parser "github.com/nikitavoloboev/markdown-parser"

type linkParseResult struct {
	issueNumber   int
	prNumber      int
	commentNumber int // 0: issue/pr; otherwise comment
	links         map[string]string
}

type linkParseRequest struct {
	issueNumber   int
	prNumber      int
	commentNumber int
	body          string
}

type normalizedIssuePr struct {
	number int
	path   string
}

func downloadAllBodyLinks(ctx context.Context, config Config, repo repository, issues []issue, prs []pr) error {
	// images and videos are _always_ public even on private repos
	publicFilesPrefix := "https://user-images.githubusercontent.com/"
	// other file types require an authenticated session to download
	privateFilesPrefix := fmt.Sprintf("https://github.com/%s/%s/files/",
		repo.Owner, repo.Name)

	requests := make(chan linkParseRequest)
	results := make(chan linkParseResult)
	numWorkers := 10
	var wg sync.WaitGroup

	go func() {
		for _, i := range issues {
			requests <- linkParseRequest{
				issueNumber:   i.Number,
				prNumber:      0,
				commentNumber: 0,
				body:          i.Body,
			}

			for _, c := range i.Comments {
				requests <- linkParseRequest{
					issueNumber:   i.Number,
					prNumber:      0,
					commentNumber: c.Number,
					body:          c.Body,
				}
			}
		}

		for _, i := range prs {
			requests <- linkParseRequest{
				issueNumber:   0,
				prNumber:      i.Number,
				commentNumber: 0,
				body:          i.Body,
			}

			for _, c := range i.Comments {
				requests <- linkParseRequest{
					issueNumber:   0,
					prNumber:      i.Number,
					commentNumber: c.Number,
					body:          c.Body,
				}
			}
		}

		close(requests)
	}()

	for j := 0; j < numWorkers; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for request := range requests {
				results <- parseBodyForLinks(request)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if len(r.links) == 0 {
			continue
		}

		var pth string
		normed := getIssueOrPr(&r)
		basepath := filepath.Join(config.BackupPath, repo.Owner,
			repo.Name, normed.path, strconv.Itoa(normed.number))

		if r.commentNumber == 0 {
			pth = filepath.Join(basepath, "files")
		} else {
			pth = filepath.Join(basepath, "comments",
				strconv.Itoa(r.commentNumber), "files")
		}

		err := os.Mkdir(pth, 0755)
		if err != nil {
			return err
		}

		for _, url := range r.links {
			if strings.HasPrefix(url, privateFilesPrefix) {
				if repo.IsPrivate {
					// TODO: we can't handle these yet
					// https://github.com/mfinelli/github-backup/issues/1
					continue
				} else {
					if !config.Quiet {
						fmt.Printf("downloading file %s\n", url)
					}

					out := filepath.Join(pth, path.Base(url))
					err = downloadPublicFile(config, out, url)
					if err != nil {
						return err
					}
				}
			} else if strings.HasPrefix(url, publicFilesPrefix) {
				if !config.Quiet {
					fmt.Printf("downloading file %s\n", url)
				}

				out := filepath.Join(pth, path.Base(url))
				err = downloadPublicFile(config, out, url)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("WARNING: unknown file type: %s\n",
					url)
				continue
			}
		}

	}

	return nil
}

func parseBodyForLinks(r linkParseRequest) linkParseResult {
	return linkParseResult{
		issueNumber:   r.issueNumber,
		prNumber:      r.prNumber,
		commentNumber: r.commentNumber,
		links:         parser.GetAllLinks(r.body),
	}
}

func downloadPublicFile(config Config, p, url string) error {
	start := time.Now()

	fp, err := os.Create(p)
	if err != nil {
		return err
	}
	defer fp.Close()

	data, err := http.Get(url)
	if err != nil {
		return err
	}
	defer data.Body.Close()

	size, err := io.Copy(fp, data.Body)
	if err != nil {
		return err
	}

	if !config.Quiet {
		fmt.Printf("downloaded %d bytes in %v\n", size,
			time.Now().Sub(start))
	}

	return nil
}

func getIssueOrPr(r *linkParseResult) normalizedIssuePr {
	if r.issueNumber == 0 {
		return normalizedIssuePr{r.prNumber, "prs"}
	} else {
		return normalizedIssuePr{r.issueNumber, "issues"}
	}
}
