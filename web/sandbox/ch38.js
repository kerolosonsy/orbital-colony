"use strict";

function ch38Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch38Label(config, key, T) {
  return ch38Text(config.labels[key], T);
}

function ch38Panel(kind, config, { element, T }) {
  const panel = element("section", `lesson-sandbox ch38-playground ch38-${kind}`);
  const title = element("h4", "ch38-title", ch38Label(config, "title", T));
  const state = element("div", "ch38-state");
  const status = element("p", "ch38-status");
  status.setAttribute("role", "status");
  status.setAttribute("aria-live", "polite");
  const controls = element("div", "ch38-controls");
  panel.append(title, state, status, controls);
  return { node: panel, state, status, controls };
}

function ch38ValueRow(config, key, displayValue, helpers) {
  return helpers.kit.valueRow({
    label: ch38Label(config, key, helpers.T),
    displayValue,
    helpers: { element: helpers.element, button: helpers.button },
  });
}

function ch38Control(panel, label, action, button) {
  const control = button(label, action);
  control.classList.add("ch38-action");
  panel.controls.appendChild(control);
  return control;
}

function ch38UpdateIterator(iterator, status) {
  iterator.count.set(iterator.index);
  iterator.view.status.textContent = ch38Label(iterator.config, status, iterator.T);
}

function ch38NextReading(iterator) {
  if (iterator.stopped) return ch38UpdateIterator(iterator, "stopped");
  if (iterator.index >= iterator.config.values.length) return ch38UpdateIterator(iterator, "ended");
  iterator.readings.push(iterator.config.values[iterator.index]);
  iterator.consumed.set(iterator.readings);
  iterator.index += 1;
  ch38UpdateIterator(iterator, "read");
}

function ch38StopReading(iterator) {
  iterator.stopped = true;
  ch38UpdateIterator(iterator, "stopped");
}

function ch38ResetIterator(iterator) {
  iterator.index = 0;
  iterator.stopped = false;
  iterator.readings = [];
  iterator.consumed.set([]);
  ch38UpdateIterator(iterator, "cleared");
}

function ch38RegisterIterator(block, helpers) {
  const { element, T, kit } = helpers;
  const config = block.sandbox.config;
  const view = ch38Panel("iterator", config, helpers);
  const iterator = { index: 0, stopped: false, readings: [], config, view, T };
  const source = kit.boxes({ label: ch38Label(config, "values", T), values: config.values, helpers: { element, button: helpers.button } });
  iterator.consumed = kit.logList({ label: ch38Label(config, "consumed", T), entries: [], helpers: { element, button: helpers.button } });
  iterator.count = ch38ValueRow(config, "consumed", 0, helpers);
  view.state.append(source, iterator.count.node, iterator.consumed.node);
  view.status.textContent = ch38Label(config, "ready", T);
  ch38Control(view, ch38Label(config, "next", T), () => ch38NextReading(iterator), helpers.button);
  ch38Control(view, ch38Label(config, "stop", T), () => ch38StopReading(iterator), helpers.button);
  ch38Control(view, ch38Label(config, "reset", T), () => ch38ResetIterator(iterator), helpers.button);
  return view.node;
}

function ch38CreateSetState(block, helpers) {
  const { element, T, kit, button } = helpers;
  const config = block.sandbox.config;
  const view = ch38Panel("set-labels", config, helpers);
  const setState = { config, view, T, kit, element, button, rowIndex: 0, seen: new Set() };
  setState.row = ch38ValueRow(config, "row", 1, helpers);
  setState.membership = ch38ValueRow(config, "selected", ch38Label(config, "empty", T), helpers);
  setState.slots = element("div", "ch38-slots");
  view.state.append(setState.row.node, setState.slots, setState.membership.node);
  ch38RenderSet(setState, "ready");
  return setState;
}

