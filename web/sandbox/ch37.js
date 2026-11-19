"use strict";

function ch37Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch37Label(config, key, T) {
  return ch37Text(config.labels[key], T);
}

function ch37Panel(kind, config, { element, T }) {
  const panel = element("section", `lesson-sandbox ch37-playground ch37-${kind}`);
  const title = element("h4", "ch37-title", ch37Label(config, "title", T));
  const state = element("div", "ch37-state");
  const status = element("p", "ch37-status");
  status.setAttribute("role", "status");
  status.setAttribute("aria-live", "polite");
  const controls = element("div", "ch37-controls");
  panel.append(title, state, status, controls);
  return { node: panel, state, status, controls };
}

function ch37ValueRow(config, key, displayValue, helpers) {
  return helpers.kit.valueRow({
    label: ch37Label(config, key, helpers.T),
    displayValue,
    helpers: { element: helpers.element, button: helpers.button },
  });
}

function ch37Control(panel, label, action, button) {
  const control = button(label, action);
  control.classList.add("ch37-action");
  panel.controls.appendChild(control);
  return control;
}

function ch37UpdateCursor(cursorState, status) {
  cursorState.pageNumber.set(cursorState.page + 1);
  cursorState.cursor.set(cursorState.config.cursors[cursorState.page]);
  cursorState.items.set(cursorState.config.pages[cursorState.page]);
  cursorState.view.status.textContent = ch37Label(cursorState.config, status, cursorState.T);
}

function ch37RequestFrame(frame, view, helpers) {
  const { element, T } = helpers;
  view.status.textContent = ch37Text(frame.caption, T);
  const frameNode = element("div", "ch37-frame");
  frameNode.append(
    element("p", "ch37-stage", ch37Text(frame.stage, T)),
    element("bdi", "ch37-value", frame.detail),
  );
  return frameNode;
}

function ch37RegisterAtomic(block, helpers) {
  const { element, T } = helpers;
  const config = block.sandbox.config;
  const view = ch37Panel("atomic-ledger", config, helpers);
  const schema = ch37ValueRow(config, "schema", config.values.schema, helpers);
  const ledger = ch37ValueRow(config, "ledger", config.values.ledger, helpers);
  view.state.append(schema.node, ledger.node);
  function update(mode) {
    schema.set(mode === "committed" ? config.values.next : config.values.schema);
    ledger.set(mode === "committed" ? config.values.next : config.values.ledger);
    view.status.textContent = ch37Label(config, mode, T);
  }
  view.status.textContent = ch37Label(config, "ready", T);
  ch37Control(view, ch37Label(config, "fail", T), () => update("failed"), helpers.button);
  ch37Control(view, ch37Label(config, "commit", T), () => update("committed"), helpers.button);
  ch37Control(view, ch37Label(config, "reset", T), () => update("ready"), helpers.button);
  return view.node;
}

function ch37RegisterSession(block, helpers) {
  const { element, T } = helpers;
  const config = block.sandbox.config;
  const view = ch37Panel("session", config, helpers);
  let active = false;
  const client = ch37ValueRow(config, "client", ch37Label(config, "empty", T), helpers);
  const store = ch37ValueRow(config, "store", ch37Label(config, "empty", T), helpers);
  view.state.append(client.node, store.node);
  function update(status, nextActive = active) {
    active = nextActive;
    client.set(active ? config.values.token : ch37Label(config, "empty", T));
    store.set(active ? config.values.digest : ch37Label(config, "empty", T));
    view.status.textContent = ch37Label(config, status, T);
  }
  view.status.textContent = ch37Label(config, "ready", T);
  ch37Control(view, ch37Label(config, "issue", T), () => update("issued", true), helpers.button);
  ch37Control(view, ch37Label(config, "logout", T), () => update("ended", false), helpers.button);
  ch37Control(view, ch37Label(config, "reset", T), () => update("cleared", false), helpers.button);
  return view.node;
}

function ch37RegisterCrew(block, helpers) {
  const { element, T } = helpers;
  const config = block.sandbox.config;
  const view = ch37Panel("crew-scope", config, helpers);
  const lookup = ch37ValueRow(config, "key", ch37Label(config, "choose", T), helpers);
  const role = ch37ValueRow(config, "role", ch37Label(config, "empty", T), helpers);
  view.state.append(lookup.node, role.node);
  function choose(index) {
    const membership = config.lookups[index];
    const displayRole = ch37Text(membership.role, T);
    lookup.set(ch37Text(membership.key, T));
    role.set(displayRole);
    view.status.textContent = displayRole;
  }
  view.status.textContent = ch37Label(config, "choose", T);
  ["blue", "gold", "red"].forEach((key, index) => {
    ch37Control(view, ch37Label(config, key, T), () => choose(index), helpers.button);
  });
  return view.node;
}

function ch37RegisterCursor(block, helpers) {
  const { element, T, kit } = helpers;
  const config = block.sandbox.config;
  const view = ch37Panel("cursor-pages", config, helpers);
  const pageNumber = ch37ValueRow(config, "page", 1, helpers);
  const cursor = ch37ValueRow(config, "cursor", config.cursors[0], helpers);
  const items = kit.logList({ label: ch37Label(config, "page", T), entries: config.pages[0], helpers: { element, button: helpers.button } });
  view.state.append(pageNumber.node, cursor.node, items.node);
  const cursorState = { page: 0, pageNumber, cursor, items, view, config, T };
  view.status.textContent = ch37Label(config, "ready", T);
  ch37Control(view, ch37Label(config, "next", T), () => {
    cursorState.page = Math.min(cursorState.page + 1, config.pages.length - 1);
    ch37UpdateCursor(cursorState, "advanced");
  }, helpers.button);
  ch37Control(view, ch37Label(config, "reset", T), () => {
    cursorState.page = 0;
    ch37UpdateCursor(cursorState, "resetStatus");
  }, helpers.button);
  return view.node;
}

function ch37RegisterRequest(block, helpers) {
  const { element, T, kit } = helpers;
  const config = block.sandbox.config;
  const view = ch37Panel("request-path", config, helpers);
  const stepper = kit.stepper({
    frames: config.frames,
    labels: {
      previous: ch37Label(config, "previous", T),
      next: ch37Label(config, "next", T),
      reset: ch37Label(config, "reset", T),
    },
    helpers: { element, button: helpers.button },
    renderFrame: (frame) => ch37RequestFrame(frame, view, helpers),
  });
  view.state.appendChild(stepper.node);
  return view.node;
}

Lessons.registerSandbox("ch37-atomic-ledger", ch37RegisterAtomic);
Lessons.registerSandbox("ch37-session-digest", ch37RegisterSession);
Lessons.registerSandbox("ch37-scoped-membership", ch37RegisterCrew);
Lessons.registerSandbox("ch37-cursor-pages", ch37RegisterCursor);
Lessons.registerSandbox("ch37-request-path", ch37RegisterRequest);
