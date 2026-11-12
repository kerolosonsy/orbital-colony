package levels

func init() {
	add(&Level{
		ID: "17-boss-rescue-route", Chapter: 17, Kind: "lab", Boss: true, Setup: seededLabRoom, Race: false,
		Title: "Boss: مسار كبسولة الإنقاذ", Goal: "اختار أقصر مسار generic واثبت العقد بتستات وbenchmark",
		Requires: []string{"tabletest", "benchmark"}, Bridge: []int{24, 27, 28, 29},
		Files: []LabFile{{Name: "main.go", Starter: `package main
func RelayPath[T comparable](links map[T][]T, start, goal T) ([]T,bool) {
    return nil,false
}
func main() {}
`}, {Name: "main_test.go", Starter: `package main
import "testing"
func TestRelayRoutes(t *testing.T) { t.Fatal("write route contract tests") }
func BenchmarkRelayRoutes(b *testing.B) {
    for i:=0;i<b.N;i++ { RelayPath(map[int][]int{1:{2}},1,2) }
}
`}},
		Lesson: `# Boss: مسار كبسولة الإنقاذ

## قبل ما تكتب
- إيه المدخلات اللي لازم تتراجع قبل أي تأثير جانبي؟
- مين مالك النتيجة وكل مورد بعد الخطأ؟

## اختبار التجميع
طوارئ المستعمرة بتجمع أدوات إنت عارفها. مفيش مفهوم جديد. راجع: 14-1-generic-max, 14-4-airlock-queue, 15-1-bfs, 16-1-mutants, 16-3-benchmark, 17-3-deep-geode.

مثال أصغر ومختلف:
~~~go
visited := map[string]bool{"airlock": true}
queue := []string{"airlock"}
_ = visited
_ = queue
~~~

## أخطاء شائعة
~~~go
if len(queue) == 0 { return []T{goal}, true } // invents a route
~~~
ده بيضيّع فشل أو بيغيّر الملكية؛ التزم بالعقد المكتوب.

## المهمة
نفّذ ~RelayPath[T comparable](links map[T][]T, start, goal T) ([]T, bool)~ لمسارات كبسولة تحت الأرض. المفاتيح محطات والجيران وصلات موجهة. مثال نداء: ~RelayPath(map[int][]int{4: {8}}, 4, 8)~. ارجع أقصر مسار بعدد الوصلات شامل الطرفين مع true، أو ~nil, false~ لو مفيش طريق. نفس الطرفين يرجّع عنصر واحد حتى لو الخريطة nil. خريطة المدخل مبتتغيّرش، والدورات متعلّقش البحث.

اكتب table-driven tests في ~main_test.go~ تمسك 3 نسخ غلط: فشل دايماً، نجاح وهمي لوجهة غايبة، واختيار الطريق الأطول. اكتب ~BenchmarkRelayRoutes~ زي 16-3؛ مش مطلوب رقم سرعة معيّن. سيب ~main~ فاضية. النجمة التالتة: عند التعادل، استخدم ترتيب الجيران الأصلي في BFS، وارجع مسار مملوك للمنادي حتى مع إحداثيات 3D كنوع generic.
`,
		Hints: []string{"استخدم BFS وخريطة seen عشان الدورات.", "خزّن الأب لكل محطة وارجع المسار بالعكس.", "اختبر نفس الطرفين والغايب وطريقين أطوالهم مختلفة؛ الـ bonus بيختبر التعادل."},
		En: &LevelText{Title: "Boss: Rescue Capsule Route", Goal: "Find a generic shortest route and prove its contract with tests and a benchmark", Files: []LabFile{{Name: "main.go", Starter: `package main
func RelayPath[T comparable](links map[T][]T, start, goal T) ([]T,bool) {
    return nil,false
}
func main() {}
`}, {Name: "main_test.go", Starter: `package main
import "testing"
func TestRelayRoutes(t *testing.T) { t.Fatal("write route contract tests") }
func BenchmarkRelayRoutes(b *testing.B) {
    for i:=0;i<b.N;i++ { RelayPath(map[int][]int{1:{2}},1,2) }
}
`}},
			Lesson: `# Boss: Rescue Capsule Route

## Before you write
- Which inputs must be checked before any side effect?
- Who owns the result and each resource after an error?

## Integration checkpoint
A colony emergency combines tools you already know. No new concept is needed. Revisit: 14-1-generic-max, 14-4-airlock-queue, 15-1-bfs, 16-1-mutants, 16-3-benchmark, 17-3-deep-geode.

A smaller, different example:
~~~go
visited := map[string]bool{"airlock": true}
queue := []string{"airlock"}
_ = visited
_ = queue
~~~

## Common mistakes
~~~go
if len(queue) == 0 { return []T{goal}, true } // invents a route
~~~
This drops a failure or changes ownership; preserve the documented contract.

## Mission
Implement ~RelayPath[T comparable](links map[T][]T, start, goal T) ([]T, bool)~ for an underground capsule network. Keys are stations and neighbors are directed links. Example call: ~RelayPath(map[int][]int{4: {8}}, 4, 8)~. Return the shortest path by edge count including both endpoints and true, or ~nil, false~ when unreachable. Equal endpoints return one element even with a nil graph. Do not mutate the graph; cycles must terminate.

Write table-driven tests in ~main_test.go~ catching three mutants: always failing, inventing a route to a missing destination, and choosing a longer route. Write ~BenchmarkRelayRoutes~ as in 16-3; no speed threshold is required. Leave ~main~ empty. Third star: break shortest-path ties using the original BFS neighbor order and return caller-owned paths, including generic 3D coordinate keys.
`,
			Hints: []string{"Use BFS and a seen map to terminate cycles.", "Store each station parent and reverse the reconstructed path.", "Test equal endpoints, unreachable goals, and unequal routes; the bonus checks ties."}},
		Lab: &Lab{ReqFile: "main_test.go", Bonus: true, Mutants: []Mutant{{Desc: "خرق عقد المسار", Code: `package main
func RelayPath[T comparable](links map[T][]T, start, goal T) ([]T,bool) {
    return nil,false
}
func main() {}
`}, {Desc: "خرق عقد المسار", Code: `package main
func RelayPath[T comparable](links map[T][]T, start, goal T) ([]T,bool) {
    return []T{start,goal},true
}
func main() {}
`}, {Desc: "خرق عقد المسار", Code: `package main
func RelayPath[T comparable](links map[T][]T, start, goal T) ([]T,bool) {
    var visit func(T, map[T]bool) []T
    visit=func(node T, seen map[T]bool) []T {
        if node==goal{return []T{node}}
        seen[node]=true
        for _,next:=range links[node] {
            if !seen[next] {if path:=visit(next,seen);path!=nil{return append([]T{node},path...)}}
        }
        return nil
    }
    path:=visit(start,map[T]bool{})
    return path,path!=nil
}
func main() {}
`}}, Hidden: `package main
import ("testing";"math/rand";"reflect")
func TestHiddenRelayShortest(t *testing.T) {
    for _,seed:=range []int64{17091,17092,17093,17094,17095} {
        rng:=rand.New(rand.NewSource(seed))
        n:=5+rng.Intn(30)
        graph:=map[int][]int{}
        for i:=0;i<n;i++ {graph[i]=[]int{i+1,i}} // reachable chain plus cycles
        graph[0]=append(graph[0],n)
        before:=append([]int(nil),graph[0]...)
        got,ok:=RelayPath(graph,0,n)
        if !ok || !reflect.DeepEqual(got,[]int{0,n}) {t.Fatal(seed,got,ok)}
        if !reflect.DeepEqual(graph[0],before) {t.Fatal("graph mutated")}
        got,ok=RelayPath(graph,0,n+1)
        if ok || got!=nil {t.Fatal("unreachable must return nil,false",got,ok)}
    }
    got,ok:=RelayPath(map[string][]string{"dock":{"shaft"},"shaft":{"dock","lab"}},"dock","lab")
    if !ok || !reflect.DeepEqual(got,[]string{"dock","shaft","lab"}){t.Fatal(got,ok)}
    got,ok=RelayPath[string](nil,"solo","solo")
    if !ok || !reflect.DeepEqual(got,[]string{"solo"}) {t.Fatal("same endpoint",got,ok)}
}
func TestHiddenBonusRelayTieAndOwnership(t *testing.T) {
    type cell struct {X,Y,Z int}
    a,b,c,d:=cell{},cell{1,0,0},cell{0,1,0},cell{0,0,1}
    graph:=map[cell][]cell{a:{c,b,c},b:{d},c:{d},d:{a}}
    got,ok:=RelayPath(graph,a,d)
    if !ok || !reflect.DeepEqual(got,[]cell{a,c,d}){t.Fatal("tie must use neighbor order",got,ok)}
    got[0]=d
    again,ok:=RelayPath(graph,a,d)
    if !ok || !reflect.DeepEqual(again,[]cell{a,c,d}){t.Fatal("result ownership",again,ok)}
}
`},
	})
}
