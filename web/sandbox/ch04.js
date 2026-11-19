"use strict";

function ch04Text(pair, T) {
  return T(pair.ar, pair.en);
}

Lessons.registerSandbox("ch04-if-branch", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const label = (key) => ch04Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch04-if-branch");
  const state = element("div", "ch04-if-state");
  const caption = element("p", "ch04-if-caption");
  caption.setAttribute("role", "status");
  const controls = element("div", "lesson-controls ch04-if-controls");
  let selected = null;
  const choiceButtons = config.choices.map((choice) => button(ch04Text(choice.label, T), () => { selected = choice; update(); }));
  const reset = button(label("reset"), () => { selected = null; update(); });
  controls.append(...choiceButtons, reset);
  panel.append(state, caption, controls);

  function update() {
    const condition = element("p", "ch04-if-value");
    condition.append(element("span", "ch04-if-label", label("condition")), element("bdi", "ch04-if-data", selected ? label(selected.ready ? "ready" : "notReady") : label("choose")));
    const body = element("p", "ch04-if-body");
    body.append(element("span", "ch04-if-label", label("body")), element("span", "ch04-if-data", selected ? ch04Text(selected.bodyState, T) : label("choose")));
    state.replaceChildren(condition, body);
    caption.textContent = selected ? ch04Text(selected.caption, T) : label("choose");
    choiceButtons.forEach((choiceButton, index) => choiceButton.setAttribute("aria-pressed", String(config.choices[index] === selected)));
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch04-switch-case", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const label = (key) => ch04Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch04-switch-case");
  const state = element("div", "ch04-switch-state");
  const caption = element("p", "ch04-switch-caption");
  caption.setAttribute("role", "status");
  let selected = null;
  const controls = element("div", "lesson-controls ch04-switch-controls");
  config.choices.forEach((choice) => controls.appendChild(button(ch04Text(choice.label, T), () => { selected = choice; update(); })));
  panel.append(state, caption, controls);

  function update() {
    const signal = element("p", "ch04-switch-value");
    signal.append(element("span", "ch04-switch-label", label("signal")), element("bdi", "ch04-switch-data", selected ? ch04Text(selected.value, T) : label("choose")));
    const selectedCase = element("p", "ch04-switch-case-value");
    selectedCase.append(element("span", "ch04-switch-label", label("case")), element("bdi", "ch04-switch-data", selected ? ch04Text(selected.case, T) : label("choose")));
    const result = element("p", "ch04-switch-result");
    result.append(element("span", "ch04-switch-label", label("result")), element("bdi", "ch04-switch-data", selected ? ch04Text(selected.result, T) : label("choose")));
    state.replaceChildren(signal, selectedCase, result);
    caption.textContent = selected ? ch04Text(selected.caption, T) : label("choose");
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch04-ordered-branches", (block, helpers) => {
  const { element, button, T } = helpers;
  const config = block.sandbox.config;
  const label = (key) => ch04Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch04-ordered-branches");
  const state = element("div", "ch04-ordered-state");
  const caption = element("p", "ch04-ordered-caption");
  caption.setAttribute("role", "status");
  const controls = element("div", "lesson-controls ch04-ordered-controls");
  let index = 0;
  const next = button(label("next"), () => { index = Math.min(index + 1, config.frames.length - 1); update(); });
  const reset = button(label("reset"), () => { index = 0; update(); });
  controls.append(next, reset);
  panel.append(state, caption, controls);

  function update() {
    const frame = config.frames[index];
    const value = element("p", "ch04-ordered-hours");
    value.append(element("span", "ch04-ordered-label", label("hours")), element("bdi", "ch04-ordered-data", String(frame.hours)));
    const checks = element("div", "ch04-ordered-checks");
    [["firstCheck", frame.first], ["secondCheck", frame.second]].forEach(([key, result]) => {
      const row = element("p", "ch04-ordered-check");
      const condition = element("span", "ch04-ordered-condition", label(key));
      condition.setAttribute("dir", "ltr");
      row.append(condition, element("span", "ch04-ordered-outcome", result === null ? label("skipped") : label(result ? "true" : "false")));
      checks.appendChild(row);
    });
    const result = element("p", "ch04-ordered-result");
    result.append(element("span", "ch04-ordered-label", label("result")), element("bdi", "ch04-ordered-data", ch04Text(frame.result, T)));
    state.replaceChildren(value, checks, result);
    caption.textContent = ch04Text(frame.caption, T);
    next.disabled = index === config.frames.length - 1;
  }

  update();
  return panel;
});
