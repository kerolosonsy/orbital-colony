package levels

func init() {
	add(&Level{ID: "29-1-sqlite-crud", Chapter: 29, Kind: "lab",
		Title: "أول قاعدة", Goal: "احفظ صناديق المؤن في SQLite بقيود واستعلامات قابلة للإلغاء", Setup: seededLabRoom,
		Requires: []string{"sel:sql.ErrNoRows"}, Deps: []string{"modernc.org/sqlite"}, Bridge: []int{61}, Race: false, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch291Starter}},
		Lesson: `# صناديق مؤن بتفضل بعد إعادة التشغيل

## قبل ما تكتب
- إيه القواعد اللي لازم القاعدة تفرضها حتى لو حد تجاوز كود Go؟
- مين مسؤول عن نتيجة الاستعلام، وإمتى الاتصال يرجع متاح؟

## الجداول والقيود
SQL بتخزن صفوف في جداول بأعمدة مسماة. استورد ~database/sql~ لواجهة Go، واعمل blank import لـ ~modernc.org/sqlite~ عشان تسجّل الدرايفر المكتوب بـ Go. ~sql.Open~ بتنشئ حوض اتصالات، مش بتثبت إن القاعدة متاحة؛ ~PingContext~ بتختبر الاتصال. صاحب القاعدة مسؤول عن ~db.Close~.
~~~go
db, err := sql.Open("sqlite", "probe.db")
if err != nil { return err }
defer db.Close()
if err := db.PingContext(ctx); err != nil { return err }
~~~
التستات بتحقن قاعدة في ملف مؤقت جديد. متفتحش قاعدة تانية جوه الدوال. قاعدة ~:memory:~ الخاصة بتخص اتصال واحد؛ اتصالات مختلفة ممكن تشوف جداول مختلفة.

~CREATE TABLE~ بتعرّف الشكل. ~INTEGER~ و ~TEXT~ أنواع أعمدة. ~PRIMARY KEY~ بيميّز الصف، و ~UNIQUE~ بترفض التكرار، و ~NOT NULL~ بترفض القيمة الغايبة، و ~CHECK~ بتفرض شرط. ~IF NOT EXISTS~ بتخلي التهيئة قابلة للتكرار. ~trim~ و ~length~ هنا دوال SQL.
~~~sql
CREATE TABLE IF NOT EXISTS probes (
 code INTEGER PRIMARY KEY CHECK(code > 0),
 name TEXT NOT NULL UNIQUE CHECK(length(trim(name)) > 0),
 charge INTEGER NOT NULL CHECK(charge >= 0)
);
~~~

## كتابة آمنة
~ExecContext~ بتنفّذ أمر من غير صفوف راجعة وبترجع ~sql.Result~. ~INSERT INTO~ بتضيف صف، و ~UPDATE ... SET~ بتغيّر أعمدة، و ~DELETE FROM~ بتحذف. ~WHERE~ بتحدد الصفوف؛ من غيرها التعديل أو الحذف يشمل الجدول كله.
~~~go
result, err := db.ExecContext(ctx,
 "INSERT INTO probes(code,name,charge) VALUES(?,?,?)", 7, "Vega's probe", 4)
if err != nil { return err }
n, err := result.RowsAffected()
~~~
كل ~?~ مكان لقيمة؛ الدرايفر بيربط القيم منفصلة عن صيغة SQL. ممنوع تلزق نص المستخدم في الاستعلام. ~RowsAffected~ بتفرّق بين غياب الصف في تعديل أو حذف (صفر) والنجاح. خطأ القيد فشل، مش نتيجة فاضية ناجحة.
~~~go
_, err = db.ExecContext(ctx, "UPDATE probes SET charge=? WHERE code=?", 0, 7)
_, err = db.ExecContext(ctx, "DELETE FROM probes WHERE code=?", 7)
~~~

## قراية صف أو أكتر
~SELECT~ بتختار أعمدة. ~QueryRowContext~ متوقعة صف واحد، و ~Scan~ بتنقل الأعمدة بالترتيب لعناوين المتغيرات. الصف الغايب بيظهر عند ~Scan~ كـ ~sql.ErrNoRows~. حوّله لخطأ المجال بدل ما ترجع كائن صفري كأنه نجاح.
~~~go
var name string
err := db.QueryRowContext(ctx, "SELECT name FROM probes WHERE code=?", 7).Scan(&name)
if err == sql.ErrNoRows { return ErrCrateMissing }
~~~
~QueryContext~ بترجع ~*sql.Rows~. ~Next~ بتتقدم، و ~Scan~ بتقرا الصف الحالي، و ~Close~ بتحرر الموارد حتى لو رجعت بدري، و ~Err~ بتكشف خطأ أثناء اللف. ترتيب SQL مش مضمون؛ اطلب ~ORDER BY~ صراحة.
~~~go
rows, err := db.QueryContext(ctx, "SELECT name FROM probes ORDER BY code")
if err != nil { return err }
defer rows.Close()
for rows.Next() {
 var name string
 if err := rows.Scan(&name); err != nil { return err }
}
return rows.Err()
~~~
كل عملية بتاخد ~context.Context~ من المنادي. مرّره لكل نداء قاعدة؛ الإلغاء وعطل التخزين لازم يفضلوا أخطاء. للنجمة التالتة: ~errors.Is(err, sql.ErrNoRows)~ بتعرف الخطأ حتى لو ملفوف.

## أخطاء شائعة
~~~go
query := "SELECT name FROM probes WHERE name='" + input + "'" // حقن SQL
rows, _ := db.QueryContext(context.Background(), query) // ضياع إلغاء المنادي
for rows.Next() { } // مفيش Close مضمون ولا فحص Err النهائي
~~~

## المهمة
سيب ~main~ فاضية ونفّذ توقيعات الـ starter على ~*sql.DB~ المحقونة. ~InitCrates~ تنشئ ~crates(id,label,units)~ وقابلة للتكرار: مفتاح صحيح ~id > 0~، و ~label~ غير غايب وفريد ومش فاضي بعد SQL trim، و ~units~ صحيح غير غايب و ~units >= 0~. احفظ الاسم زي ما هو من غير تطبيع.

~PutCrate~ تضيف ~Crate{ID,Label,Units}~ المعطاة. ~GetCrate~ ترجع التلات حقول بالـ ID. ~ListCrates~ ترجع الكل بترتيب ID تصاعدي وslice فاضية مش ~nil~ لو الجدول فاضي. ~ResizeCrate~ تغيّر الوحدات بس، و ~DeleteCrate~ تحذف الـ ID بس. غياب صف في القراية أو التعديل أو الحذف يرجع ~ErrCrateMissing~، حتى لو الـ ID مش موجب. أي خطأ قاعدة تاني يفضل خطأ؛ الأمر الفاشل ميغيّرش الحالة. مرّر السياق واقفل الصفوف وافحص خطأ اللف النهائي. المنادي بيهيّئ الجدول قبل CRUD وبيملك عمر القاعدة.

~CrateError(err error) error~ تحوّل ~sql.ErrNoRows~ بالظبط لـ ~ErrCrateMissing~ وتحافظ على الباقي، حتى ~nil~. النجمة التالتة: حوّل الخطأ الملفوف كمان بـ ~errors.Is~. التستات ببيانات seeded وقيود وإلغاء وقاعدة مقفولة؛ مفيش شبكة ولا ~Report~.
`,
		Hints: []string{"ابدأ بـ InitCrates وجدول فيه PRIMARY KEY وUNIQUE وCHECK.", "مرّر القيم كـ placeholders؛ استخدم RowsAffected عشان تكتشف الغياب في التعديل والحذف.", "في ListCrates اقفل rows بـ defer وافحص rows.Err؛ للنجمة رجّع sql.ErrNoRows من الأخطاء الملفوفة كمان."},
		En: &LevelText{Title: "First Database", Goal: "Persist supply crates in SQLite with constraints and cancelable queries", Files: []LabFile{{Name: "main.go", Starter: ch291Starter}},
			Lesson: `# Supply crates that survive a restart

## Before you write
- Which rules belong in the database even when a caller bypasses Go?
- Who owns a returned result set, and when is its connection reusable?

## Tables and constraints
SQL stores rows in tables with named columns. Import ~database/sql~ for Go's database API and blank-import ~modernc.org/sqlite~ to register the pure-Go driver. ~sql.Open~ creates a pool, not proof of connectivity; ~PingContext~ checks it. The caller owns ~db.Close~.
~~~go
db, err := sql.Open("sqlite", "probe.db")
if err != nil { return err }
defer db.Close()
if err := db.PingContext(ctx); err != nil { return err }
~~~
The tests inject a pool over a fresh temporary file. Do not open another database inside your functions. A private ~:memory:~ database belongs to one connection, so multiple connections can see different tables.

~CREATE TABLE~ defines a schema. ~INTEGER~ and ~TEXT~ are column types. ~PRIMARY KEY~ identifies a row; ~UNIQUE~ rejects duplicates; ~NOT NULL~ rejects missing values; ~CHECK~ enforces an expression. ~IF NOT EXISTS~ makes this bootstrap repeatable. ~trim~ and ~length~ below are SQL functions.
~~~sql
CREATE TABLE IF NOT EXISTS probes (
 code INTEGER PRIMARY KEY CHECK(code > 0),
 name TEXT NOT NULL UNIQUE CHECK(length(trim(name)) > 0),
 charge INTEGER NOT NULL CHECK(charge >= 0)
);
~~~

## Writing safely
~ExecContext~ runs a statement without returned rows and yields a ~sql.Result~. ~INSERT INTO~ adds a row, ~UPDATE ... SET~ replaces columns, and ~DELETE FROM~ removes rows. ~WHERE~ limits affected rows. Without it, an update or delete touches the entire table.
~~~go
result, err := db.ExecContext(ctx,
 "INSERT INTO probes(code,name,charge) VALUES(?,?,?)", 7, "Vega's probe", 4)
if err != nil { return err }
n, err := result.RowsAffected()
~~~
Each ~?~ is a placeholder for one value; the driver binds arguments separately from SQL syntax. Never concatenate user text into SQL. ~RowsAffected~ distinguishes a missing update/delete (zero rows) from success. Constraint errors are failures, not success with an empty result.
~~~go
_, err = db.ExecContext(ctx, "UPDATE probes SET charge=? WHERE code=?", 0, 7)
_, err = db.ExecContext(ctx, "DELETE FROM probes WHERE code=?", 7)
~~~

## Reading one row or many
~SELECT~ chooses columns. ~QueryRowContext~ expects one row and ~Scan~ copies columns in order into pointers. Missing rows appear at ~Scan~ as ~sql.ErrNoRows~. Map that to a domain error rather than returning a zero object as success.
~~~go
var name string
err := db.QueryRowContext(ctx, "SELECT name FROM probes WHERE code=?", 7).Scan(&name)
if err == sql.ErrNoRows { return ErrCrateMissing }
~~~
~QueryContext~ returns ~*sql.Rows~. ~Next~ advances, ~Scan~ reads the current row, ~Close~ releases resources even on an early return, and ~Err~ reports iteration failures. SQL has no default ordering; use ~ORDER BY~ explicitly.
~~~go
rows, err := db.QueryContext(ctx, "SELECT name FROM probes ORDER BY code")
if err != nil { return err }
defer rows.Close()
for rows.Next() {
 var name string
 if err := rows.Scan(&name); err != nil { return err }
}
return rows.Err()
~~~
Every operation receives the caller's ~context.Context~. Pass it through all database calls; cancellation and storage failures must remain errors. Third star: ~errors.Is(err, sql.ErrNoRows)~ also recognizes wrapped missing-row errors.

## Common mistakes
~~~go
query := "SELECT name FROM probes WHERE name='" + input + "'" // Injection
rows, _ := db.QueryContext(context.Background(), query) // Lost caller cancellation
for rows.Next() { } // No guaranteed Close or terminal Err check
~~~

## Your mission
Keep ~main~ empty. Implement the starter signatures over the injected ~*sql.DB~. ~InitCrates~ creates ~crates(id,label,units)~ idempotently: integer primary key ~id > 0~, unique non-null ~label~ with nonblank SQL-trimmed content, and non-null integer ~units >= 0~. Store labels exactly, without normalization.

~PutCrate~ inserts the supplied ~Crate{ID,Label,Units}~. ~GetCrate~ returns all three stored fields by ID. ~ListCrates~ returns all rows by ascending ID and a non-nil empty slice for an empty table. ~ResizeCrate~ changes only units; ~DeleteCrate~ removes only that ID. Missing get/update/delete returns ~ErrCrateMissing~ (including nonpositive IDs). All other database errors propagate; failed statements leave state unchanged. Pass the supplied context everywhere, close rows, and check terminal iteration errors. The caller initializes the schema before CRUD and owns the database lifetime.

~CrateError(err error) error~ maps exactly ~sql.ErrNoRows~ to ~ErrCrateMissing~ and preserves other errors including ~nil~. Third star: also map wrapped ~sql.ErrNoRows~ via ~errors.Is~. Tests use seeded crate data, constraints, cancellation, and closed storage; no network or ~Report~.
`,
			Hints: []string{"Start with InitCrates and a table with PRIMARY KEY, UNIQUE, and CHECK.", "Bind values as placeholders; use RowsAffected to detect missing updates and deletes.", "In ListCrates defer rows.Close and check rows.Err; for the star recognize wrapped sql.ErrNoRows too."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch291Hidden},
	})
}

const ch291Starter = `package main

import (
	"context"
	"database/sql"
	"errors"
)

type Crate struct {
	ID    int64
	Label string
	Units int
}

var ErrCrateMissing = errors.New("crate missing")

func InitCrates(ctx context.Context, db *sql.DB) error        { return errors.New("TODO") }
func PutCrate(ctx context.Context, db *sql.DB, c Crate) error { return errors.New("TODO") }
func GetCrate(ctx context.Context, db *sql.DB, id int64) (Crate, error) {
	return Crate{}, errors.New("TODO")
}
func ListCrates(ctx context.Context, db *sql.DB) ([]Crate, error) { return nil, errors.New("TODO") }
func ResizeCrate(ctx context.Context, db *sql.DB, id int64, units int) error {
	return errors.New("TODO")
}
func DeleteCrate(ctx context.Context, db *sql.DB, id int64) error { return errors.New("TODO") }
func CrateError(err error) error                                  { return err }
func main()                                                       {}
`

