"use strict";

function ch03Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch03Frames(block, helpers, kind, renderState) {
  const { element, button, T } = helpers;
  const config = block.sandbox.config;
  const label = (key) => ch03Text(config.labels[key], T);
  const panel = element("section", `lesson-sandbox ${kind}`);
  const state = element("div", `${kind}-state`);
  const caption = element("p", `${kind}-caption`);
  caption.setAttribute("role", "status");
  const controls = element("div", `lesson-controls ${kind}-controls`);
  let index = 0;
  const next = button(label("next"), () => { index = Math.min(index + 1, config.frames.length - 1); update(); });
  const reset = button(label("reset"), () => { index = 0; update(); });
  controls.append(next, reset);
  panel.append(state, caption, controls);

  function update() {
    const frame = config.frames[index];
    state.replaceChildren(renderState(frame, config, label, element, T));
    caption.textContent = ch03Text(frame.caption, T);
    next.disabled = index === config.frames.length - 1;
  }

  update();
  return panel;
}

Lessons.registerSandbox("ch03-repeat-counter", (block, helpers) => ch03Frames(block, helpers, "ch03-repeat-counter", (frame, config, label, element) => {
  const state = element("div", "ch03-repeat-state");
  const counter = element("p", "ch03-repeat-counter-value");
  counter.append(element("span", "ch03-repeat-label", label("counter")), element("bdi", "ch03-repeat-number", String(frame.i)));
  const output = element("div", "ch03-repeat-output");
  output.appendChild(element("span", "ch03-repeat-label", label("prints")));
  const values = element("div", "ch03-repeat-values");
  values.setAttribute("dir", "ltr");
  frame.prints.forEach((value) => values.appendChild(element("span", "ch03-repeat-value", String(value))));
  output.appendChild(values);
  state.append(counter, output);
  return state;
}));

Lessons.registerSandbox("ch03-nested-grid", (block, helpers) => ch03Frames(block, helpers, "ch03-nested-grid", (frame, config, label, element) => {
  const state = element("div", "ch03-nested-state");
  const position = element("p", "ch03-nested-position");
  position.append(element("span", "ch03-nested-label", label("position")), element("bdi", "ch03-nested-coordinates", `${frame.row}, ${frame.column}`));
  const grid = element("div", "ch03-nested-cells");
  grid.setAttribute("dir", "ltr");
  const visited = new Set(frame.visited);
  const active = frame.row * config.columns + frame.column;
  for (let cell = 0; cell < config.rows * config.columns; cell++) {
    const node = element("span", "ch03-nested-cell", `${Math.floor(cell / config.columns)},${cell % config.columns}`);
    node.classList.toggle("visited", visited.has(cell));
    node.classList.toggle("active", cell === active);
    grid.appendChild(node);
  }
  state.append(position, grid);
  return state;
}));

Lessons.registerSandbox("ch03-condition-loop", (block, helpers) => ch03Frames(block, helpers, "ch03-condition-loop", (frame, config, label, element, T) => {
  const state = element("div", "ch03-condition-state");
  const count = element("p", "ch03-condition-count");
  count.append(element("span", "ch03-condition-label", label("count")), element("bdi", "ch03-condition-number", String(frame.count)));
  const condition = element("p", "ch03-condition-value");
  condition.append(element("span", "ch03-condition-label", label("condition")), element("span", "ch03-condition-status", label(frame.running ? "running" : "stopped")));
  const progress = element("div", "ch03-condition-progress");
  progress.setAttribute("aria-label", label("count"));
  for (let step = 0; step < config.limit; step++) {
    const dot = element("span", "ch03-condition-dot");
    dot.classList.toggle("complete", step < frame.runs);
    progress.appendChild(dot);
  }
  state.append(count, condition, progress);
  return state;
}));
