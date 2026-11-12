# خطة تطوير المستعمرة المفتوحة — Orbital Colony Open World

> **الغرض:** مواصفات تنفيذ كاملة يستلمها منفّذ (إنسان أو agent) من غير ما يرجع للمحادثة. مكمّلة لـ `PLAN.md` (المنهج) ومش بديل عنه.
> **الحالة:** خطة فقط. مفيش كود اتكتب.
> **النتيجة المستهدفة:** المستعمرة تتحول من لعبة idle (استنّى واشتري) لمشروع Go واحد طويل العمر. التقدم فيها بييجي من كود أحسن، وكل مسار في المنهج (فصول 1–37) ليه نظام مقابل لازم اللاعب يبرمجه.

## 0. التشخيص (ليه التغيير)

| المشكلة | الدليل في الكود |
|---|---|
| الترقيات أرقام مش مهارات | `econ.Catalog`: 213 مستوى، منهم 12 بس بيفتحوا حاجة تتكتب (8 `api.*` و4 `crop.*`). الباقي مضاعفات (`yield.*` 84، `grow.*` 32، `engine` 18...) |
| الـ starter بيكسب لوحده | `levels.Colony.Starter` حلقة snake شغالة للأبد؛ مفيش سبب تعدّلها |
| ثلثين المنهج غايب | `econ.ChapterReward`: رخص لفصول 4، 5، 6، 8، 10، 15 بس. فصول 18–38 كريدت فقط |
| مفيش أهداف ولا مخاطر | `NewColony`: `Check` بترجع `false` دايماً، و`Storm` مش بيتفعّل. الخسارة الوحيدة عفن البوص |
| مفيش feedback على الكود | `web/colony.js: renderStatus`: كريدت/دقيقة بس، من غير مقارنة بين التشغيلات |
| ملف واحد | `colonyRun` بياخد `code` واحد ويكتبه `main.go` |
| مفيش سياسة imports | `colonyRun` مش بينادي `restrictedImports` (المهمات والـ labs بينادوها) |
| الـ state من برنامج اللاعب | `colonyGame.event`: `c.state = event.State` من fd 3. اللاعب يقدر يكتب على fd 3 بنفسه |

## 1. المبادئ (تحكم أي قرار مش مكتوب هنا)

1. **التقدم = كود + تكلفة.** أي ميزة جديدة محتاجة الاتنين: عقد (تحدي كود) يفتح المخطط، وبعدين موارد حقيقية + كريدت عشان تتبني. مفيش حاجة مجانية، ومفيش حاجة بالكريدت لوحده بعد أول 3 فصول.
1. **دايماً فيه تحسين.** مفيش ترقية ليها حد أقصى، والتكلفة بتطلع لموارد من سلاسل أعمق. مفيش نقطة «خلّصت اللعبة».
2. **نفس قواعد `PLAN.md`:** Go حقيقي، لغة الدروس عربي مصري + `En` كامل، اتبع أسلوب الكود الموجود، مفيش refactor مش مطلوب، كل مرحلة تخلص بـ `go test ./...` أخضر في `cmd/orbital` و`colony`، و`node --test web/` للواجهة.
3. **قاعدة التسلسل:** العقد يحتاج بس مفاهيم اتشرحت في فصول اللاعب خلّصها. العقد مقفول لحد ما فصله يخلص.
4. **الحفظ القديم ميتكسرش:** `State.Version` يترفع، و`Normalize` يرحّل. اللاعب ميخسرش حاجة اشتراها.
5. **المستعمرة مش مهمة:** مفيش starter لحل العقد، ومفيش تلميحات كود. فيه brief ورابط للدرس.

## 2. خريطة الكود الحالي

| ملف | الدور | هيتغير في |
|---|---|---|
| `colony/econ/econ.go` | `State`، `Catalog`، السوق، `Buy`، `ChapterReward`/`Claim` | C1، C2، C6 |
| `colony/econ/layers.go` | الطبقات تحت الأرض وترقيات الحفر | — |
| `colony/internal/sim/colonymode.go` | `NewColony`، `colonyPlant/Harvest`، `Export`، `colonyTicker`، `colonyCommands` (fd 4)، `econSnapshot` (fd 3) | C1، C2، C3، C5 |
| `colony/internal/sim/extras.go` | `Storm`، `stormEvent`، `runStorm` (للمهمات) | C2 |
| `colony/internal/sim/sim.go:419` | `act()`: كل أوامر الدرون بتعدّي منها | C2، C3 |
| `colony/colony.go` | API اللاعب، `startColony` | C1، C2، C5 |
| `colony/levels/levels.go:198` | `levels.Colony` (الدرس والـ starter) | C1، C4 |
| `cmd/orbital/colony.go` | `colonyGame`، endpoints: status/run/stop/buy/claim/stream | الكل |
| `cmd/orbital/features.go` | `scan()`، `checkRequirements(code, reqs, lang)`، مفاتيح `requirementTexts` | C1 (استخدام بس) |
| `cmd/orbital/policy.go` | `restrictedImports(code, lvl, lang)` | C4 |
| `cmd/orbital/workspace.go` | فولدر VS Code للمهمة | C4 |
| `web/colony.js` | اللوحة: كريدت، سوق، متجر، مكافآت | الكل |

