# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

Anyone who wants to learn Go, from a first-ever line of code to advanced concurrency. The creator plans to publish it, so first-time visitors who know nothing about the game must understand it within seconds. Lessons, hints and messages are in Egyptian Arabic; code and Go terms stay in English. Players use laptops, large monitors, tablets and iPads.

## Product Purpose

Orbital Colony teaches Go through missions: the player commands drones on planet Kepler-G by writing real Go, which is compiled and run with the real `go` toolchain. Success is a player who keeps going chapter after chapter (16 chapters, 45 missions) because each mission is a challenge worth solving, and who ends up writing idiomatic, concurrent, tested Go.

## Positioning

The code is real Go, built by the real compiler: real compiler errors (with Arabic explanations), real goroutines running in parallel wall-clock time, the real race detector, real `go test` with mutation testing and benchmarks. Missions are randomized each run, concepts are enforced by AST checks, and stars reward efficient solutions, so copying a fixed answer does not work.

## Operating Context

- A local Go server (`go run ./cmd/orbital`) serves the web client on localhost and runs player code.
- Each mission: read the lesson, write code, press Run (or save in VS Code, which syncs), watch the drone animation, read the console, earn 1–3 stars.
- Lab missions (chapter 16) have no drone; they show `go test` steps.
- Progress is stored in the browser (localStorage); missions unlock in order.

## Capabilities and Constraints

- Lessons use a small markdown dialect (`~code~`, `~~~` fences, tables).
- The editor is a textarea with a syntax-highlight overlay; VS Code linking gives full gopls tooling.
- The world is drawn on a canvas: tiles (regolith, soil, ice, rock), plants, crystals, beacons, multiple coloured drones, storms.
- Sound effects are wanted, generated in code (no audio files), with a mute control.
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
