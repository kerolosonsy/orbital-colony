---
name: Orbital Colony
description: Learn Go by surveying planet Kepler-G, one geologic map sheet per mission.
colors:
  ink: "#17142a"
  ink-2: "#1d1934"
  sheet: "#221d3b"
  sheet-2: "#2a2447"
  sheet-3: "#332c55"
  rule: "#3e3666"
  rule-strong: "#5a4f8c"
  contact: "#0d0a1c"
  paper: "#efe7d6"
  paper-2: "#cfc6b2"
  muted: "#a69dc2"
  faint: "#9a92bf"
  regolith: "#d39b55"
  soil: "#a45f3e"
  ice: "#9fd2cb"
  rock: "#77709a"
  moss: "#93c45c"
  lichen: "#45b3a8"
  crystal: "#d75c9f"
  survey: "#ef5a40"
  survey-hi: "#ff7458"
  survey-ink: "#1a0d0a"
  star: "#f2b33d"
  warn: "#f0a14a"
  err: "#ff6b6b"
typography:
  display:
    fontFamily: "Readex Pro, Geeza Pro, system-ui, sans-serif"
    fontSize: "30px"
    fontWeight: 700
    lineHeight: 1.35
  headline:
    fontFamily: "Readex Pro, Geeza Pro, system-ui, sans-serif"
    fontSize: "21px"
    fontWeight: 700
    lineHeight: 1.4
  title:
    fontFamily: "Readex Pro, Geeza Pro, system-ui, sans-serif"
    fontSize: "18px"
    fontWeight: 600
    lineHeight: 1.35
  body:
    fontFamily: "Readex Pro, Geeza Pro, system-ui, sans-serif"
    fontSize: "15px"
    fontWeight: 400
    lineHeight: 1.6
  lesson:
    fontFamily: "Readex Pro, Geeza Pro, system-ui, sans-serif"
    fontSize: "15px"
    fontWeight: 400
    lineHeight: 1.85
  label:
    fontFamily: "Barlow Condensed, Readex Pro, sans-serif"
    fontSize: "15px"
    fontWeight: 700
    lineHeight: 1
    letterSpacing: "0.06em"
  numeral:
    fontFamily: "Barlow Condensed, Readex Pro, sans-serif"
    fontSize: "16px"
    fontWeight: 600
    lineHeight: 1
    letterSpacing: "0.04em"
    fontFeature: "tnum"
  sheet-code:
    fontFamily: "Barlow Condensed, Readex Pro, sans-serif"
    fontSize: "12.5px"
    fontWeight: 600
    lineHeight: 1
    letterSpacing: "0.08em"
  code:
    fontFamily: "JetBrains Mono, ui-monospace, Menlo, monospace"
    fontSize: "14px"
    fontWeight: 400
    lineHeight: 1.65
    fontFeature: "liga 0, calt 0"
rounded:
  sheet: "2px"
  inline: "4px"
  tab: "5px"
spacing:
  gutter: "10px"
  panel-x: "12px"
  lesson: "16px 18px"
  titleblock: "0 14px"
components:
  button-tool:
    backgroundColor: "transparent"
    textColor: "{colors.paper-2}"
    rounded: "{rounded.sheet}"
    padding: "0 11px"
    height: "36px"
  button-tool-hover:
    backgroundColor: "{colors.sheet-2}"
    textColor: "{colors.paper}"
  button-run:
    backgroundColor: "{colors.survey}"
    textColor: "{colors.survey-ink}"
    typography: "{typography.label}"
    rounded: "{rounded.sheet}"
    padding: "0 20px"
    height: "38px"
  button-run-hover:
    backgroundColor: "{colors.survey-hi}"
  button-run-busy:
    backgroundColor: "{colors.sheet-3}"
    textColor: "{colors.paper}"
  panel:
    backgroundColor: "{colors.sheet}"
    textColor: "{colors.paper}"
    rounded: "{rounded.sheet}"
  editor:
    backgroundColor: "{colors.ink}"
    textColor: "{colors.paper}"
    typography: "{typography.code}"
    padding: "14px 16px"
  chip:
    backgroundColor: "{colors.sheet-2}"
    textColor: "{colors.paper}"
    typography: "{typography.numeral}"
    rounded: "{rounded.sheet}"
    height: "24px"
    padding: "0 9px 0 7px"
  tab-active:
    backgroundColor: "{colors.ink}"
    textColor: "{colors.paper}"
    rounded: "{rounded.tab}"
    height: "30px"
    padding: "0 12px"
  quad-cell:
    backgroundColor: "{colors.sheet}"
    textColor: "{colors.paper}"
    padding: "9px 12px"
  quad-cell-done:
    textColor: "{colors.ink}"
  console:
    backgroundColor: "{colors.ink}"
    textColor: "{colors.paper}"
    padding: "10px 12px"
