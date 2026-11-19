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
		LessonRevision: 1, Steps: ch37FoundationArabicSteps,
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
		Hints: []string{"اعتبر جدول النسخ وصف لحالة schema اللي اتثبتت.", "اختبر كل migration بالاتجاهين، وخلي repository methods منفصلة عن HTTP.", "ارجع لتواقيع Migrate وNewServer في starter وحدد حدود كل transaction."},
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
			LessonRevision: 1, Steps: ch37FoundationEnglishSteps,
			Hints: []string{"Treat the version table as a description of the schema state that actually committed.", "Exercise each migration in both directions, and keep repository methods separate from HTTP.", "Use the Migrate and NewServer signatures in the starter to map transaction boundaries."}},
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
		LessonRevision: 1, Steps: ch37AuthArabicSteps,
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
		Hints: []string{"الجلسة تحمل سرًا للعميل؛ قاعدة البيانات تحتاج قيمة مقارنة من غير حفظ السر الخام.", "افصل تحقق كلمة المرور عن إصدار الجلسة ثم عن إلغائها.", "ارجع لتواقيع auth.go وحدد كل دالة بتنشئ أو تستهلك session."},
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
			LessonRevision: 1, Steps: ch37AuthEnglishSteps,
			Hints: []string{"A session gives the client a secret; the database needs a comparable value without the raw secret.", "Separate password verification from session creation and revocation.", "Use the auth.go signatures to find each function that creates or consumes a session."}},
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

func init() {
	add(&Level{ID: "37-3-tenants", Chapter: 37, Kind: "lab", CarryFrom: "37-2-auth",
		Title: "الأطقم المعزولة", Goal: "ضيف أطقم بأدوار owner/admin/member، وعزل كامل بين الأطقم (الغريب يشوف not found)، ومنع التصعيد، وسجل audit", Setup: seededLabRoom,
		LessonRevision: 1, Steps: ch37TenantsArabicSteps,
		Requires: []string{"sel:sql.ErrNoRows", "map"}, Deps: []string{"modernc.org/sqlite", "github.com/jackc/pgx/v5", "golang.org/x/crypto"}, Bridge: []int{100}, Race: false, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch371Starter}, {Name: "auth.go", Starter: ch372Starter}, {Name: "crews.go", Starter: ch373Starter}},
		Lesson: `# مركز القيادة: الأطقم المعزولة

## قبل ما تكتب
- لو حد برا الطاقم سأل عن مهمة جوه الطاقم، المفروض يعرف إنها موجودة؟
- مين يقدر يدّي صلاحية، وأعلى صلاحية يقدر يديها قد إيه؟

## بتبني على إيه
~main.go~ و ~auth.go~ جايين من المراحل اللي فاتت، والملف الجديد ~crews.go~. الدخول من ~37-2-auth~، والـ middleware اللي بيلف handler من فصل بوابة القاعدة (~26-3-router~)، والـ transactions من ~29-3-transfer~. الجديد قرارات الدمج بس.

## العزل = not found
الطاقم هو الـ tenant. أي طلب على طاقم إنت مش عضو فيه بيرجع 404 بالظبط زي الطاقم اللي مش موجود. لو رجّعت 403، الغريب يعرف إن الـ ID ده حقيقي ويقدر يعدّ الأطقم.
~~~go
role, err := c.Role(ctx, crewID, user.ID)
if errors.Is(err, ErrNotFound) { return ErrNotFound } // غريب = مش موجود
~~~
الـ 403 للعضو بس: هو عارف إن الطاقم موجود، بس دوره مش كفاية.

## الأدوار كرُتب
خزّن الدور نص في القاعدة بـ ~CHECK (role IN (...))~، وقارن في Go برتبة رقمية من map. قاعدة منع التصعيد: محدش يدّي دور قد دوره أو أعلى منه.
~~~go
var level = map[string]int{"cadet": 1, "pilot": 2}
if level[grant] >= level[mine] { return ErrForbidden }
~~~

## audit في نفس الـ transaction
سطر الـ audit بيتكتب في نفس transaction التغيير. لو التغيير اتلغى، السطر يتلغي معاه؛ ولو السطر فشل، التغيير ميتمّش. الترتيب برقم ~seq~ لكل طاقم (~MAX(seq) + 1~ جوه الـ transaction) مش بالوقت، لأن طلبين ممكن ياخدوا نفس الثانية.
~~~go
tx.QueryRowContext(ctx, "SELECT COALESCE(MAX(seq), 0) + 1 FROM audit WHERE crew_id = $1", id).Scan(&seq)
~~~

## أخطاء شائعة
~~~go
if role == "" { http.Error(w, "forbidden", 403) } // بيكشف إن الطاقم موجود
if mine == "admin" || mine == "owner" { grant(role) } // admin يقدر يعمل owner جديد
db.ExecContext(ctx, insertMember); db.ExecContext(ctx, insertAudit) // audit برا الـ transaction
~~~

## المهمة
ضيف ترحيل 5: ~memberships(crew_id, user_id, role)~ بمفتاح ~(crew_id, user_id)~ و ~role~ واحد من ~owner, admin, member~، و ~audit(crew_id, seq, actor_id, action, target, at)~. في ~crews.go~ على ~Crews{Auth}~:

~CreateCrew(ctx, actor, name)~ تنشئ طاقم بـ ID عشوائي، والمنشئ ~owner~، وسطر audit ~crew.create~ هدفه ID الطاقم. ~Role(ctx, crewID, userID)~ ترجع الدور أو ~ErrNotFound~. ~AddMember(ctx, actor, crewID, email, role)~: الغريب ~ErrNotFound~؛ العضو العادي ~ErrForbidden~؛ دور غير ~admin~ أو ~member~ يرجع ~ErrInvalidRole~؛ دور قد دور المنادي أو أعلى ~ErrForbidden~ (الـ owner يضيف admin أو member، والـ admin يضيف member بس)؛ إيميل مش متسجّل ~ErrNotFound~؛ عضو موجود ~ErrConflict~. النجاح بيكتب ~member.add~ هدفه ID المستخدم المضاف. ~Audit(ctx, actor, crewID)~ للـ admin والـ owner بس، بترتيب ~seq~.

~NewCommandAPI(c)~ بيخدم كل مسارات ~NewAPI~، وكل المسارات دي محتاجة كوكي ~session~ (من غيرها 401 ~unauthorized~): ~POST /v1/crews~ بـ ~{"name"}~ يرجع 201 و ~{"id","name"}~ (اسم فاضي 400 ~bad_request~). ~POST /v1/crews/{crew}/members~ بـ ~{"email","role"}~ يرجع 201. ~GET /v1/crews/{crew}/audit~ يرجع 200 و ~{"entries":[{"actor","action","target"}]}~. ~GET /v1/crews/{crew}/tasks/{task}~ دلوقتي للأعضاء بس. الأخطاء: 404 ~not_found~، و 403 ~forbidden~، و 400 ~invalid_role~، و 409 ~conflict~.
`,
		Hints: []string{"كل صلاحية تخص العضوية داخل الطاقم، مش المستخدم بمعزل عن السياق.", "رتّب خطوات التحقق قبل تنفيذ تغيير العضوية أو كتابة audit.", "ارجع لطرق Crews في starter وحدد مين actor ومين target في كل طلب."},
		En: &LevelText{Title: "Isolated Crews", Goal: "Add crews with owner/admin/member roles, full isolation between crews (outsiders see not found), escalation prevention, and an audit log", Files: []LabFile{{Name: "main.go", Starter: ch371Starter}, {Name: "auth.go", Starter: ch372Starter}, {Name: "crews.go", Starter: ch373Starter}},
			Lesson: `# Command Center: Isolated Crews

## Before you write
- If someone outside a crew asks about a task inside it, should they learn that it exists?
- Who may grant a role, and how high a role may they grant?

## What this builds on
~main.go~ and ~auth.go~ carry over from earlier stages, and the new file is ~crews.go~. Sign-in comes from ~37-2-auth~, handler-wrapping middleware from the Base Gateway chapter (~26-3-router~), and transactions from ~29-3-transfer~. Only the integration decisions are new.

## Isolation means not found
The crew is the tenant. Any request against a crew you do not belong to returns 404, exactly like a crew that does not exist. A 403 would tell an outsider the ID is real and let them enumerate crews.
~~~go
role, err := c.Role(ctx, crewID, user.ID)
if errors.Is(err, ErrNotFound) { return ErrNotFound } // outsider = missing
~~~
403 is only for members: they already know the crew exists, but their role is not enough.

## Roles as ranks
Store the role as text with ~CHECK (role IN (...))~ and compare in Go with numeric ranks from a map. The no-escalation rule: nobody grants a role equal to or above their own.
~~~go
var level = map[string]int{"cadet": 1, "pilot": 2}
if level[grant] >= level[mine] { return ErrForbidden }
~~~

## Audit in the same transaction
The audit row is written in the same transaction as the change. If the change rolls back, so does the row; if the row fails, the change does not happen. Order by a per-crew ~seq~ (~MAX(seq) + 1~ inside the transaction), not by time, since two requests can share a second.
~~~go
tx.QueryRowContext(ctx, "SELECT COALESCE(MAX(seq), 0) + 1 FROM audit WHERE crew_id = $1", id).Scan(&seq)
~~~

## Common mistakes
~~~go
if role == "" { http.Error(w, "forbidden", 403) } // reveals the crew exists
if mine == "admin" || mine == "owner" { grant(role) } // an admin can mint a new owner
db.ExecContext(ctx, insertMember); db.ExecContext(ctx, insertAudit) // audit outside the transaction
~~~

## Your mission
Add migration 5: ~memberships(crew_id, user_id, role)~ keyed by ~(crew_id, user_id)~ with ~role~ one of ~owner, admin, member~, and ~audit(crew_id, seq, actor_id, action, target, at)~. In ~crews.go~ on ~Crews{Auth}~:

~CreateCrew(ctx, actor, name)~ creates a crew with a random ID, makes the creator ~owner~, and writes audit ~crew.create~ targeting the crew ID. ~Role(ctx, crewID, userID)~ returns the role or ~ErrNotFound~. ~AddMember(ctx, actor, crewID, email, role)~: an outsider gets ~ErrNotFound~; a plain member ~ErrForbidden~; a role other than ~admin~ or ~member~ returns ~ErrInvalidRole~; a role equal to or above the caller's returns ~ErrForbidden~ (an owner adds admins or members, an admin adds members only); an unregistered email ~ErrNotFound~; an existing member ~ErrConflict~. Success writes ~member.add~ targeting the added user's ID. ~Audit(ctx, actor, crewID)~ is for admins and owners only, in ~seq~ order.

~NewCommandAPI(c)~ serves every ~NewAPI~ route, and all of the following need a ~session~ cookie (without one: 401 ~unauthorized~): ~POST /v1/crews~ with ~{"name"}~ returns 201 and ~{"id","name"}~ (a blank name is 400 ~bad_request~). ~POST /v1/crews/{crew}/members~ with ~{"email","role"}~ returns 201. ~GET /v1/crews/{crew}/audit~ returns 200 and ~{"entries":[{"actor","action","target"}]}~. ~GET /v1/crews/{crew}/tasks/{task}~ is now members only. Errors: 404 ~not_found~, 403 ~forbidden~, 400 ~invalid_role~, 409 ~conflict~.
`,
			LessonRevision: 1, Steps: ch37TenantsEnglishSteps,
			Hints: []string{"Every permission belongs to a membership in a crew, not to the user in isolation.", "Order authorization checks before changing membership or writing audit data.", "Use the Crews signatures in the starter to identify each request's actor and target."}},
		Lab: &Lab{ReqFile: "crews.go", Bonus: false, Hidden: ch373Hidden},
	})
}

