# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

Anyone who wants to learn Go, from a first-ever line of code through concurrency, testing, data tools, services, databases, networking and systems. The curriculum culminates in a staged crews-and-tasks API. The creator plans to publish it, so first-time visitors who know nothing about the game must understand it within seconds. Lessons, hints, quizzes and messages support Egyptian Arabic and English; code and Go terms stay in English. Players use laptops, large monitors, tablets and iPads.

## Product Purpose

Orbital Colony teaches Go through missions: the player commands drones on planet Kepler-G, then builds tools and services for the colony by writing real Go, compiled and run with the real `go` toolchain. The implemented curriculum has 181 missions across 38 chapters: 174 missions in chapters 1–37 and 7 optional Modern Go labs in chapter 38. It includes 58 drone missions and 123 labs, with 8 track bosses counted among the labs. Success is a player who keeps going because each mission is worth solving, and who can build and test a Go API using the concepts learned along the way.

## Positioning

The code is real Go, built by the real compiler: real compiler errors (with localized explanations), real goroutines running in parallel wall-clock time, the real race detector, real `go test` with mutation testing and benchmarks. Seeded worlds and test data vary the challenges; AST checks enforce required source features. Stars reward mission-specific goals, benchmark speedups or bonus checks. The promise is hands-on practice with the real tools.

## Operating Context

- A local Go server (`go run ./cmd/orbital`) serves the web client on localhost and runs player code.
- Each mission: read the lesson, write code, press Run (or save in VS Code, which syncs), inspect the animation or lab results, read the console, earn 1–3 stars.
- Labs appear throughout the curriculum, including chapter 16 and track bosses. They show verification steps rather than drone animation: player tests, hidden tests, mutants, benchmarks, fuzzing, bonus checks or vet, as defined by the mission. Labs marked `Race` run with `-race`.
- Player profiles persist progress, code, preferences, best quiz scores and review metadata on the local server under `saves/profiles/`; localStorage keeps a browser copy. Older saves can load without quiz or review fields.
- Missions normally unlock in curriculum order. A mission immediately following a service-backed mission can be opened without completing that service mission. Track groups organize the index; they do not independently unlock their chapters.
- External lab libraries are pinned in the separate `labdeps` module and prepared with `make deps` before offline player builds. `Level.Deps` controls allowed external imports. `make services-up` starts local Postgres, MongoDB and Kafka-compatible Redpanda for service labs; SQLite runs in-process and Redis exercises use miniredis. The game runs locally and executes player code; import restrictions do not provide isolation for untrusted code.

## Curriculum

Counts come from `colony/levels` and include bosses. Quizzes, review attempts and the open colony mode are separate from the mission count.

| Track | Chapters | Missions | Culmination |
|---|---|---|---|
| Fundamentals | 1–9 | 37 | Chapter 9 boss |
| Concurrency | 10–13, 24–25 | 27 | Chapter 25 boss |
| Craft (generics, algorithms, tests) | 14–17 | 19 | Chapter 17 boss |
| Data and Tools | 18–23 | 27 | Chapter 23 boss |
| Services | 26–28 | 16 | Chapter 28 boss |
| Databases | 29–30 | 11 | Chapter 30 boss |
| Networking | 31–33 | 17 | Chapter 33 boss |
| Systems | 34–36 | 15 | Chapter 36 boss |
| Capstone | 37 | 5 | One API developed in stages |
| Extras | 38 | 7 | Optional Modern Go labs |

The eight bosses combine earlier concepts in each of the first eight tracks and require no external services. The five-stage Capstone carries code forward through schema and migrations, opaque sessions, isolated crews and roles, idempotent versioned tasks with cursor pagination, and operations with request IDs, metrics, rate limiting and graceful shutdown. Its hidden tests always exercise SQLite and also exercise Postgres when available; passing offline does not verify Postgres integration.

Each of the 38 chapters has five bilingual, no-code self-assessment questions with explanations. The best score is saved independently of mission stars. After every five newly completed missions, spaced review offers an older completed mission with a fresh seed; players can start or skip it. Drone reviews have an empty starter; lab reviews retain their own starter files and do not carry saved solutions forward. Reviews preserve mission stars and saved solutions.

Chapter 38 unlocks after the last mission of chapter 17, advances in its own sequence and is not required for the main curriculum. Its seven labs cover iterators and range-over-function, sets and labels, Once, errgroup, allocation measurement and escape analysis, vet, and CPU/heap profiling with pprof.

## Capabilities and Constraints

- Lessons use a small markdown dialect (`~code~`, `~~~` fences, tables).
- The editor is a textarea with a syntax-highlight overlay; VS Code linking gives full gopls tooling.
- The world is drawn on a canvas: tiles (regolith, soil, ice, rock), plants, crystals, beacons, multiple coloured drones, storms.
- Sound effects are generated in code (no audio files), with a mute control.
- Must work from tablet width to large monitors.

## Brand Commitments

Name: Orbital Colony. Setting: a space colony on planet Kepler-G. Voice: friendly Egyptian Arabic, encouraging, a little playful.

## Evidence on Hand

No testimonials, user counts or press exist. Do not invent any.

## Product Principles

1. Real Go, never a toy imitation.
2. Learn by doing and by hitting the bug yourself before being told the answer.
3. Every mission should feel like a game worth one more try.
4. Clear to a total beginner, still challenging for an expert.

## Accessibility & Inclusion

Right-to-left Arabic text alongside left-to-right code. Readable contrast for long sessions; sound must be mutable.
