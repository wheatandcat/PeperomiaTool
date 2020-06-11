package main

import (
	"context"
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
}

// ResposeType Graphqlのタイプ
type ResposeType struct {
	Repository Repository `json:"repository"`
}

// Repository Repositoryのタイプ
type Repository struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Milestone Milestone `json:"milestone"`
}

// Milestone Milestoneのタイプ
type Milestone struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
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

	client := graphql.NewClient("https://api.github.com/graphql")

	// make a request
	req := graphql.NewRequest(`
query Repository($owner: String!, $name: String!, $labels: [String!]) {
  repository(owner: $owner, name: $name) {
    id
    name
    url
    milestone(number: 9) {
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

	req.Var("owner", config.GitHub.Owner)
	req.Var("name", config.GitHub.Repository)
	req.Var("labels", [1]string{config.GitHub.Label})

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+config.GitHub.Token)

	ctx := context.Background()

	var respData ResposeType
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	for i, pr := range respData.Repository.Milestone.PullRequests.Nodes {
		start := strings.Index(pr.Body, "## 対応内容")
		if start == -1 {
			start = 0
		} else {
			start += len("## 対応内容")
		}
		end := strings.Index(pr.Body, "## その他")
		if end == -1 {
			end = len(pr.Body)
		}

		body := pr.Body[start:end]
		respData.Repository.Milestone.PullRequests.Nodes[i].Body = body
	}

	text := `
# {{ .Repository.Name }} {{ .Repository.Milestone.Title }} リリースノート

## GitHub

[{{ .Repository.Name }}]({{ .Repository.URL }})

## マイルストーン

[{{ .Repository.Milestone.Title }}]({{ .Repository.Milestone.URL }})

## 対応内容
{{range $index, $pr := .Repository.Milestone.PullRequests.Nodes}}
- [{{ $pr.Title }}]({{ $pr.URL }}){{end}}

## リリース詳細
{{range $index, $pr := .Repository.Milestone.PullRequests.Nodes}}
### {{ $pr.Title }}
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
