//go:build unix

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"orbital/colony/econ"
)

func TestProcessColonyShutdownSavesOrLosesCargo(t *testing.T) {
	r := testRunner(t)
	for _, tc := range []struct {
		name, code string
		healthy    bool
		grace      bool
	}{
		{"no registration", `package main
import "orbital/colony"
func main(){colony.Report("ready");for{colony.Wait()}}`, false, false},
		{"returns and unloads", `package main
import "orbital/colony"
func main(){stop:=colony.Shutdown();colony.Move(colony.East);colony.Report("ready");<-stop;for !colony.Move(colony.West){colony.Repair()};colony.Unload()}`, true, true},
		{"ignores notification", `package main
import "orbital/colony"
func main(){colony.Shutdown();colony.Report("ready");for{colony.Wait()}}`, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			prog := buildColonyProgram(t, r, tc.code)
			st := econ.New(11)
			st.Contracts["c13-reed-clock"] = &econ.ContractRun{Best: 1}
			if err := st.Accept(econ.CleanStopContract); err != nil {
				t.Fatal(err)
			}
			st.Fleet[0].Cargo["moss"] = 3
			blob, _ := json.Marshal(st)
			save := filepath.Join(t.TempDir(), "colony.json")
			if err := os.WriteFile(save, blob, 0600); err != nil {
				t.Fatal(err)
			}
			reader, writer, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, prog)
			cmd.Env = append(os.Environ(), "ORBITAL_LEVEL=colony", "ORBITAL_COLONY="+save, "ORBITAL_EVENTS_FD=3", "ORBITAL_COLONY_UNIT=2ms")
			cmd.ExtraFiles = []*os.File{writer}
			if err := cmd.Start(); err != nil {
				writer.Close()
				t.Fatal(err)
			}
			writer.Close()
			defer func() { _ = cmd.Process.Kill() }()
			scanner := bufio.NewScanner(reader)
			scanner.Buffer(make([]byte, 4096), 8<<20)
			sent := time.Time{}
			var final *econ.State
			for scanner.Scan() {
				var ev struct {
					Op    string
					State *econ.State
				}
				if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
					t.Fatal(err)
				}
				if ev.Op == "report" && sent.IsZero() {
					sent = time.Now()
					if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
						t.Fatal(err)
					}
				}
				if ev.State != nil {
					final = ev.State
				}
			}
			if err := scanner.Err(); err != nil {
				t.Fatal(err)
			}
			if err := cmd.Wait(); err != nil {
				t.Fatal(err)
			}
			if sent.IsZero() || final == nil {
				t.Fatal("missing ready/final save")
			}
			if len(final.Fleet[0].Cargo) != 0 {
				t.Fatalf("cargo persisted after exit: %v", final.Fleet[0].Cargo)
			}
			final.FinishCleanStop()
			if final.Won(econ.CleanStopContract) != tc.healthy {
				t.Fatalf("clean stop won=%v, want %v", final.Won(econ.CleanStopContract), tc.healthy)
			}
			if tc.healthy && final.Inv["moss"] != 3 {
				t.Fatal("unloaded cargo was lost")
			}
			if !tc.healthy && final.Inv["moss"] != 0 {
				t.Fatal("ignored cargo was rescued")
			}
			elapsed := time.Since(sent)
			if !tc.grace && elapsed > time.Second {
				t.Fatalf("unregistered shutdown waited %s", elapsed)
			}
			if tc.name == "ignores notification" && (elapsed < 1900*time.Millisecond || elapsed > 3*time.Second) {
				t.Fatalf("grace period %s", elapsed)
			}
		})
	}
}

func TestProcessColonyServerSignalsCleanStopAndDiscardsExitCargo(t *testing.T) {
	t.Setenv("ORBITAL_COLONY_UNIT", "2ms")
	r := testRunner(t)
	for _, clean := range []bool{false, true} {
		t.Run(map[bool]string{false: "natural exit loses cargo", true: "server signals and saves unload"}[clean], func(t *testing.T) {
			code := `package main;import "orbital/colony";func main(){colony.Report("leaving")}`
			if clean {
				body, err := os.ReadFile("testdata/colony/c13b-clean-stop/main.go")
				if err != nil {
					t.Fatal(err)
				}
				code = string(body)
			}
			prog := buildColonyProgram(t, r, code)
			c, err := newColonyGame(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			c.state.Contracts["c13-reed-clock"] = &econ.ContractRun{Best: 1}
			if err := c.state.Accept(econ.CleanStopContract); err != nil {
				t.Fatal(err)
			}
			c.state.Fleet[0].Cargo["moss"] = 3
			if err := c.start(filepath.Dir(prog), "en"); err != nil {
				t.Fatal(err)
			}
			defer c.stop()
			c.mu.Lock()
			done := c.active.done
			c.mu.Unlock()
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("server never stopped contract process")
			}
			blob, err := os.ReadFile(c.path)
			if err != nil {
				t.Fatal(err)
			}
			var saved econ.State
			if err := json.Unmarshal(blob, &saved); err != nil {
				t.Fatal(err)
			}
			if len(saved.Fleet[0].Cargo) != 0 || saved.Won(econ.CleanStopContract) != clean {
				t.Fatalf("unexpected exit state: %+v", saved)
			}
			if clean && saved.Inv["moss"] < 3 {
				t.Fatal("server lost unloaded cargo")
			}
			if !clean && saved.Contracts[econ.CleanStopContract].CargoLost != 3 {
				t.Fatal("loss across runs was forgotten")
			}
		})
	}
}
