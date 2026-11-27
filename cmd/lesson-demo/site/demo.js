"use strict";

const directions = Object.freeze({ North: [0, 1], East: [1, 0], South: [0, -1], West: [-1, 0] });
const maxCommands = 12;

function runRoute(mission, commands) {
  let { X, Y } = mission.Start;
  let blocked = 0;
  for (const direction of commands) {
    if (!Object.hasOwn(directions, direction)) throw new Error("Unknown direction");
    const [dx, dy] = directions[direction];
    const x = X + dx;
    const y = Y + dy;
    if (x < 0 || y < 0 || x >= mission.Width || y >= mission.Height) {
      blocked++;
      continue;
    }
    X = x;
    Y = y;
  }
  const success = X === mission.Beacon.X && Y === mission.Beacon.Y;
  return { X, Y, blocked, success, optimal: success && commands.length <= mission.Optimal };
}

function goSource(commands) {
  const calls = commands.map((direction) => "    colony.Move(colony." + direction + ")");
  return ["package main", "", 'import "orbital/colony"', "", "func main() {", ...calls, "}"].join("\n");
}

function mountDemo(mission, document) {
  const select = (id) => document.getElementById(id);
  const moveButtons = [...document.querySelectorAll("[data-move]")];
  const map = select("map");
  let arabic = false;
  const commands = ["East"];
  let result;
  const t = (en, ar) => arabic ? ar : en;

  map.style.setProperty("--columns", String(mission.Width));
  const cells = [];
  for (let y = mission.Height - 1; y >= 0; y--) {
    for (let x = 0; x < mission.Width; x++) {
      const cell = document.createElement("span");
      cell.className = "cell";
      map.append(cell);
      cells.push({ cell, x, y });
    }
  }

  function render() {
    const locale = arabic ? mission.Arabic : mission.English;
    document.documentElement.lang = arabic ? "ar" : "en";
    document.documentElement.dir = arabic ? "rtl" : "ltr";
    document.title = locale.Title + " · Orbital Colony";
    const language = select("language");
    language.lang = arabic ? "en" : "ar";
    language.textContent = t("العربية", "English");
    const labels = {
      "mission-title": locale.Title,
      goal: locale.Goal,
      "map-label": t("Drone map", "خريطة الدرون"),
      "drone-label": t("Drone", "الدرون"),
      "beacon-label": t("Beacon", "المنارة"),
      "editor-label": t("Your Go commands", "أوامر Go بتاعتك"),
      instructions: t("Add direction commands, then run them in order.", "ضيف أوامر الحركة وبعدين شغّلها بالترتيب."),
      run: t("Run commands", "شغّل الأوامر"),
      undo: t("Undo", "تراجع"),
      reset: t("Reset", "ابدأ من الأول"),
      scope: t("This browser lesson simulates direction commands. Open the game locally to write and run Go.", "الدرس ده بيحاكي أوامر الحركة في المتصفح. شغّل اللعبة على جهازك عشان تكتب وتشغّل كود Go."),
      "full-game": t("Full game and setup", "اللعبة الكاملة وطريقة التشغيل"),
    };
    for (const [id, label] of Object.entries(labels)) select(id).textContent = label;

    select("source").textContent = goSource(commands);
    select("run").disabled = commands.length === 0;
    select("undo").disabled = commands.length === 0;
    for (const button of moveButtons) {
      button.disabled = commands.length >= maxCommands;
      button.setAttribute("aria-label", t(button.dataset.move, ({
        North: "شمال", East: "شرق", South: "جنوب", West: "غرب",
      })[button.dataset.move]));
    }
    const { X, Y } = result || mission.Start;
    const target = mission.Beacon;
    for (const { cell, x, y } of cells) {
      const isTarget = x === target.X && y === target.Y;
      const hasDrone = x === X && y === Y;
      cell.classList.toggle("target", isTarget);
      cell.replaceChildren();
      for (const [show, kind] of [[isTarget, "beacon"], [hasDrone, "drone"]]) {
        if (show) {
          const icon = document.createElement("i");
          icon.className = kind;
          cell.append(icon);
        }
      }
    }

    let feedback = t("Add a move to continue the starter.", "ضيف حركة عشان تكمّل الكود.");
    if (commands.length >= maxCommands) feedback = t("Command list full. Run it or undo a move.", "الأوامر كتير. شغّلها أو ارجع حركة.");
    if (result?.success) {
      feedback = result.optimal
        ? t("Signal received in the fewest steps!", "وصلت للمنارة بأقل عدد خطوات!")
        : t("Signal received. Can you reach it with fewer steps?", "وصلت للمنارة. تقدر توصل بخطوات أقل؟");
    } else if (result) {
      feedback = result.blocked
        ? t("The map edge blocked a move. Try again.", "حدود الخريطة وقفت حركة. جرّب تاني.")
        : t("The drone stopped before reaching the beacon.", "الدرون وقف قبل ما يوصل للمنارة.");
    }
    select("status").textContent = feedback + " " + t(
      "Drone (" + X + ", " + Y + "). Beacon (" + target.X + ", " + target.Y + ").",
      "الدرون (" + X + "، " + Y + "). المنارة (" + target.X + "، " + target.Y + ").",
    );
  }

  for (const button of moveButtons) button.addEventListener("click", () => {
    if (commands.length === maxCommands) return;
    commands.push(button.dataset.move);
    result = undefined;
    render();
  });
  select("undo").addEventListener("click", () => {
    commands.pop();
    result = undefined;
    render();
  });
  select("reset").addEventListener("click", () => {
    commands.splice(0, commands.length, "East");
    result = undefined;
    render();
  });
  select("run").addEventListener("click", () => {
    result = runRoute(mission, commands);
    render();
  });
  select("language").addEventListener("click", () => {
    arabic = !arabic;
    render();
  });
  render();
}

if (typeof module === "object" && module.exports) module.exports = { runRoute, goSource };
else mountDemo(mission, document);