---

# Design System: Orbital Colony

## Overview

**Creative North Star: "The Survey Sheet"**

Every mission is one sheet of a planetary geologic survey of Kepler-G, and the player's Go code is the instrument that maps the ground. The whole client is drawn as cartographic print on a deep map-ink ground: a title block across the top, framed work sheets below, and a map with neatline, graticule ticks, collar text, north arrow, scale bar and legend. It rejects the dark-IDE-with-neon-glow look that coding games default to.

Colour is flat and matte, and every hue does a job. The four ground units (regolith, soil, ice, rock) and the map symbols (moss, lichen, crystal, beacon) are the palette. Survey red is the single action colour. Depth comes from rules, not light: hairline borders, double neatlines and contact lines where units meet. Texture is drawn in code (stipple, furrows, sheen, cross-hatch) and never comes from images.

Density is a working tool's density. Three sheets sit side by side on a 10px gutter. Lessons read right to left in Arabic, while code, the map and its numerals read left to right. Numbers are cartographic condensed figures and flip digit by digit when they change. Mission progress is a quadrangle index whose cells fill with their chapter's unit colour.

**Key Characteristics:**
- Map-ink ground, flat unit colours, no glow, no smooth gradients.
- Double-ruled frames: an outer strong rule and an inner hairline, separated by a band of ink.
- Condensed tabular numerals for every count, tick, code and coordinate.
- One action colour (survey red), used only for Run and the next-mission button.
- Symbols keep one scale. A plant's ring fills as it grows and never resizes.
- Arabic text RTL, code and map LTR, both on the same sheet.

## Colors

The palette is a geologic map key on a violet-black ground. Unit colours are used at full strength as fills and hold their meaning everywhere they appear.

### Primary
- **Survey Red** (`survey`): the only action colour. It fills the Run button, the next-mission button and the start button, outlines the current cell in the quadrangle index, and draws the beacon symbol and the SURVEYED stamp. Hover lifts it to **Survey Red Hi** (`survey-hi`), which is also the editor caret. Text on it is **Survey Ink** (`survey-ink`), a near-black red.

### Secondary
- **Regolith Ochre** (`regolith`): the default ground unit and the interface's secondary accent. It covers focus rings, links, list markers, lesson h3 swatches, the goal callout tint, the overall-progress bar, number tokens in code, and inline code text.
- **Star Amber** (`star`): earned stars, hint callouts, console tips and the crystal-harvest label.

### Tertiary (geologic units and map symbols)
- **Soil Rust** (`soil`), **Ice Pale Teal** (`ice`), **Rock Slate Violet** (`rock`): ground units, each paired with its own texture. Ice also carries info roles: the linked-VS-Code state, console reports and package names in code.
- **Moss Green** (`moss`): the plant symbol, the drone energy bar and success states (console ok, ok banner).
- **Lichen Teal** (`lichen`): the second plant symbol, marked with a dot at its centre.
- **Crystal Magenta** (`crystal`): the crystal diamond symbol.
- **Drone palette** (canvas only): `#ff7a52` for the primary drone, then `#f2d16b`, `#8fd3ff`, `#c7a2ff`, `#efe7d6`, `#93c45c`, indexed by drone id.
- **Chapter units** (quadrangle index): chapters cycle through regolith, ice, moss, crystal, soil, lichen, star and `#c7a2ff`.

### Neutral
- **Map Ink** (`ink`): the page ground, the canvas ground, the editor and the console.
- **Ink Raised** (`ink-2`): the title block, the editor panel and the quad-row name column.
- **Sheet** (`sheet`, `sheet-2`, `sheet-3`): panel surface, then hover, then empty track or busy state.
- **Rule** (`rule`) and **Rule Strong** (`rule-strong`): hairlines and outer frames. On the canvas, Rule Strong draws the outer neatline and graticule ticks.
- **Contact Black** (`contact`): canvas only. It draws geologic contact lines, symbol outlines and the swatch keylines in the legend.
- **Paper** (`paper`): primary text, the inner neatline, the north arrow and the hover label plate.
- **Paper 2** (`paper-2`): lesson body text and secondary copy.
- **Muted** (`muted`) and **Faint** (`faint`): sheet codes, graticule labels, read-only tags and unpressed toggles.
- **Warn** (`warn`) and **Error** (`err`): low energy and console warnings, and compile errors, the error banner and the error gutter line.