مفاتيح `Requires` الموجودة: `for for2 if switch func errcheck slice range map struct method pointer ptrrecv nilcheck stringer typeswitch go waitgroup chan select mutex generic generictype benchmark tabletest rune iota const defer recover errwrap errtype funcvalue closure embed variadic bufchan close recursion fuzz subtest iface` + `sel:pkg.Name`.

## 3. المرحلة C1 — لوحة العقود

### C1.1 البيانات (`colony/econ/contracts.go` جديد)

```go
type Contract struct {
    ID       string   // "c05-careful-harvest"
    Chapter  int      // لازم يكون مكتمل عشان يتفتح
    Title, TitleEn, Brief, BriefEn string
    Requires []string // مفاتيح features.go؛ بتتفحص على كود اللاعب قبل الـ build
    After    []string // عقود لازم تخلص قبله
    Window   int64    // ticks من أول Run؛ 0 = من غير وقت
    Goal     Goal
    Limits   Limits
    Reward   Reward   // النوع الموجود: Credits + Grants
}
type Goal struct {
    Deliver  map[string]float64 // كميات تتباع بـ Export جوه الـ Window
    Earn     float64            // أو كريدت مكتسبة
    Submit   string             // أو مفتاح colony.Submit (C5)
    MaxLost  int                // أقصى نباتات اتدمرت/عفّنت
    MinDrones int               // عدد درونات لازم يشتغلوا فعلاً (≥ N أمر لكل واحد)
}
type Limits struct {
    MaxDrones int      // 0 = من غير حد
    Deny      []string // ترقيات api.* بتتقفل مؤقتاً
}
type ContractRun struct {
    Status  string  // "", "active", "done", "failed"
    Start   int64   // tick
    Sold    map[string]float64
    Earned  float64
    Lost    int
    Best    int64   // أقل ticks خلص فيها
    Tries   int
}
```

`State` يزيد: `Contracts map[string]*ContractRun`، `Active string`. `Version` = 2.

دوال pure (كلها بتستات في `contracts_test.go`): `Contracts []*Contract`، `FindContract(id)`، `(s) Available(done map[int]bool) []*Contract`، `(s) Accept(id) error`، `(s) Progress(tick) (done, failed bool, why string)`، `(s) Settle(tick) (Reward, error)`.

### C1.2 التدفق

1. اللاعب يدوس «اقبل» والمستعمرة واقفة → `POST /api/colony/contract {id}` → `Accept` (عقد واحد نشط بس).
2. `colonyRun`: لو فيه عقد نشط، `checkRequirements(code, c.Requires, lang)` قبل الـ build. لو ناقص → `{"phase":"requires","missing":[...]}` بنفس نصوص المهمات، ومفيش تشغيل.
3. الـ sim بيعدّ: `Export` يزوّد `Sold`/`Earned`، تدمير نبات غير ناضج وعفن البوص يزوّدوا `Lost`، وعدّاد أوامر لكل درون.
4. `colonyTicker` بينادي `Progress` كل ثانية: نجاح → `Settle` + event `{"op":"contract","ok":true}`؛ انتهاء `Window` → `failed` + event. الفشل ميخسّرش حاجة غير المحاولة.
5. `POST /api/colony/contract {id:""}` يلغي العقد النشط.

### C1.3 إعادة التوازن

- رخص `api.*` و`crop.*` بقت **مخططات**: العقد بيحط المخطط في `State.Blueprints`، والشراء بعد كده بموارد + كريدت (قسم 3A.2). `Status()` يرفض الشراء من غير مخطط برسالة «خلّص عقد كذا الأول». `ChapterReward` يفضل يدفع كريدت بس.
- المضاعفات (`yield.*`، `grow.*`، `engine`...) من غير حد أقصى وبتكلفة موارد (قسم 3A.3).
- **ترحيل v1 → v2 في `Normalize`:** أي مستوى مملوك يفضل زي ما هو، وأي `api.*`/`crop.*` مملوكة تاخد مخططها تلقائي. تست: save v1 فيه `yield.moss: 12` و`api.radar: 1` يفضلوا بعد التحميل.

### C1.4 العقود (فصول 2–17)

