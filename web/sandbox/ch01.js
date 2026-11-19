"use strict";

function ch01Text(pair, T) {
  return T(pair.ar, pair.en);
}

function ch01StepperLabels(config, T) {
  return Object.fromEntries(["previous", "next", "reset"].map((key) => [key, ch01Text(config.labels[key], T)]));
}

Lessons.registerSandbox("ch01-route-stepper", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch01Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch01-sandbox ch01-route-stepper");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: ch01StepperLabels(config, T),
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch01-route-frame");
      const command = element("div", "ch01-route-command");
      command.append(element("span", "ch01-label", label("command")), element("bdi", "ch01-value", ch01Text(frame.command, T)));
      const position = element("div", "ch01-route-position");
      const coordinate = element("bdi", "ch01-coordinate", `(${frame.x}, ${frame.y})`);
      coordinate.dir = "ltr";
      position.append(element("span", "ch01-label", label("position")), coordinate);
      view.append(command, position, element("p", "ch01-caption", ch01Text(frame.caption, T)));
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch01-harvest-call", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch01Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch01-sandbox ch01-harvest-call");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: ch01StepperLabels(config, T),
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch01-harvest-frame");
      const board = element("div", "ch01-tile-row");
      config.tiles.forEach((tile, index) => {
        const cell = element("div", "ch01-tile" + (frame.position === index ? " occupied" : ""));
        cell.appendChild(element("span", "ch01-tile-name", ch01Text(tile.label, T)));
        if (frame.position === index) cell.appendChild(element("span", "ch01-tile-state", label("drone")));
        if (tile.plant) {
          const harvested = frame.harvested.includes(index);
          cell.appendChild(element("span", "ch01-tile-state", label(harvested ? "collected" : "plant")));
        }
        board.appendChild(cell);
      });
      const action = element("div", "ch01-route-command");
      action.append(element("span", "ch01-label", label("command")), element("bdi", "ch01-value", ch01Text(frame.action, T)));
      const count = element("div", "ch01-route-position");
      count.append(element("span", "ch01-label", label("harvested")), element("bdi", "ch01-coordinate", String(frame.harvested.length)));
      view.append(board, action, count, element("p", "ch01-caption", ch01Text(frame.caption, T)));
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch01-planting-state", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch01Text(config.labels[key], T);
  const panel = element("section", "lesson-sandbox ch01-sandbox ch01-planting-state");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: ch01StepperLabels(config, T),
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch01-plant-frame");
      const ground = element("div", "ch01-state-row");
      ground.append(element("span", "ch01-label", label("ground")), element("bdi", "ch01-value", label(frame.ground)));
      const plant = element("div", "ch01-state-row");
      plant.append(element("span", "ch01-label", label("plant")), element("bdi", "ch01-value", label(frame.plant)));
      const action = element("div", "ch01-route-command");
      action.append(element("span", "ch01-label", label("action")), element("bdi", "ch01-value", ch01Text(frame.action, T)));
      view.append(ground, plant, action, element("p", "ch01-caption", ch01Text(frame.caption, T)));
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});
