# Orbital Colony — كل الأوامر في مكان واحد. اكتب: make  (أو make help)
#
# الإعدادات (تقدر تغيّرها: make start PORT=8080)
PORT    ?= 7777
ADDR    ?= 127.0.0.1:$(PORT)
SAVES   ?= saves
# تسجيل الدخول الاختياري: make start PUBLIC=http://127.0.0.1:7777
PUBLIC  ?=
# مفاتيح الدخول بـGoogle على جهازك: حطها في .orbital.env (مش بيترفع على git)
#   GOOGLE_CLIENT_ID=...
#   GOOGLE_CLIENT_SECRET=...
-include .orbital.env
export GOOGLE_CLIENT_ID GOOGLE_CLIENT_SECRET
BIN     := ./orbital
PIDFILE := .orbital.pid
LOGFILE := .orbital.log
BACKUPS := backups
OPEN    ?= open
STAMP    = $(shell date +%Y%m%d-%H%M%S)

.DEFAULT_GOAL := help
.PHONY: help run start stop restart status logs open build test test-quick test-web demo check fmt vet lint stats \
        backup restore backups players reset-all clean doctor

help: ## اعرض كل الأوامر
	@echo "Orbital Colony"
	@echo
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  make %-11s %s\n", $$1, $$2}'
	@echo
	@echo "  PORT=$(PORT)   مثال: make start PORT=8080"

# ---------- التشغيل ----------

run: build ## شغّل اللعبة هنا في الـ terminal (Ctrl+C يقفلها)
	@echo "→ http://$(ADDR)"
	$(BIN) -addr $(ADDR) -saves $(SAVES) -public "$(PUBLIC)"

start: build ## شغّل اللعبة في الخلفية وافتح المتصفح
	@if [ -f $(PIDFILE) ] && kill -0 $$(cat $(PIDFILE)) 2>/dev/null; then \
		echo "اللعبة شغالة أصلاً (PID $$(cat $(PIDFILE))). استخدم: make restart"; exit 0; fi; \
	if lsof -iTCP:$(PORT) -sTCP:LISTEN >/dev/null 2>&1; then \
		echo "البورت $(PORT) مشغول ببرنامج تاني. جرّب: make start PORT=7778"; exit 1; fi; \
	nohup $(BIN) -addr $(ADDR) -saves $(SAVES) -public "$(PUBLIC)" > $(LOGFILE) 2>&1 & echo $$! > $(PIDFILE); \
	for i in 1 2 3 4 5 6 7 8 9 10; do \
		lsof -iTCP:$(PORT) -sTCP:LISTEN >/dev/null 2>&1 && break; \
		kill -0 $$(cat $(PIDFILE)) 2>/dev/null || break; sleep 0.5; done; \
	if lsof -iTCP:$(PORT) -sTCP:LISTEN >/dev/null 2>&1; then \
		echo "✓ اللعبة شغالة: http://$(ADDR)  (PID $$(cat $(PIDFILE)))"; $(OPEN) http://$(ADDR); \
	else echo "✗ اللعبة ماشتغلتش:"; cat $(LOGFILE); rm -f $(PIDFILE); exit 1; fi

stop: ## اقفل اللعبة اللي شغالة في الخلفية
	@if [ -f $(PIDFILE) ] && kill -0 $$(cat $(PIDFILE)) 2>/dev/null; then \
		kill $$(cat $(PIDFILE)) && echo "✓ اتقفلت"; else echo "مفيش لعبة شغالة من make start"; fi; \
	rm -f $(PIDFILE)

restart: ## ابني من جديد وشغّل (استخدمها بعد أي تحديث)
	@$(MAKE) --no-print-directory stop
	@$(MAKE) --no-print-directory start

status: ## هل اللعبة شغالة؟
	@if [ -f $(PIDFILE) ] && kill -0 $$(cat $(PIDFILE)) 2>/dev/null; then \
		echo "✓ شغالة: http://$(ADDR)  (PID $$(cat $(PIDFILE)))"; else echo "✗ مش شغالة. شغّلها: make start"; fi

logs: ## تابع رسايل السيرفر (Ctrl+C للخروج)
	@touch $(LOGFILE); tail -f $(LOGFILE)

open: ## افتح اللعبة في المتصفح
	$(OPEN) http://$(ADDR)

build: ## ابني السيرفر
	go build -o $(BIN) ./cmd/orbital

# ---------- الاختبارات ----------

test: ## كل الاختبارات: كل مهمة لازم تتحل (حوالي 8 دقايق)
	cd colony && go test ./...
	go test ./cmd/orbital/ -count=1 -timeout 30m

test-quick: ## اختبارات سريعة (أقل من دقيقة)
	cd colony && go test ./...
	go test ./cmd/orbital/ -count=1 -run '^(TestLevelOrder|TestPreviewHasNoTimers|TestGoalMustHold|TestCompileErrorTips|Test.*Profile.*|Test.*Lang.*|TestListenConfig.*|TestStatic.*|TestLevels.*)$$/^(1-1-|1-2-|3-3-)'

