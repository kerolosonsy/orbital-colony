"use strict";

Lessons.registerSandbox("ch33-http-flow", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (pair) => T(pair.ar, pair.en);
  const panel = element("section", "lesson-sandbox ch33-http-flow");
  const title = element("h3", "ch33-flow-title", text(config.title));
  const runLabel = element("p", "ch33-flow-run");
  const caption = element("p", "ch33-flow-caption");
  const metrics = element("div", "ch33-flow-metrics");
  const metricRows = [0, 1].map(() => {
    const row = element("div", "ch33-flow-metric");
    const name = element("span", "ch33-flow-label");
    const value = element("bdi", "ch33-flow-value");
    value.setAttribute("dir", "auto");
    row.append(name, value);
    metrics.appendChild(row);
    return { name, value };
  });
  let runIndex = 0;
  let frameIndex = 0;
  const update = () => {
    const run = config.runs[runIndex];
    const frame = run.frames[frameIndex];
    advance.disabled = frameIndex === run.frames.length - 1;
    runLabel.textContent = text(run.label);
    caption.textContent = text(frame.caption);
    metricRows[0].name.textContent = text(frame.leftLabel);
    metricRows[0].value.textContent = text(frame.leftValue);
    metricRows[1].name.textContent = text(frame.rightLabel);
    metricRows[1].value.textContent = text(frame.rightValue);
  };
  const advance = button(text(config.labels.advance), () => {
    frameIndex = Math.min(config.runs[runIndex].frames.length - 1, frameIndex + 1);
    update();
  });
  const controls = element("div", "ch33-flow-controls");
  controls.appendChild(advance);
  if (config.runs.length > 1) {
    controls.appendChild(button(text(config.labels.repeat), () => {
      runIndex = (runIndex + 1) % config.runs.length;
      frameIndex = 0;
      update();
    }));
  }
  controls.appendChild(button(text(config.labels.reset), () => {
    runIndex = 0;
    frameIndex = 0;
    update();
  }));
  panel.append(title, runLabel, caption, metrics, controls);
  update();
  return panel;
});
