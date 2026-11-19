"use strict";

Lessons.registerSandbox("ch14-generic-max", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => T(config.labels[key].ar, config.labels[key].en);
  const text = (pair) => T(pair.ar, pair.en);
  const panel = element("section", "lesson-sandbox ch14-sandbox ch14-max");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch14-frame");
      view.appendChild(element("p", "ch14-caption", text(frame.caption)));
      view.appendChild(element("p", "ch14-readout", `${label("type")}: ${text(frame.type)}`));
      const values = element("div", "ch14-values");
      frame.values.forEach((value, index) => {
        const cell = element("bdi", "ch14-value", String(value));
        cell.setAttribute("dir", "ltr");
        if (index === frame.current) cell.classList.add("current");
        if (index === frame.best) cell.classList.add("best");
        values.appendChild(cell);
      });
      view.appendChild(values);
      const current = frame.current < 0 ? label("none") : String(frame.values[frame.current]);
      const best = frame.best < 0 ? label("none") : String(frame.values[frame.best]);
      view.append(
        element("p", "ch14-readout", `${label("current")}: ${current}`),
        element("p", "ch14-readout", `${label("best")}: ${best}`),
      );
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});

Lessons.registerSandbox("ch14-stack", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const label = (key) => T(config.labels[key].ar, config.labels[key].en);
  const text = (pair) => T(pair.ar, pair.en);
  const panel = element("section", "lesson-sandbox ch14-sandbox ch14-stack");
  const generic = element("p", "ch14-readout ch14-code");
  generic.textContent = `${label("generic")}: ${text(config.generic)}`;
  const examples = element("div", "ch14-examples");
  config.examples.forEach((entry) => {
    examples.appendChild(element("p", "ch14-readout ch14-code", `${text(entry.type)} = ${text(entry.value)}`));
  });
  const history = element("div", "ch14-history");
  const top = element("p", "ch14-readout");
  const status = element("p", "ch14-caption");
  status.hidden = true;
  const controls = element("div", "ch14-controls");
  let route = [...config.initial];
  let arrivalIndex = 0;

  const push = button(label("push"), () => {
    if (arrivalIndex < config.arrivals.length) route.push(config.arrivals[arrivalIndex++]);
    status.textContent = label("pushed");
    status.hidden = false;
    update();
  });
  const peek = button(label("peek"), () => {
    status.textContent = `${label("top")}: ${route.length ? text(route[route.length - 1]) : label("empty")}`;
    status.hidden = false;
    update();
  });
  const pop = button(label("pop"), () => {
    if (route.length) route.pop();
    status.textContent = label("removed");
    status.hidden = false;
    update();
  });
  const reset = button(label("reset"), () => {
    route = [...config.initial];
    arrivalIndex = 0;
    status.textContent = label("resetDone");
    status.hidden = false;
    update();
  });
  controls.append(push, peek, pop, reset);
  panel.append(generic, element("p", "ch14-section-label", label("cells")), examples, element("p", "ch14-section-label", label("route")), history, top, status, controls);

  function update() {
    history.replaceChildren(...route.map((value) => {
      const cell = element("bdi", "ch14-stack-cell", text(value));
      cell.setAttribute("dir", "auto");
      return cell;
    }));
    top.textContent = `${label("top")}: ${route.length ? text(route[route.length - 1]) : label("empty")}`;
    push.disabled = arrivalIndex === config.arrivals.length;
    pop.disabled = route.length === 0;
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch14-type-constraints", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const label = (key) => T(config.labels[key].ar, config.labels[key].en);
  const panel = element("section", "lesson-sandbox ch14-sandbox ch14-constraints");
  const constraint = element("p", "ch14-readout ch14-code", `${label("constraint")}: ${config.constraint}`);
  const type = element("p", "ch14-readout ch14-code");
  const underlying = element("p", "ch14-readout ch14-code");
  const result = element("p", "ch14-result");
  const controls = element("div", "ch14-controls");
  let index = 0;
  const next = button(label("next"), () => { index = (index + 1) % config.types.length; update(); });
  const reset = button(label("reset"), () => { index = 0; update(); });
  controls.append(next, reset);
  panel.append(constraint, type, underlying, result, controls);

  function update() {
    const candidate = config.types[index];
    type.textContent = `${label("type")}: ${candidate.name}`;
    underlying.textContent = `${label("underlying")}: ${candidate.underlying}`;
    result.textContent = label(candidate.accepted ? "accepted" : "rejected");
    result.classList.toggle("accepted", candidate.accepted);
    result.classList.toggle("rejected", !candidate.accepted);
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch14-queue", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const label = (key) => T(config.labels[key].ar, config.labels[key].en);
  const text = (pair) => T(pair.ar, pair.en);
  const panel = element("section", "lesson-sandbox ch14-sandbox ch14-queue");
  const type = element("p", "ch14-readout ch14-code", `${text(config.elementType)}`);
  const queueLine = element("div", "ch14-queue-items");
  const length = element("p", "ch14-readout");
  const front = element("p", "ch14-readout");
  const status = element("p", "ch14-caption");
  status.hidden = true;
  const controls = element("div", "ch14-controls");
  let queue = [...config.initial];
  let arrivalIndex = 0;

  const arrive = button(label("add"), () => {
    if (arrivalIndex < config.arrivals.length) queue.push(config.arrivals[arrivalIndex++]);
    status.textContent = label("added");
    status.hidden = false;
    update();
  });
  const read = button(label("read"), () => {
    status.textContent = queue.length ? `${label("front")}: ${text(queue[0])}` : label("noRead");
    status.hidden = false;
    update();
  });
  const depart = button(label("depart"), () => {
    if (queue.length) queue.shift();
    status.textContent = label("departed");
    status.hidden = false;
    update();
  });
  const reset = button(label("reset"), () => {
    queue = [...config.initial];
    arrivalIndex = 0;
    status.textContent = label("resetDone");
    status.hidden = false;
    update();
  });
  controls.append(arrive, read, depart, reset);
  panel.append(type, element("p", "ch14-section-label", label("queue")), queueLine, length, front, status, controls);

  function update() {
    queueLine.replaceChildren(...queue.map((value) => {
      const cell = element("bdi", "ch14-queue-item", text(value));
      cell.setAttribute("dir", "auto");
      return cell;
    }));
    if (!queue.length) queueLine.appendChild(element("span", "ch14-empty", label("empty")));
    length.textContent = `${label("length")}: ${queue.length}`;
    front.textContent = `${label("front")}: ${queue.length ? text(queue[0]) : label("empty")}`;
    arrive.disabled = arrivalIndex === config.arrivals.length;
    depart.disabled = queue.length === 0;
  }

  update();
  return panel;
});

Lessons.registerSandbox("ch14-recursion", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => T(config.labels[key].ar, config.labels[key].en);
  const text = (pair) => T(pair.ar, pair.en);
  const panel = element("section", "lesson-sandbox ch14-sandbox ch14-recursion");
  const stepper = kit.stepper({
    frames: config.frames,
    labels: { previous: label("previous"), next: label("next"), reset: label("reset") },
    helpers: { element, button },
    renderFrame(frame) {
      const view = element("div", "ch14-frame");
      view.appendChild(element("p", "ch14-caption", text(frame.caption)));
      view.appendChild(element("p", "ch14-section-label", label("calls")));
      const calls = element("ol", "ch14-call-stack");
      frame.calls.forEach((call) => {
        const item = element("li", "ch14-code", call);
        item.setAttribute("dir", "ltr");
        calls.appendChild(item);
      });
      view.append(calls, element("p", "ch14-readout ch14-code", `${label("returned")}: ${frame.returned}`));
      return view;
    },
  });
  panel.appendChild(stepper.node);
  return panel;
});
