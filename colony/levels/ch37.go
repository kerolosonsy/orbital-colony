package levels

// Chapter 37 is the capstone: one command-center API grown across five
// stages. Every stage carries the player's files from the previous one.
// Contract tests always run against SQLite; the same suite also runs against
// Postgres when ORBITAL_PG_DSN reaches a live server, and skips otherwise.
// Level.Service stays empty on purpose: the runner's service precheck would
// block the SQLite path whenever Postgres is down.

func init() {
	add(&Level{ID: "37-1-schema", Chapter: 37, Kind: "lab",
		Title: "الأساس", Goal: "ابنِ schema مركز القيادة بترحيلات up/down وrepository خلف interface وسيرفر HTTP صغير", Setup: seededLabRoom,
		Requires: []string{"sel:sql.ErrNoRows", "defer"}, Deps: []string{"modernc.org/sqlite", "github.com/jackc/pgx/v5"}, Bridge: []int{100}, Race: false, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch371Starter}},
		Lesson: `# مركز القيادة: الأساس

## قبل ما تكتب
- لو نفس الكود اشتغل على SQLite عندك وعلى Postgres في الإنتاج، إيه اللي لازم يفضل مشترك بينهم؟
- لو ترحيل وقع في النص، القاعدة لازم تبقى في أنهي نسخة؟

## المشروع الختامي
الفصل ده مفيهوش مفهوم جديد. هتبني API واحد لمركز قيادة المستعمرة على خمس مراحل: الأساس، الدخول، الأطقم، المهام، والتشغيل. كل مرحلة بتبدأ من ملفاتك في المرحلة اللي قبلها (~CarryFrom~)، فاللي تكتبه النهارده هتكمّل عليه بكرة.

المرحلة دي بتجمع ثلاث مهمات قديمة: الجداول والقيود من ~29-1-sqlite-crud~، والترحيلات وجدول النسخ من ~29-2-migrations~، والـ repository خلف interface من ~29-5-sqlx~. الجديد هنا قرارات الدمج بس.

## SQL واحدة لدرايفرين
التستات المخفية بتشغّل نفس العقد على ~modernc.org/sqlite~ دايماً، وعلى Postgres كمان لو ~ORBITAL_PG_DSN~ بيوصل لسيرفر شغال. يعني الـ SQL بتاعك لازم يتقبل من الاتنين.

الـ placeholder المشترك هو ~$1~ و ~$2~ بالترتيب. SQLite بتفهم ~$N~ كرقم المعامل زي Postgres بالظبط، أما ~?~ فـ Postgres مبيفهمهاش.
~~~go
db.QueryRowContext(ctx, "SELECT name FROM crews WHERE id = $1", id)
~~~
الأنواع المشتركة: ~TEXT~ و ~BIGINT~ و ~BOOLEAN~. المفاتيح هنا نصوص بيولّدها المنادي، فمش محتاج ~AUTOINCREMENT~ ولا ~SERIAL~ اللي كل قاعدة بتكتبهم بشكل.

## التكرار من غير ما تقرا رسالة الدرايفر
رسالة خطأ ~UNIQUE~ مختلفة بين الدرايفرين. بدل ما تحلل النص، خلي القاعدة ترفض بهدوء: ~ON CONFLICT DO NOTHING~ بتتجاهل الصف المكرر، و ~RowsAffected~ بصفر معناها «موجود قبل كده».
~~~go
res, err := db.ExecContext(ctx, "INSERT INTO beacons(id) VALUES ($1) ON CONFLICT DO NOTHING", id)
n, _ := res.RowsAffected()
if n == 0 { return ErrConflict }
~~~

## مفاتيح أجنبية
~REFERENCES crews(id)~ بتمنع مهمة لطاقم مش موجود، و ~ON DELETE CASCADE~ بتمسح مهام الطاقم لما الطاقم يتمسح. SQLite مبتفرضهاش إلا لو الاتصال فاتح ~foreign_keys~؛ التستات بتفتحه في الـ DSN بـ ~_pragma=foreign_keys(1)~.
~~~sql
CREATE TABLE logs (
 id TEXT PRIMARY KEY,
 crew_id TEXT NOT NULL REFERENCES crews(id) ON DELETE CASCADE
)
~~~

## ترحيل = transaction
كل ترحيل ومعاه سطر جدول النسخ في transaction واحدة (~BeginTx~ ثم ~defer tx.Rollback()~ ثم ~Commit~). الاتنين بيدعموا DDL جوه transaction، فالترحيل اللي يقع في النص مبيسيبش جدول نصه اتعمل.
~~~go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback() // من غير أثر بعد Commit
~~~

## مسارات فيها متغيرات
~http.ServeMux~ بيقبل نمط فيه طريقة ومتغيرات: ~"GET /v1/crews/{crew}"~. ~r.PathValue("crew")~ بترجع الجزء اللي اتطابق. أي طريقة تانية على نفس المسار بترجع 405 لوحدها.
~~~go
mux.HandleFunc("GET /v1/docks/{dock}", func(w http.ResponseWriter, r *http.Request) {
 fmt.Fprint(w, r.PathValue("dock"))
})
~~~

## أخطاء شائعة
~~~go
db.ExecContext(ctx, "INSERT INTO crews VALUES (?, ?)", id, name) // ? مش شغالة على Postgres
if strings.Contains(err.Error(), "UNIQUE") { return ErrConflict } // نص الدرايفر بيختلف
tx.ExecContext(ctx, m.Up); db.ExecContext(ctx, record) // التسجيل برا الـ transaction
~~~

## المهمة
املا ~Migrations~ بالنسخ ~1, 2, 3~ بالترتيب (ممنوع فجوات): ~users(id, email, password_hash)~ بـ ~email~ فريد مش فاضي، ~crews(id, name)~ بـ ~name~ مش فاضي بعد ~trim~، و ~tasks(id, crew_id, title, done, version)~ بـ ~crew_id~ مرجع لـ ~crews~ بـ ~ON DELETE CASCADE~، و ~title~ مش فاضي بعد ~trim~، و ~done~ افتراضي ~FALSE~، و ~version >= 1~. كل ~Down~ بتشيل اللي عملته ~Up~ بتاعتها بالظبط.

~SchemaVersion~ بتنشئ ~schema_migrations(version INTEGER PRIMARY KEY)~ لو مش موجود وترجع أكبر نسخة أو ~0~. ~Migrate(ctx, db, target)~ تطلع أو تنزل لحد ~target~ خطوة خطوة، كل خطوة في transaction مع سطرها. ~target~ برا ~0..len(Migrations)~، أو قاعدة نسختها أحدث من الكود، أو نسخ مش متتالية: خطأ. خطوة فاشلة ترجع خطأ والنسخة تفضل آخر خطوة نجحت.

~Store~ ينفّذ ~Repository~. ~CreateUser~ تحفظ الإيميل بعد ~TrimSpace~ و ~ToLower~، و ~UserByEmail~ تدوّر بنفس التطبيع. التكرار في أي ~Create~ يرجع ~ErrConflict~، و ~CreateTask~ لطاقم مش موجود ترجع ~ErrNotFound~. ~GetTask(ctx, crewID, id)~ ترجع ~ErrNotFound~ لو المهمة مش في الطاقم ده بالذات. أي خطأ تاني من القاعدة يفضل خطأ.

~NewServer(s)~: ~GET /healthz~ يرجع 200 و ~{"schema":N}~، أو 503 و ~{"error":"unavailable"}~ لو القاعدة وقعت. ~GET /v1/crews/{crew}/tasks/{task}~ يرجع 200 والمهمة JSON بالمفاتيح ~id, crew_id, title, done, version~، أو 404 و ~{"error":"not_found"}~. كل الردود ~application/json~. النجمة التالتة: مسح طاقم بيمسح مهامه. ~main~ فاضية.
`,
		Hints: []string{"ابدأ بـ Migrations وSchemaVersion؛ جرّب Migrate لـ 3 ثم 0 ثم 3 تاني.", "اعمل دالة خطوة واحدة: BeginTx، defer Rollback، Exec للـ SQL، Exec لسطر النسخة، Commit.", "استخدم $1 في كل الاستعلامات وON CONFLICT DO NOTHING مع RowsAffected؛ في CreateTask اتأكد من الطاقم الأول جوه نفس الـ transaction."},
		En: &LevelText{Title: "Foundation", Goal: "Build the command-center schema with up/down migrations, a repository behind an interface, and a small HTTP server", Files: []LabFile{{Name: "main.go", Starter: ch371Starter}},
			Lesson: `# Command Center: Foundation

## Before you write
- If the same code runs on SQLite locally and on Postgres in production, what must stay common to both?
- If a migration fails halfway, which version must the database be left at?

## The capstone
This chapter introduces no new concept. You build one API for the colony command center across five stages: foundation, sign-in, crews, tasks, and operations. Each stage starts from your files in the previous stage (~CarryFrom~), so what you write today you extend tomorrow.

This stage combines three earlier missions: tables and constraints from ~29-1-sqlite-crud~, migrations and a version ledger from ~29-2-migrations~, and a repository behind an interface from ~29-5-sqlx~. Only the integration decisions are new.

## One SQL dialect for two drivers
The hidden tests always run the same contract on ~modernc.org/sqlite~, and also on Postgres when ~ORBITAL_PG_DSN~ reaches a live server. Your SQL must be accepted by both.

The shared placeholder style is ~$1~, ~$2~ in order. SQLite reads ~$N~ as a parameter number just like Postgres; Postgres does not understand ~?~.
~~~go
db.QueryRowContext(ctx, "SELECT name FROM crews WHERE id = $1", id)
~~~
Shared types: ~TEXT~, ~BIGINT~, and ~BOOLEAN~. Keys here are caller-generated text, so you need neither ~AUTOINCREMENT~ nor ~SERIAL~, which each database spells differently.

## Duplicates without parsing driver messages
A ~UNIQUE~ violation reads differently in each driver. Instead of parsing text, let the database decline quietly: ~ON CONFLICT DO NOTHING~ skips the duplicate row, and a zero ~RowsAffected~ means "already exists".
~~~go
res, err := db.ExecContext(ctx, "INSERT INTO beacons(id) VALUES ($1) ON CONFLICT DO NOTHING", id)
n, _ := res.RowsAffected()
if n == 0 { return ErrConflict }
~~~

## Foreign keys
~REFERENCES crews(id)~ blocks a task for a missing crew, and ~ON DELETE CASCADE~ removes a crew's tasks when the crew is deleted. SQLite enforces them only when the connection enables ~foreign_keys~; the tests enable it in the DSN with ~_pragma=foreign_keys(1)~.
~~~sql
CREATE TABLE logs (
 id TEXT PRIMARY KEY,
 crew_id TEXT NOT NULL REFERENCES crews(id) ON DELETE CASCADE
)
~~~

## A migration is a transaction
Each migration and its ledger row share one transaction (~BeginTx~, then ~defer tx.Rollback()~, then ~Commit~). Both databases support DDL inside a transaction, so a migration that fails halfway never leaves a half-built table.
~~~go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback() // harmless after Commit
~~~

## Paths with variables
~http.ServeMux~ accepts a pattern with a method and wildcards: ~"GET /v1/crews/{crew}"~. ~r.PathValue("crew")~ returns the matched segment. Other methods on the same path get 405 automatically.
~~~go
mux.HandleFunc("GET /v1/docks/{dock}", func(w http.ResponseWriter, r *http.Request) {
 fmt.Fprint(w, r.PathValue("dock"))
})
~~~

## Common mistakes
~~~go
db.ExecContext(ctx, "INSERT INTO crews VALUES (?, ?)", id, name) // ? fails on Postgres
if strings.Contains(err.Error(), "UNIQUE") { return ErrConflict } // driver text differs
tx.ExecContext(ctx, m.Up); db.ExecContext(ctx, record) // ledger row outside the transaction
~~~

## Your mission
Fill ~Migrations~ with versions ~1, 2, 3~ in order (no gaps): ~users(id, email, password_hash)~ with a unique nonempty ~email~; ~crews(id, name)~ with a ~name~ nonblank after ~trim~; and ~tasks(id, crew_id, title, done, version)~ where ~crew_id~ references ~crews~ with ~ON DELETE CASCADE~, ~title~ is nonblank after ~trim~, ~done~ defaults to ~FALSE~, and ~version >= 1~. Each ~Down~ removes exactly what its ~Up~ created.

~SchemaVersion~ creates ~schema_migrations(version INTEGER PRIMARY KEY)~ if missing and returns the highest version or ~0~. ~Migrate(ctx, db, target)~ walks up or down to ~target~ one step at a time, each step in a transaction with its ledger row. A ~target~ outside ~0..len(Migrations)~, a database newer than the code, or nonconsecutive versions: error. A failing step returns an error and the version stays at the last successful step.

~Store~ implements ~Repository~. ~CreateUser~ stores the email after ~TrimSpace~ and ~ToLower~, and ~UserByEmail~ looks up with the same normalization. A duplicate in any ~Create~ returns ~ErrConflict~, and ~CreateTask~ for a missing crew returns ~ErrNotFound~. ~GetTask(ctx, crewID, id)~ returns ~ErrNotFound~ unless the task belongs to that exact crew. Any other database error stays an error.

~NewServer(s)~: ~GET /healthz~ returns 200 and ~{"schema":N}~, or 503 and ~{"error":"unavailable"}~ if the database is down. ~GET /v1/crews/{crew}/tasks/{task}~ returns 200 and the task as JSON with keys ~id, crew_id, title, done, version~, or 404 and ~{"error":"not_found"}~. Every reply is ~application/json~. Third star: deleting a crew deletes its tasks. Leave ~main~ empty.
`,
			Hints: []string{"Start with Migrations and SchemaVersion; try Migrate to 3, then 0, then 3 again.", "Write a single-step helper: BeginTx, defer Rollback, Exec the SQL, Exec the ledger row, Commit.", "Use $1 in every query and ON CONFLICT DO NOTHING with RowsAffected; in CreateTask check the crew first inside the same transaction."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch371Hidden},
	})
}

