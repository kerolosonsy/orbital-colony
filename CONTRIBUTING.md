# Contributing to Orbital Colony

Lesson corrections, language improvements, accessibility fixes, tests, and simulation changes are welcome. Read the [Code of Conduct](CODE_OF_CONDUCT.md) and [Security Policy](SECURITY.md) before opening your first issue.

## Get started

1. Pick a [good first issue](https://github.com/kerolosonsy/orbital-colony/issues?q=is%3Aopen+label%3A%22good+first+issue%22), or propose a small change using an issue template.
2. Fork the repository, clone your fork, and create a branch for the issue.
3. Install Go 1.27.1 or newer for the full suite and Node.js 24 for web tests. Start the app with `go run ./cmd/orbital`.
4. Run the checks for your change. Include the commands and results in your pull request.

| Change | Checks |
| --- | --- |
| Go | `make test-quick` and `make vet` |
| Lessons or simulation | `cd colony && go test ./...` |
| Web client or browser demo | `make test-web` and `make demo` |
| Full local run | `make deps` then `make test` |

Apply `gofmt` to each changed Go file. The full suite needs permission to open loopback sockets for networking labs. Local service integrations use Docker and the separate setup in the [play guide](docs/PLAYING.ar.md).

## Improve a lesson

Start at `colony/levels/` and find the lesson by its mission ID. Give each learning step one clear goal. Keep English and Egyptian Arabic instructions, starter code, and hints aligned. Exercise new behavior against a seeded world or a small deterministic input. Verify the lesson at its actual entry point in the game.

Keep student-facing hints focused on reasoning. Test fixtures may contain reference implementations.

## Open a pull request

Link the issue, describe the learner-visible change, and list what you verified. Add before-and-after images for visual changes where useful. Share only code, text, and assets you have the right to contribute under the [project license](LICENSE). Credit external material and retain any required notices.
