# Lesson playgrounds

## Ownership

Each chapter owns `colony/levels/chNN_lessons.go` and its `chNN.js`, `chNN.css`, and `chNN.test.js` files here. Keep chapter work in those files so chapter agents can work in parallel. The shared page wiring, lesson renderer, kit, and harness are maintained centrally.

## Config and naming

Put language-neutral values and display strings in `LessonSandbox.Config`. Store localized text under string keys named `ar` and `en`, such as `labels.next.ar` and `labels.next.en`. The Arabic and English lesson configs must have the same keys and the same non-string values; string values may differ. Prefix playground IDs with `chNN-`, and name kinds `chNN-<name>`, for example `ch07-route-order`.

Use illustrative data. Never solve or reveal the mission in a playground. Keep it within the roughly 290px lesson panel, use the `T` helper for any renderer-owned labels, and check both Arabic RTL and English LTR layouts.

## Register, build, and test

Register a renderer with `Lessons.registerSandbox(kind, renderer)`. The renderer receives `(block, { element, button, T, kit })`; read data from `block.sandbox.config`. `element` builds nodes, `button` creates a real button, and `T(ar, en)` selects the active language. All widget labels should be passed in the active language.

The shared kit API is: `kit.stepper({frames, renderFrame, labels, helpers})` — previous/next/reset controls and a live frame counter; `kit.toggle({label, onLabel, offLabel, isOn, onChange, helpers})` — two-state button; `kit.slider({label, currentValue, min, max, step, onChange, helpers})` — range input and live value; `kit.valueRow({label, displayValue, helpers})` — labeled value row with change flash; `kit.boxes({label, values, helpers})` — visual boxes row; `kit.logList({label, entries, helpers})` — live log list; `kit.codeView({lines, activeLine, helpers})` — code with a zero-based highlighted line.

The chapter harness loads `lessons.js`, `kit.js`, and the chapter file, then renders a sandbox from a config fixture. Run the full gate with `node --test web/*.test.js web/sandbox/*.test.js`.

Example renderer for `ch07-route-order`:

```js
Lessons.registerSandbox("ch07-route-order", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const text = (pair) => T(pair.ar, pair.en);
  const label = (key) => text(config.labels[key]);
  const panel = element("section", "lesson-sandbox ch07-route-order");
  const stepper = kit.stepper({
    frames: config.routes,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(route, index) {
      return element("p", "route-frame", `${index + 1}. ${text(route)}`);
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});
```