const ch291Hidden = `package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"
	"testing"

	_ "modernc.org/sqlite"
)

func hiddenCrateDB(t *testing.T) *sql.DB {
	t.Helper()
	db, e := sql.Open("sqlite", filepath.Join(t.TempDir(), "crates.db"))
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	if e = InitCrates(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	return db
}
func TestHiddenCrateCRUD(t *testing.T) {
	ctx := context.Background()
	db := hiddenCrateDB(t)
	if e := InitCrates(ctx, db); e != nil {
		t.Fatal(e)
	}
	empty, e := ListCrates(ctx, db)
	if e != nil || empty == nil || len(empty) != 0 {
		t.Fatal("empty", empty, e)
	}
	rng := rand.New(rand.NewSource(2901))
	want := []Crate{}
	for i := 0; i < 25; i++ {
		c := Crate{int64(i + 1), fmt.Sprintf("مؤن '%d; --", rng.Int63()), rng.Intn(200)}
		want = append(want, c)
	}
	for _, i := range rng.Perm(len(want)) {
		if e := PutCrate(ctx, db, want[i]); e != nil {
			t.Fatal(e)
		}
	}
	got, e := ListCrates(ctx, db)
	if e != nil || !reflect.DeepEqual(got, want) {
		t.Fatal("ordered list", got, e)
	}
	for _, c := range want {
		if got, e := GetCrate(ctx, db, c.ID); e != nil || got != c {
			t.Fatal(got, e)
		}
	}
	for _, c := range []Crate{{1, "other", 2}, {99, want[0].Label, 2}, {0, "bad", 2}, {99, " ", 2}, {99, "bad", -1}} {
		if e := PutCrate(ctx, db, c); e == nil {
			t.Fatal("constraint accepted", c)
		}
	}
	if e := ResizeCrate(ctx, db, 1, -1); e == nil {
		t.Fatal("negative update")
	}
	if e := ResizeCrate(ctx, db, 1, 0); e != nil {
		t.Fatal(e)
	}
	c, e := GetCrate(ctx, db, 1)
	if e != nil || c.Units != 0 || c.Label != want[0].Label {
		t.Fatal(c, e)
	}
	if e := DeleteCrate(ctx, db, 1); e != nil {
		t.Fatal(e)
	}
	for _, e := range []error{DeleteCrate(ctx, db, 1), ResizeCrate(ctx, db, 1, 2)} {
		if !errors.Is(e, ErrCrateMissing) {
			t.Fatal("missing mutation", e)
		}
	}
	if _, e := GetCrate(ctx, db, 1); !errors.Is(e, ErrCrateMissing) {
		t.Fatal("missing get", e)
	}
	if db.Stats().InUse != 0 {
		t.Fatal("rows leaked")
	}
}
func TestHiddenCrateCancellationAndStorageFailure(t *testing.T) {
	db := hiddenCrateDB(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, get := GetCrate(ctx, db, 1)
	_, list := ListCrates(ctx, db)
	for _, e := range []error{InitCrates(ctx, db), PutCrate(ctx, db, Crate{1, "a", 1}), get, list, ResizeCrate(ctx, db, 1, 1), DeleteCrate(ctx, db, 1)} {
		if !errors.Is(e, context.Canceled) {
			t.Fatal("cancellation", e)
		}
	}
	ctx = context.Background()
	if _, e := db.ExecContext(ctx, "INSERT INTO crates VALUES(1,'x',-1)"); e == nil {
		t.Fatal("schema bypass accepted")
	}
	db.Close()
	if _, e := ListCrates(ctx, db); e == nil {
		t.Fatal("closed database")
	}
}
func TestHiddenBonusWrappedMissing(t *testing.T) {
	if !errors.Is(CrateError(fmt.Errorf("lookup: %w", sql.ErrNoRows)), ErrCrateMissing) {
		t.Fatal("wrapped missing")
	}
	sentinel := errors.New("storage failure")
	if CrateError(sentinel) != sentinel || CrateError(nil) != nil {
		t.Fatal("lost error")
	}
}
`

func init() {
	add(&Level{ID: "29-2-migrations", Chapter: 29, Kind: "lab",
		Title: "الترحيلات", Goal: "طوّر مخطط مخزن المستعمرة بترحيلات up/down وسجل نسخ ذرّي", Setup: seededLabRoom,
		Requires: []string{"sel:sql.Tx"}, Deps: []string{"modernc.org/sqlite"}, Bridge: []int{64}, Race: false, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch292Starter}},
		Lesson: `# تطوير مخطط المساكن

## قبل ما تكتب
- لو الخطوة التانية فشلت بعد إنشاء جدول، إيه اللي المفروض يفضل؟
- هترفض أرقام النسخ المكررة إزاي قبل أي تعديل؟

## سجل النسخ
الترحيل زوج SQL: ~Up~ بيطوّر المخطط، و ~Down~ بيرجع الخطوة. مثال عنصر مصدر: ~Migration{Version: 3, Up: "CREATE TABLE antennas(id INTEGER)", Down: "DROP TABLE antennas"}~. ~DROP TABLE~ بتحذف الجدول وبياناته؛ الرجوع مش نسخة احتياطية.

احفظ النسخ المطبقة في جدول مخصص، مش متغير Go. ~SELECT version FROM habitat_versions ORDER BY version~ بيفضل بعد إعادة التشغيل. السجل لازم يكون بادئة من المصدر المرتب: النسختين ~3, 9~ صح لو هما أول عنصرين؛ تسجيل ~9~ لوحدها مش متسق. ارفض النسخة المجهولة بدل تخمين الرجوع.
~~~sql
CREATE TABLE IF NOT EXISTS habitat_versions(version INTEGER PRIMARY KEY);
INSERT INTO habitat_versions(version) VALUES(3);
DELETE FROM habitat_versions WHERE version=3;
~~~

## المعاملة قبل مشغّل الترحيلات
~sql.Tx~ بتحجز اتصال وبتجمع أوامر في وحدة واحدة. ~BeginTx(ctx, nil)~ بتبدأ بالاختيارات الافتراضية. نفّذ كل أمر على المعاملة نفسها، مش على ~db~. ~Commit~ بتنشر الكل، و ~Rollback~ بتلغي الكل. الرجوع المؤجل آمن بعد نجاح الحفظ.
~~~go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback()
if _, err := tx.ExecContext(ctx, "CREATE TABLE antennas(id INTEGER)"); err != nil { return err }
if _, err := tx.ExecContext(ctx, "INSERT INTO habitat_versions VALUES(?)", 3); err != nil { return err }
return tx.Commit()
~~~
SQLite هنا بتدعم أوامر المخطط جوه المعاملة: إنشاء الجدول وتسجيل النسخة بيتحفظوا مع بعض. نصوص SQL المحقونة موثوقة ومكتوبة من المؤلف، مش مدخل مستخدم. النص ممكن يحتوي أوامر متعددة، بس ممنوع أوامر التحكم في المعاملة؛ المشغّل هو اللي بيملكها.

## تقدّم جزئي صريح
استخدم معاملة لكل نسخة. لو التالتة فشلت، أول اتنين يفضلوا مطبقين. فشل الرجوع يسيب النسخة مسجلة ويرجع تغييرات مخطط الخطوة. مرّر السياق لكل نداء قاعدة ومعاملة. متدّعيش قفل بين عمليات نظام مختلفة؛ المنادي بينفّذ النداءات بالتتابع.

افحص المصدر كله قبل إنشاء السجل: الأرقام موجبة ومتزايدة، والنصين مش فاضيين. كده التكرار والترتيب النازل وغياب نص الرجوع بيترفضوا قبل أي أثر. ~strings.TrimSpace(m.Up) == ""~ بتكشف النص اللي كله مسافات. الدرس بيتتبع هوية النسخة بس؛ البصمات وتعديل نص مطبق برا العقد.

## أخطاء شائعة
~~~go
tx, _ := db.BeginTx(ctx, nil)
db.ExecContext(ctx, migration.Up) // برا المعاملة، وممكن يقفل حوض باتصال واحد
tx.Commit() // تجاهل الخطأ وتسجيل النسخة مش ذرّي
~~~

## المهمة
نفّذ ~MigrateHabitat(ctx context.Context, db *sql.DB, source []Migration, down bool) error~. المدخل نوع الـ starter: ~Migration{Version int; Up, Down string}~. افحص المصدر كله قبل أي تغيير؛ الفجوات مسموحة، والتكرار والأرقام غير الموجبة والترتيب الغلط والنص الفاضي أخطاء. هيّئ ~habitat_versions(version INTEGER PRIMARY KEY)~ بشكل قابل للتكرار. اقرا السجل وارفض أي تاريخ مش بادئة مطابقة للمصدر.

مع ~down=false~ طبّق كل النسخ الجديدة بترتيب المصدر، معاملة لكل نسخة تشمل تسجيلها. تكرار up من غير جديد مبيعملش حاجة. مع ~down=true~ ارجع آخر نسخة مطبقة بس واحذفها من السجل في نفس المعاملة؛ السجل الفاضي نجاح من غير تغيير. عند فشل خطوة ارجعها وارجع خطأ وسيب الخطوات السابقة المحفوظة. مرّر أخطاء الإلغاء والقاعدة. اقفل نتائج القراية قبل فتح معاملة. المنادي بيملك القاعدة وبيسلسل النداءات وبيوفر نصوص موثوقة من غير أوامر تحكم في المعاملة. سيب ~main~ فاضية؛ مفيش شبكة ولا ~Report~.
`,
		Hints: []string{"افحص المصدر كله قبل إنشاء سجل النسخ: أرقام موجبة ومتزايدة ونصوص مش فاضية.", "اقرا النسخ المطبقة وتأكد إنها بادئة من المصدر، وبعدها اختار pending للـ up أو آخر نسخة للـ down.", "في كل خطوة افتح BeginTx؛ نفّذ SQL وتعديل السجل على tx نفسها، واعمل defer Rollback وبعد النجاح Commit."},
		En: &LevelText{Title: "Migrations", Goal: "Evolve the colony store schema with up/down migrations and an atomic version ledger", Files: []LabFile{{Name: "main.go", Starter: ch292Starter}},
			Lesson: `# Evolving the habitat schema

## Before you write
- If step two fails after creating a table, which changes should survive?
- How can you reject duplicate versions before any schema mutation?

## A version ledger
A migration is a pair of SQL scripts: ~Up~ advances the schema and ~Down~ reverses that step. ~Migration{Version: 3, Up: "CREATE TABLE antennas(id INTEGER)", Down: "DROP TABLE antennas"}~ is one source item. ~DROP TABLE~ removes a table and its data; down is destructive, not a backup.

Keep applied versions in a dedicated table, not a variable in Go. ~SELECT version FROM habitat_versions ORDER BY version~ survives process restarts. Applied history must be a prefix of the ordered source: versions ~3, 9~ are fine if those are the first two source entries; recording ~9~ alone is inconsistent. Reject unknown versions rather than guessing a rollback.
~~~sql
CREATE TABLE IF NOT EXISTS habitat_versions(version INTEGER PRIMARY KEY);
INSERT INTO habitat_versions(version) VALUES(3);
DELETE FROM habitat_versions WHERE version=3;
~~~

## Transactions come before migration runners
A ~sql.Tx~ reserves a connection and groups statements into one unit. ~BeginTx(ctx, nil)~ starts with default options. Run every statement through the transaction, not through ~db~. ~Commit~ publishes all changes; ~Rollback~ discards them. A deferred rollback is harmless after a successful commit.
~~~go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback()
if _, err := tx.ExecContext(ctx, "CREATE TABLE antennas(id INTEGER)"); err != nil { return err }
if _, err := tx.ExecContext(ctx, "INSERT INTO habitat_versions VALUES(?)", 3); err != nil { return err }
return tx.Commit()
~~~
SQLite supports transactional DDL here: table creation and version insertion commit together. The injected SQL scripts are trusted author-written code, not user text. A script can contain multiple statements. It must not contain transaction-control statements: the runner owns the transaction.

## Honest partial progress
Use one transaction per version, not one transaction for the entire source. If the third up fails, the first two stay applied. A down failure keeps its version recorded and restores that step's schema changes. Check the context before database work and pass it through each transaction. Do not claim cross-process locking: callers serialize this small runner.

Validate the entire source before bootstrap: versions positive and strictly increasing, both scripts nonblank. This rejects duplicate versions, descending order, and missing rollback scripts before effects. ~strings.TrimSpace(m.Up) == ""~ detects a whitespace-only script. This lesson tracks version identity only; checksums and editing already applied scripts are outside its contract.

## Common mistakes
~~~go
tx, _ := db.BeginTx(ctx, nil)
db.ExecContext(ctx, migration.Up) // Outside tx; may deadlock a one-connection pool
tx.Commit() // Error ignored, and no version recorded atomically
~~~

## Your mission
Implement ~MigrateHabitat(ctx context.Context, db *sql.DB, source []Migration, down bool) error~. Inputs use the starter's ~Migration{Version int; Up, Down string}~. Fully validate source before any database change; gaps are allowed, duplicate/nonpositive/out-of-order versions or blank scripts are errors. Bootstrap ~habitat_versions(version INTEGER PRIMARY KEY)~ idempotently. Read applied versions and reject history that is not an exact prefix of source.

With ~down=false~, apply every pending up in source order, one transaction each, including insertion of its version. Repeating up is a no-op. With ~down=true~, revert only the latest applied migration and delete its version in the same transaction; an empty history is a successful no-op. On a failed step roll back that step and return an error, preserving earlier committed steps. Propagate cancellation and database errors. Close any result sets before opening a transaction. The caller owns the database, serializes calls, and supplies trusted scripts without transaction commands. Keep ~main~ empty; no network or ~Report~.
`,
			Hints: []string{"Validate the whole source before bootstrapping the ledger: positive increasing versions and nonblank SQL.", "Read applied versions and require a prefix of the source; then choose pending up steps or the latest down step.", "For each step BeginTx, execute SQL and the ledger mutation on that same tx, defer Rollback, then Commit on success."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch292Hidden},
	})
}