| ID | فصل | المطلوب | Requires | يفتح |
|---|---|---|---|---|
| `c02-price-watch` | 2 | اكسب 150 كريدت والبيع كله بسعر ≥ 90% من الأساسي (يعني لازم تقرا السعر وتستنى) | `if` | `api.prices` |
| `c03-full-sweep` | 3 | سلّم 60 moss في 900 tick على أرض 5×5 | `for2` | `crop.lichen` |
| `c04-right-ground` | 4 | سلّم 30 moss و30 lichen، `MaxLost: 0` | `switch` | `api.scanner` |
| `c05-careful-harvest` | 5 | سلّم 80 من أي محصول، `MaxLost: 0`، `Deny: api.scanner` (يعني لازم `TryHarvest` وتعالج الخطأ) | `errcheck`,`func` | `api.try` |
| `c06-ripe-ledger` | 6 | سلّم 120 moss في 700 tick (محتاج تسجّل وقت النضج لكل خانة بدل ما تلف على الكل) | `map`,`slice` | `crop.crystal`,`api.measure` |
| `c07-plot-type` | 7 | سلّم 40 glowcap بمتوسط حصاد ≥ 3 للفطرة (عناقيد) | `struct`,`method` | `crop.glowcap` |
| `c08-big-veins` | 8 | سلّم 60 crystal من عروق حجم ≥ 4 بس | `pointer`,`nilcheck` | `api.radar` |
| `c09-crop-protocol` | 9 | سلّم 3 محاصيل مختلفة، 40 من كل واحد، في عقد واحد | `iface` | `crop.reed` |
| `c10-fleet` | 10 | سلّم 200 moss في 600 tick، `MinDrones: 3` | `go`,`waitgroup` | `api.spawn`,`hangar` |
| `c11-dispatch` | 11 | سلّم 100 crystal، `MinDrones: 3`، درون واحد بس مسموحله `Export` | `chan`,`select` | `hangar` |
| `c12-one-ledger` | 12 | سلّم 300 moss، `MinDrones: 4`، والتشغيل بـ `-race` نضيف (C3.3) | `mutex` | `market.logistics` مستوى |
| `c13-reed-clock` | 13 | سلّم 50 reed، `MaxLost: 2` | `sel:context.WithTimeout` | `market.contracts` مستوى |
| `c14-generic-queue` | 14 | سلّم 150 من محصولين، والكود فيه نوع generic بيتعرّف | `generictype` | 500 كريدت |
| `c15-shortest-run` | 15 | سلّم 80 crystal بمسافة كلية ≤ حد (عدّاد حركة من C3.1) | `recursion` أو `sel:heap.Push` | `api.swap` |
| `c17-deep-haul` | 17 | سلّم 20 geode | `for2`,`slice` | `api.dig`,`drill` |

الأرقام مبدئية: المنفّذ يظبطها بحيث الحل المرجعي يخلص في 40–70% من الـ `Window`، والـ starter ميخلصش.

### C1.5 الواجهة والتحقق

- قسم «العقود» فوق المتجر في `#col-panel`: المتاح، النشط (شريط تقدم + وقت باقي)، المكتمل (أحسن وقت)، والمقفول («خلّص الفصل N»).
- عقد نشط يظهر في `#goal-strip-text` بدل الهدف العام.
- **تم لما:** لكل عقد حل مرجعي في `cmd/orbital/testdata/colony/<id>/` + تست `TestContractsSolvable` بيشغّل الـ sim بساعة سريعة (`Unit` = 1ms) ويتأكد إن الحل بيخلص والـ starter لأ؛ تست i18n للعقود؛ تست الترحيل v1→v2.

## 3A. الاقتصاد الجديد (المراحل E1–E5)

قرارات صاحب المشروع (نهائية): الترقيات بموارد حقيقية + كريدت مع بعض؛ سلاسل تصنيع وخامات جديدة؛ عقد يفتح وموارد تبني؛ مستويات لا نهائية؛ ومحتوى كتير. والتحسين المستمر من أربع مصادر: صيانة واستهلاك، حدود تخزين وطاقة، أهداف كفاءة بتكبر، وبحث علمي.

### 3A.1 الموارد (E1)

**خامات** (الـ 7 الموجودين + 7 جداد):

| خامة | بتيجي منين | أمر | تتفتح |
|---|---|---|---|
| moss، lichen، glowcap، reed | زراعة (موجود) | `Plant`/`Harvest` | موجود |
| crystal | عروق على الصخر (موجود) | `Harvest` | موجود |
| iron، geode | الطبقات 1+ و2+ (موجود) | `Dig` | موجود |
| water | تسييح جليد؛ الخانة بتتحول regolith وبترجع جليد بعد 600 tick | `Melt()` | فصل 4 |
| silica | غربلة regolith؛ كل خانة ليها مخزون بيخلص وبيتجدد ببطء | `Sift()` | فصل 6 |
| copper | الطبقة 2+ | `Dig` | فصل 17 |
| sulfur | الطبقة 3+، جنب فوّهات بتأذي الدرون لو وقف عليها أكتر من 20 tick | `Dig` | فصل 17 |
| titanium | الطبقة 4+ | `Dig` | فصل 17 |
| helium | جامع غاز (مبنى) على السطح، بينتج بالليل بس | مبنى | فصل 19 |
| uranium | الطبقة 5، محتاج درون بدرع | `Dig` | فصل 25 |

توليد الخامات تحت الأرض بنفس `undergroundAt` (hash بالإحداثيات)، فالحفظ القديم بيكتشفها في الخانات اللي لسه متحفرتش بس.

**منتجات مصنّعة** (`econ.Recipes`):