### Named Rules
**The Unit Means Something Rule.** A unit colour always stands for its unit. Regolith is ochre in the tile, the legend, the inventory chip and the code token, so never reuse a unit hue for an unrelated decoration.

**The One Red Rule.** Survey red marks the action that moves the survey forward: Run, next mission, start, and the current sheet. Toolbar buttons stay ruled and transparent.

## Typography

**Display Font:** Readex Pro (with Geeza Pro, system-ui)
**Label/Numeral Font:** Barlow Condensed (with Readex Pro)
**Mono Font:** JetBrains Mono (with ui-monospace, Menlo)

**Character:** Readex Pro is a round, open Arabic-and-Latin face that carries the friendly Egyptian Arabic voice. Barlow Condensed is the cartographer's lettering, used for map collars, sheet codes, tick labels, counts and the Run label. JetBrains Mono, with ligatures off, keeps code literal.

**Open item:** all three families load from the Google Fonts CDN (`fonts.googleapis.com`) and are not self-hosted yet. Self-hosting is pending the user's permission to download the font files. The fallbacks above are the real stacks in use until then.

### Hierarchy
- **Display** (700, fluid 30–42px, 1.25, −0.01em): the welcome heading, balanced wrap. It drops to 24px on phones. The result title uses 28px at 700.
- **Headline** (700, 21px, 1.4): lesson h2. Overlay sheet heads use 21px at 600.
- **Title** (600, 19px): the mission title in the title block, 16px on phones. The goal line is 17px at 700 with 1.5 line height, in paper.
- **Body** (400, 15px, 1.6): interface copy. **Lesson** body uses 1.85 line height at a 70ch maximum width.
- **Label** (Barlow Condensed 700, 15px, 0.06em, uppercase): the Run button and its siblings, 18px on the big start button.
- **Numeral** (Barlow Condensed 600, 15–16px, tabular): the overall count, the tick counter and inventory chips. Quad numbers use 30px at 700, in the chapter's unit colour.
- **Sheet code** (Barlow Condensed 600, 12–12.5px, 0.06–0.08em): sheet codes in the title block and in index cells.
- **Code** (JetBrains Mono 400, 14px, 1.65): the editor. The console is 12.5px and lesson pre and inline code are 13px.

### Named Rules
**The Condensed Numeral Rule.** Every number a player tracks (ticks, counts, coordinates, sheet codes, scale-bar text) is set in Barlow Condensed with tabular figures. Changing numbers flip digit by digit; they never cross-fade as whole strings.

**The No Ligature Rule.** Code, the editor and the console turn off `liga` and `calt`, so the player sees exactly the characters they typed.

## Layout

The desktop page is one fixed viewport with no page scroll. A 60px title block sits on top. Below it, three sheets sit in a grid of `minmax(290px, 24%) | minmax(0, 1fr) | minmax(400px, 39%)` (lesson, editor, map) on a 10px gutter with 10px outer padding. Each sheet has a 38–42px head strip ruled off with a hairline. The map sheet stacks HUD, map frame, banner, legend, controls and console, and the console flexes to fill what is left.

- **≤1240px:** toolbar labels collapse to icons and the progress bar narrows to 60px.
- **≤1100px (tablet):** the lesson becomes a right-hand drawer, up to 440px wide, that slides in on the lesson button. A goal strip appears under the title block, and the grid becomes editor `1fr` beside a map of `minmax(320px, 42%)`. The map frame is sized by an aspect-ratio variable (`--map-ar`, set per world from the grid size) up to 62vh, so the frame follows the sheet's proportions.
- **≤760px (phone):** a single column with page scroll. The title block wraps and Run stretches to fill its row. The editor is 380px tall and the map sheet 640px. Index rows stack, with two cells per row.

**Open item:** letterboxing of the map at tablet and phone widths is only partly addressed. The aspect-ratio frame fixed the main case, but the 62vh cap and the fixed 640px phone panel can still leave ink bands around the sheet.

The canvas reserves a margin around the map for its furniture: 30px left for y ticks, 46px below for x ticks and the scale bar, 34px above for the collar and north arrow, and 16px right. Cells are capped at 80px.

## Elevation & Depth