const ch373Starter = `package main

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrForbidden   = errors.New("forbidden")
	ErrInvalidRole = errors.New("invalid role")
)

type AuditEntry struct {
	Actor  string ` + "`json:\"actor\"`" + `
	Action string ` + "`json:\"action\"`" + `
	Target string ` + "`json:\"target\"`" + `
}

type Crews struct{ Auth *Auth }

func (c *Crews) CreateCrew(ctx context.Context, actor User, name string) (Crew, error) {
	return Crew{}, errors.New("TODO")
}
func (c *Crews) Role(ctx context.Context, crewID, userID string) (string, error) {
	return "", errors.New("TODO")
}
func (c *Crews) AddMember(ctx context.Context, actor User, crewID, email, role string) error {
	return errors.New("TODO")
}
func (c *Crews) Audit(ctx context.Context, actor User, crewID string) ([]AuditEntry, error) {
	return nil, errors.New("TODO")
}

func NewCommandAPI(c *Crews) http.Handler { return NewAPI(c.Auth) }
`

// ch37People is shared by the hidden tests of stages 3 to 5: a signed-in
// cast of users and a cookie header for each.
const ch37People = `
type hiddenClock struct{ now time.Time }

func (c *hiddenClock) Now() time.Time { return c.now }

type hiddenPerson struct {
	user   User
	cookie map[string]string
}

func hiddenCommand(t *testing.T, db *sql.DB, seed int64, names ...string) (*Crews, *hiddenClock, map[string]hiddenPerson) {
	t.Helper()
	ctx := context.Background()
	if err := Migrate(ctx, db, len(Migrations)); err != nil {
		t.Fatal(err)
	}
	clock := &hiddenClock{time.Unix(1_900_000_000, 0)}
	a := &Auth{Store: NewStore(db), Now: clock.Now, Random: rand.New(rand.NewSource(seed)), TTL: 2 * time.Hour, Cost: bcrypt.MinCost}
	people := map[string]hiddenPerson{}
	for _, name := range names {
		u, err := a.Register(ctx, name+"@colony", name+"-password")
		if err != nil {
			t.Fatal(err)
		}
		token, err := a.Login(ctx, name+"@colony", name+"-password")
		if err != nil {
			t.Fatal(err)
		}
		people[name] = hiddenPerson{u, map[string]string{"Cookie": SessionCookie + "=" + token}}
	}
	return &Crews{Auth: a}, clock, people
}
`

