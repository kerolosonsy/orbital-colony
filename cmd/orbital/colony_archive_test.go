package main

import (
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"orbital/colony/econ"
	"orbital/colony/labkit"
)

func TestSystemOwnershipEnvironmentAndArchivePersistence(t *testing.T) {
	labkit.LeakCheck(t)
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	for attempt, archive := range []bool{false, true, true} {
		if archive {
			c.state.Levels["sys.archive"] = 1
		}
		cmd := exec.Command("/not-started")
		cmd.Env = []string{"OTHER=kept", "ORBITAL_COLONY_DIR=inherited", "ORBITAL_MARKET_URL=http://inherited"}
		s, err := c.startArchiveMarket(cmd, writer, "en")
		if err != nil {
			t.Fatal(err)
		}
		env := strings.Join(cmd.Env, "\n")
		if strings.Contains(env, "inherited") || !strings.Contains(env, "OTHER=kept") || s.market != nil || strings.Contains(env, "ORBITAL_MARKET_URL=") {
			t.Fatal(env)
		}
		if archive {
			if !strings.Contains(env, "ORBITAL_COLONY_DIR="+s.dir) || !strings.HasPrefix(s.dir, filepath.Dir(c.path)) {
				t.Fatal(env)
			}
			marker := filepath.Join(s.dir, "player-plan.txt")
			if attempt == 2 {
				body, err := os.ReadFile(marker)
				if err != nil || string(body) != "keep across runs" {
					t.Fatal("archive did not persist", err)
				}
			}
			if err := os.WriteFile(marker, []byte("keep across runs"), 0o600); err != nil {
				t.Fatal(err)
			}
		} else if s.dir != "" || strings.Contains(env, "ORBITAL_COLONY_DIR=") {
			t.Fatal("archive started without ownership")
		}
		s.closeUnstarted(cmd)
	}
	if err := c.persist(); err != nil {
		t.Fatal(err)
	}
	if c.state.Has("sys.market") {
		t.Fatal("archive unexpectedly granted market")
	}
	reloaded, err := newColonyGame(filepath.Dir(c.path))
	if err != nil || reloaded.state.Seed != c.state.Seed || !reloaded.state.Has("sys.archive") {
		t.Fatal("old save no longer loads", err)
	}
}

func TestArchiveJournalRejectsPartialAndAcceptsCompleteVersions(t *testing.T) {
	for _, version := range []string{"old", "new", "partial", "checksum-only", "early-exit"} {
		t.Run(version, func(t *testing.T) {
			c, err := newColonyGame(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			c.state.Levels["sys.archive"] = 1
			if err := c.state.Accept("d22-atomic-journal"); err != nil {
				t.Fatal(err)
			}
			s := &archiveMarketRun{dir: filepath.Join(filepath.Dir(c.path), "colony-data"), seed: c.state.Seed, journal: true, killed: version != "early-exit", lang: "en"}
			if err := s.prepareArchive(); err != nil {
				t.Fatal(err)
			}
			plan := journalContents(s.seed, version)
			if version == "partial" {
				plan = plan[:len(plan)/2]
			}
			if version == "checksum-only" {
				plan = []byte("sha256:000000\n")
			}
			if err := os.WriteFile(filepath.Join(s.dir, "plan.txt"), plan, 0o600); err != nil {
				t.Fatal(err)
			}
			line := c.finishArchiveMarket(&colonyRun{am: s})
			var event struct{ OK bool }
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				t.Fatal(err)
			}
			want := version == "old" || version == "new"
			if event.OK != want || c.state.Won("d22-atomic-journal") != want || c.state.Active != "" {
				t.Fatal("incorrect journal settlement", line)
			}
		})
	}
}

