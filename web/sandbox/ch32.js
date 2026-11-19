"use strict";

function ch32Text(pair, T) { return T(pair.ar, pair.en); }
function ch32Panel(element, kind) { return element("section", `lesson-sandbox ch32-sandbox ${kind}`); }
function ch32Label(config, key, T) { return ch32Text(config.labels[key], T); }
function ch32Row(element, label, value, className = "") {
  const row = element("p", `ch32-row ${className}`.trim());
  const data = element("bdi", "ch32-value", String(value));
  data.dir = "auto";
  row.append(element("span", "ch32-label", label), data);
  row.setAttribute("role", "status");
  return row;
}

Lessons.registerSandbox("ch32-message-pumps", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch32Panel(element, "ch32-message-pumps");
  const events = element("ol", "ch32-events");
  let runIndex = 0;
  function render() {
    events.replaceChildren(...config.runs[runIndex].map((event) => {
      const action = ch32Text(config.eventLabels[event.kind], T);
      return element("li", "ch32-event", `${action}: ${event.type} · ${event.payload}`);
    }));
  }
  const run = button(ch32Label(config, "run", T), () => {
    runIndex = (runIndex + 1) % config.runs.length;
    render();
  });
  const controls = element("div", "ch32-controls");
  controls.appendChild(run);
  const sample = config.runs[0];
  panel.append(element("p", "ch32-heading", ch32Label(config, "title", T)),
    ch32Row(element, ch32Label(config, "type", T), sample.map((event) => event.type).join(" · ")),
    ch32Row(element, ch32Label(config, "payload", T), sample.map((event) => event.payload).join(" · ")),
    element("p", "ch32-label", ch32Label(config, "events", T)), events,
    element("p", "ch32-caption", ch32Label(config, "caption", T)), controls);
  render();
  return panel;
});

Lessons.registerSandbox("ch32-unary-rpc", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch32Panel(element, "ch32-unary-rpc");
  let index = -1;
  const query = ch32Row(element, ch32Label(config, "query", T), ch32Label(config, "empty", T));
  const reply = ch32Row(element, ch32Label(config, "reply", T), ch32Label(config, "empty", T));
  const status = element("p", "ch32-caption", ch32Label(config, "initial", T));
  status.setAttribute("role", "status");
  function render() {
    if (index < 0) {
      query.children[1].textContent = ch32Label(config, "empty", T);
      reply.children[1].textContent = ch32Label(config, "empty", T);
      status.textContent = ch32Label(config, "initial", T);
      return;
    }
    const item = config.queries[index];
    query.children[1].textContent = `${item.sector}`;
    reply.children[1].textContent = `${item.sector} · ${item.pressure}`;
    status.textContent = `${ch32Label(config, "sector", T)}: ${item.sector} · ${ch32Label(config, "pressure", T)}: ${item.pressure}`;
  }
  const send = button(ch32Label(config, "send", T), () => {
    index = (index + 1) % config.queries.length;
    render();
  });
  const reset = button(ch32Label(config, "reset", T), () => { index = -1; render(); });
  const controls = element("div", "ch32-controls");
  controls.append(send, reset);
  panel.append(element("p", "ch32-heading", ch32Label(config, "title", T)), query, reply, status, controls);
  return panel;
});

Lessons.registerSandbox("ch32-stream-interleave", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch32Panel(element, "ch32-stream-interleave");
  const events = element("ol", "ch32-events");
  let runIndex = 0;
  function render() {
    events.replaceChildren(...config.runs[runIndex].map((event) =>
      element("li", "ch32-event", `${ch32Label(config, "stream", T)} ${event.stream} · ${ch32Label(config, "sequence", T)} ${event.sequence} · ${ch32Label(config, "delta", T)} ${event.delta}`)));
  }
  const run = button(ch32Label(config, "run", T), () => {
    runIndex = (runIndex + 1) % config.runs.length;
    render();
  });
  const controls = element("div", "ch32-controls");
  controls.appendChild(run);
  panel.append(element("p", "ch32-heading", ch32Label(config, "title", T)), events,
    element("p", "ch32-caption", ch32Label(config, "caption", T)), controls);
  render();
  return panel;
});

Lessons.registerSandbox("ch32-peer-trust", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch32Panel(element, "ch32-peer-trust");
  let index = 0;
  const peer = ch32Row(element, ch32Label(config, "peer", T), "");
  const expects = ch32Row(element, ch32Label(config, "checks", T), "");
  const presents = ch32Row(element, ch32Label(config, "peerPresents", T), "");
  function render() {
    const side = config.sides[index];
    peer.children[1].textContent = ch32Text(side.peer, T);
    expects.children[1].textContent = side.expects;
    presents.children[1].textContent = side.presents;
  }
  const flip = button(ch32Label(config, "flip", T), () => { index = (index + 1) % config.sides.length; render(); });
  const controls = element("div", "ch32-controls");
  controls.appendChild(flip);
  panel.append(element("p", "ch32-heading", ch32Label(config, "title", T)), peer, expects, presents, controls);
  render();
  return panel;
});

Lessons.registerSandbox("ch32-alpn-reads", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch32Panel(element, "ch32-alpn-reads");
  let handshakeIndex = 0;
  let readIndex = 0;
  const offered = ch32Row(element, ch32Label(config, "offered", T), config.offered.join(" · "));
  const selected = ch32Row(element, ch32Label(config, "selected", T), config.negotiated[handshakeIndex]);
  const reads = element("ol", "ch32-events");
  function renderReads() {
    reads.replaceChildren(...config.readOrders[readIndex].map((key) =>
      element("li", "ch32-event", ch32Text(config.readLabels[key], T))));
  }
  const handshake = button(ch32Label(config, "handshake", T), () => {
    handshakeIndex = (handshakeIndex + 1) % config.negotiated.length;
    selected.children[1].textContent = config.negotiated[handshakeIndex];
  });
  const run = button(ch32Label(config, "run", T), () => {
    readIndex = (readIndex + 1) % config.readOrders.length;
    renderReads();
  });
  const controls = element("div", "ch32-controls");
  controls.append(handshake, run);
  panel.append(element("p", "ch32-heading", ch32Label(config, "title", T)), offered, selected,
    ch32Row(element, ch32Label(config, "connection", T), config.connectionValue),
    element("p", "ch32-label", ch32Label(config, "reads", T)), reads, controls);
  renderReads();
  return panel;
});

Lessons.registerSandbox("ch32-ssh-allowlist", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const panel = ch32Panel(element, "ch32-ssh-allowlist");
  let index = 0;
  const choice = ch32Row(element, ch32Label(config, "choice", T), "");
  const action = ch32Row(element, ch32Label(config, "action", T), "");
  function render() {
    const command = config.commands[index];
    choice.children[1].textContent = command.choice;
    action.children[1].textContent = ch32Text(command.action, T);
  }
  const next = button(ch32Label(config, "next", T), () => { index = (index + 1) % config.commands.length; render(); });
  const controls = element("div", "ch32-controls");
  controls.appendChild(next);
  panel.append(element("p", "ch32-heading", ch32Label(config, "title", T)), choice, action, controls);
  render();
  return panel;
});
