"use strict";

function ch12Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch12Label(config, key, T) {
  return ch12Text(config.labels[key], T);
}

function ch12Caption(element, text) {
  const node = element("p", "ch12-sandbox-caption", text);
  node.setAttribute("role", "status");
  return node;
}

Lessons.registerSandbox("ch12-shared-counter", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = element("section", "lesson-sandbox ch12-shared-counter");
  let locked = false;
  let runIndex = 0;
  let eventIndex = 0;
  const schedule = kit.valueRow({ label: ch12Label(config, "scenario", T), displayValue: "", helpers: { element, button } });
  const total = kit.valueRow({ label: ch12Label(config, "total", T), displayValue: "", helpers: { element, button } });
  const waitGroup = kit.valueRow({ label: ch12Label(config, "waitGroup", T), displayValue: "", helpers: { element, button } });
  const caption = ch12Caption(element, "");
  const lockButton = button("", () => {
    locked = !locked;
    eventIndex = 0;
    update();
  });
  const next = button(ch12Label(config, "next", T), () => {
    eventIndex = (eventIndex + 1) % events().length;
    update();
  });
  const again = button(ch12Label(config, "again", T), () => {
    runIndex = (runIndex + 1) % config.runs.length;
    eventIndex = 0;
    update();
  });
  const controls = element("div", "ch12-sandbox-controls");
  controls.append(lockButton, next, again);
  panel.append(schedule.node, caption, total.node, waitGroup.node, controls);

  function events() {
    const run = config.runs[runIndex];
    return locked ? run.locked : run.unlocked;
  }

  function update() {
    const run = config.runs[runIndex];
    const frame = events()[eventIndex];
    schedule.set(ch12Text(run.name, T));
    caption.textContent = ch12Text(frame.caption, T);
    total.set(frame.total);
    waitGroup.set(ch12Label(config, frame.done ? "done" : "waiting", T));
    lockButton.textContent = `${ch12Label(config, "mutex", T)}: ${ch12Label(config, locked ? "on" : "off", T)}`;
    lockButton.setAttribute("aria-pressed", String(locked));
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch12-registry-lock", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = element("section", "lesson-sandbox ch12-registry-lock");
  let orderIndex = 0;
  let eventIndex = 0;
  const schedule = kit.valueRow({ label: ch12Label(config, "scenario", T), displayValue: "", helpers: { element, button } });
  const sharedLock = kit.valueRow({ label: ch12Label(config, "sharedLock", T), displayValue: ch12Label(config, "shared", T), helpers: { element, button } });
  const caption = ch12Caption(element, "");
  const map = element("div", "ch12-sandbox-boxes");
  const log = element("div", "ch12-sandbox-boxes");
  const waitGroup = kit.valueRow({ label: ch12Label(config, "waitGroup", T), displayValue: "", helpers: { element, button } });
  const next = button(ch12Label(config, "next", T), () => {
    eventIndex = (eventIndex + 1) % (config.writers.length + 1);
    update();
  });
  const again = button(ch12Label(config, "again", T), () => {
    orderIndex = (orderIndex + 1) % config.orders.length;
    eventIndex = 0;
    update();
  });
  const controls = element("div", "ch12-sandbox-controls");
  controls.append(next, again);
  panel.append(schedule.node, sharedLock.node, caption, map, log, waitGroup.node, controls);

  function update() {
    const order = config.orders[orderIndex];
    const recorded = order.sequence.slice(0, eventIndex).map((writerIndex) => config.writers[writerIndex]);
    const latest = recorded[recorded.length - 1];
    schedule.set(ch12Text(order.name, T));
    if (eventIndex === 0) caption.textContent = ch12Label(config, "start", T);
    else if (eventIndex === config.writers.length) caption.textContent = ch12Label(config, "completion", T);
    else caption.textContent = ch12Text(latest.caption, T);
    const sortedMap = [...recorded].sort((a, b) => a.key.localeCompare(b.key));
    map.replaceChildren(kit.boxes({
      label: ch12Label(config, "map", T),
      values: sortedMap.length ? sortedMap.map((writer) => `${writer.key}=${writer.value}`) : [ch12Label(config, "empty", T)],
      helpers: { element, button },
    }));
    log.replaceChildren(kit.boxes({
      label: ch12Label(config, "log", T),
      values: recorded.length ? recorded.map((writer) => `${ch12Text(writer.worker, T)}: ${writer.key}`) : [ch12Label(config, "empty", T)],
      helpers: { element, button },
    }));
    waitGroup.set(ch12Label(config, eventIndex === config.writers.length ? "done" : "waiting", T));
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch12-deadlock-order", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = element("section", "lesson-sandbox ch12-deadlock-order");
  let fixedOrder = false;
  let scheduleIndex = 0;
  let eventIndex = 0;
  const schedule = kit.valueRow({ label: ch12Label(config, "schedule", T), displayValue: "", helpers: { element, button } });
  const order = kit.valueRow({ label: ch12Label(config, "lockOrder", T), displayValue: "", helpers: { element, button } });
  const caption = ch12Caption(element, "");
  const drill = kit.valueRow({ label: ch12Label(config, "drill", T), displayValue: "", helpers: { element, button } });
  const cart = kit.valueRow({ label: ch12Label(config, "cart", T), displayValue: "", helpers: { element, button } });
  const state = kit.valueRow({ label: ch12Label(config, "status", T), displayValue: "", helpers: { element, button } });
  const next = button(ch12Label(config, "next", T), () => {
    eventIndex = (eventIndex + 1) % events().length;
    update();
  });
  const again = button(ch12Label(config, "again", T), () => {
    scheduleIndex = (scheduleIndex + 1) % schedules().length;
    eventIndex = 0;
    update();
  });
  const toggle = button("", () => {
    fixedOrder = !fixedOrder;
    scheduleIndex = 0;
    eventIndex = 0;
    update();
  });
  const controls = element("div", "ch12-sandbox-controls");
  controls.append(next, again, toggle);
  panel.append(schedule.node, order.node, caption, drill.node, cart.node, state.node, controls);

  function schedules() {
    return fixedOrder ? config.fixed : config.mixed;
  }

  function events() {
    return schedules()[scheduleIndex].events;
  }

  function holder(key) {
    if (key === "a") return ch12Label(config, "workerA", T);
    if (key === "b") return ch12Label(config, "workerB", T);
    return ch12Label(config, "free", T);
  }

  function update() {
    const scenario = schedules()[scheduleIndex];
    const frame = events()[eventIndex];
    schedule.set(ch12Text(scenario.name, T));
    order.set(ch12Label(config, fixedOrder ? "fixed" : "mixed", T));
    caption.textContent = ch12Text(frame.caption, T);
    drill.set(holder(frame.drill));
    cart.set(holder(frame.cart));
    state.set(ch12Label(config, frame.status, T));
    toggle.textContent = `${ch12Label(config, "orderButton", T)}: ${ch12Label(config, fixedOrder ? "fixed" : "mixed", T)}`;
    toggle.setAttribute("aria-pressed", String(fixedOrder));
  }

  update();
  return panel;
});
