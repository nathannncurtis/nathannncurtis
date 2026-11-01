package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	username = "nathannncurtis"
	svgWidth = 820
	weeks    = 52
	hotWeeks = 4 // trailing weeks drawn in the accent colour
)

// public repos shown by name; every other contribution (including private work) is "everything else"
var featured = map[string]bool{
	"orca":            true,
	"dokkit":          true,
	"steddi":          true,
	"coil":            true,
	"mdview-zig":      true,
	"nathancurtis.me": true,
}

var langColors = map[string]string{
	"Python":     "#4A9A9A",
	"Rust":       "#C97B5E",
	"Zig":        "#E9A14B",
	"TypeScript": "#7DB4E8",
	"JavaScript": "#7DB4E8",
	"Swift":      "#A78BFA",
	"C++":        "#6A7A8A",
	"C":          "#6A7A8A",
	"C#":         "#6A7A8A",
	"Go":         "#5CB8B8",
	"Elixir":     "#C084B0",
	"PHP":        "#6A7A8A",
}

const defaultLangColor = "#8B949E"

// palette holds the colours baked into one SVG file. Backgrounds are never painted,
// so each file sits on whatever page colour the reader's GitHub theme uses.
type palette struct {
	Faint  string // eyebrow, month labels
	Muted  string // legend text, totals
	Rest   string // "everything else" segment
	RestOp string
	Year   string // weekly bars
	YearOp string
	Hot    string // trailing weeks
	Base   string // zero-week stub
}

var dark = palette{
	Faint: "rgba(201,214,227,0.38)", Muted: "rgba(201,214,227,0.66)",
	Rest: "rgba(201,214,227,0.36)", RestOp: "0.6",
	Year: "#4FB0A8", YearOp: "0.7", Hot: "#E3A867", Base: "rgba(201,214,227,0.18)",
}

var light = palette{
	Faint: "#8c959f", Muted: "#59636e",
	Rest: "#8c959f", RestOp: "0.6",
	Year: "#2E8A83", YearOp: "0.7", Hot: "#B8762F", Base: "#d1d9e0",
}

type project struct {
	Name  string
	Color string
	Count int
}

func getToken() string {
	if t := os.Getenv("GITHUB_TOKEN"); t != "" {
		return t
	}
	out, err := exec.Command("gh", "auth", "token").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	fmt.Fprintln(os.Stderr, "No GitHub token found. Set GITHUB_TOKEN or install gh CLI.")
	os.Exit(1)
	return ""
}

const query = `query($login: String!) {
  user(login: $login) {
    contributionsCollection {
      contributionCalendar { weeks { contributionDays { date contributionCount } } }
      commitContributionsByRepository(maxRepositories: 100) {
        repository { name isPrivate owner { login } primaryLanguage { name } }
        contributions { totalCount }
      }
    }
  }
}`