const ch373Hidden = `package main

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
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)
` + ch37Backends + ch37People + `
func TestHiddenCrewRoles(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		ctx := context.Background()
		c, _, p := hiddenCommand(t, db, 3731, "owner", "admin", "member", "outsider", "spare", "extra")
		crew, err := c.CreateCrew(ctx, p["owner"].user, "Orion")
		if err != nil || crew.ID == "" || crew.Name != "Orion" {
			t.Fatal(crew, err)
		}
		if role, err := c.Role(ctx, crew.ID, p["owner"].user.ID); err != nil || role != "owner" {
			t.Fatal("creator role", role, err)
		}
		if _, err := c.Role(ctx, crew.ID, p["outsider"].user.ID); !errors.Is(err, ErrNotFound) {
			t.Fatal("outsider role", err)
		}
		steps := []struct {
			actor, email, role string
			want               error
		}{
			{"owner", "admin@colony", "admin", nil},
			{"admin", "member@colony", "member", nil},
			{"admin", "spare@colony", "admin", ErrForbidden},
			{"admin", "spare@colony", "owner", ErrInvalidRole},
			{"owner", "spare@colony", "owner", ErrInvalidRole},
			{"owner", "spare@colony", "captain", ErrInvalidRole},
			{"member", "spare@colony", "member", ErrForbidden},
			{"outsider", "spare@colony", "member", ErrNotFound},
			{"owner", "ghost@colony", "member", ErrNotFound},
			{"owner", "member@colony", "member", ErrConflict},
		}
		for _, s := range steps {
			err := c.AddMember(ctx, p[s.actor].user, crew.ID, s.email, s.role)
			if (s.want == nil && err != nil) || (s.want != nil && !errors.Is(err, s.want)) {
				t.Fatal("add member", s, err)
			}
		}
		if role, _ := c.Role(ctx, crew.ID, p["spare"].user.ID); role != "" {
			t.Fatal("denied grant still created membership", role)
		}
		entries, err := c.Audit(ctx, p["admin"].user, crew.ID)
		want := []AuditEntry{{p["owner"].user.ID, "crew.create", crew.ID}, {p["owner"].user.ID, "member.add", p["admin"].user.ID}, {p["admin"].user.ID, "member.add", p["member"].user.ID}}
		if err != nil || fmt.Sprint(entries) != fmt.Sprint(want) {
			t.Fatal("audit", entries, err)
		}
		if _, err := c.Audit(ctx, p["member"].user, crew.ID); !errors.Is(err, ErrForbidden) {
			t.Fatal("member audit", err)
		}
		if _, err := c.Audit(ctx, p["outsider"].user, crew.ID); !errors.Is(err, ErrNotFound) {
			t.Fatal("outsider audit", err)
		}
		other, err := c.CreateCrew(ctx, p["outsider"].user, "Lyra")
		if err != nil {
			t.Fatal(err)
		}
		if entries, _ := c.Audit(ctx, p["outsider"].user, other.ID); len(entries) != 1 {
			t.Fatal("audit leaks across crews", entries)
		}
		if err := c.AddMember(ctx, p["owner"].user, other.ID, "extra@colony", "member"); !errors.Is(err, ErrNotFound) {
			t.Fatal("owner of one crew acted on another", err)
		}
	})
}

func TestHiddenCrewHTTPIsolation(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		c, _, p := hiddenCommand(t, db, 3732, "vega", "rigel", "deneb")
		srv := httptest.NewServer(NewCommandAPI(c))
		defer srv.Close()
		cl := srv.Client()
		code, _, body := hiddenGetJSON(t, cl, "POST", srv.URL+"/v1/crews", "{\"name\":\"Vega Wing\"}", nil)
		if code != 401 || body["error"] != "unauthorized" {
			t.Fatal("anonymous crew create", code, body)
		}
		code, _, body = hiddenGetJSON(t, cl, "POST", srv.URL+"/v1/crews", "{\"name\":\"  \"}", p["vega"].cookie)
		if code != 400 || body["error"] != "bad_request" {
			t.Fatal("blank crew", code, body)
		}
		code, _, body = hiddenGetJSON(t, cl, "POST", srv.URL+"/v1/crews", "{\"name\":\"Vega Wing\"}", p["vega"].cookie)
		crewID, _ := body["id"].(string)
		if code != 201 || crewID == "" || body["name"] != "Vega Wing" {
			t.Fatal("create crew", code, body)
		}
		if err := c.Auth.Store.CreateTask(context.Background(), Task{"secret", crewID, "خريطة المخزن", false, 1}); err != nil {
			t.Fatal(err)
		}
		code, _, body = hiddenGetJSON(t, cl, "POST", srv.URL+"/v1/crews/"+crewID+"/members", "{\"email\":\"rigel@colony\",\"role\":\"member\"}", p["vega"].cookie)
		if code != 201 {
			t.Fatal("add member", code, body)
		}
		missing := "/v1/crews/no-such-crew"
		for _, path := range []string{"/v1/crews/" + crewID, missing} {
			code, _, body = hiddenGetJSON(t, cl, "GET", srv.URL+path+"/tasks/secret", "", p["deneb"].cookie)
			if code != 404 || body["error"] != "not_found" {
				t.Fatal("outsider task", path, code, body)
			}
			code, _, body = hiddenGetJSON(t, cl, "GET", srv.URL+path+"/audit", "", p["deneb"].cookie)
			if code != 404 || body["error"] != "not_found" {
				t.Fatal("outsider audit", path, code, body)
			}
			code, _, body = hiddenGetJSON(t, cl, "POST", srv.URL+path+"/members", "{\"email\":\"deneb@colony\",\"role\":\"member\"}", p["deneb"].cookie)
			if code != 404 || body["error"] != "not_found" {
				t.Fatal("outsider self-invite", path, code, body)
			}
		}
		code, _, body = hiddenGetJSON(t, cl, "GET", srv.URL+"/v1/crews/"+crewID+"/tasks/secret", "", p["rigel"].cookie)
		if code != 200 || body["title"] != "خريطة المخزن" {
			t.Fatal("member task", code, body)
		}
		code, _, _ = hiddenGetJSON(t, cl, "GET", srv.URL+"/v1/crews/"+crewID+"/tasks/secret", "", nil)
		if code != 401 {
			t.Fatal("anonymous task read", code)
		}
		for _, c2 := range []struct {
			body, err string
			code      int
		}{{"{\"email\":\"deneb@colony\",\"role\":\"member\"}", "forbidden", 403}, {"{\"email\":\"deneb@colony\",\"role\":\"owner\"}", "forbidden", 403}} {
			code, _, body = hiddenGetJSON(t, cl, "POST", srv.URL+"/v1/crews/"+crewID+"/members", c2.body, p["rigel"].cookie)
			if code != c2.code || body["error"] != c2.err {
				t.Fatal("member escalation", c2, code, body)
			}
		}
		code, _, body = hiddenGetJSON(t, cl, "POST", srv.URL+"/v1/crews/"+crewID+"/members", "{\"email\":\"deneb@colony\",\"role\":\"pilot\"}", p["vega"].cookie)
		if code != 400 || body["error"] != "invalid_role" {
			t.Fatal("invalid role", code, body)
		}
		code, _, body = hiddenGetJSON(t, cl, "POST", srv.URL+"/v1/crews/"+crewID+"/members", "{\"email\":\"rigel@colony\",\"role\":\"member\"}", p["vega"].cookie)
		if code != 409 || body["error"] != "conflict" {
			t.Fatal("duplicate member", code, body)
		}
		code, _, body = hiddenGetJSON(t, cl, "GET", srv.URL+"/v1/crews/"+crewID+"/audit", "", p["rigel"].cookie)
		if code != 403 || body["error"] != "forbidden" {
			t.Fatal("member audit", code, body)
		}
		code, _, body = hiddenGetJSON(t, cl, "GET", srv.URL+"/v1/crews/"+crewID+"/audit", "", p["vega"].cookie)
		entries, _ := body["entries"].([]any)
		if code != 200 || len(entries) != 2 {
			t.Fatal("owner audit", code, body)
		}
		last, _ := entries[1].(map[string]any)
		if last["action"] != "member.add" || last["actor"] != p["vega"].user.ID || last["target"] != p["rigel"].user.ID {
			t.Fatal("audit entry", last)
		}
	})
}
`

