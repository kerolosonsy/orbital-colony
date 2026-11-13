package econ

var testedPlannerContract = &Contract{ID: "c16-tested-planner", Chapter: 16, Requires: []string{"sel:plan.Next"}, Window: 1500,
	Title: "مخطط متختبر", TitleEn: "Tested Planner",
	Brief:   "سلّم 40 من أي محصول. اكتب في plan/plan.go دالة pure اسمها Next(ready []int) int بترجّع index أصغر وقت، أول index عند التعادل، و-1 للمدخل الفاضي من غير تغيير المدخل. اكتب تستاتك عشان تمسك 3 mutants. استخدم player/plan في برنامجك. الدرس: فصل 16.",
	BriefEn: "Deliver 40 of any crop. In plan/plan.go write a pure Next(ready []int) int returning the earliest-ready index, the first index on ties, and -1 for empty input, without changing the input. Write tests that catch 3 mutants. Use player/plan in your program. Lesson: chapter 16.",
	Goal:    Goal{Any: 40, MaxLost: NoLimit}, Reward: Reward{Credits: 500}}
