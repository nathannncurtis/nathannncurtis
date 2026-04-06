# nathan 👋

I build internal tools. The kind that connect systems no one designed to talk to each other, replace the manual glue holding production workflows together, and survive environments where "just use a SaaS" isn't an option. I'm usually the person diagnosing why something broke before I'm the one building the fix.

Most of my work lives in private repos and under my org. Here's what I can talk about.

---

### right now

**[mdview-zig](https://github.com/nathannncurtis/mdview-zig)** — a fast, native markdown viewer in Zig. Platform-native text rendering (DirectWrite on Windows, Cairo/Pango on Linux, CoreText on macOS), no webview, no runtime — ~285KB standalone binary across all three platforms. Rewrite of an earlier Rust/WebView2 version that dropped the 200MB browser dependency entirely.

**[Study Aggregator v4.0](https://github.com/nathannncurtis/Study-Aggregator)** — DICOM processor for legal/medical imaging workflows. Rewrote the parsing hot path as a native Rust engine called via subprocess from the PyQt5 GUI. Zero-copy memory-mapped parsing, parallel directory walking and per-file parsing (rayon), streaming ZIP extraction so multi-GB archives don't balloon RAM. Roughly **100–400x faster** than the previous pydicom-based implementation on typical office hardware.

**[commit-summarizer](https://github.com/nathannncurtis/commit-summarizer)** — webhook service that verifies GitHub push events, runs commit metadata through a local Ollama model, and posts plain-English summaries to Slack. Entirely on-prem, no data egress to third-party LLM APIs.

**[obsidian-vault-sync](https://github.com/nathannncurtis/obsidian-vault-sync)** — self-hosted real-time Obsidian sync. FastAPI server + TypeScript plugin, WebSocket push with polling fallback, token auth, Dockerized. Keeps vaults in sync across devices with no third-party service in the loop.

**[steddi](https://steddi.io)** — iOS navigation app for daily commuters who already know their routes. Swift, SwiftUI, MapKit, CarPlay support. Only reroutes when it actually matters. Currently in development.

---

### the kind of stuff I build

Document processing pipelines for legal and medical workflows: batch OCR, PDF compression, image classification. Built around the actual hardware, network paths, and legacy software they have to survive in.

DICOM viewers and x-ray processing tools that handle the gap between clinical systems and what people actually need to do with the data.

Dashboards, operational automation, and small utilities. The stuff that saves 20 minutes a day and no one thinks to automate.

---

### tech

![Python](https://img.shields.io/badge/Python-3776AB?style=flat&logo=python&logoColor=white) ![Swift](https://img.shields.io/badge/Swift-F05138?style=flat&logo=swift&logoColor=white) ![Rust](https://img.shields.io/badge/Rust-000000?style=flat&logo=rust&logoColor=white) ![Zig](https://img.shields.io/badge/Zig-F7A41D?style=flat&logo=zig&logoColor=white) ![C](https://img.shields.io/badge/C-A8B9CC?style=flat&logo=c&logoColor=black) ![C++](https://img.shields.io/badge/C++-00599C?style=flat&logo=cplusplus&logoColor=white) ![C#](https://img.shields.io/badge/C%23-239120?style=flat&logo=csharp&logoColor=white) ![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=flat&logo=typescript&logoColor=white) ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white) ![PHP](https://img.shields.io/badge/PHP-777BB4?style=flat&logo=php&logoColor=white) ![PowerShell](https://img.shields.io/badge/PowerShell-5391FE?style=flat&logo=powershell&logoColor=white)

![React](https://img.shields.io/badge/React-61DAFB?style=flat&logo=react&logoColor=black) ![SwiftUI](https://img.shields.io/badge/SwiftUI-F05138?style=flat&logo=swift&logoColor=white) ![Electron](https://img.shields.io/badge/Electron-47848F?style=flat&logo=electron&logoColor=white) ![FastAPI](https://img.shields.io/badge/FastAPI-009688?style=flat&logo=fastapi&logoColor=white) ![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white) ![SQLite](https://img.shields.io/badge/SQLite-003B57?style=flat&logo=sqlite&logoColor=white) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white)

---

more at [nathancurtis.to](https://nathancurtis.to)
