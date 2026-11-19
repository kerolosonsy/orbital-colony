"use strict";

function ch20Text(T, pair) {
  return typeof pair === "string" ? pair : T(pair.ar, pair.en);
}

Lessons.registerSandbox("ch20-bulletin-state", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch20Text(T, config.labels[key]);
  const panel = element("section", "lesson-sandbox ch20-bulletin-state");
  const source = kit.valueRow({ label: label("source"), displayValue: config.source, helpers: { element } });
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const caption = element("p", "ch20-playground-caption", ch20Text(T, frame.caption));
      const piece = kit.valueRow({ label: label("piece"), displayValue: ch20Text(T, frame.piece), helpers: { element } });
      const state = kit.valueRow({ label: label("state"), displayValue: ch20Text(T, frame.state), helpers: { element } });
      const output = kit.valueRow({ label: label("output"), displayValue: frame.output, helpers: { element } });
      const view = element("div", "ch20-bulletin-frame");
      view.append(caption, piece.node, state.node, output.node);
      return view;
    },
  });
  panel.append(source.node, stepper.node);
  return panel;
});

Lessons.registerSandbox("ch20-error-join", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch20Text(T, config.labels[key]);
  const panel = element("section", "lesson-sandbox ch20-error-join");
  const controls = element("div", "ch20-field-controls");
  const causes = kit.logList({ label: label("causes"), entries: [], helpers: { element } });
  const invalidFields = config.fields.map((field) => field.invalid);
  const refresh = () => {
    causes.set(config.fields.filter((field, index) => invalidFields[index]).map((field) => ch20Text(T, field.cause)));
  };
  config.fields.forEach((field, index) => {
    const toggle = kit.toggle({
      label: ch20Text(T, field.label),
      onLabel: label("invalid"),
      offLabel: label("valid"),
      isOn: field.invalid,
      helpers: { element, button },
      onChange(isInvalid) {
        invalidFields[index] = isInvalid;
        refresh();
      },
    });
    controls.appendChild(toggle.node);
  });
  panel.append(controls, causes.node);
  refresh();
  return panel;
});

Lessons.registerSandbox("ch20-lexer-offset", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch20Text(T, config.labels[key]);
  const panel = element("section", "lesson-sandbox ch20-lexer-offset");
  const source = kit.valueRow({ label: label("source"), displayValue: config.source, helpers: { element } });
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const caption = element("p", "ch20-playground-caption", ch20Text(T, frame.caption));
      const cursor = kit.valueRow({ label: label("cursor"), displayValue: frame.cursor, helpers: { element } });
      const token = kit.valueRow({ label: label("token"), displayValue: ch20Text(T, frame.token), helpers: { element } });
      const offset = kit.valueRow({ label: label("offset"), displayValue: frame.offset, helpers: { element } });
      const view = element("div", "ch20-lexer-frame");
      view.append(caption, cursor.node, token.node, offset.node);
      return view;
    },
  });
  panel.append(source.node, stepper.node);
  return panel;
});

Lessons.registerSandbox("ch20-expression-tree", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch20Text(T, config.labels[key]);
  const panel = element("section", "lesson-sandbox ch20-expression-tree");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const caption = element("p", "ch20-playground-caption", ch20Text(T, frame.caption));
      const expression = kit.valueRow({ label: label("expression"), displayValue: frame.expression, helpers: { element } });
      const value = kit.valueRow({ label: label("value"), displayValue: ch20Text(T, frame.value), helpers: { element } });
      const view = element("div", "ch20-expression-frame");
      view.append(caption, expression.node, value.node);
      return view;
    },
  });
  panel.append(stepper.node);
  return panel;
});
