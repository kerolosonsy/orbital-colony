"use strict";

function ch30Text(value, T) {
  return T(value.ar, value.en);
}

function ch30Panel(element, kind) {
  return element("section", `lesson-sandbox ch30-playground ${kind}`);
}

function ch30Row(element, label, value) {
  const row = element("div", "ch30-row");
  row.append(element("span", "ch30-label", label), element("bdi", "ch30-value", String(value)));
  return row;
}

function ch30Cell(row) {
  return row.children[1];
}

function ch30Status(element, initial = "") {
  const status = element("p", "ch30-status", initial);
  status.setAttribute("role", "status");
  status.setAttribute("aria-live", "polite");
  return status;
}

Lessons.registerSandbox("ch30-cache-aside", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch30Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch30Panel(element, "ch30-cache-aside");
  const key = ch30Row(element, label("key"), text(config.key));
  const cache = ch30Row(element, label("cache"), "");
  const store = ch30Row(element, label("store"), text(config.value));
  const result = ch30Row(element, label("result"), "—");
  const status = ch30Status(element);
  let cachePresent = config.cachePresent;
  let lastSource = "";
  function update() {
    ch30Cell(cache).textContent = cachePresent ? text(config.value) : label("miss");
    ch30Cell(result).textContent = lastSource ? text(config.value) : "—";
    if (lastSource === "store") status.textContent = `${label("storeSource")} · ${label("filled")}`;
    else if (lastSource === "cache") status.textContent = label("cacheSource");
    else status.textContent = cachePresent ? label("hit") : label("miss");
  }
  const controls = element("div", "ch30-controls");
  controls.append(
    button(label("read"), () => { lastSource = cachePresent ? "cache" : "store"; cachePresent = true; update(); }),
    button(label("evict"), () => { cachePresent = false; lastSource = ""; update(); }),
  );
  panel.append(key, cache, store, result, status, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch30-document-pages", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch30Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch30Panel(element, "ch30-document-pages");
  const page = ch30Row(element, label("page"), "1");
  const heading = element("p", "ch30-label", label("documents"));
  const documents = element("ol", "ch30-document-list");
  const more = ch30Status(element);
  let offset = 0;
  function update() {
    const pageRows = config.rows.slice(offset, offset + config.pageSize);
    ch30Cell(page).textContent = String(Math.floor(offset / config.pageSize) + 1);
    documents.replaceChildren(...pageRows.map((row) => element("li", "", text(row))));
    more.textContent = `${label("more")}: ${label(offset + config.pageSize < config.rows.length ? "yes" : "no")}`;
    next.disabled = offset + config.pageSize >= config.rows.length;
  }
  const controls = element("div", "ch30-controls");
  const next = button(label("next"), () => { if (offset + config.pageSize < config.rows.length) offset += config.pageSize; update(); });
  controls.append(next, button(label("reset"), () => { offset = 0; update(); }));
  panel.append(page, heading, documents, more, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch30-search-index", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch30Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch30Panel(element, "ch30-search-index");
  const term = ch30Row(element, label("term"), "");
  const heading = element("p", "ch30-label", label("postings"));
  const postings = element("ol", "ch30-document-list");
  const note = ch30Status(element, label("mapNote"));
  let index = 0;
  function update() {
    const entry = config.terms[index];
    ch30Cell(term).textContent = text(entry.term);
    postings.replaceChildren(...entry.postings.map((item) => element("li", "", text(item))));
  }
  const controls = element("div", "ch30-controls");
  controls.appendChild(button(label("lookup"), () => { index = (index + 1) % config.terms.length; update(); }));
  panel.append(term, heading, postings, note, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch30-read-write-window", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch30Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch30Panel(element, "ch30-read-write-window");
  const clock = ch30Row(element, label("clock"), "");
  const ownRoute = ch30Row(element, label("ownSession"), "");
  const otherRoute = ch30Row(element, label("otherSession"), label("replica"));
  let now = config.now;
  let stickyUntil = config.now;
  function update() {
    ch30Cell(clock).textContent = `${now} ${label("tick")}`;
    ch30Cell(ownRoute).textContent = label(now < stickyUntil ? "primary" : "replica");
  }
  const controls = element("div", "ch30-controls");
  controls.append(
    button(label("write"), () => { stickyUntil = now + config.stickyTicks; update(); }),
    button(label("advance"), () => { now++; update(); }),
    button(label("reset"), () => { now = config.now; stickyUntil = config.now; update(); }),
  );
  panel.append(clock, ownRoute, otherRoute, controls);
  update();
  return panel;
});

Lessons.registerSandbox("ch30-commit-cache-boundary", (block, { element, button, T }) => {
  const config = block.sandbox.config;
  const text = (value) => ch30Text(value, T);
  const label = (key) => text(config.labels[key]);
  const panel = ch30Panel(element, "ch30-commit-cache-boundary");
  const draft = ch30Row(element, label("draft"), config.proposed);
  const saved = ch30Row(element, label("saved"), config.before);
  const cache = ch30Row(element, label("cache"), label("fresh"));
  const boundary = ch30Status(element, label("notCommitted"));
  let committed = false;
  let failed = false;
  function update() {
    ch30Cell(saved).textContent = String(committed ? config.proposed : config.before);
    ch30Cell(cache).textContent = label(failed ? "warning" : "fresh");
    boundary.textContent = committed ? label("savedBoundary") : label("notCommitted");
    commit.disabled = committed;
    failCache.disabled = !committed || failed;
  }
  const controls = element("div", "ch30-controls");
  const commit = button(label("commit"), () => { committed = true; update(); });
  const failCache = button(label("failCache"), () => { failed = true; update(); });
  controls.append(commit, failCache, button(label("reset"), () => { committed = false; failed = false; update(); }));
  panel.append(draft, saved, cache, boundary, controls);
  update();
  return panel;
});