const ch292Starter = `package main

import (
	"context"
	"database/sql"
	"errors"
)

type Migration struct {
	Version  int
	Up, Down string
}

func MigrateHabitat(ctx context.Context, db *sql.DB, source []Migration, down bool) error {
	return errors.New("TODO")
}
func main() {}
`

const ch292Hidden = `package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func hiddenMigrationDB(t *testing.T) *sql.DB {
	t.Helper()
	db, e := sql.Open("sqlite", filepath.Join(t.TempDir(), "habitat.db"))
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	return db
}
func hiddenVersionCount(t *testing.T, db *sql.DB, want int) {
	t.Helper()
	var n int
	if e := db.QueryRow("SELECT count(*) FROM habitat_versions").Scan(&n); e != nil || n != want {
		t.Fatal("versions", n, e, "want", want)
	}
}
func TestHiddenMigrationRoundTrip(t *testing.T) {
	db := hiddenMigrationDB(t)
	ctx := context.Background()
	rng := rand.New(rand.NewSource(2902))
	source := []Migration{}
	names := []string{}
	for i := 1; i <= 6; i++ {
		name := fmt.Sprintf("module_%d", rng.Intn(1000000))
		names = append(names, name)
		source = append(source, Migration{i * 3, "CREATE TABLE " + name + "(id INTEGER)", "DROP TABLE " + name})
	}
	for i := 0; i < 2; i++ {
		if e := MigrateHabitat(ctx, db, source, false); e != nil {
			t.Fatal(e)
		}
		hiddenVersionCount(t, db, 6)
		for _, name := range names {
			hiddenMigrationTable(t, db, name, true)
		}
	}
	for i := 5; i >= 0; i-- {
		if e := MigrateHabitat(ctx, db, source, true); e != nil {
			t.Fatal(e)
		}
		hiddenVersionCount(t, db, i)
		for j, name := range names {
			hiddenMigrationTable(t, db, name, j < i)
		}
	}
	if e := MigrateHabitat(ctx, db, source, true); e != nil {
		t.Fatal("empty down", e)
	}
	if e := MigrateHabitat(ctx, db, nil, false); e != nil {
		t.Fatal("empty up", e)
	}
	if e := MigrateHabitat(ctx, db, source, false); e != nil {
		t.Fatal("reapply", e)
	}
}
func TestHiddenMigrationValidation(t *testing.T) {
	for _, source := range [][]Migration{{{0, "SELECT 1", "SELECT 1"}}, {{2, "SELECT 1", "SELECT 1"}, {2, "SELECT 2", "SELECT 2"}}, {{3, "SELECT 1", "SELECT 1"}, {1, "SELECT 2", "SELECT 2"}}, {{1, " ", "SELECT 1"}}, {{1, "SELECT 1", ""}}} {
		db := hiddenMigrationDB(t)
		if e := MigrateHabitat(context.Background(), db, source, false); e == nil {
			t.Fatal("invalid source accepted", source)
		}
		var n int
		if e := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE name='habitat_versions'").Scan(&n); e != nil || n != 0 {
			t.Fatal("validation mutated schema", n, e)
		}
	}
}
func TestHiddenMigrationFailureAndHistory(t *testing.T) {
	ctx := context.Background()
	db := hiddenMigrationDB(t)
	source := []Migration{{1, "CREATE TABLE core(id INTEGER)", "DROP TABLE core"}, {4, "CREATE TABLE partial(id INTEGER); INSERT INTO absent VALUES(1)", "DROP TABLE partial"}}
	if e := MigrateHabitat(ctx, db, source, false); e == nil {
		t.Fatal("failure swallowed")
	}
	hiddenVersionCount(t, db, 1)
	var n int
	if e := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE name='partial'").Scan(&n); e != nil || n != 0 {
		t.Fatal("DDL not rolled back", n, e)
	}
	source[1].Up = "CREATE TABLE partial(id INTEGER)"
	if e := MigrateHabitat(ctx, db, source, false); e != nil {
		t.Fatal(e)
	}
	source[1].Down = "DROP TABLE partial; INSERT INTO absent VALUES(1)"
	if e := MigrateHabitat(ctx, db, source, true); e == nil {
		t.Fatal("down failure swallowed")
	}
	hiddenVersionCount(t, db, 2)
	if e := db.QueryRow("SELECT count(*) FROM partial").Scan(&n); e != nil {
		t.Fatal("down DDL not rolled back", e)
	}
	if e := MigrateHabitat(ctx, db, source[:1], false); e == nil {
		t.Fatal("unknown version accepted")
	}
	if _, e := db.Exec("DELETE FROM habitat_versions WHERE version=1"); e != nil {
		t.Fatal(e)
	}
	if e := MigrateHabitat(ctx, db, source, false); e == nil {
		t.Fatal("non-prefix accepted")
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if e := MigrateHabitat(canceled, db, source, false); !errors.Is(e, context.Canceled) {
		t.Fatal("cancel", e)
	}
}

func hiddenMigrationTable(t *testing.T, db *sql.DB, name string, present bool) {
	t.Helper()
	var n int
	if e := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", name).Scan(&n); e != nil {
		t.Fatal(e)
	}
	if (n == 1) != present {
		t.Fatal("schema state", name, n, present)
	}
}
`

func init() {
	add(&Level{ID: "29-3-transfer", Chapter: 29, Kind: "lab",
		Title: "تحويل ذرّي", Goal: "انقل وحدات الأكسجين بين خزانات من غير فقد أو تكرار تحت التزامن", Setup: seededLabRoom,
		Requires: []string{"sel:sql.Tx"}, Deps: []string{"modernc.org/sqlite"}, Bridge: []int{66}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch293Starter}},
		Lesson: `# نقل أكسجين بحفظ واحد

## قبل ما تكتب
- إيه اللي يحصل لو الصمام فشل بعد الخصم وقبل الإضافة؟
- هل محاولتين بنفس request ID ممكن يصرفوا الأكسجين مرتين؟

## حفظ الإجمالي هو العقد
التحويل بيغيّر صفين وبيسجّل طلب واحد. مجموع أرصدة الخزانات لازم يفضل ثابت. أدوات المعاملة من الترحيلات بتنفع للبيانات العادية كمان. ~OxygenMove{Request: "air-7", From: "east", To: "west", Units: 4}~ بتحدد عملية واحدة غير قابلة للتقسيم. ~InitOxygen(ctx, db)~ بتهيّئ الجداول بس؛ المنادي بيملأ الأرصدة.

افحص سجل الطلبات جوه المعاملة. ~SELECT request FROM oxygen_moves WHERE request=?~ بيكشف الطلب المحفوظ قبل كده. ارجع ~ErrReplay~ من غير تغيير الأرصدة. المفتاح عام؛ إعادة استخدامه لمدخل مختلف برضه تكرار.
~~~go
var existing string
err := tx.QueryRowContext(ctx,
 "SELECT request FROM oxygen_moves WHERE request=?", "air-7").Scan(&existing)
if err == nil { return ErrReplay }
if !errors.Is(err, sql.ErrNoRows) { return err }
~~~

## احرس الخصم في SQL
قراية رصيد وبعدين تعديله برا معاملة بتفتح سباق. حط الشرط في التعديل وافحص ~RowsAffected~. صفر صفوف يعني مصدر غايب أو رصيد مش كفاية؛ الاتنين هنا ~ErrFuel~.
~~~go
result, err := tx.ExecContext(ctx,
 "UPDATE oxygen_tanks SET units=units-? WHERE name=? AND units>=?", 4, "east", 4)
if err != nil { return err }
n, err := result.RowsAffected()
if err != nil { return err }
if n != 1 { return ErrFuel }
~~~
الإضافة لازم تأثر في صف واحد كمان؛ الوجهة الغايبة لازم ترجع الخصم. ~UPDATE oxygen_tanks SET units=units+? WHERE name=?~ بتغيّر الوجهة بس. بعد الكتابتين، ~INSERT INTO oxygen_moves(request) VALUES(?)~ بتسجّل الاكتمال، وبعدها ~tx.Commit()~. كل رجوع بدري محمي بـ ~defer tx.Rollback()~.

## حافظ على الحساب الصحيح
أعداد SQLite الصحيحة موقّعة بحجم 64 بت. الجمع اللي يتجاوز الحد ممكن يتحوّل لقيمة عشرية بدل الفشل. ~math.MaxInt64~ أكبر عدد صحيح موقّع بحجم 64 بت؛ احسب السعة الباقية بالطرح قبل الإضافة.
~~~go
capacity := int64(math.MaxInt64) - int64(4)
~~~
احرس الإضافة بشرط ~AND units<=?~، واربط ~math.MaxInt64-m.Units~. الوجهة المليانة ترجع ~ErrFuel~ ويتراجع الخصم؛ الوصول للحد بالظبط مسموح.

## حقن الفشل من غير اعتماد على الجهاز
~afterDebit func() error~ نقطة حقن بعد نجاح الخصم وقبل الإضافة. لو مش ~nil~ بتتنادى مرة؛ خطأ منها يلغي العملية ويفضل قابل للكشف بـ ~errors.Is~. لو ~nil~ مفيش نداء. ممنوع الـ hook تستخدم نفس الحوض؛ المعاملة حاجزة الاتصال الوحيد.
~~~go
fault := errors.New("valve unavailable")
err := MoveOxygen(ctx, db, move, func() error { return fault })
if errors.Is(err, fault) { /* الرصيدين زي ما كانوا */ }
~~~
المنادي بيضبط حوض SQLite على اتصال واحد قبل الشغل. ده بيسلسل المعاملات المتزامنة في المعمل؛ متدّعيش دعم إعادة المحاولة بين عمليات نظام أو اتصالات متعددة. ~CHECK(units>=0)~ حارس إضافي في القاعدة. الأسماء والوحدات محقونة؛ التنفيذ مش محتاج ساعة ولا عشوائية.

## أخطاء شائعة
~~~go
db.ExecContext(ctx, "UPDATE oxygen_tanks SET units=units-4 WHERE name='east'")
return errors.New("credit failed") // الخصم اتحفظ لوحده خلاص
~~~

## المهمة
نفّذ ~InitOxygen(ctx context.Context, db *sql.DB) error~ بشكل قابل للتكرار لإنشاء ~oxygen_tanks(name TEXT PRIMARY KEY, units INTEGER NOT NULL CHECK(units>=0))~ و ~oxygen_moves(request TEXT PRIMARY KEY)~. متضيفش خزانات.

نفّذ ~MoveOxygen(ctx context.Context, db *sql.DB, m OxygenMove, afterDebit func() error) error~. ارفض الطلب أو المصدر أو الوجهة الفاضية، وتساوي الخزانين، والوحدات غير الموجبة بـ ~ErrTransfer~ قبل الكتابة. طلب صالح بمفتاح مسجل يرجع ~ErrReplay~ حتى لو باقي المدخل مختلف. غير كده اخصم، ونادي الـ hook الاختيارية، وضيف، وسجّل الطلب في معاملة واحدة. الخزان الغايب أو رصيد المصدر غير الكافي أو تجاوز عدد الوجهة الصحيح يرجع ~ErrFuel~. تحويل الرصيد كله ورصيد وجهة صفر مسموحين. مرّر السياق لكل SQL ومرّر باقي الأخطاء. أي فشل يسيب الرصيدين والسجل زي ما كانوا ويسمح بإعادة استخدام الطلب. احفظ مرة واحدة عند النجاح. المنادي بيستخدم SQLite باتصال واحد بيسلسل المعاملات حتى مع نداءات متزامنة وبيملك عمر القاعدة. مفيش شبكة؛ سيب ~main~ فاضية.

النجمة التالتة: إلغاء السياق من الـ hook، حتى لو رجّعت ~nil~، لازم يفشل بخطأ يطابق ~context.Canceled~، ويرجع كل التغييرات، ويسمح بإعادة الطلب بسياق جديد.
`,
		Hints: []string{"ابدأ بمعاملة واقرا request ID جواها قبل الخصم.", "استخدم UPDATE بشرط units>=? وافحص RowsAffected، وبعد الخصم نفّذ afterDebit لو موجودة.", "الإضافة وتسجيل الطلب جوه نفس tx؛ أي فشل يرجع من غير Commit، والـ defer يعمل Rollback حتى مع إلغاء السياق."},
		En: &LevelText{Title: "Atomic Transfer", Goal: "Transfer oxygen units between tanks without loss or duplication under concurrency", Files: []LabFile{{Name: "main.go", Starter: ch293Starter}},
			Lesson: `# One oxygen move, one commit

## Before you write
- What happens if the valve fails after debit but before credit?
- Can two retries with the same request ID both spend oxygen?

## Conservation is the contract
A transfer changes two rows and records one request. The sum of tank balances must stay constant. The transaction tools from the migrations mission apply to ordinary data too. ~OxygenMove{Request: "air-7", From: "east", To: "west", Units: 4}~ names one indivisible operation. ~InitOxygen(ctx, db)~ only bootstraps tables; tank balances come from the caller.

Check the request ledger inside the transaction. ~SELECT request FROM oxygen_moves WHERE request=?~ detects an already committed request. Return ~ErrReplay~ without touching balances. The request key is global: reusing it for different input is still a replay.
~~~go
var existing string
err := tx.QueryRowContext(ctx,
 "SELECT request FROM oxygen_moves WHERE request=?", "air-7").Scan(&existing)
if err == nil { return ErrReplay }
if !errors.Is(err, sql.ErrNoRows) { return err }
~~~

## Guard the debit in SQL
Reading a balance and later updating it outside a transaction leaves a race. Put the condition in the update and inspect ~RowsAffected~. Zero affected rows means missing source or insufficient balance, both ~ErrFuel~ here.
~~~go
result, err := tx.ExecContext(ctx,
 "UPDATE oxygen_tanks SET units=units-? WHERE name=? AND units>=?", 4, "east", 4)
if err != nil { return err }
n, err := result.RowsAffected()
if err != nil { return err }
if n != 1 { return ErrFuel }
~~~
Credit must also affect exactly one row: a missing destination must roll back the debit. ~UPDATE oxygen_tanks SET units=units+? WHERE name=?~ changes only that destination. Only after both writes should ~INSERT INTO oxygen_moves(request) VALUES(?)~ record completion, followed by ~tx.Commit()~. Every early return is protected by ~defer tx.Rollback()~.

## Keep integer arithmetic exact
SQLite integers are signed 64-bit. Overflowing addition can become a floating-point value instead of failing. ~math.MaxInt64~ is the largest signed 64-bit integer; compute remaining capacity by subtraction before adding.
~~~go
capacity := int64(math.MaxInt64) - int64(4)
~~~
Guard the credit with ~AND units<=?~, binding ~math.MaxInt64-m.Units~. A full destination returns ~ErrFuel~ and rolls back the debit; reaching the exact limit is valid.

## Inject failure without machine dependencies
~afterDebit func() error~ is a boundary hook called after the successful debit and before credit. A non-nil hook runs once; a returned error aborts and remains discoverable with ~errors.Is~. A nil hook does nothing. The hook must not call the same pool: with one connection the transaction already owns it.
~~~go
fault := errors.New("valve unavailable")
err := MoveOxygen(ctx, db, move, func() error { return fault })
if errors.Is(err, fault) { /* both balances unchanged */ }
~~~
The caller configures the SQLite pool to one connection before any calls. This serializes concurrent transactions for this lab; do not promise multi-process or multiple-connection retry handling. A ~CHECK(units>=0)~ is an additional database guard. IDs and units are supplied, so no clock or randomness belongs in the implementation.

## Common mistakes
~~~go
db.ExecContext(ctx, "UPDATE oxygen_tanks SET units=units-4 WHERE name='east'")
return errors.New("credit failed") // Debit already committed separately
~~~

## Your mission
Implement ~InitOxygen(ctx context.Context, db *sql.DB) error~ idempotently, creating ~oxygen_tanks(name TEXT PRIMARY KEY, units INTEGER NOT NULL CHECK(units>=0))~ and ~oxygen_moves(request TEXT PRIMARY KEY)~. It inserts no tanks.

Implement ~MoveOxygen(ctx context.Context, db *sql.DB, m OxygenMove, afterDebit func() error) error~. Reject empty request/source/destination, equal tank names, or nonpositive units as ~ErrTransfer~, before writes. Valid existing request IDs return ~ErrReplay~ even when the payload differs. Otherwise atomically debit, call the optional hook, credit, and record the request. Missing tanks, insufficient source units, or destination integer overflow return ~ErrFuel~. Exact-balance transfers and a zero destination balance are valid. Pass context through all SQL calls; propagate other errors. Any failure leaves both balances and the ledger unchanged and permits retrying the same request. Commit once on success. The caller uses a one-connection SQLite pool, serializing transactions even for concurrent calls, and owns its lifetime. No network; keep ~main~ empty.

Third star: cancellation triggered by the hook, even if it returns ~nil~, must abort with an error matching ~context.Canceled~, roll back, and leave the request reusable with a fresh context.
`,
			Hints: []string{"Begin a transaction and check the request ID inside it before debiting.", "Use an UPDATE guarded by units>=? and check RowsAffected; call afterDebit after the debit when non-nil.", "Credit and record the request in the same tx; return without Commit on any failure so deferred Rollback also covers cancellation."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch293Hidden},
	})
}