function ch38RenderSet(setState, status) {
  const { config, view, T, kit, element, button } = setState;
  const shownRow = setState.rowIndex === 0 ? 0 : Math.min(setState.rowIndex - 1, config.rows.length - 1);
  setState.row.set(Math.max(1, setState.rowIndex));
  setState.slots.replaceChildren(kit.boxes({ label: ch38Label(config, "slots", T), values: config.rows[shownRow], helpers: { element, button } }));
  setState.membership.set([...setState.seen].join(", ") || ch38Label(config, "empty", T));
  view.status.textContent = ch38Label(config, status, T);
}

function ch38InspectNextRow(setState) {
  if (setState.rowIndex >= setState.config.rows.length) return ch38RenderSet(setState, "done");
  const slotsInRow = setState.config.rows[setState.rowIndex];
  let added = 0;
  let duplicate = false;
  let skipped = false;
  for (const slot of slotsInRow) {
    if (slot === 0) {
      skipped = true;
      break;
    }
    if (setState.seen.has(slot)) {
      duplicate = true;
      continue;
    }
    setState.seen.add(slot);
    added += 1;
  }
  setState.rowIndex += 1;
  const outcome = skipped ? "skipped" : added ? "picked" : duplicate ? "duplicate" : "skipped";
  ch38RenderSet(setState, outcome);
}

function ch38ResetSet(setState) {
  setState.rowIndex = 0;
  setState.seen = new Set();
  ch38RenderSet(setState, "cleared");
}

function ch38RegisterSet(block, helpers) {
  const setState = ch38CreateSetState(block, helpers);
  const { config, view, T } = setState;
  ch38Control(view, ch38Label(config, "next", T), () => ch38InspectNextRow(setState), helpers.button);
  ch38Control(view, ch38Label(config, "reset", T), () => ch38ResetSet(setState), helpers.button);
  return view.node;
}

function ch38RegisterOnce(block, helpers) {
  const { T } = helpers;
  const config = block.sandbox.config;
  const view = ch38Panel("once", config, helpers);
  let runs = 0;
  const count = ch38ValueRow(config, "runs", runs, helpers);
  const cachedValue = ch38ValueRow(config, "value", config.labels.zero, helpers);
  view.state.append(count.node, cachedValue.node);
  view.status.textContent = ch38Label(config, "ready", T);
  ch38Control(view, ch38Label(config, "call", T), () => {
    const status = runs === 0 ? "loaded" : "reused";
    runs = 1;
    count.set(runs);
    view.status.textContent = ch38Label(config, status, T);
  }, helpers.button);
  ch38Control(view, ch38Label(config, "reset", T), () => {
    runs = 0;
    count.set(runs);
    view.status.textContent = ch38Label(config, "cleared", T);
  }, helpers.button);
  return view.node;
}

function ch38RegisterErrgroup(block, helpers) {
  const { element, T, kit } = helpers;
  const config = block.sandbox.config;
  const view = ch38Panel("errgroup-orders", config, helpers);
  let orderIndex = 0;
  const order = kit.logList({ label: ch38Label(config, "jobs", T), entries: config.orders[orderIndex], helpers: { element, button: helpers.button } });
  view.state.appendChild(order.node);
  view.status.textContent = ch38Label(config, "possible", T);
  ch38Control(view, ch38Label(config, "run", T), () => {
    orderIndex = (orderIndex + 1) % config.orders.length;
    order.set(config.orders[orderIndex]);
    view.status.textContent = ch38Label(config, "another", T);
  }, helpers.button);
  ch38Control(view, ch38Label(config, "reset", T), () => {
    orderIndex = 0;
    order.set(config.orders[orderIndex]);
    view.status.textContent = ch38Label(config, "restored", T);
  }, helpers.button);
  return view.node;
}

