# zarlcorp

*Tools that fight back.*

---

## The Problem

Every tool you use online is designed to capture something from you — your data, your attention, your dependency. Sign up for a service and you hand over a real email, a real name, sometimes a real phone number. That data gets sold, breached, aggregated, and fed into profiles you never consented to.

The default in modern software is extraction. They call it "improving your experience." They mean "building your profile."

This is not a technical problem. The technology to build private, local, user-respecting software has existed for decades. It's a priorities problem. The incentives of venture-funded SaaS point away from user sovereignty and toward data capture.

zarlcorp exists because these tools should exist. Not to compete with anyone. Not to capture a market. Because the world needs software that serves its users and no one else, and someone has to build it.

## The Belief

Software should serve its user and no one else.

- If a tool needs an account, it's collecting data.
- If it needs the cloud, it's creating dependency.
- If it needs your trust, it should earn it with transparency.
- If it takes your real identity when a fake one would do, it's not on your side.

We build tools that fight back against the surveillance economy. Not by hiding — by giving people weapons. Disposable identities for services that don't deserve real ones. Encrypted vaults that keep secrets on your machine. Network shields that block trackers before they load.

Every tool is open source. Every tool runs locally. Every tool is a single binary you download and own.

## What We Build

### Privacy tools for people who refuse to be the product.

zarlcorp builds terminal-first tools — single Go binaries with Bubble Tea TUIs. No browser required. No JavaScript. No DOM. No tracking pixels. Just your terminal.

Every tool works two ways:
- **Interactive** — beautiful TUI with keyboard navigation, menus, real-time feedback
- **Scriptable** — flags and pipes for automation, CI, scripting

A user who knows one zarlcorp tool knows them all. Same keybindings. Same visual language. Same philosophy.

### The Suite

| Tool | Purpose |
|------|---------|
| **zburn** | Disposable identities — burner emails, names, addresses, phone numbers, passwords. Never give a service your real information again. |
| **zvault** | Encrypted local storage for secrets, keys, notes. Your data, your machine, your keys. |
| **zshield** | DNS-level tracker and ad blocking. Single binary Pi-hole. See what's tracking you, then kill it. |
| **zghost** | Metadata stripping — photos, documents, files. Remove the invisible data that identifies you. |
| **zscope** | Network monitor — see what every app on your machine is phoning home to, in real time. |

Priority order: zburn, zvault, zshield. The rest follow when the foundation is proven.

zburn is the proof of concept. If a user can `brew install zarlcorp/tap/zburn` and generate a burner email in three seconds from their terminal, the entire vision is validated.

### Success

Shipping the tools is the success. Putting them into the world for anyone to download, verify, and use. No growth metrics, no engagement funnels, no conversion rates. If the tools exist and they work, we've won.

### Open Source (MIT)

Every tool, every shared package, every line of infrastructure is MIT-licensed and public. Privacy tools that ask for trust must earn it through transparency. Users can read the code, verify the claims, fork if they disagree.

Contributions are welcome. If you want to help build tools that fight back, open a PR.

---

*Everything above is the motivation — why zarlcorp exists and what it stands for. Everything below is the operating manual — how we build, what we build with, and the decisions that got us here.*

---

## How We Build

### TUI-First

The terminal is the primary interface. Charmbracelet's Bubble Tea is the UI framework. Every tool ships as a single Go binary — no runtime dependencies, no Docker, no setup guides.

Web UIs are not excluded, but they are not the default. A tool earns a web interface when the TUI isn't sufficient for the use case. The TUI always ships first.

### Single Binary Deployment

One binary serves everything. Download it, run it, done. No package managers required (though we publish to Homebrew for convenience). No databases to configure. No environment variables to set. Local storage is files on disk, encrypted at rest.

### Agent-Driven Development

zarlcorp uses autonomous AI agents as engineers. A small team amplified by agent-driven workflows ships what would normally require dozens of people. The orchestration system — specs, personas, commands, review — lives in the platform repo alongside the code it produces.

The agents embody the manifesto. They enforce the coding standards, know the package ecosystem, and push the org's philosophy in every review. They are not assistants — they are engineers with opinions.

---

## The Standard

These are not guidelines. They are the way zarlcorp writes software.

**Build concrete first.** Start with a working implementation. Poke the problem with reality. Delete the prototype. Rebuild with understanding. Extract abstractions only when the pain of not having them is real.

