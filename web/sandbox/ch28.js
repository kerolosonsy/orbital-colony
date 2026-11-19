"use strict";

function ch28Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch28Row(element, label, value) {
  const row = element("div", "ch28-row");
  const output = value && typeof value === "object" ? value : element("bdi", "ch28-value", String(value));
  row.append(element("span", "ch28-label", label), output);
  return row;
}

function ch28Caption(element, T, pair) {
  const node = element("p", "ch28-caption", ch28Text(pair, T));
  node.setAttribute("role", "status");
  node.setAttribute("aria-live", "polite");
  return node;
}

Lessons.registerSandbox("ch28-gin-flow", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch28-sandbox ch28-gin-flow");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: {
      previous: ch28Text(labels.previous, T),
      next: ch28Text(labels.next, T),
      reset: ch28Text(labels.reset, T),
    },
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch28-frame");
      view.append(
        ch28Row(element, ch28Text(labels.step, T), frame.step),
        ch28Row(element, ch28Text(labels.value, T), frame.value),
        ch28Caption(element, T, frame.caption),
      );
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch28-fiber-copy", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch28-sandbox ch28-fiber-copy");
  const buffer = element("bdi", "ch28-value");
  const saved = element("bdi", "ch28-value", config.values[0]);
  let changed = false;
  function update() {
    buffer.textContent = config.values[changed ? 1 : 0];
  }
  panel.append(
    ch28Row(element, ch28Text(labels.buffer, T), buffer),
    ch28Row(element, ch28Text(labels.saved, T), saved),
    ch28Caption(element, T, labels.caption),
    button(ch28Text(labels.reuse, T), () => { changed = !changed; update(); }),
  );
  update();
  return panel;
});

Lessons.registerSandbox("ch28-json-shape", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch28-sandbox ch28-json-shape");
  const value = element("bdi", "ch28-value");
  const type = element("bdi", "ch28-value");
  const number = element("bdi", "ch28-value");
  const caption = ch28Caption(element, T, config.samples[0].caption);
  let index = 0;
  function update() {
    const sample = config.samples[index];
    value.textContent = sample.value;
    type.textContent = sample.type;
    number.textContent = sample.number;
    caption.textContent = ch28Text(sample.caption, T);
  }
  panel.append(
    ch28Row(element, ch28Text(labels.field, T), value),
    ch28Row(element, ch28Text(labels.jsonType, T), type),
    ch28Row(element, ch28Text(labels.number, T), number),
    caption,
    button(ch28Text(labels.next, T), () => { index = (index + 1) % config.samples.length; update(); }),
  );
  update();
  return panel;
});

Lessons.registerSandbox("ch28-shutdown-order", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch28-sandbox ch28-shutdown-order");
  let orderIndex = 0;
  const timeline = kit.logList({
    label: ch28Text(labels.timeline, T),
    entries: config.orders[orderIndex].map((entry) => ch28Text(entry, T)),
    helpers: { element },
  });
  panel.append(
    timeline.node,
    ch28Caption(element, T, labels.caption),
    button(ch28Text(labels.runAgain, T), () => {
      orderIndex = (orderIndex + 1) % config.orders.length;
      timeline.set(config.orders[orderIndex].map((entry) => ch28Text(entry, T)));
    }),
  );
  return panel;
});

Lessons.registerSandbox("ch28-response-commit", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch28-sandbox ch28-response-commit");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: {
      previous: ch28Text(labels.previous, T),
      next: ch28Text(labels.next, T),
      reset: ch28Text(labels.reset, T),
    },
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch28-frame");
      view.append(
        ch28Row(element, ch28Text(labels.status, T), frame.status),
        ch28Row(element, ch28Text(labels.header, T), frame.header),
        ch28Row(element, ch28Text(labels.body, T), frame.body),
        ch28Row(element, ch28Text(labels.state, T), frame.state),
        ch28Caption(element, T, frame.caption),
      );
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});
