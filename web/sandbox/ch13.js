"use strict";

Lessons.registerSandbox("ch13-recall", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const label = (key) => T(config.labels[key].ar, config.labels[key].en);
  const panel = element("section", "lesson-sandbox ch13-sandbox ch13-recall");
  const signal = element("p", "ch13-readout");
  const order = element("p", "ch13-order");
  const workers = element("ul", "ch13-workers");
  const joined = element("p", "ch13-readout");
  const controls = element("div", "ch13-controls");
  let orderIndex = 0;
  let canceled = false;
  let isJoined = false;
  let states = config.workers.map(() => "working");
  let schedule = [...config.orders[orderIndex]];

  const recall = button(label("recall"), () => { canceled = true; update(); });
  const advance = button(label("advance"), () => {
    const next = schedule.find((index) => states[index] !== "atBase");
    if (next !== undefined) states[next] = "atBase";
    update();
  });
  const shuffle = button(label("shuffleReset"), () => {
    orderIndex = (orderIndex + 1) % config.orders.length;
    schedule = [...config.orders[orderIndex]];
    canceled = false;
    isJoined = false;
    states = config.workers.map(() => "working");
    update();
  });
  const wait = button(label("wait"), () => { isJoined = true; update(); });
  controls.append(recall, advance, shuffle, wait);
  panel.append(signal, order, workers, joined, element("p", "ch13-note", label("explanation")), controls);

  function update() {
    signal.textContent = `${label("signal")}: ${label(canceled ? "canceled" : "active")}`;
    order.textContent = `${label("order")}: ${schedule.map((index) => config.workers[index]).join(" → ")}`;
    workers.replaceChildren(...config.workers.map((name, index) => {
      const row = element("li", "ch13-worker");
      const code = element("bdi", "ch13-code", name);
      code.setAttribute("dir", "ltr");
      const status = element("span", "", label(states[index] === "atBase" ? "atBase" : "working"));
      row.append(code, status);
      return row;
    }));
    const allReturned = states.every((state) => state === "atBase");
    joined.textContent = `${label("workersLabel")}: ${label(isJoined ? "joined" : "notJoined")}`;
    recall.disabled = canceled;
    advance.disabled = !canceled || allReturned;
    wait.disabled = !canceled || !allReturned || isJoined;
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch13-timeout", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const label = (key) => T(config.labels[key].ar, config.labels[key].en);
  const panel = element("section", "lesson-sandbox ch13-sandbox ch13-timeout");
  const state = element("p", "ch13-readout");
  const time = element("p", "ch13-readout ch13-code");
  const done = element("p", "ch13-readout ch13-code");
  const error = element("p", "ch13-readout ch13-code");
  const controls = element("div", "ch13-controls");
  let elapsed = 0;
  let reason = "active";

  const tick = button(label("tick"), () => {
    if (reason !== "active") return;
    elapsed++;
    if (elapsed >= config.deadlineTicks) reason = "deadlineReached";
    update();
  });
  const cancel = button(label("cancel"), () => {
    if (reason === "active") reason = "manuallyStopped";
    update();
  });
  const reset = button(label("reset"), () => { elapsed = 0; reason = "active"; update(); });
  controls.append(tick, cancel, reset);
  panel.append(state, time, done, error, controls);

  function update() {
    state.textContent = `${label("context")}: ${label(reason)}`;
    time.textContent = `${label("elapsed")}: ${elapsed} / ${config.deadlineTicks} ${label("deadline")}`;
    done.textContent = `${label("done")}: ${label(reason === "active" ? "open" : "closed")}`;
    const err = reason === "active" ? "nil" : reason === "deadlineReached" ? "deadlineError" : "canceledError";
    error.textContent = `${label("err")}: ${label(err)}`;
    tick.disabled = reason !== "active";
    cancel.disabled = reason !== "active";
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch13-first-wins", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const label = (key) => T(config.labels[key].ar, config.labels[key].en);
  const panel = element("section", "lesson-sandbox ch13-sandbox ch13-race");
  const order = element("p", "ch13-order");
  const signal = element("p", "ch13-readout");
  const workersLabel = element("p", "ch13-section-label", label("workersLabel"));
  const workers = element("ul", "ch13-workers");
  const controls = element("div", "ch13-controls");
  let orderIndex = 0;
  let schedule = [...config.orders[orderIndex]];
  let winner = null;
  let states = config.workers.map(() => "searching");

  const raceStep = button(label("raceStep"), () => {
    winner = schedule[0];
    states[winner] = "found";
    update();
  });
  const respond = button(label("respond"), () => {
    const next = schedule.find((index) => index !== winner && states[index] === "searching");
    if (next !== undefined) states[next] = "stopped";
    update();
  });
  const again = button(label("shuffleAgain"), () => {
    orderIndex = (orderIndex + 1) % config.orders.length;
    schedule = [...config.orders[orderIndex]];
    winner = null;
    states = config.workers.map(() => "searching");
    update();
  });
  controls.append(raceStep, respond, again);
  panel.append(order, signal, workersLabel, workers, element("p", "ch13-note", label("explanation")), controls);

  function update() {
    order.textContent = `${label("order")}: ${schedule.map((index) => config.workers[index]).join(" → ")}`;
    signal.textContent = `${label("signal")}: ${label(winner === null ? "open" : "canceled")}`;
    workers.replaceChildren(...config.workers.map((name, index) => {
      const row = element("li", "ch13-worker");
      const code = element("bdi", "ch13-code", name);
      code.setAttribute("dir", "ltr");
      const statusKey = states[index] === "found" ? "found" : states[index] === "stopped" ? "stopped" : "searching";
      row.append(code, element("span", "", label(statusKey)));
      return row;
    }));
    raceStep.disabled = winner !== null;
    respond.disabled = winner === null || states.every((value, index) => index === winner || value === "stopped");
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch13-context-chain", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => T(config.labels[key].ar, config.labels[key].en);
  const panel = element("section", "lesson-sandbox ch13-sandbox ch13-chain");
  const caption = element("p", "ch13-note");
  const context = element("p", "ch13-readout ch13-code");
  const error = element("p", "ch13-readout ch13-code");
  const counter = element("span", "ch13-counter");
  const controls = element("div", "ch13-controls");
  let index = 0;
  let deadline = false;
  let reasonControl;

  const frame = element("div", "ch13-chain-frame");
  const previous = button(label("previous"), () => { index--; update(); });
  const next = button(label("next"), () => { index++; update(); });
  const reset = button(label("reset"), () => { index = 0; update(); });
  reasonControl = kit.toggle({
    label: label("reason"),
    onLabel: label("deadline"),
    offLabel: label("canceled"),
    isOn: false,
    onChange(isDeadline) { deadline = isDeadline; update(); },
    helpers: { element, button },
  });
  controls.append(previous, counter, next, reset, reasonControl.node);
  panel.append(frame, caption, context, error, controls);

  function update() {
    const step = config.frames[index];
    const code = kit.codeView({ lines: config.code, activeLine: step.line, helpers: { element, button } });
    frame.replaceChildren(code.node);
    caption.textContent = T(step.caption.ar, step.caption.en);
    context.textContent = `${label("context")}: ${label("same")}`;
    const returned = index >= 2;
    error.textContent = `${label("err")}: ${returned ? label(deadline ? "deadline" : "canceled") : label("nil")}`;
    counter.textContent = `${index + 1} / ${config.frames.length}`;
    previous.disabled = index === 0;
    next.disabled = index === config.frames.length - 1;
  }

  update();
  return panel;
});