const ch371Starter = `package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
)

type Migration struct {
	Version  int
	Up, Down string
}

var Migrations = []Migration{}

type User struct{ ID, Email, PasswordHash string }
type Crew struct{ ID, Name string }
type Task struct {
	ID      string ` + "`json:\"id\"`" + `
	CrewID  string ` + "`json:\"crew_id\"`" + `
	Title   string ` + "`json:\"title\"`" + `
	Done    bool   ` + "`json:\"done\"`" + `
	Version int64  ` + "`json:\"version\"`" + `
}

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type Repository interface {
	CreateUser(ctx context.Context, u User) error
	UserByEmail(ctx context.Context, email string) (User, error)
	CreateCrew(ctx context.Context, c Crew) error
	CreateTask(ctx context.Context, t Task) error
	GetTask(ctx context.Context, crewID, id string) (Task, error)
}

type Store struct{ db *sql.DB }

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

func SchemaVersion(ctx context.Context, db *sql.DB) (int, error) { return 0, errors.New("TODO") }
func Migrate(ctx context.Context, db *sql.DB, target int) error  { return errors.New("TODO") }

func (s *Store) CreateUser(ctx context.Context, u User) error { return errors.New("TODO") }
func (s *Store) UserByEmail(ctx context.Context, email string) (User, error) {
	return User{}, errors.New("TODO")
}
func (s *Store) CreateCrew(ctx context.Context, c Crew) error { return errors.New("TODO") }
func (s *Store) CreateTask(ctx context.Context, t Task) error { return errors.New("TODO") }
func (s *Store) GetTask(ctx context.Context, crewID, id string) (Task, error) {
	return Task{}, errors.New("TODO")
}

func NewServer(s *Store) http.Handler { return http.NotFoundHandler() }

func main() {}
`

