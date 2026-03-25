# Introduction

**Go Initializer** is a web-based scaffolding tool that generates production-ready Go project starter kits in seconds.

## What it does

Select your project type, Go version, and preferred framework — Go Initializer assembles a clean project structure with:

- Correct `go.mod` configuration
- Framework-specific routing boilerplate
- Optional Docker support
- Optional add-ons (databases, caches, logging, vector stores)

## Why use it?

Setting up a new Go project involves a lot of repetitive decisions: folder layout, module naming, dependency wiring. Go Initializer encodes community best practices so you can skip the boilerplate and start building immediately.

## Project types

| Type | Description |
|------|-------------|
| Simple Project | A single-binary HTTP service with minimal structure |
| Microservice | A production-oriented service with health probes and structured internals |
| API Server | A REST API layout with middleware composability |
| CLI Application | A sub-command CLI binary with Cobra / urfave scaffolding |
| AI Agent | A full LLM-backed agent project with tool-calling, vector store, and provider wiring |

## AI & Workflow integration

Go Initializer ships first-class support for **AI Agent** projects. Pick an LLM provider (LangChainGo, OpenAI, Gemini, or Ollama), optionally add a vector store (pgvector, Qdrant, chromem), and the generated project is already wired to run as a standalone agent or drop into a larger orchestration pipeline.

See the [AI Agent Guide](ai-agent) for a walkthrough.

---

Ready to dive in? Head to [Quick Start](quick-start) or explore the generator directly.
