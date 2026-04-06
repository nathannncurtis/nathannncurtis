package main

import (
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
	username   = "nathannncurtis"
	monthsBack = 6
	svgWidth   = 820
	rowH       = 28
	topPad     = 52
	bottomPad  = 24
	tlLeft     = 160
	rightPad   = 28
)

var orgs = []string{"ronsin-lss"}

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
	"PHP":        "#6A7A8A",
}

const defaultColor = "#556270"
const orgColor = "#3A4450"

type ghRepo struct {
	Name     string `json:"name"`
	Private  bool   `json:"private"`
	Archived bool   `json:"archived"`
	PushedAt string `json:"pushed_at"`
	Language *string `json:"language"`
}

type entry struct {
	Label    string
	Lang     string
	Pushed   time.Time
	Color    string
	IsOrg    bool
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

func apiGet(path, token string) ([]byte, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com"+path, nil)
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "timeline-generator")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func fetchRepos(token string) []ghRepo {
	var all []ghRepo

	// personal repos
	for page := 1; ; page++ {
		data, err := apiGet(fmt.Sprintf("/user/repos?type=owner&per_page=100&page=%d", page), token)
		if err != nil {
			break
		}
		var batch []ghRepo
		if json.Unmarshal(data, &batch) != nil || len(batch) == 0 {
			break
		}
		all = append(all, batch...)
	}

	// org repos
	for _, org := range orgs {
		for page := 1; ; page++ {
			data, err := apiGet(fmt.Sprintf("/orgs/%s/repos?per_page=100&page=%d", org, page), token)
			if err != nil {
				break
			}
			var batch []ghRepo
			if json.Unmarshal(data, &batch) != nil || len(batch) == 0 {
				break
			}
			for i := range batch {
				batch[i].Name = "org:" + batch[i].Name
			}
			all = append(all, batch...)
		}
	}

	return all
}

func buildEntries(repos []ghRepo, cutoff time.Time) (personal []entry, orgEntries []entry) {
	for _, r := range repos {
		if r.Archived || r.PushedAt == "" {
			continue
		}
		pushed, err := time.Parse(time.RFC3339, r.PushedAt)
		if err != nil {
			continue
		}
		if pushed.Before(cutoff) {
			continue
		}

		lang := ""
		if r.Language != nil {
			lang = *r.Language
		}
		color := defaultColor
		if c, ok := langColors[lang]; ok {
			color = c
		}

		isOrg := strings.HasPrefix(r.Name, "org:")
		name := strings.TrimPrefix(r.Name, "org:")

		e := entry{
			Label:  name,
			Lang:   lang,
			Pushed: pushed,
			Color:  color,
			IsOrg:  isOrg,
		}

		if isOrg {
			e.Color = orgColor
			orgEntries = append(orgEntries, e)
		} else {
			personal = append(personal, e)
		}
	}

	sort.Slice(personal, func(i, j int) bool {
		return personal[i].Pushed.After(personal[j].Pushed)
	})
	sort.Slice(orgEntries, func(i, j int) bool {
		return orgEntries[i].Pushed.After(orgEntries[j].Pushed)
	})

	return
}

func xPos(dt, cutoff time.Time, totalDays float64) float64 {
	days := dt.Sub(cutoff).Hours() / 24
	frac := days / totalDays
	return float64(tlLeft) + frac*float64(svgWidth-tlLeft-rightPad)
}

func generateSVG(personal, orgEntries []entry, cutoff, now time.Time) string {
	rowCount := len(personal)
	if len(orgEntries) > 0 {
		rowCount++
	}
	height := topPad + rowCount*rowH + bottomPad
	totalDays := now.Sub(cutoff).Hours() / 24
	if totalDays < 1 {
		totalDays = 1
	}

	var b strings.Builder
	w := func(s string) { b.WriteString(s + "\n") }

	w(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, svgWidth, height, svgWidth, height))
	w(fmt.Sprintf(`  <rect width="%d" height="%d" rx="10" fill="#0d1117"/>`, svgWidth, height))
	w(fmt.Sprintf(`  <rect x="0" y="0" width="%d" height="1" rx="0.5" fill="#4A9A9A" opacity="0.3"/>`, svgWidth))
	w(`  <text x="28" y="28" fill="#556270" font-family="monospace" font-size="10" letter-spacing="2">RECENT ACTIVITY</text>`)

	// month gridlines
	d := time.Date(cutoff.Year(), cutoff.Month(), 1, 0, 0, 0, 0, time.UTC)
	for d.Before(now) || d.Equal(now) {
		if !d.Before(cutoff) {
			mx := xPos(d, cutoff, totalDays)
			label := strings.ToUpper(d.Format("Jan"))
			w(fmt.Sprintf(`  <line x1="%.0f" y1="%d" x2="%.0f" y2="%d" stroke="#1a1f26" stroke-width="1"/>`, mx, topPad-8, mx, height-bottomPad))
			w(fmt.Sprintf(`  <text x="%.0f" y="%d" fill="#3a4450" font-family="monospace" font-size="9" text-anchor="middle">%s</text>`, mx, topPad-14, label))
		}
		if d.Month() == 12 {
			d = time.Date(d.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)
		} else {
			d = time.Date(d.Year(), d.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		}
	}

	// today marker
	tx := xPos(now, cutoff, totalDays)
	w(fmt.Sprintf(`  <line x1="%.0f" y1="%d" x2="%.0f" y2="%d" stroke="#4A9A9A" stroke-width="1" opacity="0.3"/>`, tx, topPad-8, tx, height-bottomPad))

	// personal rows
	for i, e := range personal {
		y := topPad + i*rowH + rowH/2
		ex := xPos(e.Pushed, cutoff, totalDays)

		label := e.Label
		if len(label) > 18 {
			label = label[:17] + "…"
		}

		w(fmt.Sprintf(`  <text x="%d" y="%d" fill="#7a8a96" font-family="'Segoe UI',system-ui,sans-serif" font-size="12" text-anchor="end">%s</text>`, tlLeft-12, y+4, label))
		w(fmt.Sprintf(`  <line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#151a20" stroke-width="1"/>`, tlLeft, y, svgWidth-rightPad, y))
		w(fmt.Sprintf(`  <circle cx="%.0f" cy="%d" r="5" fill="%s" opacity="0.9"/>`, ex, y, e.Color))

		if e.Lang != "" {
			w(fmt.Sprintf(`  <text x="%.0f" y="%d" fill="#3a4450" font-family="monospace" font-size="9">%s</text>`, ex+10, y+4, e.Lang))
		}
	}

	// consolidated org row
	if len(orgEntries) > 0 {
		i := len(personal)
		y := topPad + i*rowH + rowH/2
		label := fmt.Sprintf("%d private", len(orgEntries))

		w(fmt.Sprintf(`  <text x="%d" y="%d" fill="#3a4450" font-family="'Segoe UI',system-ui,sans-serif" font-size="12" text-anchor="end" font-style="italic">%s</text>`, tlLeft-12, y+4, label))
		w(fmt.Sprintf(`  <line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#151a20" stroke-width="1"/>`, tlLeft, y, svgWidth-rightPad, y))

		for _, oe := range orgEntries {
			ox := xPos(oe.Pushed, cutoff, totalDays)
			w(fmt.Sprintf(`  <circle cx="%.0f" cy="%d" r="3" fill="%s" opacity="0.7"/>`, ox, y, orgColor))
		}
	}

	w(`</svg>`)
	return b.String()
}

func main() {
	token := getToken()

	now := time.Now().UTC()
	cutoff := now.AddDate(0, -monthsBack, 0)

	fmt.Printf("Fetching repos for %s + orgs %v...\n", username, orgs)
	repos := fetchRepos(token)
	fmt.Printf("  Found %d repos total\n", len(repos))

	personal, orgEntries := buildEntries(repos, cutoff)
	fmt.Printf("  %d personal + %d org active in last %d months\n", len(personal), len(orgEntries), monthsBack)

	svg := generateSVG(personal, orgEntries, cutoff, now)

	scriptDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	outPath := filepath.Join(scriptDir, "..", "assets", "timeline.svg")
	// fallback if running via `go run`
	if _, err := os.Stat(filepath.Dir(outPath)); os.IsNotExist(err) {
		outPath = filepath.Join("assets", "timeline.svg")
	}
	os.MkdirAll(filepath.Dir(outPath), 0755)
	os.WriteFile(outPath, []byte(svg), 0644)
	fmt.Printf("  Written to %s\n", outPath)
}
