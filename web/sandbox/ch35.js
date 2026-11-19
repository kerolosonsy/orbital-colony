"use strict";

Lessons.registerSandbox("ch35-event-trace", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (pair) => T(pair.ar, pair.en);
  const labels = config.labels;
  const variants = config.variants;
  const panel = element("section", "lesson-sandbox ch35-event-trace");
  const rows = element("div", "ch35-trace-rows");
  const caption = element("p", "ch35-trace-caption");
  caption.setAttribute("role", "status");
  caption.setAttribute("aria-live", "polite");
  let variant = 0;
  let index = 0;

  function render() {
    const frame = variants[variant][index];
    rows.replaceChildren(...frame.rows.map((row) => {
      const node = element("div", "ch35-trace-row");
      node.append(element("span", "ch35-trace-label", text(row.label)), element("bdi", "ch35-trace-value", text(row.value)));
      return node;
    }));
    caption.textContent = text(frame.caption);
  }

  const controls = element("div", "ch35-trace-controls");
  controls.append(
    button(text(labels.advance), () => { index = (index + 1) % variants[variant].length; render(); }),
    button(text(labels.again), () => { variant = 1 - variant; index = 0; render(); }),
    button(text(labels.reset), () => { variant = 0; index = 0; render(); }),
  );
  panel.append(rows, caption, controls);
  render();
  return panel;
});
