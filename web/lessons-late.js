"use strict";

Lessons.registerSandbox("sort-compare", (block, { element, button, T }) => {
  const spec = block.sandbox;
  const panel = element("section", "lesson-sandbox lesson-sort-compare");
  const source = element("figure", "lesson-sort-source");
  source.appendChild(element("figcaption", "", spec.sourceLabel));
  const row = element("ol", "lesson-sort-row");
  row.setAttribute("aria-label", T("كروت الأحجام والمواقع", "Size and position cards"));
  const startingPositions = spec.sequence.map((_, position) => position).reverse();
  const cards = new Map(spec.sequence.map((size, position) => {
    const card = element("li", "lesson-sort-card");
    card.dataset.position = String(position);
    card.setAttribute("aria-label", T(`الموقع ${position}، الحجم ${size}`, `Position ${position}, size ${size}`));
    card.append(element("b", "lesson-sort-size", String(size)), element("span", "lesson-sort-position", `x=${position}`));
    return [position, card];
  }));
  startingPositions.forEach((position) => row.appendChild(cards.get(position)));

  let descending = false;
  let tieBreak = false;
  const direction = button("", () => { descending = !descending; update(true); });
  const ties = button("", () => { tieBreak = !tieBreak; update(true); });
  const reset = button(spec.resetLabel, () => { descending = false; tieBreak = false; update(false, true); });
  const controls = element("div", "lesson-controls lesson-sort-controls");
  controls.append(direction, ties, reset);
  source.appendChild(row);
  panel.append(source, controls);

  function update(animate = false, restoreInput = false) {
    const previous = [...row.children];
    const bounds = new Map();
    if (animate && typeof globalThis.matchMedia === "function" && !globalThis.matchMedia("(prefers-reduced-motion: reduce)").matches) {
      for (const card of previous) {
        if (typeof card.getBoundingClientRect === "function") bounds.set(card, card.getBoundingClientRect());
      }
    }
    const current = restoreInput ? startingPositions.map((position) => cards.get(position)) : previous;
    const ordered = current.sort((left, right) => {
      const a = Number(left.dataset.position), b = Number(right.dataset.position);
      if (spec.sequence[a] !== spec.sequence[b]) {
        if (descending) return spec.sequence[a] < spec.sequence[b] ? 1 : -1;
        return spec.sequence[a] < spec.sequence[b] ? -1 : 1;
      }
      if (tieBreak) return a - b;
      return 0;
    });
    row.replaceChildren(...ordered);
    const tiedSizes = new Set(spec.sequence.filter((size) => spec.sequence.filter((other) => other === size).length > 1));
    for (const card of ordered) card.classList.toggle("same-size", tiedSizes.has(spec.sequence[Number(card.dataset.position)]));
    if (bounds.size) {
      for (const card of ordered) {
        const before = bounds.get(card);
        if (!before || typeof card.animate !== "function") continue;
        const after = card.getBoundingClientRect();
        const deltaX = before.left - after.left;
        const deltaY = before.top - after.top;
        if (deltaX || deltaY) card.animate(
          [{ transform: `translate(${deltaX}px, ${deltaY}px)` }, { transform: "translate(0, 0)" }],
          { duration: 200, easing: "ease-out" },
        );
      }
    }
    direction.textContent = `${spec.sortDirectionLabel}: ${descending ? spec.descendingLabel : spec.ascendingLabel}`;
    ties.textContent = `${spec.sortTieLabel}: ${tieBreak ? spec.tieOnLabel : spec.tieOffLabel}`;
    direction.setAttribute("aria-pressed", String(descending));
    ties.setAttribute("aria-pressed", String(tieBreak));
  }

  update(false);
  return panel;
});

