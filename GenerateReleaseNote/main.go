package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/machinebox/graphql"
)

// Config 設定ファイルのタイプ
type Config struct {
	GitHub GitHubConfig
}

// GitHubConfig GitHub設定ファイルのタイプ
type GitHubConfig struct {
	Token      string `toml:"token"`
	Owner      string `toml:"owner"`
	Repository string `toml:"repository"`
	Label      string `toml:"label"`
	StartText  string `toml:"startText"`
	EndText    string `toml:"endText"`
}

// ResposeType Graphqlのタイプ
type ResposeType struct {
	Repository Repository `json:"repository"`
}

// Repository Repositoryのタイプ
type Repository struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	URL        string         `json:"url"`
	Milestone  Milestone      `json:"milestone"`
	Milestones MilestoneNodes `json:"milestones"`
}

// MilestoneNodes MilestoneNodesのタイプ
type MilestoneNodes struct {
	Nodes []Milestone `json:"nodes"`
}

// Milestone Milestoneのタイプ
type Milestone struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Number       int          `json:"number"`
	URL          string       `json:"url"`
	PullRequests PullRequests `json:"pullRequests"`
}

// PullRequests PullRequestsのタイプ
type PullRequests struct {
	Nodes []PullRequest `json:"nodes"`
}

// PullRequest PullRequestのタイプ
type PullRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
}

func main() {
	var config Config
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	mn, err := config.GitHub.getMilestoneNumber(os.Getenv("MILESTONE"))
	if err != nil {
		log.Fatal(err)
	}

	respData, err := config.GitHub.getPullRequest(mn)
	if err != nil {
		log.Fatal(err)
	}

	for i, pr := range respData.Repository.Milestone.PullRequests.Nodes {
		start := strings.Index(pr.Body, config.GitHub.StartText)
		if start == -1 {
			start = 0
		} else {
			start += len(config.GitHub.StartText)
		}
		end := strings.Index(pr.Body, config.GitHub.EndText)
		if end == -1 {
			end = len(pr.Body)
		}

		body := pr.Body[start:end]
		respData.Repository.Milestone.PullRequests.Nodes[i].Body = body
	}

	text := `
タイトル: {{ .Repository.Name }} {{ .Repository.Milestone.Title }} リリースノート

# GitHub

[{{ .Repository.Name }}]({{ .Repository.URL }})

# マイルストーン

[{{ .Repository.Milestone.Title }}]({{ .Repository.Milestone.URL }})

# 対応内容
{{range $index, $pr := .Repository.Milestone.PullRequests.Nodes}}
- [{{ $pr.Title }}](#{{ $pr.ID }}){{end}}

# リリース詳細
{{range $index, $pr := .Repository.Milestone.PullRequests.Nodes}}<a id="{{ $pr.ID }}"></a>
## [{{ $pr.Title }}]({{ $pr.URL }})
{{ $pr.Body}}
{{end}}
`

	tpl, err := template.New("").Parse(text)
	if err != nil {
		log.Fatal(err)
	}

	if err := tpl.Execute(os.Stdout, respData); err != nil {
		log.Fatal(err)
	}
}

func (c *GitHubConfig) getMilestoneNumber(mt string) (int, error) {
	client := graphql.NewClient("https://api.github.com/graphql")
	req := graphql.NewRequest(`
query Repository($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    id
    milestones(first: 20, orderBy: {field: CREATED_AT, direction: DESC}) {
      nodes {
        id
        number
        title
      }
    }
  }
}
`)
	req.Var("owner", c.Owner)
	req.Var("name", c.Repository)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData ResposeType
	if err := client.Run(ctx, req, &respData); err != nil {
		return 0, err
	}

	mn := 0
	for _, m := range respData.Repository.Milestones.Nodes {
		if m.Title == mt {
			mn = m.Number
			break
		}
	}

	if mn == 0 {
		return 0, fmt.Errorf("Error: '%s' milestone title not found", mt)
	}

	return mn, nil
}

func (c *GitHubConfig) getPullRequest(mn int) (ResposeType, error) {
	client := graphql.NewClient("https://api.github.com/graphql")

	req := graphql.NewRequest(`
query Repository($owner: String!, $name: String!, $labels: [String!], $milestoneNumber: Int!) {
  repository(owner: $owner, name: $name) {
    id
    name
    url
    milestone(number: $milestoneNumber) {
      id
      title
      url
      pullRequests(first: 100, labels: $labels) {
        nodes {
          id
          title
          url
          body
        }
      }
    }
  }
}
`)

	req.Var("owner", c.Owner)
	req.Var("name", c.Repository)
	req.Var("labels", [1]string{c.Label})
	req.Var("milestoneNumber", mn)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData ResposeType
	err := client.Run(ctx, req, &respData)
	return respData, err
}
