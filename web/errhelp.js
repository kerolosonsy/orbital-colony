"use strict";

window.ErrorHelp = {
  render(tip) {
    if (!tip?.meaning) return null;
    const details = document.createElement("details");
    details.className = "err-help";
    details.dir = "auto";

    const summary = document.createElement("summary");
    summary.textContent = tip.meaning;
    details.appendChild(summary);

    const body = document.createElement("div");
    body.className = "err-help-body";
    for (const [label, value] of [
      [T("السبب غالباً", "Likely cause"), tip.cause],
      [T("راجع", "Check"), tip.check],
    ]) {
      if (!value) continue;
      const row = document.createElement("p");
      const heading = document.createElement("strong");
      heading.textContent = label + ": ";
      const text = document.createElement("span");
      text.textContent = value;
      row.appendChild(heading);
      row.appendChild(text);
      body.appendChild(row);
    }
    details.appendChild(body);
    return details;
  },
};