func init() {
	add(&Level{ID: "37-4-tasks", Chapter: 37, Kind: "lab", CarryFrom: "37-3-tenants",
		Title: "لوحة المهام", Goal: "ضيف مهام بـ Idempotency-Key ونسخة متفائلة وصفحات بـ cursor ثابت من غير تكرار", Setup: seededLabRoom,
		LessonRevision: 1, Steps: ch37TasksArabicSteps,
		Requires: []string{"sel:strconv.Atoi", "sel:base64.RawURLEncoding"}, Deps: []string{"modernc.org/sqlite", "github.com/jackc/pgx/v5", "golang.org/x/crypto"}, Bridge: []int{100}, Race: false, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch371Starter}, {Name: "auth.go", Starter: ch372Starter}, {Name: "crews.go", Starter: ch373Starter}, {Name: "tasks.go", Starter: ch374Starter}},
		Lesson: `# مركز القيادة: لوحة المهام

## قبل ما تكتب
- الشبكة قطعت بعد ما السيرفر أنشأ المهمة وقبل ما الرد يوصل. العميل هيعيد الطلب: المفروض يحصل إيه؟
- اتنين فتحوا نفس المهمة وعدّلوها في نفس الوقت. مين تعديله يكسب؟

## بتبني على إيه
كل الملفات اللي فاتت جاية زي ما هي، والجديد ~tasks.go~. العضوية والـ 404 للغريب من ~37-3-tenants~، والـ transactions من ~29-3-transfer~ (اللي فيها request ID مكرر)، و ~sha256~ من فصل الخزنة. الجديد قرارات الدمج بس.

## مفتاح التكرار
العميل بيبعت هيدر ~Idempotency-Key~ بقيمة هو اختارها. أول طلب بالمفتاح ده بيتنفّذ ويتحفظ رده؛ أي إعادة بنفس المفتاح ونفس الجسم بترجع نفس الرد المحفوظ بالظبط من غير مهمة تانية. نفس المفتاح بجسم مختلف غلطة من العميل: 422.
~~~go
key := r.Header.Get("Idempotency-Key")
sum := sha256.Sum256(raw) // بصمة الجسم
~~~
المفتاح بيخص المستخدم والطاقم: مفتاح Vega مبيأثرش على مفتاح Rigel.

## احجز المفتاح أول حاجة
الحجز (~INSERT ... ON CONFLICT DO NOTHING~) هو أول أمر في الـ transaction. كده الطلب التاني بنفس المفتاح بيستنى لحد ما الأول يعمل commit، وبعدين يلاقي الصف ويرجّع الرد المحفوظ. لو قريت الأول وبعدين كتبت، طلبين متوازيين هيشوفوا «مفيش» ويعملوا مهمتين.

## نسخة متفائلة
كل مهمة ليها ~version~. التعديل بيبعت النسخة اللي شافها، والـ ~UPDATE~ بيشترطها:
~~~sql
UPDATE beacons SET label = $1, version = version + 1 WHERE id = $2 AND version = $3
~~~
~RowsAffected~ بصفر معناها حد سبقك: رجّع 409 ومعاه النسخة الحالية عشان العميل يقرا تاني ويقرر.

## صفحات بـ cursor
~OFFSET~ بيكرر أو يفوّت صفوف لو حد ضاف مهمة بين صفحتين. بدل كده رتّب بمفتاح ثابت (~id~) واطلب اللي بعد آخر واحد شفته: ~WHERE id > $2 ORDER BY id LIMIT $3~. الـ cursor هو آخر ~id~ متغلّف بـ ~base64.RawURLEncoding~ عشان العميل يعامله كنص معتم.
~~~go
next := base64.RawURLEncoding.EncodeToString([]byte(last.ID))
~~~
اطلب ~limit + 1~ صف: لو جالك الزيادة، فيه صفحة كمان. ~strconv.Atoi~ بتحوّل ~limit~ من الـ query لرقم.

## أخطاء شائعة
~~~go
if exists(key) { return saved } ; insert(key) // قراية ثم كتابة: طلبين متوازيين يعدّوا
db.ExecContext(ctx, "UPDATE tasks SET title = $1 WHERE id = $2", title, id) // آخر واحد يمسح تعديل اللي قبله
db.QueryContext(ctx, "SELECT ... LIMIT $1 OFFSET $2", n, page*n) // صفوف بتتكرر لما حد يضيف
~~~

## المهمة
ضيف ترحيل 6: جدول للمفاتيح بمفتاح ~(crew_id, user_id, key)~ فيه البصمة والحالة والجسم المحفوظ. في ~tasks.go~ اعمل ~NewTaskAPI(c)~ اللي بيخدم كل مسارات ~NewCommandAPI~ ويضيف (كلهم بكوكي، والغريب 404 ~not_found~):

~POST /v1/crews/{crew}/tasks~ بـ ~{"title"}~ وهيدر ~Idempotency-Key~ (من 1 لـ 128 حرف، وإلا 400 ~missing_idempotency_key~). عنوان فاضي أو JSON مكسور: 400 ~bad_request~. النجاح 201 والمهمة بـ ~version~ 1 و ~done~ false. الإعادة بنفس المفتاح والجسم: نفس الحالة ونفس الجسم بالبايت، وهيدر ~Idempotent-Replayed: true~، حتى لو الطلبات متوازية. نفس المفتاح بجسم تاني: 422 ~idempotency_mismatch~.

~PATCH /v1/crews/{crew}/tasks/{task}~ بـ ~{"title"?, "done"?, "version"}~: ~version~ لازم 1 أو أكتر، والعنوان لو موجود مش فاضي، وإلا 400 ~bad_request~. النسخة الصح: 200 والمهمة بعد التعديل والنسخة زادت 1. النسخة القديمة: 409 و ~{"error":"version_conflict","version":N}~ بالنسخة الحالية. مهمة مش موجودة: 404.

~GET /v1/crews/{crew}/tasks?limit=&cursor=~: ~limit~ افتراضي 20 ومن 1 لـ 100، و cursor مش base64 صح: 400 ~bad_request~. يرجع 200 و ~{"tasks":[...],"next":"..."}~ بترتيب ~id~، و ~next~ فاضي في آخر صفحة. اللف على كل الصفحات بيرجّع كل مهمة مرة واحدة، حتى لو اتضافت مهام بين الصفحات.
`,
		Hints: []string{"فكر في التكرار والتعارض كحالتين مختلفتين من طلبات الشبكة.", "قسّم POST وPATCH وGET pages إلى تدفقات منفصلة قبل مشاركة أجزاء التخزين.", "ارجع إلى NewTaskAPI في starter وحدد البيانات التي تخص كل route."},
		En: &LevelText{Title: "Task Board", Goal: "Add tasks with Idempotency-Key, optimistic versions, and stable cursor pages without duplicates", Files: []LabFile{{Name: "main.go", Starter: ch371Starter}, {Name: "auth.go", Starter: ch372Starter}, {Name: "crews.go", Starter: ch373Starter}, {Name: "tasks.go", Starter: ch374Starter}},
			Lesson: `# Command Center: Task Board

## Before you write
- The network drops after the server created a task but before the reply arrives. The client retries: what should happen?
- Two people opened the same task and edited it at once. Whose edit wins?

## What this builds on
All earlier files carry over as they are, and the new file is ~tasks.go~. Membership and outsider 404s come from ~37-3-tenants~, transactions from ~29-3-transfer~ (with its repeated request ID), and ~sha256~ from the Vault chapter. Only the integration decisions are new.

## The idempotency key
The client sends an ~Idempotency-Key~ header with a value it chose. The first request with that key runs and its reply is saved; any retry with the same key and the same body returns exactly that saved reply without a second task. The same key with a different body is a client bug: 422.
~~~go
key := r.Header.Get("Idempotency-Key")
sum := sha256.Sum256(raw) // body fingerprint
~~~
Keys belong to a user and a crew: Vega's key never affects Rigel's.

## Claim the key first
The claim (~INSERT ... ON CONFLICT DO NOTHING~) is the first statement in the transaction. A second request with the same key then waits until the first commits, finds the row, and returns the saved reply. If you read first and write later, two concurrent requests both see "nothing" and create two tasks.

## Optimistic versions
Every task has a ~version~. An edit sends the version it saw, and the ~UPDATE~ requires it:
~~~sql
UPDATE beacons SET label = $1, version = version + 1 WHERE id = $2 AND version = $3
~~~
Zero ~RowsAffected~ means someone got there first: return 409 with the current version so the client can re-read and decide.

## Cursor pages
~OFFSET~ repeats or skips rows when someone adds a task between pages. Instead, order by a stable key (~id~) and ask for what comes after the last one you saw: ~WHERE id > $2 ORDER BY id LIMIT $3~. The cursor is that last ~id~ wrapped in ~base64.RawURLEncoding~ so the client treats it as opaque text.
~~~go
next := base64.RawURLEncoding.EncodeToString([]byte(last.ID))
~~~
Fetch ~limit + 1~ rows: if the extra one arrives, there is another page. ~strconv.Atoi~ turns ~limit~ from the query into a number.

## Common mistakes
~~~go
if exists(key) { return saved } ; insert(key) // read then write: two concurrent requests both pass
db.ExecContext(ctx, "UPDATE tasks SET title = $1 WHERE id = $2", title, id) // the last writer erases the previous edit
db.QueryContext(ctx, "SELECT ... LIMIT $1 OFFSET $2", n, page*n) // rows repeat when someone inserts
~~~

## Your mission
Add migration 6: a key table keyed by ~(crew_id, user_id, key)~ holding the fingerprint, status, and saved body. In ~tasks.go~ write ~NewTaskAPI(c)~, which serves every ~NewCommandAPI~ route and adds (all need the cookie, and outsiders get 404 ~not_found~):

~POST /v1/crews/{crew}/tasks~ with ~{"title"}~ and an ~Idempotency-Key~ header (1 to 128 characters, else 400 ~missing_idempotency_key~). A blank title or malformed JSON: 400 ~bad_request~. Success: 201 and the task with ~version~ 1 and ~done~ false. A retry with the same key and body: the same status and byte-identical body, plus ~Idempotent-Replayed: true~, even when requests are concurrent. The same key with another body: 422 ~idempotency_mismatch~.

~PATCH /v1/crews/{crew}/tasks/{task}~ with ~{"title"?, "done"?, "version"}~: ~version~ must be 1 or more and a present title must be nonblank, else 400 ~bad_request~. The right version: 200 and the updated task with the version incremented by 1. A stale version: 409 and ~{"error":"version_conflict","version":N}~ with the current version. A missing task: 404.

~GET /v1/crews/{crew}/tasks?limit=&cursor=~: ~limit~ defaults to 20 and must be 1 to 100, and a cursor that is not valid base64 is 400 ~bad_request~. Returns 200 and ~{"tasks":[...],"next":"..."}~ in ~id~ order, with an empty ~next~ on the last page. Walking every page returns each task exactly once, even when tasks are added between pages.
`,
			LessonRevision: 1, Steps: ch37TasksEnglishSteps,
			Hints: []string{"Treat retries and conflicts as different outcomes of network requests.", "Separate POST, PATCH, and GET-page flows before sharing storage helpers.", "Use NewTaskAPI in the starter to identify the data each route needs."}},
		Lab: &Lab{ReqFile: "tasks.go", Bonus: false, Hidden: ch374Hidden},
	})
}

const ch374Starter = `package main

import "net/http"

func NewTaskAPI(c *Crews) http.Handler { return NewCommandAPI(c) }
`