const ch293Starter = `package main

import (
	"context"
	"database/sql"
	"errors"
)

var ErrFuel = errors.New("unavailable oxygen")
var ErrReplay = errors.New("request already applied")
var ErrTransfer = errors.New("invalid transfer")

type OxygenMove struct {
	Request, From, To string
	Units             int64
}

func InitOxygen(ctx context.Context, db *sql.DB) error { return errors.New("TODO") }
func MoveOxygen(ctx context.Context, db *sql.DB, m OxygenMove, afterDebit func() error) error {
	return errors.New("TODO")
}
func main() {}
`

const ch293Hidden = `package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"path/filepath"
	"sync"
	"testing"

	_ "modernc.org/sqlite"
)

func hiddenOxygenDB(t *testing.T) *sql.DB {
	t.Helper()
	db, e := sql.Open("sqlite", filepath.Join(t.TempDir(), "oxygen.db"))
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	if e = InitOxygen(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec("INSERT INTO oxygen_tanks VALUES('a',10000),('b',0)"); e != nil {
		t.Fatal(e)
	}
	return db
}
func hiddenOxygenBalance(t *testing.T, db *sql.DB, a, b, moves int) {
	t.Helper()
	for name, want := range map[string]int{"a": a, "b": b} {
		var n int
		if e := db.QueryRow("SELECT units FROM oxygen_tanks WHERE name=?", name).Scan(&n); e != nil || n != want {
			t.Fatal("balance", name, n, want, e)
		}
	}
	var n int
	if e := db.QueryRow("SELECT count(*) FROM oxygen_moves").Scan(&n); e != nil || n != moves {
		t.Fatal("ledger", n, moves, e)
	}
}
func TestHiddenOxygenTransfers(t *testing.T) {
	db := hiddenOxygenDB(t)
	ctx := context.Background()
	rng := rand.New(rand.NewSource(2903))
	total := 0
	for i := 0; i < 20; i++ {
		units := rng.Intn(50) + 1
		m := OxygenMove{fmt.Sprint(i), "a", "b", int64(units)}
		if e := MoveOxygen(ctx, db, m, nil); e != nil {
			t.Fatal(e)
		}
		total += units
		if e := MoveOxygen(ctx, db, m, nil); !errors.Is(e, ErrReplay) {
			t.Fatal("replay", e)
		}
	}
	hiddenOxygenBalance(t, db, 10000-total, total, 20)
	for _, m := range []OxygenMove{{"x", "a", "b", 0}, {"x", "a", "b", -1}, {"", "a", "b", 1}, {"x", "a", "a", 1}, {"x", "", "b", 1}} {
		if e := MoveOxygen(ctx, db, m, nil); !errors.Is(e, ErrTransfer) {
			t.Fatal("invalid", m, e)
		}
	}
	for _, m := range []OxygenMove{{"x", "a", "b", 10001}, {"x", "missing", "b", 1}, {"x", "a", "missing", 1}} {
		if e := MoveOxygen(ctx, db, m, nil); !errors.Is(e, ErrFuel) {
			t.Fatal("unavailable", m, e)
		}
	}
	sentinel := errors.New("valve jammed")
	calls := 0
	m := OxygenMove{"failed", "a", "b", 5}
	if e := MoveOxygen(ctx, db, m, func() error { calls++; return sentinel }); !errors.Is(e, sentinel) || calls != 1 {
		t.Fatal("hook", calls, e)
	}
	hiddenOxygenBalance(t, db, 10000-total, total, 20)
	if e := MoveOxygen(ctx, db, m, nil); e != nil {
		t.Fatal("retry after failure", e)
	}
	hiddenOxygenBalance(t, db, 9995-total, total+5, 21)
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if e := MoveOxygen(canceled, db, OxygenMove{"cancel", "a", "b", 1}, nil); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
func TestHiddenOxygenConcurrentReplay(t *testing.T) {
	db := hiddenOxygenDB(t)
	var wg sync.WaitGroup
	results := make(chan error, 32)
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- MoveOxygen(context.Background(), db, OxygenMove{"same", "a", "b", 10000}, nil)
		}()
	}
	wg.Wait()
	close(results)
	success := 0
	for e := range results {
		if e == nil {
			success++
		} else if !errors.Is(e, ErrReplay) {
			t.Fatal(e)
		}
	}
	if success != 1 {
		t.Fatal("successes", success)
	}
	hiddenOxygenBalance(t, db, 0, 10000, 1)
}
func TestHiddenBonusOxygenMidCancel(t *testing.T) {
	db := hiddenOxygenDB(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m := OxygenMove{"retry", "a", "b", 5}
	if e := MoveOxygen(ctx, db, m, func() error { cancel(); return nil }); !errors.Is(e, context.Canceled) {
		t.Fatal("mid-transaction cancel", e)
	}
	hiddenOxygenBalance(t, db, 10000, 0, 0)
	if e := MoveOxygen(context.Background(), db, m, nil); e != nil {
		t.Fatal("retry canceled request", e)
	}
	hiddenOxygenBalance(t, db, 9995, 5, 1)
}

func TestHiddenOxygenCreditFailureRollsBack(t *testing.T) {
	db := hiddenOxygenDB(t)
	ctx := context.Background()
	if _, e := db.Exec("CREATE TRIGGER jammed_credit BEFORE UPDATE ON oxygen_tanks WHEN NEW.name='b' BEGIN SELECT RAISE(ABORT, 'jammed'); END"); e != nil {
		t.Fatal(e)
	}
	m := OxygenMove{"valve", "a", "b", 7}
	if e := MoveOxygen(ctx, db, m, nil); e == nil {
		t.Fatal("credit failure swallowed")
	}
	hiddenOxygenBalance(t, db, 10000, 0, 0)
	if _, e := db.Exec("DROP TRIGGER jammed_credit"); e != nil {
		t.Fatal(e)
	}
	if e := MoveOxygen(ctx, db, m, nil); e != nil {
		t.Fatal("retry after SQL failure", e)
	}
	hiddenOxygenBalance(t, db, 9993, 7, 1)
}
func TestHiddenOxygenConcurrentDistinctRequests(t *testing.T) {
	db := hiddenOxygenDB(t)
	var wg sync.WaitGroup
	results := make(chan error, 32)
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results <- MoveOxygen(context.Background(), db, OxygenMove{fmt.Sprint(i), "a", "b", 500}, nil)
		}(i)
	}
	wg.Wait()
	close(results)
	success := 0
	for e := range results {
		if e == nil {
			success++
		} else if !errors.Is(e, ErrFuel) {
			t.Fatal(e)
		}
	}
	if success != 20 {
		t.Fatal("overspent source", success)
	}
	hiddenOxygenBalance(t, db, 0, 10000, 20)
}

func TestHiddenOxygenIntegerCapacity(t *testing.T) {
	db := hiddenOxygenDB(t)
	ctx := context.Background()
	if _, e := db.Exec("UPDATE oxygen_tanks SET units=? WHERE name='b'", int64(math.MaxInt64)); e != nil {
		t.Fatal(e)
	}
	if e := MoveOxygen(ctx, db, OxygenMove{"overflow", "a", "b", 1}, nil); !errors.Is(e, ErrFuel) {
		t.Fatal("overflow must roll back", e)
	}
	hiddenOxygenBalance(t, db, 10000, math.MaxInt64, 0)
	if _, e := db.Exec("UPDATE oxygen_tanks SET units=? WHERE name='b'", int64(math.MaxInt64-10000)); e != nil {
		t.Fatal(e)
	}
	if e := MoveOxygen(ctx, db, OxygenMove{"boundary", "a", "b", 10000}, nil); e != nil {
		t.Fatal("exact capacity", e)
	}
	hiddenOxygenBalance(t, db, 0, math.MaxInt64, 1)
}
`

