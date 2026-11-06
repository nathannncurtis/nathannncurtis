"""Redraws the year line (assets/year-*.svg) from the last twelve months of contributions.

Run daily by .github/workflows/timeline.yml. Needs GITHUB_TOKEN with access to private
contribution counts. Only the line and its end dot change; fonts, colors and timing stay
as they are in the committed SVGs.
"""
import json, math, os, re, urllib.request, pathlib

QUERY = '{ viewer { contributionsCollection { contributionCalendar { weeks { contributionDays { contributionCount } } } } } }'

def days():
    req = urllib.request.Request("https://api.github.com/graphql", data=json.dumps({"query": QUERY}).encode(),
                                 headers={"Authorization": f"bearer {os.environ['GITHUB_TOKEN']}", "Content-Type": "application/json"})
    weeks = json.load(urllib.request.urlopen(req))["data"]["viewer"]["contributionsCollection"]["contributionCalendar"]["weeks"]
    return [d["contributionCount"] for w in weeks for d in w["contributionDays"]]

def path(counts, width):
    k = 9; sm = []
    for i in range(len(counts)):
        js = [j for j in range(-k, k + 1) if 0 <= i + j < len(counts)]
        sm.append(sum(counts[i + j] * (k + 1 - abs(j)) for j in js) / sum(k + 1 - abs(j) for j in js))
    vals = [math.sqrt(v) for v in sm[::2]]; top = max(vals) or 1
    base, amp, x0, x1 = 52, 46, 4, width - 4
    pts = [(x0 + (x1 - x0) * i / (len(vals) - 1), base - amp * v / top) for i, v in enumerate(vals)]
    d = f"M{pts[0][0]:.1f} {pts[0][1]:.1f}"
    for i in range(len(pts) - 1):
        p0 = pts[i - 1] if i else pts[i]; p1, p2 = pts[i], pts[i + 1]; p3 = pts[i + 2] if i + 2 < len(pts) else p2
        c1 = (p1[0] + (p2[0] - p0[0]) / 6, p1[1] + (p2[1] - p0[1]) / 6); c2 = (p2[0] - (p3[0] - p1[0]) / 6, p2[1] - (p3[1] - p1[1]) / 6)
        d += f" C{c1[0]:.1f} {c1[1]:.1f} {c2[0]:.1f} {c2[1]:.1f} {p2[0]:.1f} {p2[1]:.1f}"
    return d, pts[-1]

counts = days()
for f in pathlib.Path("assets").glob("year-*.svg"):
    svg = f.read_text()
    width = float(re.search(r'<svg[^>]*\swidth="([\d.]+)"', svg).group(1))
    d, (ex, ey) = path(counts, width)
    svg = re.sub(r'(<path class="yl" d=")[^"]*', lambda m: m.group(1) + d, svg)
    svg = re.sub(r'(<circle class="yd" cx=")[^"]*(" cy=")[^"]*', lambda m: f'{m.group(1)}{ex:.1f}{m.group(2)}{ey:.1f}', svg)
    f.write_text(svg)
