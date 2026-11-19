"use strict";

function ch07Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch07Panel(element, kind) {
  return element("section", "lesson-sandbox ch07-sandbox ch07-" + kind);
}

function ch07Caption(element, text) {
  const caption = element("p", "ch07-caption", text);
  caption.setAttribute("role", "status");
  caption.setAttribute("aria-live", "polite");
  return caption;
}

Lessons.registerSandbox("ch07-point-fields", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = ch07Panel(element, "point-fields");
  const heading = element("p", "ch07-heading", ch07Text(config.labels.object, T));
  const fields = element("div", "ch07-fields");
  const output = kit.valueRow({
    label: ch07Text(config.labels.output, T),
    displayValue: ch07Text(config.labels.waiting, T),
    helpers: { element },
  });
  const caption = ch07Caption(element, ch07Text(config.labels.waiting, T));
  const controls = element("div", "ch07-controls");

  config.fields.forEach((field) => {
    const row = element("div", "ch07-field");
    row.append(
      element("bdi", "ch07-field-name", ch07Text(field.name, T)),
      element("bdi", "ch07-field-value", String(field.value)),
    );
    fields.appendChild(row);
    const read = button(ch07Text(field.read, T), () => {
      output.set(String(field.value));
      caption.textContent = ch07Text(field.caption, T);
      controls.querySelectorAll("button").forEach((item) => item.setAttribute("aria-pressed", String(item === read)));
    });
    read.setAttribute("aria-pressed", "false");
    controls.appendChild(read);
  });
  panel.append(heading, fields, output.node, caption, controls);
  return panel;
});

Lessons.registerSandbox("ch07-nav-pointer", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const text = (pair) => ch07Text(pair, T);
  const panel = ch07Panel(element, "nav-pointer");
  const position = kit.valueRow({ label: text(config.labels.position), displayValue: config.start, helpers: { element } });
  const caption = ch07Caption(element, text(config.captions.ready));
  const controls = element("div", "ch07-controls");
  let value = config.start;
  const update = () => position.set(String(value));
  controls.append(
    button(text(config.labels.step), () => {
      value += config.increment;
      update();
      caption.textContent = text(config.captions.updated);
    }),
    button(text(config.labels.reset), () => {
      value = config.start;
      update();
      caption.textContent = text(config.captions.reset);
    }),
  );
  panel.append(position.node, caption, controls);
  return panel;
});

Lessons.registerSandbox("ch07-manhattan", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const text = (pair) => ch07Text(pair, T);
  const panel = ch07Panel(element, "manhattan");
  let x = config.initialX;
  let y = config.initialY;
  const xGap = kit.valueRow({ label: text(config.labels.x), displayValue: x, helpers: { element } });
  const yGap = kit.valueRow({ label: text(config.labels.y), displayValue: y, helpers: { element } });
  const distance = kit.valueRow({ label: text(config.labels.total), displayValue: x + y, helpers: { element } });
  const caption = ch07Caption(element, text(config.captions.ready));
  const controls = element("div", "ch07-controls");
  const update = () => {
    xGap.set(String(x));
    yGap.set(String(y));
    distance.set(String(x + y));
  };
  controls.append(
    button(text(config.labels.addX), () => {
      x += 1;
      update();
      caption.textContent = text(config.captions.x);
    }),
    button(text(config.labels.addY), () => {
      y += 1;
      update();
      caption.textContent = text(config.captions.y);
    }),
    button(text(config.labels.reset), () => {
      x = config.initialX;
      y = config.initialY;
      update();
      caption.textContent = text(config.captions.reset);
    }),
  );
  panel.append(xGap.node, yGap.node, distance.node, caption, controls);
  return panel;
});

Lessons.registerSandbox("ch07-embedded-pack", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const text = (pair) => ch07Text(pair, T);
  const panel = ch07Panel(element, "embedded-pack");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: {
      previous: text(config.labels.previous),
      next: text(config.labels.next),
      reset: text(config.labels.reset),
    },
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch07-embedded-frame");
      view.append(
        element("pre", "ch07-state-code", text(frame.value)),
        ch07Caption(element, text(frame.caption)),
      );
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch07-range-questions", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const text = (pair) => ch07Text(pair, T);
  const panel = ch07Panel(element, "range-questions");
  const stepper = kit.stepper({
    frames: config.questions,
    labels: {
      previous: text(config.labels.previous),
      next: text(config.labels.next),
      reset: text(config.labels.reset),
    },
    helpers: { element, button },
    renderFrame(question) {
      const view = element("div", "ch07-question-frame");
      view.append(
        element("p", "ch07-prompt", text(question.prompt)),
        kit.valueRow({ label: text(config.labels.question), displayValue: String(question.index), helpers: { element } }).node,
        kit.valueRow({ label: text(config.labels.choices), displayValue: String(question.choices), helpers: { element } }).node,
        ch07Caption(element, text(question.caption)),
      );
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});
