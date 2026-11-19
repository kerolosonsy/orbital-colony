"use strict";

function ch10Text(localizedText, T) {
  return T(localizedText.ar, localizedText.en);
}

function ch10AppendGroups(panel, config, element, T) {
  if (!config.groups) return;
  panel.appendChild(element("p", "ch10-sandbox-label", ch10Text(config.groupsLabel, T)));
  const groups = element("ul", "ch10-sandbox-groups");
  config.groups.forEach((group) => {
    const groupRow = element("li", "");
    groupRow.appendChild(element("bdi", "", ch10Text(group, T)));
    groups.appendChild(groupRow);
  });
  panel.appendChild(groups);
}

function ch10AdvanceTimeline(config, T, log, timeline) {
  let events = config.runs[timeline.runIndex].events;
  if (timeline.eventIndex === events.length) {
    timeline.runIndex = (timeline.runIndex + 1) % config.runs.length;
    timeline.eventIndex = 0;
    timeline.visibleEvents = [];
    events = config.runs[timeline.runIndex].events;
  }
  timeline.visibleEvents.push(ch10Text(events[timeline.eventIndex], T));
  timeline.eventIndex++;
  log.set(timeline.visibleEvents);
  log.list.children.forEach((entry) => entry.setAttribute("dir", "auto"));
}

function ch10Timeline(kind, block, { element, button, T, kit }) {
  const config = block.sandbox.config;
  const panel = element("section", `lesson-sandbox ${kind}`);
  const intro = element("p", "ch10-sandbox-intro", ch10Text(config.intro, T));
  intro.setAttribute("dir", "auto");
  panel.appendChild(intro);
  ch10AppendGroups(panel, config, element, T);
  const log = kit.logList({ label: ch10Text(config.logLabel, T), entries: [], helpers: { element } });
  log.node.querySelector(".sandbox-kit-label").setAttribute("dir", "auto");
  const timeline = { runIndex: 0, eventIndex: 0, visibleEvents: [] };
  const next = button(ch10Text(config.button, T), () => ch10AdvanceTimeline(config, T, log, timeline));
  panel.append(log.node, next);
  return panel;
}

Lessons.registerSandbox("ch10-waitgroup-order", (block, helpers) =>
  ch10Timeline("ch10-waitgroup-order", block, helpers));

Lessons.registerSandbox("ch10-swarm-order", (block, helpers) =>
  ch10Timeline("ch10-swarm-order", block, helpers));

Lessons.registerSandbox("ch10-groups-order", (block, helpers) =>
  ch10Timeline("ch10-groups-order", block, helpers));
