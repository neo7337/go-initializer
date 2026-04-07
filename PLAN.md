# go-initializer — Marketing Plan

## Goal

Grow both **end-users** (Go developers using the tool) and **open-source contributors** (developers submitting PRs), simultaneously and sustainably, through an evergreen content and community strategy.

---

## Key Differentiators (use in all messaging)

These are the angles that separate go-initializer from every other scaffolding tool:

1. **AI Agent scaffolding** — first-class LangChainGo, OpenAI, Gemini, Ollama support with vector store wiring. Unique in the Go ecosystem.
2. **Three access modes** — web UI (zero install), CLI (interactive TUI or CI flag mode), REST API. Covers every workflow from beginner to pipeline.
3. **In-memory generation** — no temp files, no disk side effects.
4. **Identical output** — CLI and web API share the same `internal/generator`; no surprises between interfaces.
5. **Production-ready defaults** — Makefile, multi-stage Dockerfile, `.gitignore`, and a starter README are always included.

---

## Phase 1 — Prepare the Repo for Growth

*One-time work. These are silent, permanent conversion layers — every visitor is funneled toward starring, using, or contributing.*

### 1.1 README overhaul ✅
- Centered hero with badge row (CI, release, Go version, license, Docker pulls)
- Bold one-liner pitch above the fold
- Homebrew as the first install option
- "Why go-initializer?" comparison table vs. manual setup and generic generators
- REST API `curl` example as a third access mode
- Expanded "What gets generated" with every addon path spelled out
- Contributors wall via `contrib.rocks`

### 1.2 GitHub labels and "good first issue" triage
- Add labels: `good first issue`, `help wanted`, `documentation`, `backend`, `frontend`, `ai-agent`, `enhancement`
- Retroactively label 5–10 existing or newly filed issues as `good first issue`
- Pin a "Contributing" issue linking to `CONTRIBUTING.md` and listing open `good first issue` items

### 1.3 Public roadmap
- Create a GitHub Project (kanban) with columns: `Planned`, `In Progress`, `Released`
- List the next 6–10 planned features in priority order
- Answers the "is this project alive?" question for potential contributors

### 1.4 Re-enable the GoReleaser CLI release pipeline
- The Homebrew binary distribution is paused (`if: false` in `.github/workflows/release.yml`)
- Re-enabling it removes a significant friction point — `brew install goini` is the lowest-effort install path

### 1.5 Enable GitHub Discussions
- Low-friction community space better than issues for Q&A, showcasing generated projects, and general feedback
- Categories to create: `Announcements`, `Q&A`, `Show and Tell`, `Ideas`

---

## Phase 2 — Launch Content

*One-time, high-effort posts that create an initial spike in awareness and backlinks.*

### 2.1 "Show HN" post on Hacker News

**Title:** `Show HN: go-initializer – scaffold any Go project (microservice, API, CLI, AI agent) in seconds`

**Body outline:**
- Para 1: The problem — every new Go project starts with the same tedious boilerplate: `go mod init`, a Makefile, a Dockerfile, a `.gitignore`, a logger setup...
- Para 2: The solution — go-initializer generates all of it in one command or a few browser clicks. Three access modes: web UI, CLI, REST API.
- Para 3: The differentiator — first-class AI Agent scaffolding (LangChainGo, OpenAI, Gemini, Ollama) with vector store wiring. One of very few tools that generates complete, runnable LLM agent boilerplate.
- Para 4: Links — repo, hosted web app, docs.

**Logistics:**
- Post on a Tuesday–Thursday, 8–10 AM US Eastern time
- Stay online for 2–3 hours to respond to all comments
- Do not post on a day a major tech story is breaking

### 2.2 Dev.to article: "From Zero to Production-Ready Go Service in 60 Seconds"

**Outline:**
1. The problem: Go boilerplate is tedious and error-prone
2. Walk through `goini new` step by step (screenshots of TUI)
3. Show the generated file tree and highlight key files
4. Run `make run` — it works immediately
5. Add an addon: `--addon database=gorm` — show the generated `internal/database/` package
6. Closing: roadmap teaser + link to repo + CTA to contribute

**Cross-post** to Hashnode after 3 days (canonical URL points to Dev.to for SEO).

### 2.3 Twitter/X launch thread

**Structure:**
- Tweet 1 (hook): "You shouldn't have to write `go mod init`, a Makefile, and a Dockerfile every time you start a Go project. So I built something."
- Tweet 2: Screenshot of `goini new` interactive TUI
- Tweet 3: Screenshot of the generated file tree
- Tweet 4: "It also has a web UI — no install required." + screenshot of go-initializer.dev
- Tweet 5: "And it scaffolds AI Agents — LangChainGo, OpenAI, Gemini, Ollama — with vector store wiring out of the box." + code snippet
- Tweet 6 (CTA): Link to repo + hosted app + "try it in 30 seconds"

