"use strict";

function ch21Text(T, pair) {
  return T(pair.ar, pair.en);
}

function ch21Label(T, config, key) {
  return ch21Text(T, config.labels[key]);
}

function ch21Panel(element, kind) {
  return element("section", `lesson-sandbox ${kind}`);
}

function ch21Caption(element, text) {
  const caption = element("p", "ch21-caption", text);
  caption.setAttribute("role", "status");
  caption.dir = "auto";
  return caption;
}

Lessons.registerSandbox("ch21-flags-env", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch21Label(T, config, key);
  const panel = ch21Panel(element, "ch21-flags-env");
  let enabled = false;
  let stateIndex = 0;
  const lookup = kit.valueRow({ label: label("lookup"), displayValue: label("environment"), helpers: { element } });
  const flag = kit.valueRow({ label: label("flag"), displayValue: label("off"), helpers: { element } });
  const environmentName = kit.valueRow({ label: label("name"), displayValue: "", helpers: { element } });
  const presence = kit.valueRow({ label: label("presence"), displayValue: label("missing"), helpers: { element } });
  const environmentValue = kit.valueRow({ label: label("value"), displayValue: "", helpers: { element } });
  const toggle = button(label("toggle"), () => { enabled = !enabled; update(); });
  const next = button(label("next"), () => { stateIndex = (stateIndex + 1) % config.states.length; update(); });
  const controls = element("div", "lesson-controls");
  controls.append(toggle, next);
  function update() {
    const state = config.states[stateIndex];
    flag.set(label(enabled ? "on" : "off"));
    environmentName.set(ch21Text(T, state.label));
    presence.set(label(state.present ? "present" : "missing"));
    environmentValue.set(ch21Text(T, state.value));
    toggle.setAttribute("aria-pressed", String(enabled));
  }
  panel.append(ch21Caption(element, label("caption")), lookup.node, flag.node, environmentName.node, presence.node, environmentValue.node, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch21-reader-lines", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch21Label(T, config, key);
  const panel = ch21Panel(element, "ch21-reader-lines");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch21-reader-frame");
      view.append(
        ch21Caption(element, ch21Text(T, frame.caption)),
        kit.valueRow({ label: label("current"), displayValue: ch21Text(T, frame.current), helpers: { element } }).node,
        kit.logList({ label: label("written"), entries: frame.written.map((line) => ch21Text(T, line)), helpers: { element } }).node,
      );
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch21-injected-draw", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch21Label(T, config, key);
  const panel = ch21Panel(element, "ch21-injected-draw");
  let varied = false;
  let calls = 0;
  let selected = null;
  const source = kit.toggle({
    label: label("source"), onLabel: label("varied"), offLabel: label("fixed"),
    onChange(next) { varied = next; calls = 0; selected = null; update(); }, helpers: { element, button },
  });
  const bound = kit.valueRow({ label: label("bound"), displayValue: String(config.choices.length), helpers: { element } });
  const index = kit.valueRow({ label: label("index"), displayValue: label("none"), helpers: { element } });
  const result = kit.valueRow({ label: label("selected"), displayValue: label("none"), helpers: { element } });
  const count = kit.valueRow({ label: label("calls"), displayValue: "0", helpers: { element } });
  const draw = button(label("draw"), () => {
    const results = varied ? config.variedResults : config.fixedResults;
    selected = results[calls % results.length];
    calls++;
    update();
  });
  const controls = element("div", "lesson-controls");
  controls.append(source.node, draw);
  function update() {
    index.set(selected === null ? label("none") : String(selected));
    result.set(selected === null ? label("none") : ch21Text(T, config.choices[selected]));
    count.set(String(calls));
  }
  panel.append(ch21Caption(element, label("caption")), bound.node, index.node, result.node, count.node, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch21-maneuver-values", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch21Label(T, config, key);
  const panel = ch21Panel(element, "ch21-maneuver-values");
  let commandIndex = 0;
  const command = kit.valueRow({ label: label("command"), displayValue: "", helpers: { element } });
  const named = kit.valueRow({ label: label("named"), displayValue: "", helpers: { element } });
  const ordinal = kit.valueRow({ label: label("ordinal"), displayValue: "", helpers: { element } });
  const next = button(label("next"), () => { commandIndex = (commandIndex + 1) % config.commands.length; update(); });
  const controls = element("div", "lesson-controls");
  controls.appendChild(next);
  function update() {
    const example = config.commands[commandIndex];
    command.set(ch21Text(T, example.text));
    named.set(ch21Text(T, example.name));
    ordinal.set(String(example.ordinal));
  }
  panel.append(ch21Caption(element, label("caption")), command.node, named.node, ordinal.node, controls);
  update();
  return panel;
});