// ch37Backends is shared by every stage's hidden tests: SQLite always,
// Postgres (in a throwaway schema) only when ORBITAL_PG_DSN answers a ping.
const ch37Backends = `
type hiddenBackend struct {
	name string
	open func(t *testing.T) *sql.DB
}

func hiddenBackends(t *testing.T) []hiddenBackend {
	backends := []hiddenBackend{{"sqlite", func(t *testing.T) *sql.DB {
		t.Helper()
		db, err := sql.Open("sqlite", "file:"+filepath.Join(t.TempDir(), "command.db")+"?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { db.Close() })
		return db
	}}}
	return append(backends, hiddenBackend{"postgres", hiddenPostgres})
}

func hiddenPostgres(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("ORBITAL_PG_DSN")
	if dsn == "" {
		t.Skip("ORBITAL_PG_DSN not set")
	}
	admin, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Skip("postgres unavailable: ", err)
	}
	defer admin.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := admin.PingContext(ctx); err != nil {
		t.Skip("postgres unavailable: ", err)
	}
	schema := fmt.Sprintf("orbital37_%d", time.Now().UnixNano())
	if _, err := admin.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	db, err := sql.Open("pgx", u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
		if admin, err := sql.Open("pgx", dsn); err == nil {
			admin.Exec("DROP SCHEMA " + schema + " CASCADE")
			admin.Close()
		}
	})
	return db
}

func eachBackend(t *testing.T, f func(t *testing.T, db *sql.DB)) {
	for _, b := range hiddenBackends(t) {
		t.Run(b.name, func(t *testing.T) { f(t, b.open(t)) })
	}
}

func hiddenGetJSON(t *testing.T, client *http.Client, method, target string, body string, header map[string]string) (int, http.Header, map[string]any) {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, target, reader)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	out := map[string]any{}
	if len(raw) > 0 && json.Unmarshal(raw, &out) != nil {
		t.Fatalf("%s %s: body is not a JSON object: %q", method, target, raw)
	}
	return resp.StatusCode, resp.Header, out
}
`

