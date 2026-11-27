"use strict";
// Sign-in page for the hosted game: sign in, make an account (with an invite
// code), or continue with Google. On success the page reloads into the game.
(function () {
  const $ = (s) => document.querySelector(s);
  const TOKEN = document.querySelector('meta[name="orbital-token"]').content;
  let lang = localStorage.getItem("orbital.login.lang") === "en" ? "en" : "ar";
  let mode = "login";
  const T = (ar, en) => (lang === "en" ? en : ar);

  const ERRORS = {
    bad_email: ["الإيميل ده مش مكتوب صح.", "That email doesn't look right."],
    weak_password: ["الباسورد لازم يكون ٨ حروف أو أكتر.", "The password needs 8 characters or more."],
    bad_invite: ["كود الدعوة مش صح.", "That invite code isn't right."],
    email_taken: ["فيه حساب بالإيميل ده. ادخل بيه بدل كده.", "There's already an account with this email. Sign in instead."],
    wrong_login: ["الإيميل أو الباسورد غلط.", "Wrong email or password."],
    too_many: ["محاولات كتير. استنى ربع ساعة وجرّب تاني.", "Too many tries. Wait 15 minutes and try again."],
    invite: ["علشان تعمل حساب جديد بـGoogle، اكتب كود الدعوة في \"حساب جديد\" الأول.", "To make a new account with Google, type your invite code under \"New account\" first."],
    unverified: ["إيميل Google ده مش متأكد منه.", "That Google email isn't verified."],
    expired: ["الدخول بـGoogle خد وقت طويل. جرّب تاني.", "Google sign-in took too long. Try again."],
    cancelled: ["الدخول بـGoogle اتلغى.", "Google sign-in was cancelled."],
    failed: ["حصلت مشكلة. جرّب تاني.", "Something went wrong. Try again."],
  };

  function showError(code) {
    const box = $("#login-error");
    if (!code) { box.hidden = true; return; }
    const [ar, en] = ERRORS[code] || ERRORS.failed;
    box.textContent = T(ar, en);
    box.hidden = false;
  }

  function applyLang() {
    document.documentElement.lang = lang;
    document.documentElement.dir = lang === "en" ? "ltr" : "rtl";
    document.querySelectorAll("[data-en]").forEach((el) => {
      if (el.dataset.ar === undefined) el.dataset.ar = el.textContent;
      el.textContent = lang === "en" ? el.dataset.en : el.dataset.ar;
    });
    $("#login-lang").textContent = lang === "en" ? "العربية" : "English";
    setMode(mode);
  }

  function setMode(next) {
    mode = next;
    const signup = mode === "signup";
    $("#tab-login").setAttribute("aria-selected", String(!signup));
    $("#tab-signup").setAttribute("aria-selected", String(signup));
    document.querySelectorAll(".signup-only").forEach((el) => (el.hidden = !signup));
    $("#submit-label").textContent = signup ? T("اعمل الحساب", "Create account") : T("ادخل", "Sign in");
    $("#f-password").autocomplete = signup ? "new-password" : "current-password";
    showError(null);
    googleHref();
  }

  function googleHref() {
    const q = new URLSearchParams({ lang, invite: $("#f-invite").value.trim() });
    $("#google-link").href = "/auth/google?" + q;
  }

  async function submit(e) {
    e.preventDefault();
    const body = { email: $("#f-email").value.trim(), password: $("#f-password").value, lang };
    if (mode === "signup") Object.assign(body, { name: $("#f-name").value.trim(), invite: $("#f-invite").value.trim() });
    const button = $("#login-submit");
    button.disabled = true;
    try {
      const res = await fetch(mode === "signup" ? "/api/auth/signup" : "/api/auth/login", {
        method: "POST", headers: { "X-Orbital-Token": TOKEN, "Content-Type": "application/json" }, body: JSON.stringify(body),
      });
      const out = await res.json().catch(() => ({}));
      if (!res.ok) { showError(out.error); return; }
      if (mode === "signup") localStorage.setItem(`orbital.lang.${out.id}.default`, lang);
      location.replace("/");
    } catch { showError("failed"); }
    finally { button.disabled = false; }
  }

  $("#tab-login").onclick = () => setMode("login");
  $("#tab-signup").onclick = () => setMode("signup");
  $("#login-lang").onclick = () => { lang = lang === "en" ? "ar" : "en"; localStorage.setItem("orbital.login.lang", lang); applyLang(); };
  $("#f-invite").oninput = googleHref;
  $("#login-form").onsubmit = submit;

  fetch("/api/auth/me", { headers: { "X-Orbital-Token": TOKEN } }).then((r) => r.json()).then((info) => {
    if (!info.googleEnabled) $("#google-box").hidden = true;
    if (info.inviteRequired) $("#google-invite-note").hidden = false;
    if (info.inviteRequired === false) {
      $("#f-invite").closest("label").classList.remove("signup-only");
      $("#f-invite").closest("label").hidden = true;
    }
  }).catch(() => {});

  applyLang();
  const failed = new URLSearchParams(location.search).get("auth");
  if (failed) {
    if (failed === "invite") setMode("signup");
    showError(failed);
    history.replaceState(null, "", "/");
  }
})();