The working surfaces are flat. Depth comes from rules and tonal steps (ink, then ink-2, then sheet, then sheet-2), not shadows. Panels use a double-rule frame drawn with inset shadows: a 1px `rule-strong` border, a 3px band of ink, then a 1px `rule` hairline (`box-shadow: inset 0 0 0 3px var(--ink), inset 0 0 0 4px var(--rule)`). The mission title box uses the same device at a smaller scale. On the canvas, the map neatline mirrors it: a 1px paper inner rule and a 2px `rule-strong` outer rule 5px out.

### Shadow Vocabulary
- **Overlay lift** (`box-shadow: 0 30px 80px -20px #000a`): modal sheets (index, result, welcome) only, over an ink scrim at 82% with a 3px backdrop blur.
- **Drawer lift** (`box-shadow: 0 24px 70px -10px #000c`): the tablet lesson drawer only.
- **Pressed** (`inset 0 3px 0 #00000040`, added to the Run keyline): the active state of Run.
- **Drone hover shadow** (canvas, `#0008` ellipse under the drone): the one cast shadow on the map. It shows the drone flying above the ground plane.

### Named Rules
**The Rules Not Light Rule.** Sheets separate by ruled frames and tone, never by glow or blur. Soft shadows belong only to things that float over the work (modal sheets and the drawer).

## Shapes

The form language is square-cornered print. The sheet radius (2px) applies to panels, buttons, chips, callouts and banners. Legend swatches, progress and energy bars, and index cells are fully square (0). Inline code and focus rings soften to 4px, and editor tabs to 5px on their top corners. Round shapes are reserved for map symbols: plant rings, drone bodies and rotors, the beacon pulse, and the legend dots for symbols. Dashed lines mean provisional or tentative: hint callouts, the hints divider, the hovered tile outline and the drone's surveyed path. A tile outline is dashed `4 3`; a path is dashed `6 5`.

**The Square Sheet Rule.** Anything that is a sheet, cell or swatch has square or 2px corners. Rounding belongs to the symbols drawn on the map.

## Components

### Buttons
- **Tool** (ruled, transparent): 36px tall, 2px corners, a 1px `rule` border and `paper-2` text at 13.5px/500 with an 18px stroked icon. On hover the fill turns `sheet-2`, the border turns `rule-strong` and the text turns `paper`. On press it drops 1px. A small variant is 30px tall. An unpressed toggle (sound off) uses `faint`, and the linked VS Code state is outlined and lettered in `ice`.
- **Tool strip:** related tools (sound, VS Code, language, player, reset) share one ruled strip: square inner joints, 2px outer corners, borders collapsed by −1px.
- **Run** (survey red): 40px tall with 22px padding, joined to the Step tool on its trailing side (Step keeps the Run keyline colour as its border), condensed 700 uppercase label at 0.06em, an inset 1px `#9e2f1d` keyline and a filled play icon. Hover is `survey-hi`. While busy it turns `sheet-3` and the icon pulses. A big variant (48px, 18px label) starts the game.
- **Transitions:** 0.2s on `cubic-bezier(.16, 1, .3, 1)` for colour, 0.12s for press.

### Chips
- **Inventory chip:** 24px tall, 2px corners, `sheet-2` fill, a `rule` border, a condensed tabular count and a 9px square swatch in the item's unit colour. When the count changes the chip bumps once (scales to 1.18).

### Cards / Containers
- **Panel (work sheet):** `sheet` fill with the double-rule frame (see Elevation) and 2px corners. The editor panel uses `ink-2`.
- **Goal callout:** regolith mixed 14% into `sheet`, with a regolith-tinted border and a target icon in regolith. It is the first thing on the lesson sheet.
- **Hint callout:** star mixed 10% into `sheet`, with a dashed star-tinted border. It rises in over 0.45s.
- **Banner:** error or ok colour mixed 16% into `ink`, with a tinted border, RTL.
- **Overlay sheet:** `sheet` fill, 2px corners, the double-rule frame (`rule-strong` outside, a band of `sheet`, a `rule` hairline inside) and overlay lift, entering with a 0.4s lift. The welcome sheet is a wide survey window (up to 1120px): the drawn map on one side drifts slowly (scale 1.04 → 1.16 over 28s, alternating) and the copy rises in once, staggered by about 0.08s.

### Inputs / Fields
- **Editor:** a transparent textarea over a highlighted `<pre>`, on `ink`. The caret is `survey-hi`, the selection is `#4a3f7a`, and the line gutter is 46px wide with `#6f679c` numerals. An error line inverts to `err`. Syntax colours: keyword `#ff8f7a`, string `#b6dd7c`, number regolith, comment `#7a72a6` italic, package ice, function `#f2cf7a`, type `#c7a2ff`.
- **Tabs:** mono 12.5px in `muted`. The active tab is filled `ink` with a `rule` border, open at the bottom.
- **Range and checkbox:** native controls with `accent-color` set to regolith.
- **Focus:** a 2px regolith outline, offset 2px, with 4px corners, on every focusable element.