function ch38RegisterBuilder(block, helpers) {
  const { T } = helpers;
  const config = block.sandbox.config;
  const view = ch38Panel("builder", config, helpers);
  let index = 0;
  let buffer = "";
  const content = ch38ValueRow(config, "buffer", ch38Label(config, "empty", T), helpers);
  const length = ch38ValueRow(config, "length", 0, helpers);
  view.state.append(content.node, length.node);
  view.status.textContent = ch38Label(config, "ready", T);
  ch38Control(view, ch38Label(config, "write", T), () => {
    if (index < config.parts.length) {
      buffer += config.parts[index];
      index += 1;
    }
    content.set(buffer || ch38Label(config, "empty", T));
    length.set(buffer.length);
    view.status.textContent = ch38Label(config, "written", T);
  }, helpers.button);
  ch38Control(view, ch38Label(config, "reset", T), () => {
    index = 0;
    buffer = "";
    content.set(ch38Label(config, "empty", T));
    length.set(0);
    view.status.textContent = ch38Label(config, "cleared", T);
  }, helpers.button);
  return view.node;
}

function ch38RegisterVet(block, helpers) {
  const { element, T } = helpers;
  const config = block.sandbox.config;
  const view = ch38Panel("vet-lenses", config, helpers);
  let testsDone = false;
  let vetDone = false;
  const tests = ch38ValueRow(config, "tests", ch38Label(config, "waiting", T), helpers);
  const vet = ch38ValueRow(config, "vet", ch38Label(config, "waiting", T), helpers);
  view.state.append(tests.node, vet.node);
  function update(status) {
    tests.set(testsDone ? ch38Label(config, "complete", T) : ch38Label(config, "waiting", T));
    vet.set(vetDone ? ch38Label(config, "complete", T) : ch38Label(config, "waiting", T));
    view.status.textContent = ch38Label(config, status, T);
  }
  view.status.textContent = ch38Label(config, "ready", T);
  ch38Control(view, ch38Label(config, "runTests", T), () => {
    testsDone = true;
    update(vetDone ? "bothDone" : "testsDone");
  }, helpers.button);
  ch38Control(view, ch38Label(config, "runVet", T), () => {
    vetDone = true;
    update(testsDone ? "bothDone" : "vetDone");
  }, helpers.button);
  ch38Control(view, ch38Label(config, "reset", T), () => {
    testsDone = false;
    vetDone = false;
    update("cleared");
  }, helpers.button);
  return view.node;
}

function ch38RegisterProfile(block, helpers) {
  const { element, T } = helpers;
  const config = block.sandbox.config;
  const view = ch38Panel("profile-focus", config, helpers);
  const focus = ch38ValueRow(config, "focus", ch38Label(config, "choose", T), helpers);
  const description = element("p", "ch38-profile-description");
  view.state.append(focus.node, description);
  view.status.textContent = ch38Label(config, "choose", T);
  ch38Control(view, ch38Label(config, "cpu", T), () => {
    focus.set(ch38Label(config, "cpu", T));
    description.textContent = ch38Label(config, "cpuFocus", T);
    view.status.textContent = ch38Label(config, "cpuFocus", T);
  }, helpers.button);
  ch38Control(view, ch38Label(config, "heap", T), () => {
    focus.set(ch38Label(config, "heap", T));
    description.textContent = ch38Label(config, "heapFocus", T);
    view.status.textContent = ch38Label(config, "heapFocus", T);
  }, helpers.button);
  return view.node;
}

Lessons.registerSandbox("ch38-iterator-stop", ch38RegisterIterator);
Lessons.registerSandbox("ch38-set-labels", ch38RegisterSet);
Lessons.registerSandbox("ch38-once-cache", ch38RegisterOnce);
Lessons.registerSandbox("ch38-errgroup-orders", ch38RegisterErrgroup);
Lessons.registerSandbox("ch38-builder-buffer", ch38RegisterBuilder);
Lessons.registerSandbox("ch38-vet-lenses", ch38RegisterVet);
Lessons.registerSandbox("ch38-profile-focus", ch38RegisterProfile);
