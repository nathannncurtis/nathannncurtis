# nathan

I build internal tools — automation, document processing, desktop apps, and the infrastructure around them. Most of my work lives in private repos for environments where nothing talks to anything else and "just use a SaaS" isn't an option.

I'm usually the person diagnosing why a workflow broke before I'm the one building the thing that fixes it.

---

### what I work with

**languages** — Python, TypeScript, JavaScript, Go, PowerShell, PHP

**frontend** — React, Electron, HTML/CSS

**backend** — FastAPI, PyQt

**data** — SQLite, PostgreSQL

**tooling** — Tesseract, Ghostscript, ImageMagick, Pillow, Docker

**environment** — Windows, WSL, Ubuntu, macOS, Arch, VS Code, Git

---

### what I've built

**document processing** — batch OCR, PDF compression, image pipelines, and classification tools for legal/medical workflows. Built around the actual hardware, network paths, and legacy software they have to survive in.

**medical imaging** — DICOM viewers, study aggregators, and x-ray processing tools that handle the gap between clinical systems and what people actually need to do with the data.

**image optimization** — [Feather](https://github.com/nathannncurtis/Feather), a lightweight image optimizer for dynamically resizing and compressing TIFF/JPEG in bulk. Built for workflows where thousands of scanned images need to shrink before archival or transfer.

**operational automation** — connecting systems that weren't designed to talk to each other, stabilizing processes that keep breaking, and replacing the manual glue that holds production workflows together.

**dashboards** — internal performance tracking, employee metrics, and operational reporting without spinning up a whole BI platform.

**small utilities** — file renaming, format conversion, record lookups, folder monitoring, archive management. The stuff that saves 20 minutes a day and no one thinks to automate.

---

### current project

[coil](https://github.com/nathannncurtis/coil) — a Python-to-executable compiler. Point it at a project directory, get a standalone .exe back. No spec files, no hook scripts.
