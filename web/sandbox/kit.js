"use strict";

function updateSandboxStepper({ frames, renderFrame, index, previous, next, counter, frame }) {
  previous.disabled = index === 0;
  next.disabled = index === frames.length - 1;
  counter.textContent = `${index + 1} / ${frames.length}`;
  frame.replaceChildren(renderFrame(frames[index], index));
}

function sandboxKitStepper({ frames, renderFrame, labels, helpers }) {
  const root = helpers.element("div", "sandbox-kit sandbox-kit-stepper");
  const frame = helpers.element("div", "sandbox-kit-frame");
  const counter = helpers.element("span", "sandbox-kit-counter");
  counter.setAttribute("role", "status");
  counter.setAttribute("aria-live", "polite");
  let index = 0;
  function setIndex(nextIndex) {
    index = Math.max(0, Math.min(frames.length - 1, nextIndex));
    updateSandboxStepper({ frames, renderFrame, index, previous, next, counter, frame });
  }
  const previous = helpers.button(labels.previous, () => setIndex(index - 1));
  const next = helpers.button(labels.next, () => setIndex(index + 1));
  const reset = helpers.button(labels.reset, () => setIndex(0));
  const controls = helpers.element("div", "sandbox-kit-controls");
  controls.append(previous, counter, next, reset);
  root.append(frame, controls);
  setIndex(index);
  return { node: root, previous, next, reset, get index() { return index; } };
}

function sandboxKitToggle({ label, onLabel, offLabel, isOn = false, onChange = () => {}, helpers }) {
  const node = helpers.button("", () => {
    isOn = !isOn;
    update();
    onChange(isOn);
  });
  function update() {
    node.textContent = `${label}: ${isOn ? onLabel : offLabel}`;
    node.setAttribute("aria-pressed", String(isOn));
  }
  update();
  return { node, get isOn() { return isOn; }, set(nextIsOn) { isOn = Boolean(nextIsOn); update(); } };
}

function sandboxKitSlider({ label, currentValue, min, max, step = 1, onChange = () => {}, helpers }) {
  const root = helpers.element("label", "sandbox-kit sandbox-kit-slider");
  const name = helpers.element("span", "sandbox-kit-label", label);
  const input = helpers.element("input", "sandbox-kit-range");
  input.type = "range";
  input.min = String(min);
  input.max = String(max);
  input.step = String(step);
  input.value = String(currentValue);
  input.setAttribute("aria-label", label);
  const output = helpers.element("output", "sandbox-kit-value", String(currentValue));
  output.setAttribute("aria-live", "polite");
  input.oninput = () => {
    currentValue = Number(input.value);
    output.textContent = String(currentValue);
    onChange(currentValue);
  };
  root.append(name, input, output);
  return { node: root, input, output, get currentValue() { return currentValue; } };
}

function sandboxKitValueRow({ label, displayValue, helpers }) {
  const root = helpers.element("div", "sandbox-kit sandbox-kit-value-row");
  const name = helpers.element("span", "sandbox-kit-label", label);
  const output = helpers.element("bdi", "sandbox-kit-value", String(displayValue));
  output.setAttribute("aria-live", "polite");
  root.append(name, output);
  return {
    node: root,
    output,
    set(nextDisplayValue) {
      output.textContent = String(nextDisplayValue);
      output.classList.remove("sandbox-kit-value-changed");
      void output.offsetWidth;
      output.classList.add("sandbox-kit-value-changed");
    },
  };
}

function sandboxKitBoxes({ label, values, helpers }) {
  const root = helpers.element("div", "sandbox-kit sandbox-kit-boxes");
  const heading = helpers.element("p", "sandbox-kit-label", label);
  const row = helpers.element("ol", "sandbox-kit-box-row");
  row.setAttribute("aria-label", label);
  values.forEach((cellValue) => row.appendChild(helpers.element("li", "sandbox-kit-box", String(cellValue))));
  root.append(heading, row);
  return root;
}

function sandboxKitLogList({ label, entries, helpers }) {
  const root = helpers.element("div", "sandbox-kit sandbox-kit-log");
  const heading = helpers.element("p", "sandbox-kit-label", label);
  const list = helpers.element("ol", "sandbox-kit-log-list");
  list.setAttribute("aria-live", "polite");
  function set(nextItems) {
    list.replaceChildren(...nextItems.map((entryText) => helpers.element("li", "", String(entryText))));
  }
  set(entries);
  root.append(heading, list);
  return { node: root, list, set };
}

function sandboxKitCodeView({ lines, activeLine = -1, helpers }) {
  const pre = helpers.element("pre", "sandbox-kit sandbox-kit-code");
  const code = helpers.element("code", "");
  const rows = lines.map((line, index) => {
    const row = helpers.element("span", "sandbox-kit-code-line", line || " ");
    row.dataset.line = String(index + 1);
    code.appendChild(row);
    return row;
  });
  pre.appendChild(code);
  function setActiveLine(nextLine) {
    activeLine = nextLine;
    rows.forEach((row, index) => row.classList.toggle("current", index === activeLine));
  }
  setActiveLine(activeLine);
  return { node: pre, setActiveLine };
}

Lessons.kit = Object.freeze({
  stepper: sandboxKitStepper,
  toggle: sandboxKitToggle,
  slider: sandboxKitSlider,
  valueRow: sandboxKitValueRow,
  boxes: sandboxKitBoxes,
  logList: sandboxKitLogList,
  codeView: sandboxKitCodeView,
});
