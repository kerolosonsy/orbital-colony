"use strict";

function ch26Panel(className, element) {
  return element("section", `lesson-sandbox ${className}`);
}

function ch26Text(config, T, key) {
  const pair = config.labels[key];
  return T(pair.ar, pair.en);
}

function ch26Display(T, pair) {
  return T(pair.ar, pair.en);
}

Lessons.registerSandbox("ch26-handler-flow", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch26Panel("ch26-handler-flow", element);
  const labels = { previous: ch26Text(config, T, "previous"), next: ch26Text(config, T, "next"), reset: ch26Text(config, T, "reset") };
  const stepper = kit.stepper({ frames: config.frames, renderFrame(frame) {
    const view = element("div", "ch26-frame-state");
    view.appendChild(element("p", "ch26-caption", ch26Display(T, frame.caption)));
    for (const key of ["header", "status", "body"]) {
      const row = element("div", "ch26-row");
      row.append(element("span", "ch26-label", ch26Text(config, T, key)), element("bdi", "ch26-value", String(frame[key])));
      view.appendChild(row);
    }
    return view;
  }, labels, helpers: { element, button } });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch26-rest-atomic", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch26Panel("ch26-rest-atomic", element);
  let id = config.nextID;
  let orderIndex = 0;
  let position = 0;
  let records = [];
  const nextID = kit.valueRow({ label: ch26Text(config, T, "nextID"), displayValue: id, helpers: { element } });
  const note = element("p", "ch26-note", ch26Text(config, T, "order"));
  const recordsLog = kit.logList({ label: ch26Text(config, T, "records"), entries: [], helpers: { element } });
  function update() {
    nextID.set(id);
    recordsLog.set(records.map((record) => `${record.id} · ${record.name}`));
  }
  const create = button(ch26Text(config, T, "next"), () => {
    const order = config.orders[orderIndex];
    if (position >= order.length) return;
    const record = config.records[order[position++]];
    id += 1;
    records.push({ id, name: ch26Display(T, record.name) });
    update();
  });
  const again = button(ch26Text(config, T, "again"), () => {
    orderIndex = (orderIndex + 1) % config.orders.length;
    position = 0;
    id = config.nextID;
    records = [];
    update();
  });
  const reset = button(ch26Text(config, T, "reset"), () => {
    orderIndex = 0;
    position = 0;
    id = config.nextID;
    records = [];
    update();
  });
  panel.append(note, nextID.node, recordsLog.node, create, again, reset);
  update();
  return panel;
});

Lessons.registerSandbox("ch26-middleware-order", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch26Panel("ch26-middleware-order", element);
  const stepper = kit.stepper({ frames: config.frames, renderFrame(frame) {
    const view = element("div", "ch26-frame-state");
    view.appendChild(element("p", "ch26-caption", ch26Display(T, frame.caption)));
    const trace = element("ol", "ch26-trace");
    trace.replaceChildren(...frame.trace.map((entry) => element("li", "ch26-trace-item", ch26Display(T, entry))));
    view.appendChild(trace);
    return view;
  }, labels: { previous: ch26Text(config, T, "previous"), next: ch26Text(config, T, "next"), reset: ch26Text(config, T, "reset") }, helpers: { element, button } });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch26-error-envelope", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch26Panel("ch26-error-envelope", element);
  const stepper = kit.stepper({ frames: config.frames, renderFrame(frame) {
    const view = element("div", "ch26-frame-state");
    view.appendChild(element("p", "ch26-caption", ch26Display(T, frame.caption)));
    view.appendChild(element("p", "ch26-label", ch26Text(config, T, "chain")));
    const chain = element("ol", "ch26-chain");
    chain.replaceChildren(...frame.chain.map((errorPart) => element("li", "ch26-chain-item", ch26Display(T, errorPart))));
    const outcomeRow = element("div", "ch26-row");
    outcomeRow.append(element("span", "ch26-label", ch26Text(config, T, "result")), element("span", "ch26-value", ch26Display(T, frame.result)));
    view.append(chain, outcomeRow);
    return view;
  }, labels: { previous: ch26Text(config, T, "previous"), next: ch26Text(config, T, "next"), reset: ch26Text(config, T, "reset") }, helpers: { element, button } });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch26-static-route", (block, { element, button, T }) => {
  const { config } = block.sandbox;
  const panel = ch26Panel("ch26-static-route", element);
  const selected = element("bdi", "ch26-value");
  const routeStatus = element("p", "ch26-state", ch26Text(config, T, "choose"));
  function choose(path) {
    selected.textContent = path;
    routeStatus.textContent = path === config.allowedPath ? ch26Text(config, T, "allowed") : ch26Text(config, T, "denied");
  }
  const exact = button(config.allowedPath, () => choose(config.allowedPath));
  const slash = button(config.nearPath, () => choose(config.nearPath));
  const selectedRow = element("div", "ch26-row");
  selectedRow.append(element("span", "ch26-label", ch26Text(config, T, "route")), selected);
  panel.append(selectedRow, routeStatus, exact, slash);
  return panel;
});

Lessons.registerSandbox("ch26-upload-stream", (block, { element, button, T }) => {
  const { config } = block.sandbox;
  const panel = ch26Panel("ch26-upload-stream", element);
  const source = element("bdi", "ch26-value", config.payload);
  const staging = element("bdi", "ch26-value");
  const copied = element("bdi", "ch26-value");
  const eof = element("p", "ch26-state", "");
  let count = config.count;
  let ended = false;
  function update() {
    staging.textContent = config.payload.slice(0, count) || ch26Text(config, T, "empty");
    copied.textContent = `${count}/${config.payload.length}`;
    eof.textContent = ended ? ch26Text(config, T, "eof") : "";
  }
  const read = button(ch26Text(config, T, "next"), () => {
    if (count < config.payload.length) count += 1;
    else ended = true;
    update();
  });
  const reset = button(ch26Text(config, T, "reset"), () => { count = config.count; ended = false; update(); });
  const sourceRow = element("div", "ch26-row");
  sourceRow.append(element("span", "ch26-label", ch26Text(config, T, "source")), source);
  const copiedRow = element("div", "ch26-row");
  copiedRow.append(element("span", "ch26-label", ch26Text(config, T, "copied")), copied);
  const stagingRow = element("div", "ch26-row");
  stagingRow.append(element("span", "ch26-label", ch26Text(config, T, "staging")), staging);
  panel.append(sourceRow, copiedRow, stagingRow, eof, read, reset);
  update();
  return panel;
});
