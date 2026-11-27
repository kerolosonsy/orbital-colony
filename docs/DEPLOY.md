# Run the game and publish the browser lesson

## Local game

Start from the repository root with `go run ./cmd/orbital`. The app listens at [127.0.0.1:7777](http://127.0.0.1:7777) by default and runs learner programs on your computer. Connecting the Go server to a public network through a proxy or tunnel would expose that execution capability. Share the static browser lesson with other people.

Optional local accounts can be enabled with `make run PUBLIC=http://127.0.0.1:7777`. Google sign-in uses a Web OAuth client with redirect URI `http://127.0.0.1:7777/auth/google/callback`. Set `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` in your shell or a local `.orbital.env` file. Existing email and password accounts keep their existing sign-in method. Use the same loopback origin each time to retain session cookies.

## Build First Signal

```sh
make demo
```

Open `dist/index.html` to try the same static assets used for deployment. The generator reads the First Signal map and localized text from the Go curriculum. The browser plays a bounded sequence of direction calls using JavaScript. The demo needs no accounts, runtime, keys, analytics, or API calls.

## Publish to GitHub Pages

In the repository's **Settings → Pages**, select **GitHub Actions** as the source. Run the **Browser demo** workflow from the default branch to publish the first build. It also runs on pushes to the default branch when the demo or curriculum changes. Verify that the workflow reports a deployment URL and that the site loads before sharing it. [GitHub Pages workflow documentation](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages)
