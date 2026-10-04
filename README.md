# nathan 👋

I build internal tools and native apps. The kind that connect systems no one designed to talk to each other, replace manual processes, and survive environments where "just use a SaaS" isn't an option. Most of my work lives in private repos under my org — here's what I can talk about.

Note: Spikes in contributions (150+ in a day) are bulk Dependabot updates across my personal and org repos

---

### now

**[mdview-zig](https://github.com/nathannncurtis/mdview-zig)** · Native markdown viewer in Zig. DirectWrite / Cairo / CoreText. ~285KB, no webview, no runtime.

**[Study Aggregator](https://github.com/nathannncurtis/Study-Aggregator)** · DICOM processor with a Rust engine. Zero-copy mmap parsing, parallel walks (rayon + jwalk). **100–400× faster** than pydicom.

**[commit-summarizer](https://github.com/nathannncurtis/commit-summarizer)** · GitHub → local Ollama → Slack. Commit summaries for non-technical stakeholders, entirely on-prem.

**[obsidian-vault-sync](https://github.com/nathannncurtis/obsidian-vault-sync)** · Self-hosted real-time Obsidian sync. FastAPI + WebSocket + TypeScript plugin + Docker.

**[steddi](https://steddi.io)** · iOS navigation app for daily commuters. Swift, SwiftUI, MapKit, CarPlay. In development.

**[coil](https://github.com/nathannncurtis/coil)** · Python-to-executable compiler. Directory in, `.exe` out. Published on [PyPI](https://pypi.org/project/coil-compiler/).

**[canopy](https://github.com/nathannncurtis/canopy)** · Fast Windows disk space analyzer. Direct NTFS MFT enumeration, squarified treemap, Fluent Design UI.

---

### what I build

Document processing for legal and medical workflows: batch OCR, PDF compression, DICOM aggregation, image classification. Built around the actual hardware, network paths, and legacy software they have to survive in.

On-prem document AI: local vision-LLM inference (vLLM, FP8) that verifies legal record productions the way a clerk does. **1,700 pages in 18 seconds** on a single workstation GPU, because the documents can't leave the building.

Cross-platform native tools where performance matters and a webview doesn't cut it.

Dashboards, operational automation, and the small utilities that save 20 minutes a day.

---

### tech

![Python](https://img.shields.io/badge/Python-3776AB?style=flat&logo=python&logoColor=white) ![Swift](https://img.shields.io/badge/Swift-F05138?style=flat&logo=swift&logoColor=white) ![Rust](https://img.shields.io/badge/Rust-000000?style=flat&logo=rust&logoColor=white) ![Zig](https://img.shields.io/badge/Zig-F7A41D?style=flat&logo=zig&logoColor=white) ![C](https://img.shields.io/badge/C-A8B9CC?style=flat&logo=c&logoColor=black) ![C++](https://img.shields.io/badge/C++-00599C?style=flat&logo=cplusplus&logoColor=white) ![C#](https://img.shields.io/badge/C%23-239120?style=flat&logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyBmaWxsPSJ3aGl0ZSIgcm9sZT0iaW1nIiB2aWV3Qm94PSIwIDAgMjQgMjQiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI%2BPHRpdGxlPkMjPC90aXRsZT48cGF0aCBkPSJNMS4xOTQgNy41NDN2OC45MTNjMCAxLjEwMy41ODggMi4xMjIgMS41NDQgMi42NzRsNy43MTggNC40NTZhMy4wODYgMy4wODYgMCAwIDAgMy4wODggMGw3LjcxOC00LjQ1NmEzLjA4NyAzLjA4NyAwIDAgMCAxLjU0NC0yLjY3NFY3LjU0M2EzLjA4NCAzLjA4NCAwIDAgMC0xLjU0NC0yLjY3M0wxMy41NDQuNDE0YTMuMDg2IDMuMDg2IDAgMCAwLTMuMDg4IDBMMi43MzggNC44N2EzLjA4NSAzLjA4NSAwIDAgMC0xLjU0NCAyLjY3M1ptNS40MDMgMi45MTR2My4wODdhLjc3Ljc3IDAgMCAwIC43NzIuNzcyLjc3My43NzMgMCAwIDAgLjc3Mi0uNzcyLjc3My43NzMgMCAwIDEgMS4zMTctLjU0Ni43NzUuNzc1IDAgMCAxIC4yMjYuNTQ2IDIuMzE0IDIuMzE0IDAgMSAxLTQuNjMxIDB2LTMuMDg3YzAtLjYxNS4yNDQtMS4yMDMuNjc5LTEuNjM3YTIuMzEyIDIuMzEyIDAgMCAxIDMuMjc0IDBjLjQzNC40MzQuNjc4IDEuMDIzLjY3OCAxLjYzN2EuNzY5Ljc2OSAwIDAgMS0uMjI2LjU0NS43NjcuNzY3IDAgMCAxLTEuMDkxIDAgLjc3Ljc3IDAgMCAxLS4yMjYtLjU0NS43Ny43NyAwIDAgMC0uNzcyLS43NzIuNzcxLjc3MSAwIDAgMC0uNzcyLjc3MlptMTIuMzUgMy4wODdhLjc3Ljc3IDAgMCAxLS43NzIuNzcyaC0uNzcydi43NzJhLjc3My43NzMgMCAwIDEtMS41NDQgMHYtLjc3MmgtMS41NDR2Ljc3MmEuNzczLjc3MyAwIDAgMS0xLjMxNy41NDYuNzc1Ljc3NSAwIDAgMS0uMjI2LS41NDZ2LS43NzJIMTJhLjc3MS43NzEgMCAxIDEgMC0xLjU0NGguNzcydi0xLjU0M0gxMmEuNzcuNzcgMCAxIDEgMC0xLjU0NGguNzcydi0uNzcyYS43NzMuNzczIDAgMCAxIDEuMzE3LS41NDYuNzc1Ljc3NSAwIDAgMSAuMjI2LjU0NnYuNzcyaDEuNTQ0di0uNzcyYS43NzMuNzczIDAgMCAxIDEuNTQ0IDB2Ljc3MmguNzcyYS43NzIuNzcyIDAgMCAxIDAgMS41NDRoLS43NzJ2MS41NDNoLjc3MmEuNzc2Ljc3NiAwIDAgMSAuNzcyLjc3MlptLTMuMDg4LTIuMzE1aC0xLjU0NHYxLjU0M2gxLjU0NHYtMS41NDNaIi8%2BPC9zdmc%2B) ![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=flat&logo=typescript&logoColor=white) ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white) ![Kotlin](https://img.shields.io/badge/Kotlin-7F52FF?style=flat&logo=kotlin&logoColor=white) ![PHP](https://img.shields.io/badge/PHP-777BB4?style=flat&logo=php&logoColor=white) ![PowerShell](https://img.shields.io/badge/PowerShell-5391FE?style=flat&logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyBmaWxsPSJ3aGl0ZSIgcm9sZT0iaW1nIiB2aWV3Qm94PSIwIDAgMjQgMjQiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI%2BPHRpdGxlPlBvd2VyU2hlbGw8L3RpdGxlPjxwYXRoIGQ9Ik0yMy4xODEgMi45NzRjLjU2OCAwIC45MjMuNDYzLjc5MiAxLjAzNWwtMy42NTkgMTUuOTgyYy0uMTMuNTcyLS42OTcgMS4wMzUtMS4yNjUgMS4wMzVILjgxOWMtLjU2OCAwLS45MjMtLjQ2My0uNzkyLTEuMDM1TDMuNjg2IDQuMDA5Yy4xMy0uNTcyLjY5Ny0xLjAzNSAxLjI2NS0xLjAzNXptLTguMzc1IDkuMzQ2Yy4yNTEtLjM5NC4yMjctLjkwNS0uMDktMS4yNDNMOS4xMjIgNS4xMjVjLS4zOC0uNDA0LTEuMDM3LS40MDctMS40NjYtLjAwMy0uNDI5LjQwMi0uNDY4IDEuMDU2LS4wODggMS40Nmw0LjY2MiA0Ljk2di4xMWwtNy40MiA1LjM3NGMtLjQ1LjMyNy0uNTMzLjk3Ny0uMTg3IDEuNDUzLjM0Ni40NzYuOTkxLjU5NyAxLjQ0LjI3bDguMjI5LTUuOTFjLjI4LS4xOTYuNDM4LS4zNjUuNTE0LS41MnptLTIuNzk2IDQuMzk5YS45MjguOTI4IDAgMDAtLjkzNC45MjNjMCAuNTEuNDE4LjkyMy45MzQuOTIzaDQuNDMzYS45MjguOTI4IDAgMDAuOTM0LS45MjMuOTI4LjkyOCAwIDAwLS45MzQtLjkyM3oiLz48L3N2Zz4%3D) ![Bash](https://img.shields.io/badge/Bash-4EAA25?style=flat&logo=gnubash&logoColor=white) ![AutoHotkey](https://img.shields.io/badge/AutoHotkey-334455?style=flat&logo=autohotkey&logoColor=white)

![PyTorch](https://img.shields.io/badge/PyTorch-EE4C2C?style=flat&logo=pytorch&logoColor=white) ![ONNX](https://img.shields.io/badge/ONNX-005CED?style=flat&logo=onnx&logoColor=white) ![Hugging Face](https://img.shields.io/badge/Hugging%20Face-FFD21E?style=flat&logo=huggingface&logoColor=black) ![vLLM](https://img.shields.io/badge/vLLM-30364D?style=flat) ![Ollama](https://img.shields.io/badge/Ollama-000000?style=flat&logo=ollama&logoColor=white) ![CUDA](https://img.shields.io/badge/CUDA-76B900?style=flat&logo=nvidia&logoColor=white) ![llama.cpp](https://img.shields.io/badge/llama.cpp-30364D?style=flat) ![OpenVINO](https://img.shields.io/badge/OpenVINO-0071C5?style=flat) ![timm](https://img.shields.io/badge/timm-EE4C2C?style=flat) ![LoRA / QLoRA](https://img.shields.io/badge/LoRA%20%2F%20QLoRA-7B2CBF?style=flat) ![OpenCV](https://img.shields.io/badge/OpenCV-5C3EE8?style=flat&logo=opencv&logoColor=white) ![PaddleOCR](https://img.shields.io/badge/PaddleOCR-0062B0?style=flat&logo=paddlepaddle&logoColor=white) ![NumPy](https://img.shields.io/badge/NumPy-013243?style=flat&logo=numpy&logoColor=white)

![Mistral AI](https://img.shields.io/badge/Mistral%20AI-FA520F?style=flat&logo=mistralai&logoColor=white) ![Whisper](https://img.shields.io/badge/Whisper-412991?style=flat) ![Qwen](https://img.shields.io/badge/Qwen-615CED?style=flat) ![Gemma](https://img.shields.io/badge/Gemma-4285F4?style=flat) ![DINOv2](https://img.shields.io/badge/DINOv2-0467DF?style=flat) ![GLiNER](https://img.shields.io/badge/GLiNER-6E3FF3?style=flat)

![React](https://img.shields.io/badge/React-61DAFB?style=flat&logo=react&logoColor=black) ![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-06B6D4?style=flat&logo=tailwindcss&logoColor=white) ![Vite](https://img.shields.io/badge/Vite-9135FF?style=flat&logo=vite&logoColor=white) ![SwiftUI](https://img.shields.io/badge/SwiftUI-F05138?style=flat&logo=swift&logoColor=white) ![Jetpack Compose](https://img.shields.io/badge/Jetpack%20Compose-4285F4?style=flat&logo=jetpackcompose&logoColor=white) ![Electron](https://img.shields.io/badge/Electron-47848F?style=flat&logo=electron&logoColor=white) ![.NET](https://img.shields.io/badge/.NET-512BD4?style=flat&logo=dotnet&logoColor=white) ![FastAPI](https://img.shields.io/badge/FastAPI-009688?style=flat&logo=fastapi&logoColor=white) ![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat&logo=docker&logoColor=white) ![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-2088FF?style=flat&logo=githubactions&logoColor=white) ![Caddy](https://img.shields.io/badge/Caddy-1F88C0?style=flat&logo=caddy&logoColor=white) ![WireGuard](https://img.shields.io/badge/WireGuard-88171A?style=flat&logo=wireguard&logoColor=white) ![Cloudflare Workers](https://img.shields.io/badge/Cloudflare%20Workers-F38020?style=flat&logo=cloudflareworkers&logoColor=white) ![Inno Setup](https://img.shields.io/badge/Inno%20Setup-264B99?style=flat) ![SQLite](https://img.shields.io/badge/SQLite-003B57?style=flat&logo=sqlite&logoColor=white) ![SQL Server](https://img.shields.io/badge/SQL%20Server-CC2927?style=flat) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white) ![Redis](https://img.shields.io/badge/Redis-FF4438?style=flat&logo=redis&logoColor=white)

![DICOM](https://img.shields.io/badge/DICOM-0B5394?style=flat) ![Win32 API](https://img.shields.io/badge/Win32%20API-0078D4?style=flat) ![Tesseract OCR](https://img.shields.io/badge/Tesseract%20OCR-5B5B5B?style=flat) ![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=flat&logo=ubuntu&logoColor=white) ![systemd](https://img.shields.io/badge/systemd-201A26?style=flat)

---

<img src="assets/timeline.svg" alt="Recent activity" />

---

<sub>get to know me better at [nathancurtis.to](https://nathancurtis.to)</sub>