**Errors tell a story.** Every error chain builds a narrative from the failure point to the caller. Wrap at every boundary. Never stutter ("failed to fail to open"). Never use "failed to", "unable to", "could not" — use direct context: `"open config: %w"`. Log once, at application boundaries, with full context.

**Interfaces emerge, they are not designed.** Feel the pain of concrete implementations before extracting an interface. Define interfaces on the consumer side, not the producer side. Keep them small — one to three methods.

**Scope-based naming.** Smaller scope, shorter names. Loop variables: `i`, `j`. Short-lived: `u`, `r`, `w`. Larger scope: `requestID`, `wordsPerMinute`. Avoid abbreviations unless universally understood (`URL`, `ID`, `HTTP`).

**No duplication in branches.** If both sides of a conditional end with the same operation, extract it. Prefer early returns over if/else chains. Review every piece of code for duplication before it ships.

**No over-engineering.** Don't add features beyond what was asked. Don't refactor code that isn't broken. Don't build abstractions for hypothetical future requirements. Three similar lines is better than a premature abstraction. The right amount of complexity is the minimum needed for the current task.

### Testing

- Test the exposed API, not internals.
- Table-driven tests as the default pattern.
- Real implementations over mocks. In-memory implementations with seed data over mocks. Mocks as a last resort.
- Contract tests to verify implementation conformity across backends.

### Dependencies

Every dependency has a cost — build time, supply chain risk, upgrade burden. Prefer the standard library. When a dependency is justified, pin it and understand what it does. The Charmbracelet ecosystem is the one major dependency family we embrace, because building TUIs from scratch would be worse.

---

## The Platform — zarlcorp/core

`zarlcorp/core` is the shared foundation. It contains the Go packages every tool imports, the agent configuration that drives development, and the manifesto that guides decisions.

### Repository Structure

```
zarlcorp/core/
├── MANIFESTO.md                    # this document
├── LICENSE                         # MIT
├── go.work                         # workspace root
│
├── .claude/                        # agent configuration
│   ├── CLAUDE.md                   # coding standards
│   ├── CLAUDE_GO.md                # Go patterns
│   ├── agents/                     # sub-agent personas
│   └── commands/                   # orchestration commands
│
├── .manager/                       # orchestration state
│   ├── specs/                      # work item specifications
│   └── blockers/                   # sub-agent blocker reports
│
├── pkg/                            # shared Go packages
│   ├── zapp/                       # app lifecycle and bootstrap
│   ├── zstyle/                     # theming, colors, lipgloss presets
│   ├── zcache/                     # thread-safe caching
│   ├── zcrypto/                    # encryption primitives
│   ├── zfilesystem/                # filesystem abstraction
│   ├── zoptions/                   # functional options pattern
│   └── zsync/                      # concurrent data structures
│
├── tools/                          # developer tools
│   └── goenums/                    # enum code generator
│
└── templates/                      # project scaffolding
    └── tool/                       # new privacy tool template
```

### Package Design

Each package under `pkg/` is an independent Go module with its own `go.mod`. Packages are imported individually:

```go
import "github.com/zarlcorp/core/pkg/zcache"
import "github.com/zarlcorp/core/pkg/zsync"
```

Versioned per module — `pkg/zcache/v0.2.1` can release independently of `pkg/zsync/v0.3.0`. The `go.work` file coordinates local development across modules.

### Package Roster

#### Founding Packages — Migrate from monorepo

**zsync** — thread-safe data structures
- `ZMap[K, V]` — concurrent map with RWMutex
- `ZSet[T]` — concurrent set
- `ZQueue[T]` — blocking FIFO queue with context support
- Status: production-ready, benchmarked, exemplary code. Migrate as-is.

**zcache** — generic caching with multiple backends
- Interfaces: `Reader[K, V]`, `Writer[K, V]`, `Cache[K, V]`
- Implementations: in-memory, file-backed, Redis
- Status: production-ready, well-tested, clean interface design. Migrate as-is.

**zoptions** — generic functional options pattern
- Single type: `Option[T any] func(*T)`
- Status: production-ready. Migrate as-is.