type gqlResp struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					Weeks []struct {
						ContributionDays []struct {
							Date              string `json:"date"`
							ContributionCount int    `json:"contributionCount"`
						} `json:"contributionDays"`
					} `json:"weeks"`
				} `json:"contributionCalendar"`
				CommitContributionsByRepository []struct {
					Repository struct {
						Name      string `json:"name"`
						IsPrivate bool   `json:"isPrivate"`
						Owner     struct {
							Login string `json:"login"`
						} `json:"owner"`
						PrimaryLanguage *struct {
							Name string `json:"name"`
						} `json:"primaryLanguage"`
					} `json:"repository"`
					Contributions struct {
						TotalCount int `json:"totalCount"`
					} `json:"contributions"`
				} `json:"commitContributionsByRepository"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func fetch(token string) gqlResp {
	body, _ := json.Marshal(map[string]any{"query": query, "variables": map[string]string{"login": username}})
	req, _ := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(body))
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("User-Agent", "timeline-generator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "GraphQL request failed:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var out gqlResp
	if err := json.Unmarshal(raw, &out); err != nil || len(out.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "GraphQL error: %v %s\n", err, raw)
		os.Exit(1)
	}
	return out
}

func comma(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

func render(p palette, total int, projects []project, rest int, weekTotals []int, weekStarts []time.Time) string {
	var b strings.Builder
	w := func(format string, a ...any) { fmt.Fprintf(&b, format+"\n", a...) }
	const height = 200
	mono := "ui-monospace,SFMono-Regular,Menlo,Consolas,monospace"
	sans := "system-ui,-apple-system,'Segoe UI',Helvetica,Arial,sans-serif"

	w(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, svgWidth, height, svgWidth, height)
	w(`  <text x="0" y="14" fill="%s" font-family="%s" font-size="10" letter-spacing="2">WHERE THE LAST 12 MONTHS WENT</text>`, p.Faint, mono)
	w(`  <text x="%d" y="14" fill="%s" font-family="%s" font-size="11" text-anchor="end">%s contributions</text>`, svgWidth, p.Muted, mono, comma(total))

	type part struct {
		label, color, op string
		count            int
	}
	var parts []part
	for _, pr := range projects {
		parts = append(parts, part{pr.Name, pr.Color, "", pr.Count})
	}
	parts = append(parts, part{"everything else, mostly private", p.Rest, p.RestOp, rest})

	opAttr := func(op string) string {
		if op == "" {
			return ""
		}
		return fmt.Sprintf(` opacity="%s"`, op)
	}

	// split bar
	x := 0.0
	for _, pt := range parts {
		pw := float64(pt.count) / float64(total) * svgWidth
		w(`  <rect x="%.1f" y="30" width="%.1f" height="26" rx="2" fill="%s"%s/>`, x, math.Max(1.5, pw-2), pt.color, opAttr(pt.op))
		x += pw
	}

	// legend
	lx := 0.0
	for _, pt := range parts {
		pct := float64(pt.count) / float64(total) * 100
		ps := fmt.Sprintf("%.0f%%", pct)
		if pct < 1 {
			ps = "&lt;1%"
		}
		label := pt.label + " " + ps
		w(`  <rect x="%.1f" y="70" width="9" height="9" rx="2" fill="%s"%s/>`, lx, pt.color, opAttr(pt.op))
		w(`  <text x="%.1f" y="78" fill="%s" font-family="%s" font-size="11">%s</text>`, lx+14, p.Muted, sans, label)
		lx += 30 + float64(len(strings.ReplaceAll(label, "&lt;", "<")))*6.2
	}

	// weekly volume strip
	maxW := 1
	for _, t := range weekTotals {
		if t > maxW {
			maxW = t
		}
	}
	bw := float64(svgWidth) / float64(len(weekTotals))
	for k, t := range weekTotals {
		h := math.Max(1, float64(t)/float64(maxW)*60)
		color, op := p.Year, p.YearOp
		if t == 0 {
			color, op = p.Base, ""
		} else if k >= len(weekTotals)-hotWeeks {
			color, op = p.Hot, ""
		}
		w(`  <rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="1.5" fill="%s"%s/>`, float64(k)*bw+1, 170-h, bw-2, h, color, opAttr(op))
	}

	// month labels at the first week containing the 1st
	for k, ws := range weekStarts {
		for d := 0; d < 7; d++ {
			day := ws.AddDate(0, 0, d)
			if day.Day() == 1 {
				mx := (float64(k) + float64(d)/7) * bw
				w(`  <text x="%.1f" y="188" fill="%s" font-family="%s" font-size="10" letter-spacing="2">%s</text>`, mx, p.Faint, mono, strings.ToUpper(day.Format("Jan")))
			}
		}
	}

	w(`</svg>`)
	return b.String()
}

func main() {
	token := getToken()
	fmt.Printf("Fetching contributions for %s...\n", username)
	cc := fetch(token).Data.User.ContributionsCollection

	// last 52 full calendar weeks (the calendar starts with a partial week)
	calWeeks := cc.ContributionCalendar.Weeks
	if len(calWeeks) > weeks {
		calWeeks = calWeeks[len(calWeeks)-weeks:]
	}
	var weekTotals []int
	var weekStarts []time.Time
	total := 0
	for _, wk := range calWeeks {
		sum := 0
		for _, d := range wk.ContributionDays {
			sum += d.ContributionCount
		}
		start, _ := time.Parse("2006-01-02", wk.ContributionDays[0].Date)
		weekTotals = append(weekTotals, sum)
		weekStarts = append(weekStarts, start)
		total += sum
	}

	// featured public repos only; private and org repos are never named
	var projects []project
	named := 0
	for _, r := range cc.CommitContributionsByRepository {
		repo := r.Repository
		if repo.IsPrivate || repo.Owner.Login != username || !featured[repo.Name] || r.Contributions.TotalCount == 0 {
			continue
		}
		color := defaultLangColor
		if repo.PrimaryLanguage != nil {
			if c, ok := langColors[repo.PrimaryLanguage.Name]; ok {
				color = c
			}
		}
		projects = append(projects, project{repo.Name, color, r.Contributions.TotalCount})
		named += r.Contributions.TotalCount
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].Count > projects[j].Count })
	if total < named {
		total = named
	}
	rest := total - named
	fmt.Printf("  %d contributions over %d weeks; %d in %d featured repos\n", total, len(weekTotals), named, len(projects))

	os.MkdirAll("assets", 0755)
	for _, f := range []struct {
		path string
		p    palette
	}{{"timeline.svg", dark}, {"timeline-light.svg", light}} {
		out := filepath.Join("assets", f.path)
		os.WriteFile(out, []byte(render(f.p, total, projects, rest, weekTotals, weekStarts)), 0644)
		fmt.Printf("  Written to %s\n", out)
	}
}
