"use strict";

function renderCh15Frame(frame, { element, button, kit, localize }) {
  const view = element("div", "ch15-trace-frame");
  const caption = element("p", "ch15-trace-caption", localize(frame.caption));
  caption.setAttribute("role", "status");
  view.appendChild(caption);
  frame.rows.forEach((row) => {
    view.appendChild(kit.boxes({ label: localize(row.label), values: row.values, helpers: { element, button } }));
  });
  return view;
}

function renderCh15Trace(block, { element, button, T, kit }) {
  const config = block.sandbox.config;
  const localize = (pair) => T(pair.ar, pair.en);
  const stepper = kit.stepper({
    frames: config.frames,
    labels: {
      previous: localize(config.labels.previous),
      next: localize(config.labels.next),
      reset: localize(config.labels.reset),
    },
    helpers: { element, button },
    renderFrame: (frame) => renderCh15Frame(frame, { element, button, kit, localize }),
  });
  const panel = element("section", "lesson-sandbox ch15-algorithm-trace");
  panel.appendChild(stepper.node);
  return panel;
}

Lessons.registerSandbox("ch15-algorithm-trace", renderCh15Trace);
