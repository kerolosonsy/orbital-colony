"use strict";

Lessons.registerSandbox("ch34-storage-flow", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (pair) => T(pair.ar, pair.en);
  const panel = element("section", "lesson-sandbox ch34-storage-flow");
  const title = element("h3", "ch34-flow-title", text(config.title));
  const caption = element("p", "ch34-flow-caption");
  const metrics = element("div", "ch34-flow-metrics");
  const metricRows = [0, 1].map(() => {
    const row = element("div", "ch34-flow-metric");
    const name = element("span", "ch34-flow-label");
    const value = element("bdi", "ch34-flow-value");
    value.setAttribute("dir", "auto");
    row.append(name, value);
    metrics.appendChild(row);
    return { name, value };
  });
  let frameIndex = 0;
  let advance;
  const update = () => {
    const frame = config.frames[frameIndex];
    if (advance) advance.disabled = frameIndex === config.frames.length - 1;
    caption.textContent = text(frame.caption);
    metricRows[0].name.textContent = text(frame.leftLabel);
    metricRows[0].value.textContent = text(frame.leftValue);
    metricRows[1].name.textContent = text(frame.rightLabel);
    metricRows[1].value.textContent = text(frame.rightValue);
  };
  const controls = element("div", "ch34-flow-controls");
  controls.append(
    advance = button(text(config.labels.advance), () => {
      frameIndex = Math.min(config.frames.length - 1, frameIndex + 1);
      update();
    }),
    button(text(config.labels.reset), () => {
      frameIndex = 0;
      update();
    }),
  );
  panel.append(title, caption, metrics, controls);
  update();
  return panel;
});
