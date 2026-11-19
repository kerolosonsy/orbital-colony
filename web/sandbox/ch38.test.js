"use strict";

const assert = require("node:assert/strict");
const test = require("node:test");
const { createSandboxHarness } = require("./harness");

const pair = (ar, en) => ({ ar, en });
const fixtures = [
  {
    id: "c381.iterator-playground",
    kind: "ch38-iterator-stop",
    config: {
      labels: {
        title: pair("المستهلك يوقف", "The consumer can stop"), values: pair("القيم المصدر", "Source values"), consumed: pair("اتستهلك", "Consumed"),
        next: pair("هات قيمة", "Next value"), stop: pair("وقّف", "Stop"), reset: pair("ابدأ تاني", "Reset"),
        ready: pair("جاهز للاستهلاك", "Ready to consume"), read: pair("قيمة جديدة وصلت", "A new value arrived"), stopped: pair("المستهلك وقف", "The consumer stopped"), cleared: pair("السجل اتمسح", "The log was cleared"), ended: pair("كل القيم وصلت", "All values arrived"),
      },
      values: [2, 5, 8, 11],
    },
  },
  {
    id: "c382.set-playground",
    kind: "ch38-set-labels",
    config: {
      labels: {
        title: pair("عضوية وتخطّي صف", "Membership and row skipping"), row: pair("الصف", "Row"), slots: pair("الخانات", "Slots"), selected: pair("الـset", "Set"),
        next: pair("افحص الصف", "Check next row"), reset: pair("من الأول", "Reset"),
        empty: pair("—", "—"), ready: pair("جاهز لفحص الصفوف", "Ready to inspect rows"), picked: pair("قيم جديدة دخلت الـset", "New values joined the set"), skipped: pair("الـlabel تخطّى باقي الصف", "The label skipped the rest of this row"), duplicate: pair("القيمة موجودة قبل كده", "The value was already present"), done: pair("خلصت الصفوف", "All rows were inspected"), cleared: pair("الـset اتمسح", "The set was cleared"),
      },
      rows: [[7, 2], [0, 5], [2, 8], [7]],
    },
  },
  {
    id: "c383.once-playground",
    kind: "ch38-once-cache",
    config: {
      labels: {
        title: pair("النتيجة الصفرية برضه نتيجة", "Zero is still a result"), runs: pair("مرات التهيئة", "Initializer runs"), value: pair("القيمة", "Value"),
        call: pair("نادِ الـloader", "Call loader"), reset: pair("ابدأ تاني", "Reset"),
        ready: pair("جاهز لأول نداء", "Ready for the first call"), loaded: pair("التهيئة اشتغلت", "The initializer ran"), reused: pair("رجّع نفس القيمة", "The same result was reused"), cleared: pair("العداد رجع للصفر", "The counter returned to zero"), zero: 0,
      },
    },
  },
  {
    id: "c384.group-playground",
    kind: "ch38-errgroup-orders",
    config: {
      labels: {
        title: pair("ترتيب ممكن يتغيّر", "A changing possible order"), jobs: pair("ترتيب الإكمال", "Completion order"),
        run: pair("شغّل ترتيب تاني", "Run another sample"), reset: pair("ارجع للمثال الأول", "Return to first sample"),
        possible: pair("ده ترتيب ممكن", "One possible order"), another: pair("مثال تاني؛ الترتيب مش ثابت", "Another sample; order is not fixed"), restored: pair("رجعنا لمثال توضيحي", "Returned to an illustrative sample"),
      },
      orders: [["job-amber", "job-blue", "job-copper"], ["job-blue", "job-amber", "job-copper"], ["job-copper", "job-blue", "job-amber"]],
    },
  },
  {
    id: "c385.builder-playground",
    kind: "ch38-builder-buffer",
    config: {
      labels: {
        title: pair("الـBuilder بيجمع النص", "A Builder collects text"), buffer: pair("المحتوى", "Buffer"), length: pair("الطول", "Length"),
        write: pair("اكتب الجزء الجاي", "Write next part"), reset: pair("فضّي الـbuffer", "Clear buffer"),
        empty: pair("فاضي", "Empty"), ready: pair("الـbuffer فاضي", "The buffer is empty"), written: pair("النص اتجمع", "Text was collected"), cleared: pair("الـbuffer اتمسح", "The buffer was cleared"),
      },
      parts: ["echo", "-", "orbit"],
    },
  },
  {
    id: "c386.vet-playground",
    kind: "ch38-vet-lenses",
    config: {
      labels: {
        title: pair("اختبار وسكان المصدر", "Tests and source inspection"), tests: pair("سلوك وقت التشغيل", "Runtime behavior"), vet: pair("فحص المصدر", "Source inspection"),
        runTests: pair("شغّل فحص السلوك", "Run behavior check"), runVet: pair("افحص المصدر", "Inspect source"), reset: pair("ابدأ تاني", "Reset"),
        ready: pair("مفيش فحوصات لسه", "No checks have run"), testsDone: pair("السلوك اتجرّب", "Behavior was exercised"), vetDone: pair("المصدر اتفحص", "Source was inspected"), bothDone: pair("الفحصين كل واحد له دور", "Each check adds different evidence"), cleared: pair("الفحوصات اتمسحت", "Check results were cleared"), waiting: pair("لسه", "Waiting"), complete: pair("اتراجع", "Reviewed"),
      },
    },
  },
  {
    id: "c387.profile-playground",
    kind: "ch38-profile-focus",
    config: {
      labels: {
        title: pair("اختار زاوية القياس", "Choose a profiling view"), focus: pair("الملف", "Profile"),
        cpu: pair("CPU", "CPU"), heap: pair("Heap", "Heap"),
        choose: pair("اختار نوع profile", "Choose a profile type"), cpuFocus: pair("وقت التنفيذ فين بيروح", "Where execution time goes"), heapFocus: pair("إيه اللي لسه في الذاكرة", "What remains in memory"),
      },
    },
  },
];

for (const fixture of fixtures) {
  for (const lang of ["ar", "en"]) {
    test(`${fixture.kind} renders and responds in ${lang}`, () => {
      const harness = createSandboxHarness({ chapter: 38, lang });
      const sandbox = harness.render(fixture.kind, fixture.config, { id: fixture.id, text: fixture.config.labels.title[lang] });
      assert.ok(sandbox, `missing renderer for ${fixture.kind}`);
      assert.equal(sandbox.querySelector(".ch38-title").textContent, fixture.config.labels.title[lang]);

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