func init() {
	add(&Level{ID: "29-4-pool", Chapter: 29, Kind: "lab",
		Title: "حوض الاتصالات", Goal: "اضبط حوض قارئ بوابات الترحيل وحرر الصفوف حتى عند الفشل", Setup: seededLabRoom,
		Requires: []string{"sel:sql.DB", "defer"}, Deps: []string{"modernc.org/sqlite"}, Bridge: []int{68}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch294Starter}},
		Lesson: `# حوض محدود لبوابات الترحيل

## قبل ما تكتب
- هل قفل نتيجة الاستعلام بيقفل القاعدة كلها؟
- التست يثبت الزحمة إزاي من غير افتراض سرعة SQL؟

## القاعدة حوض اتصالات
~*sql.DB~ آمنة للمشاركة بين goroutines؛ مش اتصال واحد. ~SetMaxOpenConns(n)~ بتحط حد للاتصالات المتزامنة. الصفر معناه غير محدود، فالمهمة بترفضه. ~SetMaxIdleConns(n)~ بتحدد الاتصالات الفاضية المحتفظ بيها، مش عدد الاستعلامات الجارية. اضبط قبل تشغيل العمال.
~~~go
db.SetMaxOpenConns(2)
db.SetMaxIdleConns(2)
~~~
ارفض الإعداد الغلط كله قبل أي setter؛ وإلا طلب مرفوض ممكن يغيّر حالة الحوض. المعمل بيستخدم ملف مؤقت عشان الاتصالات تشوف نفس القاعدة. مفيش سيرفر ولا socket.

## الإحصائيات لقطات وعدّادات
~db.Stats()~ بترجع ~sql.DBStats~. ~MaxOpenConnections~ الحد المضبوط، و ~OpenConnections~ العدد الحالي، و ~InUse~ الاتصالات المشغولة، و ~Idle~ الفاضية دلوقتي. ~Idle~ مش الحد الأقصى المضبوط للفاضي.
~~~go
var before sql.DBStats = db.Stats()
// العمال يحجزوا اتصالات هنا
after := db.Stats()
waits := after.WaitCount - before.WaitCount
_ = waits
~~~
~WaitCount~ تراكمي؛ قارن الفرق بدل افتراض البداية صفر. ~WaitDuration~ وقت حقيقي متراكم ومش معيار نجاح حتمي. اللقطة ممكن تتغير بعد قرايتها فوراً.

~~~go
stats := db.Stats()
fmt.Println(stats.MaxOpenConnections, stats.OpenConnections, stats.InUse, stats.Idle)
fmt.Println(stats.WaitDuration) // للملاحظة بس؛ مش مقارنة بقيمة ثابتة.
~~~

## إثبات الزحمة بملكية المورد
~db.Conn(ctx)~ بتحجز ~*sql.Conn~ لحد ما ~conn.Close()~ ترجعها للحوض. التستات بتحجز العدد المضبوط كله، وتبدأ قارئ، وتلاحظ زيادة عداد الانتظار، وبعدين تحرر اتصال. ده حاجز بيتحكم فيه امتلاك المورد، مش نوم.
~~~go
conn, err := db.Conn(ctx)
if err != nil { return err }
defer conn.Close()
~~~
~QueryContext~ المستنية لازم تستجيب لإلغاء المنادي. تحرير الاتصال المحجوز يسمح لاستعلام تاني يكمل. كل goroutine ليها سياق تست محدود؛ التنظيف بيملك الاتصالات والقاعدة.

## الصفوف بتحجز سعة كمان
قارئ 29-1 لازم يقفل الصفوف عند خطأ القراية زي النجاح. خلية بايظة ممكن تفشّل ~Scan~ قبل نهاية اللف، فالوصول للنهاية مش كفاية للتنظيف.
~~~go
rows, err := db.QueryContext(ctx, "SELECT code FROM relay_ports WHERE code>=? ORDER BY code", 10)
if err != nil { return err }
defer rows.Close()
for rows.Next() {
 var code int
 if err := rows.Scan(&code); err != nil { return err }
}
return rows.Err()
~~~

## أخطاء شائعة
~~~go
db.SetMaxOpenConns(0) // غير محدود، مش متوقف
rows, _ := db.QueryContext(ctx, "SELECT code FROM relay_ports")
if rows.Next() { return nil } // رجوع بدري بيسرّب سعة من غير Close مؤجل
~~~

## المهمة
نفّذ ~ConfigureRelayPool(db *sql.DB, max int) error~. ارفض ~max <= 0~ من غير تغيير الإعداد. غير كده اضبط حد الاتصالات المفتوحة وحد الفاضية على ~max~. المنادي بيوفر حوض مش ~nil~ وبيضبطه قبل الاستخدام؛ متفتحش ولا تقفل القاعدة بنفسك.

نفّذ ~ReadRelayPorts(ctx context.Context, db *sql.DB, min int) ([]int, error)~. المنادي بيوفر ~relay_ports(code)~؛ اختار الأكواد الأكبر من أو المساوية للحد المربوط كقيمة، ورتّب تصاعدي مع الحفاظ على التكرار. النتيجة الفاضية slice مش ~nil~. عند خطأ الاستعلام أو القراية أو نهاية اللف ارجع ~nil,error~ مع الحفاظ على الإلغاء. اقفل الصفوف دايماً، حتى عند خطأ ~Scan~ المبكر. التستات بتضيف خلية نص مش قابلة للتحويل لعدد، وبتحجز سعة حقيقية، وبتلغي المنتظرين، وبتثبت إعادة استخدام الاتصال. سيب ~main~ فاضية. التنفيذ مش محتاج ~Report~ أو شبكة أو نوم أو ساعة أو إعدادات جهاز.
`,
		Hints: []string{"ارفض max غير الموجب قبل أي setter؛ اضبط SetMaxOpenConns وSetMaxIdleConns على نفس الحد.", "استخدم QueryContext وWHERE code>=? وORDER BY code؛ خلّي النتيجة الفاضية slice مش nil.", "حط defer rows.Close فور نجاح QueryContext وافحص كل Scan وبعد اللف rows.Err؛ الفشل يرجع nil slice."},
		En: &LevelText{Title: "Connection Pool", Goal: "Bound the relay-port reader pool and release rows even on failure", Files: []LabFile{{Name: "main.go", Starter: ch294Starter}},
			Lesson: `# A bounded pool for relay ports

## Before you write
- Does closing a result set also close the entire database?
- How can a test prove contention without assuming how fast SQL runs?

## A DB is a pool
~*sql.DB~ is safe to share among goroutines; it is not one connection. ~SetMaxOpenConns(n)~ bounds simultaneous connections. Zero means unlimited, so this mission rejects it. ~SetMaxIdleConns(n)~ controls how many unused connections may be retained, not how many may execute queries. Configure before workers start.
~~~go
db.SetMaxOpenConns(2)
db.SetMaxIdleConns(2)
~~~
Reject the whole invalid configuration before either setter; otherwise a rejected change could still mutate the pool. This lab uses a temporary file so connections share the same database. No server or socket is involved.

## Stats are snapshots and counters
~db.Stats()~ returns ~sql.DBStats~. ~MaxOpenConnections~ is the configured open limit, ~OpenConnections~ is the current number, ~InUse~ is occupied connections, and ~Idle~ is unused connections right now. ~Idle~ is not the configured maximum idle count.
~~~go
var before sql.DBStats = db.Stats()
// Workers acquire connections here.
after := db.Stats()
waits := after.WaitCount - before.WaitCount
_ = waits
~~~
~WaitCount~ is cumulative: compare deltas instead of requiring an initial zero. ~WaitDuration~ measures accumulated real time and is not a deterministic success criterion. A snapshot may change as soon as it is read.

~~~go
stats := db.Stats()
fmt.Println(stats.MaxOpenConnections, stats.OpenConnections, stats.InUse, stats.Idle)
fmt.Println(stats.WaitDuration) // Observation only; never an exact assertion.
~~~

## Prove contention using ownership
~db.Conn(ctx)~ reserves one ~*sql.Conn~ until ~conn.Close()~ returns it to the pool. The tests hold exactly the configured number, start a reader, observe an increased wait count, then release one. This is a barrier controlled by resource ownership, not a sleep.
~~~go
conn, err := db.Conn(ctx)
if err != nil { return err }
defer conn.Close()
~~~
A waiting ~QueryContext~ must observe its caller's cancellation. Releasing held connections should let another query proceed. Every goroutine gets a bounded test context; test cleanup owns connections and the database.

## Rows hold capacity too
The reader from 29-1 must close rows on scan failures as well as successful iteration. A malformed cell can fail ~Scan~ before the iterator finishes, so reaching the end cannot be the only cleanup mechanism.
~~~go
rows, err := db.QueryContext(ctx, "SELECT code FROM relay_ports WHERE code>=? ORDER BY code", 10)
if err != nil { return err }
defer rows.Close()
for rows.Next() {
 var code int
 if err := rows.Scan(&code); err != nil { return err }
}
return rows.Err()
~~~

## Common mistakes
~~~go
db.SetMaxOpenConns(0) // Unlimited, not stopped
rows, _ := db.QueryContext(ctx, "SELECT code FROM relay_ports")
if rows.Next() { return nil } // Early return leaks capacity without deferred Close
~~~

## Your mission
Implement ~ConfigureRelayPool(db *sql.DB, max int) error~. Reject ~max <= 0~ without changing configuration. Otherwise call both open and idle limit setters with ~max~. The caller supplies a non-nil pool and configures it before use; do not open or close the database yourself.

Implement ~ReadRelayPorts(ctx context.Context, db *sql.DB, min int) ([]int, error)~. The caller supplies ~relay_ports(code)~; select codes greater than or equal to the bound minimum, ordered ascending, preserving duplicates. Return a non-nil empty slice for no rows. On query, scan, or terminal iteration failure return ~nil,error~, preserving cancellation. Always close rows, including early scan failure. Tests deliberately insert a text cell that cannot scan into an integer, hold real pool capacity, cancel waiters, and verify connection reuse. Keep ~main~ empty. No ~Report~, network, sleep, clock reads, or machine configuration belong in your implementation.
`,
			Hints: []string{"Reject nonpositive max before any setter; set SetMaxOpenConns and SetMaxIdleConns to that limit.", "Use QueryContext with WHERE code>=? and ORDER BY code; return a non-nil empty slice.", "Defer rows.Close immediately after QueryContext succeeds, check each Scan and final rows.Err; return a nil slice on failure."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch294Hidden},
	})
}

const ch294Starter = `package main

import (
	"context"
	"database/sql"
	"errors"
)

func ConfigureRelayPool(db *sql.DB, max int) error { return errors.New("TODO") }
func ReadRelayPorts(ctx context.Context, db *sql.DB, min int) ([]int, error) {
	return nil, errors.New("TODO")
}
func main() {}
`

const ch294Hidden = `package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"math/rand"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func hiddenRelayDB(t *testing.T) *sql.DB {
	t.Helper()
	db, e := sql.Open("sqlite", filepath.Join(t.TempDir(), "relay.db"))
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	if _, e = db.Exec("CREATE TABLE relay_ports(code)"); e != nil {
		t.Fatal(e)
	}
	return db
}
func hiddenWaitForPool(t *testing.T, db *sql.DB, baseline int64) {
	t.Helper()
	deadline := time.NewTimer(3 * time.Second)
	defer deadline.Stop()
	for db.Stats().WaitCount <= baseline {
		select {
		case <-deadline.C:
			t.Fatal("query did not wait for the occupied pool")
		default:
			runtime.Gosched()
		}
	}
}
func TestHiddenRelayBoundAndWait(t *testing.T) {
	db := hiddenRelayDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	limit := rand.New(rand.NewSource(2904)).Intn(3) + 1
	if e := ConfigureRelayPool(db, limit); e != nil {
		t.Fatal(e)
	}
	for _, bad := range []int{0, -1} {
		if e := ConfigureRelayPool(db, bad); e == nil || db.Stats().MaxOpenConnections != limit {
			t.Fatal("invalid configuration changed pool", e)
		}
	}
	held := []*sql.Conn{}
	defer func() {
		for _, c := range held {
			c.Close()
		}
	}()
	for i := 0; i < limit; i++ {
		c, e := db.Conn(ctx)
		if e != nil {
			t.Fatal(e)
		}
		held = append(held, c)
	}
	baseline := db.Stats()
	if baseline.InUse != limit || baseline.OpenConnections > limit || baseline.MaxOpenConnections != limit {
		t.Fatal("ceiling", baseline)
	}
	done := make(chan error, 1)
	go func() { _, e := ReadRelayPorts(ctx, db, 0); done <- e }()
	hiddenWaitForPool(t, db, baseline.WaitCount)
	if s := db.Stats(); s.OpenConnections > limit || s.InUse > limit {
		t.Fatal("exceeded ceiling", s)
	}
	held[0].Close()
	select {
	case e := <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-ctx.Done():
		t.Fatal("release did not unblock reader")
	}
	for _, c := range held {
		c.Close()
	}
	held = nil
	if db.Stats().Idle != limit {
		t.Fatal("idle retention limit not applied", db.Stats())
	}
	if e := ConfigureRelayPool(db, 1); e != nil {
		t.Fatal(e)
	}
	c, e := db.Conn(ctx)
	if e != nil {
		t.Fatal(e)
	}
	defer c.Close()
	waiting, stop := context.WithCancel(ctx)
	defer stop()
	base := db.Stats().WaitCount
	go func() { _, e := ReadRelayPorts(waiting, db, 0); done <- e }()
	hiddenWaitForPool(t, db, base)
	stop()
	select {
	case e := <-done:
		if !errors.Is(e, context.Canceled) {
			t.Fatal("wait cancel", e)
		}
	case <-ctx.Done():
		t.Fatal("canceled waiter stuck")
	}
	c.Close()
	if _, e := ReadRelayPorts(ctx, db, 0); e != nil {
		t.Fatal("pool reuse", e)
	}
}
func TestHiddenRelayRowsSuccessAndScanFailure(t *testing.T) {
	db := hiddenRelayDB(t)
	if e := ConfigureRelayPool(db, 1); e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	got, e := ReadRelayPorts(ctx, db, 0)
	if e != nil || got == nil || len(got) != 0 {
		t.Fatal("empty", got, e)
	}
	rng := rand.New(rand.NewSource(2940))
	want := []int{}
	for i := 0; i < 30; i++ {
		want = append(want, i*3)
	}
	for _, i := range rng.Perm(len(want)) {
		if _, e := db.ExecContext(ctx, "INSERT INTO relay_ports VALUES(?)", want[i]); e != nil {
			t.Fatal(e)
		}
	}
	got, e = ReadRelayPorts(ctx, db, 0)
	if e != nil || !reflect.DeepEqual(got, want) || db.Stats().InUse != 0 {
		t.Fatal("ports/rows", got, e)
	}
	got, e = ReadRelayPorts(ctx, db, 45)
	if e != nil || !reflect.DeepEqual(got, want[15:]) {
		t.Fatal("inclusive minimum", got, e)
	}
	if _, e = db.ExecContext(ctx, "INSERT INTO relay_ports VALUES('broken')"); e != nil {
		t.Fatal(e)
	}
	if got, e = ReadRelayPorts(ctx, db, 0); e == nil || got != nil {
		t.Fatal("scan failure", got, e)
	}
	if db.Stats().InUse != 0 {
		t.Fatal("scan failure leaked rows")
	}
	if _, e = db.ExecContext(ctx, "DELETE FROM relay_ports WHERE code='broken'"); e != nil {
		t.Fatal("pool stuck", e)
	}
	canceled, stop := context.WithCancel(ctx)
	stop()
	if _, e := ReadRelayPorts(canceled, db, 0); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if _, e = db.ExecContext(ctx, "DROP TABLE relay_ports"); e != nil {
		t.Fatal(e)
	}
	if _, e := ReadRelayPorts(ctx, db, 0); e == nil {
		t.Fatal("query failure swallowed")
	}
}

var hiddenRelayFault = errors.New("relay cursor failed")

type hiddenRelayDriver struct{}
type hiddenRelayConn struct{}
type hiddenRelayRows struct{ read bool }

func (hiddenRelayDriver) Open(string) (driver.Conn, error)  { return hiddenRelayConn{}, nil }
func (hiddenRelayConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unused Prepare") }
func (hiddenRelayConn) Close() error                        { return nil }
func (hiddenRelayConn) Begin() (driver.Tx, error)           { return nil, errors.New("unused Begin") }
func (hiddenRelayConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &hiddenRelayRows{}, nil
}
func (*hiddenRelayRows) Columns() []string { return []string{"code"} }
func (*hiddenRelayRows) Close() error      { return nil }
func (r *hiddenRelayRows) Next(dest []driver.Value) error {
	if r.read {
		return hiddenRelayFault
	}
	r.read = true
	dest[0] = int64(7)
	return nil
}
func init() { sql.Register("hidden-relay-fault", hiddenRelayDriver{}) }
func TestHiddenRelayTerminalError(t *testing.T) {
	db, e := sql.Open("hidden-relay-fault", "")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	got, e := ReadRelayPorts(context.Background(), db, 0)
	if !errors.Is(e, hiddenRelayFault) || got != nil {
		t.Fatal("partial cursor failure reported as success", got, e)
	}
	if db.Stats().InUse != 0 {
		t.Fatal("failed cursor retained connection")
	}
}
`

