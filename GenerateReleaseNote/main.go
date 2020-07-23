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
	Token        string   `toml:"token"`
	Owner        string   `toml:"owner"`
	Repositories []string `toml:"repositories"`
	Label        string   `toml:"label"`
	StartText    string   `toml:"startText"`
	EndText      string   `toml:"endText"`
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

// Template Templateのタイプ
type Template struct {
	Mileston            string
	AppMilestoneURL     string
	App                 []PullRequest
	BackendMilestoneURL string
	Backend             []PullRequest
	WebMilestoneURL     string
	Web                 []PullRequest
	HelpMilestoneURL    string
	Help                []PullRequest
	ToolMilestoneURL    string
	Tool                []PullRequest
	LPMilestoneURL      string
	LP                  []PullRequest
}

func main() {
	var config Config
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	var tmp Template

	for _, r := range config.GitHub.Repositories {
		prs, mu, err := config.GitHub.getPullRequestByMilestone(r)
		if err != nil {
			continue
		}

		switch r {
		case "Peperomia":
			tmp.App = prs
			tmp.AppMilestoneURL = mu
		case "PeperomiaBackend":
			tmp.Backend = prs
			tmp.BackendMilestoneURL = mu
		case "PeperomiaWeb":
			tmp.Web = prs
			tmp.WebMilestoneURL = mu
		case "PeperomiaHelp":
			tmp.Help = prs
			tmp.HelpMilestoneURL = mu
		case "PeperomiaWebSite":
			tmp.LP = prs
			tmp.LPMilestoneURL = mu
		case "PeperomiaTool":
			tmp.Tool = prs
			tmp.ToolMilestoneURL = mu
		}
	}

	tmp.Mileston = os.Getenv("MILESTONE")

	text := `
タイトル: ペペロミア {{ .Mileston }} リリースノート

# GitHub

[Peperomia](https://github.com/wheatandcat/Peperomia)

# マイルストーン

{{ if .AppMilestoneURL }}[Peperomia {{ .Mileston }}]({{ .AppMilestoneURL }})
{{ end }}
{{ if .BackendMilestoneURL }}[PeperomiaBackend {{ .Mileston }}]({{ .BackendMilestoneURL }})
{{ end }}
{{ if .WebMilestoneURL }}[PeperomiaWeb {{ .Mileston }}]({{ .WebMilestoneURL }})
{{ end }}
{{ if .HelpMilestoneURL }}[PeperomiaHelp {{ .Mileston }}]({{ .HelpMilestoneURL }})
{{ end }}
{{ if .LPMilestoneURL }}[Peperomia LPサイト {{ .Mileston }}]({{ .LPMilestoneURL }})
{{ end }}
{{ if .ToolMilestoneURL }}[PeperomiaTool {{ .Mileston }}]({{ .ToolMilestoneURL }})
{{ end }}

# 対応内容
{{range $index, $pr := .App}}
- [[Peperomia]{{ $pr.Title }}](#{{ $pr.ID }}){{end}}
{{range $index, $pr := .Backend}}
- [[PeperomiaBackend]{{ $pr.Title }}](#{{ $pr.ID }}){{end}}
{{range $index, $pr := .Web}}
- [[PeperomiaWeb]{{ $pr.Title }}](#{{ $pr.ID }}){{end}}
{{range $index, $pr := .Help}}
- [[PeperomiaHelp]{{ $pr.Title }}](#{{ $pr.ID }}){{end}}
{{range $index, $pr := .LP}}
- [[PeperomiaWebSite]{{ $pr.Title }}](#{{ $pr.ID }}){{end}}
{{range $index, $pr := .Tool}}
- [[PeperomiaTool]{{ $pr.Title }}](#{{ $pr.ID }}){{end}}

# リリース詳細
{{range $index, $pr := .App}}<a id="{{ $pr.ID }}"></a>
## [[Peperomia]{{ $pr.Title }}]({{ $pr.URL }})
{{ $pr.Body}}
{{end}}
{{range $index, $pr := .Backend}}<a id="{{ $pr.ID }}"></a>
## [[PeperomiaBackend]{{ $pr.Title }}]({{ $pr.URL }})
{{ $pr.Body}}
{{end}}
{{range $index, $pr := .Web}}<a id="{{ $pr.ID }}"></a>
## [[PeperomiaWeb]{{ $pr.Title }}]({{ $pr.URL }})
{{ $pr.Body}}
{{end}}
{{range $index, $pr := .Help}}<a id="{{ $pr.ID }}"></a>
## [[PeperomiaHelp]{{ $pr.Title }}]({{ $pr.URL }})
{{ $pr.Body}}
{{end}}
{{range $index, $pr := .Tool}}<a id="{{ $pr.ID }}"></a>
## [[PeperomiaWebSite]{{ $pr.Title }}]({{ $pr.URL }})
{{ $pr.Body}}
{{end}}
{{range $index, $pr := .LP}}<a id="{{ $pr.ID }}"></a>
## [[PeperomiaTool]{{ $pr.Title }}]({{ $pr.URL }})
{{ $pr.Body}}
{{end}}
`

	tpl, err := template.New("").Parse(text)
	if err != nil {
		log.Fatal(err)
	}

	if err := tpl.Execute(os.Stdout, tmp); err != nil {
		log.Fatal(err)
	}
}

func (c *GitHubConfig) getPullRequestByMilestone(repository string) ([]PullRequest, string, error) {
	var ps []PullRequest

	mn, err := c.getMilestoneNumber(os.Getenv("MILESTONE"), repository)
	if err != nil {
		return ps, "", err
	}

	respData, err := c.getPullRequest(mn, repository)
	if err != nil {
		return ps, "", err
	}

	for i, pr := range respData.Repository.Milestone.PullRequests.Nodes {
		start := strings.Index(pr.Body, c.StartText)
		if start == -1 {
			start = 0
		} else {
			start += len(c.StartText)
		}
		end := strings.Index(pr.Body, c.EndText)
		if end == -1 {
			end = len(pr.Body)
		}

		body := pr.Body[start:end]
		respData.Repository.Milestone.PullRequests.Nodes[i].Body = body
	}

	return respData.Repository.Milestone.PullRequests.Nodes, respData.Repository.Milestone.URL, nil
}

func (c *GitHubConfig) getMilestoneNumber(mt string, repository string) (int, error) {
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
	req.Var("name", repository)

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

func (c *GitHubConfig) getPullRequest(mn int, repository string) (ResposeType, error) {
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
	req.Var("name", repository)
	req.Var("labels", [1]string{c.Label})
	req.Var("milestoneNumber", mn)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+c.Token)

	ctx := context.Background()

	var respData ResposeType
	err := client.Run(ctx, req, &respData)
	return respData, err
}
