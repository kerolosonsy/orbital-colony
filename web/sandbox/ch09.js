"use strict";

function ch09Text(localizedText, T) {
  return T(localizedText.ar, localizedText.en);
}

function ch09Panel(element, kind, config, T) {
  const panel = element("section", `lesson-sandbox ${kind}`);
  const intro = element("p", "ch09-sandbox-intro", ch09Text(config.intro, T));
  intro.setAttribute("dir", "auto");
  panel.appendChild(intro);
  return panel;
}

function ch09ValueRow(kit, element, rowLabel, localizedValue, T) {
  return kit.valueRow({
    label: ch09Text(rowLabel, T),
    displayValue: ch09Text(localizedValue, T),
    helpers: { element },
  });
}

Lessons.registerSandbox("ch09-stringer", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = ch09Panel(element, "ch09-stringer", config, T);
  let showingString = false;
  const output = ch09ValueRow(kit, element, config.outputLabel, config.default, T);
  const action = button(ch09Text(config.showString, T), () => {
    showingString = !showingString;
    output.set(ch09Text(showingString ? config.string : config.default, T));
    action.textContent = ch09Text(showingString ? config.showDefault : config.showString, T);
  });
  panel.append(output.node, action);
  return panel;
});

Lessons.registerSandbox("ch09-plan-callback", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = ch09Panel(element, "ch09-plan-callback", config, T);
  const tiles = config.tiles;
  let index = 0;
  const tile = ch09ValueRow(kit, element, config.tileLabel, tiles[0].tile, T);
  const actionCount = ch09ValueRow(kit, element, config.resultLabel, tiles[0].result, T);
  const caption = element("p", "ch09-sandbox-caption", ch09Text(tiles[0].caption, T));
  caption.setAttribute("dir", "auto");
  const next = button(ch09Text(config.button, T), () => {
    index = (index + 1) % tiles.length;
    tile.set(ch09Text(tiles[index].tile, T));
    actionCount.set(ch09Text(tiles[index].result, T));
    caption.textContent = ch09Text(tiles[index].caption, T);
  });
  panel.append(tile.node, actionCount.node, caption, next);
  return panel;
});

Lessons.registerSandbox("ch09-type-switch", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = ch09Panel(element, "ch09-type-switch", config, T);
  const values = config.values;
  let index = 0;
  const currentValue = ch09ValueRow(kit, element, config.valueLabel, values[0].value, T);
  const type = ch09ValueRow(kit, element, config.typeLabel, values[0].type, T);
  const selectedCase = ch09ValueRow(kit, element, config.caseLabel, values[0].case, T);
  const effect = ch09ValueRow(kit, element, config.effectLabel, values[0].effect, T);
  const next = button(ch09Text(config.button, T), () => {
    index = (index + 1) % values.length;
    currentValue.set(ch09Text(values[index].value, T));
    type.set(ch09Text(values[index].type, T));
    selectedCase.set(ch09Text(values[index].case, T));
    effect.set(ch09Text(values[index].effect, T));
  });
  panel.append(currentValue.node, type.node, selectedCase.node, effect.node, next);
  return panel;
});

Lessons.registerSandbox("ch09-trim-space", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = ch09Panel(element, "ch09-trim-space", config, T);
  let trimmed = false;
  const before = ch09ValueRow(kit, element, config.inputLabel, config.raw, T);
  const after = ch09ValueRow(kit, element, config.outputLabel, config.raw, T);
  const action = button(ch09Text(config.trim, T), () => {
    trimmed = !trimmed;
    after.set(ch09Text(trimmed ? config.trimmed : config.raw, T));
    action.textContent = ch09Text(trimmed ? config.restore : config.trim, T);
  });
  panel.append(before.node, after.node, action);
  return panel;
});

Lessons.registerSandbox("ch09-map-updates", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const panel = ch09Panel(element, "ch09-map-updates", config, T);
  const steps = config.steps;
  let index = 0;
  const map = ch09ValueRow(kit, element, config.mapLabel, steps[0].map, T);
  const caption = element("p", "ch09-sandbox-caption", ch09Text(steps[0].caption, T));
  caption.setAttribute("dir", "auto");
  const next = button(ch09Text(config.button, T), () => {
    index = (index + 1) % steps.length;
    map.set(ch09Text(steps[index].map, T));
    caption.textContent = ch09Text(steps[index].caption, T);
  });
  panel.append(map.node, caption, next);
  return panel;
});
