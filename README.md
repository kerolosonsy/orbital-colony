# Orbital Colony

Orbital Colony is a bilingual game for learning Go by programming drones and building tools for a space colony. Start with function calls and control flow, then work through concurrency, testing, HTTP, databases, and a capstone API. Play in English or Egyptian Arabic. [الدليل بالعربي](docs/PLAYING.ar.md)

## Play locally

Install [Go 1.25 or newer](https://go.dev/dl/), then run these commands from a terminal.

```sh
git clone https://github.com/kerolosonsy/orbital-colony.git
cd orbital-colony
go run ./cmd/orbital
```

Open [127.0.0.1:7777](http://127.0.0.1:7777), choose **First Signal**, complete the starter, and select **Run**. Use the language button to switch languages. The game saves progress locally.

Orbital Colony uses your local Go compiler and runs programs with your desktop user's access. Run code you trust and keep the game on your own computer. [Security and reporting](SECURITY.md)

### What you can play

- **Drone missions.** Write Go to explore maps, harvest materials, and launch automated colony systems. Step through your program and compare code, variables, and map state.
- **Coding labs.** Build and test complete programs with the Go toolchain. Some labs use race detection, fuzzing, benchmarks, or hidden test cases.
- **Guided lessons.** Work through explanations, small interactive exercises, hints, and chapter quizzes. The curriculum has 181 missions across 38 chapters.

Use Go **1.27.1 or newer** and run `make deps` to prepare the optional third-party lab dependencies. For labs that need local services, start Docker and run `make services-up`. [Arabic play and setup guide](docs/PLAYING.ar.md)

## Browser demo

**First Signal** lets you assemble Go direction calls and see how their order changes the drone's position. It runs as a static lesson in your browser. Build it with `make demo` and open `dist/index.html`. [Demo and GitHub Pages setup](docs/DEPLOY.md)

## Keep learning Go

Use the guided missions for practice, then design and implement an independent project in [Golang-100-projects](https://github.com/kerolosonsy/Golang-100-projects). The challenge follows its own learner-written implementation rules. Start with [Project 001, Hello CLI](https://github.com/kerolosonsy/Golang-100-projects/blob/main/golang-100-projects/01-foundations/001_hello_cli/README.md) to apply the basics of `main`, output, and error handling.

## Contribute

Pick an issue or propose one. You can improve a translation, review a lesson, fix an accessibility problem, or develop simulation behavior.

| Path | What belongs there |
| --- | --- |
| `colony/levels/` | Mission rules, starter code, lessons, and quizzes |
| `colony/internal/sim/` and `colony/econ/` | Drone simulation and colony economy |
| `cmd/orbital/` | Local HTTP server, Go runner, debugger, and profiles |
| `web/` | Game UI and interactive lesson widgets |
| `cmd/lesson-demo/` | The static First Signal demo |

[Contribution guide](CONTRIBUTING.md) · [Code of conduct](CODE_OF_CONDUCT.md) · [Security policy](SECURITY.md)

## License

Orbital Colony's source and original documentation are released under the [MIT License](LICENSE). Third-party dependencies and externally hosted fonts retain their respective licenses.
