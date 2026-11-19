"use strict";

function renderCh18Sandbox(block, { element, button, T }) {
  const config = block.sandbox.config;
  const { labels, data, mode } = config;
  const label = (key) => T(labels[key].ar, labels[key].en);
  const localized = (pair) => T(pair.ar, pair.en);
  const panel = element("section", `lesson-sandbox ch18-sandbox ch18-${mode}`);
  panel.setAttribute("dir", T("rtl", "ltr"));
  const state = element("div", "ch18-sandbox-state");
  state.setAttribute("role", "status");
  state.setAttribute("aria-live", "polite");
  const controls = element("div", "ch18-sandbox-controls");
  panel.append(state, controls);

  function row(name, shownValue, direction = "auto") {
    const line = element("div", "ch18-sandbox-row");
    line.appendChild(element("span", "ch18-sandbox-label", name));
    const output = element("bdi", "ch18-sandbox-value", String(shownValue));
    output.setAttribute("dir", direction);
    line.appendChild(output);
    return line;
  }

  function caption(text) {
    return element("p", "ch18-sandbox-caption", text);
  }

  function list(items) {
    const output = element("ol", "ch18-sandbox-list");
    for (const item of items) output.appendChild(element("li", "ch18-sandbox-item", item));
    return output;
  }

  if (mode === "index") {
    let index = 0;
    const action = button(label("action"), () => {
      index = (index + 1) % data.values.length;
      update();
    });
    controls.appendChild(action);
    function update() {
      state.replaceChildren(
        row(label("values"), `[${data.values.join(", ")}]`, "ltr"),
        row(label("read"), `values[${index}] = ${data.values[index]}`, "ltr"),
        caption(localized(labels.caption[index])),
      );
    }
    update();
  } else if (mode === "json") {
    let encoded = false;
    const action = button(label("action"), () => {
      encoded = !encoded;
      update();
    });
    controls.appendChild(action);
    function update() {
      state.replaceChildren(
        row(label("source"), data.source, "ltr"),
        row(label("result"), encoded ? data.encoded : label("before"), "ltr"),
        caption(label(encoded ? "after" : "before")),
      );
    }
    update();
  } else if (mode === "csv") {
    let rowIndex = 0;
    const action = button(label("action"), () => {
      rowIndex = (rowIndex + 1) % data.rows.length;
      update();
    });
    controls.appendChild(action);
    function update() {
      const current = data.rows[rowIndex];
      state.replaceChildren(
        row(label("source"), localized(current.source), "auto"),
        element("p", "ch18-sandbox-heading", label("fields")),
        list(current.fields.map(localized)),
        caption(localized(labels.captions[rowIndex])),
      );
    }
    update();
  } else if (mode === "unicode") {
    let split = false;
    let lower = false;
    const splitButton = button(label("split"), () => {
      split = !split;
      splitButton.setAttribute("aria-pressed", String(split));
      update();
    });
    const lowerButton = button(label("lower"), () => {
      lower = !lower;
      lowerButton.setAttribute("aria-pressed", String(lower));
      update();
    });
    splitButton.setAttribute("aria-pressed", "false");
    lowerButton.setAttribute("aria-pressed", "false");
    controls.append(splitButton, lowerButton);
    function update() {
      const shownSource = lower ? data.lowerRaw : data.raw;
      const words = split ? (lower ? data.lowerWords : data.words).map(localized) : null;
      state.replaceChildren(
        row(label("source"), localized(shownSource), "auto"),
        ...(words ? [element("p", "ch18-sandbox-heading", label("result")), list(words)] : []),
        caption(lower ? label("lowerCaption") : split ? label("splitCaption") : label("rawCaption")),
      );
    }
    update();
  } else if (mode === "splitn") {
    let limitIndex = 0;
    const action = button(label("action"), () => {
      limitIndex = (limitIndex + 1) % data.limits.length;
      update();
    });
    controls.appendChild(action);
    function splitAtLimit(raw, separator, limit) {
      const parts = [];
      let start = 0;
      while (parts.length < limit - 1) {
        const next = raw.indexOf(separator, start);
        if (next < 0) break;
        parts.push(raw.slice(start, next));
        start = next + separator.length;
      }
      parts.push(raw.slice(start));
      return parts;
    }
    function update() {
      const limit = data.limits[limitIndex];
      state.replaceChildren(
        row(label("source"), data.raw, "ltr"),
        row(label("limit"), limit, "ltr"),
        element("p", "ch18-sandbox-heading", label("parts")),
        list(splitAtLimit(data.raw, data.separator, limit)),
        caption(label("caption")),
      );
    }
    update();
  } else if (mode === "search") {
    let searched = false;
    const action = button(label("action"), () => {
      searched = !searched;
      update();
    });
    controls.appendChild(action);
    function update() {
      state.replaceChildren(
        row(label("name"), searched ? data.normalized : data.name, "ltr"),
        row(label("query"), data.query, "ltr"),
        row(label("result"), searched ? data.match : label("pending"), "ltr"),
        caption(label("caption")),
      );
    }
    update();
  }

  return panel;
}

Lessons.registerSandbox("ch18-slice-index", renderCh18Sandbox);
Lessons.registerSandbox("ch18-json-marshal", renderCh18Sandbox);
Lessons.registerSandbox("ch18-csv-quoted-field", renderCh18Sandbox);
Lessons.registerSandbox("ch18-unicode-words", renderCh18Sandbox);
Lessons.registerSandbox("ch18-splitn-limit", renderCh18Sandbox);
Lessons.registerSandbox("ch18-normalized-search", renderCh18Sandbox);
