package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
)

// public repos shown as rows; private and org work only appears inside the "everything" total
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
	Text  string // project names, numbers
	Muted string // last push, header labels
	Faint string // eyebrow, language
	Rule  string // header rule
	Line  string // row rules
	Hot   string // "everything" row
}

var dark = palette{
	Text: "rgba(230,237,243,0.92)", Muted: "rgba(201,214,227,0.62)", Faint: "rgba(201,214,227,0.38)",
	Rule: "rgba(201,214,227,0.18)", Line: "rgba(201,214,227,0.08)", Hot: "#E3A867",
}

var light = palette{
	Text: "#1f2328", Muted: "#59636e", Faint: "#8c959f",
	Rule: "#d1d9e0", Line: "#eff2f5", Hot: "#B8762F",
}

type project struct {
	Name   string
	Lang   string
	Color  string
	Count  int
	Weeks  []int
	Pushed time.Time
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

// gql posts a GraphQL query and decodes data into out.
func gql(token, q string, vars map[string]any, out any) {
	body, _ := json.Marshal(map[string]any{"query": q, "variables": vars})
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
	var env struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(raw, &env); err != nil || len(env.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "GraphQL error: %v %s\n", err, raw)
		os.Exit(1)
	}
	if err := json.Unmarshal(env.Data, out); err != nil {
		fmt.Fprintln(os.Stderr, "GraphQL decode:", err)
		os.Exit(1)
	}
}

const calendarQuery = `query($login: String!) {
  user(login: $login) { contributionsCollection { contributionCalendar { weeks { contributionDays { date contributionCount } } } } }
}`

// per-repo daily commit counts; queried in windows short enough that no repo has more than 100 active days
const commitsQuery = `query($login: String!, $from: DateTime!, $to: DateTime!) {
  user(login: $login) { contributionsCollection(from: $from, to: $to) {
    commitContributionsByRepository(maxRepositories: 100) {
      repository { name isPrivate pushedAt owner { login } primaryLanguage { name } }
      contributions(first: 100) { nodes { occurredAt commitCount } }
    }
  } }
}`

type calendarResp struct {
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
		} `json:"contributionsCollection"`
	} `json:"user"`
}

type commitsResp struct {
	User struct {
		ContributionsCollection struct {
			CommitContributionsByRepository []struct {
				Repository struct {
					Name      string    `json:"name"`
					IsPrivate bool      `json:"isPrivate"`
					PushedAt  time.Time `json:"pushedAt"`
					Owner     struct {
						Login string `json:"login"`
					} `json:"owner"`
					PrimaryLanguage *struct {
						Name string `json:"name"`
					} `json:"primaryLanguage"`
				} `json:"repository"`
				Contributions struct {
					Nodes []struct {
						OccurredAt  time.Time `json:"occurredAt"`
						CommitCount int       `json:"commitCount"`
					} `json:"nodes"`
				} `json:"contributions"`
			} `json:"commitContributionsByRepository"`
		} `json:"contributionsCollection"`
	} `json:"user"`
}

