package main

import (
	"context"
	"fmt"
	"sync"
)

import parser "github.com/nikitavoloboev/markdown-parser"

type linkParseResult struct {
	issueNumber   int
	commentNumber int // 0: issue; otherwise comment
	links         map[string]string
}

func downloadAllBodyLinks(ctx context.Context, config Config, repo repository, issues []issue) error {
	results := make(chan linkParseResult)
	var wg sync.WaitGroup

	for _, i := range issues {
		wg.Add(1)
		go func(iNumber, cNumber int, body string) {
			defer wg.Done()
			results <- parseBodyForLinks(iNumber, cNumber, body)
		}(i.Number, 0, i.Body)

		for _, c := range i.Comments {
			wg.Add(1)
			go func(iNo, cNo int, body string) {
				defer wg.Done()
				results <- parseBodyForLinks(iNo, cNo, body)
			}(i.Number, c.Number, c.Body)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.commentNumber == 0 {
			fmt.Printf("checked issue %d\n", r.issueNumber)
		} else {
			fmt.Printf("checked comment %d on issue %d\n", r.commentNumber, r.issueNumber)
		}

		if len(r.links) == 0 {
			fmt.Println("there were no links")
		} else {
			fmt.Println(r.links)
		}
	}

	return nil
}

func parseBodyForLinks(issueNumber, commentNumber int, body string) linkParseResult {
	return linkParseResult{
		issueNumber:   issueNumber,
		commentNumber: commentNumber,
		links:         parser.GetAllLinks(body),
	}
}
