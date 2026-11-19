"use strict";

const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");

class TestElement {
  constructor(tag, { animations = false, animationCalls = [], onFocus = () => {} } = {}) {
    this.tagName = tag.toLowerCase();
    this.children = [];
    this.dataset = {};
    this.attributes = {};
    this.className = "";
    this._text = "";
    this._html = "";
    this.parentElement = null;
    this.hidden = false;
    this.disabled = false;
    this.scrollTop = 0;
    this.offsetTop = 0;
    this.offsetHeight = 20;
    this.clientHeight = 100;
    const properties = new Map();
    this.style = {
      setProperty: (name, value) => properties.set(name, value),
      removeProperty: (name) => properties.delete(name),
      getPropertyValue: (name) => properties.get(name) || "",
    };
    this.classList = {
      add: (...names) => { this.className = [...new Set([...this.className.split(/\s+/).filter(Boolean), ...names])].join(" "); },
      remove: (...names) => { this.className = this.className.split(/\s+/).filter((name) => name && !names.includes(name)).join(" "); },
      toggle: (name, force) => {
        const present = this.classList.contains(name);
        const next = force === undefined ? !present : Boolean(force);
        if (next) this.classList.add(name); else this.classList.remove(name);
        return next;
      },
      contains: (name) => this.className.split(/\s+/).includes(name),
    };
    if (animations) {
      this.getBoundingClientRect = () => ({ left: (this.parentElement?.children.indexOf(this) || 0) * 64, top: 0 });
      this.animate = (...animation) => animationCalls.push([this, ...animation]);
    }
    this.onFocus = onFocus;
  }

  set textContent(value) {
    this.replaceChildren();
    this._text = String(value);
    this._html = "";
  }

  get textContent() {
    return this._text + this._html.replace(/<[^>]*>/g, "") + this.children.map((child) => child.textContent).join("");
  }

  set innerHTML(value) {
    this.replaceChildren();
    this._html = String(value);
    this._text = "";
  }

  get innerHTML() { return this._html; }

  append(...nodes) {
    for (const node of nodes) {
      if (node.parentElement) node.parentElement.children = node.parentElement.children.filter((child) => child !== node);
      node.parentElement = this;
      this.children.push(node);
    }
  }

  appendChild(node) { this.append(node); return node; }

  replaceChildren(...nodes) {
    for (const child of this.children) child.parentElement = null;
    this.children = [];
    this._text = "";
    this._html = "";
    this.append(...nodes);
  }

  setAttribute(name, value) { this.attributes[name] = String(value); }
  getAttribute(name) { return this.attributes[name] ?? null; }
  focus() { this.onFocus(this); }
  click() { if (!this.disabled) this.onclick?.(); }

  querySelectorAll(selector) {
    const selectors = selector.split(",").map((part) => part.trim());
    const matches = (node, query) => {
      const attribute = query.match(/^([a-z][\w-]*)?\[([^\]=]+)(?:=["']?([^\]"']+)["']?)?\]$/i);
      if (attribute) return (!attribute[1] || node.tagName === attribute[1].toLowerCase()) && node.attributes[attribute[2]] !== undefined && (attribute[3] === undefined || node.attributes[attribute[2]] === attribute[3]);
      if (query.startsWith(".")) return node.classList.contains(query.slice(1));
      return node.tagName === query.toLowerCase();
    };
    return this.children.flatMap((child) => [
      ...(selectors.some((query) => matches(child, query)) ? [child] : []),
      ...child.querySelectorAll(selector),
    ]);
  }

  querySelector(selector) { return this.querySelectorAll(selector)[0] || null; }
}

function sandboxChapterSources(chapterNumber) {
  const web = path.resolve(__dirname, "..");
  return {
    lesson: fs.readFileSync(path.join(web, "lessons.js"), "utf8"),
    kit: fs.readFileSync(path.join(__dirname, "kit.js"), "utf8"),
    chapter: fs.readFileSync(path.join(__dirname, `ch${chapterNumber}.js`), "utf8"),
  };
}

function sandboxPageRoots() {
  return Object.fromEntries(["#lesson", "#lesson-panel", "#btn-lesson", "#hint-list"].map((id) => [id, new TestElement("section")]));
}

function createSandboxContext({ roots, store, lang }) {
  return vm.createContext({
    ResizeObserver: class { observe() {} disconnect() {} },
    document: {
      body: new TestElement("body"),
      createElement: (tag) => new TestElement(tag),
      querySelector: (selector) => roots["#lesson"].querySelector(selector),
    },
    $: (selector) => roots[selector],
    ed: new TestElement("textarea"),
    store,
    S: { review: false },
    T: (arabic, english) => lang === "en" ? english : arabic,
    md: (text) => String(text),
    inline: (text) => String(text),
    highlight: (text) => String(text),
    Sound: { play() {} },
  });
}

function loadSandboxScripts({ context, chapterNumber, sources }) {
  vm.runInContext(`${sources.lesson}\nthis.lessonUI = Lessons;`, context, { filename: "web/lessons.js" });
  vm.runInContext(sources.kit, context, { filename: "web/sandbox/kit.js" });
  vm.runInContext(sources.chapter, context, { filename: `web/sandbox/ch${chapterNumber}.js` });
}

function renderSandboxBlock({ context, roots, chapterNumber, kind, config, id, text }) {
  const blockID = id || `ch${chapterNumber}-playground`;
  const level = {
    id: blockID,
    lessonRevision: 1,
    steps: [{ id: `${blockID}-step`, title: "Try it", blocks: [{ id: blockID, kind: "sandbox", text, sandbox: { kind, config } }] }],
  };
  context.lessonUI.open(level);
  return roots["#lesson"].querySelector(".lesson-sandbox");
}

function clickSandboxButton(node) {
  if (!node || node.tagName !== "button" || node.disabled) throw new Error("expected an enabled button");
  node.click();
}

function createSandboxHarness({ chapter, lang = "en" }) {
  const chapterNumber = String(chapter).padStart(2, "0");
  const roots = sandboxPageRoots();
  const store = { blob: { lessons: {} }, save() {} };
  const context = createSandboxContext({ roots, store, lang });
  loadSandboxScripts({ context, chapterNumber, sources: sandboxChapterSources(chapterNumber) });

  return {
    context,
    roots,
    article: roots["#lesson"],
    register: (kind, render) => context.lessonUI.registerSandbox(kind, render),
    click: clickSandboxButton,
    render(kind, config, options = {}) {
      return renderSandboxBlock({ context, roots, chapterNumber, kind, config, ...options });
    },
  };
}

module.exports = { TestElement, createSandboxHarness };