const ch374Hidden = `package main

import (
	"bytes"
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
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)
` + ch37Backends + ch37People + `
type hiddenLockedRandom struct {
	mu sync.Mutex
	r  io.Reader
}

func (l *hiddenLockedRandom) Read(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.r.Read(p)
}

func hiddenRaw(t *testing.T, c *http.Client, method, target, body string, header map[string]string) (int, http.Header, []byte) {
	t.Helper()
	req, _ := http.NewRequest(method, target, strings.NewReader(body))
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, resp.Header, raw
}

func hiddenWith(h map[string]string, k, v string) map[string]string {
	out := map[string]string{k: v}
	for a, b := range h {
		out[a] = b
	}
	return out
}

func hiddenTaskSetup(t *testing.T, db *sql.DB, seed int64) (*httptest.Server, map[string]hiddenPerson, string) {
	t.Helper()
	c, _, p := hiddenCommand(t, db, seed, "vega", "rigel", "deneb")
	c.Auth.Random = &hiddenLockedRandom{r: c.Auth.Random}
	crew, err := c.CreateCrew(context.Background(), p["vega"].user, "Board")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.AddMember(context.Background(), p["vega"].user, crew.ID, "rigel@colony", "member"); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(NewTaskAPI(c))
	t.Cleanup(srv.Close)
	return srv, p, crew.ID
}

func TestHiddenIdempotentCreate(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		srv, p, crew := hiddenTaskSetup(t, db, 3741)
		cl := srv.Client()
		base := srv.URL + "/v1/crews/" + crew + "/tasks"
		vega := p["vega"].cookie
		code, _, raw := hiddenRaw(t, cl, "POST", base, "{\"title\":\"x\"}", vega)
		if code != 400 || !bytes.Contains(raw, []byte("missing_idempotency_key")) {
			t.Fatal("missing key", code, string(raw))
		}
		code, _, raw = hiddenRaw(t, cl, "POST", base, "{\"title\":\"x\"}", hiddenWith(vega, "Idempotency-Key", strings.Repeat("k", 129)))
		if code != 400 || !bytes.Contains(raw, []byte("missing_idempotency_key")) {
			t.Fatal("long key", code, string(raw))
		}
		for _, bad := range []string{"{\"title\":\"  \"}", "{oops"} {
			code, _, raw = hiddenRaw(t, cl, "POST", base, bad, hiddenWith(vega, "Idempotency-Key", "bad-"+bad))
			if code != 400 || !bytes.Contains(raw, []byte("bad_request")) {
				t.Fatal("bad body", bad, code, string(raw))
			}
		}
		code, _, raw = hiddenRaw(t, cl, "POST", base, "{\"title\":\"x\"}", hiddenWith(p["deneb"].cookie, "Idempotency-Key", "k"))
		if code != 404 {
			t.Fatal("outsider create", code, string(raw))
		}
		body := "{\"title\":\"إصلاح الهوائي\"}"
		key := hiddenWith(vega, "Idempotency-Key", fmt.Sprintf("op-%d", rand.New(rand.NewSource(3741)).Int63()))
		code, h, first := hiddenRaw(t, cl, "POST", base, body, key)
		var task map[string]any
		json.Unmarshal(first, &task)
		if code != 201 || task["title"] != "إصلاح الهوائي" || task["version"] != float64(1) || task["done"] != false || task["crew_id"] != crew || h.Get("Idempotent-Replayed") != "" {
			t.Fatal("create", code, string(first), h)
		}
		for i := 0; i < 3; i++ {
			code, h, again := hiddenRaw(t, cl, "POST", base, body, key)
			if code != 201 || !bytes.Equal(bytes.TrimSpace(again), bytes.TrimSpace(first)) || h.Get("Idempotent-Replayed") != "true" {
				t.Fatal("replay", code, string(again), h)
			}
		}
		code, _, raw = hiddenRaw(t, cl, "POST", base, "{\"title\":\"other\"}", key)
		if code != 422 || !bytes.Contains(raw, []byte("idempotency_mismatch")) {
			t.Fatal("mismatch", code, string(raw))
		}
		code, _, raw = hiddenRaw(t, cl, "POST", base, body, hiddenWith(p["rigel"].cookie, "Idempotency-Key", key["Idempotency-Key"]))
		if code != 201 || bytes.Equal(bytes.TrimSpace(raw), bytes.TrimSpace(first)) {
			t.Fatal("keys must be per user", code, string(raw))
		}
		var wg sync.WaitGroup
		results := make([]string, 8)
		codes := make([]int, 8)
		for i := range results {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				req, _ := http.NewRequest("POST", base, strings.NewReader("{\"title\":\"parallel\"}"))
				req.Header.Set("Cookie", vega["Cookie"])
				req.Header.Set("Idempotency-Key", "parallel-key")
				resp, err := cl.Do(req)
				if err != nil {
					return
				}
				defer resp.Body.Close()
				raw, _ := io.ReadAll(resp.Body)
				codes[i], results[i] = resp.StatusCode, strings.TrimSpace(string(raw))
			}(i)
		}
		wg.Wait()
		for i := range results {
			if codes[i] != 201 || results[i] != results[0] {
				t.Fatal("concurrent retries", codes, results)
			}
		}
		var n int
		if err := db.QueryRow("SELECT COUNT(*) FROM tasks WHERE crew_id = $1 AND title = 'parallel'", crew).Scan(&n); err != nil || n != 1 {
			t.Fatal("concurrent retries created tasks:", n, err)
		}
	})
}

func TestHiddenOptimisticVersion(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		srv, p, crew := hiddenTaskSetup(t, db, 3742)
		cl := srv.Client()
		base := srv.URL + "/v1/crews/" + crew + "/tasks"
		_, _, raw := hiddenRaw(t, cl, "POST", base, "{\"title\":\"seal hatch\"}", hiddenWith(p["vega"].cookie, "Idempotency-Key", "v1"))
		var task map[string]any
		json.Unmarshal(raw, &task)
		id, _ := task["id"].(string)
		item := base + "/" + id
		code, _, body := hiddenGetJSON(t, cl, "PATCH", item, "{\"done\":true,\"version\":1}", p["rigel"].cookie)
		if code != 200 || body["done"] != true || body["title"] != "seal hatch" || body["version"] != float64(2) {
			t.Fatal("patch", code, body)
		}
		code, _, body = hiddenGetJSON(t, cl, "PATCH", item, "{\"title\":\"stale\",\"version\":1}", p["vega"].cookie)
		if code != 409 || body["error"] != "version_conflict" || body["version"] != float64(2) {
			t.Fatal("stale patch", code, body)
		}
		code, _, body = hiddenGetJSON(t, cl, "GET", item, "", p["vega"].cookie)
		if code != 200 || body["title"] != "seal hatch" || body["version"] != float64(2) {
			t.Fatal("stale patch changed the task", code, body)
		}
		for _, bad := range []string{"{\"title\":\"x\"}", "{\"title\":\" \",\"version\":2}", "{nope"} {
			code, _, body = hiddenGetJSON(t, cl, "PATCH", item, bad, p["vega"].cookie)
			if code != 400 || body["error"] != "bad_request" {
				t.Fatal("bad patch", bad, code, body)
			}
		}
		code, _, _ = hiddenGetJSON(t, cl, "PATCH", base+"/missing", "{\"version\":1}", p["vega"].cookie)
		if code != 404 {
			t.Fatal("missing task", code)
		}
		code, _, _ = hiddenGetJSON(t, cl, "PATCH", item, "{\"version\":2}", p["deneb"].cookie)
		if code != 404 {
			t.Fatal("outsider patch", code)
		}
		var wg sync.WaitGroup
		wins := make([]int, 6)
		for i := range wins {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				req, _ := http.NewRequest("PATCH", item, strings.NewReader(fmt.Sprintf("{\"title\":\"racer %d\",\"version\":2}", i)))
				req.Header.Set("Cookie", p["vega"].cookie["Cookie"])
				resp, err := cl.Do(req)
				if err == nil {
					wins[i] = resp.StatusCode
					resp.Body.Close()
				}
			}(i)
		}
		wg.Wait()
		ok := 0
		for _, code := range wins {
			if code == 200 {
				ok++
			} else if code != 409 {
				t.Fatal("concurrent patch status", wins)
			}
		}
		if ok != 1 {
			t.Fatal("exactly one concurrent edit must win", wins)
		}
	})
}

func TestHiddenCursorPages(t *testing.T) {
	eachBackend(t, func(t *testing.T, db *sql.DB) {
		srv, p, crew := hiddenTaskSetup(t, db, 3743)
		cl := srv.Client()
		base := srv.URL + "/v1/crews/" + crew + "/tasks"
		rng := rand.New(rand.NewSource(3743))
		total := 25 + rng.Intn(10)
		for i := 0; i < total; i++ {
			code, _, raw := hiddenRaw(t, cl, "POST", base, fmt.Sprintf("{\"title\":\"t%d\"}", i), hiddenWith(p["vega"].cookie, "Idempotency-Key", fmt.Sprint("p", i)))
			if code != 201 {
				t.Fatal(code, string(raw))
			}
		}
		for _, q := range []string{"?limit=0", "?limit=101", "?limit=x", "?cursor=!!!"} {
			code, _, body := hiddenGetJSON(t, cl, "GET", base+q, "", p["vega"].cookie)
			if code != 400 || body["error"] != "bad_request" {
				t.Fatal("bad query", q, code, body)
			}
		}
		code, _, body := hiddenGetJSON(t, cl, "GET", base, "", p["vega"].cookie)
		if tasks, _ := body["tasks"].([]any); code != 200 || len(tasks) != 20 || body["next"] == "" {
			t.Fatal("default limit", code, len(tasks))
		}
		code, _, _ = hiddenGetJSON(t, cl, "GET", base, "", p["deneb"].cookie)
		if code != 404 {
			t.Fatal("outsider list", code)
		}
		seen := map[string]bool{}
		cursor, pages, added := "", 0, 0
		for {
			target := base + "?limit=4"
			if cursor != "" {
				target += "&cursor=" + url.QueryEscape(cursor)
			}
			code, _, body := hiddenGetJSON(t, cl, "GET", target, "", p["rigel"].cookie)
			tasks, _ := body["tasks"].([]any)
			if code != 200 || len(tasks) > 4 {
				t.Fatal("page", code, body)
			}
			for _, raw := range tasks {
				task := raw.(map[string]any)
				id := task["id"].(string)
				if seen[id] || task["crew_id"] != crew {
					t.Fatal("duplicate across pages", id)
				}
				seen[id] = true
			}
			pages++
			if pages == 2 {
				for i := 0; i < 3; i++ {
					hiddenRaw(t, cl, "POST", base, fmt.Sprintf("{\"title\":\"late%d\"}", i), hiddenWith(p["vega"].cookie, "Idempotency-Key", fmt.Sprint("late", i)))
				}
				added = 3
			}
			next, _ := body["next"].(string)
			if next == "" {
				break
			}
			cursor = next
		}
		var n int
		db.QueryRow("SELECT COUNT(*) FROM tasks WHERE crew_id = $1", crew).Scan(&n)
		if n != total+added {
			t.Fatal("task count", n)
		}
		if len(seen) < total { // every task that existed before the walk was seen
			t.Fatal("pages ended early")
		}
		_ = errors.Is
	})
}
`

