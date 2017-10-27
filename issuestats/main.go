package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	owner := "syncthing"
	repo := "syncthing"
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	for i := 2550; i < 3980; i++ {
		issue, _, err := client.Issues.Get(owner, repo, i)
		if err != nil {
			log.Fatal(err)
		}
		if *issue.State != "closed" {
			continue
		}

		evs, _, err := client.Issues.ListIssueEvents(owner, repo, i, nil)
		if err != nil {
			log.Fatal(err)
		}

		withCommit := false
		closed := false
		var delay time.Duration
		for _, ev := range evs {
			if *ev.Event == "closed" {
				closed = true
				if ev.CommitID != nil {
					withCommit = true
				}
				delay = ev.CreatedAt.Sub(*issue.CreatedAt)
			}
		}
		if closed {
			fmt.Printf("%d,%v,%d\n", i, withCommit, int(delay.Hours()/24))
		}
	}
}
