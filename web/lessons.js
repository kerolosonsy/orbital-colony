"use strict";

const Lessons = (() => {
  let navigationObserver = null;
  function element(tag, className, text) {
    const node = document.createElement(tag);
    if (className) node.className = className;
    if (text !== undefined) node.textContent = text;
    return node;
  }

  function button(label, action) {
    const node = element("button", "tool", label);
    node.type = "button";
    node.onclick = action;
    return node;
  }

  function createSession(level, options) {
    const questions = new Map(level.steps.flatMap((step) => step.blocks.filter((b) => b.question).map((b) => [b.id, b.question])));
    const record = restoreLessonRecord(level, options.review ? null : options.saved, questions);
    const persist = () => {
      if (!options.review) options.save({ ...record, answers: { ...record.answers }, wrongAttempts: { ...record.wrongAttempts } });
    };
    return {
      record,
      get index() { return level.steps.findIndex((step) => step.id === record.step); },
      go(index) { record.step = level.steps[index].id; persist(); },
      answer(id, choice) {
        record.answers[id] = choice;
        if (choice !== questions.get(id).answer) record.wrongAttempts[id] = Math.min(2, (record.wrongAttempts[id] || 0) + 1);
        persist();
      },
      retry(id) { delete record.answers[id]; persist(); },
    };
  }

  function restoreLessonRecord(level, saved, questions) {
    const record = { revision: level.lessonRevision, step: level.steps[0].id, answers: {}, wrongAttempts: {} };
    if (saved?.revision !== record.revision) return record;
    if (level.steps.some((step) => step.id === saved.step)) record.step = saved.step;
    for (const [id, choice] of Object.entries(saved.answers || {})) {
      if (Number.isInteger(choice) && choice >= 0 && choice < (questions.get(id)?.options.length || 0)) record.answers[id] = choice;
    }
    for (const [id, question] of questions) {
      const count = saved.wrongAttempts?.[id];
      if (Number.isInteger(count) && count >= 0) record.wrongAttempts[id] = Math.min(count, 2);
      if (record.answers[id] !== undefined && record.answers[id] !== question.answer) record.wrongAttempts[id] = Math.max(1, record.wrongAttempts[id] || 0);
    }
    return record;
  }

  function prose(text) {
    const node = element("div", "lesson-copy");
    node.innerHTML = md(text);
    return node;
  }

  function codeView(code) {
    const pre = element("pre", "lesson-code");
    const source = element("code");
    code.split("\n").forEach((line, index) => {
      const row = element("span", "lesson-code-line");
      row.dataset.line = String(index + 1);
      row.innerHTML = highlight(line) || " ";
      source.appendChild(row);
    });
    pre.appendChild(source);
    return pre;
  }

  function codeContext(block) {
    return element("p", "lesson-model-label", T("جوه main؛ import: ", "Inside main; import: ") + block.imports.join(", "));
  }

  function isolated(text, direction, className) {
    const label = element("bdi", className, text);
    label.dir = direction;
    return label;
  }

  function sliceSummary(slice) {
    const caption = element("figcaption");
    caption.append(isolated(`${slice.name}:`, "ltr", "lesson-summary-go"), isolated(`len=${slice.length}`, "ltr", "lesson-summary-go"));
    if (slice.elementsOnly) return caption;
    if (slice.capacity === undefined) {
      const capacity = element("span", "lesson-summary-part");
      capacity.dir = "ltr";
      capacity.append(isolated("cap=", "ltr", "lesson-summary-go"), isolated(T("غير محددة", "not specified"), "auto"));
      caption.appendChild(capacity);
    } else caption.appendChild(isolated(`cap=${slice.capacity}`, "ltr", "lesson-summary-go"));
    if (slice.store) {
      const store = element("span", "lesson-summary-part");
      store.append(isolated(T("مخزن", "Store"), "auto"), isolated(slice.store, "ltr", "lesson-summary-go"));
      caption.appendChild(store);
    }
    if (slice.nil) caption.appendChild(isolated("nil", "ltr", "lesson-summary-go"));
    return caption;
  }

  function sliceView(slice) {
    const view = element("figure", "lesson-slice");
    view.appendChild(sliceSummary(slice));
    const cells = element("div", "lesson-cells");
    const count = slice.capacity === undefined ? slice.length : slice.capacity;
    for (let index = 0; index < count; index++) {
      const reserved = index >= slice.length;
      const cell = element("div", "lesson-cell" + (reserved ? " reserved" : ""));
      cell.append(element("span", "lesson-index", String(index)), element("b", "", reserved ? "—" : String(slice.values[index])));
      cells.appendChild(cell);
    }
    if (!count) cells.appendChild(element("p", "", T("مفيش عناصر.", "No elements.")));
    view.appendChild(cells);
    return view;
  }

  function frameView(frame) {
    const state = element("div", "lesson-frame");
    state.appendChild(element("p", "lesson-line-label", frame.line ? T(`السطر الحالي: ${frame.line}`, `Current line: ${frame.line}`) : T("قبل التنفيذ", "Before execution")));
    const caption = element("p", "lesson-caption", frame.caption);
    caption.setAttribute("role", "status");
    state.appendChild(caption);
    if (frame.variables?.length) {
      const vars = element("dl", "lesson-variables");
      frame.variables.forEach((variable) => vars.append(element("dt", "", variable.name), element("dd", "", variable.value)));
      state.appendChild(vars);
    }
    if (frame.slice) state.appendChild(sliceView(frame.slice));
    state.appendChild(element("p", "lesson-output-label", T("الخرج لحد دلوقتي", "Output so far")));
    state.appendChild(element("pre", "lesson-output", frame.output || T("مفيش خرج.", "No output.")));
    if (frame.panic) state.appendChild(element("p", "lesson-panic", `Panic · ${T("التنفيذ وقف هنا", "Execution stopped here")}: ${frame.panic}`));
    return state;
  }

  function walkthrough(block) {
    const panel = element("section", "lesson-walkthrough");
    const code = codeView(block.code);
    const state = element("div");
    const count = element("span", "lesson-position");
    let index = 0;
    const frames = block.walkthrough.frames;
    const previous = walkthroughButton(T("السطر اللي فات", "Previous line"), () => { index--; update(); });
    const next = walkthroughButton(T("السطر الجاي", "Next line"), () => { index++; update(); });
    const reset = walkthroughButton(T("من الأول", "Restart"), () => { index = 0; update(); });
    function update() {
      const frame = frames[index];
      highlightLine(code, frame.line);
      state.replaceChildren(frameView(frame));
      count.textContent = `${index + 1} / ${frames.length}`;
      previous.disabled = index === 0;
      next.disabled = index === frames.length - 1;
    }
    const controls = element("div", "lesson-controls");
    controls.setAttribute("role", "group");
    controls.setAttribute("aria-label", T("شرح الكود خطوة بخطوة", "Code walkthrough"));
    controls.append(previous, count, next, reset);
    panel.append(element("p", "lesson-model-label", T("شرح توضيحي خطوة بخطوة.", "Illustrated walkthrough.")), codeContext(block), code, controls, state);
    update();
    return panel;
  }

  function walkthroughButton(label, action) {
    const control = button(label, action);
    control.setAttribute("aria-label", T(`شرح الكود: ${label}`, `Walkthrough: ${label}`));
    return control;
  }

  function highlightLine(code, number) {
    code.querySelectorAll(".lesson-code-line").forEach((line) => {
      const current = Number(line.dataset.line) === number;
      line.classList.toggle("current", current);
      if (!current) return;
      // Only the code box scrolls; instant changes also honor reduced motion.
      if (line.offsetTop < code.scrollTop) code.scrollTop = line.offsetTop;
      else if (line.offsetTop + line.offsetHeight > code.scrollTop + code.clientHeight) code.scrollTop = line.offsetTop + line.offsetHeight - code.clientHeight;
    });
    if (!number) code.scrollTop = 0;
  }

  function question(block, session) {
    const field = element("fieldset", "lesson-question");
    field.appendChild(element("legend", "", T("اختار إجابتك", "Choose your answer")));
    const feedback = element("p", "lesson-feedback");
    feedback.setAttribute("role", "status");
    const choices = block.question.options.map((option, choice) => {
      const node = button("", () => { session.answer(block.id, choice); update(); });
      node.classList.add("lesson-choice");
      node.innerHTML = inline(option);
      return node;
    });
    const retry = button(T("جرّب تاني", "Try again"), () => { session.retry(block.id); update(); choices[0].focus(); });
    function update() {
      const answer = session.record.answers[block.id];
      const answered = answer !== undefined;
      const revealed = answer === block.question.answer || session.record.wrongAttempts[block.id] >= 2;
      choices.forEach((choice, index) => updateQuestionChoice(choice, block.question, index, { answer, revealed }));
      feedback.textContent = answered ? `${answer === block.question.answer ? T("صح!", "Correct!") : T("مش صح.", "Incorrect.")} ${block.question.feedback[answer]}` : "";
      retry.hidden = !answered;
    }
    field.append(...choices, feedback, retry);
    update();
    return field;
  }

  function updateQuestionChoice(choice, question, index, state) {
    const incorrect = index === state.answer && state.answer !== question.answer;
    choice.innerHTML = inline(question.options[index]);
    choice.disabled = state.answer !== undefined;
    choice.classList.toggle("correct", state.revealed && index === question.answer);
    choice.classList.toggle("incorrect", incorrect);
    if (state.revealed && index === question.answer) {
      choice.appendChild(element("span", "lesson-answer-label", T("الإجابة الصحيحة", "Correct answer")));
      if (state.answer !== index) choice.appendChild(element("span", "lesson-answer-explanation", question.feedback[index]));
    }
    if (incorrect) choice.appendChild(element("span", "lesson-answer-label", T("إجابتك", "Your answer")));
  }

  function createSandbox(spec) {
    let presetIndex = 0;
    let slice;
    let storeNumber;
    function reset(index = presetIndex) {
      presetIndex = index;
      storeNumber = spec.presets[index].capacity ? 1 : 0;
      slice = { name: "s", values: [...spec.presets[index].values], length: spec.presets[index].values.length, capacity: spec.presets[index].capacity, store: storeNumber ? "A" : "" };
    }
    reset();
    return {
      get slice() { return slice; },
      get presetIndex() { return presetIndex; },
      reset,
      append() {
        if (slice.length === spec.limit) return;
        if (slice.length === slice.capacity) {
          slice.capacity = Math.max(slice.length + 1, Math.max(1, 2 * slice.capacity));
          slice.store = String.fromCharCode(65 + storeNumber++);
        }
        slice.values.push(spec.appendValue);
        slice.length++;
      },
    };
  }

  function sandbox(block) {
    const spec = block.sandbox;
    if (spec.kind === "slice-append-sequence") return sequenceSandbox(spec);
    if (spec.kind === "slice-sort-clone") return sortCloneSandbox(spec);
    if (spec.kind === "map-count") return mapCountSandbox(spec);
    const model = createSandbox(spec);
    const panel = element("section", "lesson-sandbox");
    const presets = element("div", "lesson-presets");
    const view = element("div");
    const status = element("p", "lesson-caption", spec.presets[0].text);
    status.setAttribute("role", "status");
    const limit = element("p", "lesson-limit", spec.limitLabel);
    const append = button(T("ضيف 7", "Append 7"), () => { model.append(); status.textContent = spec.appendText; update(); });
    const reset = button(T("ابدأ من الأول", "Reset"), () => { model.reset(); status.textContent = spec.resetText; update(); });
    spec.presets.forEach((preset, index) => presets.appendChild(button(preset.label, () => { model.reset(index); status.textContent = preset.text; update(); })));
    function update() {
      view.replaceChildren(sliceView(model.slice));
      append.disabled = model.slice.length === spec.limit;
      limit.hidden = !append.disabled;
      presets.querySelectorAll("button").forEach((node, index) => node.setAttribute("aria-pressed", String(index === model.presetIndex)));
    }
    const controls = element("div", "lesson-controls");
    controls.append(append, reset);
    panel.append(presets, view, element("p", "lesson-reserved-label", spec.reservedLabel), element("p", "lesson-growth-label", spec.growthLabel), status, limit, controls);
    update();
    return panel;
  }

  function sortCloneSandbox(spec) {
    const state = { source: [...spec.sequence], copy: null, status: "" };
    const view = sortCloneView(spec);
    const refresh = () => renderSortClone(state, view);
    view.controls.append(
      button(spec.cloneSortLabel, () => { state.copy = [...state.source].sort((a, b) => a - b); state.status = spec.cloneSortText; refresh(); }),
      button(spec.originalSortLabel, () => { state.source.sort((a, b) => a - b); state.status = spec.originalSortText; refresh(); }),
      button(spec.resetLabel, () => { state.source = [...spec.sequence]; state.copy = null; state.status = spec.resetText; refresh(); }),
    );
    refresh();
    return view.panel;
  }

  function sortCloneView(spec) {
    const panel = element("section", "lesson-sandbox lesson-sort-clone");
    const source = element("figure", "lesson-sequence-source");
    const sourceRow = element("ol", "lesson-crystals");
    source.append(element("figcaption", "", spec.sourceLabel), sourceRow);
    const controls = element("div", "lesson-controls");
    const copyFigure = element("figure", "lesson-sort-copy");
    const copyRow = element("ol", "lesson-crystals");
    copyFigure.append(element("figcaption", "", T("النسخة المرتبة", "Sorted copy")), copyRow);
    const status = element("p", "lesson-caption");
    status.setAttribute("role", "status");
    panel.append(source, controls, copyFigure, status);
    return { panel, sourceRow, copyFigure, copyRow, controls, status };
  }

  function renderSortClone(state, view) {
    view.sourceRow.replaceChildren(...state.source.map(sortCloneSourceCell));
    view.copyRow.replaceChildren(...(state.copy || []).map(sortCloneCopyCell));
    view.copyFigure.hidden = state.copy === null;
    view.status.textContent = state.status;
    view.status.hidden = !state.status;
  }

  function sortCloneSourceCell(size, index) {
    const cell = element("li", "lesson-crystal");
    cell.append(element("b", "", String(size)), element("span", "lesson-index", String(index)));
    cell.setAttribute("aria-label", T(`index ${index}، حجم ${size}`, `index ${index}, size ${size}`));
    return cell;
  }

  function sortCloneCopyCell(size) {
    const cell = element("li", "lesson-crystal");
    cell.appendChild(element("b", "", String(size)));
    return cell;
  }

  function mapCountSandbox(spec) {
    const state = { counts: Object.create(null), readResult: null };
    const view = mapCountView(spec);
    const refresh = () => renderMapCount(state, view, spec);
    spec.kinds.forEach((kind) => view.kinds.appendChild(button(kind, () => {
      state.counts[kind] = (state.counts[kind] || 0) + 1;
      refresh();
    })));
    view.controls.append(
      button(spec.missingReadLabel, () => { state.readResult = "0, false"; refresh(); }),
      button(spec.resetLabel, () => { state.counts = Object.create(null); state.readResult = null; refresh(); }),
    );
    refresh();
    return view.panel;
  }

  function mapCountView(spec) {
    const panel = element("section", "lesson-sandbox lesson-map-count");
    const kinds = element("div", "lesson-map-kinds");
    const summary = element("p", "lesson-position");
    const order = element("p", "lesson-model-label", spec.mapOrderLabel);
    const table = element("table", "lesson-map-table");
    const head = element("thead");
    const headerRow = element("tr");
    headerRow.append(element("th", "", T("النوع", "Kind")), element("th", "", T("العدد", "Count")));
    head.appendChild(headerRow);
    const body = element("tbody");
    const empty = element("p", "lesson-map-empty", T("مفيش مفاتيح لسه.", "No keys yet."));
    const result = element("p", "lesson-map-readout");
    result.setAttribute("role", "status");
    const controls = element("div", "lesson-controls");
    table.append(head, body);
    panel.append(element("p", "lesson-model-label", spec.sourceLabel), kinds, summary, table, empty, order, controls, result);
    return { panel, kinds, summary, body, empty, result, controls };
  }

  function renderMapCount(state, view, spec) {
    const keys = Object.keys(state.counts).sort();
    view.summary.textContent = `counts · len=${keys.length}`;
    view.body.replaceChildren(...keys.map((key) => mapCountRow(key, state.counts[key])));
    view.empty.hidden = keys.length > 0;
    view.result.hidden = state.readResult === null;
    view.result.textContent = state.readResult === null ? "" : `${spec.missingKey}: ${state.readResult}`;
  }

  function mapCountRow(key, count) {
    const row = element("tr");
    row.append(element("td", "", key), element("td", "", String(count)));
    return row;
  }

  function sequenceSource(spec) {
    const source = element("figure", "lesson-sequence-source");
    source.appendChild(element("figcaption", "", spec.sourceLabel));
    const row = element("ol", "lesson-crystals");
    spec.sequence.forEach((size, index) => {
      const cell = element("li", "lesson-crystal");
      cell.append(element("b", "", String(size)), element("span", "lesson-index", String(index)));
      cell.setAttribute("aria-label", T(`index ${index}، حجم ${size}`, `index ${index}, size ${size}`));
      row.appendChild(cell);
    });
    source.appendChild(row);
    return source;
  }

  function sequenceSandbox(spec) {
    const panel = element("section", "lesson-sandbox lesson-sequence");
    const source = sequenceSource(spec);
    const view = element("div"), counter = element("span", "lesson-position");
    const done = element("p", "lesson-caption", spec.doneText);
    done.setAttribute("role", "status");
    counter.setAttribute("role", "status");
    const values = [];
    const add = button(spec.appendLabel, () => { if (values.length < spec.sequence.length) values.push(spec.sequence[values.length]); update(); });
    const reset = button(spec.resetLabel, () => { values.length = 0; update(); });
    function update() {
      view.replaceChildren(sliceView({ name: "sizes", values, length: values.length, elementsOnly: true }));
      counter.textContent = `${values.length} / ${spec.sequence.length}`;
      add.disabled = values.length === spec.sequence.length;
      done.hidden = !add.disabled;
    }
    const controls = element("div", "lesson-controls");
    controls.append(add, counter, reset);
    panel.append(source, view, controls, done);
    update();
    return panel;
  }

  function blockView(block, session) {
    const wrapper = element(block.optional ? "details" : "section", "lesson-block");
    wrapper.dataset.block = block.id;
    if (block.optional) wrapper.appendChild(element("summary", "", T("تفاصيل أكتر: تقدير ticks", "More detail: tick estimate")));
    wrapper.appendChild(prose(block.text));
    if (block.note) {
      const detail = element("details", "lesson-note");
      detail.append(element("summary", "", T("تفاصيل أكتر", "More detail")), prose(block.note));
      wrapper.appendChild(detail);
    }
    if (block.kind === "walkthrough") wrapper.appendChild(walkthrough(block));
    if (block.kind === "question") {
      if (block.code) wrapper.append(codeContext(block), codeView(block.code));
      wrapper.appendChild(question(block, session));
    }
    if (block.kind === "sandbox") wrapper.appendChild(sandbox(block));
    return wrapper;
  }

  function returnToCode() {
    document.body.classList.remove("lesson-open");
    $("#btn-lesson").setAttribute("aria-expanded", "false");
    ed.focus();
  }

  function navigation(level, navigate) {
    const nav = element("nav", "lesson-nav");
    nav.setAttribute("aria-label", T("خطوات الدرس", "Lesson steps"));
    const select = element("select", "lesson-step-select");
    select.setAttribute("aria-label", T("اختار خطوة", "Choose a step"));
    level.steps.forEach((step, index) => {
      const option = element("option", "", step.title);
      option.value = String(index); select.appendChild(option);
    });
    const position = element("span", "lesson-position");
    position.setAttribute("role", "status");
    const previous = button(T("السابق", "Back"), () => navigate(Number(select.value) - 1));
    const next = button(T("التالي", "Next"), () => navigate(Number(select.value) + 1));
    select.onchange = () => navigate(Number(select.value));
    const controls = element("div", "lesson-controls");
    controls.append(previous, position, next);
    nav.append(select, controls);
    return { node: nav, update(index) {
      select.value = String(index);
      position.textContent = T(`خطوة ${index + 1} من ${level.steps.length}`, `Step ${index + 1} of ${level.steps.length}`);
      previous.disabled = index === 0; next.disabled = index === level.steps.length - 1;
    } };
  }

  function open(level) {
    navigationObserver?.disconnect();
    const lessonPanel = $("#lesson-panel");
    lessonPanel.style.removeProperty("--lesson-nav-offset");
    const article = $("#lesson");
    if (!level.steps?.length) { article.innerHTML = md(level.lesson); return; }
    const session = createSession(level, {
      review: S.review, saved: store.blob.lessons[level.id],
      save(record) { store.blob.lessons[level.id] = record; store.save(); },
    });
    const title = element("h2", "lesson-step-title");
    title.tabIndex = -1;
    const body = element("div");
    const nav = navigation(level, navigate);
    function render() {
      const index = session.index;
      nav.update(index);
      title.textContent = level.steps[index].title;
      body.replaceChildren(...level.steps[index].blocks.map((block) => blockView(block, session)));
    }
    function navigate(index) {
      session.go(index); render(); title.focus();
      $("#lesson-panel").scrollTop = 0;
    }
    article.replaceChildren(nav.node, title, body, button(T("ارجع للكود", "Return to code"), returnToCode));
    render();
    navigationObserver = trackNavigationOffset(nav.node);
  }

  function trackNavigationOffset(nav) {
    const panel = $("#lesson-panel");
    const update = () => panel.style.setProperty("--lesson-nav-offset", `${nav.offsetHeight + 16}px`);
    const observer = new ResizeObserver(update);
    observer.observe(nav);
    update();
    return observer;
  }

  return { open, createSession };
})();