func init() {
	add(&Level{ID: "37-5-operate", Chapter: 37, Kind: "lab", CarryFrom: "37-4-tasks",
		Title: "غرفة التشغيل", Goal: "لف الـ API بـ request ID وشكل أخطاء ثابت وmetrics وحد طلبات وسجل من غير أسرار وإغلاق بيصفّي الطلبات الجارية", Setup: seededLabRoom,
		LessonRevision: 1, Steps: ch37OpsArabicSteps,
		Requires: []string{"recover", "mutex"}, Deps: []string{"modernc.org/sqlite", "github.com/jackc/pgx/v5", "golang.org/x/crypto"}, Bridge: []int{100, 91, 96}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch371Starter}, {Name: "auth.go", Starter: ch372Starter}, {Name: "crews.go", Starter: ch373Starter}, {Name: "tasks.go", Starter: ch374Starter}, {Name: "ops.go", Starter: ch375Starter}},
		Lesson: `# مركز القيادة: غرفة التشغيل

## قبل ما تكتب
- مستخدم بعتلك رسالة «جالي خطأ». إزاي تلاقي طلبه بالظبط في السجل؟
- السيرفر بيتقفل للتحديث وفيه طلب في النص. يكمّل ولا يتقطع؟

## بتبني على إيه
كل ملفات المراحل اللي فاتت جاية زي ما هي، والجديد ~ops.go~. الـ middleware وظرف الأخطاء من ~26-4-envelope~، وحد الطلبات من ~27-4-rate-limit~، والإغلاق المرتب من ~28-4-graceful~، و ~recover~ من فصل الدوال والأخطاء، و ~sync.Mutex~ من فصل السباق (12). الـ metrics والسجلات المنظمة هي نفس فكرة فصل التشغيل والمراقبة (36)، والحماية من الضغط نفس روح فصل الوسطاء (33). الجديد قرارات الدمج بس.

## request ID في كل رد
كل طلب بياخد ID: لو العميل بعت ~X-Request-ID~ سليم (حروف وأرقام وشرطة، لحد 64) نستخدمه، وإلا نولّد واحد من ~Ops.NewID~ المحقونة. الـ ID بيرجع في هيدر الرد، وفي جسم أي خطأ، وفي سطر السجل. كده رسالة المستخدم توصّلك لسطر واحد.
~~~go
var safeID = regexp.MustCompile("^[a-z0-9]{1,8}$")
if !safeID.MatchString(id) { id = newID() }
~~~
متقبلش أي ID من العميل: نص فيه سطر جديد ممكن يزوّر سطر في السجل.

## شكل خطأ واحد
كل المراحل اللي فاتت بترجع ~{"error":"..."}~. الـ middleware بيضيف ~request_id~ لأي رد حالته 400 أو أكتر. عشان يعمل كده، بيلف ~ResponseWriter~: الردود الناجحة بتعدّي زي ما هي، وردود الأخطاء بتتجمّع في buffer وتتكتب من جديد.
~~~go
type capture struct {
 http.ResponseWriter
 status int
}
func (c *capture) WriteHeader(s int) { c.status = s; c.ResponseWriter.WriteHeader(s) }
~~~
الـ panic في أي handler بيتحوّل لـ 500 و ~{"error":"internal"}~ بـ ~defer~ و ~recover~، والسيرفر يكمّل.

## metrics كنص
~GET /metrics~ بيرجع نص بصيغة Prometheus المكتوبة باليد: سطر لكل حالة وعدّاد. مش محتاجين مكتبة؛ الصيغة بسيطة، والعدّادات محمية بـ ~sync.Mutex~ لأن الطلبات متوازية.
~~~text
command_requests_total{code="200"} 12
command_requests_in_flight 1
~~~

## حد الطلبات بساعة محقونة
نافذة ثابتة لكل عنوان IP (من ~net.SplitHostPort(r.RemoteAddr)~، مش من هيدر العميل يقدر يزوّره). النافذة بتبدأ مع أول طلب وعمرها ~Window~. الزيادة 429 مع ~Retry-After~ بالثواني الباقية مقرّبة لفوق.

## سجل من غير أسرار
سطر JSON لكل طلب فيه ~request_id~ و ~method~ و ~path~ و ~status~ بس. ~r.URL.Path~ من غير الـ query، ومن غير أي هيدر: الكوكي والـ tokens والباسوردات عمرهم ما يدخلوا السجل.

## إغلاق بنافذة تصفية
~Serve(ctx, ln, h, drain)~: لما ~ctx~ يتلغي، ~Shutdown~ بيقفل الـ listener فوراً ويستنى الطلبات الجارية لحد ~drain~. لو خلصت، ~Serve~ ترجع ~nil~. لو النافذة خلصت، ~Close~ بيقطع الباقي و ~Serve~ ترجع الخطأ.
~~~go
shutdownCtx, cancel := context.WithTimeout(context.Background(), drain)
defer cancel()
err := srv.Shutdown(shutdownCtx)
~~~

## أخطاء شائعة
~~~go
log.Printf("%s %s %v", r.Method, r.URL, r.Header) // الكوكي والـ query في السجل
key := r.Header.Get("X-Forwarded-For") // العميل يختار مفتاح الحد بنفسه
counts[status]++ // من غير قفل: data race
srv.Shutdown(ctx) // ctx اللي اتلغى أصلاً: مفيش تصفية خالص
~~~

## المهمة
في ~ops.go~ اكتب ~Operate(h, o)~ و ~Serve~. ~Ops~ فيها ~Now~ و ~Limit~ و ~Window~ و ~Log~ و ~NewID~؛ ضيف أي حقول خاصة محتاجها، و ~Operate~ هي اللي تجهزها.

كل رد فيه ~X-Request-ID~: القيمة الجاية لو بتطابق ~^[A-Za-z0-9-]{1,64}$~، وإلا ~o.NewID()~. ~GET /metrics~ (مش بيتعد ولا بيتحد) يرجع 200 نص فيه ~command_requests_total{code="N"} count~ لكل حالة اترجعت، و ~command_requests_in_flight N~ للطلبات اللي جوه ~h~ دلوقتي. كل رد حالته 400 أو أكتر جسمه ~{"error":"...","request_id":"..."}~ بـ ~application/json~: كود الخطأ من جسم ~h~ لو فيه ~error~، وإلا اسم الحالة صغير بـ ~_~ (زي ~not_found~). الردود الناجحة بالبايت زي ما ~h~ كتبها. الـ panic: 500 و ~internal~. أكتر من ~Limit~ طلب من نفس IP في ~Window~ (تبدأ من أول طلب): 429 و ~rate_limited~ و ~Retry-After~ بالثواني الباقية لفوق، من غير ما ~h~ تتنادي. كل طلب (غير ~/metrics~) يكتب سطر JSON واحد في ~o.Log~ بالمفاتيح ~request_id, method, path, status~ بس. كل ده آمن مع طلبات متوازية (~-race~).

~Serve(ctx, ln, h, drain)~ تخدم على ~ln~ لحد ما ~ctx~ يتلغي، وبعدها تبطّل تقبل اتصالات وتصفّي الجاري: ~nil~ لو خلص في ~drain~، وخطأ لو لأ (بعد ما تقطع الباقي). النجمة التالتة: ~command_requests_in_flight~ صح وطلب شغال.
`,
		Hints: []string{"الـmiddleware يضيف معلومات الطلب حول handler من غير نقل ملكية الرد.", "قسّم عمليات HTTP عن الحالة المحمية بالتزامن وإدارة عمر السيرفر.", "ابدأ من تواقيع Operate وServe في starter وحدد موارد كل واحدة."},
		En: &LevelText{Title: "Operations Room", Goal: "Wrap the API with request IDs, a fixed error shape, metrics, rate limiting, secret-free logs, and shutdown that drains in-flight requests", Files: []LabFile{{Name: "main.go", Starter: ch371Starter}, {Name: "auth.go", Starter: ch372Starter}, {Name: "crews.go", Starter: ch373Starter}, {Name: "tasks.go", Starter: ch374Starter}, {Name: "ops.go", Starter: ch375Starter}},
			Lesson: `# Command Center: Operations Room

## Before you write
- A user writes "I got an error". How do you find their exact request in the log?
- The server is shutting down for an update while a request is halfway done. Does it finish or get cut off?

## What this builds on
Every file from the earlier stages carries over as is, and the new file is ~ops.go~. Middleware and the error envelope come from ~26-4-envelope~, rate limiting from ~27-4-rate-limit~, graceful shutdown from ~28-4-graceful~, ~recover~ from the Functions and Errors chapter, and ~sync.Mutex~ from The Race (12). Metrics and structured logs are the same idea as the Operations and Monitoring chapter (36), and load protection shares the spirit of the Brokers chapter (33). Only the integration decisions are new.

## A request ID on every reply
Each request gets an ID: if the client sent a safe ~X-Request-ID~ (letters, digits, and dashes, up to 64) we use it, otherwise we generate one from the injected ~Ops.NewID~. The ID goes back in the reply header, in every error body, and in the log line. A user's message then leads you to one line.
~~~go
var safeID = regexp.MustCompile("^[a-z0-9]{1,8}$")
if !safeID.MatchString(id) { id = newID() }
~~~
Never accept any client ID: text with a newline could forge a log line.

## One error shape
Every earlier stage returns ~{"error":"..."}~. The middleware adds ~request_id~ to any reply with status 400 or above. To do that it wraps the ~ResponseWriter~: successful replies pass straight through, error replies are collected in a buffer and rewritten.
~~~go
type capture struct {
 http.ResponseWriter
 status int
}
func (c *capture) WriteHeader(s int) { c.status = s; c.ResponseWriter.WriteHeader(s) }
~~~
A panic in any handler becomes 500 and ~{"error":"internal"}~ through ~defer~ and ~recover~, and the server keeps running.

## Metrics as text
~GET /metrics~ returns hand-written Prometheus text format: one line per status and counter. No library is needed; the format is simple, and the counters are guarded by ~sync.Mutex~ because requests run concurrently.
~~~text
command_requests_total{code="200"} 12
command_requests_in_flight 1
~~~

## Rate limiting with an injected clock
A fixed window per IP address (from ~net.SplitHostPort(r.RemoteAddr)~, not from a header the client could forge). The window starts with the first request and lasts ~Window~. Extra requests get 429 with ~Retry-After~ in remaining seconds, rounded up.

## A log without secrets
One JSON line per request with only ~request_id~, ~method~, ~path~, and ~status~. ~r.URL.Path~ without the query, and no headers: cookies, tokens, and passwords never reach the log.

## Shutdown with a drain window
~Serve(ctx, ln, h, drain)~: when ~ctx~ is canceled, ~Shutdown~ closes the listener at once and waits for in-flight requests up to ~drain~. If they finish, ~Serve~ returns ~nil~. If the window runs out, ~Close~ cuts the rest and ~Serve~ returns the error.
~~~go
shutdownCtx, cancel := context.WithTimeout(context.Background(), drain)
defer cancel()
err := srv.Shutdown(shutdownCtx)
~~~

## Common mistakes
~~~go
log.Printf("%s %s %v", r.Method, r.URL, r.Header) // cookie and query in the log
key := r.Header.Get("X-Forwarded-For") // the client picks its own limit key
counts[status]++ // no lock: data race
srv.Shutdown(ctx) // an already-canceled ctx: no drain at all
~~~

## Your mission
In ~ops.go~ write ~Operate(h, o)~ and ~Serve~. ~Ops~ holds ~Now~, ~Limit~, ~Window~, ~Log~, and ~NewID~; add any private fields you need, and let ~Operate~ initialize them.

Every reply carries ~X-Request-ID~: the incoming value if it matches ~^[A-Za-z0-9-]{1,64}$~, otherwise ~o.NewID()~. ~GET /metrics~ (neither counted nor limited) returns 200 text with ~command_requests_total{code="N"} count~ for every status returned, and ~command_requests_in_flight N~ for requests currently inside ~h~. Every reply with status 400 or above has the body ~{"error":"...","request_id":"..."}~ as ~application/json~: the error code from ~h~'s body if it has ~error~, otherwise the lowercase status text with ~_~ (like ~not_found~). Successful replies are byte-identical to what ~h~ wrote. A panic: 500 and ~internal~. More than ~Limit~ requests from one IP within ~Window~ (starting at the first request): 429, ~rate_limited~, and ~Retry-After~ in remaining seconds rounded up, without calling ~h~. Every request (except ~/metrics~) writes exactly one JSON line to ~o.Log~ with only the keys ~request_id, method, path, status~. All of this is safe under concurrent requests (~-race~).

~Serve(ctx, ln, h, drain)~ serves on ~ln~ until ~ctx~ is canceled, then stops accepting connections and drains in-flight requests: ~nil~ if they finish within ~drain~, an error otherwise (after cutting the rest). Third star: ~command_requests_in_flight~ is correct while a request is running.
`,
			LessonRevision: 1, Steps: ch37OpsEnglishSteps,
			Hints: []string{"Middleware adds request context around a handler without taking ownership of its reply.", "Separate HTTP operations from synchronized state and server lifetime management.", "Start with Operate and Serve in the starter and map each function's resources."}},
		Lab: &Lab{ReqFile: "ops.go", Bonus: true, Hidden: ch375Hidden},
	})
}