| درجة | منتج | المدخلات | مبنى | ticks |
|---|---|---|---|---|
| T1 | biomass | 4 moss أو 2 lichen | Composter | 40 |
| T1 | plate | 2 iron | Smelter | 60 |
| T1 | wire | 2 copper | Smelter | 60 |
| T1 | glass | 3 silica | Kiln | 80 |
| T1 | fuel | 2 biomass + 1 water | Refinery | 100 |
| T2 | fertilizer | 2 biomass + 1 sulfur | Refinery | 120 |
| T2 | circuit | 2 wire + 1 glass + 1 crystal | Assembler | 200 |
| T2 | alloy | 2 plate + 1 titanium | Smelter | 180 |
| T2 | battery | 1 wire + 1 sulfur + 1 plate | Assembler | 200 |
| T2 | coolant | 2 water + 1 helium | Refinery | 160 |
| T3 | processor | 3 circuit + 1 geode | Assembler | 500 |
| T3 | frame | 2 alloy + 1 circuit + 1 battery | Assembler | 500 |
| T3 | cell | 1 uranium + 1 alloy + 2 coolant | Reactor Lab | 800 |
| T4 | core | 2 processor + 2 frame + 1 cell | Assembler | 2000 |

كل المنتجات بتتباع في السوق (السعر ≈ 1.6× قيمة المدخلات) عشان البيع يفضل اختيار مقابل الصرف على الترقيات.

### 3A.2 المباني واللوجستيات (E2)

- المبنى بيتبني على خانة ويشغلها (`Tile.B`)، فبينافس أرض الزراعة. فيه 9 أنواع: Composter، Smelter، Kiln، Refinery، Assembler، Reactor Lab، Silo، Generator (3 أنواع في E3)، Lab.
- كل مبنى ليه مخزن دخل ومخزن خرج بسعة (8 في المستوى 0). مليان = واقف. **الدرونات هي اللي بتنقل**، وده مصدر تحديات التفكير: كام درون، أنهي مسار، مين يتصنّع الأول.
- API جديد (`colony/build.go`):

```go
func Build(k BuildingKind) error                 // على خانة الدرون؛ بيخصم التكلفة
func Here() *Building                            // nil لو مفيش مبنى
func Buildings() []*Building
func (b *Building) Kind() BuildingKind
func (b *Building) Pos() Point
func (b *Building) SetRecipe(out Entity) error
func (b *Building) Put(e Entity, n int) (int, error)   // الدرون لازم يكون على الخانة
func (b *Building) Take(e Entity, n int) (int, error)
func (b *Building) Status() BuildingStatus       // In، Out، Progress، Wear، Powered، Level
func Carry() map[Entity]int                      // cargo الدرون
```

- الدرون ليه cargo بسعة (10 في الأول) من أول ما أول مبنى يتبني. قبل كده الحصاد بيروح المخزن مباشرة زي دلوقتي، عشان الفصول 1–6 متتعقدش.
- الأخطاء typed: `ErrNoInput`، `ErrFull`، `ErrNotHere`، `ErrNoPower`، `ErrWorn`.
- كل مبنى بيتفتح بعقد (مخطط) وبيتبني بموارد، وليه مستويات لا نهائية (سرعة، سعة) بنفس معادلة 3A.3.

### 3A.3 معادلة التكلفة والمستويات اللا نهائية (E1)

`Upgrade.Max` يتشال (0 = من غير حد) لكل المضاعفات والمباني. `Upgrade.cost` يتبدل بوصف بيانات:

```go
type CostRule struct {
    Credits, Growth float64   // credits = Credits * Growth^level
    Tiers   [][]string        // Tiers[t] = الموارد المطلوبة في الشريحة t
    Qty, QtyGrowth float64
}
```

- الشريحة `t = level / 5`. كل 5 مستويات المورد المطلوب بيطلع درجة: خامة ← T1 ← T2 ← T3 ← T4. بعد آخر شريحة: نفس موارد T4 والكمية بتكبر.
- مثال `yield.moss`: مستويات 0–4 بـ moss، 5–9 بـ biomass، 10–14 بـ fertilizer، 15–19 بـ processor، 20+ بـ core.
- **عائد متناقص** عشان الأرقام متنفجرش: `mult = 1 + a * level^0.75`. `engine`: `ms = 260 / (1 + 0.09 * level^0.7)` بحد أدنى 20ms.
- المستويات 0–2 من أي ترقية أساسية (moss، engine، field) بكريدت + خامة واحدة بس، عشان لاعب الفصول 1–3 يقدر يتقدم.
- تست: التكلفة monotonic لأول 500 مستوى ومفيش `Inf`/`NaN`؛ كل مورد مطلوب ليه طريق إنتاج متفتح قبل الشريحة اللي بتطلبه.

### 3A.4 الطاقة والتخزين والصيانة (E3)

