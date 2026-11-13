package main

import (
	"context"

	"orbital/colony/levels"
)

func labMutantFile(lab *levels.Lab) string {
	if lab.MutantFile != "" {
		return lab.MutantFile
	}
	return "main.go"
}

func (r *runner) checkColonyPlanner(ctx context.Context, files map[string]string) (*runResponse, error) {
	level := &levels.Level{Kind: "lab", Lab: &levels.Lab{
		MutantFile: "plan/plan.go",
		Mutants: []levels.Mutant{
			{Desc: tr(contextLang(ctx), "المدخل الفاضي", "Empty input"), Code: "package plan\nfunc Next(ready []int) int { if len(ready)==0{return 0}; best:=0; for i:=range ready{if ready[i]<ready[best]{best=i}};return best }"},
			{Desc: tr(contextLang(ctx), "اختيار آخر تعادل", "Choosing the last tie"), Code: "package plan\nfunc Next(ready []int) int { best:=-1;for i:=range ready{if best<0||ready[i]<=ready[best]{best=i}};return best }"},
			{Desc: tr(contextLang(ctx), "اختيار الأبطأ", "Choosing the latest"), Code: "package plan\nfunc Next(ready []int) int { best:=-1;for i:=range ready{if best<0||ready[i]>ready[best]{best=i}};return best }"},
		},
		HiddenFiles: map[string]string{"plan/zz_hidden_test.go": plannerHiddenTests},
	}}
	for _, name := range projectNames(files) {
		level.Files = append(level.Files, levels.LabFile{Name: name, Starter: files[name]})
	}
	return r.RunLab(ctx, level, files)
}

const plannerHiddenTests = `package plan
import ("testing"; "slices")
func TestHiddenPlanner(t *testing.T) {
 for _, ready := range [][]int{nil, {}, {4}, {5,2,2,8}, {-2,-5,0}, {3,1,0}} {
  before := slices.Clone(ready)
  want := -1
  for i := range ready { if want<0 || ready[i]<ready[want] { want=i } }
  if got:=Next(ready);got!=want {t.Fatalf("Next(%v)=%d want %d",ready,got,want)}
  if !slices.Equal(ready,before) {t.Fatal("Next must not change its input")}
 }
}
`
