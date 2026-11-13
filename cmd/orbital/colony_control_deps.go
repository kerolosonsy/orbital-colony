package main

import (
	"net/http"
	"strings"

	"orbital/colony/levels"
)

func colonyControlLevel() *levels.Level {
	level := *levels.Colony
	level.Deps = append(append([]string{}, level.Deps...), "modernc.org/sqlite", "github.com/prometheus/client_golang")
	return &level
}

func (s *server) checkControlImports(w http.ResponseWriter, r *http.Request, project *colonyProject) bool {
	c := s.requestGame(w, r)
	if c == nil {
		return false
	}
	level := &levels.Level{Kind: "lab"}
	c.mu.Lock()
	for node, dep := range map[string]string{"sys.warehouse": "modernc.org/sqlite", "sys.metrics": "github.com/prometheus/client_golang"} {
		if c.state.Has(node) {
			level.Deps = append(level.Deps, dep)
		}
	}
	c.mu.Unlock()
	for _, name := range projectNames(project.Files) {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		if missing := restrictedImports(project.Files[name], level, requestLang(r)); len(missing) > 0 {
			writeJSON(w, rejectedImports(missing, requestLang(r)))
			return false
		}
	}
	return true
}

func (c *colonyGame) controlRequirements(lang string) []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	node := ""
	switch c.state.Active {
	case "h26-control-api", "h28-drain":
		node = "sys.control"
	case "s29-warehouse", "s29-transfer":
		node = "sys.warehouse"
	case "o36-dashboard":
		node = "sys.metrics"
	}
	if node != "" && !c.state.Has(node) {
		return []string{tr(lang, "اشتري عقدة النظام الأول: ", "Buy the system node first: ") + node}
	}
	return nil
}