| نظام | القاعدة | اللي بيتبرمج |
|---|---|---|
| طاقة | كل مبنى شغال بيسحب kW. المولدات: Solar (نهار بس، دورة 1200 tick)، Burner (بيحرق fuel)، Reactor (بيحرق cell). لو السحب > الإنتاج، المباني بتقف بالأولوية | `colony.Power() (supply, demand float64)`، `b.SetPriority(n)`، `b.Enable(bool)`. الكود لازم يطفّي مباني بالليل أو يخزّن batteries |
| تخزين | Silo لكل فئة (محاصيل/خامات/مصنّعات) بسعة. الزيادة بتترفض (`ErrFull`) ومش بتتحط | الكود لازم يبيع أو يصنّع قبل ما يتملي |
| وقود الدرون | بعد أول Refinery: كل أمر بيسحب شحنة. `Charge()` عند القاعدة بتحرق fuel. درون شحنته صفر بيمشي بنص السرعة | `d.Battery() int`، جدولة الشحن |
| تآكل المباني | `Wear` بيزيد مع كل دورة إنتاج. فوق 80 السرعة بتقل، وعند 100 بيقف | `b.Repair()` بتاخد plate؛ الكود يعمل صيانة وقائية |
| مصاريف ثابتة | كل درون وكل مبنى ليه إيجار كريدت/دقيقة (صغير). كريدت سالب = المباني بتقف، مش خسارة | يخلي «اشتري كل حاجة» قرار غلط |

كل نظام بيبدأ بس بعد ما عقده يخلص، وبيتشرح في الـ brief بتاع العقد ده.

### 3A.5 البحث وشجرة التقنيات (E4)

- **نقط البحث (RP)** من مصدرين: مبنى Lab بيحرق منتجات (كل درجة ليها قيمة RP)، و**تجارب**: تحديات كود قصيرة متكررة بـ seed جديد كل مرة (نفس آلية العقود، `Window` قصير، ومكافأة RP بتقل 20% مع كل تكرار في نفس اليوم).
- **عقدة في الشجرة** = `{ID, Branch, After []string, Contract string, RP float64, Cost CostRule, Effect}`. الفتح محتاج التلاتة: العقد، الـ RP، وتكلفة البناء.
- الفروع على مسارات المنهج، والمنفّذ يكتب العقد بنفس تعريف «تم» في قسم 9:

| فرع | فصول | عدد العقد | أمثلة |
|---|---|---|---|
| الزراعة | 1–6 | 16 | ري بالـ water (+نمو للجيران)، تسميد، دورة محاصيل، محصول هجين |
| التعدين | 8، 15، 17 | 14 | رادار عميق، حفار متعدد، دعامات، درع uranium |
| التصنيع | 7، 9 | 16 | كل مبنى + مخزن أكبر + recipe بديلة أرخص |
| الأسطول | 10–13، 24–25 | 18 | درونات متخصصة (ناقل/حفّار/مزارع)، cargo أكبر، محطة شحن تانية، طابور مهام مشترك |
| الطاقة | 19، 25 | 10 | بطاريات شبكة، reactor، توزيع ذكي |
| البيانات | 18–23 | 12 | ledger، تقارير، تشفير شحنات |
| الخدمات | 26–30 | 14 | السوق HTTP، مركز القيادة، المخزن SQL (قسم 7) |
| الشبكات | 31–33 | 10 | الجيران، طرق تجارة، بوابة |
| الأنظمة | 34–37 | 12 | مراقبة، نسخ احتياطي، مستعمرة تانية |

المجموع ≈ 122 عقدة مختلفة، وكل واحدة ليها مستويات لا نهائية. الأنظمة في قسم 7 (C5) بقت عقد في الفروع دي بدل ما تتفتح بالفصل لوحده.

### 3A.6 رتبة المستعمرة: أهداف كفاءة بتكبر (E5)

- رتبة `N` من غير حد. كل رتبة ليها هدف يتحقق على نافذة 3000 tick متواصلة والـ sim بيقيسه.
- الرتب الفردية: **إنتاج** (قيمة المنتجات/دقيقة) والهدف بيتضرب في 1.5 كل مرة.
- الرتب الزوجية: **كفاءة** = قيمة المنتجات بالأسعار الأساسية ومن غير مضاعفات `yield` ÷ عدد أوامر الدرونات. الترقيات مبتأثرش فيها، فمتتعدّاش غير بكود أذكى (مسارات أقصر، انتظار أقل، نقل أقل). الهدف بيتضرب في 1.15.
- الرتبة بتقفل: شرائح التكلفة (الشريحة `t` محتاجة رتبة `2t`)، حجم الأرض، وعدد الدرونات. يعني الفلوس لوحدها مش كفاية.
- الواجهة: شريط الرتبة فوق اللوحة بالهدف الحالي والقياس الحي (من عدّادات C3.1).

### 3A.7 الحالة والتحقق

`State` يزيد: `Blueprints map[string]bool`، `Buildings []Building`، `RP float64`، `Tech map[string]int`، `Rank int`، `Power`، `Silos`. `Goods` و`BasePrice` ياخدوا الموارد الجديدة. `Version` = 3 (بعد C1).

**تم لما (لكل مرحلة E):** تستات pure في `econ` لكل recipe ولمعادلة التكلفة؛ تست sim لكل مبنى (دخل → خرج بالوقت الصح، وقوف عند الامتلاء/انقطاع الطاقة/التآكل)؛ حل مرجعي بيوصل رتبة 4 بالساعة السريعة؛ تست ترحيل v2 → v3؛ `En` كامل لكل اسم ووصف.