const ch375Starter = `package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

type Ops struct {
	Now    func() time.Time
	Limit  int
	Window time.Duration
	Log    io.Writer
	NewID  func() string
}

func Operate(h http.Handler, o *Ops) http.Handler { return h }

func Serve(ctx context.Context, ln net.Listener, h http.Handler, drain time.Duration) error {
	return errors.New("TODO")
}
`

const ch375Hidden = `package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

type hiddenOpsClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *hiddenOpsClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *hiddenOpsClock) Add(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

type hiddenLog struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (l *hiddenLog) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.Write(p)
}

func (l *hiddenLog) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.buf.String()
}

func hiddenOps(limit int, seed int64) (*Ops, *hiddenOpsClock, *hiddenLog) {
	clock := &hiddenOpsClock{now: time.Unix(1_900_000_000, 0)}
	log := &hiddenLog{}
	var mu sync.Mutex
	rng := rand.New(rand.NewSource(seed))
	return &Ops{Now: clock.Now, Limit: limit, Window: 10 * time.Second, Log: log, NewID: func() string {
		mu.Lock()
		defer mu.Unlock()
		return fmt.Sprintf("gen-%d", rng.Int63())
	}}, clock, log
}

func hiddenInner() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ok", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"status\":\"green\",\"n\":1}\n")
	})
	mux.HandleFunc("GET /gone", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
	})
	mux.HandleFunc("GET /teapot", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		io.WriteString(w, "plain words")
	})
	mux.HandleFunc("GET /boom", func(w http.ResponseWriter, r *http.Request) { panic("reactor fault") })
	return mux
}

func hiddenCall(t *testing.T, c *http.Client, target string, header map[string]string) (int, http.Header, []byte) {
	t.Helper()
	req, _ := http.NewRequest("GET", target, nil)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, resp.Header, raw
}

func TestHiddenRequestIDsAndErrorShape(t *testing.T) {
	o, _, log := hiddenOps(1000, 3751)
	srv := httptest.NewServer(Operate(hiddenInner(), o))
	defer srv.Close()
	c := srv.Client()
	code, h, raw := hiddenCall(t, c, srv.URL+"/ok", map[string]string{"X-Request-ID": "trace-Abc-42"})
	if code != 200 || h.Get("X-Request-ID") != "trace-Abc-42" || string(raw) != "{\"status\":\"green\",\"n\":1}\n" {
		t.Fatal("success passthrough", code, h, string(raw))
	}
	ids := map[string]bool{}
	for _, bad := range []string{"", "has space", "line\\nbreak", strings.Repeat("a", 65), "semi;colon"} {
		_, h, _ := hiddenCall(t, c, srv.URL+"/ok", map[string]string{"X-Request-ID": bad})
		id := h.Get("X-Request-ID")
		if !strings.HasPrefix(id, "gen-") || ids[id] {
			t.Fatal("generated id", bad, id)
		}
		ids[id] = true
	}
	for _, cs := range []struct {
		path, code string
		status     int
	}{{"/gone", "not_found", 404}, {"/teapot", "i'm_a_teapot", 418}, {"/boom", "internal", 500}, {"/nowhere", "not_found", 404}} {
		code, h, raw := hiddenCall(t, c, srv.URL+cs.path, nil)
		var body map[string]string
		if json.Unmarshal(raw, &body) != nil || code != cs.status || body["error"] != cs.code || body["request_id"] != h.Get("X-Request-ID") || body["request_id"] == "" || !strings.HasPrefix(h.Get("Content-Type"), "application/json") || len(body) != 2 {
			t.Fatal("error shape", cs, code, string(raw), h)
		}
	}
	code, _, _ = hiddenCall(t, c, srv.URL+"/ok", nil)
	if code != 200 {
		t.Fatal("server died after a panic", code)
	}
	if lines := strings.Count(log.String(), "\n"); lines != 11 {
		t.Fatal("one log line per request, got", lines)
	}
}

func TestHiddenRateLimitAndMetrics(t *testing.T) {
	o, clock, _ := hiddenOps(3, 3752)
	srv := httptest.NewServer(Operate(hiddenInner(), o))
	defer srv.Close()
	c := srv.Client()
	for i := 0; i < 3; i++ {
		if code, _, _ := hiddenCall(t, c, srv.URL+"/ok", nil); code != 200 {
			t.Fatal("under limit", i, code)
		}
	}
	code, h, raw := hiddenCall(t, c, srv.URL+"/ok", nil)
	if code != 429 || h.Get("Retry-After") != "10" || !bytes.Contains(raw, []byte("\"rate_limited\"")) || !bytes.Contains(raw, []byte(h.Get("X-Request-ID"))) {
		t.Fatal("limit", code, h, string(raw))
	}
	clock.Add(3500 * time.Millisecond)
	if code, h, _ := hiddenCall(t, c, srv.URL+"/boom", nil); code != 429 || h.Get("Retry-After") != "7" {
		t.Fatal("limited request must not reach the handler", code, h.Get("Retry-After"))
	}
	code, _, raw = hiddenCall(t, c, srv.URL+"/metrics", nil)
	text := string(raw)
	if code != 200 || !strings.Contains(text, "command_requests_total{code=\"200\"} 3\n") || !strings.Contains(text, "command_requests_total{code=\"429\"} 2\n") || strings.Contains(text, "code=\"500\"") || !strings.Contains(text, "command_requests_in_flight 0\n") {
		t.Fatal("metrics", code, text)
	}
	clock.Add(6500 * time.Millisecond)
	if code, _, _ := hiddenCall(t, c, srv.URL+"/ok", nil); code != 200 {
		t.Fatal("new window", code)
	}

	o2, _, _ := hiddenOps(5, 3753)
	srv2 := httptest.NewServer(Operate(hiddenInner(), o2))
	defer srv2.Close()
	var wg sync.WaitGroup
	var mu sync.Mutex
	ok := 0
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := srv2.Client().Get(srv2.URL + "/ok")
			if err != nil {
				return
			}
			resp.Body.Close()
			if resp.StatusCode == 200 {
				mu.Lock()
				ok++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if ok != 5 {
		t.Fatal("concurrent limit let through", ok)
	}
	_, _, raw = hiddenCall(t, srv2.Client(), srv2.URL+"/metrics", nil)
	if !strings.Contains(string(raw), "command_requests_total{code=\"200\"} 5\n") || !strings.Contains(string(raw), "command_requests_total{code=\"429\"} 19\n") {
		t.Fatal("concurrent metrics", string(raw))
	}
}

func TestHiddenLogRedaction(t *testing.T) {
	o, _, log := hiddenOps(100, 3754)
	srv := httptest.NewServer(Operate(hiddenInner(), o))
	defer srv.Close()
	secret := fmt.Sprintf("s3cr3t%d", rand.New(rand.NewSource(3754)).Int63())
	hiddenCall(t, srv.Client(), srv.URL+"/gone?token="+secret+"&password="+secret, map[string]string{"Cookie": "session=" + secret, "Authorization": "Bearer " + secret, "X-Request-ID": "probe-1"})
	text := log.String()
	if strings.Contains(text, secret) {
		t.Fatal("secret leaked into the log:", text)
	}
	sc := bufio.NewScanner(strings.NewReader(text))
	for sc.Scan() {
		var line map[string]any
		if err := json.Unmarshal(sc.Bytes(), &line); err != nil || len(line) != 4 || line["request_id"] != "probe-1" || line["method"] != "GET" || line["path"] != "/gone" || line["status"] != float64(404) {
			t.Fatal("log line", sc.Text(), err)
		}
	}
}

func hiddenServe(t *testing.T, h http.Handler, drain time.Duration) (string, context.CancelFunc, chan error) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, ln, h, drain) }()
	return "http://" + ln.Addr().String(), cancel, done
}

func TestHiddenGracefulDrain(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	slow := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
		io.WriteString(w, "landed")
	})
	base, cancel, done := hiddenServe(t, slow, 3*time.Second)
	got := make(chan string, 1)
	go func() {
		resp, err := http.Get(base + "/dock")
		if err != nil {
			got <- "error: " + err.Error()
			return
		}
		defer resp.Body.Close()
		raw, _ := io.ReadAll(resp.Body)
		got <- string(raw)
	}()
	<-entered
	cancel()
	deadline := time.Now().Add(2 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", strings.TrimPrefix(base, "http://"), 100*time.Millisecond)
		if err != nil {
			break
		}
		conn.Close()
		if time.Now().After(deadline) {
			t.Fatal("listener still accepting after shutdown began")
		}
		time.Sleep(10 * time.Millisecond)
	}
	select {
	case err := <-done:
		t.Fatal("Serve returned before draining", err)
	default:
	}
	close(release)
	if body := <-got; body != "landed" {
		t.Fatal("in-flight request was cut:", body)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatal("clean drain must return nil", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Serve did not return")
	}
}

func TestHiddenDrainWindowExpires(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	defer close(release)
	stuck := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
	})
	base, cancel, done := hiddenServe(t, stuck, 150*time.Millisecond)
	go func() {
		if resp, err := http.Get(base + "/stuck"); err == nil {
			resp.Body.Close()
		}
	}()
	<-entered
	start := time.Now()
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expired drain must return an error")
		}
		if time.Since(start) < 100*time.Millisecond {
			t.Fatal("drain window ignored")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Serve waited past the drain window")
	}
}

func TestHiddenBonusInFlightGauge(t *testing.T) {
	o, _, _ := hiddenOps(100, 3755)
	entered, release := make(chan struct{}), make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hold", func(w http.ResponseWriter, r *http.Request) {
		entered <- struct{}{}
		<-release
	})
	srv := httptest.NewServer(Operate(mux, o))
	defer srv.Close()
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if resp, err := srv.Client().Get(srv.URL + "/hold"); err == nil {
				resp.Body.Close()
			}
		}()
		<-entered
	}
	_, _, raw := hiddenCall(t, srv.Client(), srv.URL+"/metrics", nil)
	close(release)
	wg.Wait()
	if !strings.Contains(string(raw), "command_requests_in_flight 2\n") {
		t.Fatal("in-flight gauge", string(raw))
	}
	_, _, raw = hiddenCall(t, srv.Client(), srv.URL+"/metrics", nil)
	if !strings.Contains(string(raw), "command_requests_in_flight 0\n") {
		t.Fatal("in-flight gauge after", string(raw))
	}
}
`
