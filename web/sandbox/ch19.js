"use strict";

function ch19Text(T, pair) {
  return T(pair.ar, pair.en);
}

function ch19Clock(totalMinutes, dayLabel) {
  const day = Math.floor(totalMinutes / 1440);
  const minuteOfDay = ((totalMinutes % 1440) + 1440) % 1440;
  const hour = String(Math.floor(minuteOfDay / 60)).padStart(2, "0");
  const minute = String(minuteOfDay % 60).padStart(2, "0");
  return `${hour}:${minute}${day > 0 ? ` (${dayLabel})` : ""}`;
}

Lessons.registerSandbox("ch19-duration-add", (block, { element, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch19Text(T, config.labels[key]);
  const panel = element("section", "lesson-sandbox ch19-duration-add");
  const start = config.startHour * 60 + config.startMinute;
  const startRow = kit.valueRow({ label: label("start"), displayValue: ch19Clock(start, ""), helpers: { element } });
  const arrival = kit.valueRow({ label: label("arrival"), displayValue: "", helpers: { element } });
  const duration = kit.slider({
    label: label("duration"),
    currentValue: config.duration.value,
    min: config.duration.min,
    max: config.duration.max,
    step: config.duration.step,
    helpers: { element },
    onChange(value) {
      arrival.set(ch19Clock(start + value, label("nextDay")));
    },
  });
  arrival.set(ch19Clock(start + config.duration.value, label("nextDay")));
  panel.append(duration.node, startRow.node, arrival.node);
  return panel;
});

Lessons.registerSandbox("ch19-zone-view", (block, { element, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch19Text(T, config.labels[key]);
  const panel = element("section", "lesson-sandbox ch19-zone-view");
  const sameInstant = element("p", "ch19-playground-caption", label("instant"));
  const utc = kit.valueRow({ label: label("utc"), displayValue: "", helpers: { element } });
  const west = kit.valueRow({ label: label("west"), displayValue: "", helpers: { element } });
  const east = kit.valueRow({ label: label("east"), displayValue: "", helpers: { element } });
  const renderHour = (hour) => {
    const format = (localHour) => {
      const day = Math.floor(localHour / 24);
      const clockHour = ((localHour % 24) + 24) % 24;
      const suffix = day < 0 ? ` (${label("previousDay")})` : day > 0 ? ` (${label("nextDay")})` : "";
      return `${String(clockHour).padStart(2, "0")}:00${suffix}`;
    };
    utc.set(format(hour));
    west.set(format(hour + config.offsets.west));
    east.set(format(hour + config.offsets.east));
  };
  const slider = kit.slider({
    label: label("hour"),
    currentValue: config.hour.value,
    min: config.hour.min,
    max: config.hour.max,
    step: config.hour.step,
    helpers: { element },
    onChange: renderHour,
  });
  renderHour(config.hour.value);
  panel.append(sameInstant, slider.node, utc.node, west.node, east.node);
  return panel;
});

Lessons.registerSandbox("ch19-select-race", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch19Text(T, config.labels[key]);
  const panel = element("section", "lesson-sandbox ch19-select-race");
  const caption = element("p", "ch19-playground-caption", label("caption"));
  const log = kit.logList({ label: label("schedule"), entries: [], helpers: { element } });
  let runIndex = 0;
  const render = () => log.set(config.runs[runIndex].events.map((event) => ch19Text(T, event)));
  const runAgain = button(label("runAgain"), () => {
    runIndex = (runIndex + 1) % config.runs.length;
    render();
  });
  const reset = button(label("reset"), () => {
    runIndex = 0;
    render();
  });
  const controls = element("div", "ch19-playground-controls");
  controls.append(runAgain, reset);
  panel.append(caption, log.node, controls);
  render();
  return panel;
});

Lessons.registerSandbox("ch19-worker-order", (block, { element, button, T, kit }) => {
  const config = block.sandbox.config;
  const label = (key) => ch19Text(T, config.labels[key]);
  const panel = element("section", "lesson-sandbox ch19-worker-order");
  const caption = element("p", "ch19-playground-caption", label("caption"));
  const finish = kit.logList({ label: label("finishOrder"), entries: [], helpers: { element } });
  const requests = kit.logList({
    label: label("requestOrder"),
    entries: config.requestOrder.map((id) => ch19Text(T, config.workers.find((worker) => worker.id === id).label)),
    helpers: { element },
  });
  let runIndex = 0;
  const render = () => {
    const workersByID = new Map(config.workers.map((worker) => [worker.id, worker]));
    finish.set(config.runs[runIndex].finishOrder.map((id) => ch19Text(T, workersByID.get(id).label)));
  };
  const runAgain = button(label("runAgain"), () => {
    runIndex = (runIndex + 1) % config.runs.length;
    render();
  });
  const controls = element("div", "ch19-playground-controls");
  controls.appendChild(runAgain);
  panel.append(caption, finish.node, requests.node, controls);
  render();
  return panel;
});
