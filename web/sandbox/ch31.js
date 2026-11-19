"use strict";

function ch31Text(pair, T) { return T(pair.ar, pair.en); }
function ch31Panel(element, kind) { return element("section", `lesson-sandbox ch31-sandbox ${kind}`); }
function ch31Label(config, key, T) { return ch31Text(config.labels[key], T); }
function ch31Row(element, label, value, className = "") {
  const row = element("p", `ch31-row ${className}`.trim());
  const data = element("bdi", "ch31-value", String(value));
  data.dir = "auto";
  row.append(element("span", "ch31-label", label), data);
  row.setAttribute("role", "status");
  return row;
}

Lessons.registerSandbox("ch31-stream-buffer", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch31Panel(element, "ch31-stream-buffer");
  const chunk = ch31Row(element, ch31Label(config, "chunk", T), "");
  const pending = ch31Row(element, ch31Label(config, "pending", T), "");
  const lines = element("ol", "ch31-lines");
  const status = element("p", "ch31-caption", ch31Label(config, "initial", T));
  status.setAttribute("role", "status");
  let index = 0;
  let buffered = "";
  const complete = [];
  const read = button(ch31Label(config, "read", T), () => {
    if (index >= config.chunks.length) return;
    const next = config.chunks[index++];
    chunk.children[1].textContent = next.text + (next.delimiter ? ` ${config.delimiterMark}` : "");
    buffered += next.text;
    if (next.delimiter) {
      complete.push(buffered);
      buffered = "";
    }
    pending.children[1].textContent = buffered || ch31Label(config, "empty", T);
    lines.replaceChildren(...complete.map((line) => element("li", "ch31-line", line)));
    status.textContent = ch31Text(next.caption, T);
    read.disabled = index === config.chunks.length;
  });
  const reset = button(ch31Label(config, "reset", T), () => {
    index = 0;
    buffered = "";
    complete.length = 0;
    chunk.children[1].textContent = ch31Label(config, "empty", T);
    pending.children[1].textContent = ch31Label(config, "empty", T);
    lines.replaceChildren();
    status.textContent = ch31Label(config, "initial", T);
    read.disabled = false;
  });
  chunk.children[1].textContent = ch31Label(config, "empty", T);
  pending.children[1].textContent = ch31Label(config, "empty", T);
  const controls = element("div", "ch31-controls");
  controls.append(read, reset);
  panel.append(element("p", "ch31-heading", ch31Label(config, "title", T)), chunk, pending,
    element("p", "ch31-label", ch31Label(config, "lines", T)), lines, status, controls);
  return panel;
});

Lessons.registerSandbox("ch31-mailbox-owner", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch31Panel(element, "ch31-mailbox-owner");
  const events = element("ol", "ch31-lines");
  const mailboxes = element("div", "ch31-mailboxes");
  const status = element("p", "ch31-caption", ch31Label(config, "caption", T));
  status.setAttribute("role", "status");
  let runIndex = 0;
  function render() {
    const queues = { slow: [...config.initialSlow], fast: [] };
    const eventRows = [];
    config.runs[runIndex].forEach((event) => {
      const queue = queues[event.recipient];
      const accepted = queue.length < config.capacity;
      if (accepted) queue.push(event.message);
      const who = ch31Label(config, event.recipient, T);
      const result = ch31Text(config.events[event.recipient], T);
      eventRows.push(element("li", "ch31-line", `${who}: ${event.message} · ${result}`));
    });
    events.replaceChildren(...eventRows);
    mailboxes.replaceChildren(
      ch31Row(element, ch31Label(config, "slow", T), queues.slow.join(", ") || ch31Label(config, "empty", T)),
      ch31Row(element, ch31Label(config, "fast", T), queues.fast.join(", ") || ch31Label(config, "empty", T)),
    );
  }
  const run = button(ch31Label(config, "run", T), () => {
    runIndex = (runIndex + 1) % config.runs.length;
    render();
  });
  const controls = element("div", "ch31-controls");
  controls.appendChild(run);
  panel.append(element("p", "ch31-heading", ch31Label(config, "title", T)),
    element("p", "ch31-label", ch31Label(config, "events", T)), events, mailboxes, status, controls);
  render();
  return panel;
});

