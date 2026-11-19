"use strict";

function ch29Text(value, T) {
  return T(value.ar, value.en);
}

function ch29Panel(element, kind) {
  return element("section", `lesson-sandbox ch29-playground ${kind}`);
}

function ch29Row(element, label, value) {
  const row = element("div", "ch29-row");
  row.append(element("span", "ch29-label", label), element("bdi", "ch29-value", String(value)));
  return row;
}

function ch29Cell(row) {
  return row.children[1];
}

function ch29Status(element, initial = "") {
  const status = element("p", "ch29-status", initial);
  status.setAttribute("role", "status");
  status.setAttribute("aria-live", "polite");
  return status;
}

Lessons.registerSandbox("ch29-sql-binding", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch29Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch29Panel(element, "ch29-sql-binding");
  const query = ch29Row(element, label("query"), config.query);
  const parameter = ch29Row(element, label("parameter"), "");
  const result = ch29Row(element, label("result"), "");
  const status = ch29Status(element);
  let index = 0;
  function update() {
    const row = config.rows[index];
    ch29Cell(parameter).textContent = text(row.key);
    ch29Cell(result).textContent = text(row.value);
    status.textContent = label(row.found ? "found" : "missing");
  }
  const controls = element("div", "ch29-controls");
  controls.appendChild(button(label("button"), () => { index = (index + 1) % config.rows.length; update(); }));
  panel.append(query, parameter, result, status, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch29-migration-boundary", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch29Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch29Panel(element, "ch29-migration-boundary");
  const saved = ch29Row(element, label("saved"), text(config.schema));
  const version = ch29Row(element, label("version"), config.version - 1);
  const proposal = ch29Row(element, label("proposal"), label("pending"));
  const status = ch29Status(element, label("pending"));
  let staged = false;
  let committed = false;
  function update() {
    ch29Cell(saved).textContent = text(committed ? config.nextSchema : config.schema);
    ch29Cell(version).textContent = String(committed ? config.version : config.version - 1);
    ch29Cell(proposal).textContent = text(staged ? config.nextSchema : config.labels.pending);
    status.textContent = label(committed ? "committed" : staged ? "pending" : "pending");
    commit.disabled = !staged || committed;
  }
  const controls = element("div", "ch29-controls");
  const stage = button(label("stage"), () => { staged = true; committed = false; update(); });
  const commit = button(label("commit"), () => { if (!staged) return; committed = true; staged = false; update(); });
  controls.append(stage, commit, button(label("reset"), () => { staged = false; committed = false; update(); }));
  panel.append(saved, version, proposal, status, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch29-atomic-transfer", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch29Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch29Panel(element, "ch29-atomic-transfer");
  const order = ch29Status(element);
  const source = ch29Row(element, label("source"), "");
  const north = ch29Row(element, label("north"), "");
  const south = ch29Row(element, label("south"), "");
  const log = element("ol", "ch29-transfer-log");
  const note = element("p", "ch29-note", label("vary"));
  let orderIndex = 0;
  function update() {
    let remaining = config.startingBalance;
    const destinations = { north: 0, south: 0 };
    const entries = [];
    for (const requestIndex of config.orders[orderIndex]) {
      const request = config.requests[requestIndex];
      const moved = remaining >= request.amount;
      if (moved) {
        remaining -= request.amount;
        destinations[request.target] += request.amount;
      }
      entries.push(`${text(request.name)} · ${request.amount} · ${label(moved ? "accepted" : "skipped")}`);
    }
    order.textContent = label("possibleOrder");
    ch29Cell(source).textContent = String(remaining);
    ch29Cell(north).textContent = String(destinations.north);
    ch29Cell(south).textContent = String(destinations.south);
    log.replaceChildren(...entries.map((entry) => element("li", "", entry)));
  }
  const controls = element("div", "ch29-controls");
  controls.appendChild(button(label("run"), () => { orderIndex = (orderIndex + 1) % config.orders.length; update(); }));
  panel.append(order, note, source, north, south, log, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch29-pool-state", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch29Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch29Panel(element, "ch29-pool-state");
  const maximum = ch29Row(element, label("max"), config.maxOpen);
  const inUse = ch29Row(element, label("inUse"), "");
  const idle = ch29Row(element, label("idle"), "");
  let active = config.inUse;
  let available = config.idle;
  function update() {
    ch29Cell(inUse).textContent = String(active);
    ch29Cell(idle).textContent = String(available);
    acquire.disabled = active === config.maxOpen;
    release.disabled = active === 0;
  }
  const controls = element("div", "ch29-controls");
  const acquire = button(label("acquire"), () => {
    if (available > 0) available--;
    active++;
    update();
  });
  const release = button(label("release"), () => {
    if (active > 0) { active--; available++; }
    update();
  });
  controls.append(acquire, release, button(label("reset"), () => { active = config.inUse; available = config.idle; update(); }));
  panel.append(maximum, inUse, idle, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch29-repository-boundary", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch29Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch29Panel(element, "ch29-repository-boundary");
  const contract = ch29Row(element, label("contract"), config.contract);
  const call = ch29Row(element, label("call"), config.call);
  const implementation = ch29Row(element, label("implementation"), "");
  const effect = ch29Status(element);
  let index = 0;
  function update() {
    ch29Cell(implementation).textContent = text(config.implementations[index].name);
    effect.textContent = text(config.implementations[index].effect);
  }
  const controls = element("div", "ch29-controls");
  controls.appendChild(button(label("switch"), () => { index = (index + 1) % config.implementations.length; update(); }));
  panel.append(contract, call, implementation, effect, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch29-association", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch29Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch29Panel(element, "ch29-association");
  const parent = ch29Row(element, label("parent"), text(config.parentName));
  const childrenLabel = element("p", "ch29-label", label("children"));
  const children = element("ol", "ch29-child-list");
  const status = ch29Status(element, label("unloaded"));
  let loaded = false;
  const controls = element("div", "ch29-controls");
  const toggle = button(label("load"), () => { loaded = !loaded; update(); });
  function update() {
    children.hidden = !loaded;
    children.replaceChildren(...(loaded ? config.children.map((child) => element("li", "", text(child))) : []));
    status.textContent = label(loaded ? "loaded" : "unloaded");
    toggle.textContent = label(loaded ? "hide" : "load");
  }
  controls.appendChild(toggle);
  panel.append(parent, childrenLabel, children, status, controls);
  update();
  return panel;
});
