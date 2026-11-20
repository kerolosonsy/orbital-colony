"use strict";

const InlineDebugger = (() => {
  const esc = (htmlText) => String(htmlText).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/\"/g, "&quot;");
  const breakpointKey = (file, line) => `${file}:${line}`;

  function translate(ar, en) {
    return typeof document !== "undefined" && document.documentElement.lang === "en" ? en : ar;
  }

  function createSession() {
    let response = null;
    let active = false;
    return {
      begin(traceResponse) {
        response = traceResponse || null;
        active = !!response?.available && !!response.frames?.length;
      },
      edit() {
        response = null;
        active = false;
      },
      clear() { this.edit(); },
      get active() { return active; },
      get response() { return response; },
    };
  }

  function toggleBreakpointList(breakpoints, file, line) {
    const next = new Set(breakpoints || []);
    const key = breakpointKey(file, Number(line));
    if (next.has(key)) next.delete(key);
    else next.add(key);
    return [...next];
  }

  function frameGoroutine(frame) {
    return frame?.goroutine ?? 1;
  }

  function createStepper(frames, onChange = () => {}) {
    const list = Array.isArray(frames) ? frames : [];
    let index = 0;
    const set = (next) => {
      index = list.length ? Math.max(0, Math.min(list.length - 1, Number(next) || 0)) : 0;
      onChange(list[index] || null, index);
      return list[index] || null;
    };
    const advance = (predicate) => findNextFrame(list, index, set, predicate);
    return createStepperControls(list, set, advance, () => index);
  }

  function findNextFrame(frames, index, set, predicate) {
    if (!frames[index]) return null;
    for (let next = index + 1; next < frames.length; next++) if (predicate(frames[next], frames[index])) return set(next);
    return set(frames.length - 1);
  }

  function createStepperControls(frames, set, advance, currentIndex) {
    return {
      get index() { return currentIndex(); },
      get frame() { return frames[currentIndex()] || null; },
      set,
      restart() { return set(0); },
      back() { return set(currentIndex() - 1); },
      into() { return set(currentIndex() + 1); },
      over() {
        return advance((frame, from) => sameGoroutine(frame, from) && Number(frame.depth || 0) <= Number(from.depth || 0));
      },
      out() {
        return advance((frame, from) => sameGoroutine(frame, from) && Number(frame.depth || 0) < Number(from.depth || 0));
      },
      continueTo(breakpoints) { return continueAtBreakpoint(frames, currentIndex(), set, breakpoints); },
    };
  }

  function continueAtBreakpoint(frames, index, set, breakpoints) {
    const keys = breakpoints instanceof Set ? breakpoints : new Set(breakpoints || []);
    for (let next = index + 1; next < frames.length; next++) {
      if (keys.has(breakpointKey(frames[next].file || "main.go", frames[next].line))) return set(next);
    }
    return set(frames.length - 1);
  }

  function sameGoroutine(left, right) {
    return frameGoroutine(left) === frameGoroutine(right);
  }

  function inlineValueText(variables, limit = 84) {
    const text = (variables || []).slice(0, 5).map((variable) => `${variable.name} = ${variable.value}`).join(", ");
    if (text.length <= limit) return text;
    return text.slice(0, Math.max(0, limit - 1)) + "…";
  }

  function renderInlineValues(variables, changed = new Set(), compact = false) {
    const list = (variables || []).slice(0, compact ? 2 : 5);
    if (!list.length) return "";
    let budget = compact ? 34 : 72;
    const chips = [];
    for (const variable of list) {
      let value = String(variable.value);
      const room = budget - variable.name.length - 3;
      if (room <= 1) break;
      if (value.length > room) value = value.slice(0, Math.max(0, room - 1)) + "…";
      budget -= variable.name.length + value.length + 3;
      const cls = changed.has(variable.name) ? "debug-chip is-changed" : "debug-chip";
      chips.push(`<span class="${cls}"><b>${esc(variable.name)}</b> ${esc(value)}</span>`);
    }
    return `<span class="debug-inline-values" aria-hidden="true">${chips.join("")}</span>`;
  }

  function editorMarkup(sourceCode, debugState, highlightCode) {
    if (!debugState?.active) return highlightCode(sourceCode) + "\n ";
    return splitHighlightedLines(highlightCode(sourceCode)).map((part, index) => {
      const number = index + 1;
      const active = number === debugState.line;
      const variables = debugState.inlineValuesByLine?.[number] || [];
      const openTags = part.open.join("");
      const closeTags = "</span>".repeat(part.close.length);
      const extra = active ? "" : number === debugState.previousLine ? " is-trail" : "";
      const changed = active ? new Set(debugState.changed || []) : new Set();
      return `<span class="debug-line${active ? " is-current" : ""}${extra}" data-debug-line="${number}">${openTags}${part.content}${closeTags}${renderInlineValues(variables, changed, !active)}</span>`;
    }).join("");
  }

  function splitHighlightedLines(markup) {
    const lines = [];
    const state = { open: [], active: [], content: "" };
    let index = 0;
    while (index < markup.length) {
      if (markup[index] === "\n") {
        lines.push({ content: state.content, open: state.open, close: state.active.slice() });
        state.content = "";
        state.open = state.active.slice();
        index++;
      } else {
        index = appendHighlightedMarkup(markup, index, state);
      }
    }
    lines.push({ content: state.content, open: state.open, close: state.active });
    return lines;
  }

  function appendHighlightedMarkup(markup, index, state) {
    if (markup.startsWith("<span ", index)) {
      const end = markup.indexOf(">", index) + 1;
      const tag = markup.slice(index, end);
      state.content += tag;
      state.active.push(tag);
      return end;
    }
    if (markup.startsWith("</span>", index)) {
      state.content += "</span>";
      state.active.pop();
      return index + 7;
    }
    state.content += markup[index];
    return index + 1;
  }

  function gutterMarkup(lineCount, errLines, breakpoints, currentLine, visited = new Set()) {
    const breakSet = breakpoints instanceof Set ? breakpoints : new Set(breakpoints || []);
    let markup = "";
    for (let line = 1; line <= lineCount; line++) {
      const error = errLines.get(line);
      const isBreakpoint = breakSet.has(line);
      const classes = ["gutter-line", error ? "err" : "", isBreakpoint ? "has-breakpoint" : "", line === currentLine ? "debug-active" : "", visited.has(line) && line !== currentLine ? "debug-visited" : ""].filter(Boolean).join(" ");
      const label = translate(`السطر ${line}${isBreakpoint ? "، نقطة توقف" : ""}${error ? "، فيه خطأ" : ""}`, `Line ${line}${isBreakpoint ? ", breakpoint" : ""}${error ? ", error" : ""}`);
      const tip = error || label + translate(isBreakpoint ? " — دوس تشيلها (F9)" : " — دوس تحط نقطة توقف (F9)", isBreakpoint ? " — click to remove (F9)" : " — click for a breakpoint (F9)");
      markup += `<button type="button" class="${classes}" data-debug-line="${line}" aria-label="${esc(label)}" title="${esc(tip)}"><span>${line}</span></button>`;
    }
    return markup;
  }

  function callStackAt(frames, index, goroutine = frameGoroutine(frames?.[index])) {
    const stack = [];
    const end = Math.min(Number(index), (frames || []).length - 1);
    for (let i = 0; i <= end; i++) {
      const frame = frames[i];
      if (frameGoroutine(frame) !== goroutine) continue;
      if (frame.event === "call") stack.push({ func: frame.func || "?", file: frame.file || "main.go", line: frame.line, depth: frame.depth || 0, goroutine, index: i });
      if (frame.event === "return" && i < end) {
        const at = stack.map((entry) => entry.func).lastIndexOf(frame.func || "?");
        if (at >= 0) stack.splice(at, 1);
      }
    }
    return stack;
  }

  function stdoutAtOffset(stdout, offset, cachedBytes) {
    const bytes = cachedBytes || new TextEncoder().encode(String(stdout || ""));
    const count = Math.max(0, Math.min(bytes.length, Number(offset) || 0));
    return new TextDecoder().decode(bytes.slice(0, count));
  }

  let session = createSession();
  let elements = null;
  let callbacks = {};
  let stepper = null;
  let request = null;
  let levelID = "";
  let response = null;
  let responseStdoutBytes = null;
  let lastVariables = new Map();
  let collapsed = false;
  let playTimer = null;
  let hooks = {};
  let endFired = false;
  let speed = 1;
  let userCollapsed = false;
  let lastStdoutOffset = 0;

  const labels = {
    restart: ["إعادة من الأول", "Restart"], back: ["خطوة لورا", "Step back"], into: ["ادخل جوه الدالة", "Step into"],
    over: ["عدّي الدالة", "Step over"], out: ["اخرج من الدالة", "Step out"], go: ["كمّل لحد نقطة توقف", "Continue"],
    play: ["شغّل خطوة بخطوة", "Play"], pause: ["وقّف التشغيل", "Pause"], exit: ["اقفل التتبّع", "Stop debugging"],
    variables: ["المتغيرات", "Variables"], stack: ["Call stack", "Call stack"], output: ["Output", "Output"],
    collapse: ["اقفل اللوحة", "Collapse panel"], expand: ["افتح اللوحة", "Expand panel"],
  };

  function bind(options) {
    callbacks = options || {};
    elements = findElements();
    bindStepButtons();
    bindTabButtons();
    elements.slider.oninput = () => { stopPlay(); stepper?.set(elements.slider.value); };
    const speedSelect = document.querySelector("#debugger-speed");
    if (speedSelect) {
      speed = Number(speedSelect.value) || 1;
      speedSelect.onchange = () => {
        speed = Number(speedSelect.value) || 1;
        if (playTimer) { stopPlay(); startPlay(); }
      };
    }
    elements.collapse.onclick = () => { userCollapsed = true; setCollapsed(!collapsed); };
    document.addEventListener("keydown", onKeydown);
    setCollapsed(matchMedia("(max-width: 760px)").matches);
  }

  function findElements() {
    return {
      toolbar: document.querySelector("#debugger-toolbar"),
      controls: document.querySelector("#debugger-controls"),
      status: document.querySelector("#debugger-status"),
      counter: document.querySelector("#debugger-counter"),
      slider: document.querySelector("#debugger-slider"),
      panel: document.querySelector("#debugger-panel"),
      variables: document.querySelector("#debugger-variables"),
      stack: document.querySelector("#debugger-stack"),
      output: document.querySelector("#debugger-output"),
      order: document.querySelector("#debugger-order"),
      collapse: document.querySelector("#debugger-collapse"),
    };
  }

  function bindStepButtons() {
    const shortcuts = { into: "F11", over: "F10", out: "Shift+F11", go: "F5", exit: "Shift+F5" };
    document.querySelectorAll("[data-debug-command]").forEach((button) => {
      const command = button.dataset.debugCommand;
      const shortcut = shortcuts[command];
      const label = translate(...labels[command]);
      button.setAttribute("aria-label", label);
      if (shortcut) button.setAttribute("aria-keyshortcuts", shortcut);
      button.title = shortcut ? `${label} (${shortcut})` : label;
      button.onclick = () => commandStep(command);
    });
  }

  function bindTabButtons() {
    document.querySelectorAll("[data-debug-tab]").forEach((button) => {
      const tab = button.dataset.debugTab;
      button.textContent = translate(...labels[tab]);
      button.onclick = () => selectTab(tab);
    });
  }

  function setCollapsed(collapsedState) {
    collapsed = !!collapsedState;
    if (!elements) return;
    elements.panel.classList.toggle("is-collapsed", collapsed);
    const [ar, en] = labels[collapsed ? "expand" : "collapse"];
    elements.collapse.textContent = collapsed ? "⌃" : "⌄";
    elements.collapse.title = translate(ar, en);
    elements.collapse.setAttribute("aria-label", translate(ar, en));
    elements.collapse.setAttribute("aria-expanded", String(!collapsed));
  }

  async function start(levelIDForTrace, tracePayload) {
    clear();
    levelID = levelIDForTrace;
    showPending();
    request = new AbortController();
    const activeRequest = request;
    try {
      const traceResponse = await callbacks.requestTrace(tracePayload, activeRequest.signal);
      if (request === activeRequest) acceptTraceResponse(traceResponse);
    } catch (error) {
      if (request === activeRequest && !activeRequest.signal.aborted) showTraceFailure(error);
    }
  }

  // startWith opens a session on frames that arrived with the run itself.
  // hooks: onFrame(frame, index, total), onEnd() once when the last frame is
  // reached, onExit() when the learner closes the debugger, onClear() when
  // the session is dropped from outside (edit, new run, other level).
  function startWith(levelIDForTrace, traceResponse, options = {}) {
    clear();
    levelID = levelIDForTrace;
    hooks = options;
    endFired = false;
    acceptTraceResponse(traceResponse);
    if (options.autoplay && session.active) startPlay();
  }

  function finishNow() {
    if (!stepper || !session.active) return;
    stopPlay();
    stepper.set(response.frames.length - 1);
  }

  function acceptTraceResponse(traceResponse) {
    response = traceResponse || {};
    responseStdoutBytes = new TextEncoder().encode(String(response.stdout || ""));
    if (!response.available) return showUnavailableTrace();
    session.begin(response);
    showAvailableTrace();
  }

  function showUnavailableTrace() {
    session.edit();
    showReason(response.message || translate("التتبّع مش متاح للكود ده.", "Debugger is unavailable for this code."));
    callbacks.renderEditor?.();
  }

  function showAvailableTrace() {
    const frames = response.frames || [];
    setGoroutineNote(frames);
    stepper = createStepper(frames, renderFrame);
    elements.toolbar.hidden = false;
    elements.controls.hidden = false;
    elements.panel.hidden = false;
    setFocus(true);
    if (!userCollapsed) setCollapsed(matchMedia("(max-width: 760px)").matches);
    if (!frames.length) return showEmptyTrace();
    setStatus(traceStatus(response));
    stepper.restart();
  }

  function setGoroutineNote(frames) {
    const goroutines = new Set(frames.map((frame) => frameGoroutine(frame)));
    elements.order.hidden = goroutines.size < 2;
    elements.order.textContent = translate("ترتيب goroutines المعروض واحد من الترتيبات الممكنة.", "The goroutine order shown is one possible order.");
  }

  function traceStatus(traceResponse) {
    if (traceResponse.message) return traceResponse.message;
    const notices = [];
    if (traceResponse.completed === false) notices.push(translate("البرنامج وقف بدري؛ الخطوات اللي اتسجلت لسه متاحة.", "The program stopped early; captured steps are available."));
    if (traceResponse.truncated) notices.push(translate("بعض الخطوات أو القيم اتقصّت عشان نحافظ على الحد.", "Some steps or values were shortened to stay within the limit."));
    return notices.join(" ");
  }

  function showEmptyTrace() {
    showReason(response.message || translate("مفيش خطوات متاحة.", "No steps are available."));
    elements.panel.hidden = false;
  }

  function showTraceFailure(error) {
    session.edit();
    stepper = null;
    response = null;
    showReason(translate("طلب التتبّع فشل: ", "Debugger request failed: ") + error.message);
    callbacks.renderEditor?.();
  }

  function showPending() {
    elements.toolbar.hidden = false;
    elements.controls.hidden = true;
    elements.panel.hidden = true;
    setStatus(translate("بنشغّل نسخة التتبّع…", "Preparing the debugger…"));
  }

  function showReason(message) {
    elements.toolbar.hidden = false;
    elements.controls.hidden = true;
    elements.panel.hidden = true;
    setStatus(message);
  }

  function unavailable(message) {
    clear();
    showReason(message);
  }

  function setStatus(message) {
    if (elements.status) elements.status.textContent = message || "";
  }

  function clear() {
    const dropped = hooks;
    hooks = {};
    endFired = false;
    dropped.onClear?.();
    resetSession();
    if (!elements) return;
    hideDebuggerUI();
    callbacks.renderEditor?.();
  }

  function resetSession() {
    stopPlay();
    lastStdoutOffset = 0;
    request?.abort();
    request = null;
    session.clear();
    response = null;
    responseStdoutBytes = null;
    stepper = null;
    lastVariables = new Map();
  }

  // Focus mode gives the code the lesson's column while a session is open.
  function setFocus(on) {
    if (typeof document === "undefined" || !document.body) return;
    document.body.classList.toggle("debug-focus", on);
    if (!on) document.body.classList.remove("lesson-open");
    callbacks.onFocusChange?.(on);
  }

  function hideDebuggerUI() {
    setFocus(false);
    elements.toolbar.hidden = true;
    elements.panel.hidden = true;
    elements.controls.hidden = true;
    elements.status.textContent = "";
    elements.variables.textContent = "";
    elements.stack.textContent = "";
    elements.output.textContent = "";
    elements.order.hidden = true;
  }

  function commandStep(command) {
    if (!stepper || !session.active) {
      if (command === "exit") clear();
      return;
    }
    if (command === "play") return togglePlay();
    stopPlay();
    const breakpoints = callbacks.getBreakpoints?.(levelID) || [];
    if (command === "restart") stepper.restart();
    else if (command === "back") stepper.back();
    else if (command === "into") stepper.into();
    else if (command === "over") stepper.over();
    else if (command === "out") stepper.out();
    else if (command === "go") stepper.continueTo(breakpoints);
    else if (command === "exit") {
      hooks.onExit?.();
      clear();
    }
  }

  function onKeydown(event) {
    if (!session.active) return;
    let command = null;
    if (event.key === "F5" && event.shiftKey) command = "exit";
    else if (event.key === "F5") command = "go";
    else if (event.key === "F10") command = "over";
    else if (event.key === "F11" && event.shiftKey) command = "out";
    else if (event.key === "F11") command = "into";
    // keys that work without Fn on laptops: Alt+← step back, Alt+→ step, Alt+Space play/pause
    else if (event.altKey && !event.metaKey && !event.ctrlKey && event.key === "ArrowLeft") command = "back";
    else if (event.altKey && !event.metaKey && !event.ctrlKey && event.key === "ArrowRight") command = "into";
    else if (event.altKey && !event.metaKey && !event.ctrlKey && event.code === "Space") command = "play";
    if (command) {
      event.preventDefault();
      commandStep(command);
    }
  }

  function playInterval() {
    const total = response?.frames?.length || 1;
    // a typical mission plays in a few seconds: not slow, still readable
    const base = Math.max(60, Math.min(400, 5000 / total));
    return base / speed;
  }

  function togglePlay() {
    if (playTimer) return stopPlay();
    if (stepper.index >= response.frames.length - 1) stepper.restart();
    startPlay();
  }

  function startPlay() {
    if (!stepper) return;
    setPlaying(true);
    playTimer = setInterval(() => {
      const before = stepper.index;
      const frame = stepper.into();
      const keys = new Set(callbacks.getBreakpoints?.(levelID) || []);
      const atBreak = frame && keys.has(breakpointKey(frame.file || "main.go", frame.line));
      if (stepper.index === before || stepper.index >= response.frames.length - 1 || atBreak) stopPlay();
    }, playInterval());
  }

  function stopPlay() {
    if (playTimer) clearInterval(playTimer);
    playTimer = null;
    setPlaying(false);
  }

  function setPlaying(on) {
    if (!elements) return;
    elements.toolbar.classList.toggle("is-playing", on);
    const button = document.querySelector('[data-debug-command="play"]');
    if (button) {
      const label = translate(...labels[on ? "pause" : "play"]);
      button.setAttribute("aria-label", label);
      button.title = label;
      button.setAttribute("aria-pressed", String(on));
    }
  }

  function whereText(frame) {
    if (!frame) return "";
    const fn = `<code>${esc(frame.func || "main")}()</code>`;
    const line = translate(`سطر ${frame.line}`, `line ${frame.line}`);
    if (frame.event === "call") return translate(`دخلنا ${fn}`, `Entered ${fn}`) + ` · ${line}`;
    if (frame.event === "return") {
      const values = (frame.returns || []).map((r) => esc(r.value)).join(", ");
      return translate(`رجعنا من ${fn}`, `Returned from ${fn}`) + (values ? ` → <code>${values}</code>` : "") + ` · ${line}`;
    }
    return `${fn} · ${line}`;
  }

  function renderTimeline(index) {
    const svg = document.querySelector("#debugger-depth");
    if (!svg || !response) return;
    const frames = response.frames;
    const n = frames.length;
    const W = 1000, H = 26;
    const maxDepth = Math.max(1, ...frames.map((f) => Number(f.depth || 0)));
    const x = (i) => (n < 2 ? 0 : (i / (n - 1)) * W);
    const y = (d) => H - 3 - (Number(d || 0) / maxDepth) * (H - 8);
    if (svg.dataset.frames !== String(n)) {
      let d = `M0 ${H}`;
      frames.forEach((f, i) => { d += ` L${x(i).toFixed(1)} ${y(f.depth).toFixed(1)}`; });
      d += ` L${W} ${H} Z`;
      svg.setAttribute("viewBox", `0 0 ${W} ${H}`);
      svg.innerHTML = `<path class="dbg-depth-area" d="${d}"/><g class="dbg-ticks"></g><line class="dbg-cursor" x1="0" x2="0" y1="0" y2="${H}"/>`;
      svg.dataset.frames = String(n);
    }
    const keys = new Set(callbacks.getBreakpoints?.(levelID) || []);
    svg.querySelector(".dbg-ticks").innerHTML = frames.map((f, i) => keys.has(breakpointKey(f.file || "main.go", f.line)) ? `<rect x="${(x(i) - 1.5).toFixed(1)}" y="0" width="3" height="5"/>` : "").join("");
    const cursor = svg.querySelector(".dbg-cursor");
    cursor.setAttribute("x1", x(index).toFixed(1));
    cursor.setAttribute("x2", x(index).toFixed(1));
  }

  function renderFrame(frame, index) {
    if (!elements || !response) return;
    const total = response.frames.length;
    renderStepControls(index, total);
    const where = document.querySelector("#debugger-where-text");
    if (where) where.innerHTML = whereText(frame);
    renderTimeline(index);
    renderVariables(frame, index);
    renderStack(frame, index, response.frames);
    renderOutput(frame);
    flashChangedVariables(frame);
    showFrameFile(frame);
    hooks.onFrame?.(frame, index, total);
    if (index >= total - 1 && !endFired) {
      endFired = true;
      hooks.onEnd?.();
    }
  }

  function renderStepControls(index, total) {
    elements.counter.textContent = `${index + 1} / ${total}`;
    elements.slider.max = String(Math.max(0, total - 1));
    elements.slider.value = String(index);
    elements.slider.disabled = total < 2;
    for (const button of document.querySelectorAll("[data-debug-command]")) {
      const command = button.dataset.debugCommand;
      button.disabled = command === "back" && index === 0 || ["into", "over", "out", "go"].includes(command) && index >= total - 1;
    }
  }

  function renderOutput(frame) {
    const outputOffset = Number(frame?.stdout) || 0;
    const capturedBytes = responseStdoutBytes?.length || 0;
    const outputText = stdoutAtOffset(response.stdout, outputOffset, responseStdoutBytes);
    const oldText = stdoutAtOffset(response.stdout, Math.min(lastStdoutOffset, outputOffset), responseStdoutBytes);
    const fresh = outputText.slice(oldText.length);
    lastStdoutOffset = outputOffset;
    const tail = response.stdoutTruncated && outputOffset > capturedBytes ? translate("\n…(الناتج اتقص)", "\n…(output truncated)") : "";
    elements.output.innerHTML = esc(oldText) + (fresh ? `<mark class="debug-fresh">${esc(fresh)}</mark>` : "") + esc(tail);
    if (!outputText) elements.output.innerHTML = `<span class="debug-empty">${esc(translate("لسه مفيش حاجة اتطبعت.", "Nothing printed yet."))}</span>`;
  }

  function flashChangedVariables(frame) {
    const prior = new Map(lastVariables);
    lastVariables = new Map((frame?.variables || []).map((variable) => [variable.name, variable.value]));
    elements.variables.querySelectorAll("[data-variable-name]").forEach((row) => {
      const name = row.dataset.variableName;
      row.classList.toggle("changed", prior.has(name) && prior.get(name) !== lastVariables.get(name));
    });
  }

  function showFrameFile(frame) {
    const file = frame?.file || "main.go";
    if (callbacks.getFile?.() !== file) callbacks.showFile?.(file);
    else callbacks.renderEditor?.();
  }

  function valueKind(value) {
    const text = String(value);
    if (/^-?\d+(\.\d+)?(e[-+]?\d+)?$/i.test(text)) return "num";
    if (text === "true" || text === "false") return "bool";
    if (text === "<nil>" || text === "nil") return "nil";
    if (/^[\[{]|^map\[/.test(text)) return "agg";
    return "str";
  }

  // The closing "}" or a return frame has no locals of its own: keep showing
  // the last values this function had so the final state stays visible.
  function lastKnownVariables(frame, index) {
    for (let i = index - 1; i >= 0; i--) {
      const earlier = response.frames[i];
      if (frameGoroutine(earlier) === frameGoroutine(frame) && earlier.func === frame.func && earlier.variables?.length) return earlier.variables;
    }
    return [];
  }

  function renderVariables(frame, index = 0) {
    let variables = frame?.variables || [];
    const prior = lastVariables;
    const returns = frame?.returns || [];
    let note = "";
    if (!variables.length && response) {
      variables = lastKnownVariables(frame, index);
      if (variables.length) note = `<p class="debug-note">${translate("آخر قيم قبل ما الدالة تخلص:", "Last values before the function ended:")}</p>`;
    }
    if (!variables.length && !returns.length) {
      elements.variables.innerHTML = `<p>${translate("مفيش متغيرات محلية في الخطوة دي.", "No local variables in this step.")}</p>`;
      return;
    }
    const header = `<thead><tr><th>${translate("المتغير", "Variable")}</th><th>${translate("القيمة", "Value")}</th></tr></thead>`;
    elements.variables.innerHTML = note + (variables.length
      ? `<table>${header}<tbody>${variables.map((variable) => `<tr data-variable-name="${esc(variable.name)}"><th dir="ltr">${esc(variable.name)}</th><td dir="ltr"><span class="debug-val debug-val-${valueKind(variable.value)}">${esc(variable.value)}</span>${prior.has(variable.name) && prior.get(variable.name) !== variable.value ? `<s class="debug-was">${esc(prior.get(variable.name))}</s>` : ""}</td></tr>`).join("")}</tbody></table>`
      : "") + (returns.length
      ? `<h4>${translate("قيم الرجوع", "Return values")}</h4><table>${header}<tbody>${returns.map((variable) => `<tr><th dir="ltr">${esc(variable.name)}</th><td dir="ltr">${esc(variable.value)}</td></tr>`).join("")}</tbody></table>`
      : "");
  }

  function renderStack(frame, index, frames) {
    const entries = callStackAt(frames, index, frameGoroutine(frame));
    const goroutine = frameGoroutine(frame);
    elements.stack.innerHTML = `<p class="debug-goroutine">goroutine ${esc(goroutine)}</p>` + (entries.length
      ? entries.map((entry, i) => `<button type="button" class="debug-stack-frame${i === entries.length - 1 ? " is-top" : ""}" style="--depth:${i}" data-frame-index="${entry.index}"><span>${esc(entry.func)}()</span><small dir="ltr">${esc(entry.file)}:${i === entries.length - 1 ? frame.line : entry.line}</small></button>`).reverse().join("")
      : `<p>${translate("مفيش call stack متاح.", "No call stack is available.")}</p>`);
    elements.stack.querySelectorAll("[data-frame-index]").forEach((button) => {
      button.onclick = () => {
        const nextIndex = Number(button.dataset.frameIndex);
        stepper?.set(nextIndex);
      };
    });
  }

  function selectTab(tab) {
    for (const button of document.querySelectorAll("[data-debug-tab]")) {
      const selected = button.dataset.debugTab === tab;
      button.setAttribute("aria-selected", String(selected));
    }
    for (const name of ["variables", "stack", "output"]) {
      document.querySelector(`#debugger-${name}`).hidden = name !== tab;
    }
  }

  function editorState(file) {
    if (!session.active || !stepper?.frame || stepper.frame.file !== file) return { active: false };
    const current = stepper.frame;
    const map = new Map();
    let callStart = 0;
    for (let i = stepper.index; i >= 0; i--) {
      const frame = response.frames[i];
      if (frame.event === "call" && frame.func === current.func && Number(frame.depth || 0) === Number(current.depth || 0) && frameGoroutine(frame) === frameGoroutine(current)) { callStart = i; break; }
    }
    const frames = response.frames.slice(0, stepper.index + 1);
    // Other lines show only what they changed the last time they ran (the
    // values right after them), so the code stays readable; the current
    // line shows everything in scope.
    const own = frames.slice(callStart).filter((frame) => frame.file === file && frameGoroutine(frame) === frameGoroutine(current)
      && frame.func === current.func && Number(frame.depth || 0) === Number(current.depth || 0) && (frame.event === "line" || frame.event === "return"));
    for (let k = 0; k + 1 < own.length; k++) {
      const ran = own[k];
      if (ran.event !== "line") continue;
      const beforeLine = new Map((ran.variables || []).map((variable) => [variable.name, variable.value]));
      const effect = (own[k + 1].variables || []).filter((variable) => beforeLine.get(variable.name) !== variable.value);
      if (effect.length) map.set(ran.line, effect);
      else map.delete(ran.line);
    }
    if (current.event === "line") map.set(current.line, current.variables || []);
    const visited = new Set(frames.filter((frame) => frame.file === file && frame.event === "line").map((frame) => frame.line));
    let previousLine = 0;
    for (let i = stepper.index - 1; i >= 0; i--) {
      const frame = response.frames[i];
      if (frame.file === file && frame.event === "line" && frame.line !== current.line && frameGoroutine(frame) === frameGoroutine(current)) { previousLine = frame.line; break; }
    }
    const before = new Map();
    for (let i = stepper.index - 1; i >= 0; i--) {
      const frame = response.frames[i];
      if (frame.func !== current.func || frameGoroutine(frame) !== frameGoroutine(current)) continue;
      for (const variable of frame.variables || []) if (!before.has(variable.name)) before.set(variable.name, variable.value);
      if (before.size >= (current.variables || []).length) break;
    }
    const changed = (current.variables || []).filter((variable) => before.has(variable.name) && before.get(variable.name) !== variable.value).map((variable) => variable.name);
    return { active: true, line: current.line, previousLine, visited, changed, inlineValuesByLine: Object.fromEntries([...map].map(([line, values]) => [line, values])) };
  }

  function breakpointLines(level, file) {
    const stored = callbacks.getBreakpoints?.(level) || [];
    return new Set(stored.filter((key) => key.startsWith(`${file}:`)).map((key) => Number(key.slice(key.lastIndexOf(":") + 1))));
  }

  function toggleLevelBreakpoint(level, file, line) {
    const current = callbacks.getBreakpoints?.(level) || [];
    callbacks.saveBreakpoints?.(level, toggleBreakpointList(current, file, line));
    if (stepper && session.active) renderTimeline(stepper.index);
  }

  return {
    bind, start, startWith, finishNow, clear, unavailable, createSession, createStepper, toggleBreakpoint: toggleLevelBreakpoint, toggleBreakpointList,
    breakpointKey, breakpointLines, editorState, editorMarkup, gutterMarkup, inlineValueText, renderInlineValues,
    callStackAt, stdoutAtOffset, get active() { return session.active; }, get playing() { return !!playTimer; },
  };
})();

globalThis.InlineDebugger = InlineDebugger;
if (typeof module !== "undefined") module.exports = InlineDebugger;