**Pin this thread to your profile.**

---

## Phase 3 — Sustained Content Engine

*Repeatable formats. Target: 2–3 pieces of content per month.*

### 3.1 Feature-spotlight posts (Dev.to / Hashnode)

One short article per significant feature shipped. Examples:

| Feature | Article title |
|---|---|
| AI Agent addon | "Scaffolding an AI Agent in Go: LangChainGo, OpenAI, Gemini, Ollama — all wired up" |
| Vector store addon | "pgvector, Qdrant, and chromem in one command with go-initializer" |
| Self-hosting | "Self-hosting go-initializer with Docker Compose" |
| REST API mode | "Using go-initializer as a REST API in your CI pipeline" |

### 3.2 Twitter/X micro-content ("did you know" format)

Low-effort single tweets between bigger posts:

- "Did you know `goini new` supports full flag mode for CI pipelines? No TUI, no prompts."
- "go-initializer generates a multi-stage Dockerfile automatically. Zero config."
- "You can call go-initializer's REST API directly with cURL. No CLI install needed."
- "go-initializer uses in-memory generation — no temp files, no disk writes during scaffolding."

Cadence: 2–3 per week.

### 3.3 Engage existing community content

- When Go blog, r/golang, or other developers post about starting Go projects, reply with go-initializer as a resource — genuinely, not spammily
- Answer Stack Overflow questions about Go project structure and link the tool where it is the correct answer
- Monitor `#golang` and `#go-lang` on Twitter/X and engage

### 3.4 Contributor spotlight (monthly, once contributors exist)

- Monthly tweet or short Dev.to post highlighting a merged PR and the contributor
- Builds contributor incentive and social proof for others considering contributing

---

## Phase 4 — Community Seeding

*One-time posts in existing Go communities.*

### 4.1 r/golang post

- Flair: "Show and Tell"
- Keep it short: what it does, a GIF or screenshot, link to repo and hosted app
- Do not use promotional language — r/golang responds well to genuine developer tools
- Be available to reply to comments for the first few hours

### 4.2 Go newsletters (submit, free)

| Newsletter | Submission |
|---|---|
| **Golang Weekly** | golangweekly.com — use the "submit a link" form |
| **Go Newsletter** | Any open Go community newsletter aggregator |

These newsletters are read by exactly the target audience and require no ongoing effort after submission.

### 4.3 Gophers Slack

- Post in `#show-and-tell` in the Gophers Slack (invite: gophers.slack.com)
- Keep the post to 2–3 sentences + link
- Be available to answer questions in thread

---

## Metrics to Track

| Metric | Starting baseline | Target (3 months) | Tool |
|---|---|---|---|
| GitHub stars | Record current | +200 | GitHub Insights |
| Docker Hub pulls | Record current | +500 | Docker Hub dashboard |
| go-initializer.dev unique visitors | Set up analytics | Establish trend | Plausible or GA4 |
| `good first issue` items claimed | 0 | ≥ 3/month | GitHub |
| External (non-owner) PRs merged | 0 | ≥ 2/month | GitHub |
| Mentions / backlinks | 0 | Track organically | Google Alerts: "go-initializer" |

---

## Execution Timeline

| When | Action |
|---|---|
| **Week 1** | 1.1 README ✅ — complete the remaining Phase 1 items (1.2 labels, 1.3 roadmap, 1.4 GoReleaser, 1.5 Discussions) |
| **Week 2** | Write Show HN post copy + Dev.to article draft |
| **Week 3** | Publish Show HN post + Twitter/X thread; submit to Golang Weekly; r/golang post |
| **Month 2** | First feature-spotlight article; begin micro-content cadence on Twitter/X |
| **Month 3+** | First contributor spotlight; re-evaluate metrics and adjust |
| **Ongoing** | Feature spotlight on each release; micro-content 2–3×/week; community engagement |

---

## Notes

- All content should emphasize **the AI Agent scaffolding** as the primary differentiator — it is timely, unique, and will resonate on HN, Twitter/X, and Dev.to.
- Avoid promotional language in community posts (r/golang, Slack). Let the tool speak for itself.
- The hosted web app at go-initializer.dev is the best entry point for new users — prioritize it in all CTAs.
- The highest single-day leverage action is the **Show HN post** — everything else compounds from the initial spike it generates.