func TestSystemRunStopClosesSidecarAndPersistsArchive(t *testing.T) {
	labkit.LeakCheck(t)
	requireMarketSocket(t)
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.state.Levels["sys.archive"], c.state.Levels["sys.market"] = 1, 1
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	code := `package main
import "orbital/colony"
func main(){ for { colony.Wait() } }
`
	prog := buildColonyProgram(t, testRunner(t), code)
	startFastSystemProgram(t, c, prog)
	waitSystemSnapshot(t, ch)
	c.mu.Lock()
	active := c.active
	s := active.am
	c.mu.Unlock()
	c.stop()
	select {
	case <-active.done:
	default:
		t.Fatal("stop returned before run cleanup")
	}
	select {
	case <-s.market.done:
	default:
		t.Fatal("HTTP serve goroutine survived stop")
	}
	conn, err := net.DialTimeout("tcp", strings.TrimPrefix(s.market.url, "http://"), 100*time.Millisecond)
	if err == nil {
		conn.Close()
		t.Fatal("listener survived stopped run")
	}
	if _, err := os.Stat(filepath.Join(s.dir, "ledger.csv")); err != nil {
		t.Fatal("persistent archive removed with build directory", err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active != nil {
		t.Fatal("stopped run is still active")
	}
	if on := c.archiveMarketSummary()["market"].(map[string]any)["on"]; on != false {
		t.Fatal("summary says market is running")
	}
}

func TestArchiveAtomicHarnessRejectsPartialPlayerJournal(t *testing.T) {
	labkit.LeakCheck(t)
	c := econ.FindContract("d22-atomic-journal")
	project := colonyProject{Files: map[string]string{"main.go": `package main
import("os";"path/filepath";"orbital/colony")
func main(){ _=os.Rename; _=os.WriteFile(filepath.Join(os.Getenv("ORBITAL_COLONY_DIR"),"plan.txt"),[]byte("half a journal"),0600); for {colony.Wait()} }
`}}
	if got := playArchiveMarketContract(t, testRunner(t), c, project); got.OK {
		t.Fatal("kill harness accepted partial player journal")
	}
}

func TestSystemWithoutNodesDoesNotInjectLiveEvent(t *testing.T) {
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	s := &archiveMarketRun{game: c}
	s.publish()
	select {
	case line := <-ch:
		t.Fatalf("system event without ownership displaced drone events: %s", line)
	default:
	}
}

func TestArchiveRunStopRetainsFilesAndReapsReaders(t *testing.T) {
	labkit.LeakCheck(t)
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.state.Levels["sys.archive"] = 1
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	prog := buildColonyProgram(t, testRunner(t), `package main
import "orbital/colony"
func main(){for{colony.Wait()}}
`)
	startFastSystemProgram(t, c, prog)
	c.mu.Lock()
	active := c.active
	c.mu.Unlock()
	waitSystemSnapshot(t, ch)
	c.stop()
	select {
	case <-active.done:
	default:
		t.Fatal("readers survived stop")
	}
	for _, name := range []string{"ledger.csv", "events.jsonl"} {
		if _, err := os.Stat(filepath.Join(active.am.dir, name)); err != nil {
			t.Fatal("archive removed on stop", err)
		}
	}
	if _, err := os.Stat(filepath.Dir(prog)); !os.IsNotExist(err) {
		t.Fatal("build directory survived stopped run")
	}
}

func waitSystemSnapshot(t *testing.T, ch <-chan string) {
	t.Helper()
	deadline := time.NewTimer(3 * time.Second)
	defer deadline.Stop()
	for {
		select {
		case line := <-ch:
			var event struct{ Op string }
			if json.Unmarshal([]byte(line), &event) == nil && event.Op == "econ" {
				return
			}
		case <-deadline.C:
			t.Fatal("system run did not initialize")
		}
	}
}

func TestArchiveRelativeProfilePathUsesAbsolutePlayerDirectory(t *testing.T) {
	root := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	rel, err := filepath.Rel(cwd, root)
	if err != nil {
		t.Fatal(err)
	}
	s := &archiveMarketRun{dir: filepath.Join(rel, "colony-data")}
	if err := s.prepareArchive(); err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(s.dir) || s.dir != filepath.Join(root, "colony-data") {
		t.Fatal("player would resolve archive relative to its build directory", s.dir)
	}
}