const ch371Hidden = `package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)
` + ch37Backends + `
var _ Repository = (*Store)(nil)

func hiddenSchemaTable(db *sql.DB, table string) bool {
	_, err := db.Exec("SELECT 1 FROM " + table + " WHERE 1 = 0")
	return err == nil
}

func TestHiddenSchemaRoundTrip(t *testing.T) {
	if len(Migrations) < 3 {
		t.Fatal("need at least migrations 1..3, got", len(Migrations))
	}
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		ctx := context.Background()
		if v, err := SchemaVersion(ctx, db); err != nil || v != 0 {
			t.Fatal("fresh version", v, err)
		}
		for round := 0; round < 2; round++ {
			if err := Migrate(ctx, db, len(Migrations)); err != nil {
				t.Fatal("up", err)
			}
			if err := Migrate(ctx, db, len(Migrations)); err != nil {
				t.Fatal("repeat up must be a no-op", err)
			}
			if v, _ := SchemaVersion(ctx, db); v != len(Migrations) {
				t.Fatal("version after up", v)
			}
			for _, table := range []string{"users", "crews", "tasks"} {
				if !hiddenSchemaTable(db, table) {
					t.Fatal("missing table", table)
				}
			}
			if err := Migrate(ctx, db, 0); err != nil {
				t.Fatal("down", err)
			}
			if v, _ := SchemaVersion(ctx, db); v != 0 {
				t.Fatal("version after down", v)
			}
			for _, table := range []string{"users", "crews", "tasks"} {
				if hiddenSchemaTable(db, table) {
					t.Fatal("down left table", table)
				}
			}
		}
		if err := Migrate(ctx, db, 1); err != nil {
			t.Fatal("partial up", err)
		}
		if v, _ := SchemaVersion(ctx, db); v != 1 || !hiddenSchemaTable(db, "users") || hiddenSchemaTable(db, "tasks") {
			t.Fatal("partial target", v)
		}
		for _, bad := range []int{-1, len(Migrations) + 1} {
			if err := Migrate(ctx, db, bad); err == nil {
				t.Fatal("accepted target", bad)
			}
		}
		if _, err := db.Exec("INSERT INTO schema_migrations(version) VALUES ($1)", len(Migrations)+5); err != nil {
			t.Fatal(err)
		}
		if err := Migrate(ctx, db, len(Migrations)); err == nil {
			t.Fatal("database newer than code was accepted")
		}
	})
}

func TestHiddenMigrationFailureKeepsVersion(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		ctx := context.Background()
		saved := Migrations
		defer func() { Migrations = saved }()
		n := len(saved)
		Migrations = append(append([]Migration{}, saved...), Migration{n + 1, "CREATE TABLE probe_half (id TEXT); INSERT INTO probe_missing VALUES (1)", "DROP TABLE probe_half"})
		if err := Migrate(ctx, db, n+1); err == nil {
			t.Fatal("broken migration succeeded")
		}
		if v, _ := SchemaVersion(ctx, db); v != n {
			t.Fatal("version after failed step", v, "want", n)
		}
		if hiddenSchemaTable(db, "probe_half") {
			t.Fatal("failed migration left a table behind")
		}
		Migrations = append(append([]Migration{}, saved...), Migration{n + 2, "SELECT 1", "SELECT 1"})
		if err := Migrate(ctx, db, n+1); err == nil {
			t.Fatal("gap in versions accepted")
		}
		ctx2, cancel := context.WithCancel(ctx)
		cancel()
		Migrations = saved
		if err := Migrate(ctx2, db, 0); err == nil {
			t.Fatal("canceled context ignored")
		}
	})
}

func hiddenSeededStore(t *testing.T, db *sql.DB) *Store {
	t.Helper()
	if err := Migrate(context.Background(), db, len(Migrations)); err != nil {
		t.Fatal(err)
	}
	return NewStore(db)
}

func TestHiddenRepository(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		ctx := context.Background()
		s := hiddenSeededStore(t, db)
		rng := rand.New(rand.NewSource(3701))
		email := fmt.Sprintf("  Pilot%d@Colony.Example ", rng.Intn(1000))
		if err := s.CreateUser(ctx, User{"u1", email, "digest-1"}); err != nil {
			t.Fatal(err)
		}
		u, err := s.UserByEmail(ctx, strings.ToUpper(email))
		if err != nil || u.ID != "u1" || u.Email != strings.ToLower(strings.TrimSpace(email)) || u.PasswordHash != "digest-1" {
			t.Fatal("user lookup", u, err)
		}
		if err := s.CreateUser(ctx, User{"u2", strings.ToLower(email), "x"}); !errors.Is(err, ErrConflict) {
			t.Fatal("duplicate email", err)
		}
		if err := s.CreateUser(ctx, User{"u1", "other@colony", "x"}); !errors.Is(err, ErrConflict) {
			t.Fatal("duplicate id", err)
		}
		if _, err := s.UserByEmail(ctx, "ghost@colony"); !errors.Is(err, ErrNotFound) {
			t.Fatal("missing user", err)
		}
		crews := []string{}
		for i := 0; i < 3; i++ {
			id := fmt.Sprintf("crew-%d-%d", i, rng.Intn(1e6))
			crews = append(crews, id)
			if err := s.CreateCrew(ctx, Crew{id, fmt.Sprintf("طاقم %d", i)}); err != nil {
				t.Fatal(err)
			}
		}
		if err := s.CreateCrew(ctx, Crew{crews[0], "again"}); !errors.Is(err, ErrConflict) {
			t.Fatal("duplicate crew", err)
		}
		if err := s.CreateCrew(ctx, Crew{"blank", "   "}); err == nil || errors.Is(err, ErrConflict) {
			t.Fatal("blank crew name", err)
		}
		want := map[string]Task{}
		for i := 0; i < 20; i++ {
			task := Task{fmt.Sprintf("t%d", i), crews[rng.Intn(len(crews))], fmt.Sprintf("فحص '%d; --", rng.Int63()), rng.Intn(2) == 0, int64(1 + rng.Intn(5))}
			want[task.ID] = task
			if err := s.CreateTask(ctx, task); err != nil {
				t.Fatal(err)
			}
		}
		for id, task := range want {
			got, err := s.GetTask(ctx, task.CrewID, id)
			if err != nil || got != task {
				t.Fatal("get task", got, task, err)
			}
			for _, crew := range crews {
				if crew != task.CrewID {
					if _, err := s.GetTask(ctx, crew, id); !errors.Is(err, ErrNotFound) {
						t.Fatal("task visible from another crew", err)
					}
				}
			}
		}
		if err := s.CreateTask(ctx, Task{"t0", crews[0], "dup", false, 1}); !errors.Is(err, ErrConflict) {
			t.Fatal("duplicate task", err)
		}
		if err := s.CreateTask(ctx, Task{"orphan", "no-such-crew", "x", false, 1}); !errors.Is(err, ErrNotFound) {
			t.Fatal("task for missing crew", err)
		}
		for _, bad := range []Task{{"b1", crews[0], "  ", false, 1}, {"b2", crews[0], "ok", false, 0}} {
			if err := s.CreateTask(ctx, bad); err == nil || errors.Is(err, ErrConflict) {
				t.Fatal("constraint accepted", bad, err)
			}
		}
		if _, err := db.Exec("INSERT INTO tasks(id, crew_id, title, version) VALUES ('raw', 'no-such-crew', 'x', 1)"); err == nil {
			t.Fatal("foreign key not enforced by the schema")
		}
		var done bool
		if _, err := db.Exec("INSERT INTO tasks(id, crew_id, title, version) VALUES ('dflt', $1, 'x', 1)", crews[0]); err != nil {
			t.Fatal(err)
		}
		if err := db.QueryRow("SELECT done FROM tasks WHERE id = 'dflt'").Scan(&done); err != nil || done {
			t.Fatal("done default", done, err)
		}
		if db.Stats().InUse != 0 {
			t.Fatal("connection leaked")
		}
		db.Close()
		if _, err := s.GetTask(ctx, crews[0], "t0"); err == nil || errors.Is(err, ErrNotFound) {
			t.Fatal("closed database must be an error, not not-found", err)
		}
	})
}

func TestHiddenHTTPContract(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		ctx := context.Background()
		s := hiddenSeededStore(t, db)
		rng := rand.New(rand.NewSource(3702))
		crew := fmt.Sprintf("c%d", rng.Intn(1e6))
		task := Task{fmt.Sprintf("k%d", rng.Intn(1e6)), crew, "سد التسريب", true, int64(1 + rng.Intn(9))}
		if err := s.CreateCrew(ctx, Crew{crew, "Vega"}); err != nil {
			t.Fatal(err)
		}
		if err := s.CreateCrew(ctx, Crew{"rival", "Rival"}); err != nil {
			t.Fatal(err)
		}
		if err := s.CreateTask(ctx, task); err != nil {
			t.Fatal(err)
		}
		srv := httptest.NewServer(NewServer(s))
		defer srv.Close()
		c := srv.Client()
		code, h, body := hiddenGetJSON(t, c, "GET", srv.URL+"/healthz", "", nil)
		if code != 200 || body["schema"] != float64(len(Migrations)) || !strings.HasPrefix(h.Get("Content-Type"), "application/json") {
			t.Fatal("healthz", code, body, h)
		}
		code, h, body = hiddenGetJSON(t, c, "GET", srv.URL+"/v1/crews/"+crew+"/tasks/"+task.ID, "", nil)
		if code != 200 || body["id"] != task.ID || body["crew_id"] != crew || body["title"] != task.Title || body["done"] != true || body["version"] != float64(task.Version) {
			t.Fatal("get task", code, body)
		}
		for _, path := range []string{"/v1/crews/rival/tasks/" + task.ID, "/v1/crews/" + crew + "/tasks/missing"} {
			code, h, body = hiddenGetJSON(t, c, "GET", srv.URL+path, "", nil)
			if code != 404 || body["error"] != "not_found" || !strings.HasPrefix(h.Get("Content-Type"), "application/json") {
				t.Fatal("not found shape", path, code, body)
			}
		}
		db.Close()
		code, _, body = hiddenGetJSON(t, c, "GET", srv.URL+"/healthz", "", nil)
		if code != 503 || body["error"] != "unavailable" {
			t.Fatal("healthz with database down", code, body)
		}
	})
}

func TestHiddenBonusCascade(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		ctx := context.Background()
		s := hiddenSeededStore(t, db)
		if err := s.CreateCrew(ctx, Crew{"gone", "Gone"}); err != nil {
			t.Fatal(err)
		}
		if err := s.CreateTask(ctx, Task{"g1", "gone", "x", false, 1}); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec("DELETE FROM crews WHERE id = 'gone'"); err != nil {
			t.Fatal(err)
		}
		if _, err := s.GetTask(ctx, "gone", "g1"); !errors.Is(err, ErrNotFound) {
			t.Fatal("tasks survived their crew", err)
		}
	})
}
`