Lessons.registerSandbox("ch31-udp-packets", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch31Panel(element, "ch31-udp-packets");
  const packet = ch31Row(element, ch31Label(config, "packet", T), config.packets[0]);
  const status = element("p", "ch31-caption", ch31Label(config, "initial", T));
  status.setAttribute("role", "status");
  let index = 0;
  const receive = button(ch31Label(config, "receive", T), () => {
    index = Math.min(index + 1, config.packets.length - 1);
    packet.children[1].textContent = config.packets[index];
    status.textContent = `${index + 1} / ${config.packets.length}`;
    receive.disabled = index === config.packets.length - 1;
  });
  const reset = button(ch31Label(config, "reset", T), () => {
    index = 0;
    packet.children[1].textContent = config.packets[0];
    status.textContent = ch31Label(config, "initial", T);
    receive.disabled = false;
  });
  const controls = element("div", "ch31-controls");
  controls.append(receive, reset);
  panel.append(element("p", "ch31-heading", ch31Label(config, "title", T)), packet, status, controls);
  return panel;
});

Lessons.registerSandbox("ch31-byte-order", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch31Panel(element, "ch31-byte-order");
  let order = 0;
  const bytes = ch31Row(element, ch31Label(config, "wire", T), config.bytes.join(" · "));
  const value = ch31Row(element, ch31Label(config, "value", T), config.orders[order].value);
  const remaining = ch31Row(element, ch31Label(config, "remaining", T), config.bytes.slice(2).join(" · ") || ch31Label(config, "empty", T));
  const swap = button(ch31Label(config, "swap", T), () => {
    order = (order + 1) % config.orders.length;
    value.children[1].textContent = `${config.orders[order].name}: ${config.orders[order].value}`;
  });
  value.children[1].textContent = `${config.orders[order].name}: ${config.orders[order].value}`;
  const controls = element("div", "ch31-controls");
  controls.appendChild(swap);
  panel.append(element("p", "ch31-heading", ch31Label(config, "title", T)), bytes, value, remaining, controls);
  return panel;
});

Lessons.registerSandbox("ch31-dns-outcomes", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch31Panel(element, "ch31-dns-outcomes");
  let index = 0;
  const outcome = config.outcomes[index];
  const name = ch31Row(element, ch31Label(config, "case", T), ch31Text(outcome.name, T));
  const result = ch31Row(element, ch31Label(config, "result", T), outcome.result);
  const status = element("p", "ch31-caption", ch31Text(outcome.caption, T));
  status.setAttribute("role", "status");
  function update() {
    const selected = config.outcomes[index];
    name.children[1].textContent = ch31Text(selected.name, T);
    result.children[1].textContent = selected.result;
    status.textContent = ch31Text(selected.caption, T);
  }
  const next = button(ch31Label(config, "next", T), () => { index = (index + 1) % config.outcomes.length; update(); });
  const reset = button(ch31Label(config, "reset", T), () => { index = 0; update(); });
  const controls = element("div", "ch31-controls");
  controls.append(next, reset);
  panel.append(element("p", "ch31-heading", ch31Label(config, "title", T)), name, result, status, controls);
  return panel;
});

Lessons.registerSandbox("ch31-result-slots", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch31Panel(element, "ch31-result-slots");
  const order = element("ol", "ch31-lines");
  const slots = element("ol", "ch31-lines");
  let runIndex = 0;
  function render() {
    const completion = config.orders[runIndex];
    order.replaceChildren(...completion.map((slot) => element("li", "ch31-line", config.ports[slot])));
    slots.replaceChildren(...config.ports.map((port, index) => element("li", "ch31-line", `${port}: ${config.results[index]}`)));
  }
  const run = button(ch31Label(config, "run", T), () => {
    runIndex = (runIndex + 1) % config.orders.length;
    render();
  });
  const controls = element("div", "ch31-controls");
  controls.appendChild(run);
  panel.append(element("p", "ch31-heading", ch31Label(config, "title", T)),
    element("p", "ch31-label", ch31Label(config, "completion", T)), order,
    element("p", "ch31-label", ch31Label(config, "slots", T)), slots,
    element("p", "ch31-caption", ch31Label(config, "caption", T)), controls);
  render();
  return panel;
});
