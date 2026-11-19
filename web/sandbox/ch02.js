"use strict";

function ch02Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch02StepperLabels(config, T) {
  return Object.fromEntries(["previous", "next", "reset"].map((key) => [key, ch02Text(config.labels[key], T)]));
}

Lessons.registerSandbox("ch02-inventory-flow", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch02Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch02-sandbox ch02-inventory-flow");
  const destination = kit.valueRow({ label: label("output"), displayValue: label("console"), helpers: { element } });
  const stepper = kit.stepper({
    frames: config.frames,
    labels: ch02StepperLabels(config, T),
    helpers: { element, button },
    renderFrame(frame) {
      const values = element("div", "ch02-inventory-values");
      for (const key of ["moss", "crystal", "total"]) {
        values.appendChild(kit.valueRow({
          label: label(key),
          displayValue: frame[key] === null ? label("empty") : frame[key],
          helpers: { element },
        }).node);
      }
      values.appendChild(element("p", "ch02-caption", ch02Text(frame.caption, T)));
      return values;
    },
  });
  const outputToggle = kit.toggle({
    label: label("output"),
    onLabel: label("base"),
    offLabel: label("console"),
    helpers: { element, button },
    onChange(isBase) { destination.set(label(isBase ? "base" : "console")); },
  });
  panel.append(stepper.node, destination.node, outputToggle.node);
  return panel;
});

Lessons.registerSandbox("ch02-return-pair", (block, { element, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch02Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch02-sandbox ch02-return-pair");
  const grid = element("div", "ch02-survey-grid");
  const returned = element("div", "ch02-returned-values");
  const widthValue = kit.valueRow({ label: label("firstReturn"), displayValue: config.width.value, helpers: { element } });
  const heightValue = kit.valueRow({ label: label("secondReturn"), displayValue: config.height.value, helpers: { element } });
  const drawMap = (width, height) => {
    grid.replaceChildren();
    grid.style.gridTemplateColumns = `repeat(${width}, minmax(0, 1fr))`;
    for (let index = 0; index < width * height; index++) grid.appendChild(element("span", "ch02-survey-cell"));
  };
  drawMap(config.width.value, config.height.value);
  returned.append(widthValue.node, heightValue.node);
  const widthSlider = kit.slider({
    label: label("width"),
    currentValue: config.width.value,
    min: config.width.min,
    max: config.width.max,
    helpers: { element },
    onChange(value) { widthValue.set(value); drawMap(value, heightSlider.currentValue); },
  });
  const heightSlider = kit.slider({
    label: label("height"),
    currentValue: config.height.value,
    min: config.height.min,
    max: config.height.max,
    helpers: { element },
    onChange(value) { heightValue.set(value); drawMap(widthSlider.currentValue, value); },
  });
  panel.append(element("p", "ch02-caption", label("caption")), grid, returned, widthSlider.node, heightSlider.node);
  return panel;
});

Lessons.registerSandbox("ch02-utf8-span", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch02Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch02-sandbox ch02-utf8-span");
  let selected = 0;
  const selectedValue = element("div", "ch02-utf8-selected");
  const byteRow = element("div", "ch02-utf8-byte-row");
  const choices = element("div", "ch02-utf8-choices");
  const markers = element("div", "ch02-utf8-markers");
  const counts = element("div", "ch02-utf8-counts");
  counts.append(
    kit.valueRow({ label: label("bytes"), displayValue: config.runes.reduce((sum, rune) => sum + rune.bytes.length, 0), helpers: { element } }).node,
    kit.valueRow({ label: label("runes"), displayValue: config.runes.length, helpers: { element } }).node,
  );
  function marker(key, runeIndex) {
    const row = element("div", "ch02-utf8-marker" + (selected === runeIndex ? " active" : ""));
    const value = element("bdi", "ch02-rune-value", config.runes[runeIndex].text);
    value.dir = "ltr";
    row.append(element("span", "ch02-label", label(key)), value);
    return row;
  }
  function refresh() {
    const selectedRune = config.runes[selected];
    selectedValue.replaceChildren(
      element("span", "ch02-label", label("selected")),
      element("bdi", "ch02-rune-value", selectedRune.text),
    );
    byteRow.replaceChildren();
    config.runes.forEach((rune, runeIndex) => {
      const group = element("div", "ch02-utf8-group" + (selected === runeIndex ? " selected" : ""));
      rune.bytes.forEach((hex) => group.appendChild(element("bdi", "ch02-utf8-byte", hex)));
      byteRow.appendChild(group);
    });
    choices.replaceChildren();
    config.runes.forEach((rune, runeIndex) => {
      const choice = button(`${label("select")} ${rune.text}`, () => {
        selected = runeIndex;
        refresh();
      });
      choice.setAttribute("aria-pressed", String(selected === runeIndex));
      choices.appendChild(choice);
    });
    markers.replaceChildren(
      marker("first", 0),
      marker("last", config.runes.length - 1),
    );
  }
  panel.append(element("p", "ch02-caption", label("caption")), selectedValue, counts, byteRow, choices, markers);
  refresh();
  return panel;
});

Lessons.registerSandbox("ch02-unit-align", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch02Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch02-sandbox ch02-unit-align");
  const source = element("div", "ch02-cargo-source");
  config.cargo.forEach((item) => {
    source.appendChild(kit.valueRow({ label: ch02Text(item.label, T), displayValue: item.count, helpers: { element } }).node);
  });
  const choices = element("div", "ch02-unit-choices");
  const conversions = element("div", "ch02-unit-conversions");
  const empty = element("p", "ch02-caption ch02-unit-empty", label("empty"));
  const total = kit.valueRow({ label: label("total"), displayValue: label("empty"), helpers: { element } });
  function choose(unit) {
    empty.hidden = true;
    choices.querySelectorAll("button").forEach((node) => node.setAttribute("aria-pressed", String(node.getAttribute("data-unit") === unit.id)));
    conversions.replaceChildren();
    const totalCells = config.cargo.reduce((sum, item) => sum + item.count * item.cellsEach, 0);
    config.cargo.forEach((item) => {
      const row = element("div", "ch02-unit-conversion-row");
      const sourceValue = element("bdi", "ch02-unit-number", `${item.count} ${ch02Text(item.label, T)}`);
      const equals = element("span", "ch02-unit-equals", config.labels.equals);
      const targetValue = element("bdi", "ch02-unit-number", `${(item.count * item.cellsEach) / unit.cellsPerUnit} ${ch02Text(unit.label, T)}`);
      row.append(sourceValue, equals, targetValue);
      conversions.appendChild(row);
    });
    total.set(`${totalCells / unit.cellsPerUnit} ${ch02Text(unit.label, T)}`);
  }
  config.units.forEach((unit) => {
    const choice = button(ch02Text(unit.label, T), () => choose(unit));
    choice.classList.add("ch02-unit-choice");
    choice.setAttribute("data-unit", unit.id);
    choice.setAttribute("aria-pressed", "false");
    choices.appendChild(choice);
  });
  panel.append(
    element("p", "ch02-caption", label("caption")),
    source,
    choices,
    empty,
    conversions,
    total.node,
  );
  return panel;
});