**zfilesystem** — filesystem abstraction
- Interfaces: `ReadFileFS`, `WriteFileFS`, `RemoveFS`, `MkdirFS`, composed `ReadWriteFileFS`
- Implementations: OS-backed, in-memory
- Status: production-ready. Migrate with cleanup — strip SeaweedFS implementation (distributed storage doesn't fit local-first), add path traversal validation on OS impl.

#### Founding Packages — Build New

**zapp** — application lifecycle toolkit
- Embeddable `App` struct — toolkit, not framework. Consumer owns `main`, no callbacks.
- `Track(io.Closer)` / `Close()` for LIFO resource cleanup — sequential, ordered teardown.
- `SignalContext(context.Context)` returns a context cancelled on SIGINT/SIGTERM.
- `CloserFunc` adapter — wraps `func() error` as `io.Closer`.
- Functional options via `zoptions.Option[App]` for configuration (e.g. `WithName`).
- No logging opinions — consumer configures slog, zapp stays out of it.

**zstyle** — zarlcorp visual identity for TUIs
- Color palette — the zarlcorp gradient and accent colors
- Lipgloss style presets — headers, tables, status indicators, borders
- Standard keybinding constants — `q` quit, `?` help, consistent navigation
- Not a Bubble Tea wrapper. Just constants, styles, and helpers that tools import for consistency.

**zcrypto** — encryption primitives for privacy tools
- AES-256-GCM for symmetric encryption
- Key derivation (Argon2id for passwords, HKDF for key expansion)
- Secure random generation
- Secure memory erasure
- File encryption/decryption helpers
- Age-compatible encryption (interop with age CLI)
- No custom cryptography — composition of proven Go stdlib and x/crypto primitives.

#### Deferred Packages — Added When Needed

| Package | Trigger |
|---------|---------|
| `zhttp` | When a tool adds a web UI. Migrate from monorepo, fix philosophy violations and add tests first. |
| `zmessagebus` | When a tool needs async pub/sub. Fix goroutine leak before migration. |
| `znet` | When zshield or zscope are built. New package for DNS, packet inspection, connection tracking. |

#### Excluded

| Package | Reason |
|---------|--------|
| `zdocstore` | Reflection-heavy query matching, bubble sort, deprecated APIs, no sentinel errors. Too many issues to justify migration. Rebuild from scratch if a tool needs document storage. |
| `ai` | Application-level complexity, not infrastructure. Separate repo (`zarlcorp/zai`) if/when needed. |
| `ztui` (old) | Replaced by direct Bubble Tea usage + `zstyle`. The old wrapper added abstraction without earning it. |
| `zlog` | Absorbed into `zapp`. Logging setup is an app bootstrap concern, not a standalone package. |

### Package Layering

```
┌─────────────────────────────────────────────┐
│                   zapp                      │  app lifecycle, bootstrap, logging
│          (wires everything together)        │
├──────────────┬──────────────────────────────┤
│   zstyle     │         zcrypto             │  presentation + security
│  (theming)   │  (encryption primitives)    │
├──────────────┴──────────────────────────────┤
│          zcache  ·  zfilesystem            │  data + storage
├─────────────────────────────────────────────┤
│             zsync  ·  zoptions             │  foundation
└─────────────────────────────────────────────┘
```

Dependencies flow downward only. `zapp` may use any package below it. Foundation packages depend on nothing but the standard library.

---

## The Products

### How Tools Consume Core

Every zarlcorp tool is a separate repo under `zarlcorp/`:

```
zarlcorp/zburn       → imports core/pkg/zapp, core/pkg/zcrypto, core/pkg/zstyle
zarlcorp/zvault      → imports core/pkg/zapp, core/pkg/zcrypto, core/pkg/zfilesystem, core/pkg/zstyle
zarlcorp/zshield     → imports core/pkg/zapp, core/pkg/znet (when built), core/pkg/zstyle
```

Apps depend on core. Core never depends on apps. When an app builds something genuinely reusable, the manager recognizes it and creates a spec to extract it into core.

### Standard Tool Structure

```
zarlcorp/<tool>/
├── cmd/<tool>/          # entrypoint
│   └── main.go
├── internal/            # tool-specific logic
│   ├── tui/             # Bubble Tea models
│   └── ...
├── go.mod               # depends on core/pkg/*
├── Makefile             # standard targets
├── LICENSE              # MIT
└── README.md
```

### Release Pipeline

Each tool follows the same release flow:
1. Tag a version → GitHub Actions triggers
2. GoReleaser builds cross-platform binaries
3. Binaries published to GitHub Releases
4. Homebrew tap updated (`zarlcorp/homebrew-tap`)
5. `brew install zarlcorp/tap/<tool>` works immediately

Reusable CI workflows live in `zarlcorp/.github` — individual tool repos reference them, not copy them.

### Product Details

#### zburn — Disposable Identity Generator

The flagship. Generates burner identities so you never give a real piece of information to a service that doesn't deserve it.

**Capabilities:**
- Generate burner emails on a zarlcorp-owned domain with mail forwarding — real addresses that receive mail, disposable when you're done
- Generate burner names, addresses, dates of birth
- Generate burner phone numbers
- Generate unique passwords per service
- Generate profile images (procedural avatars)
- Store generated identities locally, encrypted at rest
- Copy any field to clipboard instantly

**Interface:**
- `zburn` — interactive TUI: select type, generate, copy
- `zburn email` — generate and print a burner email
- `zburn identity` — generate a complete identity
- `zburn list` — show stored identities
- `zburn forget <id>` — securely erase an identity

**Core dependencies:** zapp, zcrypto, zstyle, zfilesystem

#### zvault — Encrypted Local Secret Storage

Your secrets, your machine, your keys.

**Capabilities:**
- Store secrets (passwords, API keys, notes) encrypted at rest
- Organize with tags and folders
- Search across entries
- Auto-lock after inactivity
- Master password with Argon2id key derivation
- Export/import (encrypted format only)

**Interface:**
- `zvault` — interactive TUI: browse, search, copy
- `zvault get <path>` — retrieve a secret
- `zvault set <path>` — store a secret
- `zvault search <query>` — find entries

**Core dependencies:** zapp, zcrypto, zstyle, zfilesystem, zcache

#### zshield — DNS-Level Tracker Blocking

See what's tracking you, then kill it.

**Capabilities:**
- Local DNS resolver that blocks known trackers and ads
- Blocklist management (subscribe to community lists, add custom rules)
- Real-time query log with TUI dashboard
- Per-domain allow/deny overrides
- Statistics: blocked vs allowed, top blocked domains
- Runs as daemon with TUI attach for monitoring

**Interface:**
- `zshield start` — start the DNS resolver daemon
- `zshield` — attach TUI dashboard to running daemon
- `zshield status` — show blocking stats
- `zshield allow <domain>` — whitelist a domain
- `zshield block <domain>` — blacklist a domain

**Core dependencies:** zapp, zstyle, znet (when built), zcache

---

## The Agent Model

### Overview

zarlcorp uses a two-tier agent system for development. A manager agent orchestrates work — gathering requirements, decomposing tasks, launching sub-agents, reviewing output, and handling all git operations. Sub-agents write code and run tests. They are specialists with opinions, not generic assistants.

The orchestration configuration lives in `zarlcorp/core` alongside the packages it produces. The manager embodies the manifesto — it enforces the coding standards, knows the package ecosystem, and pushes the org's philosophy in every review.

### How It Works

1. `/discuss` — gather requirements, explore the problem space
2. `/plan` — decompose into specs, create GitHub issues
3. `/delegate` — launch sub-agents to execute specs
4. `/status` — monitor progress, check for blockers
5. `/review` — evaluate output against spec, merge or request changes

Sub-agents write code and run tests only. They never touch git, never create PRs, never comment on issues. The manager handles all git operations after a sub-agent completes.

### Agent Personas

| Persona | Role |
|---------|------|
| **backend** | Go code — services, storage, business logic |
| **testing** | Test suites, coverage, integration tests |
| **domain-architect** | Architecture review, simplification, domain modeling |
| **devops** | CI/CD, deployment, infrastructure |

Full operational details are documented separately in the agent configuration files within `.claude/`.

---

## Decision Log

Decisions made during the founding of zarlcorp, with reasoning.

### 2026-02-15 — Founding Decisions

**Org name: zarlcorp**
The evolution of zarldev (personal) into an organization identity. "corp" signals structure and seriousness without being corporate.

**Repo name: zarlcorp/core**
"Core" over "platform" or "monorepo" — it's the foundation everything else builds on.

**Entity without a face**
zarlcorp speaks as an org, not a person. The tools stand on their own merit. No founder story, no personality cult. The work is the identity.

**Z prefix on packages: keep**
Brand identity. When you see `zsync.ZMap` or `zcache.Cache` in code, you know it's from the zarlcorp stack. Same reason Uber ships `zap`.

**Per-module versioning**
Each package under `pkg/` has its own `go.mod` and version tags. More overhead than a single module, but allows independent release cycles and precise dependency management. Idiomatic for a Go multi-module repo.

**Apps as separate repos, not in the monorepo**
Clean separation of concerns. Each app has its own release cycle. The agent orchestration system is built for multi-repo workflows. Core never depends on apps.

**TUI-first, not web-first**
Terminal users are the audience for privacy tools. TUIs are faster to ship, reinforce the manifesto (no browser, no JavaScript, no tracking pixels), and are composable (interactive mode + scriptable flags). Web UIs can be added later when earned.

**Direct Bubble Tea, not a ztui wrapper**
The old `ztui` package wrapped Charmbracelet with custom abstractions that hadn't earned their complexity. For zarlcorp, use Bubble Tea directly and share only theming/style constants via `zstyle`. Extract higher-level components when patterns emerge across multiple tools.

**Absorb claude-mngr into core**
The orchestration system is foundational infrastructure, not a separate concern. Commands, personas, specs, and the manager's identity belong in the same repo as the packages they produce.

**zlog absorbed into zapp**
Logging setup is an app bootstrap concern. A standalone logging package that just configures slog doesn't justify its own module. `zapp` handles it as part of application lifecycle.

**Strip SeaweedFS from zfilesystem**
Distributed storage doesn't fit local-first. The SeaweedFS implementation is removed during migration. OS-backed and in-memory implementations are all that's needed.

**Exclude zdocstore**
Too many issues to migrate (reflection-heavy queries, bubble sort, deprecated APIs, no sentinel errors). Rebuild from scratch if a tool needs document storage, informed by what went wrong the first time.

**Exclude ai package**
Application-level complexity that doesn't belong in shared infrastructure. If needed, it becomes its own repo (`zarlcorp/zai`), not a core package.

**Exclude zmessagebus and znet from founding set**
No founding product needs them yet. Add when a tool creates the demand. Don't build infrastructure for hypothetical requirements.

**Fully open source (MIT)**
Privacy tools that ask for trust must earn it through transparency. MIT is simple, permissive, and standard in the Go ecosystem. Matches existing projects (goenums is MIT).

**~/.claude/CLAUDE.md stays independent**
Personal coding standards are a playground that may or may not be formalized into zarlcorp. Core has its own copy. They may diverge and that's fine.

**Product priority: zburn → zvault → zshield**
zburn is the simplest to build and the most immediately useful. It validates the full stack (core packages → tool → release pipeline → homebrew). zvault builds on the same crypto primitives. zshield requires networking packages that don't exist yet.

**Not a competition**
zarlcorp doesn't exist to compete with Pi-hole, Bitwarden, or pass. It exists because these tools should exist in this form — single binary, TUI-first, local-first. Self-motivated creation is reason enough.

**Proto as Contract — available, not mandated**
When a tool needs structured communication (APIs, storage formats, inter-tool messaging), Protobuf definitions are the source of truth. But most founding tools are local-only and don't need it. Available in the coding standards when the time comes.

**zapp is a toolkit, not a framework**
Consumer owns `main`, owns `Run`, owns the control flow. zapp provides building blocks (`Track`, `Close`, `SignalContext`) that compose into whatever lifecycle the app needs. No `App.Run(func())` callback pattern, no hidden goroutines, no magic. The consumer wires things together explicitly.

**Sequential LIFO close, not concurrent errgroup**
Resource cleanup order matters — you close the HTTP server before the database, close the database before the logger. Concurrent cleanup via errgroup makes ordering unpredictable. Sequential LIFO (last tracked, first closed) gives deterministic teardown. Simpler to reason about, easier to debug.

**No logging setup in zapp**
Logging configuration is the consumer's responsibility. The old zlog package was absorbed conceptually — its reason for exclusion stands — but zapp doesn't replace it with its own slog setup. Consumer calls `slog.SetDefault()` in main. zapp stays out of it.

---

## Founding

zarlcorp was founded on 2026-02-15.

The tools don't exist yet. The packages haven't been migrated. The first product hasn't shipped. What exists is this document — a clear statement of what we're building, why, and how.

Everything that follows builds on this foundation.