func init() {
	add(&Level{ID: "37-2-auth", Chapter: 37, Kind: "lab", CarryFrom: "37-1-schema",
		Title: "بوابة الدخول", Goal: "ضيف تسجيل ودخول بجلسات معتمة متخزّن منها digest بس، وخروج، وإلغاء الجلسات عند تغيير الباسورد", Setup: seededLabRoom,
		Requires: []string{"sel:bcrypt.GenerateFromPassword", "sel:sha256.Sum256"}, Deps: []string{"modernc.org/sqlite", "github.com/jackc/pgx/v5", "golang.org/x/crypto"}, Bridge: []int{100}, Race: false, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch371Starter}, {Name: "auth.go", Starter: ch372Starter}},
		Lesson: `# مركز القيادة: بوابة الدخول

## قبل ما تكتب
- لو حد سرق نسخة من القاعدة، يقدر يدخل بالجلسات اللي فيها؟
- لما الباسورد يتغيّر، مين لسه المفروض يفضل داخل؟

## بتبني على إيه
~main.go~ جاي من المرحلة اللي فاتت كما هو، والملف الجديد ~auth.go~. هتستخدم ~bcrypt~ من ~27-1-signed-pass~، والجلسات والكوكيز من ~27-2-sessions~، و ~crypto/rand~ و ~sha256~ من فصل الخزنة (23). الجديد قرارات الدمج بس.

## ترحيل جديد، مش تعديل قديم
القاعدة عند اللاعبين التانيين واقفة على نسخة 3. ممنوع تعدّل ترحيل اتطبّق؛ ضيف نسخة 4 في آخر ~Migrations~.
~~~go
{4, "CREATE TABLE badges (id TEXT PRIMARY KEY)", "DROP TABLE badges"},
~~~

## token معتم وdigest بس
الـ token اللي في الكوكي 32 بايت عشوائية من ~Auth.Random~ (التستات بتحقن مصدر ثابت). القاعدة بتخزن ~sha256~ بتاعه بس؛ لو القاعدة اتسرقت، الـ digest مينفعش يتحط في كوكي.
~~~go
sum := sha256.Sum256([]byte(token))
key := hex.EncodeToString(sum[:])
~~~
الدخول بيدوّر بالـ digest، والخروج بيمسح الصف. البحث بالـ digest مبيكشفش بايتات الـ token، فمش محتاج مقارنة constant-time هنا.

## الوقت والعمر محقونين
~Auth.Now~ ساعة محقونة و ~Auth.TTL~ عمر الجلسة. ~expires_at~ بالثواني (~Unix()~)، والجلسة منتهية لو ~Now >= expires_at~. الكوكي بياخد نفس العمر في ~MaxAge~.

## إلغاء عند تغيير الباسورد
تحديث الهاش ومسح كل جلسات المستخدم في transaction واحدة، فمفيش لحظة الهاش الجديد فيها شغال والجلسات القديمة لسه عايشة.

## أخطاء شائعة
~~~go
db.ExecContext(ctx, "INSERT INTO sessions VALUES ($1, $2, $3)", token, id, exp) // الـ token نفسه في القاعدة
if errors.Is(err, ErrNotFound) { return "", ErrNotFound } // بيكشف إن الإيميل مش متسجّل
http.SetCookie(w, &http.Cookie{Name: "session", Value: token}) // من غير HttpOnly وSecure
~~~

## المهمة
ضيف ترحيل 4: ~sessions(digest TEXT PRIMARY KEY, user_id ..., expires_at BIGINT)~ و ~user_id~ مرجع لـ ~users~. في ~auth.go~:

~Register(ctx, email, password)~: باسورد أقل من 8 أو أكتر من 72 بايت يرجع ~ErrWeakPassword~؛ غير كده هاش ~bcrypt~ بتكلفة ~a.Cost~ و ID عشوائي من ~a.Random~، والإيميل المكرر ~ErrConflict~. ~Login~ ترجع token جديد بـ ~base64.RawURLEncoding~ لـ 32 بايت؛ إيميل مش موجود أو باسورد غلط: نفس ~ErrBadCredentials~. ~Authenticate(ctx, token)~ ترجع المستخدم أو ~ErrUnauthorized~ للمجهول والمنتهي. ~Logout~ تمسح الجلسة. ~ChangePassword(ctx, userID, old, new)~ تتأكد من القديم (~ErrBadCredentials~)، وتطبّق نفس حد القوة، وتلغي كل جلسات المستخدم.

~NewAPI(a)~ بيخدم كل مسارات ~NewServer~ كمان، وبيضيف: ~POST /v1/register~ بـ ~{"email","password"}~ يرجع 201 و ~{"id","email"}~، أو 409 ~conflict~، أو 400 ~weak_password~. ~POST /v1/login~ يرجع 204 وكوكي ~session~ بـ ~HttpOnly~ و ~Secure~ و ~SameSite=Lax~ و ~Path=/~ و ~MaxAge~ بثواني ~TTL~. ~GET /v1/me~ بالكوكي يرجع 200 و ~{"id","email"}~. ~POST /v1/logout~ يرجع 204 ويمسح الكوكي (~MaxAge~ سالب). ~POST /v1/password~ بـ ~{"old","new"}~ يرجع 204. أي فشل دخول: 401 و ~{"error":"unauthorized"}~. JSON مكسور: 400 و ~{"error":"bad_request"}~. النجمة التالتة: الجلسة المنتهية تتمسح من القاعدة أول ما حد يحاول يستخدمها.
`,
		Hints: []string{"ابدأ بترحيل 4 ودالة digest(token) string بـ sha256 وhex.", "Login: UserByEmail ثم bcrypt.CompareHashAndPassword؛ الغياب والغلط نفس ErrBadCredentials.", "ChangePassword: UPDATE للهاش وDELETE لجلسات المستخدم في نفس الـ transaction؛ وفي Authenticate امسح الصف المنتهي للنجمة."},
		En: &LevelText{Title: "Sign-in Gate", Goal: "Add registration and sign-in with opaque sessions stored only as digests, logout, and session revocation on password change", Files: []LabFile{{Name: "main.go", Starter: ch371Starter}, {Name: "auth.go", Starter: ch372Starter}},
			Lesson: `# Command Center: Sign-in Gate

## Before you write
- If someone steals a copy of the database, can they sign in with the sessions inside it?
- When a password changes, who should still be signed in?

## What this builds on
~main.go~ carries over from the previous stage as is, and the new file is ~auth.go~. You use ~bcrypt~ from ~27-1-signed-pass~, sessions and cookies from ~27-2-sessions~, and ~crypto/rand~ and ~sha256~ from the Vault chapter (23). Only the integration decisions are new.

## A new migration, not an edited one
Other installations already sit at version 3. Never edit an applied migration; append version 4 to ~Migrations~.
~~~go
{4, "CREATE TABLE badges (id TEXT PRIMARY KEY)", "DROP TABLE badges"},
~~~

## Opaque token, digest only
The cookie token is 32 random bytes from ~Auth.Random~ (tests inject a fixed source). The database stores only its ~sha256~; if the database leaks, a digest cannot be replayed as a cookie.
~~~go
sum := sha256.Sum256([]byte(token))
key := hex.EncodeToString(sum[:])
~~~
Sign-in looks up by digest and logout deletes the row. Looking up by digest reveals no token bytes, so no constant-time compare is needed here.

## Injected time and lifetime
~Auth.Now~ is an injected clock and ~Auth.TTL~ the session lifetime. ~expires_at~ is in seconds (~Unix()~), and a session is expired when ~Now >= expires_at~. The cookie gets the same lifetime in ~MaxAge~.

## Revocation on password change
Updating the hash and deleting all of the user's sessions share one transaction, so there is no moment where the new hash works while old sessions are still alive.

## Common mistakes
~~~go
db.ExecContext(ctx, "INSERT INTO sessions VALUES ($1, $2, $3)", token, id, exp) // raw token in the database
if errors.Is(err, ErrNotFound) { return "", ErrNotFound } // reveals the email is not registered
http.SetCookie(w, &http.Cookie{Name: "session", Value: token}) // no HttpOnly or Secure
~~~

## Your mission
Add migration 4: ~sessions(digest TEXT PRIMARY KEY, user_id ..., expires_at BIGINT)~ with ~user_id~ referencing ~users~. In ~auth.go~:

~Register(ctx, email, password)~: a password under 8 or over 72 bytes returns ~ErrWeakPassword~; otherwise a ~bcrypt~ hash at cost ~a.Cost~ and a random ID from ~a.Random~, and a duplicate email returns ~ErrConflict~. ~Login~ returns a fresh token, ~base64.RawURLEncoding~ of 32 bytes; an unknown email or a wrong password both return ~ErrBadCredentials~. ~Authenticate(ctx, token)~ returns the user, or ~ErrUnauthorized~ for unknown and expired tokens. ~Logout~ deletes the session. ~ChangePassword(ctx, userID, old, new)~ verifies the old password (~ErrBadCredentials~), applies the same strength rule, and revokes all of the user's sessions.

~NewAPI(a)~ serves every ~NewServer~ route too, and adds: ~POST /v1/register~ with ~{"email","password"}~ returns 201 and ~{"id","email"}~, or 409 ~conflict~, or 400 ~weak_password~. ~POST /v1/login~ returns 204 and a ~session~ cookie with ~HttpOnly~, ~Secure~, ~SameSite=Lax~, ~Path=/~, and ~MaxAge~ of ~TTL~ in seconds. ~GET /v1/me~ with the cookie returns 200 and ~{"id","email"}~. ~POST /v1/logout~ returns 204 and clears the cookie (negative ~MaxAge~). ~POST /v1/password~ with ~{"old","new"}~ returns 204. Every sign-in failure: 401 and ~{"error":"unauthorized"}~. Malformed JSON: 400 and ~{"error":"bad_request"}~. Third star: an expired session is deleted from the database as soon as someone tries to use it.
`,
			Hints: []string{"Start with migration 4 and a digest(token) string helper using sha256 and hex.", "Login: UserByEmail, then bcrypt.CompareHashAndPassword; missing and wrong both return ErrBadCredentials.", "ChangePassword: UPDATE the hash and DELETE the user's sessions in one transaction; for the star, delete the expired row in Authenticate."}},
		Lab: &Lab{ReqFile: "auth.go", Bonus: true, Hidden: ch372Hidden},
	})
}

