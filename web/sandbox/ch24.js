"use strict";

function ch24Text(localizedText, T) {
  return T(localizedText.ar, localizedText.en);
}

function updateCh24Frame(frame, frameView) {
  frameView.caption.textContent = ch24Text(frame.caption, frameView.T);
  frameView.rows.replaceChildren(...frame.rows.map((row) => frameView.kit.valueRow({
    label: ch24Text(row.label, frameView.T),
    displayValue: row.value,
    helpers: { element: frameView.element },
  }).node));
}

function renderCh24Frames(block, renderTools) {
  const { element, button, T, kit } = renderTools;
  const config = block.sandbox.config;
  const panel = element("section", "lesson-sandbox ch24-playground");
  panel.setAttribute("dir", T("rtl", "ltr"));
  const caption = element("p", "ch24-playground-caption");
  caption.setAttribute("aria-live", "polite");
  const rows = element("div", "ch24-playground-rows");
  const frameView = { caption, rows, element, T, kit };
  let index = 0;
  const advanceButton = button(ch24Text(config.action, T), () => {
    index = (index + 1) % config.frames.length;
    updateCh24Frame(config.frames[index], frameView);
  });
  panel.append(caption, rows, advanceButton);
  updateCh24Frame(config.frames[index], frameView);
  return panel;
}

Lessons.registerSandbox("ch24-channel-handoff", renderCh24Frames);
Lessons.registerSandbox("ch24-buffer-queue", renderCh24Frames);
Lessons.registerSandbox("ch24-typed-reply", renderCh24Frames);
Lessons.registerSandbox("ch24-drain-close", renderCh24Frames);
Lessons.registerSandbox("ch24-http-result", renderCh24Frames);
Lessons.registerSandbox("ch24-write-offset", renderCh24Frames);
Lessons.registerSandbox("ch24-image-config", renderCh24Frames);
