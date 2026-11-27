"use strict";

// Chat: the in-game helper. It asks /api/chat, which answers through the
// player's own ChatGPT sign-in (Codex CLI). Each mission keeps its own
// conversation for this session.
(function () {
  const escape = (s) => String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");

  // render turns a reply's small Markdown (fences, `code`, **bold**, bullets,
  // numbered lines) into safe HTML: everything is escaped first.
  function render(text) {
    const parts = String(text).split(/```[\w-]*\n?/);
    return parts.map((part, i) => {
      if (i % 2) return `<pre dir="ltr"><code>${escape(part.replace(/\n$/, ""))}</code></pre>`;
      const lines = escape(part).split("\n");
      let html = "", list = null;
      const close = () => { if (list) { html += `</${list}>`; list = null; } };
      for (const raw of lines) {
        const line = raw
          .replace(/`([^`]+)`/g, '<code dir="ltr">$1</code>')
          .replace(/\*\*([^*]+)\*\*/g, "<b>$1</b>");
        const bullet = /^\s*[-*•]\s+(.*)$/.exec(line), number = /^\s*\d+[.)]\s+(.*)$/.exec(line);
        if (bullet || number) {
          const kind = bullet ? "ul" : "ol";
          if (list !== kind) { close(); html += `<${kind}>`; list = kind; }
          html += `<li>${(bullet || number)[1]}</li>`;
        } else {
          close();
          if (line.trim()) html += `<p>${line}</p>`;
        }
      }
      close();
      return html;
    }).join("");
  }

  // suggestions are one-tap questions; the last run decides the first one.
  function suggestions(T, hasError) {
    return [
      hasError ? T("ليه الكود بيطلع error؟", "Why does my code fail?") : T("مش فاهم المهمة", "I don't get the mission"),
      T("إديني تلميح صغير", "Give me a small hint"),
      T("الكود بتاعي ناقصه إيه؟", "What is my code missing?"),
    ];
  }

  // errorText explains a failed answer in the player's language.
  function errorText(T, result) {
    if (result?.error === "limit") return T("خلّصت أسئلة المساعد النهارده. بكرة تقدر تسأل تاني.", "You've used today's helper questions. You can ask again tomorrow.");
    if (result?.setup) return T("المساعد محتاج Codex CLI متسطب ومتسجّل بحساب ChatGPT. في الـterminal: brew install codex وبعدين codex login.", "The helper needs the Codex CLI signed in with ChatGPT. In a terminal: brew install codex, then codex login.");
    return T("المساعد مردّش المرة دي. جرّب تاني كمان شوية.", "The helper didn't answer this time. Try again in a moment.");
  }

  // ---- the panel (browser only) ----
  // host: { T, level() → id, code() → {file, code}, console() → text, hasError(), api(path, body) }
  function mount(host) {
    const panel = document.getElementById("chat");
    const list = panel.querySelector(".chat-log");
    const form = panel.querySelector("form");
    const input = form.querySelector("textarea");
    const send = form.querySelector("button[type=submit]");
    const chips = panel.querySelector(".chat-chips");
    const threads = new Map(); // level id → [{role, text}]
    let busy = false;

    const thread = () => {
      const id = host.level();
      if (!threads.has(id)) threads.set(id, []);
      return threads.get(id);
    };

    function draw() {
      const messages = thread();
      list.innerHTML = messages.length ? "" : `<p class="chat-empty">${escape(host.T("اسألني عن المهمة أو عن الكود بتاعك. هديك تلميحات، مش الحل على طول.", "Ask me about the mission or your code. I'll give hints, not the answer straight away."))}</p>`;
      for (const m of messages) {
        const div = document.createElement("div");
        div.className = `chat-msg ${m.role}${m.error ? " error" : ""}`;
        div.innerHTML = m.role === "assistant" ? render(m.text) : `<p>${escape(m.text)}</p>`;
        list.appendChild(div);
      }
      if (busy) {
        const typing = document.createElement("div");
        typing.className = "chat-msg assistant typing";
        typing.innerHTML = `<span></span><span></span><span></span><em>${escape(host.T("بيفكّر…", "Thinking…"))}</em>`;
        list.appendChild(typing);
      }
      chips.innerHTML = "";
      if (!messages.length) for (const q of suggestions(host.T, host.hasError())) {
        const b = document.createElement("button");
        b.type = "button"; b.className = "tool small"; b.textContent = q;
        b.onclick = () => ask(q);
        chips.appendChild(b);
      }
      send.disabled = busy;
      list.scrollTop = list.scrollHeight;
    }

    async function ask(text) {
      text = text.trim();
      if (!text || busy) return;
      const messages = thread();
      messages.push({ role: "user", text });
      busy = true; input.value = ""; draw();
      const { file, code } = host.code();
      let result;
      try {
        result = await host.api("/api/chat", {
          level: host.level(), file, code, console: host.console(),
          messages: messages.filter((m) => !m.error).map(({ role, text }) => ({ role, text })),
        });
      } catch (e) { result = { error: String(e) }; }
      busy = false;
      if (result?.reply) messages.push({ role: "assistant", text: result.reply });
      else messages.push({ role: "assistant", text: errorText(host.T, result), error: true });
      draw();
      if (!panel.hidden) input.focus();
    }

    form.addEventListener("submit", (e) => { e.preventDefault(); ask(input.value); });
    input.addEventListener("keydown", (e) => {
      if (e.key === "Enter" && !e.shiftKey && !e.isComposing) { e.preventDefault(); ask(input.value); }
      if (e.key === "Escape") close();
    });
    panel.querySelector(".chat-close").onclick = () => close();

    function open() {
      panel.hidden = false;
      document.body.classList.add("chat-open");
      draw();
      input.focus();
    }
    function close() {
      panel.hidden = true;
      document.body.classList.remove("chat-open");
      document.getElementById("btn-chat")?.focus();
    }
    return { open, close, toggle: () => (panel.hidden ? open() : close()), refresh: () => { if (!panel.hidden) draw(); } };
  }

  const api = { render, suggestions, errorText, mount };
  if (typeof window !== "undefined") window.Chat = api;
  if (typeof module !== "undefined") module.exports = api;
})();