func init() {
	add(&Level{ID: "29-5-sqlx", Chapter: 29, Kind: "lab",
		Title: "sqlx", Goal: "افصل شحن المؤن عن repository حقيقي بـ named queries وstruct scan على Postgres", Setup: seededLabRoom,
		Requires: []string{"iface", "sel:sqlx.NamedExecContext", "sel:sqlx.SelectContext"}, Deps: []string{"github.com/jmoiron/sqlx", "github.com/jackc/pgx/v5"}, Bridge: []int{62}, Race: true, Service: "postgres",
		Files: []LabFile{{Name: "main.go", Starter: ch295Starter}},
		Lesson: `# مؤن بأسماء مع sqlx

## قبل ما تكتب
- إيه السلوك اللي مخزن وهمي يثبته، وإيه اللي محتاج PostgreSQL؟
- مين بيختار نص الاتصال وبينضف مخطط التست؟

## sqlx بتسيب SQL صريحة
استورد ~github.com/jmoiron/sqlx~. ~*sqlx.DB~ بتغلّف ~*sql.DB~ وبتحتفظ بسلوك السياق والحوض، مع مساعدات إضافية. استورد ~github.com/jackc/pgx/v5/stdlib~ لتسجيل الدرايفر لما تفتح PostgreSQL. مسار الـ module هو ~github.com/jackc/pgx/v5~.
~~~go
db, err := sqlx.Open("pgx", os.Getenv("ORBITAL_PG_DSN"))
if err != nil { return err }
defer db.Close()
if err := db.PingContext(ctx); err != nil { return err }
~~~
المخزن عندنا بياخد قاعدة مفتوحة. تجهيز الاتصال في التست بس بيقرا ~ORBITAL_PG_DSN~. تستات الوحدة مش محتاجة خدمة؛ تستات التكامل بتعمل skip لو المتغير غايب. كل جولة تكامل بتنشىء وتحذف schema باسم فريد بتاعها، مش جداول مشتركة.

## اربط أسماء بدل عدّ الأماكن
~sqlx.NamedExecContext~ بتاخد placeholders مسماة وstruct أو map. الـ tag المسمى ~db~ بيربط حقل Go باسم عمود SQL. sqlx بتحوّل ~:name~ لصيغة الدرايفر؛ PostgreSQL بتستخدم ~$1~ و ~$2~ بدل ~?~ بتاعة SQLite.
~~~go
type Beacon struct {
 Name string ` + "`" + `db:"name"` + "`" + `
 Power int ` + "`" + `db:"power"` + "`" + `
}
_, err := sqlx.NamedExecContext(ctx, db,
 "INSERT INTO beacons(name,power) VALUES(:name,:power)", Beacon{"Vega", 0})
~~~
الأسماء دي قيم مربوطة، مش أسماء جداول. متستخدمش named parameter مكان اسم الجدول ومتجمعش SQL مع اسم خليج من المستخدم. فشل القيد يفضل خطأ قاعدة.

## قراية struct
~sqlx.SelectContext~ بتقرا الصفوف في slice من structs بنفس الـ tags. بتقفل الصفوف وبتفحص أخطاء اللف داخلياً. ~sqlx.GetContext~ المقابل لصف واحد وبترجع ~sql.ErrNoRows~ لو غايب.
~~~go
out := []Beacon{}
err := sqlx.SelectContext(ctx, db, &out,
 "SELECT name,power FROM beacons WHERE power >= $1 ORDER BY name", 0)
var one Beacon
err = sqlx.GetContext(ctx, db, &one, "SELECT name,power FROM beacons WHERE name=$1", "Vega")
~~~
هيّئ الـ slice لو العقد واعد بقائمة فاضية مش ~nil~. عند الفشل ارجع ~nil,error~ عشان البيانات المقروءة جزئياً متظهرش كنجاح.

## اختبر الحد بصدق
~RationStore~ واجهة المنادي الصغيرة: ~Put(context.Context, Ration) error~ و ~List(context.Context, int) ([]Ration, error)~. نقدر نختبر ~StockRation(ctx, store, Ration{Bay: " east ", Units: 0})~ بتنفيذ في الذاكرة عشان نثبت إزالة المسافات والتحقق والتعامل مع الخطأ. ده مش دليل على ربط الأسماء ولا قيود PostgreSQL؛ التكامل بيختبر ~SQLRations~ الحقيقي.

## أخطاء شائعة
~~~go
sqlx.NamedExecContext(ctx, db, "INSERT INTO beacons(name) VALUES(':name')", in) // placeholder جوه نص
sqlx.SelectContext(context.Background(), db, &out, "SELECT * FROM beacons") // سياق المنادي ضاع
~~~

## المهمة
احتفظ بأنواع الـ starter وسيب ~main~ فاضية. نفّذ ~StockRation(ctx context.Context, store RationStore, r Ration) error~: افحص السياق الأول، وشيل المسافات من ~Bay~ بـ ~strings.TrimSpace~، وارفض الاسم الفاضي أو ~Units~ السالبة بـ ~ErrRation~؛ غير كده نادي ~store.Put~ مرة بنفس السياق والقيمة المطبعة. صفر وحدات مسموح. حافظ على أخطاء المخزن ومتندهش له عند فشل التحقق أو الإلغاء.

~NewSQLRations(db *sqlx.DB) *SQLRations~ بتغلّف الحوض المحقون. ~InitRations(ctx, db)~ تنشئ ~ration_bins(bay TEXT PRIMARY KEY CHECK(length(trim(bay))>0), units INTEGER NOT NULL CHECK(units>=0))~ بشكل قابل للتكرار. ~Put~ تضيف الحقول زي ما هي باستخدام ~sqlx.NamedExecContext~ وtags ~db~؛ من غير تطبيع، وبترجع أخطاء التكرار والقيود. ~List(ctx, min)~ تستخدم ~sqlx.SelectContext~ عشان ترجع الخلجان اللي وحداتها على الأقل ~min~، مرتبة بالاسم، بكل الحقول، وslice فاضية مش ~nil~. عند الخطأ ارجع ~nil,error~. مرّر السياق لكل استعلام وحافظ على الأخطاء وسيب عمر القاعدة للمنادي. التكامل على PostgreSQL حقيقي من ~ORBITAL_PG_DSN~؛ تستات المخزن الوهمي بتشتغل دايماً من غيره. مفيش نداءات شبكة صريحة في المخزن.
`,
		Hints: []string{"خلّي StockRation تتحقق من السياق والاسم والوحدات قبل نداء المخزن المحقون.", "اعمل repository شايلة *sqlx.DB؛ اربط :bay و:units بـ sqlx.NamedExecContext وdb tags.", "استخدم sqlx.SelectContext على slice مهيّأة وplaceholder $1، وسيب التست يعمل schema معزولة من ORBITAL_PG_DSN."},
		En: &LevelText{Title: "sqlx", Goal: "Separate ration stocking from a real Postgres repository using named queries and struct scanning", Files: []LabFile{{Name: "main.go", Starter: ch295Starter}},
			Lesson: `# Named supplies with sqlx

## Before you write
- Which behavior can a fake store prove, and which needs PostgreSQL?
- Who chooses the database connection string and cleans up the test schema?

## sqlx keeps SQL explicit
Import ~github.com/jmoiron/sqlx~. ~*sqlx.DB~ wraps ~*sql.DB~, retaining its context and pool behavior while adding helpers. Import ~github.com/jackc/pgx/v5/stdlib~ for its driver registration when opening PostgreSQL. The module dependency is ~github.com/jackc/pgx/v5~.
~~~go
db, err := sqlx.Open("pgx", os.Getenv("ORBITAL_PG_DSN"))
if err != nil { return err }
defer db.Close()
if err := db.PingContext(ctx); err != nil { return err }
~~~
Our repository receives an already opened database. Only the test connection setup reads ~ORBITAL_PG_DSN~. Unit tests need no service; integration tests skip when the variable is absent. Each integration run creates and drops its own uniquely named schema, never shared tables.

## Bind names instead of counting positions
~sqlx.NamedExecContext~ accepts named placeholders and a struct or map. A ~db~ tag maps a Go field to its SQL name. sqlx rewrites ~:name~ placeholders to the driver's format; PostgreSQL positional parameters use ~$1~, ~$2~, rather than SQLite's ~?~.
~~~go
type Beacon struct {
 Name string ` + "`" + `db:"name"` + "`" + `
 Power int ` + "`" + `db:"power"` + "`" + `
}
_, err := sqlx.NamedExecContext(ctx, db,
 "INSERT INTO beacons(name,power) VALUES(:name,:power)", Beacon{"Vega", 0})
~~~
Names are still bound values, not table identifiers. Do not use named parameters to interpolate a table name. Do not build SQL by concatenating a bay name. Constraint failures remain database errors.

## Struct scanning
~sqlx.SelectContext~ scans every row into a slice of structs, using those same tags. It closes rows and checks iteration errors internally. ~sqlx.GetContext~ is its single-row counterpart and returns ~sql.ErrNoRows~ for absence.
~~~go
out := []Beacon{}
err := sqlx.SelectContext(ctx, db, &out,
 "SELECT name,power FROM beacons WHERE power >= $1 ORDER BY name", 0)
var one Beacon
err = sqlx.GetContext(ctx, db, &one, "SELECT name,power FROM beacons WHERE name=$1", "Vega")
~~~
Keep the empty slice initialized if the contract promises ~[]~ rather than ~nil~. Return ~nil,error~ on failure so partially scanned data cannot look like a successful result.

## Test the boundary honestly
~RationStore~ is the consumer's small interface: ~Put(context.Context, Ration) error~ and ~List(context.Context, int) ([]Ration, error)~. ~StockRation(ctx, store, Ration{Bay: " east ", Units: 0})~ can be tested with an in-memory implementation to verify trimming, validation, and error handling. That fake does not prove named binding or PostgreSQL constraints: integration tests exercise the concrete ~SQLRations~ repository.

## Common mistakes
~~~go
sqlx.NamedExecContext(ctx, db, "INSERT INTO beacons(name) VALUES(':name')", in) // Quoted placeholder
sqlx.SelectContext(context.Background(), db, &out, "SELECT * FROM beacons") // Caller cancellation lost
~~~

## Your mission
Keep the starter types and empty ~main~. Implement ~StockRation(ctx context.Context, store RationStore, r Ration) error~: check context first, trim ~Bay~ with ~strings.TrimSpace~, reject empty bay or negative ~Units~ as ~ErrRation~, and otherwise call ~store.Put~ once with the normalized value and original context. Zero units is valid. Preserve store errors and do not call the store on validation/cancellation failure.

~NewSQLRations(db *sqlx.DB) *SQLRations~ wraps the injected pool. ~InitRations(ctx, db)~ idempotently creates ~ration_bins(bay TEXT PRIMARY KEY CHECK(length(trim(bay))>0), units INTEGER NOT NULL CHECK(units>=0))~. ~Put~ inserts exact supplied fields using ~sqlx.NamedExecContext~ and ~db~ tags; it performs no normalization and returns duplicate/constraint errors. ~List(ctx, min)~ uses ~sqlx.SelectContext~ to return bays with units at least ~min~, ordered by bay, with all fields populated and a non-nil empty slice. On error return ~nil,error~. Pass context through every query, preserve errors, and leave database lifetime to the caller. Integration uses real PostgreSQL via ~ORBITAL_PG_DSN~; fake-store tests always run without it. No network calls belong in the repository itself.
`,
			Hints: []string{"Have StockRation check context, bay, and units before calling the injected store.", "Keep *sqlx.DB in the repository; bind :bay and :units with sqlx.NamedExecContext and db tags.", "Use sqlx.SelectContext into an initialized slice and a $1 placeholder; let tests isolate a schema from ORBITAL_PG_DSN."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch295Hidden},
	})
}

const ch295Starter = `package main

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Ration struct {
	Bay   string ` + "`" + `db:"bay"` + "`" + `
	Units int    ` + "`" + `db:"units"` + "`" + `
}
type RationStore interface {
	Put(context.Context, Ration) error
	List(context.Context, int) ([]Ration, error)
}

var ErrRation = errors.New("invalid ration")

func StockRation(ctx context.Context, store RationStore, r Ration) error { return errors.New("TODO") }

type SQLRations struct{ db *sqlx.DB }

func NewSQLRations(db *sqlx.DB) *SQLRations                   { return &SQLRations{db: db} }
func InitRations(ctx context.Context, db *sqlx.DB) error      { return errors.New("TODO") }
func (s *SQLRations) Put(ctx context.Context, r Ration) error { return errors.New("TODO") }
func (s *SQLRations) List(ctx context.Context, min int) ([]Ration, error) {
	return nil, errors.New("TODO")
}
func main() {}
`

const ch295Hidden = `package main

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type hiddenRationStore struct {
	items []Ration
	err   error
	ctx   context.Context
}

func (s *hiddenRationStore) Put(ctx context.Context, r Ration) error {
	s.ctx = ctx
	if s.err != nil {
		return s.err
	}
	s.items = append(s.items, r)
	return nil
}
func (s *hiddenRationStore) List(ctx context.Context, min int) ([]Ration, error) {
	return nil, errors.New("unused fake read")
}
func TestHiddenRationsUnit(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, "identity")
	fake := &hiddenRationStore{}
	rng := rand.New(rand.NewPCG(29, 5))
	for i := 0; i < 25; i++ {
		r := Ration{fmt.Sprintf(" bay-%d ", rng.Int64()), rng.IntN(200)}
		if e := StockRation(ctx, fake, r); e != nil {
			t.Fatal(e)
		}
		want := r
		want.Bay = strings.TrimSpace(want.Bay)
		if fake.items[i] != want || fake.ctx != ctx {
			t.Fatal("normalization/context", fake.items[i], fake.ctx)
		}
	}
	for _, r := range []Ration{{"", 0}, {" \t\n", 1}, {"bay", -1}} {
		if e := StockRation(ctx, fake, r); !errors.Is(e, ErrRation) {
			t.Fatal("invalid", e)
		}
	}
	if len(fake.items) != 25 {
		t.Fatal("invalid input stored")
	}
	if e := StockRation(ctx, fake, Ration{"zero", 0}); e != nil {
		t.Fatal("zero boundary", e)
	}
	sentinel := errors.New("storage offline")
	fake.err = sentinel
	if e := StockRation(ctx, fake, Ration{"a", 1}); !errors.Is(e, sentinel) {
		t.Fatal("store failure", e)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if e := StockRation(canceled, fake, Ration{"a", 1}); !errors.Is(e, context.Canceled) {
		t.Fatal("cancel before store", e)
	}
}
func hiddenRationPostgres(t *testing.T) *sqlx.DB {
	t.Helper()
	dsn := os.Getenv("ORBITAL_PG_DSN")
	if dsn == "" {
		t.Skip("set ORBITAL_PG_DSN for isolated Postgres integration")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cfg, e := pgx.ParseConfig(dsn)
	if e != nil {
		t.Fatal("invalid Postgres configuration")
	}
	admin := sqlx.NewDb(stdlib.OpenDB(*cfg), "pgx")
	t.Cleanup(func() { admin.Close() })
	var token [16]byte
	if _, e := cryptorand.Read(token[:]); e != nil {
		t.Fatal(e)
	}
	schema := "orbital_rations_" + hex.EncodeToString(token[:])
	if _, e = admin.ExecContext(ctx, "CREATE SCHEMA "+schema); e != nil {
		t.Fatal("create isolated schema:", e)
	}
	t.Cleanup(func() {
		clean, stop := context.WithTimeout(context.Background(), 10*time.Second)
		defer stop()
		if _, e := admin.ExecContext(clean, "DROP SCHEMA "+schema+" CASCADE"); e != nil {
			t.Error("drop owned schema:", e)
		}
	})
	scoped := cfg.Copy()
	scoped.RuntimeParams["search_path"] = schema
	db := sqlx.NewDb(stdlib.OpenDB(*scoped), "pgx")
	db.SetMaxOpenConns(4)
	t.Cleanup(func() { db.Close() })
	return db
}
func TestHiddenRationsPostgresIntegration(t *testing.T) {
	db := hiddenRationPostgres(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if e := InitRations(ctx, db); e != nil {
		t.Fatal(e)
	}
	if e := InitRations(ctx, db); e != nil {
		t.Fatal(e)
	}
	repo := NewSQLRations(db)
	empty, e := repo.List(ctx, 0)
	if e != nil || empty == nil || len(empty) != 0 {
		t.Fatal("empty", empty, e)
	}
	rng := rand.New(rand.NewPCG(295, 29))
	want := []Ration{}
	for i := 0; i < 20; i++ {
		want = append(want, Ration{fmt.Sprintf("bay-%02d-'--", i), rng.IntN(200)})
	}
	var wg sync.WaitGroup
	failures := make(chan error, len(want))
	for _, r := range want {
		wg.Add(1)
		go func(r Ration) { defer wg.Done(); failures <- StockRation(ctx, repo, r) }(r)
	}
	wg.Wait()
	close(failures)
	for e := range failures {
		if e != nil {
			t.Fatal(e)
		}
	}
	got, e := repo.List(ctx, 0)
	if e != nil || !reflect.DeepEqual(got, want) {
		t.Fatal("struct scan/order", got, e)
	}
	filtered := []Ration{}
	for _, r := range want {
		if r.Units >= 100 {
			filtered = append(filtered, r)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Bay < filtered[j].Bay })
	got, e = repo.List(ctx, 100)
	if e != nil || !reflect.DeepEqual(got, filtered) {
		t.Fatal("filter", got, e)
	}
	if e := repo.Put(ctx, want[0]); e == nil {
		t.Fatal("duplicate accepted")
	}
	for _, r := range []Ration{{" ", 2}, {"negative", -1}} {
		if e := repo.Put(ctx, r); e == nil {
			t.Fatal("database constraint", r)
		}
	}
	canceled, stop := context.WithCancel(ctx)
	stop()
	if _, e := repo.List(canceled, 0); !errors.Is(e, context.Canceled) {
		t.Fatal("read cancellation", e)
	}
	if e := repo.Put(canceled, Ration{"canceled", 1}); !errors.Is(e, context.Canceled) {
		t.Fatal("write cancellation", e)
	}
	db.Close()
	if _, e := repo.List(ctx, 0); e == nil {
		t.Fatal("closed storage swallowed")
	}
}
`

func init() {
	add(&Level{ID: "29-6-gorm", Chapter: 29, Kind: "lab",
		Title: "GORM", Goal: "اربط قباب الزراعة بأحواضها وحمّل العلاقات من غير N+1", Setup: seededLabRoom,
		Requires: []string{"sel:gorm.DB"}, Deps: []string{"gorm.io/gorm", "github.com/glebarez/sqlite"}, Bridge: []int{63}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch296Starter}},
		Lesson: `# قباب الزراعة وأحواضها

## قبل ما تكتب
- هل القبة ممكن تفضل محفوظة لو إضافة حوض فشلت؟
- محتاج كام استعلام عشان تحمّل خمس قباب بكل أحواضها؟

## الـ models ربط صريح
استورد ~gorm.io/gorm~ و ~github.com/glebarez/sqlite~ لدرايفر SQLite مكتوب بـ Go. ~gorm.Open~ بترجع ~*gorm.DB~، مقبض الـ ORM. المنادي بيضبط الحوض والمفاتيح الخارجية. مفيش حقول وقت هنا، فالموديلات مش محتاجة قراية ساعة تلقائية.
~~~go
db, err := gorm.Open(sqlite.Open("garden.db"), &gorm.Config{})
if err != nil { return err }
pool, err := db.DB()
if err != nil { return err }
defer pool.Close()
pool.SetMaxOpenConns(1)
db.Exec("PRAGMA foreign_keys=ON")
~~~
GORM بتربط حقول struct المصدّرة بأعمدة وأسماء جداول جمع بـ snake case. ~gorm:"primaryKey"~ بتحدد المفتاح. ~not null~ بترفض الغياب، و ~uniqueIndex~ بتفرض الفرادة. ~AutoMigrate~ بتنشئ المخطط الناقص من الموديلات. مناسبة للمعمل؛ مش بديل تاريخ النسخ القابل للرجوع في 29-2.
~~~go
type Rack struct {
 ID uint ` + "`" + `gorm:"primaryKey"` + "`" + `
 Name string ` + "`" + `gorm:"uniqueIndex;not null"` + "`" + `
 Slots []Slot ` + "`" + `gorm:"foreignKey:RackID;references:ID;constraint:OnDelete:RESTRICT"` + "`" + `
}
type Slot struct {
 ID uint ` + "`" + `gorm:"primaryKey"` + "`" + `
 RackID uint ` + "`" + `gorm:"not null;uniqueIndex:rack_slot"` + "`" + `
 Label string ` + "`" + `gorm:"not null;uniqueIndex:rack_slot"` + "`" + `
}
err := db.WithContext(ctx).AutoMigrate(&Rack{}, &Slot{})
~~~
~foreignKey:RackID~ بتشير لحقل الابن، و ~references:ID~ لمفتاح الأب. ~OnDelete:RESTRICT~ بتمنع حذف أب ليه أبناء. استخدام نفس اسم ~uniqueIndex~ على حقلين بيخلي التركيبة فريدة، مش كل عمود لوحده.

## سياق وعلاقات ذرّية
~WithContext(ctx)~ بترجع مقبض استعلام بسياق المنادي. العملية النهائية بتوفر ~Error~؛ متتجاهلوش. ~Create(&rack)~ بتضيف الأب وبتكتب الـ ID المتولّد فيه، وبتحفظ الأبناء بالـ ID ده. ~Transaction~ بتحفظ بس لو الـ callback رجّعت ~nil~؛ ارجع خطأ أي ابن عشان الأب يتراجع كمان.
~~~go
rack := Rack{Name: "Vega", Slots: []Slot{{Label: "north"}}}
err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
 return tx.Create(&rack).Error
})
~~~
افحص كل ابن قبل الكتابة، بس تعامل كمان مع فشل القاعدة أثناء حفظ العلاقات. تكرار أسماء الأبناء ممكن يفشل بعد إضافة الأب؛ نجاح التحقق لوحده مش دليل على الذرّية.

## Preload ومشكلة N+1
~Find~ بتحمّل قايمة، و ~Order~ بيحدد ترتيبها. حلقة بتقرا أبناء كل أب لوحدهم تعمل استعلام آباء ومعاه N استعلام أبناء: N+1. ~Preload~ بتسجل التحميل المسبق قبل ~Find~ النهائية؛ استعلام آباء واحد واستعلام أبناء مجمّع واحد.
~~~go
out := []Rack{}
err := db.WithContext(ctx).Preload("Slots", func(q *gorm.DB) *gorm.DB {
 return q.Order("id ASC")
}).Order("id ASC").Find(&out).Error
~~~
نتيجة آباء مش فاضية تستخدم أمرين بالظبط حتى لو بعضهم ملهوش أبناء. الفاضية تستخدم واحد عشان مفيش مفاتيح آباء للتحميل. التستات بتحقن logger لـ GORM وبتعد نداءات ~Trace~، وبتصفّر العداد قبل القراية عشان تجهيز الجداول ميتحسبش. ~GardenQueryCount(5, true)~ المفروض ترجع 2، و ~GardenQueryCount(5, false)~ ترجع 6. الدالة توثيق قابل للتنفيذ لميزانية الاستعلامات.

## القيم الصفرية في التعديل
~Updates(struct)~ عادة بتتخطى القيم الصفرية ومنها ~false~. استخدم ~Update~ لعمود واحد عشان تكتب false. ~Model~ بتحدد الجدول، و ~Where~ الصف، و ~RowsAffected~ بتكشف الغياب؛ بنمثّله هنا بـ ~gorm.ErrRecordNotFound~.
~~~go
result := db.WithContext(ctx).Model(&GardenBed{}).Where("id = ?", id).Update("active", false)
if result.Error != nil { return result.Error }
if result.RowsAffected == 0 { return gorm.ErrRecordNotFound }
~~~

## أخطاء شائعة
~~~go
db.Find(&domes).Preload("Beds") // متأخر: القراية حصلت
for _, d := range domes { db.Where("dome_id=?", d.ID).Find(&beds) } // N+1
 db.Model(&bed).Updates(GardenBed{Active:false}) // false هتتخطى
~~~

## المهمة
احتفظ بأسماء أنواع وحقول الـ starter. ضيف tags صريحة لـ ~GardenDome~ و ~GardenBed~: مفاتيح ID موجبة النوع unsigned ومتولدة من القاعدة، واسم قبة فريد وغير غايب، و ~DomeID~ و ~Label~ للابن غير غايبين، وفهرس فريد مشترك للحقلين دول، وعلاقة ~Beds~ صريحة من ~DomeID~ لـ ~ID~ مع منع الحذف. الجداول الافتراضية ~garden_domes~ و ~garden_beds~. ~InitGarden(ctx, db)~ تعمل auto-migrate للاتنين بشكل قابل للتكرار. المنادي بيفعّل المفاتيح الخارجية وبيملك مقبض SQLite/GORM المحقون.

~CreateGardenDome(ctx, db, name, labels)~ تشيل مسافات الاسم وكل label وترفض الفاضي بأي خطأ قبل الكتابة. متغيّرش slice المنادي. قائمة labels الفاضية مسموحة. احفظ المجموعة كلها ذرّياً والأحواض نشطة، وارجع المفاتيح المتولدة والخارجية. تكرار اسم قبة أو اسم حوض داخل نفس القبة يفشل من غير صفوف جزئية. عند أي خطأ ارجع ~GardenDome~ صفرية. تكرار label في قبتين مختلفتين مسموح.

~ListGarden(ctx, db)~ ترجع القباب مرتبة بالـ ID بأحواض محمّلة مسبقاً ومرتبة بالـ ID وslices فاضية مش ~nil~ للآباء والأبناء. استخدم أمرين SQL بالظبط للقائمة المش فاضية، وواحد للفاضية. عند الفشل ارجع ~nil,error~. ~GardenQueryCount(domes, preload)~ بتوثّق الميزانية: للأعداد غير السالبة، 1 للصفر، وإلا 2 مع التحميل المسبق أو ~domes+1~ من غيره.

~SetGardenBed(ctx, db, id, active)~ تغيّر نشاط الحوض ده بس، حتى false؛ الغايب يرجع ~gorm.ErrRecordNotFound~. حافظ على باقي الأخطاء ومرّر السياق لكل العمليات. قيود القاعدة ترفض الابن اليتيم وحذف أب عنده أبناء. كل الدوال تدعم استخدام المقبض المحقون بالتزامن. سيب ~main~ فاضية؛ مفيش شبكة ولا ~Report~.
`,
		Hints: []string{"عرّف مفاتيح ID وعلاقة Beds بـ foreignKey:DomeID وreferences:ID وقيد RESTRICT، وفهرس مشترك على DomeID وLabel.", "اتحقق من الاسم وكل labels قبل Create؛ اعمل إنشاء القبة والأحواض في Transaction واحدة.", "في ListGarden استخدم Preload قبل Find ورتّب القباب والأحواض بالـ ID؛ Update لعمود active بيحفظ false كمان."},
		En: &LevelText{Title: "GORM", Goal: "Relate garden domes to beds and load associations without N+1 queries", Files: []LabFile{{Name: "main.go", Starter: ch296Starter}},
			Lesson: `# Garden domes and their beds

## Before you write
- Can a dome remain saved when one child insert fails?
- How many queries load five domes and all their beds?

## Models are an explicit mapping
Import ~gorm.io/gorm~ and ~github.com/glebarez/sqlite~ for a pure-Go SQLite dialect. ~gorm.Open~ returns ~*gorm.DB~, the ORM handle. The caller configures its underlying pool and foreign keys. This lab has no timestamp fields, so no automatic clock readings are needed in your models.
~~~go
db, err := gorm.Open(sqlite.Open("garden.db"), &gorm.Config{})
if err != nil { return err }
pool, err := db.DB()
if err != nil { return err }
defer pool.Close()
pool.SetMaxOpenConns(1)
db.Exec("PRAGMA foreign_keys=ON")
~~~
GORM maps exported struct fields to columns and plural snake-case table names. A ~gorm:"primaryKey"~ tag declares a key. ~not null~ forbids missing values; ~uniqueIndex~ creates a uniqueness rule. ~AutoMigrate~ creates missing schema from models. It is convenient for this lab; it is not the reversible version history built in 29-2.
~~~go
type Rack struct {
 ID uint ` + "`" + `gorm:"primaryKey"` + "`" + `
 Name string ` + "`" + `gorm:"uniqueIndex;not null"` + "`" + `
 Slots []Slot ` + "`" + `gorm:"foreignKey:RackID;references:ID;constraint:OnDelete:RESTRICT"` + "`" + `
}
type Slot struct {
 ID uint ` + "`" + `gorm:"primaryKey"` + "`" + `
 RackID uint ` + "`" + `gorm:"not null;uniqueIndex:rack_slot"` + "`" + `
 Label string ` + "`" + `gorm:"not null;uniqueIndex:rack_slot"` + "`" + `
}
err := db.WithContext(ctx).AutoMigrate(&Rack{}, &Slot{})
~~~
~foreignKey:RackID~ points to the child's field; ~references:ID~ names the parent's key. ~OnDelete:RESTRICT~ prevents deleting a parent with children. Giving two fields the same named ~uniqueIndex~ makes their combination unique, rather than making each column unique alone.

## Context and atomic associations
~WithContext(ctx)~ returns a query handle carrying the caller's cancellation. A terminal operation exposes ~Error~; never ignore it. ~Create(&rack)~ inserts the parent, writes its generated ID back, and saves the association with that ID on each child. ~Transaction~ commits only when its callback returns ~nil~; return any child error to roll back the parent too.
~~~go
rack := Rack{Name: "Vega", Slots: []Slot{{Label: "north"}}}
err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
 return tx.Create(&rack).Error
})
~~~
Validate every child before writing, but also handle database failure during association creation. Duplicate child labels can fail after the parent insert; successful validation alone does not prove atomicity.

## Preload and the N+1 problem
~Find~ loads a list; ~Order~ makes its order explicit. A loop querying children separately makes one parent query plus N child queries: N+1. ~Preload~ registers eager loading before the terminal ~Find~. It loads parents once and children with one grouped query.
~~~go
out := []Rack{}
err := db.WithContext(ctx).Preload("Slots", func(q *gorm.DB) *gorm.DB {
 return q.Order("id ASC")
}).Order("id ASC").Find(&out).Error
~~~
A nonempty parent result uses exactly two statements, even when some parents have no children. An empty result uses one because there are no parent IDs to preload. Tests inject a GORM logger and count ~Trace~ callbacks, resetting immediately before the read so schema setup is excluded. ~GardenQueryCount(5, true)~ should report 2; ~GardenQueryCount(5, false)~ should report 6. This function is executable documentation of the query budget.

## Zero values in updates
~Updates(struct)~ normally skips zero fields, including ~false~. Use a single-column ~Update~ to write a false value. ~Model~ selects the table; ~Where~ selects the row; ~RowsAffected~ detects absence, represented here by ~gorm.ErrRecordNotFound~.
~~~go
result := db.WithContext(ctx).Model(&GardenBed{}).Where("id = ?", id).Update("active", false)
if result.Error != nil { return result.Error }
if result.RowsAffected == 0 { return gorm.ErrRecordNotFound }
~~~

## Common mistakes
~~~go
db.Find(&domes).Preload("Beds") // Too late: Find already ran
for _, d := range domes { db.Where("dome_id=?", d.ID).Find(&beds) } // N+1
 db.Model(&bed).Updates(GardenBed{Active:false}) // False skipped
~~~

## Your mission
Keep starter type and field names. Add explicit tags to ~GardenDome~ and ~GardenBed~: generated unsigned primary IDs, unique non-null dome name, non-null child ~DomeID~ and ~Label~, a joint unique index over those two child fields, and an explicit ~Beds~ association referencing ~DomeID~ → ~ID~ with delete restriction. Default tables are ~garden_domes~ and ~garden_beds~. ~InitGarden(ctx, db)~ auto-migrates both, idempotently. The caller enables foreign keys and owns the injected SQLite/GORM handle.

~CreateGardenDome(ctx, db, name, labels)~ trims name and every label and rejects blanks with any error before writes. Do not mutate the caller's slice. Empty labels are valid. Create the whole aggregate atomically with every bed initially active; return generated IDs and foreign keys. Duplicate dome names or duplicate labels within one dome fail with no partial rows. On any error return a zero ~GardenDome~. Labels may repeat in different domes.

~ListGarden(ctx, db)~ returns all domes ordered by ID, with preloaded beds ordered by ID and non-nil empty slices (parents and child associations). Use exactly two SQL statements for any nonempty list, one when empty. On failure return ~nil,error~. ~GardenQueryCount(domes, preload)~ documents that budget: for nonnegative counts, 1 when zero, otherwise 2 with preload or ~domes+1~ without it.

~SetGardenBed(ctx, db, id, active)~ updates just that bed's active flag, including false; missing IDs return ~gorm.ErrRecordNotFound~. Preserve other errors and pass context through all operations. Database constraints reject orphan children and parent deletion while children exist. All functions support concurrent use of the injected handle. Keep ~main~ empty; no network or ~Report~.
`,
			Hints: []string{"Define ID keys, a Beds relation with foreignKey:DomeID and references:ID, RESTRICT, and a joint DomeID/Label unique index.", "Validate the name and every label before Create; create the dome and beds inside one Transaction.", "In ListGarden register Preload before Find and order both domes and beds by ID; a single-column Update also persists false."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch296Hidden},
	})
}

