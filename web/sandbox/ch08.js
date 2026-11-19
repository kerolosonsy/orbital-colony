"use strict";

function ch08Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch08Panel(element, kind) {
  return element("section", "lesson-sandbox ch08-sandbox ch08-" + kind);
}

function ch08Caption(element, text) {
  const caption = element("p", "ch08-caption", text);
  caption.setAttribute("role", "status");
  caption.setAttribute("aria-live", "polite");
  return caption;
}

Lessons.registerSandbox("ch08-copy-pointer", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const text = (pair) => ch08Text(pair, T);
  const panel = ch08Panel(element, "copy-pointer");
  let caller = config.start;
  let local = text(config.labels.noLocal);
  const callerValue = kit.valueRow({ label: text(config.labels.caller), displayValue: caller, helpers: { element } });
  const localValue = kit.valueRow({ label: text(config.labels.local), displayValue: local, helpers: { element } });
  const caption = ch08Caption(element, text(config.captions.ready));
  const controls = element("div", "ch08-controls");
  const update = () => {
    callerValue.set(String(caller));
    localValue.set(local);
  };
  controls.append(
    button(text(config.labels.copy), () => {
      caller = config.start;
      local = String(config.changed);
      update();
      caption.textContent = text(config.captions.copy);
    }),
    button(text(config.labels.pointer), () => {
      caller = config.changed;
      local = text(config.labels.noLocal);
      update();
      caption.textContent = text(config.captions.pointer);
    }),
    button(text(config.labels.reset), () => {
      caller = config.start;
      local = text(config.labels.noLocal);
      update();
      caption.textContent = text(config.captions.reset);
    }),
  );
  panel.append(callerValue.node, localValue.node, caption, controls);
  return panel;
});

Lessons.registerSandbox("ch08-receivers", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const text = (pair) => ch08Text(pair, T);
  const panel = ch08Panel(element, "receivers");
  let caller = config.start;
  let receiver = text(config.labels.notCalled);
  const callerValue = kit.valueRow({ label: text(config.labels.caller), displayValue: caller, helpers: { element } });
  const receiverValue = kit.valueRow({ label: text(config.labels.receiver), displayValue: receiver, helpers: { element } });
  const caption = ch08Caption(element, text(config.captions.ready));
  const controls = element("div", "ch08-controls");
  const update = () => {
    callerValue.set(String(caller));
    receiverValue.set(receiver);
  };
  controls.append(
    button(text(config.labels.value), () => {
      caller = config.start;
      receiver = String(config.start + config.amount);
      update();
      caption.textContent = text(config.captions.value);
    }),
    button(text(config.labels.pointer), () => {
      caller = config.start + config.amount;
      receiver = String(caller);
      update();
      caption.textContent = text(config.captions.pointer);
    }),
    button(text(config.labels.reset), () => {
      caller = config.start;
      receiver = text(config.labels.notCalled);
      update();
      caption.textContent = text(config.captions.reset);
    }),
  );
  panel.append(callerValue.node, receiverValue.node, caption, controls);
  return panel;
});

Lessons.registerSandbox("ch08-nil-guard", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const text = (pair) => ch08Text(pair, T);
  const panel = ch08Panel(element, "nil-guard");
  const pointer = kit.valueRow({ label: text(config.labels.pointer), displayValue: text(config.labels.nilPointer), helpers: { element } });
  const result = kit.valueRow({ label: text(config.labels.result), displayValue: text(config.labels.emptyResult), helpers: { element } });
  const caption = ch08Caption(element, text(config.captions.empty));
  const controls = element("div", "ch08-controls");
  const toggle = kit.toggle({
    label: text(config.labels.target),
    onLabel: text(config.labels.on),
    offLabel: text(config.labels.off),
    isOn: false,
    onChange(isPresent) {
      pointer.set(isPresent ? text(config.labels.valuePointer) : text(config.labels.nilPointer));
      result.set(isPresent ? String(config.value) : text(config.labels.emptyResult));
      caption.textContent = text(isPresent ? config.captions.value : config.captions.empty);
    },
    helpers: { element, button },
  });
  controls.appendChild(toggle.node);
  panel.append(pointer.node, result.node, caption, controls);
  return panel;
});

Lessons.registerSandbox("ch08-cart-chain", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const text = (pair) => ch08Text(pair, T);
  const panel = ch08Panel(element, "cart-chain");
  const heading = element("p", "ch08-heading", text(config.labels.chain));
  const chain = element("div", "ch08-chain");
  config.nodes.forEach((node, index) => {
    if (index > 0) chain.appendChild(element("span", "ch08-link", text(config.link)));
    const tag = element("bdi", "ch08-node", text(node.tag));
    chain.appendChild(tag);
  });
  const pointer = kit.valueRow({ label: text(config.labels.pointer), displayValue: text(config.nodes[0].tag), helpers: { element } });
  const caption = ch08Caption(element, text(config.captions.ready));
  const controls = element("div", "ch08-controls");
  let index = 0;
  const update = () => {
    pointer.set(text(config.nodes[index].tag));
    chain.querySelectorAll(".ch08-node").forEach((node, nodeIndex) => node.classList.toggle("current", nodeIndex === index));
  };
  const follow = button(text(config.labels.follow), () => {
    if (index < config.nodes.length - 1) {
      index += 1;
      caption.textContent = text(config.captions.moved);
      update();
    } else {
      caption.textContent = text(config.captions.end);
    }
  });
  const reset = button(text(config.labels.reset), () => {
    index = 0;
    caption.textContent = text(config.captions.reset);
    update();
  });
  controls.append(follow, reset);
  panel.append(heading, chain, pointer.node, caption, controls);
  update();
  return panel;
});