Lessons.registerSandbox("grid-index", (block, { element, button, T }) => {
  const spec = block.sandbox;
  const panel = element("section", "lesson-sandbox lesson-grid-index");
  const source = element("figure", "lesson-grid-source");
  source.appendChild(element("figcaption", "", spec.sourceLabel));
  const expression = element("p", "lesson-grid-expression");
  expression.setAttribute("role", "status");
  expression.setAttribute("dir", "ltr");
  const views = element("div", "lesson-grid-views");
  const originalView = element("section", "lesson-grid-view");
  const copyView = element("section", "lesson-grid-view");
  originalView.appendChild(element("h4", "", T("الأصل", "Original")));
  copyView.appendChild(element("h4", "", T("النسخة", "Copy")));
  const originalTable = element("div", "lesson-grid-table-wrap");
  const copyTable = element("div", "lesson-grid-table-wrap");
  originalView.appendChild(originalTable);
  copyView.appendChild(copyTable);
  views.append(originalView, copyView);

  let mode = "outer";
  let selected = { row: 0, column: 0 };
  let original = [];
  let copy = [];
  const outerClone = button(spec.gridOuterCloneLabel, () => setMode("outer"));
  const rowClone = button(spec.gridRowCloneLabel, () => setMode("rows"));
  const edit = button(spec.gridEditCopyLabel, () => { copy[selected.row][selected.column] = 99; update(); });
  const reset = button(spec.resetLabel, resetState);
  const modes = element("div", "lesson-presets lesson-grid-modes");
  modes.append(outerClone, rowClone);
  const controls = element("div", "lesson-controls lesson-grid-controls");
  controls.append(edit, reset);

  function cloneRows(grid) { return grid.map((row) => [...row]); }
  function createCopy() {
    original = cloneRows(spec.grid);
    copy = mode === "outer" ? original.slice() : cloneRows(original);
  }
  function setMode(next) {
    mode = next;
    createCopy();
    update();
  }
  function resetState() {
    mode = "outer";
    selected = { row: 0, column: 0 };
    createCopy();
    update();
  }
  function gridTable(values, interactive, name) {
    const table = element("table", "lesson-grid-table");
    table.setAttribute("aria-label", name);
    const head = element("thead");
    const heading = element("tr");
    heading.appendChild(element("th", "lesson-grid-corner", "r/c"));
    values[0].forEach((_, column) => {
      const cell = element("th", "lesson-grid-column", String(column));
      cell.setAttribute("scope", "col");
      if (column === selected.column) cell.classList.add("selected-axis");
      heading.appendChild(cell);
    });
    head.appendChild(heading);
    table.appendChild(head);
    const body = element("tbody");
    values.forEach((row, rowIndex) => {
      const tr = element("tr");
      const rowHead = element("th", "lesson-grid-row", String(rowIndex));
      rowHead.setAttribute("scope", "row");
      if (rowIndex === selected.row) rowHead.classList.add("selected-axis");
      tr.appendChild(rowHead);
      row.forEach((value, column) => {
        const td = element("td", "lesson-grid-cell");
        if (rowIndex === selected.row && column === selected.column) td.classList.add("selected-cell");
        if (interactive) {
          const cell = button(String(value), () => { selected = { row: rowIndex, column }; update(); });
          cell.classList.add("lesson-grid-cell-button");
          cell.setAttribute("aria-label", T(`grid صف ${rowIndex} عمود ${column}: ${value}`, `grid row ${rowIndex} column ${column}: ${value}`));
          cell.setAttribute("aria-pressed", String(rowIndex === selected.row && column === selected.column));
          td.appendChild(cell);
        } else {
          td.appendChild(element("span", "lesson-grid-value", String(value)));
        }
        tr.appendChild(td);
      });
      body.appendChild(tr);
    });
    table.appendChild(body);
    return table;
  }
  function update() {
    outerClone.setAttribute("aria-pressed", String(mode === "outer"));
    rowClone.setAttribute("aria-pressed", String(mode === "rows"));
    expression.textContent = `grid[${selected.row}][${selected.column}] = ${original[selected.row][selected.column]}`;
    originalTable.replaceChildren(gridTable(original, true, T("الشبكة الأصلية", "Original grid")));
    copyTable.replaceChildren(gridTable(copy, false, T("الشبكة المنسوخة", "Cloned grid")));
  }
  createCopy();
  update();
  source.appendChild(expression);
  source.appendChild(views);
  panel.append(source, modes, controls);
  return panel;
});