## 4. المرحلة C2 — عالم حي

### C2.1 مجدول الأحداث (`colony/internal/sim/colonyevents.go` جديد)

حتمي من `st.Seed` ورقم الـ tick (نفس أسلوب `groundAt`)، بيتنادى من `colonyTicker`. الأحداث بتبدأ بعد ما عقد فصلها يخلص، عشان اللاعب ميقابلش حاجة متشرحتش.

| حدث | يبدأ بعد | السلوك | API |
|---|---|---|---|
| عاصفة | `c13-reed-clock` | تحذير 60 tick قبلها على `colony.StormWarning()` (موجودة). وقت الضربة: أي درون برّه `(0,0)` بيتعطّل 150 tick، والنباتات غير الناضجة على السطح بتتدمر. **مفيش `finish("fail")`** زي المهمات | `StormWarning()`، `StormClear()` |
| عطل درون | `c05-careful-harvest` | احتمال 1% لكل أمر: الأمر يفشل والدرون يفضل معطّل لحد `Repair()` | `colony.Fault() error`، `colony.Repair() bool` (+ نسخ `*Drone`) |
| صدمة سوق | `c02-price-watch` | `Demand[g]` لسلعة واحدة ينزل لـ 0.4 أو يطلع لـ 1.6 لمدة 300 tick | `Price` الموجودة |
| طلب عاجل | `c06-ripe-ledger` | عقد صغير تلقائي: سلعة × كمية × 400 tick بسعر مضاعف، بيتقبل تلقائي | `colony.Rush() (Entity, int, bool)` |

`stormEvent` الحالية بتنادي `finish("fail")`؛ اعمل فرع `w.Econ != nil` بدل ما تغيّر سلوك المهمات. تست يتأكد إن مهمات فصل 13 لسه زي ما هي.

### C2.2 الإغلاق المرتب

- دلوقتي `RunColony` بيمسك SIGTERM ويحفظ ويخرج لوحده. الجديد: `colony.Shutdown() <-chan struct{}`. لو اللاعب ناداها مرة، SIGTERM بيقفل القناة ويدّي البرنامج 2 ثانية (السيرفر بيستنى 3 قبل الـ kill في `colonyGame.stop`) قبل الحفظ.
- بعد عقد `c13`، الدرونات ليها `Capacity` (حقل موجود في `World`): المحصول في الـ cargo لحد `Unload()` عند القاعدة. اللي في الـ cargo وقت الخروج بيضيع. ده اللي بيدّي للإغلاق المرتب معنى.
- عقد `c13b-clean-stop` (فصل 13): السيرفر بيبعت SIGTERM في وقت عشوائي جوه العقد؛ النجاح = صفر cargo ضايع. `Requires: select`.

**تم لما:** تست لكل حدث (حتمي بنفس الـ seed)، تست إن برنامج بيتجاهل العاصفة بيخسر والبرنامج المرجعي لأ، وتست الـ SIGTERM في `process_unix_test.go`.

## 5. المرحلة C3 — قياس الكود

### C3.1 العدّادات (في `act()`)

لكل درون: عدد الأوامر لكل `op`، الأوامر الفاشلة، ticks في `Wait`، مسافة الحركة، حركات من غير شغل بعدها. بتتبعت في `econSnapshot` الكامل (كل 3 ثواني) كـ `stats`.

### C3.2 سجل التشغيلات

`saves/profiles/<id>/colony-runs.json`: آخر 50 تشغيل `{codeHash (sha256 للملفات), started, ticks, earned, perMin, stats, contract}`. بيتكتب في goroutine الخروج جوه `colonyGame.start`. `GET /api/colony/runs`.

الواجهة: كارت «آخر تشغيل» بعد كل Stop: كريدت/دقيقة مقابل التشغيل اللي قبله (بالنسبة المئوية)، ونسبة وقت كل درون (شغل / حركة / انتظار / فشل) كـ bar، وأكتر أمر فشل.

### C3.3 وضع `-race`

toggle «افحص السباقات» جنب Run: الـ build بـ `-race`. سطر `WARNING: DATA RACE` في stderr → event `race` بنفس tip الـ labs، والتشغيل بيتعلّم «فيه سباق». العقود اللي عليها شرط race (`c12`) بتفرض الـ toggle.

**تم لما:** تست إن العدّادات صح على سيناريو ثابت، تست إن السجل بيتقص على 50، تست إن برنامج فيه race بيتمسك.

## 6. المرحلة C4 — مشروع حقيقي

