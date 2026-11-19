"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { createSandboxHarness } = require("./harness");

const pair = (ar, en) => ({ ar, en });
const fixtures = [
  {
    id: "c371.atomic-playground",
    kind: "ch37-atomic-ledger",
    config: {
      labels: {
        title: pair("تغييرين مع بعض", "Two changes together"),
        schema: pair("schema", "schema"), ledger: pair("سجل النسخ", "version ledger"),
        fail: pair("محاكاة فشل", "Simulate failure"), commit: pair("محاكاة نجاح", "Simulate commit"), reset: pair("ابدأ تاني", "Reset"),
        ready: pair("جاهز لتغيير تجريبي", "Ready for a sample change"), failed: pair("فشل؛ الاتنين زي ما هم", "Failed; both values stay"), committed: pair("نجاح؛ القيمتين اتحدثوا سوا", "Committed; both values moved"),
      },
      values: { schema: 4, ledger: 4, next: 5 },
    },
  },
  {
    id: "c372.session-playground",
    kind: "ch37-session-digest",
    config: {
      labels: {
        title: pair("التوكن ومخزونه", "Token and stored digest"),
        client: pair("مع العميل", "Client holds"), store: pair("المخزن", "Store holds"),
        issue: pair("إصدار جلسة", "Issue session"), logout: pair("خروج", "Log out"), reset: pair("ابدأ تاني", "Reset"),
        empty: pair("مفيش", "None"), ready: pair("مفيش جلسة", "No active session"), issued: pair("العميل ماسك التوكن", "Client holds the token"), ended: pair("الجلسة اتقفلت", "Session ended"), cleared: pair("الحالة اتصفّرت", "Session state cleared"),
      },
      values: { token: "token-amber", digest: "digest-amber" },
    },
  },
  {
    id: "c373.scope-playground",
    kind: "ch37-scoped-membership",
    config: {
      labels: {
        title: pair("المفتاح فيه نطاق", "The key includes scope"),
        key: pair("المفتاح المركّب", "Composite key"), role: pair("الدور", "Role"), choose: pair("اختار نطاق", "Choose a scope"),
        blue: pair("blue + mira", "blue + mira"), gold: pair("gold + mira", "gold + mira"), red: pair("red + mira", "red + mira"), empty: pair("—", "—"),
      },
      lookups: [
        { key: pair("blue + mira", "blue + mira"), role: pair("مالك", "owner") },
        { key: pair("gold + mira", "gold + mira"), role: pair("عضو", "member") },
        { key: pair("red + mira", "red + mira"), role: pair("مش موجود", "not found") },
      ],
    },
  },
  {
    id: "c374.cursor-playground",
    kind: "ch37-cursor-pages",
    config: {
      labels: {
        title: pair("صفحات بمؤشر ثابت", "Pages with a stable marker"), page: pair("الصفحة", "Page"), cursor: pair("cursor", "cursor"),
        next: pair("الصفحة الجاية", "Next page"), reset: pair("ارجع للبداية", "Back to start"),
        ready: pair("أول صفحة", "First page"), advanced: pair("المؤشر اتحرك مع الصفحة", "The marker moved with the page"), resetStatus: pair("رجعنا لأول صفحة", "Returned to the first page"),
      },
      pages: [["task-amber", "task-blue"], ["task-copper", "task-violet"]],
      cursors: ["start", "after-amber"],
    },
  },
  {
    id: "c375.request-playground",
    kind: "ch37-request-path",
    config: {
      labels: {
        title: pair("مسار الطلب", "Request path"), previous: pair("اللي قبلها", "Previous"), next: pair("اللي بعدها", "Next"), reset: pair("من الأول", "Reset"),
      },
      frames: [
        { stage: pair("وصول", "Arrival"), detail: "GET /v1/tasks?sample=1", caption: pair("طلب تجريبي وصل", "A sample request arrived") },
        { stage: pair("هوية الطلب", "Request identity"), detail: "req-42", caption: pair("اتضاف request ID", "A request ID was added") },
        { stage: pair("الرد", "Response"), detail: "202 Accepted", caption: pair("الرد رجع للعميل", "The response returned to the client") },
        { stage: pair("السجل", "Log"), detail: "GET /v1/tasks · 202", caption: pair("السجل فيه path من غير query", "The log keeps the path, not query") },
      ],
    },
  },
];

for (const fixture of fixtures) {
  for (const lang of ["ar", "en"]) {
    test(`${fixture.kind} renders and responds in ${lang}`, () => {
      const harness = createSandboxHarness({ chapter: 37, lang });
      const sandbox = harness.render(fixture.kind, fixture.config, { id: fixture.id, text: fixture.config.labels.title[lang] });
      assert.ok(sandbox, `missing renderer for ${fixture.kind}`);
      assert.equal(sandbox.querySelector(".ch37-title").textContent, fixture.config.labels.title[lang]);

      const controls = sandbox.querySelectorAll("button");
      assert.ok(controls.length > 0 && controls.length <= 4);
      assert.ok(controls.every((control) => control.tagName === "button"));
      for (const control of controls) {
        if (control.disabled) {
          const next = controls.find((candidate) => candidate.textContent === fixture.config.labels.next[lang]);
          assert.ok(next, "expected a next control to advance the stepper");
          harness.click(next);
          assert.equal(control.disabled, false, "previous control should enable after advancing");
        }
        const before = sandbox.textContent;
        harness.click(control);
        const after = sandbox.textContent;
        assert.notEqual(after, before, `${fixture.kind} control did not change visible state`);
      }
    });
  }
}
