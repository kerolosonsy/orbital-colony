"use strict";

function ch05Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch05Panel(element, kind) {
  return element("section", `lesson-sandbox ch05-sandbox ch05-${kind}`);
}

function ch05Source(config, element, T) {
  return element("p", "ch05-sandbox-source", ch05Text(config.source, T));
}

function ch05Caption(config, element, T) {
  const caption = element("p", "ch05-sandbox-caption", ch05Text(config.caption, T));
  caption.setAttribute("role", "status");
  caption.setAttribute("aria-live", "polite");
  return caption;
}

Lessons.registerSandbox("ch05-params", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch05Text(config.labels[key], T);
  const panel = ch05Panel(element, "params");
  panel.appendChild(ch05Source(config, element, T));
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const content = element("div", "ch05-sandbox-frame");
      content.appendChild(ch05Caption(frame, element, T));
      content.appendChild(kit.valueRow({ label: label("call"), displayValue: frame.call, helpers: { element } }).node);
      content.appendChild(kit.valueRow({ label: label("inside"), displayValue: frame.inside, helpers: { element } }).node);
      content.appendChild(kit.valueRow({ label: label("result"), displayValue: frame.result, helpers: { element } }).node);
      return content;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch05-error-value", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch05Text(config.labels[key], T);
  const panel = ch05Panel(element, "error-value");
  panel.appendChild(ch05Source(config, element, T));
  const caption = element("p", "ch05-sandbox-caption");
  caption.setAttribute("role", "status");
  caption.setAttribute("aria-live", "polite");
  const error = kit.valueRow({ label: label("error"), displayValue: "", helpers: { element } });
  const value = kit.valueRow({ label: label("value"), displayValue: config.value, helpers: { element } });
  const action = kit.valueRow({ label: label("action"), displayValue: "", helpers: { element } });
  let hasError = true;
  function update() {
    const state = hasError ? config.failure : config.success;
    caption.textContent = ch05Text(state.caption, T);
    error.set(ch05Text(state.error, T));
    action.set(ch05Text(state.action, T));
  }
  const toggle = kit.toggle({
    label: label("outcome"),
    onLabel: label("present"),
    offLabel: label("absent"),
    isOn: hasError,
    onChange(next) { hasError = next; update(); },
    helpers: { element, button },
  });
  panel.append(toggle.node, caption, error.node, value.node, action.node);
  update();
  return panel;
});

Lessons.registerSandbox("ch05-wrap", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch05Text(config.labels[key], T);
  const panel = ch05Panel(element, "wrap");
  panel.appendChild(ch05Source(config, element, T));
  const caption = element("p", "ch05-sandbox-caption");
  caption.setAttribute("role", "status");
  caption.setAttribute("aria-live", "polite");
  const message = kit.valueRow({ label: label("message"), displayValue: "", helpers: { element } });
  const identity = kit.valueRow({ label: label("identity"), displayValue: "", helpers: { element } });
  let isWrapped = true;
  function update() {
    const state = isWrapped ? config.wrapped : config.plain;
    caption.textContent = ch05Text(state.caption, T);
    message.set(state.message);
    identity.set(state.identity);
  }
  const toggle = kit.toggle({
    label: label("format"), onLabel: config.labels.wrap, offLabel: config.labels.plain, isOn: isWrapped,
    onChange(next) { isWrapped = next; update(); },
    helpers: { element, button },
  });
  panel.append(toggle.node, caption, message.node, identity.node);
  update();
  return panel;
});

Lessons.registerSandbox("ch05-defer", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch05Text(config.labels[key], T);
  const panel = ch05Panel(element, "defer");
  panel.appendChild(ch05Source(config, element, T));
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const content = element("div", "ch05-sandbox-frame");
      content.appendChild(ch05Caption(frame, element, T));
      content.appendChild(kit.valueRow({ label: label("output"), displayValue: frame.output, helpers: { element } }).node);
      return content;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch05-closure", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch05Text(config.labels[key], T);
  const panel = ch05Panel(element, "closure");
  panel.appendChild(ch05Source(config, element, T));
  const status = element("p", "ch05-sandbox-caption", ch05Text(config.states.idle, T));
  status.setAttribute("role", "status");
  status.setAttribute("aria-live", "polite");
  let count = config.initialCount;
  const countRow = kit.valueRow({ label: label("count"), displayValue: count, helpers: { element } });
  const store = button(label("store"), () => {
    status.textContent = ch05Text(config.states.stored, T);
    store.setAttribute("aria-pressed", "true");
    call.disabled = false;
  });
  const call = button(label("call"), () => {
    count += 1;
    countRow.set(count);
    status.textContent = ch05Text(config.states.called, T);
  });
  call.disabled = true;
  store.setAttribute("aria-pressed", "false");
  const reset = button(label("reset"), () => {
    count = config.initialCount;
    countRow.set(count);
    status.textContent = ch05Text(config.states.idle, T);
    store.setAttribute("aria-pressed", "false");
    call.disabled = true;
  });
  const controls = element("div", "lesson-controls ch05-sandbox-controls");
  controls.append(store, call, reset);
  panel.append(status, countRow.node, controls);
  return panel;
});
