"use strict";

function fitView(bounds, width, height, margin) {
  const zoom = Math.min(1, Math.max(1, width - 2 * margin) / Math.max(1, bounds.maxX - bounds.minX), Math.max(1, height - 2 * margin) / Math.max(1, bounds.maxY - bounds.minY));
  return { zoom, x: width / 2 - (bounds.minX + bounds.maxX) / 2 * zoom, y: height / 2 - (bounds.minY + bounds.maxY) / 2 * zoom };
}

const TreeView = (() => {
  function edgeSatisfied(_node, prereq, byId) {
    const parent = byId[prereq.id];
    return !!parent && parent.level >= prereq.level;
  }

  function missingCount(node) { return node.checks.filter((check) => !check.ok).length; }

  function pathToNearestOwned(id, byId) {
    if (!byId[id] || byId[id].state === "owned") return [];
    const queue = [[id]], seen = new Set([id]);
    for (let i = 0; i < queue.length; i++) {
      const chain = queue[i], node = byId[chain[0]];
      if (node.state === "owned") return chain;
      for (const parent of [...node.after].sort((a, b) => a.id.localeCompare(b.id))) {
        if (byId[parent.id] && !seen.has(parent.id)) {
          seen.add(parent.id); queue.push([parent.id, ...chain]);
        }
      }
    }
    return [];
  }

  function matchesSearch(node, query) {
    const needle = query.trim().toLocaleLowerCase();
    return node.state !== "unknown" && !!needle && [node.name, node.nameEn].some((name) => (name || "").toLocaleLowerCase().includes(needle));
  }

  function buyEnabled(node) {
    return node.state === "available" || node.state === "owned" && node.checks.every((check) => check.ok);
  }

  function buyLabel(node) {
    if (node.checks.some((check) => check.kind === "max")) return T("مكتمل", "Maxed");
    return node.level === 0 ? T("ابني", "Build") : T("رقّي", "Upgrade");
  }

  function branchHue(branch, branches) {
    const ordered = [...new Set(branches.filter(Boolean))].sort();
    return 25 + Math.max(0, ordered.indexOf(branch)) % 9 * 40;
  }

  let tree = null, handlers = {}, overlay = null, svg = null, world = null, detail = null;
  let byId = {}, positions = {}, branches = [], selected = null, filter = "", query = "";
  let pan = { x: 0, y: 0 }, zoom = 1, minZoom = 0.25, width = 0, height = 0, returnFocus = null;
  const pointers = new Map();
  let moved = false, pressedNode = null;
  const local = (node, key) => T(node[key] || "", (LANG === "en" && node[key + "En"]) || node[key] || "");
  const fmt = (n) => Math.floor(n).toLocaleString("en-US");
  const attr = (s) => esc(String(s)).replace(/"/g, "&quot;").replace(/'/g, "&#39;");
  const color = (branch) => branch ? `var(--tv-branch-${(branchHue(branch, branches) - 25) / 40})` : "var(--paper-2)";

  function createOverlay() {
    overlay = document.createElement("dialog");
    overlay.className = "treeview";
    overlay.setAttribute("aria-labelledby", "tv-title");
    overlay.innerHTML = `<header class="tv-header"><h2 id="tv-title">${esc(T("شجرة التقنيات", "Tech tree"))}</h2>
      <div class="tv-stats"></div>
      <label class="tv-search"><span>${esc(T("بحث بالاسم", "Search by name"))}</span><input type="search" autocomplete="off"></label>
      <button type="button" data-action="centre">${esc(T("رجّعني للمركز", "Re-centre"))}</button>
      <button type="button" data-action="close">${esc(T("قفل", "Close"))}</button>
      <nav class="tv-branches" aria-label="${attr(T("الفروع", "Branches"))}"></nav></header>
      <div class="tv-body"><div class="tv-map"><svg class="tv-canvas" xmlns="http://www.w3.org/2000/svg" role="group" tabindex="0" aria-label="${attr(T("شجرة التقنيات: اسحب للتحريك وقرّب للتكبير", "Tech tree: drag to pan and pinch to zoom"))}"><g class="tv-world"></g></svg>
      <p class="tv-help">${esc(T("اسحب للتحريك · عجلة أو قرصة للتكبير · + / − للتكبير بالكيبورد", "Drag to pan · Wheel or pinch to zoom · + / − for keyboard zoom"))}</p></div>
      <aside class="tv-detail" aria-live="polite" aria-label="${attr(T("تفاصيل التقنية", "Technology details"))}"></aside></div>`;
    document.body.appendChild(overlay);
    svg = overlay.querySelector(".tv-canvas"); world = overlay.querySelector(".tv-world"); detail = overlay.querySelector(".tv-detail");
    bindControls(); bindPointers();
    new ResizeObserver(resize).observe(svg);
  }

  function bindControls() {
    overlay.querySelector('[data-action="close"]').onclick = close;
    overlay.querySelector('[data-action="centre"]').onclick = centre;
    overlay.querySelector("input").oninput = (ev) => { query = ev.target.value; renderGraph(); };
    overlay.querySelector(".tv-branches").onclick = (ev) => {
      const button = ev.target.closest("[data-branch]");
      if (button) { filter = button.dataset.branch; renderBranches(); renderGraph(); }
    };
    overlay.addEventListener("cancel", (ev) => { ev.preventDefault(); close(); });
    svg.addEventListener("keydown", canvasKey);
    svg.addEventListener("click", (ev) => {
      const node = ev.target.closest("[data-node]");
      if (ev.detail === 0 && node) selectNode(node.dataset.node);
    });
    svg.addEventListener("wheel", (ev) => {
      ev.preventDefault();
      const point = canvasPoint(ev), delta = ev.deltaY * (ev.deltaMode === 1 ? 16 : ev.deltaMode === 2 ? height : 1);
      zoomAt(Math.exp(-delta * 0.002), point, point);
    }, { passive: false });
    detail.addEventListener("click", detailClick);
  }

  function open(nextTree, nextHandlers) {
    if (!overlay) createOverlay();
    handlers = nextHandlers || {};
    if (!overlay.open) {
      returnFocus = document.activeElement;
      overlay.showModal();
      pan = { x: 0, y: 0 }; zoom = 1; selected = null; filter = ""; query = "";
      overlay.querySelector("input").value = "";
    }
    update(nextTree); resize(); centre();
  }

  function close() {
    if (!overlay || !overlay.open) return;
    clearPointers(); overlay.close();
    if (returnFocus && returnFocus.isConnected) returnFocus.focus();
  }

  function update(nextTree) {
    tree = nextTree;
    byId = Object.fromEntries(tree.nodes.map((node) => [node.id, node]));
    positions = layoutTree(tree.nodes);
    branches = [...new Set(tree.nodes.map((node) => node.branch).filter(Boolean))].sort();
    if (!byId[selected]) selected = null;
    if (!branches.includes(filter)) filter = "";
    if (!overlay || !overlay.open) return;
    const stats = overlay.querySelector(".tv-stats");
    if (stats) {
      const html = headerStatsMarkup(tree);
      stats.innerHTML = html;
      stats.hidden = !html;
    }
    renderBranches(); renderGraph(); renderDetail();
    minZoom = Math.min(0.25, fittedView().zoom);
  }

  function headerStatsMarkup(state = tree) {
    if (!state) return "";
    const parts = [];
    if (typeof state.rank === "number" && Number.isFinite(state.rank)) {
      parts.push(`<span class="tv-stat tv-stat-rank"><span>${esc(T("المرتبة", "Rank"))}</span><b dir="ltr">${fmt(state.rank)}</b></span>`);
    }
    if (typeof state.rp === "number" && Number.isFinite(state.rp)) {
      parts.push(`<span class="tv-stat tv-stat-rp"><span>${esc(T("نقاط البحث", "RP"))}</span><b dir="ltr">${fmt(state.rp)}</b></span>`);
    }
    return parts.join("");
  }

  function renderBranches() {
    const focusedBranch = document.activeElement && document.activeElement.dataset.branch;
    overlay.querySelector(".tv-branches").innerHTML = `<button type="button" data-branch="" aria-pressed="${!filter}">${esc(T("كل الفروع", "All branches"))}</button>` + branches.map((branch) => {
      const node = tree.nodes.find((n) => n.branch === branch);
      return `<button type="button" data-branch="${attr(branch)}" aria-pressed="${filter === branch}" style="--branch:${color(branch)}"><i aria-hidden="true"></i>${esc(local(node, "branch"))}</button>`;
    }).join("");
    if (focusedBranch !== undefined) [...overlay.querySelectorAll("[data-branch]")].find((el) => el.dataset.branch === focusedBranch)?.focus();
  }

  function stateLabel(node) {
    if (node.state === "owned") return T(`مستوى ${fmt(node.level)}`, `Level ${fmt(node.level)}`);
    if (node.state === "available") return T("متاحة", "Available");
    if (node.state === "near") return T(`ناقص ${missingCount(node)}`, `${missingCount(node)} missing`);
    if (node.state === "locked") return T("مقفولة", "Locked");
    return T("مجهولة", "Unknown");
  }

  function renderGraph() {
    const focusedId = document.activeElement && document.activeElement.dataset.node;
    const chain = selected && byId[selected].state === "locked" ? pathToNearestOwned(selected, byId) : [];
    const pathEdges = new Set(chain.slice(1).map((id, i) => JSON.stringify([chain[i], id])));
    world.innerHTML = tree.nodes.map((node) => node.after.map((prereq) => edgeMarkup(node, prereq, pathEdges)).join("")).join("") + tree.nodes.map((node) => nodeMarkup(node, chain)).join("");
    transform();
    if (focusedId) [...world.querySelectorAll("[data-node]")].find((el) => el.dataset.node === focusedId)?.focus();
  }

  function edgeMarkup(node, prereq, pathEdges) {
    const from = positions[prereq.id], to = positions[node.id];
    const satisfied = edgeSatisfied(node, prereq, byId), onPath = pathEdges.has(JSON.stringify([prereq.id, node.id]));
    const dim = filter && node.branch !== filter && byId[prereq.id].branch !== filter;
    const dx = to.x - from.x;
    return `<path class="tv-edge${satisfied ? " is-satisfied" : ""}${onPath ? " is-path" : ""}${dim ? " is-dim" : ""}" style="--branch:${color(byId[prereq.id].branch)}" d="M ${from.x} ${from.y} C ${from.x + dx / 2} ${from.y}, ${to.x - dx / 2} ${to.y}, ${to.x} ${to.y}"/>`;
  }

  function nodeMarkup(node, chain) {
    const point = positions[node.id], unknown = node.state === "unknown";
    const name = unknown ? T("؟", "؟") : local(node, "name");
    const classes = `tv-node is-${node.state}${selected === node.id ? " is-selected" : ""}${chain.includes(node.id) ? " is-path" : ""}${filter && node.branch !== filter ? " is-dim" : ""}${matchesSearch(node, query) ? " is-match" : ""}`;
    const lock = node.state === "locked" ? '<path class="tv-lock" d="M 54 8 v-3 a3 3 0 0 1 6 0 v3 m-8 0 h10 v9 h-10 z"/>' : "";
    return `<g class="${attr(classes)}" data-node="${attr(node.id)}" style="--branch:${color(node.branch)}" transform="translate(${point.x} ${point.y})" role="button" tabindex="0" aria-pressed="${selected === node.id}" aria-label="${attr(name + (unknown ? "" : ": " + stateLabel(node)))}">
      <title>${esc(name)}</title><rect x="-66" y="-22" width="132" height="44" rx="2"/>
      <svg x="-62" y="-18" width="124" height="36"><text x="62" y="${unknown ? 26 : 15}" class="tv-name" text-anchor="middle" direction="${LANG === "en" ? "ltr" : "rtl"}">${esc(name)}</text></svg>
      ${unknown ? "" : `<text class="tv-state" x="0" y="15" text-anchor="middle">${esc(stateLabel(node))}</text>${lock}`}</g>`;
  }

  function selectNode(id) {
    if (!byId[id]) return;
    selected = id; renderGraph(); renderDetail();
  }

  function renderDetail() {
    const node = byId[selected];
    if (!node) { detail.innerHTML = `<p>${esc(T("اختار عقدة عشان تشوف شروطها وترقيتها.", "Select a node to see its requirements and upgrade."))}</p>`; return; }
    if (node.state === "unknown") {
      detail.innerHTML = `<h3>${esc(T("؟", "؟"))}</h3><p>${esc(T("كمّل بناء التقنيات القريبة عشان تكشف العقدة دي.", "Build nearby technologies to reveal this node."))}</p>`;
      return;
    }
    detail.innerHTML = `<h3>${esc(local(node, "name"))}</h3><p class="tv-detail-state">${esc(stateLabel(node))}</p>
      <p>${esc(local(node, "desc"))}</p><h4>${esc(T(`الأثر في المستوى ${fmt(node.level + 1)}`, `Effect at level ${fmt(node.level + 1)}`))}</h4><p>${esc(local(node, "effect"))}</p>
      <h4>${esc(T("شروط الترقية", "Upgrade requirements"))}</h4><ul class="tv-checks">${node.checks.map((c, i) => checkMarkup(c, i, node)).join("")}</ul>
      <h4>${esc(T("التكلفة · معاك / المطلوب", "Cost · held / required"))}</h4><ul class="tv-cost">${costMarkup(node)}</ul>
      <button type="button" class="tv-buy" ${buyEnabled(node) ? "" : "disabled"}>${esc(buyLabel(node))}</button>`;
  }

  function checkLabel(check, node) {
    const raw = local(check, "label").trim();
    const num = check.need ?? check.required ?? check.target ?? check.level ?? check.rp ?? check.rank ?? check.value ?? check.ref;
    if (check.kind === "rp") {
      if (raw && !/^\d+(?:\.\d+)?$/.test(raw)) return raw;
      const amt = raw || num || (node && node.rp);
      return amt !== undefined && amt !== ""
        ? T(`نقاط البحث: ${fmt(Number(amt) || 0)}`, `Research points: ${fmt(Number(amt) || 0)}`)
        : T("نقاط بحث كافية", "Enough research points");
    }
    if (check.kind === "rank") {
      if (raw && !/^\d+(?:\.\d+)?$/.test(raw)) return raw;
      const lvl = raw || num;
      return lvl !== undefined && lvl !== ""
        ? T(`المرتبة ${fmt(Number(lvl) || 0)} أو أعلى`, `Rank ${fmt(Number(lvl) || 0)} or higher`)
        : T("مرتبة المستعمرة", "Colony rank");
    }
    return raw;
  }

  function checkMarkup(check, index, node) {
    const link = !check.ok && (check.kind === "parent" || check.kind === "contract");
    const label = esc(checkLabel(check, node)), mark = esc(check.ok ? T("✓", "✓") : T("✗", "✗"));
    const kindCls = check.kind ? ` tv-check-${attr(check.kind)}` : "";
    return `<li class="${check.ok ? "is-ok" : "is-missing"}${kindCls}"><span aria-hidden="true">${mark}</span><span class="tv-sr">${esc(check.ok ? T("مكتمل: ", "Met: ") : T("ناقص: ", "Missing: "))}</span>${link ? `<button type="button" data-check="${index}">${label}</button>` : `<span>${label}</span>`}</li>`;
  }

  function costMarkup(node, state = tree) {
    const names = typeof NAME !== "undefined" ? NAME : {};
    const rpCost = Number((node && (node.rp ?? (node.cost && node.cost.rp))) || 0);
    const costs = [[T("كريدت", "Credits"), (state && state.credits) || 0, (node.cost && node.cost.credits) || 0]];
    if (rpCost > 0) costs.push([T("نقاط البحث", "RP"), (state && state.rp) || 0, rpCost]);
    for (const [good, needed] of Object.entries((node.cost && node.cost.goods) || {})) {
      costs.push([T(names[good] || good, names[good] || good), (state && state.inv && state.inv[good]) || 0, needed]);
    }
    return costs.map(([name, held, needed]) => `<li class="${held >= needed ? "is-ok" : "is-missing"}"><span>${esc(name)}</span><b dir="ltr">${esc(T(`${fmt(held)} / ${fmt(needed)}`, `${fmt(held)} / ${fmt(needed)}`))}</b></li>`).join("");
  }

  function detailClick(ev) {
    const button = ev.target.closest("button"), node = byId[selected];
    if (!button || !node) return;
    if (button.classList.contains("tv-buy") && buyEnabled(node)) { handlers.onBuy?.(node.id); return; }
    const check = node.checks[Number(button.dataset.check)];
    if (!check || check.ok) return;
    if (check.kind === "parent") selectNode(check.ref);
    if (check.kind === "contract") handlers.onOpenContract?.(check.ref);
  }

  function resize() {
    if (!overlay.open) return;
    const rect = svg.getBoundingClientRect();
    width = rect.width; height = rect.height;
    minZoom = Math.min(0.25, fittedView().zoom);
    svg.setAttribute("viewBox", `0 0 ${width} ${height}`); transform();
  }

  function transform() {
    world.setAttribute("transform", `translate(${width / 2 + pan.x} ${height / 2 + pan.y}) scale(${zoom})`);
    overlay.classList.toggle("is-overview", zoom < 0.55);
  }
  function fittedView() {
    const known = Object.fromEntries(tree.nodes.filter((node) => node.state !== "unknown").map((node) => [node.id, positions[node.id]]));
    return fitView(layoutBounds(known), width, height, 32);
  }
  function centre() {
    const fit = fittedView();
    zoom = fit.zoom; minZoom = Math.min(0.25, zoom);
    pan = { x: fit.x - width / 2, y: fit.y - height / 2 }; transform();
  }
  function canvasPoint(ev) { const rect = svg.getBoundingClientRect(); return { x: ev.clientX - rect.left, y: ev.clientY - rect.top }; }

  function zoomAt(factor, before, after) {
    const next = Math.max(minZoom, Math.min(3, zoom * factor)), ratio = next / zoom;
    pan.x = after.x - width / 2 - (before.x - width / 2 - pan.x) * ratio;
    pan.y = after.y - height / 2 - (before.y - height / 2 - pan.y) * ratio;
    zoom = next; transform();
  }

  function canvasKey(ev) {
    const node = ev.target.closest("[data-node]");
    if (node && (ev.key === "Enter" || ev.key === " ")) { ev.preventDefault(); selectNode(node.dataset.node); return; }
    const step = { ArrowLeft: [40, 0], ArrowRight: [-40, 0], ArrowUp: [0, 40], ArrowDown: [0, -40] }[ev.key];
    if (step) { ev.preventDefault(); pan.x += step[0]; pan.y += step[1]; transform(); }
    if (["+", "=", "-"].includes(ev.key)) {
      ev.preventDefault(); const point = { x: width / 2, y: height / 2 };
      zoomAt(ev.key === "-" ? 1 / 1.2 : 1.2, point, point);
    }
  }

  function bindPointers() {
    svg.addEventListener("pointerdown", (ev) => {
      if (ev.button !== 0) return;
      if (!pointers.size) { moved = false; pressedNode = ev.target.closest("[data-node]")?.dataset.node; }
      else moved = true;
      const point = canvasPoint(ev);
      pointers.set(ev.pointerId, { ...point, startX: point.x, startY: point.y });
      svg.setPointerCapture(ev.pointerId);
    });
    svg.addEventListener("pointermove", pointerMove);
    svg.addEventListener("pointerup", pointerUp);
    svg.addEventListener("pointercancel", () => { moved = true; clearPointers(); });
    svg.addEventListener("lostpointercapture", (ev) => { pointers.delete(ev.pointerId); });
  }

  function pointerMove(ev) {
    const previous = pointers.get(ev.pointerId);
    if (!previous) return;
    const before = [...pointers.values()], point = canvasPoint(ev);
    if (Math.hypot(point.x - previous.startX, point.y - previous.startY) > 5) moved = true;
    pointers.set(ev.pointerId, { ...previous, ...point });
    const after = [...pointers.values()];
    if (before.length === 1) { pan.x += point.x - previous.x; pan.y += point.y - previous.y; transform(); return; }
    const distance = (pair) => Math.hypot(pair[0].x - pair[1].x, pair[0].y - pair[1].y);
    const midpoint = (pair) => ({ x: (pair[0].x + pair[1].x) / 2, y: (pair[0].y + pair[1].y) / 2 });
    if (distance(before) > 0) zoomAt(distance(after) / distance(before), midpoint(before), midpoint(after));
  }

  function pointerUp(ev) {
    if (!pointers.has(ev.pointerId)) return;
    pointers.delete(ev.pointerId);
    if (svg.hasPointerCapture(ev.pointerId)) svg.releasePointerCapture(ev.pointerId);
    if (!pointers.size && !moved && pressedNode) selectNode(pressedNode);
  }

  function clearPointers() {
    for (const id of pointers.keys()) if (svg.hasPointerCapture(id)) svg.releasePointerCapture(id);
    pointers.clear(); pressedNode = null;
  }

  if (typeof module !== "undefined" && module.exports) module.exports = { fitView, edgeSatisfied, missingCount, pathToNearestOwned, matchesSearch, buyEnabled, buyLabel, checkLabel, checkMarkup, costMarkup, headerStatsMarkup, branchHue };
  return { open, close, update, get isOpen() { return !!overlay && overlay.open; } };
})();