### Navigation
- **Title block:** a 60px `ink-2` bar ruled underneath. In order it holds the brand mark (an orbit glyph stroked in regolith with a survey-red moon) and a condensed wordmark, the index button, the mission title box (double-ruled, title plus sheet code behind a hairline divider), stars, the overall count with its 6px progress bar, then the sound, VS Code and reset tools and Run.

### Map Sheet (signature)
The world is drawn as a geologic unit map. Each tile is a flat unit fill plus its texture:
- **Regolith:** 1.6px stipple at `#00000030`.
- **Soil:** dashed horizontal furrows at `#00000038`.
- **Ice:** fine 45° sheen at `#ffffff55`.
- **Rock:** cross-hatch at `#00000040`.

A faint cell grid (`#00000022`) runs under round-capped `contact` lines. Contact lines are drawn only where two different units meet, and their width scales with the cell (at least 1.5px). Around the map sit:
- the double neatline;
- a graticule of x ticks along the foot and y ticks up the side, labelled in condensed 600 `muted`;
- a collar line, "KEPLER-G · SURVEY SHEET {code} · {w}×{h}", in condensed 700 `paper-2`;
- a north arrow;
- a three-cell alternating scale bar labelled "3 TILES = 3 MOVES = 3 TICKS".

Symbols:
- **Beacon:** a survey-red triangle with a paper centre dot and a fading pulse ring.
- **Crystal:** a magenta diamond with a highlight seam. Its size is printed in the cell corner.
- **Plants:** a fixed-size ring that fills as a pie with growth, plus a second, breathing ring when ripe. Lichen adds a centre dot.
- **Drones:** four paper rotors around a body in the drone colour. The primary drone has a centre dot; other drones carry their id.

The hovered tile gets a dashed paper outline and a paper label plate with ink text. Floating event labels are condensed 700 text with a 4px ink stroke. A dust storm is an ochre wash with moving hatch lines, clipped to the neatline.

### Quadrangle Index (signature)
Chapters are rows of sheets tiled edge to edge (rows overlap by −1px so rules never double). Each row has a 210px name column on `ink-2` with a large condensed number in the chapter's unit colour, followed by three cells. A locked cell is hatched at 135°. A solved cell fills with its chapter unit plus a 6px stipple, with `ink` text. The current cell carries a 2px survey-red inset outline.

### Field Notes Console
Mono 12.5px on `ink`. Each line has a 4px square bullet in its own colour: output paper, error err, tip star, info muted, report ice, ok moss (600), warn warn.

### Result Stamp
A survey-red rubber stamp ("SURVEYED" plus the sheet code, condensed 700 17px, 3px border), rotated −9° and overlapping the result sheet's corner. It stamps in from 2.2× scale and the sheet takes the hit: a 4px drop and −0.5° tilt at 0.5s. Stars (62px) pop in staggered by 0.18s.

## Do's and Don'ts

### Do:
- **Do** paint ground as flat unit colour plus its code-drawn texture (stipple, furrow, sheen, hatch). Draw texture with hard-stop gradient functions or canvas strokes.
- **Do** frame every work surface with the double rule: `rule-strong` outside, a band of ink, then a `rule` hairline inside.
- **Do** set every tracked number in Barlow Condensed with tabular figures, and flip it digit by digit when it changes.
- **Do** keep Survey Red for the action that advances the survey and for the current sheet.
- **Do** keep map symbols at one scale through every state, and show change by filling, never by resizing.
- **Do** isolate LTR code, numerals and the map inside RTL Arabic copy (`direction: ltr`, `unicode-bidi: isolate`).
- **Do** honour `prefers-reduced-motion`; all animation and transition durations collapse to 0.01ms.

### Don't:
- **Don't** use glow, neon or blurred light on working surfaces; that is the dark-IDE default this world refuses.
- **Don't** use smooth colour-to-colour gradients. Gradient functions are allowed only as hard-stop textures.
- **Don't** give a unit colour a second, unrelated meaning.
- **Don't** fill toolbar buttons with survey red or any unit colour; tools stay transparent and ruled.
- **Don't** round sheets, cells or swatches beyond 2px.
- **Don't** add raster or image textures to the map; every pattern is drawn in code.