const ch296Starter = `package main

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type GardenDome struct {
	ID   uint
	Name string
	Beds []GardenBed
}
type GardenBed struct {
	ID     uint
	DomeID uint
	Label  string
	Active bool
}

func InitGarden(ctx context.Context, db *gorm.DB) error { return errors.New("TODO") }
func CreateGardenDome(ctx context.Context, db *gorm.DB, name string, labels []string) (GardenDome, error) {
	return GardenDome{}, errors.New("TODO")
}
func ListGarden(ctx context.Context, db *gorm.DB) ([]GardenDome, error) {
	return nil, errors.New("TODO")
}
func SetGardenBed(ctx context.Context, db *gorm.DB, id uint, active bool) error {
	return errors.New("TODO")
}
func GardenQueryCount(domes int, preload bool) int { return 0 }
func main()                                        {}
`

const ch296Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type hiddenGardenLog struct {
	logger.Interface
	calls atomic.Int64
}

func (l *hiddenGardenLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.calls.Add(1)
}
func hiddenGardenDB(t *testing.T) (*gorm.DB, *hiddenGardenLog) {
	t.Helper()
	log := &hiddenGardenLog{Interface: logger.Discard}
	db, e := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "garden.db")), &gorm.Config{Logger: log, TranslateError: true})
	if e != nil {
		t.Fatal(e)
	}
	pool, e := db.DB()
	if e != nil {
		t.Fatal(e)
	}
	pool.SetMaxOpenConns(1)
	t.Cleanup(func() { pool.Close() })
	if e := db.Exec("PRAGMA foreign_keys=ON").Error; e != nil {
		t.Fatal(e)
	}
	if e := InitGarden(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	if e := InitGarden(context.Background(), db); e != nil {
		t.Fatal("repeat migration", e)
	}
	return db, log
}
func hiddenGardenCounts(t *testing.T, db *gorm.DB, domes, beds int64) {
	t.Helper()
	var a, b int64
	if e := db.Model(&GardenDome{}).Count(&a).Error; e != nil {
		t.Fatal(e)
	}
	if e := db.Model(&GardenBed{}).Count(&b).Error; e != nil {
		t.Fatal(e)
	}
	if a != domes || b != beds {
		t.Fatal("partial association state", a, b, domes, beds)
	}
}
func TestHiddenGardenPreloadAndBudget(t *testing.T) {
	db, log := hiddenGardenDB(t)
	ctx := context.Background()
	log.calls.Store(0)
	empty, e := ListGarden(ctx, db)
	if e != nil || empty == nil || len(empty) != 0 || log.calls.Load() != 1 || GardenQueryCount(0, true) != 1 || GardenQueryCount(0, false) != 1 {
		t.Fatal("empty", empty, e, log.calls.Load())
	}
	rng := rand.New(rand.NewSource(2906))
	n := rng.Intn(4) + 3
	want := []GardenDome{}
	for i := 0; i < n; i++ {
		labels := []string{fmt.Sprintf(" z-%d ", rng.Int63()), fmt.Sprintf(" a-%d ", rng.Int63())}
		before := append([]string(nil), labels...)
		d, e := CreateGardenDome(ctx, db, fmt.Sprintf(" dome-%d ", i), labels)
		if e != nil || d.ID == 0 || d.Name != fmt.Sprintf("dome-%d", i) || len(d.Beds) != 2 {
			t.Fatal("create", d, e)
		}
		if !reflect.DeepEqual(labels, before) {
			t.Fatal("mutated labels")
		}
		for j, b := range d.Beds {
			if b.ID == 0 || b.DomeID != d.ID || !b.Active || b.Label != strings.TrimSpace(labels[j]) {
				t.Fatal("association", b)
			}
		}
		want = append(want, d)
	}
	log.calls.Store(0)
	got, e := ListGarden(ctx, db)
	calls := log.calls.Load()
	if e != nil || !reflect.DeepEqual(got, want) || calls != 2 || GardenQueryCount(n, true) != int(calls) {
		t.Fatal("preload/order/budget", got, e, calls)
	}
	log.calls.Store(0)
	var naive []GardenDome
	if e := db.Order("id ASC").Find(&naive).Error; e != nil {
		t.Fatal(e)
	}
	for _, d := range naive {
		var beds []GardenBed
		if e := db.Where("dome_id=?", d.ID).Find(&beds).Error; e != nil {
			t.Fatal(e)
		}
	}
	if calls := log.calls.Load(); calls != int64(n+1) || GardenQueryCount(n, false) != int(calls) {
		t.Fatal("N+1 report", calls)
	}
	if e := SetGardenBed(ctx, db, want[0].Beds[0].ID, false); e != nil {
		t.Fatal(e)
	}
	got, e = ListGarden(ctx, db)
	if e != nil || got[0].Beds[0].Active || !got[0].Beds[1].Active {
		t.Fatal("false update", got, e)
	}
	if e := SetGardenBed(ctx, db, ^uint(0)>>1, true); !errors.Is(e, gorm.ErrRecordNotFound) {
		t.Fatal("missing", e)
	}
	var wg sync.WaitGroup
	failures := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, e := ListGarden(ctx, db); failures <- e }()
	}
	wg.Wait()
	close(failures)
	for e := range failures {
		if e != nil {
			t.Fatal(e)
		}
	}
}
func TestHiddenGardenAtomicAssociations(t *testing.T) {
	db, _ := hiddenGardenDB(t)
	ctx := context.Background()
	for _, c := range []struct {
		name   string
		labels []string
	}{{" ", nil}, {"bad", []string{"valid", " "}}} {
		if _, e := CreateGardenDome(ctx, db, c.name, c.labels); e == nil {
			t.Fatal("invalid aggregate")
		}
		hiddenGardenCounts(t, db, 0, 0)
	}
	if _, e := CreateGardenDome(ctx, db, "partial", []string{"same", "same"}); e == nil {
		t.Fatal("duplicate children accepted")
	}
	hiddenGardenCounts(t, db, 0, 0)
	d, e := CreateGardenDome(ctx, db, "empty", nil)
	if e != nil {
		t.Fatal(e)
	}
	if _, e := CreateGardenDome(ctx, db, "empty", nil); e == nil {
		t.Fatal("duplicate parent")
	}
	hiddenGardenCounts(t, db, 1, 0)
	if e := db.Create(&GardenBed{DomeID: d.ID + 100, Label: "orphan"}).Error; e == nil {
		t.Fatal("missing foreign key")
	}
	hiddenGardenCounts(t, db, 1, 0)
	parent, e := CreateGardenDome(ctx, db, "occupied", []string{"bed"})
	if e != nil {
		t.Fatal(e)
	}
	if e := db.Delete(&GardenDome{}, parent.ID).Error; e == nil {
		t.Fatal("RESTRICT missing")
	}
	hiddenGardenCounts(t, db, 2, 1)
	got, e := ListGarden(ctx, db)
	if e != nil || len(got) != 2 || got[0].Beds == nil || len(got[0].Beds) != 0 {
		t.Fatal("empty association", got, e)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, e := ListGarden(canceled, db); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if _, e := CreateGardenDome(canceled, db, "canceled", nil); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if e := SetGardenBed(canceled, db, parent.Beds[0].ID, false); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	hiddenGardenCounts(t, db, 2, 1)
}

func TestHiddenGardenLabelsAreScopedToDome(t *testing.T) {
	db, _ := hiddenGardenDB(t)
	ctx := context.Background()
	for _, name := range []string{"east", "west"} {
		if _, e := CreateGardenDome(ctx, db, name, []string{"shared"}); e != nil {
			t.Fatal("labels must be unique per dome", e)
		}
	}
	hiddenGardenCounts(t, db, 2, 2)
}
`
