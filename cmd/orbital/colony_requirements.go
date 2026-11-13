package main

import "strings"

func projectRequirements(files map[string]string, reqs []string, lang string) []string {
	combined := &features{extra: map[string]bool{}, sels: map[string]int{}, calls: map[string]int{}, userFuncs: map[string]int{}}
	declared := map[string]int{}
	calls := map[string]int{}
	for _, name := range projectNames(files) {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		ft, err := scan(files[name])
		if err != nil {
			continue
		}
		for function := range ft.userFuncs {
			declared[ft.pkg+"."+function] = 0
		}
		for function, count := range ft.calls {
			if !strings.Contains(function, ".") {
				function = ft.pkg + "." + function
			}
			calls[function] += count
		}
		combined.add(ft)
	}
	for name := range declared {
		declared[name] = calls[name]
	}
	combined.userFuncs = declared
	var missing []string
	for _, req := range reqs {
		var alternatives []string
		met := false
		for _, alt := range strings.Split(req, "|") {
			unmet := checkFeatureRequirements(combined, []string{alt}, lang)
			if len(unmet) == 0 {
				met = true
				break
			}
			alternatives = append(alternatives, unmet...)
		}
		if !met {
			missing = append(missing, strings.Join(alternatives, tr(lang, " أو ", " or ")))
		}
	}
	return missing
}

func (ft *features) add(other *features) {
	ft.loops += other.loops
	ft.ifs += other.ifs
	ft.switches += other.switches
	ft.errChecks += other.errChecks
	ft.ranges += other.ranges
	ft.slices += other.slices
	ft.maps += other.maps
	ft.structs += other.structs
	ft.methods += other.methods
	ft.ptrParams += other.ptrParams
	ft.addrOf += other.addrOf
	ft.ptrRecv += other.ptrRecv
	ft.stringers += other.stringers
	ft.typeSwitches += other.typeSwitches
	ft.goStmts += other.goStmts
	ft.selects += other.selects
	ft.chans += other.chans
	ft.genericFuncs += other.genericFuncs
	ft.genericTypes += other.genericTypes
	ft.benchmarks += other.benchmarks
	for name, present := range other.extra {
		ft.extra[name] = ft.extra[name] || present
	}
	for name, count := range other.sels {
		ft.sels[name] += count
	}
	for name, count := range other.calls {
		ft.calls[name] += count
	}
	for name, count := range other.userFuncs {
		ft.userFuncs[name] += count
	}
}