func comma(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

func ago(t, now time.Time) string {
	d := int(now.Sub(t).Hours() / 24)
	switch {
	case d <= 0:
		return "today"
	case d == 1:
		return "yesterday"
	default:
		return fmt.Sprintf("%dd ago", d)
	}
}

// sparkline draws a filled line of weekly values, scaled to its own maximum.
func sparkline(vals []int, x0, w, base, h float64, color string) string {
	mx := 1
	for _, v := range vals {
		if v > mx {
			mx = v
		}
	}
	pts := []string{fmt.Sprintf("%.1f,%.1f", x0, base)}
	for k, v := range vals {
		x := x0 + (float64(k)+0.5)/float64(len(vals))*w
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", x, base-float64(v)/float64(mx)*h))
	}
	pts = append(pts, fmt.Sprintf("%.1f,%.1f", x0+w, base))
	d := "M " + strings.Join(pts, " L ")
	return fmt.Sprintf(`  <path d="%s Z" fill="%s" fill-opacity="0.22"/>`+"\n"+`  <path d="%s" fill="none" stroke="%s" stroke-width="1.3" stroke-linejoin="round"/>`+"\n", d, color, d, color)
}

func render(p palette, total int, weekTotals []int, lastActive string, projects []project, now time.Time) string {
	var b strings.Builder
	w := func(format string, a ...any) { fmt.Fprintf(&b, format+"\n", a...) }
	mono := "ui-monospace,SFMono-Regular,Menlo,Consolas,monospace"
	sans := "system-ui,-apple-system,'Segoe UI',Helvetica,Arial,sans-serif"
	const (
		sx, sw   = 300.0, 330.0 // sparkline column
		colCount = 720
		hdrY     = 40
		rowH     = 32
	)
	height := hdrY + 30 + 12 + len(projects)*rowH + 8

	w(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, svgWidth, height, svgWidth, height)
	w(`  <text x="0" y="14" fill="%s" font-family="%s" font-size="10" letter-spacing="2">LAST 12 MONTHS</text>`, p.Faint, mono)
	for _, h := range []struct {
		label  string
		x      float64
		anchor string
	}{{"PROJECT", 0, "start"}, {"LAST 52 WEEKS", sx, "start"}, {"COMMITS", colCount, "end"}, {"LAST PUSH", svgWidth, "end"}} {
		w(`  <text x="%.0f" y="%d" fill="%s" font-family="%s" font-size="9" letter-spacing="1.5" text-anchor="%s">%s</text>`, h.x, hdrY, p.Faint, mono, h.anchor, h.label)
	}
	w(`  <line x1="0" y1="%d" x2="%d" y2="%d" stroke="%s"/>`, hdrY+8, svgWidth, hdrY+8, p.Rule)

	// everything row: the full contribution calendar, private work included as counts
	y := hdrY + 30
	w(`  <text x="0" y="%d" fill="%s" font-family="%s" font-size="13" font-weight="600">everything</text>`, y, p.Text, sans)
	w(`  <text x="110" y="%d" fill="%s" font-family="%s" font-size="10">incl. private</text>`, y, p.Faint, mono)
	b.WriteString(sparkline(weekTotals, sx, sw, float64(y+2), 18, p.Hot))
	w(`  <text x="%d" y="%d" fill="%s" font-family="%s" font-size="12" text-anchor="end">%s</text>`, colCount, y, p.Hot, mono, comma(total))
	w(`  <text x="%d" y="%d" fill="%s" font-family="%s" font-size="11" text-anchor="end">%s</text>`, svgWidth, y, p.Muted, mono, lastActive)
	w(`  <line x1="0" y1="%d" x2="%d" y2="%d" stroke="%s"/>`, y+12, svgWidth, y+12, p.Line)

	for r, pr := range projects {
		yy := y + 34 + r*rowH
		w(`  <circle cx="5" cy="%d" r="4" fill="%s"/>`, yy-4, pr.Color)
		w(`  <text x="16" y="%d" fill="%s" font-family="%s" font-size="13">%s</text>`, yy, p.Text, sans, pr.Name)
		w(`  <text x="150" y="%d" fill="%s" font-family="%s" font-size="10">%s</text>`, yy, p.Faint, mono, pr.Lang)
		b.WriteString(sparkline(pr.Weeks, sx, sw, float64(yy+2), 16, pr.Color))
		w(`  <text x="%d" y="%d" fill="%s" font-family="%s" font-size="12" text-anchor="end">%s</text>`, colCount, yy, p.Text, mono, comma(pr.Count))
		w(`  <text x="%d" y="%d" fill="%s" font-family="%s" font-size="11" text-anchor="end">%s</text>`, svgWidth, yy, p.Muted, mono, ago(pr.Pushed, now))
		w(`  <line x1="0" y1="%d" x2="%d" y2="%d" stroke="%s"/>`, yy+12, svgWidth, yy+12, p.Line)
	}

	w(`</svg>`)
	return b.String()
}

func main() {
	token := getToken()
	now := time.Now().UTC()
	vars := map[string]any{"login": username}
	fmt.Printf("Fetching contributions for %s...\n", username)

	// last 52 full calendar weeks (the calendar starts with a partial week)
	var cal calendarResp
	gql(token, calendarQuery, vars, &cal)
	calWeeks := cal.User.ContributionsCollection.ContributionCalendar.Weeks
	if len(calWeeks) > weeks {
		calWeeks = calWeeks[len(calWeeks)-weeks:]
	}
	start, _ := time.Parse("2006-01-02", calWeeks[0].ContributionDays[0].Date)
	var weekTotals []int
	total, lastDay := 0, ""
	for _, wk := range calWeeks {
		sum := 0
		for _, d := range wk.ContributionDays {
			sum += d.ContributionCount
			if d.ContributionCount > 0 {
				lastDay = d.Date
			}
		}
		weekTotals = append(weekTotals, sum)
		total += sum
	}

	// featured public repos only; private and org repos are never named or counted here
	byName := map[string]*project{}
	for from := start; from.Before(now); from = from.AddDate(0, 0, 91) {
		to := from.AddDate(0, 0, 91)
		if to.After(now) {
			to = now
		}
		var cr commitsResp
		gql(token, commitsQuery, map[string]any{"login": username, "from": from.Format(time.RFC3339), "to": to.Format(time.RFC3339)}, &cr)
		for _, r := range cr.User.ContributionsCollection.CommitContributionsByRepository {
			repo := r.Repository
			if repo.IsPrivate || repo.Owner.Login != username || !featured[repo.Name] {
				continue
			}
			pr := byName[repo.Name]
			if pr == nil {
				lang, color := "", defaultLangColor
				if repo.PrimaryLanguage != nil {
					lang = repo.PrimaryLanguage.Name
					if c, ok := langColors[lang]; ok {
						color = c
					}
				}
				pr = &project{Name: repo.Name, Lang: lang, Color: color, Weeks: make([]int, len(weekTotals)), Pushed: repo.PushedAt}
				byName[repo.Name] = pr
			}
			for _, n := range r.Contributions.Nodes {
				k := int(n.OccurredAt.Sub(start).Hours() / 24 / 7)
				if k >= 0 && k < len(pr.Weeks) {
					pr.Weeks[k] += n.CommitCount
					pr.Count += n.CommitCount
				}
			}
		}
	}

	var projects []project
	for _, pr := range byName {
		if pr.Count > 0 {
			projects = append(projects, *pr)
		}
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].Pushed.After(projects[j].Pushed) })
	// the everything row's last push is the newest of the calendar and the featured repos
	latest, _ := time.Parse("2006-01-02", lastDay)
	for _, pr := range projects {
		if pr.Pushed.After(latest) {
			latest = pr.Pushed
		}
	}
	lastActive := ago(latest, now)
	fmt.Printf("  %d contributions over %d weeks; %d featured repos active\n", total, len(weekTotals), len(projects))

	os.MkdirAll("assets", 0755)
	for _, f := range []struct {
		path string
		p    palette
	}{{"timeline.svg", dark}, {"timeline-light.svg", light}} {
		out := filepath.Join("assets", f.path)
		os.WriteFile(out, []byte(render(f.p, total, weekTotals, lastActive, projects, now)), 0644)
		fmt.Printf("  Written to %s\n", out)
	}
}