- `POST /api/colony/run` ياخد `files map[string]string` (والـ `code` القديم يفضل شغال). نفس قواعد `LabFile.Name`: `filepath.IsLocal`، امتدادات `.go .sql .json .yaml .txt .md .html`، packages فرعية مسموحة. حد 512KB للكل.
- `restrictedImports` على كل ملف `.go` بسياسة الـ labs (stdlib ما عدا `unsafe`، `os/exec`، `syscall`، `plugin`). `levels.Colony.Deps` فاضية لحد C5.
- `POST /api/colony/test`: `go test ./...` بـ timeout 90s ومن غير تشغيل المستعمرة؛ الناتج في الكونسول بنفس parser أخطاء الـ labs. `POST /api/colony/vet` نفس الفكرة.
- الواجهة: tabs للملفات + «ملف جديد» + زرارين Test و Vet. `workspace/colony/` يتزامن بكل الملفات.
- عقد `c16-tested-planner` (فصل 16): الكود فيه package `plan` بدالة pure، و`go test` بتاع اللاعب لازم يمسك 3 mutants (استخدم آلية `Lab.Mutants` الموجودة).

**الثغرة (state من برنامج اللاعب):** القرار الافتراضي إنها تفضل، زي قرار `PLAN.md` قسم 12 عن الـ hidden tests: اللعبة محلية ولمستخدم واحد، واللي يغش بيغش نفسه. البديل (نقل الـ sim للسيرفر والـ lib تبقى RPC على الـ pipes) تغيير كبير في `sim` كله ومش في الخطة دي.

**تم لما:** مشروع 3 packages بيتبني ويتشغّل ويتختبر من الواجهة ومن VS Code، وتست إن `os/exec` بيترفض.

## 7. المرحلة C5 — أنظمة المسارات

كل نظام: مكوّن في السيرفر (`cmd/orbital/colony_<name>.go`) بيتكلم مع برنامج اللاعب، متغيرات بيئة، وعقود. النظام بيتفتح لما فصوله تخلص. السيرفر بينفّذ أي أثر على الاقتصاد عن طريق قناة الأوامر fd 4 (زي `{"buy":...}` دلوقتي): `{"credit":n,"take":{good:n}}`.

| نظام | فصول | السيرفر بيوفّر | اللاعب بيكتب | عقود |
|---|---|---|---|---|
| **الأرشيف** | 18–23 | فولدر `ORBITAL_COLONY_DIR` فيه `ledger.csv` (كل بيع) و`events.jsonl` (كل حدث)، الـ sim بيكتبهم | يقرا ويحلل ويبلّغ بـ `colony.Submit(key, v)` (جديدة) | `d18-audit`: إجمالي البيع لكل سلعة من الـ ledger (فيه صفوف بايظة). `d19-shift-report`: إنتاج كل ساعة كونية. `d22-atomic-journal`: يحفظ خطته بحفظ ذرّي والسيرفر بيقتله في النص ويتأكد إن الملف سليم |
| **السوق HTTP** | 24، 26–28، 33 | `ORBITAL_MARKET_URL` على loopback: `GET /prices`، `POST /orders` (بيع بالجملة +20%، لازم `Idempotency-Key`)، بيرجّع 429 و503 أحياناً | عميل بـ timeout وretry وbackoff | `h24-bulk-seller`: 20 order من غير تكرار ولا فقد. `h33-stubborn-client`: نفس الشيء والسوق بيقع 30% |
| **مركز القيادة** | 26–28 | السيرفر بيطلب من سيرفر اللاعب على `ORBITAL_LISTEN_ADDR`: `GET /healthz`، `GET /inventory`، `POST /orders` | سيرفر `net/http` بـ handlers حقيقية جنب حلقة الدرونات | `h26-control-api`: كل الطلبات ترد صح في 200ms. `h28-drain`: طلب جاري يخلص وقت الـ SIGTERM |
| **المخزن** | 29–30 | ملف SQLite في `ORBITAL_COLONY_DIR` (`modernc.org/sqlite` من `labdeps`) | schema وترحيل ومخزون في القاعدة | `s29-warehouse`: السيرفر بيعمل restart للبرنامج 3 مرات و`/inventory` لازم يكمّل من غير فقد. `s29-transfer`: تحويل ذرّي بين مخزنين تحت طلبات متوازية |
| **الجيران** | 31–32 | مستعمرة تانية على `ORBITAL_NEIGHBOR_ADDR`: بروتوكول TCP بالسطر (`OFFER`/`ACCEPT`/`SHIP`)، بتبعت سطور طويلة وبتقطع الاتصال أحياناً | عميل TCP بـ deadlines وحدود سطور | `n31-trade-route`: 10 صفقات. `n31-bad-line`: الاتصال يفضل سليم مع مدخلات بايظة |
| **المراقبة** | 36 | السيرفر بيعمل scrape لـ `/metrics` بتاع اللاعب كل 5 ثواني ويرسمها في اللوحة | Prometheus metrics بـ labels محدودة و`/readyz` | `o36-dashboard`: 4 metrics مطلوبة بقيم مطابقة للـ sim (±2%). الأوامر بتتبعت بس و`/readyz` = 200 |

- `levels.Colony.Deps` تاخد المكتبات المطلوبة (sqlite، prometheus) و`go.mod` بتاع تشغيل المستعمرة ياخد نفس requires من `labdeps/go.mod` زي `goTest()`.
- **عقد الختام `x37-command-center`** (فصل 37): السوق + مركز القيادة + المخزن + المراقبة شغالين مع بعض 10 دقايق تحت أحداث C2، من غير فقد ولا تكرار.
- كل مكوّن سيرفر بيتقفل مع `colonyGame.stop` ومبيسمعش غير على `127.0.0.1`.

