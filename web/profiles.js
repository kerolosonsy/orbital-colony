"use strict";
// Players sharing one computer. Each profile has its own progress, colony and
// language on the server; switching saves everything and reloads the page.
const Profiles = (() => {
  let list = [];
  const me = () => list.find((p) => p.id === PROFILE);

  async function init() {
    try { list = await api("/api/profiles"); } catch { return; }
    if (!me() && list.length) return switchTo(list[0].id, list[0].lang);
    const p = me();
    if (p && (p.lang === "en" ? "en" : "ar") !== LANG) return switchTo(p.id, p.lang);
    $("#profile-name").textContent = p ? p.name : "";
  }

  // persist this player's work, then come back as someone else
  async function switchTo(id, lang) {
    flushSave();
    await store.push();
    if (Colony.running) { try { await api("/api/colony/stop", {}); } catch { /* it stops when the next run starts */ } }
    localStorage.setItem(KEYS.profile, id);
    localStorage.setItem(KEYS.lang(id), lang === "en" ? "en" : "ar");
    location.reload();
  }

  function render() {
    const ul = $("#profile-list");
    ul.innerHTML = list.map((p) => `<li class="profile${p.id === PROFILE ? " current" : ""}" data-id="${esc(p.id)}">
      <span class="avatar" aria-hidden="true">${esc([...p.name][0] || "?")}</span>
      <span class="who"><b>${esc(p.name)}</b><small>${p.lang === "en" ? "English" : "العربية"}${p.id === PROFILE ? " · " + T("بتلعب دلوقتي", "playing now") : ""}</small></span>
      <span class="acts">
        ${p.id === PROFILE ? "" : `<button class="run small" data-act="play">${T("العب", "Play")}</button>`}
        <button class="tool small" data-act="rename">${T("غيّر الاسم", "Rename")}</button>
        <button class="tool small" data-act="duplicate" title="${T("لاعب جديد بنفس النجوم والكود والمستعمرة", "A new player with the same stars, code and colony")}">${T("اعمل نسخة", "Duplicate")}</button>
        <button class="tool small" data-act="reset">${T("ابدأ من الأول", "Start over")}</button>
        ${list.length > 1 ? `<button class="tool small danger" data-act="delete">${T("امسح", "Delete")}</button>` : ""}
      </span>
    </li>`).join("");
    ul.querySelectorAll("button[data-act]").forEach((b) => (b.onclick = () => act(b.dataset.act, list.find((p) => p.id === b.closest("li").dataset.id))));
  }

  async function act(what, p) {
    const path = `/api/profiles/${encodeURIComponent(p.id)}`;
    try {
      if (what === "play") return switchTo(p.id, p.lang);
      if (what === "rename") {
        const name = prompt(T("الاسم الجديد:", "New name:"), p.name)?.trim();
        if (!name) return;
        await api(path + "/rename", { name });
      }
      if (what === "duplicate") {
        const name = prompt(T("اسم النسخة:", "Name for the copy:"), p.name + " 2")?.trim();
        if (!name) return;
        await api(path + "/duplicate", { name });
      }
      if (what === "reset") {
        if (!confirm(T(`${p.name} هيبدأ من الأول: كل النجوم والكود والمستعمرة هيتمسحوا. متأكد؟`,
          `${p.name} will start over: all stars, code and the colony will be erased. Are you sure?`))) return;
        await api(path + "/reset", {});
        localStorage.removeItem(KEYS.cache(p.id));
        if (p.id === PROFILE) { clearTimeout(store.timer); store.timer = null; clearTimeout(saveTimer); saveTimer = null; location.reload(); return; }
      }
      if (what === "delete") {
        if (!confirm(T(`تمسح ${p.name} خالص؟ مفيش رجوع.`, `Delete ${p.name} for good? This can't be undone.`))) return;
        await api(path, null, "DELETE");
        localStorage.removeItem(KEYS.cache(p.id));
        localStorage.removeItem(KEYS.lang(p.id));
        if (p.id === PROFILE) {
          clearTimeout(store.timer); store.timer = null; clearTimeout(saveTimer); saveTimer = null;
          list = await api("/api/profiles");
          localStorage.setItem(KEYS.profile, list[0].id);
          localStorage.setItem(KEYS.lang(list[0].id), list[0].lang === "en" ? "en" : "ar");
          location.reload();
          return;
        }
      }
      list = await api("/api/profiles");
      render();
      $("#profile-name").textContent = me()?.name || "";
    } catch (e) {
      alert(T("حصلت مشكلة: ", "Something went wrong: ") + e.message);
    }
  }

  $("#btn-profile").onclick = async () => {
    try { list = await api("/api/profiles"); } catch { /* show what we have */ }
    render();
    $("#profile-new-lang").value = LANG;
    showOverlay("#profiles");
  };

  $("#profile-new").onsubmit = async (e) => {
    e.preventDefault();
    const name = $("#profile-new-name").value.trim();
    if (!name) return;
    try {
      const p = await api("/api/profiles", { name, lang: $("#profile-new-lang").value });
      await switchTo(p.id, p.lang);
    } catch (err) { alert(T("مقدرتش أضيف اللاعب: ", "Couldn't add the player: ") + err.message); }
  };

  $("#btn-lang").onclick = async () => {
    const lang = LANG === "en" ? "ar" : "en";
    try { await api(`/api/profiles/${encodeURIComponent(PROFILE)}/lang`, { lang }); } catch { /* still switch locally */ }
    switchTo(PROFILE, lang);
  };

  // hosted: the account behind these players, and signing out of it
  if (HOSTED) {
    const title = $("#profiles-title");
    title.textContent = T(title.dataset.hosted, title.dataset.hostedEn);
    $("#account-box").hidden = false;
    if (document.querySelector('meta[name="orbital-vscode"]')?.content === "off") $("#btn-vscode").hidden = true;
    api("/api/auth/me").then((me) => ($("#account-email").textContent = me.email)).catch(() => {});
    $("#btn-logout").onclick = async () => {
      flushSave();
      await store.push();
      try { await api("/api/auth/logout", {}); } finally { location.replace("/"); }
    };
  }

  return { init };
})();
