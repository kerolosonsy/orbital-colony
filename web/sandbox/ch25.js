"use strict";

function ch25Panel(className, element) {
  return element("section", `lesson-sandbox ${className}`);
}

function ch25Text(config, T, key) {
  const pair = config.labels[key];
  return T(pair.ar, pair.en);
}

function ch25Display(T, pair) {
  return T(pair.ar, pair.en);
}

function ch25Line(element, className, text) {
  return element("p", className, text);
}

function ch25AdvanceRace(config, T, simulation, refresh) {
  const actions = config.schedules[simulation.scheduleIndex];
  if (simulation.actionIndex >= actions.length) return;
  const [workerID, action] = actions[simulation.actionIndex++].split("-");
  const worker = config.workers.find((candidate) => candidate.id === workerID);
  const workerName = ch25Display(T, worker.name);
  if (action === "read") {
    simulation.reads[workerID] = simulation.counter;
    simulation.entries.push(`${workerName} ${ch25Text(config, T, "read")} ${simulation.counter}`);
  } else {
    simulation.counter = simulation.reads[workerID] + worker.add;
    simulation.entries.push(`${workerName} ${ch25Text(config, T, "write")} ${simulation.counter}`);
  }
  refresh();
}

function ch25ResetRace(config, simulation, scheduleIndex, refresh) {
  simulation.scheduleIndex = scheduleIndex;
  simulation.actionIndex = 0;
  simulation.counter = config.initial;
  simulation.reads = {};
  simulation.entries = [];
  refresh();
}

Lessons.registerSandbox("ch25-locked-account", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch25Panel("ch25-locked-account", element);
  const balance = kit.valueRow({ label: ch25Text(config, T, "balance"), displayValue: config.balance, helpers: { element } });
  const state = ch25Line(element, "ch25-state", "");
  const protectedNote = ch25Line(element, "ch25-note", ch25Text(config, T, "protected"));
  const orderNote = ch25Line(element, "ch25-note", ch25Text(config, T, "order"));
  const log = kit.logList({ label: ch25Text(config, T, "history"), entries: [], helpers: { element } });
  let amountLeft = config.balance;
  let order = [0, 1];
  let position = 0;
  const transactions = config.transactions;
  function update() {
    balance.set(amountLeft);
    const done = position === order.length;
    state.textContent = done ? ch25Text(config, T, "done") : ch25Display(T, transactions[order[position]].worker);
    log.set(logEntries);
  }
  const logEntries = [];
  const next = button(ch25Text(config, T, "next"), () => {
    if (position >= order.length) return;
    const transaction = transactions[order[position++]];
    const worker = ch25Display(T, transaction.worker);
    if (amountLeft >= transaction.amount) {
      amountLeft -= transaction.amount;
      logEntries.push(`${worker} ${ch25Text(config, T, "spends")} ${transaction.amount}; ${ch25Text(config, T, "left")} ${amountLeft}`);
    } else {
      logEntries.push(`${worker} ${ch25Text(config, T, "skips")} (${ch25Text(config, T, "because")})`);
    }
    update();
  });
  const again = button(ch25Text(config, T, "again"), () => {
    order = order[0] === 0 ? [1, 0] : [0, 1];
    amountLeft = config.balance;
    position = 0;
    logEntries.length = 0;
    update();
  });
  const reset = button(ch25Text(config, T, "reset"), () => {
    order = [0, 1];
    amountLeft = config.balance;
    position = 0;
    logEntries.length = 0;
    update();
  });
  panel.append(balance.node, state, protectedNote, orderNote, log.node, next, again, reset);
  update();
  return panel;
});

Lessons.registerSandbox("ch25-race-report", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch25Panel("ch25-race-report", element);
  const counterRow = kit.valueRow({ label: ch25Text(config, T, "value"), displayValue: config.initial, helpers: { element } });
  const note = ch25Line(element, "ch25-note", ch25Text(config, T, "order"));
  const log = kit.logList({ label: ch25Text(config, T, "trace"), entries: [], helpers: { element } });
  const simulation = { scheduleIndex: 0, actionIndex: 0, counter: config.initial, reads: {}, entries: [] };
  function update() {
    counterRow.set(simulation.counter);
    log.set(simulation.entries);
  }
  const next = button(ch25Text(config, T, "next"), () => ch25AdvanceRace(config, T, simulation, update));
  const again = button(ch25Text(config, T, "again"), () => {
    const otherSchedule = (simulation.scheduleIndex + 1) % config.schedules.length;
    ch25ResetRace(config, simulation, otherSchedule, update);
  });
  const reset = button(ch25Text(config, T, "reset"), () => ch25ResetRace(config, simulation, 0, update));
  panel.append(counterRow.node, note, log.node, next, again, reset);
  update();
  return panel;
});

