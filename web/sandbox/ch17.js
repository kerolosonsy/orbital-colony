"use strict";

function renderCh17Sandbox(block, { element, button, T }) {
  const config = block.sandbox.config;
  const { labels, data, mode } = config;
  const label = (key) => T(labels[key].ar, labels[key].en);
  const value = (pair) => T(pair.ar, pair.en);
  const panel = element("section", `lesson-sandbox ch17-sandbox ch17-${mode}`);
  panel.setAttribute("dir", T("rtl", "ltr"));
  const state = element("div", "ch17-sandbox-state");
  state.setAttribute("role", "status");
  state.setAttribute("aria-live", "polite");
  const controls = element("div", "ch17-sandbox-controls");
  panel.append(state, controls);

  function row(name, shownValue) {
    const line = element("div", "ch17-sandbox-row");
    line.append(element("span", "ch17-sandbox-label", name));
    const output = element("bdi", "ch17-sandbox-value", String(shownValue));
    output.setAttribute("dir", "ltr");
    line.appendChild(output);
    return line;
  }

  function caption(text) {
    const node = element("p", "ch17-sandbox-caption", text);
    return node;
  }

  if (mode === "coordinate") {
    let layer = 0;
    const action = button(label("action"), () => {
      layer = (layer + 1) % data.layers.length;
      update();
    });
    controls.appendChild(action);
    function update() {
      state.replaceChildren(
        element("p", "ch17-sandbox-heading", label("location")),
        row(label("x"), data.x),
        row(label("y"), data.y),
        row(label("z"), data.layers[layer]),
        caption(label("caption")),
      );
    }
    update();
  } else if (mode === "volume") {
    let selected = 0;
    const action = button(label("action"), () => {
      selected = (selected + 1) % data.cells.length;
      update();
    });
    controls.appendChild(action);
    function update() {
      const cell = data.cells[selected];
      const indices = `[${cell.z}][${cell.x}][${cell.y}]`;
      state.replaceChildren(
        element("p", "ch17-sandbox-heading", label("indices")),
        row(indices, cell.value),
        caption(label("caption")),
      );
    }
    update();
  } else if (mode === "queue") {
    let queue = [...data.items];
    let taken = null;
    const action = button(label("action"), () => {
      if (queue.length) taken = queue.shift();
      update();
      if (!queue.length) action.disabled = true;
    });
    controls.appendChild(action);
    function update() {
      const list = element("ol", "ch17-sandbox-queue");
      list.setAttribute("aria-label", label("queue"));
      for (const item of queue) list.appendChild(element("li", "ch17-sandbox-item", value(item)));
      state.replaceChildren(
        row(label("current"), taken ? value(taken) : label("empty")),
        element("p", "ch17-sandbox-heading", label("queue")),
        queue.length ? list : element("p", "ch17-sandbox-empty", label("empty")),
      );
    }
    update();
  } else if (mode === "key") {
    let probeIndex = 0;
    const action = button(label("action"), () => {
      probeIndex = (probeIndex + 1) % data.probes.length;
      update();
    });
    controls.appendChild(action);
    const displayKey = (key) => `{X: ${key.x}, Z: ${key.z}}`;
    function update() {
      const probe = data.probes[probeIndex];
      const matches = probe.x === data.stored.x && probe.z === data.stored.z;
      state.replaceChildren(
        row(label("stored"), displayKey(data.stored)),
        row(label("lookup"), displayKey(probe)),
        row(label("result"), label(matches ? "found" : "missing")),
      );
    }
    update();
  }

  return panel;
}

Lessons.registerSandbox("ch17-layer-coordinate", renderCh17Sandbox);
Lessons.registerSandbox("ch17-volume-index", renderCh17Sandbox);
Lessons.registerSandbox("ch17-fifo-queue", renderCh17Sandbox);
Lessons.registerSandbox("ch17-comparable-key", renderCh17Sandbox);
