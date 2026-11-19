"use strict";

function ch23Text(localizedText, T) {
  return T(localizedText.ar, localizedText.en);
}

function updateCh23Frame(frame, frameView) {
  frameView.caption.textContent = ch23Text(frame.caption, frameView.T);
  frameView.rows.replaceChildren(...frame.rows.map((row) => frameView.kit.valueRow({
    label: ch23Text(row.label, frameView.T),
    displayValue: row.value,
    helpers: { element: frameView.element },
  }).node));
}

function renderCh23Frames(block, renderTools) {
  const { element, button, T, kit } = renderTools;
  const config = block.sandbox.config;
  const panel = element("section", "lesson-sandbox ch23-playground");
  panel.setAttribute("dir", T("rtl", "ltr"));
  const caption = element("p", "ch23-playground-caption");
  caption.setAttribute("aria-live", "polite");
  const rows = element("div", "ch23-playground-rows");
  const frameView = { caption, rows, element, T, kit };
  let index = 0;
  const advanceButton = button(ch23Text(config.action, T), () => {
    index = (index + 1) % config.frames.length;
    updateCh23Frame(config.frames[index], frameView);
  });
  panel.append(caption, rows, advanceButton);
  updateCh23Frame(config.frames[index], frameView);
  return panel;
}

Lessons.registerSandbox("ch23-read-full", renderCh23Frames);
Lessons.registerSandbox("ch23-aad-check", renderCh23Frames);
Lessons.registerSandbox("ch23-read-limit", renderCh23Frames);
Lessons.registerSandbox("ch23-stable-order", renderCh23Frames);
Lessons.registerSandbox("ch23-csv-records", renderCh23Frames);