Lessons.registerSandbox("ch25-token-clock", (block, { element, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch25Panel("ch25-token-clock", element);
  const completeIntervals = Math.floor(config.elapsed / config.interval);
  const completedRow = kit.valueRow({ label: ch25Text(config, T, "complete"), displayValue: completeIntervals, helpers: { element } });
  const permitRow = kit.valueRow({ label: ch25Text(config, T, "tokens"), displayValue: `${Math.min(config.capacity, config.tokens + completeIntervals)} / ${config.capacity}`, helpers: { element } });
  const elapsedSlider = kit.slider({
    label: ch25Text(config, T, "elapsed"), currentValue: config.elapsed, min: config.min, max: config.max,
    onChange(elapsedSeconds) {
      const intervals = Math.floor(elapsedSeconds / config.interval);
      completedRow.set(intervals);
      permitRow.set(`${Math.min(config.capacity, config.tokens + intervals)} / ${config.capacity}`);
    },
    helpers: { element },
  });
  panel.append(elapsedSlider.node, completedRow.node, permitRow.node, ch25Line(element, "ch25-note", ch25Text(config, T, "order")));
  return panel;
});

Lessons.registerSandbox("ch25-nonblocking-bus", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch25Panel("ch25-nonblocking-bus", element);
  const queueView = element("div", "ch25-box-row");
  const state = ch25Line(element, "ch25-state", "");
  const note = ch25Line(element, "ch25-note", ch25Text(config, T, "order"));
  const log = kit.logList({ label: ch25Text(config, T, "attempts"), entries: [], helpers: { element } });
  let queue = [...config.initial];
  let logEntries = [];
  function update() {
    const displayedQueue = queue.length ? queue : [ch25Text(config, T, "empty")];
    queueView.replaceChildren(kit.boxes({ label: ch25Text(config, T, "queue"), values: displayedQueue, helpers: { element } }));
    state.textContent = `${ch25Text(config, T, "queue")}: ${queue.length}/${config.capacity}`;
    log.set(logEntries);
  }
  const publish = button(ch25Text(config, T, "publish"), () => {
    if (queue.length < config.capacity) {
      queue.push(config.message);
      logEntries.push(`${config.message}: ${ch25Text(config, T, "sent")}`);
    } else {
      logEntries.push(`${config.message}: ${ch25Text(config, T, "full")}`);
    }
    update();
  });
  const receive = button(ch25Text(config, T, "receive"), () => {
    if (queue.length) queue.shift();
    update();
  });
  const reset = button(ch25Text(config, T, "reset"), () => {
    queue = [...config.initial];
    logEntries = [];
    update();
  });
  const heading = element("p", "ch25-label", ch25Text(config, T, "queue"));
  panel.append(heading, queueView, state, note, log.node, publish, receive, reset);
  update();
  return panel;
});

Lessons.registerSandbox("ch25-ttl-cache", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch25Panel("ch25-ttl-cache", element);
  const age = kit.valueRow({ label: ch25Text(config, T, "age"), displayValue: config.age, helpers: { element } });
  const ttl = kit.valueRow({ label: ch25Text(config, T, "ttl"), displayValue: config.ttl, helpers: { element } });
  const cacheStatus = ch25Line(element, "ch25-state", "");
  let seconds = config.age;
  let read = false;
  function update() {
    age.set(seconds);
    ttl.set(config.ttl);
    cacheStatus.textContent = read ? ch25Text(config, T, seconds >= config.ttl ? "expired" : "fresh") : "";
  }
  const advance = button(ch25Text(config, T, "advance"), () => {
    seconds = Math.min(config.maxAge, seconds + 1);
    read = false;
    update();
  });
  const lookup = button(ch25Text(config, T, "read"), () => {
    read = true;
    update();
  });
  const reset = button(ch25Text(config, T, "reset"), () => {
    seconds = config.age;
    read = false;
    update();
  });
  panel.append(age.node, ttl.node, cacheStatus, advance, lookup, reset);
  update();
  return panel;
});

Lessons.registerSandbox("ch25-atomic-counter", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch25Panel("ch25-atomic-counter", element);
  const counterRow = kit.valueRow({ label: ch25Text(config, T, "load"), displayValue: 0, helpers: { element } });
  const orderNote = ch25Line(element, "ch25-note", ch25Text(config, T, "order"));
  const log = kit.logList({ label: ch25Text(config, T, "trace"), entries: [], helpers: { element } });
  let orderIndex = 0;
  let position = 0;
  let total = 0;
  let entries = [];
  function update() {
    counterRow.set(total);
    log.set(entries);
  }
  const next = button(ch25Text(config, T, "next"), () => {
    const order = config.orders[orderIndex];
    if (position >= order.length) return;
    const worker = config.workers[order[position++]];
    total += worker.add;
    entries.push(`${ch25Display(T, worker.name)} +${worker.add} = ${total}`);
    update();
  });
  const shuffle = button(ch25Text(config, T, "again"), () => {
    orderIndex = (orderIndex + 1) % config.orders.length;
    position = 0;
    total = 0;
    entries = [];
    update();
  });
  const reset = button(ch25Text(config, T, "reset"), () => {
    orderIndex = 0;
    position = 0;
    total = 0;
    entries = [];
    update();
  });
  panel.append(counterRow.node, orderNote, log.node, next, shuffle, reset);
  update();
  return panel;
});

Lessons.registerSandbox("ch25-closed-drain", (block, { element, button, T, kit }) => {
  const { config } = block.sandbox;
  const panel = ch25Panel("ch25-closed-drain", element);
  const channelState = ch25Line(element, "ch25-state", "");
  const values = element("div", "ch25-box-row");
  const queueStatus = ch25Line(element, "ch25-note", "");
  let queue = [...config.values];
  let closed = false;
  function update() {
    channelState.textContent = `${ch25Text(config, T, "channel")}: ${ch25Text(config, T, closed ? "closed" : "open")}`;
    values.replaceChildren(kit.boxes({ label: ch25Text(config, T, "channel"), values: queue, helpers: { element } }));
    queueStatus.textContent = queue.length ? ch25Text(config, T, "queued") : closed ? ch25Text(config, T, "drained") : ch25Text(config, T, "empty");
  }
  const close = button(ch25Text(config, T, "close"), () => { closed = true; update(); });
  const receive = button(ch25Text(config, T, "receive"), () => { if (queue.length) queue.shift(); update(); });
  const reset = button(ch25Text(config, T, "reset"), () => { queue = [...config.values]; closed = false; update(); });
  panel.append(channelState, values, queueStatus, close, receive, reset);
  update();
  return panel;
});
