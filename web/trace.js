"use strict";

const TracePanel = (() => {
  let dialog = null;
  let request = null;
  let returnFocus = null;

  function text(ar, en) {
    return typeof T === "function" ? T(ar, en) : (globalThis.LANG === "en" ? en : ar);
  }

  function escapeHTML(htmlText) {
    return String(htmlText).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  }

  function codeMarkup(code, line) {
    return String(code || "").split("\n").map((sourceLine, index) => {
      const number = index + 1;
      const active = number === line;
      return `<span class="trace-line${active ? " is-current" : ""}" data-trace-line="${number}"><span class="trace-line-number">${number}</span><code>${escapeHTML(sourceLine) || " "}</code></span>`;
    }).join("");
  }

  function createStepper(frames, onChange) {
    let index = 0;
    function set(next) {
      index = frames.length ? Math.max(0, Math.min(frames.length - 1, Number(next) || 0)) : 0;
      onChange(frames[index] || null, index);
      return frames[index] || null;
    }
    return {
      get index() { return index; },
      get frame() { return frames[index] || null; },
      previous() { return set(index - 1); },
      next() { return set(index + 1); },
      set,
    };
  }

  function panelMarkup() {
    return `
      <section class="trace-panel">
        <header class="trace-header">
          <div><h2 id="trace-title">${text("تتبّع كودك", "Trace your code")}</h2><span class="trace-file"></span></div>
          <button type="button" class="tool small trace-close" aria-label="${text("إغلاق", "Close")}">×</button>
        </header>
        <p class="trace-status" role="status" aria-live="polite"></p>
        <p class="trace-before">${text("القيم المعروضة قبل تنفيذ السطر المظلّل.", "Values are shown before the highlighted line runs.")}</p>
        <p class="trace-order hidden"></p>
        <p class="trace-truncated hidden" role="status"></p>
        <pre class="trace-code" dir="ltr" aria-label="${text("الكود مع السطر الحالي", "Code with current line")}"><code></code></pre>
        <section class="trace-variables" aria-live="polite">
          <h3>${text("قيم المتغيرات", "Variable values")}</h3>
          <div class="trace-variable-list"></div>
        </section>
        <footer class="trace-controls">
          <button type="button" class="tool small trace-prev">${text("السابق", "Previous")}</button>
          <input class="trace-slider" type="range" min="0" max="0" value="0" aria-label="${text("خطوة التتبّع", "Trace step")}">
          <button type="button" class="tool small trace-next">${text("التالي", "Next")}</button>
          <span class="trace-counter" aria-live="polite"></span>
        </footer>
      </section>`;
  }

  function ensureDialog() {
    if (dialog) return;
    dialog = document.createElement("dialog");
    dialog.className = "overlay trace-overlay";
    dialog.setAttribute("aria-labelledby", "trace-title");
    dialog.innerHTML = panelMarkup();
    document.body.appendChild(dialog);
    bindDialog();
  }

  function bindDialog() {
    dialog.querySelector(".trace-close").onclick = () => dialog.close();
    dialog.querySelector(".trace-prev").onclick = () => stepper?.previous();
    dialog.querySelector(".trace-next").onclick = () => stepper?.next();
    dialog.querySelector(".trace-slider").oninput = (event) => stepper?.set(event.target.value);
    dialog.addEventListener("click", (event) => { if (event.target === dialog) dialog.close(); });
    dialog.addEventListener("close", () => {
      request?.abort();
      request = null;
      if (returnFocus?.isConnected) returnFocus.focus();
    });
  }

  let stepper = null;
  let currentSources = {};
  let currentResponse = null;

  function renderFrame(frame, index) {
    if (!dialog || !frame || !currentResponse) return;
    const file = frame.file || "main.go";
    const source = currentSources[file] || currentSources["main.go"] || "";
    renderSource(file, source, frame.line);
    renderVariables(frame.variables || []);
    renderStepControls(index);
    renderTruncation(currentResponse.truncated || frame.truncated);
  }

  function renderSource(file, source, line) {
    dialog.querySelector(".trace-file").textContent = file;
    const code = dialog.querySelector(".trace-code");
    code.querySelector("code").innerHTML = codeMarkup(source, line);
    code.setAttribute("aria-label", text(`السطر الحالي ${line}`, `Current line ${line}`));
    code.querySelector(".is-current")?.scrollIntoView({ block: "center" });
  }

  function renderVariables(variables) {
    const list = dialog.querySelector(".trace-variable-list");
    list.innerHTML = variables.length
      ? `<table><thead><tr><th>${text("المتغير", "Variable")}</th><th>${text("القيمة", "Value")}</th></tr></thead><tbody>${variables.map((localVariable) => `<tr><th dir="ltr">${escapeHTML(localVariable.name)}</th><td dir="ltr">${escapeHTML(localVariable.value)}</td></tr>`).join("")}</tbody></table>`
      : `<p>${text("مفيش متغيرات محلية في الخطوة دي.", "No local variables in this step.")}</p>`;
  }

  function renderStepControls(index) {
    const total = currentResponse.frames.length;
    dialog.querySelector(".trace-counter").textContent = `${index + 1} / ${total}`;
    const slider = dialog.querySelector(".trace-slider");
    slider.value = String(index);
    slider.disabled = total < 2;
    dialog.querySelector(".trace-prev").disabled = index === 0;
    dialog.querySelector(".trace-next").disabled = index >= total - 1;
  }

  function renderTruncation(truncated) {
    const notice = dialog.querySelector(".trace-truncated");
    notice.textContent = text("بعض القيم أو الخطوات اتقصّت عشان نحافظ على حدّ التتبّع.", "Some values or steps were shortened to keep the trace within its limits.");
    notice.classList.toggle("hidden", !truncated);
  }

  function showResponse(response, sources) {
    currentResponse = response;
    currentSources = sources || {};
    const frames = response.frames || [];
    showTraceStatus(response);
    showGoroutineOrder(frames);
    dialog.querySelector(".trace-slider").max = String(Math.max(0, frames.length - 1));
    if (!response.available || !frames.length) {
      showNoFrames();
      return;
    }
    stepper = createStepper(frames, renderFrame);
    stepper.set(0);
  }

  function showTraceStatus(response) {
    dialog.querySelector(".trace-status").textContent = response.message || (response.completed === false
      ? text("البرنامج وقف قبل ما يكمّل؛ دي السطور اللي اتسجلت.", "The program stopped early; these are the lines captured.")
      : "");
  }

  function showGoroutineOrder(frames) {
    const order = dialog.querySelector(".trace-order");
    const goroutines = new Set(frames.map((frame) => frame.goroutine));
    order.textContent = text("ترتيب goroutines المعروض واحد من الترتيبات الممكنة.", "The goroutine order shown is one possible order.");
    order.classList.toggle("hidden", goroutines.size < 2);
  }

  function showNoFrames() {
    dialog.querySelector(".trace-file").textContent = "";
    dialog.querySelector(".trace-code code").textContent = "";
    dialog.querySelector(".trace-variable-list").textContent = "";
    dialog.querySelector(".trace-counter").textContent = "0 / 0";
    dialog.querySelector(".trace-prev").disabled = true;
    dialog.querySelector(".trace-next").disabled = true;
    dialog.querySelector(".trace-slider").disabled = true;
    dialog.querySelector(".trace-truncated").classList.add("hidden");
    if (!dialog.querySelector(".trace-status").textContent) {
      dialog.querySelector(".trace-status").textContent = text("مفيش خطوات تتبّع متاحة.", "No trace steps are available.");
    }
  }

  async function open(snapshot) {
    ensureDialog();
    request?.abort();
    returnFocus = document.activeElement;
    if (!dialog.open) dialog.showModal();
    resetPanelForLoading();
    request = new AbortController();
    return loadTrace(snapshot, request);
  }

  function resetPanelForLoading() {
    currentResponse = null;
    dialog.querySelector(".trace-status").textContent = text("بنشغّل نسخة التتبّع…", "Running a separate trace copy…");
    dialog.querySelector(".trace-order").classList.add("hidden");
    dialog.querySelector(".trace-truncated").classList.add("hidden");
    dialog.querySelector(".trace-code code").textContent = "";
    dialog.querySelector(".trace-variable-list").textContent = "";
    dialog.querySelector(".trace-counter").textContent = "";
    dialog.querySelector(".trace-prev").disabled = true;
    dialog.querySelector(".trace-next").disabled = true;
    dialog.querySelector(".trace-slider").disabled = true;
    stepper = null;
  }

  async function loadTrace(snapshot, activeRequest) {
    let response;
    try {
      response = await api("/api/trace", snapshot.payload, null, activeRequest.signal);
    } catch (error) {
      if (!activeRequest.signal.aborted && dialog.open) {
        dialog.querySelector(".trace-status").textContent = text("حصلت مشكلة في طلب التتبّع: ", "Trace request failed: ") + error.message;
      }
      return;
    }
    if (dialog.open && request === activeRequest) showResponse(response, snapshot.sources);
  }

  return { open, createStepper, codeMarkup };
})();

globalThis.TracePanel = TracePanel;
if (typeof module !== "undefined") module.exports = TracePanel;
