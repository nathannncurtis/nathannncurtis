import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
GLYPHS = ROOT / "assets" / "glyphs"
BADGES = ROOT / "assets" / "stack"

GRAY = "#8b949e"
ICON = 20
GAP = 5
FONT_SIZE = 10
HEIGHT = ICON + GAP + 12
PAD = 2
FONT = "-apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif"

INTRO = "I build to fit the machine, not the cloud I might need someday. I'd rather open a closed system than work around it."

PROJECTS = """[ORCA](https://github.com/nathannncurtis/orca) :: a pattern for building on systems you don't own<br>
[Dokkit](https://github.com/nathannncurtis/dokkit) :: watches everything I run, flags what needs a look<br>
<sub>more at [nathancurtis.me/work](https://nathancurtis.me/work)</sub>"""

ROWS = [
    ("languages", [("rust", "Rust"), ("zig", "Zig"), ("c", "C"), ("go", "Go"), ("csharp", "C#"),
                   ("python", "Python"), ("elixir", "Elixir"), ("typescript", "TypeScript"),
                   ("swift", "Swift"), ("kotlin", "Kotlin")]),
    ("frameworks", [("dotnet", ".NET"), ("tauri", "Tauri"), ("fastapi", "FastAPI"),
                    ("phoenixframework", "Phoenix")]),
    ("ml", [("pytorch", "PyTorch"), ("llamacpp", "llama.cpp"), ("vllm", "vLLM"), ("onnx", "ONNX"),
            ("opencv", "OpenCV")]),
    ("infra", [("postgresql", "PostgreSQL"), ("sqlserver", "SQL Server"), ("sqlite", "SQLite"),
               ("redis", "Redis"), ("caddy", "Caddy"), ("wireguard", "WireGuard"),
               ("cloudflareworkers", "Cloudflare Workers")]),
    ("research", [("hxd", "HxD"), ("ghidra", "Ghidra"), ("hopper", "Hopper"), ("x64dbg", "x64dbg"),
                  ("aflplusplus", "AFL++"), ("burpsuite", "Burp Suite"), ("wireshark", "Wireshark")]),
]

WIDTHS = {**{c: 556 for c in "abdeghnopqu0123456789"}, **{c: 500 for c in "ckpsvxyz"},
          **{c: 222 for c in "ijl"}, "f": 278, "t": 278, "r": 333, "m": 833, "w": 722,
          **{c: 667 for c in "ABEKPSVXY"}, **{c: 722 for c in "CDHNRUw"}, "F": 611, "G": 778,
          "I": 278, "J": 500, "L": 556, "M": 833, "O": 778, "Q": 778, "T": 611, "W": 944, "Z": 611,
          " ": 278, ".": 278, "#": 556, "+": 584}


def text_width(s):
    return sum(WIDTHS.get(ch, 556) for ch in s) * FONT_SIZE / 1000


def esc(s):
    return s.replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;")


SPACING = 10


def badge(slug, name, left):
    src = (GLYPHS / f"{slug}.svg").read_text()
    root = re.search(r"<svg\b[^>]*>", src).group(0)
    view_box = re.search(r'viewBox="([^"]+)"', root).group(1)
    fill = re.search(r'\sfill="([^"]+)"', root)
    body = src[src.index(root) + len(root):src.rindex("</svg>")]
    body = re.sub(r"<title>.*?</title>", "", body, flags=re.S).strip()
    width = max(ICON, round(text_width(name))) + 2 * PAD
    x = left + (width - ICON) / 2
    fill_attr = f' fill="{fill.group(1)}"' if fill else ""
    return (
        f'<svg x="{x:g}" y="0" width="{ICON}" height="{ICON}" viewBox="{view_box}"{fill_attr}>{body}</svg>'
        f'<text x="{left + width / 2:g}" y="{HEIGHT - 2}">{esc(name)}</text>'
    ), width


def row_svg(items):
    parts, left = [], 0
    for slug, name in items:
        svg, width = badge(slug, name, left)
        parts.append(svg)
        left += width + SPACING
    total = left - SPACING
    names = ", ".join(name for _, name in items)
    return (
        f'<svg xmlns="http://www.w3.org/2000/svg" width="{total}" height="{HEIGHT}" '
        f'viewBox="0 0 {total} {HEIGHT}" role="img" aria-label="{esc(names)}"><title>{esc(names)}</title>'
        f'<g text-anchor="middle" font-family="{FONT}" font-size="{FONT_SIZE}" fill="{GRAY}">'
        + "".join(parts) + "</g></svg>\n"
    ), total


def main():
    BADGES.mkdir(parents=True, exist_ok=True)
    for old in BADGES.glob("*.svg"):
        old.unlink()
    parts = [INTRO, "", PROJECTS]
    for label, items in ROWS:
        svg, width = row_svg(items)
        (BADGES / f"{label}.svg").write_text(svg)
        names = ", ".join(name for _, name in items)
        parts += ["", f"<sub>{label}</sub><br>",
                  f'<img src="assets/stack/{label}.svg" width="{width}" height="{HEIGHT}" alt="{esc(names)}">']
    (ROOT / "README.md").write_text("\n".join(parts) + "\n")


if __name__ == "__main__":
    main()