test-web: ## اختبارات واجهة المتصفح (محتاج Node 20+)
	node --test web/*.test.js web/sandbox/*.test.js cmd/lesson-demo/site/*.test.js

demo: ## Build the static First Signal lesson in dist/
	go run ./cmd/lesson-demo

check: fmt vet test-quick test-web ## نسّق وافحص وجرّب بسرعة قبل ما تحفظ شغلك

fmt: ## نسّق كود Go
	gofmt -w cmd colony assets.go

vet: ## افحص الكود بـ go vet
	go vet ./...
	cd colony && go vet ./...

lint: ## افحص من غير تعديل: gofmt -l و go vet و staticcheck (لو متسطب)
	@test -z "$$(gofmt -l cmd colony assets.go)" || { echo "محتاج gofmt:"; gofmt -l cmd colony assets.go; exit 1; }
	@$(MAKE) --no-print-directory vet
	@if command -v staticcheck >/dev/null; then staticcheck ./cmd/... && cd colony && staticcheck ./...; \
	else echo "• staticcheck مش متسطب: go install honnef.co/go/tools/cmd/staticcheck@latest"; fi

stats: ## اطبع أكتر 20 ورقة اتعثر فيها اللاعبون
	go run ./cmd/orbital -stats -saves $(SAVES)

# ---------- الحفظ واللاعبين ----------

players: ## اعرض اللاعبين المتسجلين
	@found=; for f in $(SAVES)/profiles/*/profile.json $(SAVES)/users/*/profiles/*/profile.json; do \
		[ -f "$$f" ] || continue; found=1; \
		name=$$(sed -n 's/.*"name": *"\([^"]*\)".*/\1/p' "$$f"); \
		lang=$$(sed -n 's/.*"lang": *"\([^"]*\)".*/\1/p' "$$f"); \
		echo "  $$(basename $$(dirname $$f))  →  $$name  ($$lang)"; done; \
	[ -n "$$found" ] || echo "لسه مفيش لاعبين. شغّل اللعبة الأول."

backup: ## خد نسخة احتياطية من كل الحفظ والـ workspace
	@mkdir -p $(BACKUPS)
	@tar -czf $(BACKUPS)/orbital-$(STAMP).tar.gz $$(ls -d saves workspace 2>/dev/null) && \
		echo "✓ $(BACKUPS)/orbital-$(STAMP).tar.gz"

backups: ## اعرض النسخ الاحتياطية
	@ls -1t $(BACKUPS)/*.tar.gz 2>/dev/null || echo "مفيش نسخ لسه. اعمل: make backup"

restore: ## رجّع نسخة: make restore FILE=backups/orbital-....tar.gz
	@[ -n "$(FILE)" ] || { echo "حدد النسخة: make restore FILE=backups/orbital-....tar.gz"; $(MAKE) --no-print-directory backups; exit 1; }
	@[ -f "$(FILE)" ] || { echo "مش لاقي $(FILE)"; exit 1; }
	@$(MAKE) --no-print-directory stop
	@$(MAKE) --no-print-directory backup
	@rm -rf saves workspace && tar -xzf "$(FILE)" && echo "✓ رجعت $(FILE). شغّل: make start"

reset-all: ## امسح كل اللاعبين وابدأ من الصفر (بياخد نسخة الأول): make reset-all CONFIRM=yes
	@[ "$(CONFIRM)" = "yes" ] || { echo "ده هيمسح كل اللاعبين وكل التقدّم."; echo "لو متأكد: make reset-all CONFIRM=yes  (هتتاخد نسخة احتياطية الأول)"; exit 1; }
	@$(MAKE) --no-print-directory stop
	@$(MAKE) --no-print-directory backup
	@rm -rf saves workspace && echo "✓ اتمسح كل حاجة. افتح المتصفح وامسح بيانات الموقع كمان لو عاوز."

# ---------- متفرقات ----------

doctor: ## اتأكد إن كل حاجة جاهزة على الجهاز
	@command -v go >/dev/null && echo "✓ $$(go version)" || echo "✗ Go مش متسطب: https://go.dev/dl"
	@command -v code >/dev/null && echo "✓ VS Code (code)" || echo "• VS Code: افتحه ومن Command Palette اختار: Shell Command: Install 'code' command in PATH"
	@command -v gopls >/dev/null && echo "✓ gopls" || echo "• gopls مش متسطب (للـ autocomplete): go install golang.org/x/tools/gopls@latest"
	@lsof -iTCP:$(PORT) -sTCP:LISTEN >/dev/null 2>&1 && echo "• البورت $(PORT) مستخدم دلوقتي" || echo "✓ البورت $(PORT) فاضي"
	@(cd labdeps && GOPROXY=off GOFLAGS=-mod=mod go build ./... >/dev/null 2>&1) && echo "✓ مكتبات المختبرات جاهزة" || echo "• مكتبات المختبرات مش جاهزة: شغّل make deps مرة وإنت متوصل بالإنترنت"
	@command -v docker >/dev/null && echo "✓ Docker (لمهمات الخدمات: make services-up)" || echo "• Docker مش متسطب: محتاجه بس لمهمات Postgres/MongoDB/Kafka"
	@command -v node >/dev/null && echo "✓ Node $$(node --version) (لاختبارات الواجهة: make test-web)" || echo "• Node مش متسطب: محتاجه بس لـ make test-web"
	@command -v codex >/dev/null && echo "✓ Codex (للمساعد في اللعبة؛ لازم يكون متسجّل: codex login)" || echo "• Codex مش متسطب (للمساعد في اللعبة): brew install codex وبعدين codex login"

clean: ## امسح ملفات البناء والـ logs (الحفظ مش بيتمسح)
	rm -f $(BIN) $(LOGFILE)

.PHONY: deps services-up services-down services-status
deps: ## جهّز مكتبات المختبرات / Prepare lab dependencies (requires network)
	cd labdeps && GOPROXY=https://proxy.golang.org go mod download all
	cd labdeps && GOPROXY=off go build ./...

services-up: ## شغّل الخدمات المحلية واستنى لحد ما تبقى جاهزة / Start local services and wait until healthy
	docker compose -f services/compose.yaml up -d --wait

services-down: ## اقفل الخدمات المحلية / Stop local services
	docker compose -f services/compose.yaml down

services-status: ## حالة الخدمات / Service status
	docker compose -f services/compose.yaml ps
