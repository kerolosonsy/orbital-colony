"use strict";

// Tour: a guide for the first three missions. A pointer glides to each part
// of the screen, a spotlight lights it, and a short bubble says what it is.
// Some steps also act: the guide types code in front of the player.
//
// Step fields:
//   at    CSS selector, {tile: [x, y]} / {tile: "drone"} on the map,
//         {line: n, rows?} in the editor, or nothing (a centred bubble)
//   side  "above" puts the bubble over the target (to keep what's below visible)
//   ar/en the bubble text (short; ~code~ marks code)
//   act   {after, type}             type text right after the anchor `after`
//         {after, complete: [pre, post]}  on a new line after `after`: type pre,
//                                   accept the autocomplete, then type post
(function () {
  const TOURS = {
    "1-1-first-signal": [
      { ar: "أهلاً! هوريك اللعبة حتة حتة. دوس **التالي** أو Enter.", en: "Hi! Let me show you around. Press **Next** or Enter." },
      { at: "#goal-strip", ar: "ده **هدف المهمة**: توصّل الدرون للمنارة.", en: "This is the **mission goal**: get the drone to the beacon." },
      { at: { tile: "drone" }, ar: "ده **الدرون** بتاعك. بيتحرك بس بالكود اللي بتكتبه.", en: "This is your **drone**. It only moves with the code you write." },
      { at: { tile: [3, 0] }, ar: "ودي **المنارة** 📡. عاوزين نوصلها.", en: "And this is the **beacon** 📡. That's where we're going." },
      { at: "#lesson", ar: "هنا **الدرس**: شرح قصير للي محتاجه في المهمة دي.", en: "Here's the **lesson**: a short read on what this mission needs." },
      { at: ".editor", ar: "وده **المحرر**. هنا بتكتب كود Go.", en: "And this is the **editor**, where you write Go." },
      { at: { line: 6 }, ar: "السطر ده بيحرّك الدرون **خانة واحدة** ناحية الشرق.", en: "This line moves the drone **one tile** east." },
      { at: { line: 7 }, act: { after: "\tcolony.Move(colony.East)\n", type: "\tcolony.Move(colony.East)\n" },
        ar: "شوف، هكتبلك سطر كمان زيه. كده الدرون هيمشي خانتين.", en: "Watch: I'll type one more like it. Now the drone moves two tiles." },
      { at: "#btn-run", ar: "**Run** بيشغّل الكود (أو Ctrl/⌘ + Enter)، والدرون ينفّذ من فوق لتحت.", en: "**Run** runs your code (or Ctrl/⌘ + Enter); the drone follows it top to bottom." },
      { at: "#console", ar: "هنا هتشوف **الرسايل والأخطاء** بعد التشغيل.", en: "Messages and **errors** show up here after a run." },
      { at: "#lvl-stars", ar: "**النجوم**: 3 لو وصلت بأقل عدد خطوات.", en: "**Stars**: 3 if you get there in the fewest steps." },
      { at: "#btn-hint", ar: "اتلخبطت؟ **التلميح** بيديك فكرة من غير ما يحل عنك.", en: "Stuck? A **hint** nudges you without solving it." },
      { ar: "دورك: المنارة على بُعد **3 خانات**. زوّد السطر اللي ناقص ودوس **Run**!", en: "Your turn: the beacon is **3 tiles** away. Add the missing line and press **Run**!" },
    ],
    "1-2-first-harvest": [
      { ar: "مهمة جديدة، حاجات جديدة. يلا نلف لفة سريعة.", en: "New mission, new things. A quick look around." },
      { at: "#goal-strip", ar: "الهدف: تحصد **3 طحالب**.", en: "Goal: harvest **3 moss**." },
      { at: { tile: [1, 0] }, ar: "دي أول **طحلبة**، والاتنين اللي بعدها على اليمين.", en: "Here's the first **moss**; two more sit to its right." },
      { at: { line: 7 }, ar: "~colony.Harvest()~ بتحصد الخانة اللي الدرون **واقف عليها**.", en: "~colony.Harvest()~ harvests the tile the drone **stands on**." },
      { at: "#hud-inv", ar: "هنا **المخزن**: اللي بتحصده بيتعد هنا.", en: "This is your **cargo**: everything you harvest is counted here." },
      { at: "#hud-tick", ar: "ده **عدّاد الخطوات**. كل أمر بياخد خطوة.", en: "This is the **tick counter**. Every command takes one tick." },
      { at: "#lvl-par", ar: "عدد الخطوات اللي يجيبلك **3 نجوم**.", en: "The number of ticks that earns **3 stars**." },
      { at: "#speed", ar: "**السرعة**: خلّي الدرون أبطأ علشان تتابعه، أو أسرع.", en: "**Speed**: slow the drone down to follow it, or speed it up." },
      { at: "#btn-skip", ar: "وده بيوريك **النتيجة على طول** من غير ما تستنى.", en: "And this jumps **straight to the result**." },
      { ar: "دورك: **تتحرك وتحصد** لحد ما تجمع التلاتة. Run!", en: "Your turn: **move and harvest** until you have all three. Run!" },
    ],
    "1-3-planting": [
      { ar: "آخر لفة: أدوات هتريّحك وإنت بتكتب.", en: "Last tour: tools that make writing easier." },
      { at: { line: 7 }, ar: "~colony.Till()~ بتحرث الخانة علشان تبقى جاهزة للزرع.", en: "~colony.Till()~ ploughs the tile so it's ready to plant." },
      { at: { line: 6, rows: 9 }, side: "above", act: { after: "\tcolony.Till()\n", complete: ["\tcolony.Pl", "colony.Moss"] },
        ar: "اكتب ~colony.~ وهتطلعلك **قايمة بالأوامر**. Enter يختار، والشرح بيظهر تحت.", en: "Type ~colony.~ and a **list of commands** pops up. Enter picks one; its help shows below." },
      { at: "#btn-step", ar: "**خطوة** بتشغّل الكود سطر سطر، وتشوف كل سطر بيعمل إيه.", en: "**Step** runs your code line by line so you see what each line does." },
      { at: "#gutter", ar: "دوس على رقم سطر يحط **نقطة توقف** (F9): التشغيل يقف عندها.", en: "Click a line number for a **breakpoint** (F9): the run pauses there." },
      { at: "#btn-format", ar: "**{ }** بيرتّب الكود (Shift+Alt+F).", en: "**{ }** tidies your code (Shift+Alt+F)." },
      { at: "#btn-reset", ar: "بوّظت الكود؟ ده بيرجّع **كود البداية** (و Ctrl+Z يرجّع اللي اتمسح).", en: "Messed it up? This restores the **starter code** (Ctrl+Z brings yours back)." },
      { at: "#btn-missions", ar: "**الفهرس**: كل المهمات والفصول.", en: "The **index**: every mission and chapter." },
      { ar: "دوس **?** في أي وقت تشوف كل الاختصارات. دورك: ازرع التلات خانات!", en: "Press **?** anytime for all shortcuts. Your turn: plant all three tiles!" },
    ],
  };

  function stepsFor(levelID) {
    return TOURS[levelID] || null;
  }

  // bubbleBox places a w×h bubble next to target rect r inside a vw×vh view:
  // below it if it fits (above first when prefer is "above"), else beside;
  // always kept on screen.
  function bubbleBox(r, w, h, vw, vh, prefer = "below", gap = 14) {
    const clamp = (v, lo, hi) => Math.max(lo, Math.min(v, hi));
    if (!r) return { left: (vw - w) / 2, top: (vh - h) / 2, side: "center" };
    let left = clamp(r.left + r.width / 2 - w / 2, 12, vw - w - 12);
    let top, side;
    const below = r.bottom + gap + h <= vh - 12, above = r.top - gap - h >= 12;
    if (below && !(prefer === "above" && above)) { top = r.bottom + gap; side = "below"; }
    else if (above) { top = r.top - gap - h; side = "above"; }
    else {
      side = r.left > vw - r.right ? "left" : "right";
      left = side === "left" ? clamp(r.left - gap - w, 12, vw - w - 12) : clamp(r.right + gap, 12, vw - w - 12);
      top = clamp(r.top + r.height / 2 - h / 2, 12, vh - h - 12);
    }
    return { left, top, side };
  }

  // keyAction maps a key to "next" | "back" | "close" | null (arrows mirror in RTL).
  function keyAction(key, rtl) {
    if (key === "Escape") return "close";
    if (key === "Enter") return "next";
    if (key === "ArrowLeft" || key === "ArrowRight") return (key === "ArrowLeft") === rtl ? "next" : "back";
    return null;
  }

  // ---- the runner (browser only) ----
  let active = null;

  // play walks the steps. host: { locate(at) → rect|null, act(act) → Promise,
  // text(step) → html, label(ar, en), onEnd({finished, interrupted}) }.
  function play(steps, host) {
    stop();
    const root = document.createElement("div");
    root.className = "tour";
    root.innerHTML = `<div class="tour-veil"></div><div class="tour-spot"></div>
      <svg class="tour-hand" viewBox="0 0 24 24" aria-hidden="true"><path d="M5 2.5v16.2l4.2-3.9 2.9 6.6 2.9-1.3-2.9-6.5 5.8-.3z"/></svg>
      <div class="tour-bubble" role="dialog" aria-live="polite"><p class="tour-text"></p>
        <div class="tour-bar"><span class="tour-count"></span><span class="tour-btns">
          <button type="button" class="tool small tour-skip"></button>
          <button type="button" class="tool small tour-back"></button>
          <button type="button" class="tool small tour-next"></button></span></div></div>`;
    document.body.appendChild(root);
    const $ = (s) => root.querySelector(s);
    const spot = $(".tour-spot"), hand = $(".tour-hand"), bubble = $(".tour-bubble");
    $(".tour-skip").textContent = host.label("تخطّي", "Skip");
    $(".tour-back").textContent = host.label("السابق", "Back");
    // buttons must not steal the editor's focus (the autocomplete closes on blur)
    for (const b of root.querySelectorAll("button")) b.addEventListener("mousedown", (e) => e.preventDefault());

    const state = { index: -1, busy: false, done: false };
    active = state;
    let handAt = [window.innerWidth / 2, window.innerHeight / 2];
    hand.style.transform = `translate(${handAt[0]}px, ${handAt[1]}px)`;

    function place() {
      const step = steps[state.index];
      const r = step.at ? host.locate(step.at) : null;
      if (r) {
        const pad = 6;
        Object.assign(spot.style, { left: `${r.left - pad}px`, top: `${r.top - pad}px`, width: `${r.width + pad * 2}px`, height: `${r.height + pad * 2}px`, opacity: 1 });
        handAt = [r.left + Math.min(r.width / 2, 60), r.top + Math.min(r.height / 2, 40)];
        hand.classList.add("on");
        root.classList.remove("center");
      } else {
        root.classList.add("center");
        Object.assign(spot.style, { left: "50%", top: "50%", width: "0px", height: "0px", opacity: 0 });
        hand.classList.remove("on");
      }
      hand.style.transform = `translate(${handAt[0]}px, ${handAt[1]}px)`;
      const b = bubbleBox(r, bubble.offsetWidth, bubble.offsetHeight, window.innerWidth, window.innerHeight, step.side);
      Object.assign(bubble.style, { left: `${b.left}px`, top: `${b.top}px` });
      bubble.dataset.side = b.side;
    }

    async function go(index) {
      if (state.busy || state.done) return;
      // skip steps whose target is hidden in this layout
      while (index < steps.length && index >= 0 && steps[index].at && !host.locate(steps[index].at)) index += index >= state.index ? 1 : -1;
      if (index < 0) return;
      if (index >= steps.length) return end(true);
      state.index = index;
      const step = steps[index];
      $(".tour-text").innerHTML = host.text(step);
      $(".tour-count").textContent = `${index + 1} / ${steps.length}`;
      $(".tour-back").disabled = index === 0;
      $(".tour-next").textContent = index === steps.length - 1 ? host.label("يلا بينا", "Let's go") : host.label("التالي", "Next");
      place();
      if (step.act && !step.played) {
        state.busy = true;
        root.classList.add("busy");
        step.played = true;
        await new Promise((ok) => setTimeout(ok, 700)); // let the pointer arrive
        try { await host.act(step.act); } catch {}
        state.busy = false;
        root.classList.remove("busy");
        if (state.done) return;
        place();
      }
      hand.classList.remove("tap"); void hand.offsetWidth; hand.classList.add("tap");
      // after a typing demo the editor keeps focus, so its help stays on screen
      if (!step.act) $(".tour-next").focus({ preventScroll: true });
    }

    function end(finished, interrupted = false) {
      if (state.done) return;
      state.done = true;
      if (active === state) active = null;
      window.removeEventListener("keydown", onKey, true);
      window.removeEventListener("resize", onResize);
      root.remove();
      for (const s of steps) delete s.played;
      host.onEnd({ finished, interrupted });
    }

    function onKey(e) {
      const action = keyAction(e.key, document.documentElement.dir === "rtl");
      if (!action || e.altKey || e.ctrlKey || e.metaKey) return;
      e.preventDefault(); e.stopImmediatePropagation();
      if (action === "close") end(false);
      else go(state.index + (action === "next" ? 1 : -1));
    }
    const onResize = () => { if (state.index >= 0) place(); };
    window.addEventListener("keydown", onKey, true);
    window.addEventListener("resize", onResize);
    $(".tour-next").onclick = () => go(state.index + 1);
    $(".tour-back").onclick = () => go(state.index - 1);
    $(".tour-skip").onclick = () => end(false);
    state.stop = () => end(false, true); // the game moved on (another mission): not "seen"
    requestAnimationFrame(() => go(0));
  }

  function stop() {
    active?.stop?.();
  }

  const api = { stepsFor, bubbleBox, keyAction, play, stop, get running() { return !!active; } };
  if (typeof window !== "undefined") window.Tour = api;
  if (typeof module !== "undefined") module.exports = api;
})();
