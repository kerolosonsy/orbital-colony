"use strict";

function ch11Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch11Label(config, key, T) {
  return ch11Text(config.labels[key], T);
}

function ch11Caption(element, text) {
  const node = element("p", "ch11-sandbox-caption", text);
  node.setAttribute("role", "status");
  return node;
}

Lessons.registerSandbox("ch11-scout-channel", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = element("section", "lesson-sandbox ch11-scout-channel");
  const schedule = kit.valueRow({ label: ch11Label(config, "schedule", T), displayValue: "", helpers: { element, button } });
  const caption = ch11Caption(element, "");
  const channel = kit.valueRow({ label: ch11Label(config, "sender", T), displayValue: "", helpers: { element, button } });
  const queue = element("div", "ch11-sandbox-boxes");
  const received = kit.valueRow({ label: ch11Label(config, "receiver", T), displayValue: "", helpers: { element, button } });
  let scheduleIndex = 0;
  let eventIndex = 0;
  const next = button(ch11Label(config, "next", T), () => {
    const events = config.schedules[scheduleIndex].events;
    eventIndex = (eventIndex + 1) % events.length;
    update();
  });
  const again = button(ch11Label(config, "again", T), () => {
    scheduleIndex = (scheduleIndex + 1) % config.schedules.length;
    eventIndex = 0;
    update();
  });
  const controls = element("div", "ch11-sandbox-controls");
  controls.append(next, again);
  panel.append(schedule.node, caption, channel.node, queue, received.node, controls);

  function update() {
    const currentSchedule = config.schedules[scheduleIndex];
    const frame = currentSchedule.events[eventIndex];
    schedule.set(ch11Text(currentSchedule.name, T));
    caption.textContent = ch11Text(frame.caption, T);
    channel.set(ch11Label(config, frame.closed ? "closed" : "open", T));
    queue.replaceChildren(kit.boxes({
      label: ch11Label(config, "queue", T),
      values: frame.queue.length ? frame.queue : [ch11Label(config, "empty", T)],
      helpers: { element, button },
    }));
    received.set(frame.received || ch11Label(config, "none", T));
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch11-worker-pool", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = element("section", "lesson-sandbox ch11-worker-pool");
  const schedule = kit.valueRow({ label: ch11Label(config, "schedule", T), displayValue: "", helpers: { element, button } });
  const caption = ch11Caption(element, "");
  const jobs = element("div", "ch11-sandbox-boxes");
  const results = element("div", "ch11-sandbox-boxes");
  const waitGroup = kit.valueRow({ label: ch11Label(config, "waitGroup", T), displayValue: "", helpers: { element, button } });
  const resultsChannel = kit.valueRow({ label: ch11Label(config, "resultsChannel", T), displayValue: "", helpers: { element, button } });
  let scheduleIndex = 0;
  let eventIndex = 0;
  const next = button(ch11Label(config, "next", T), () => {
    const events = config.schedules[scheduleIndex].events;
    eventIndex = (eventIndex + 1) % events.length;
    update();
  });
  const again = button(ch11Label(config, "again", T), () => {
    scheduleIndex = (scheduleIndex + 1) % config.schedules.length;
    eventIndex = 0;
    update();
  });
  const controls = element("div", "ch11-sandbox-controls");
  controls.append(next, again);
  panel.append(schedule.node, caption, jobs, results, waitGroup.node, resultsChannel.node, controls);

  function update() {
    const currentSchedule = config.schedules[scheduleIndex];
    const frame = currentSchedule.events[eventIndex];
    schedule.set(ch11Text(currentSchedule.name, T));
    caption.textContent = ch11Text(frame.caption, T);
    jobs.replaceChildren(kit.boxes({
      label: ch11Label(config, "jobs", T), values: frame.jobs.length ? frame.jobs : [ch11Label(config, "empty", T)], helpers: { element, button },
    }));
    results.replaceChildren(kit.boxes({
      label: ch11Label(config, "results", T), values: frame.results.length ? frame.results : [ch11Label(config, "empty", T)], helpers: { element, button },
    }));
    waitGroup.set(ch11Label(config, frame.done ? "done" : "waiting", T));
    resultsChannel.set(ch11Label(config, frame.closed ? "closed" : "open", T));
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch11-storm-select", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = element("section", "lesson-sandbox ch11-storm-select");
  let warningReady = false;
  let calmClosed = false;
  let useDefault = true;
  let readyChoice = 0;
  let outcome = ch11Label(config, "idle", T);
  let note = "";
  const warning = kit.valueRow({ label: ch11Label(config, "warning", T), displayValue: "", helpers: { element, button } });
  const calm = kit.valueRow({ label: ch11Label(config, "calm", T), displayValue: "", helpers: { element, button } });
  const mode = kit.valueRow({ label: ch11Label(config, "mode", T), displayValue: "", helpers: { element, button } });
  const result = kit.valueRow({ label: ch11Label(config, "outcome", T), displayValue: "", helpers: { element, button } });
  const caption = ch11Caption(element, "");
  const warningButton = button("", () => {
    warningReady = !warningReady;
    resetResult();
  });
  const calmButton = button("", () => {
    calmClosed = !calmClosed;
    resetResult();
  });
  const modeButton = button("", () => {
    useDefault = !useDefault;
    resetResult();
  });
  const run = button(ch11Label(config, "run", T), () => {
    note = "";
    if (warningReady && calmClosed) {
      const choices = ["return", "resume"];
      outcome = ch11Label(config, choices[readyChoice % choices.length], T);
      readyChoice++;
      note = ch11Label(config, "bothReady", T);
    } else if (warningReady) {
      outcome = ch11Label(config, "return", T);
    } else if (calmClosed) {
      outcome = ch11Label(config, "resume", T);
    } else {
      outcome = ch11Label(config, useDefault ? "work" : "wait", T);
    }
    update();
  });
  const controls = element("div", "ch11-sandbox-controls ch11-storm-controls");
  controls.append(warningButton, calmButton, modeButton, run);
  panel.append(warning.node, calm.node, mode.node, result.node, caption, controls);

  function resetResult() {
    outcome = ch11Label(config, "idle", T);
    note = "";
    update();
  }

  function update() {
    warning.set(ch11Label(config, warningReady ? "warningReady" : "warningEmpty", T));
    calm.set(ch11Label(config, calmClosed ? "calmClosed" : "calmOpen", T));
    mode.set(ch11Label(config, useDefault ? "defaultMode" : "blockingMode", T));
    result.set(outcome);
    caption.textContent = note;
    warningButton.textContent = `${ch11Label(config, "warningButton", T)}: ${ch11Label(config, warningReady ? "warningReady" : "warningEmpty", T)}`;
    warningButton.setAttribute("aria-pressed", String(warningReady));
    calmButton.textContent = `${ch11Label(config, "calmButton", T)}: ${ch11Label(config, calmClosed ? "calmClosed" : "calmOpen", T)}`;
    calmButton.setAttribute("aria-pressed", String(calmClosed));
    modeButton.textContent = `${ch11Label(config, "modeButton", T)}: ${ch11Label(config, useDefault ? "defaultMode" : "blockingMode", T)}`;
    modeButton.setAttribute("aria-pressed", String(!useDefault));
  }

  update();
  return panel;
});
