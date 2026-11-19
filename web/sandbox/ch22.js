"use strict";

function ch22Text(T, pair) {
  return T(pair.ar, pair.en);
}

function ch22Label(T, config, key) {
  return ch22Text(T, config.labels[key]);
}

function ch22Panel(element, kind) {
  return element("section", `lesson-sandbox ${kind}`);
}

function ch22Caption(element, text) {
  const caption = element("p", "ch22-caption", text);
  caption.setAttribute("role", "status");
  caption.dir = "auto";
  return caption;
}

Lessons.registerSandbox("ch22-depth-tree", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch22Label(T, config, key);
  const panel = ch22Panel(element, "ch22-depth-tree");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch22-tree-frame");
      view.append(
        ch22Caption(element, ch22Text(T, frame.caption)),
        kit.valueRow({ label: label("depth"), displayValue: String(frame.depth), helpers: { element } }).node,
        kit.logList({ label: label("entries"), entries: frame.entries.map((entry) => ch22Text(T, entry)), helpers: { element } }).node,
      );
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch22-review-plan", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch22Label(T, config, key);
  const panel = ch22Panel(element, "ch22-review-plan");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch22-plan-frame");
      view.append(
        ch22Caption(element, ch22Text(T, frame.caption)),
        kit.logList({ label: label("plan"), entries: frame.moves.map((move) => ch22Text(T, move)), helpers: { element } }).node,
        kit.valueRow({ label: label("store"), displayValue: ch22Text(T, frame.store), helpers: { element } }).node,
      );
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch22-atomic-publish", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch22Label(T, config, key);
  const panel = ch22Panel(element, "ch22-atomic-publish");
  let targetValue = label("old");
  let stagingValue = label("empty");
  let phase = "initial";
  const caption = ch22Caption(element, "");
  const target = kit.valueRow({ label: label("target"), displayValue: targetValue, helpers: { element } });
  const staging = kit.valueRow({ label: label("staging"), displayValue: stagingValue, helpers: { element } });
  const progress = kit.valueRow({ label: label("phase"), displayValue: label(phase), helpers: { element } });
  const write = button(label("write"), () => { stagingValue = label("new"); phase = "staged"; update(); });
  const publish = button(label("publish"), () => { targetValue = label("new"); stagingValue = label("empty"); phase = "published"; update(); });
  const fail = button(label("fail"), () => { targetValue = label("old"); stagingValue = label("empty"); phase = "failed"; update(); });
  const reset = button(label("reset"), () => { targetValue = label("old"); stagingValue = label("empty"); phase = "initial"; update(); });
  const controls = element("div", "lesson-controls");
  controls.append(write, publish, fail, reset);
  function update() {
    caption.textContent = ch22Text(T, config.captions[phase]);
    target.set(targetValue);
    staging.set(stagingValue);
    progress.set(label(phase));
    publish.disabled = phase !== "staged";
    fail.disabled = phase !== "staged";
  }
  panel.append(caption, target.node, staging.node, progress.node, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch22-size-hash", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch22Label(T, config, key);
  const panel = ch22Panel(element, "ch22-size-hash");
  let pairIndex = 0;
  const firstFile = kit.valueRow({ label: label("fileA"), displayValue: "", helpers: { element } });
  const secondFile = kit.valueRow({ label: label("fileB"), displayValue: "", helpers: { element } });
  const sizes = kit.valueRow({ label: label("sizes"), displayValue: "", helpers: { element } });
  const filter = kit.valueRow({ label: label("sizeCheck"), displayValue: "", helpers: { element } });
  const hashes = kit.valueRow({ label: label("hashes"), displayValue: "", helpers: { element } });
  const content = kit.valueRow({ label: label("contentCheck"), displayValue: "", helpers: { element } });
  const next = button(label("next"), () => { pairIndex = (pairIndex + 1) % config.pairs.length; update(); });
  const controls = element("div", "lesson-controls");
  controls.appendChild(next);
  function update() {
    const pair = config.pairs[pairIndex];
    const sameSize = pair.firstSize === pair.secondSize;
    firstFile.set(ch22Text(T, pair.first));
    secondFile.set(ch22Text(T, pair.second));
    sizes.set(`${pair.firstSize} / ${pair.secondSize}`);
    filter.set(label(sameSize ? "sameSize" : "differentSize"));
    hashes.set(sameSize ? `${pair.firstHash} / ${pair.secondHash}` : label("notChecked"));
    content.set(label(!sameSize ? "notChecked" : pair.firstHash === pair.secondHash ? "sameContent" : "differentContent"));
  }
  panel.append(ch22Caption(element, label("caption")), firstFile.node, secondFile.node, sizes.node, filter.node, hashes.node, content.node, controls);
  update();
  return panel;
});