const ch372Starter = `package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	ErrWeakPassword   = errors.New("weak password")
	ErrBadCredentials = errors.New("bad credentials")
	ErrUnauthorized   = errors.New("unauthorized")
)

const SessionCookie = "session"

type Auth struct {
	Store  *Store
	Now    func() time.Time
	Random io.Reader
	TTL    time.Duration
	Cost   int
}

func (a *Auth) Register(ctx context.Context, email, password string) (User, error) {
	return User{}, errors.New("TODO")
}
func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	return "", errors.New("TODO")
}
func (a *Auth) Authenticate(ctx context.Context, token string) (User, error) {
	return User{}, errors.New("TODO")
}
func (a *Auth) Logout(ctx context.Context, token string) error { return errors.New("TODO") }
func (a *Auth) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	return errors.New("TODO")
}

func NewAPI(a *Auth) http.Handler { return NewServer(a.Store) }
`

const ch372Hidden = `package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)
` + ch37Backends + `
type hiddenClock struct{ now time.Time }

func (c *hiddenClock) Now() time.Time { return c.now }

func hiddenAuth(t *testing.T, db *sql.DB, seed int64) (*Auth, *hiddenClock) {
	t.Helper()
	if len(Migrations) < 4 {
		t.Fatal("add migration 4 for sessions")
	}
	if err := Migrate(context.Background(), db, len(Migrations)); err != nil {
		t.Fatal(err)
	}
	clock := &hiddenClock{time.Unix(1_900_000_000, 0)}
	return &Auth{Store: NewStore(db), Now: clock.Now, Random: rand.New(rand.NewSource(seed)), TTL: 2 * time.Hour, Cost: bcrypt.MinCost}, clock
}

func hiddenSessionRows(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.Query("SELECT digest, user_id FROM sessions ORDER BY digest")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var d, u string
		if err := rows.Scan(&d, &u); err != nil {
			t.Fatal(err)
		}
		out = append(out, d+"|"+u)
	}
	return out
}

func TestHiddenAuthFlow(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		ctx := context.Background()
		a, clock := hiddenAuth(t, db, 3721)
		for _, weak := range []string{"", "short7!", strings.Repeat("x", 73)} {
			if _, err := a.Register(ctx, "w@colony", weak); !errors.Is(err, ErrWeakPassword) {
				t.Fatal("weak password", len(weak), err)
			}
		}
		u, err := a.Register(ctx, " Nova@Colony ", "orbit-pass-1")
		if err != nil || u.ID == "" || u.Email != "nova@colony" {
			t.Fatal("register", u, err)
		}
		if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("orbit-pass-1")) != nil {
			t.Fatal("password hash is not bcrypt of the password")
		}
		if _, err := a.Register(ctx, "NOVA@colony", "orbit-pass-2"); !errors.Is(err, ErrConflict) {
			t.Fatal("duplicate email", err)
		}
		other, err := a.Register(ctx, "rigel@colony", "rigel-pass-1")
		if err != nil || other.ID == u.ID {
			t.Fatal("second user", other, err)
		}
		for _, bad := range [][2]string{{"nova@colony", "wrong-pass"}, {"ghost@colony", "orbit-pass-1"}} {
			if _, err := a.Login(ctx, bad[0], bad[1]); !errors.Is(err, ErrBadCredentials) {
				t.Fatal("bad login", bad, err)
			}
		}
		tokens := map[string]bool{}
		for i := 0; i < 3; i++ {
			tok, err := a.Login(ctx, "NOVA@colony", "orbit-pass-1")
			if err != nil || len(tok) != 43 || tokens[tok] {
				t.Fatal("login token", tok, err)
			}
			tokens[tok] = true
		}
		rows := hiddenSessionRows(t, db)
		if len(rows) != 3 {
			t.Fatal("session rows", rows)
		}
		for tok := range tokens {
			sum := sha256.Sum256([]byte(tok))
			found := false
			for _, row := range rows {
				if strings.Contains(row, tok) {
					t.Fatal("raw token stored in the database")
				}
				found = found || strings.HasPrefix(row, hex.EncodeToString(sum[:])+"|"+u.ID)
			}
			if !found {
				t.Fatal("session digest is not hex sha256 of the token")
			}
			if got, err := a.Authenticate(ctx, tok); err != nil || got.ID != u.ID {
				t.Fatal("authenticate", got, err)
			}
		}
		if _, err := a.Authenticate(ctx, "forged"); !errors.Is(err, ErrUnauthorized) {
			t.Fatal("forged token", err)
		}
		var first string
		for tok := range tokens {
			first = tok
			break
		}
		if err := a.Logout(ctx, first); err != nil {
			t.Fatal(err)
		}
		if _, err := a.Authenticate(ctx, first); !errors.Is(err, ErrUnauthorized) {
			t.Fatal("logout did not revoke", err)
		}
		rigel, _ := a.Login(ctx, "rigel@colony", "rigel-pass-1")
		clock.now = clock.now.Add(2*time.Hour - time.Second)
		if _, err := a.Authenticate(ctx, rigel); err != nil {
			t.Fatal("expired too early", err)
		}
		clock.now = clock.now.Add(time.Second)
		if _, err := a.Authenticate(ctx, rigel); !errors.Is(err, ErrUnauthorized) {
			t.Fatal("expired session accepted", err)
		}
		clock.now = time.Unix(1_900_000_000, 0)
		if err := a.ChangePassword(ctx, u.ID, "wrong-pass", "new-orbit-pass"); !errors.Is(err, ErrBadCredentials) {
			t.Fatal("change with wrong old", err)
		}
		if err := a.ChangePassword(ctx, u.ID, "orbit-pass-1", "short"); !errors.Is(err, ErrWeakPassword) {
			t.Fatal("change to weak", err)
		}
		keep, _ := a.Login(ctx, "rigel@colony", "rigel-pass-1")
		if err := a.ChangePassword(ctx, u.ID, "orbit-pass-1", "new-orbit-pass"); err != nil {
			t.Fatal(err)
		}
		for tok := range tokens {
			if _, err := a.Authenticate(ctx, tok); !errors.Is(err, ErrUnauthorized) {
				t.Fatal("password change kept an old session", err)
			}
		}
		if _, err := a.Authenticate(ctx, keep); err != nil {
			t.Fatal("password change revoked another user's session", err)
		}
		if _, err := a.Login(ctx, "nova@colony", "orbit-pass-1"); !errors.Is(err, ErrBadCredentials) {
			t.Fatal("old password still works", err)
		}
		if _, err := a.Login(ctx, "nova@colony", "new-orbit-pass"); err != nil {
			t.Fatal("new password", err)
		}
	})
}

func TestHiddenAuthHTTP(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		a, clock := hiddenAuth(t, db, 3722)
		srv := httptest.NewServer(NewAPI(a))
		defer srv.Close()
		jar, _ := cookiejar.New(nil)
		c := srv.Client()
		c.Jar = jar
		code, _, body := hiddenGetJSON(t, c, "POST", srv.URL+"/v1/register", "{\"email\":\"Vega@Colony\",\"password\":\"vega-pass-1\"}", nil)
		if code != 201 || body["email"] != "vega@colony" || body["id"] == "" || body["password_hash"] != nil {
			t.Fatal("register", code, body)
		}
		id := body["id"]
		for _, c2 := range []struct {
			body string
			code int
			err  string
		}{{"{\"email\":\"vega@colony\",\"password\":\"vega-pass-2\"}", 409, "conflict"}, {"{\"email\":\"x@colony\",\"password\":\"1234567\"}", 400, "weak_password"}, {"{nope", 400, "bad_request"}} {
			code, _, body = hiddenGetJSON(t, c, "POST", srv.URL+"/v1/register", c2.body, nil)
			if code != c2.code || body["error"] != c2.err {
				t.Fatal("register failure", c2, code, body)
			}
		}
		code, _, body = hiddenGetJSON(t, c, "GET", srv.URL+"/v1/me", "", nil)
		if code != 401 || body["error"] != "unauthorized" {
			t.Fatal("me without cookie", code, body)
		}
		for _, bad := range []string{"{\"email\":\"vega@colony\",\"password\":\"wrong-pass\"}", "{\"email\":\"ghost@colony\",\"password\":\"vega-pass-1\"}"} {
			code, _, body = hiddenGetJSON(t, c, "POST", srv.URL+"/v1/login", bad, nil)
			if code != 401 || body["error"] != "unauthorized" {
				t.Fatal("bad login", code, body)
			}
		}
		req, _ := http.NewRequest("POST", srv.URL+"/v1/login", strings.NewReader("{\"email\":\"vega@colony\",\"password\":\"vega-pass-1\"}"))
		resp, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		var session *http.Cookie
		for _, ck := range resp.Cookies() {
			if ck.Name == SessionCookie {
				session = ck
			}
		}
		if resp.StatusCode != 204 || session == nil || !session.HttpOnly || !session.Secure || session.SameSite != http.SameSiteLaxMode || session.Path != "/" || session.MaxAge != 7200 {
			t.Fatalf("login cookie %d %+v", resp.StatusCode, session)
		}
		u, _ := url.Parse(srv.URL)
		jar.SetCookies(u, []*http.Cookie{{Name: SessionCookie, Value: session.Value}})
		code, _, body = hiddenGetJSON(t, c, "GET", srv.URL+"/v1/me", "", nil)
		if code != 200 || body["id"] != id || body["email"] != "vega@colony" {
			t.Fatal("me", code, body)
		}
		code, _, _ = hiddenGetJSON(t, c, "GET", srv.URL+"/healthz", "", nil)
		if code != 200 {
			t.Fatal("NewAPI must still serve stage-1 routes", code)
		}
		code, _, _ = hiddenGetJSON(t, c, "POST", srv.URL+"/v1/password", "{\"old\":\"vega-pass-1\",\"new\":\"vega-pass-2\"}", nil)
		if code != 204 {
			t.Fatal("password change", code)
		}
		jar.SetCookies(u, []*http.Cookie{{Name: SessionCookie, Value: session.Value}})
		code, _, _ = hiddenGetJSON(t, c, "GET", srv.URL+"/v1/me", "", nil)
		if code != 401 {
			t.Fatal("session survived password change", code)
		}
		token, err := a.Login(context.Background(), "vega@colony", "vega-pass-2")
		if err != nil {
			t.Fatal(err)
		}
		jar.SetCookies(u, []*http.Cookie{{Name: SessionCookie, Value: token}})
		req, _ = http.NewRequest("POST", srv.URL+"/v1/logout", nil)
		resp, err = c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		cleared := false
		for _, ck := range resp.Cookies() {
			cleared = cleared || (ck.Name == SessionCookie && ck.MaxAge < 0)
		}
		if resp.StatusCode != 204 || !cleared {
			t.Fatal("logout", resp.StatusCode, resp.Cookies())
		}
		if _, err := a.Authenticate(context.Background(), token); !errors.Is(err, ErrUnauthorized) {
			t.Fatal("logout left the session usable", err)
		}
		late, _ := a.Login(context.Background(), "vega@colony", "vega-pass-2")
		clock.now = clock.now.Add(3 * time.Hour)
		jar.SetCookies(u, []*http.Cookie{{Name: SessionCookie, Value: late}})
		code, _, body = hiddenGetJSON(t, c, "GET", srv.URL+"/v1/me", "", nil)
		if code != 401 || body["error"] != "unauthorized" {
			t.Fatal("expired cookie", code, body)
		}
		_ = fmt.Sprint
	})
}

func TestHiddenBonusExpiredSessionDeleted(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		ctx := context.Background()
		a, clock := hiddenAuth(t, db, 3723)
		if _, err := a.Register(ctx, "deneb@colony", "deneb-pass-1"); err != nil {
			t.Fatal(err)
		}
		tok, err := a.Login(ctx, "deneb@colony", "deneb-pass-1")
		if err != nil {
			t.Fatal(err)
		}
		clock.now = clock.now.Add(2 * time.Hour)
		if _, err := a.Authenticate(ctx, tok); !errors.Is(err, ErrUnauthorized) {
			t.Fatal(err)
		}
		if rows := hiddenSessionRows(t, db); len(rows) != 0 {
			t.Fatal("expired session row kept", rows)
		}
	})
}
`
