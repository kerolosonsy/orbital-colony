"use strict";

function ch27Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch27Row(element, label, value) {
  const row = element("div", "ch27-row");
  const output = value && typeof value === "object" ? value : element("bdi", "ch27-value", String(value));
  row.append(element("span", "ch27-label", label), output);
  return row;
}

function ch27Caption(element, T, pair) {
  const node = element("p", "ch27-caption", ch27Text(pair, T));
  node.setAttribute("role", "status");
  node.setAttribute("aria-live", "polite");
  return node;
}

Lessons.registerSandbox("ch27-signed-pass", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch27-sandbox ch27-signed-pass");
  let messageIndex = 0;
  let keyIndex = 0;
  const message = element("bdi", "ch27-value");
  const key = element("bdi", "ch27-value");
  const tag = element("bdi", "ch27-value");
  const caption = ch27Caption(element, T, labels.caption);
  function update() {
    message.textContent = config.messages[messageIndex];
    key.textContent = config.keys[keyIndex];
    tag.textContent = config.tags[messageIndex * config.keys.length + keyIndex];
  }
  panel.append(
    ch27Row(element, ch27Text(labels.payload, T), message),
    ch27Row(element, ch27Text(labels.key, T), key),
    ch27Row(element, ch27Text(labels.tag, T), tag),
    caption,
  );
  const actions = element("div", "ch27-actions");
  actions.append(
    button(ch27Text(labels.changeMessage, T), () => { messageIndex = (messageIndex + 1) % config.messages.length; update(); }),
    button(ch27Text(labels.rotateKey, T), () => { keyIndex = (keyIndex + 1) % config.keys.length; update(); }),
  );
  panel.appendChild(actions);
  update();
  return panel;
});

Lessons.registerSandbox("ch27-session-lookup", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch27-sandbox ch27-session-lookup");
  const id = element("bdi", "ch27-value");
  const result = element("bdi", "ch27-value");
  const caption = ch27Caption(element, T, config.lookups[0].caption);
  let index = 0;
  function update() {
    const lookup = config.lookups[index];
    id.textContent = lookup.id;
    result.textContent = lookup.result;
    caption.textContent = ch27Text(lookup.caption, T);
  }
  panel.append(
    ch27Row(element, ch27Text(labels.browser, T), id),
    ch27Row(element, ch27Text(labels.server, T), result),
    caption,
    button(ch27Text(labels.next, T), () => { index = (index + 1) % config.lookups.length; update(); }),
  );
  update();
  return panel;
});

Lessons.registerSandbox("ch27-preflight", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch27-sandbox ch27-preflight");
  const method = element("bdi", "ch27-value");
  const stage = element("bdi", "ch27-value");
  const status = element("bdi", "ch27-value");
  const calls = element("bdi", "ch27-value");
  const caption = ch27Caption(element, T, config.requests[0].caption);
  let index = 0;
  function update() {
    const request = config.requests[index];
    method.textContent = request.method;
    stage.textContent = request.stage;
    status.textContent = request.status;
    calls.textContent = request.calls;
    caption.textContent = ch27Text(request.caption, T);
  }
  panel.append(
    ch27Row(element, ch27Text(labels.method, T), method),
    ch27Row(element, ch27Text(labels.stage, T), stage),
    ch27Row(element, ch27Text(labels.status, T), status),
    ch27Row(element, ch27Text(labels.innerCalls, T), calls),
    caption,
    button(ch27Text(labels.next, T), () => { index = (index + 1) % config.requests.length; update(); }),
  );
  update();
  return panel;
});

Lessons.registerSandbox("ch27-client-counts", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch27-sandbox ch27-client-counts");
  const counts = config.clients.map(() => 0);
  const values = config.clients.map(() => element("bdi", "ch27-value", "0"));
  const actions = element("div", "ch27-actions");
  panel.appendChild(ch27Row(element, ch27Text(labels.client, T), ch27Text(labels.count, T)));
  config.clients.forEach((client, index) => {
    panel.appendChild(ch27Row(element, client.address, values[index]));
    actions.appendChild(button(ch27Text(client.button, T), () => {
      counts[index] += 1;
      values[index].textContent = String(counts[index]);
    }));
  });
  panel.append(element("p", "ch27-caption", ch27Text(labels.caption, T)), actions);
  return panel;
});

Lessons.registerSandbox("ch27-base62-index", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const labels = config.labels;
  const panel = element("section", "lesson-sandbox ch27-sandbox ch27-base62-index");
  const stepper = kit.stepper({
    frames: config.lookups,
    labels: {
      previous: ch27Text(labels.previous, T),
      next: ch27Text(labels.next, T),
      reset: ch27Text(labels.reset, T),
    },
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch27-frame");
      view.append(
        ch27Row(element, ch27Text(labels.index, T), frame.index),
        ch27Row(element, ch27Text(labels.symbol, T), frame.symbol),
        ch27Caption(element, T, frame.caption),
      );
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});
