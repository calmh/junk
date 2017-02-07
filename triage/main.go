package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	addr         = ":8080"
	owner        = "syncthing"
	repo         = "syncthing"
	token        = os.Getenv("GITHUB_TOKEN")
	triageAge    = 7 * 24 * time.Hour
	oldBugCutoff = 365 * 24 * time.Hour
)

func main() {
	flag.StringVar(&addr, "addr", addr, "Listen address")
	flag.StringVar(&owner, "owner", owner, "Owner")
	flag.StringVar(&repo, "repo", repo, "Repository")
	flag.StringVar(&token, "token", token, "Token")
	flag.Parse()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		triageIssues, err := issuesRequiringTriage(client)
		if err != nil {
			log.Fatal(err)
		}
		oldBugs, err := oldBugs(client)
		if err != nil {
			log.Fatal(err)
		}

		tpl := template.New("index.tpl.html")
		tpl.Funcs(template.FuncMap{
			"daysAgo": daysAgo,
		})
		template.Must(tpl.ParseFiles("index.tpl.html"))
		tpl.Execute(w, map[string]interface{}{
			"Triage":  triageIssues,
			"OldBugs": oldBugs,
		})
	})
	http.ListenAndServe(addr, nil)
}

func daysAgo(t time.Time) string {
	days := time.Since(t).Seconds() / 86400
	if days > 365 {
		return fmt.Sprintf("%.1f years", days/365)
	}
	if days > 60 {
		return fmt.Sprintf("%.1f months", days/30)
	}
	if days > 14 {
		return fmt.Sprintf("%.0f weeks", days/7)
	}
	return fmt.Sprintf("%.0f days", days)
}

func issuesRequiringTriage(client *github.Client) ([]*github.Issue, error) {
	opts := &github.IssueListByRepoOptions{
		Milestone: "none",
		Direction: "asc",
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	}

	issues, err := getIssues(client, opts)
	if err != nil {
		return nil, err
	}

	var triage []*github.Issue
	for _, i := range issues {
		if i.PullRequestLinks != nil {
			continue
		}
		if i.CreatedAt != nil && time.Since(*i.CreatedAt) > triageAge {
			triage = append(triage, i)
		}
	}

	return triage, nil
}

func oldBugs(client *github.Client) ([]*github.Issue, error) {
	opts := &github.IssueListByRepoOptions{
		Direction: "asc",
		Labels:    []string{"bug"},
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	}

	issues, err := getIssues(client, opts)
	if err != nil {
		return nil, err
	}

	var old []*github.Issue
	for _, i := range issues {
		if i.PullRequestLinks != nil {
			continue
		}
		if i.CreatedAt != nil && time.Since(*i.CreatedAt) > oldBugCutoff {
			old = append(old, i)
		}
	}

	return old, nil
}

func getIssues(client *github.Client, opts *github.IssueListByRepoOptions) ([]*github.Issue, error) {
	var allIssues []*github.Issue
	for {
		issues, resp, err := client.Issues.ListByRepo(owner, repo, opts)
		if err != nil {
			return nil, err
		}

		allIssues = append(allIssues, issues...)

		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}

	return allIssues, nil
}