**تم لما:** لكل نظام تست integration بحل مرجعي (بيتخطّى لو الـ sandbox مانع loopback، زي `p0_test.go:261`)، وتست إن المكوّن بيتقفل من غير goroutine leak (`labkit.LeakCheck`).

## 8. المرحلة C6 — المواسم

- موسم = مستعمرة جديدة بـ seed جديد ومدة ثابتة (مثلاً 6000 tick) وقيد: «درون واحد»، «من غير رادار»، «سوق متقلب»، «عواصف كل 400 tick».
- `saves/profiles/<id>/seasons.json`: أحسن نتيجة لكل (قيد، seed). المستعمرة الأساسية متتأثرش.
- النتيجة = الكريدت المكتسبة. الواجهة تعرض نتيجتك مقابل أحسن نتيجة ليك ومقابل الحل المرجعي.
- `econ.New(seed)` + `Limits` موجودين من C1؛ المطلوب endpoint `POST /api/colony/season` وملف حفظ منفصل.

**تم لما:** موسم كامل بيتلعب بالساعة السريعة في تست، والنتيجة حتمية لنفس الكود ونفس الـ seed (في المواسم الدرون الواحد بس؛ التزامن بطبيعته مش حتمي).

## 9. تعريف «تم» لأي عقد

1. `Contract` كامل + `En`.
2. فصله في `Chapter`، و`Requires` من المفاتيح الموجودة بس.
3. حل مرجعي في `testdata/colony/<id>/` بيخلص في 40–70% من الـ `Window`.
4. الـ starter (وحل العقد اللي قبله) ميخلصوش.
5. الـ brief بيقول المطلوب وشرط النجاح بالظبط، ورابط للدرس، ومن غير كود حل.

## 10. الترتيب والاعتماديات

| مرحلة | يعتمد على | حجم تقريبي | ينفع بالتوازي مع |
|---|---|---|---|
| C1 العقود والمخططات | — | متوسط | C4 |
| E1 الموارد ومعادلة التكلفة | C1 | متوسط | C3، C4 |
| E2 المباني واللوجستيات | E1 | كبير | C3 |
| E3 طاقة وتخزين وصيانة | E2 | متوسط | C2 |
| E4 البحث والشجرة | E2 | كبير (محتوى) | E3 |
| E5 الرتب | E1، C3 | صغير | E3، E4 |
| C2 عالم حي | C1 | متوسط | C3، C4 |
| C3 القياس | C1 (للعدّادات في العقود) | صغير | C2، C4 |
| C4 مشروع حقيقي | — | متوسط | C1–C3 |
| C5 الأنظمة | C2، C4 | كبير (كل نظام لوحده) | الأنظمة مع بعض |
| C6 المواسم | C1، C3 | صغير | C5 |

التعارضات المتوقعة: `econ.State`، `colonymode.go`، `web/colony.js`، و`cmd/orbital/colony.go`. شغّل C1 الأول لوحدها عشان تثبّت شكل الـ `State`.

## 11. مخاطر معروفة

- **توازن الاقتصاد:** 14 خامة و14 منتج ومستويات لا نهائية = مساحة توازن كبيرة. اكتب simulator صغير في تست (`econ`) بيلعب بسياسة greedy ويطبع الوقت لكل رتبة، واظبط الأرقام عليه قبل اللعب الفعلي.
- **حجم المحتوى:** ≈122 عقدة بعقودها وحلولها المرجعية. فرع واحد مكتمل أحسن من تسعة ناقصين؛ ابدأ بالزراعة والتصنيع والأسطول.
- **تعقيد على المبتدئ:** المباني والـ cargo والطاقة ميظهروش قبل عقودهم. تست يتأكد إن save جديد في الفصل 3 شايف نفس واجهة النهارده + العقود.
- **التوازن:** أرقام العقود محتاجة لعب فعلي. `TestContractsSolvable` بيضمن إنها ممكنة، مش إنها ممتعة.
- **الوقت الحقيقي في التستات:** الـ sim بـ `Realtime`. التستات لازم تصغّر `Unit` و`ActionDur`، وعقود التزامن هتبقى حساسة للتحميل على الجهاز: سيب هامش 2× في الـ `Window`.
- **الثغرة في الـ state** (قسم 6): مقبولة بقرار، والعقود قابلة للتزوير بنفس الطريقة.
- **C5 حجمه كبير:** كل نظام مستقل. نظام واحد مكتمل بعقوده أحسن من ستة ناقصين.
- **loopback في البيئات المقيّدة:** تستات C5 لازم تعمل `t.Skip` زي الموجود.
- **شغل غير مرفوع:** الفرع `curriculum-expansion` فيه تعديلات مش committed في `cmd/orbital/main.go`، `runner.go`، `web/colony.js` (ميزة المراجعة). تتعمل commit قبل البدء.
