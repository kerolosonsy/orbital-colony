package levels

const surveyExample = `sizes := []int{}
sizes = append(sizes, 5)
sizes = append(sizes, 2)
sizes = append(sizes, 8)
fmt.Println(sizes, len(sizes), sizes[0])`

func surveyTrace(captions []string) *LessonWalkthrough {
	frames := []LessonFrame{{Line: 0, Caption: captions[0]}}
	values := [][]int{{}, {5}, {5, 2}, {5, 2, 8}, {5, 2, 8}}
	for i, cells := range values {
		frame := LessonFrame{Line: i + 1, Caption: captions[i+1], Slice: &LessonSlice{Name: "sizes", Values: cells, Length: len(cells), ElementsOnly: true}}
		if i == len(values)-1 {
			frame.Output = "[5 2 8] 3 5"
		}
		frames = append(frames, frame)
	}
	return &LessonWalkthrough{Frames: frames}
}

var surveyLessonArabic = []LessonStep{
	{
		ID:    "c61.v2.mission-step",
		Title: "المهمة",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.mission",
				Kind: "markdown",
				Text: "1. قيس الـ8 كريستالات من الشمال لليمين.\n2. احفظ أحجامهم بالترتيب في slice واحدة.\n3. ابعت الـslice كلها مرة واحدة بـ~colony.Report~.\n\nالـstarter قاس أول كريستالة وحفظها فعلًا.\n\nالذهب: **15 ticks أو أقل**.",
			},
		},
	},
	{
		ID:    "c61.v2.idea-step",
		Title: "الفكرة",
		Blocks: []LessonBlock{
			{
				ID:        "c61.v2.trace",
				Kind:      "walkthrough",
				Text:      "تخيل الـslice صندوق، وكل حجم جديد بيتحط في آخره بـ~append~.",
				Code:      surveyExample,
				Imports:   []string{"fmt"},
				ExampleID: "e61-crystal-log-v2",
				Walkthrough: surveyTrace([]string{
					"قبل التشغيل.",
					"سجل فاضي؛ len=0.",
					"5 عند index 0.",
					"2 بعدها، مش بدلها.",
					"8 بقت آخر عنصر.",
					"السجل، طوله، وأول حجم.",
				}),
			},
		},
	},
	{
		ID:    "c61.v2.try-step",
		Title: "جرّب",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.sequence",
				Kind: "sandbox",
				Text: "ضيف الأحجام بالترتيب؛ راقب ~len~ والـindex وهما بيتغيروا.",
				Sandbox: &LessonSandbox{
					Kind:        "slice-append-sequence",
					Sequence:    []int{5, 2, 8, 4, 1, 7, 3, 6},
					SourceLabel: "أحجام للتوضيح، مش بيانات المهمة.",
					AppendLabel: "ضيف الحجم الجاي",
					ResetLabel:  "من الأول",
					DoneText:    "اكتمل السجل: 8 أحجام بالترتيب.",
				},
			},
		},
	},
	{
		ID:    "c61.v2.watch-step",
		Title: "خد بالك",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.watch",
				Kind: "markdown",
				Text: "- احتفظ بنتيجة ~append~: ~sizes = append(...)~.\n- أرقام الـindex من 0 لحد ~len(sizes)-1~؛ براها يحصل panic.",
			},
		},
	},
	{
		ID:    "c61.v2.check-step",
		Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.plan",
				Kind: "markdown",
				Text: "التكرار محتاج loop.\n\nلكل كريستالة: قيس، ضيف، واتجه للي بعدها.\n\nالتقرير مرة واحدة بعد الـloop؛ أول قراءة موجودة فعلًا.",
			},
			{
				ID:   "c61.v2.q-create",
				Kind: "question",
				Text: "عشان السجل يتجمع، تنشئ الـslice فين؟",
				Question: &LessonQuestion{
					Options: []string{
						"قبل الـloop",
						"جوه كل لفة",
						"بعد الـloop",
					},
					Answer: 0,
					Feedback: []string{
						"سجل واحد قبل الـloop يحتفظ بكل القراءات.",
						"إنشاءها كل لفة يبدأ سجل فاضي تاني.",
						"السجل لازم يبقى موجود قبل تجميع القراءات.",
					},
				},
			},
			{
				ID:   "c61.v2.q-append",
				Kind: "question",
				Text: "أنهي أداة تضيف حجم من غير مسح الموجود؟",
				Question: &LessonQuestion{
					Options: []string{
						"len",
						"append",
						"make",
					},
					Answer: 1,
					Feedback: []string{
						"len بتعد الموجود؛ مش بتضيف قراءة.",
						"append بتضيف في الآخر؛ احتفظ بالناتج.",
						"make بتنشىء slice، مش بتضيف للسجل الحالي.",
					},
				},
			},
			{
				ID:   "c61.v2.q-moves",
				Kind: "question",
				Text: "واقف عند أول كريستالة؛ كام حركة توصّلك للتامنة؟",
				Question: &LessonQuestion{
					Options: []string{
						"8",
						"6",
						"7",
					},
					Answer: 2,
					Feedback: []string{
						"8 حركات تحاول تعدّي آخر كريستالة.",
						"6 حركات تسيبك قبل الأخيرة بمسافة.",
						"8 قياسات + 7 حركات = 15 ticks؛ Report=0 وبافتراض مفيش تكلفة تانية.",
					},
				},
			},
		},
	},
	{
		ID:    "c61.v2.turn-step",
		Title: "دورك",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.turn",
				Kind: "markdown",
				Text: "اكتب أول محاولة بنفسك، شغّلها، وراجع السجل اللي طلع.",
			},
		},
	},
}

var surveyLessonEnglish = []LessonStep{
	{
		ID:    "c61.v2.mission-step",
		Title: "Mission",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.mission",
				Kind: "markdown",
				Text: "1. Measure the 8 crystals from left to right.\n2. Keep their sizes in order in one slice.\n3. Report the whole slice once with ~colony.Report~.\n\nThe starter already saved the first measurement.\n\nGold: **15 ticks or fewer**.",
			},
		},
	},
	{
		ID:    "c61.v2.idea-step",
		Title: "Idea",
		Blocks: []LessonBlock{
			{
				ID:        "c61.v2.trace",
				Kind:      "walkthrough",
				Text:      "A slice is a box: ~append~ adds each size at its end.",
				Code:      surveyExample,
				Imports:   []string{"fmt"},
				ExampleID: "e61-crystal-log-v2",
				Walkthrough: surveyTrace([]string{
					"Before execution.",
					"Empty slice; len=0.",
					"5 at index 0.",
					"2 follows 5.",
					"8 comes last.",
					"Log, length, first size.",
				}),
			},
		},
	},
	{
		ID:    "c61.v2.try-step",
		Title: "Try it",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.sequence",
				Kind: "sandbox",
				Text: "Add sizes in order; watch ~len~ grow and indexes advance.",
				Sandbox: &LessonSandbox{
					Kind:        "slice-append-sequence",
					Sequence:    []int{5, 2, 8, 4, 1, 7, 3, 6},
					SourceLabel: "Illustrative sizes, not mission data.",
					AppendLabel: "Add next size",
					ResetLabel:  "Reset",
					DoneText:    "Complete: 8 sizes in order.",
				},
			},
		},
	},
	{
		ID:    "c61.v2.watch-step",
		Title: "Watch out",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.watch",
				Kind: "markdown",
				Text: "- Keep append's result: ~sizes = append(...)~.\n- Indexes: 0 through ~len(sizes)-1~; outside means panic.",
			},
		},
	},
	{
		ID:    "c61.v2.check-step",
		Title: "Think & check",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.plan",
				Kind: "markdown",
				Text: "Use a loop for repeated work.\n\nFor each crystal: measure, append, move toward the next.\n\nReport once afterward; the starter saved the first reading.",
			},
			{
				ID:   "c61.v2.q-create",
				Kind: "question",
				Text: "Where should this mission's slice be created?",
				Question: &LessonQuestion{
					Options: []string{
						"Before the loop",
						"Inside each round",
						"After the loop",
					},
					Answer: 0,
					Feedback: []string{
						"One shared log keeps all readings.",
						"Each round would restart the log.",
						"The log must exist before collecting readings.",
					},
				},
			},
			{
				ID:   "c61.v2.q-append",
				Kind: "question",
				Text: "What adds a size without erasing earlier ones?",
				Question: &LessonQuestion{
					Options: []string{
						"len",
						"append",
						"make",
					},
					Answer: 1,
					Feedback: []string{
						"len counts; it adds nothing.",
						"append extends the log; save its result.",
						"make creates a slice, rather than appending.",
					},
				},
			},
			{
				ID:   "c61.v2.q-moves",
				Kind: "question",
				Text: "Starting on crystal 1, how many moves reach crystal 8?",
				Question: &LessonQuestion{
					Options: []string{
						"8",
						"6",
						"7",
					},
					Answer: 2,
					Feedback: []string{
						"Eight moves go past the last crystal.",
						"Six moves stop one crystal too soon.",
						"8 measures + 7 moves = 15 ticks; Report=0, assuming no other tick costs.",
					},
				},
			},
		},
	},
	{
		ID:    "c61.v2.turn-step",
		Title: "Your turn",
		Blocks: []LessonBlock{
			{
				ID:   "c61.v2.turn",
				Kind: "markdown",
				Text: "Write your first attempt, then run it.",
			},
		},
	},
}

const topThreeExample = `sizes := []int{4, 1, 3}
for i, size := range sizes {
	fmt.Println(i, size)
}
sorted := slices.Clone(sizes)
slices.Sort(sorted)
fmt.Println(sizes)
fmt.Println(sorted)`

func topThreeTrace(captions []string) *LessonWalkthrough {
	sizes := []int{4, 1, 3}
	slice := &LessonSlice{Name: "sizes", Values: sizes, Length: len(sizes), ElementsOnly: true}
	return &LessonWalkthrough{Frames: []LessonFrame{
		{Line: 0, Caption: captions[0]},
		{Line: 1, Caption: captions[1], Slice: slice},
		{Line: 2, Caption: captions[2], Variables: []LessonVariable{{Name: "i", Value: "0"}, {Name: "size", Value: "4"}}, Slice: slice},
		{Line: 3, Caption: captions[3], Variables: []LessonVariable{{Name: "i", Value: "0"}, {Name: "size", Value: "4"}}, Slice: slice, Output: "0 4"},
		{Line: 2, Caption: captions[4], Variables: []LessonVariable{{Name: "i", Value: "1"}, {Name: "size", Value: "1"}}, Slice: slice, Output: "0 4"},
		{Line: 3, Caption: captions[5], Variables: []LessonVariable{{Name: "i", Value: "1"}, {Name: "size", Value: "1"}}, Slice: slice, Output: "0 4\n1 1"},
		{Line: 2, Caption: captions[6], Variables: []LessonVariable{{Name: "i", Value: "2"}, {Name: "size", Value: "3"}}, Slice: slice, Output: "0 4\n1 1"},
		{Line: 3, Caption: captions[7], Variables: []LessonVariable{{Name: "i", Value: "2"}, {Name: "size", Value: "3"}}, Slice: slice, Output: "0 4\n1 1\n2 3"},
		{Line: 5, Caption: captions[8], Variables: []LessonVariable{{Name: "sorted", Value: "[4 1 3]"}}, Slice: slice, Output: "0 4\n1 1\n2 3"},
		{Line: 6, Caption: captions[9], Variables: []LessonVariable{{Name: "sorted", Value: "[1 3 4]"}}, Slice: slice, Output: "0 4\n1 1\n2 3"},
		{Line: 7, Caption: captions[10], Slice: slice, Output: "0 4\n1 1\n2 3\n[4 1 3]"},
		{Line: 8, Caption: captions[11], Slice: slice, Output: "0 4\n1 1\n2 3\n[4 1 3]\n[1 3 4]"},
	}}
}

var topThreeLessonArabic = []LessonStep{
	{
		ID: "c62.v1.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c62.v1.mission", Kind: "markdown",
			Text: "1. الـstarter قاس العشرة وسجل الحجم حسب index؛ الدرون عند x=9.\n2. اختار أكبر 3؛ حمولة الدرون 3، وبعدها احصدهم.\n3. هدف الذهب بيتغير حسب المواقع؛ شوف المعروض.",
		}},
	},
	{
		ID: "c62.v1.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{{
			ID: "c62.v1.trace", Kind: "walkthrough",
			Text: "في slice تدريب، range بتدي index وقيمة؛ Clone تنسخ، وSort ترتب.",
			Code: topThreeExample, Imports: []string{"fmt", "slices"}, ExampleID: "e62-range-clone-sort",
			Walkthrough: topThreeTrace([]string{
				"قبل التشغيل.", "أحجام للتوضيح؛ كل index ثابت.", "range أعطت أول index وقيمة.", "طبعناهم مع بعض.",
				"index التالي وقيمته.", "آخر عنصر في range.", "طبعنا آخر زوج.", "اللفة خلصت.",
				"Clone عملت slice جديدة.", "Sort رتبت النسخة بس.", "ترتيب الأصل فضل زي ما هو.", "النسخة المرتبة اتطبعت.",
			}),
		}},
	},
	{
		ID: "c62.v1.try-step", Title: "جرّب",
		Blocks: []LessonBlock{{
			ID: "c62.v1.sort-playground", Kind: "sandbox",
			Text: "رتّب الأصل أو نسخة؛ راقب الأحجام والـindex.",
			Sandbox: &LessonSandbox{
				Kind: "slice-sort-clone", Sequence: []int{5, 2, 8, 4, 1, 7, 3, 6},
				SourceLabel:    "أحجام للتوضيح، مش بيانات المهمة.",
				CloneSortLabel: "رتّب نسخة (Clone ثم Sort)", OriginalSortLabel: "رتّب الأصل", ResetLabel: "من الأول",
				CloneSortText:    "النسخة المرتبة تحت؛ الأصل فضل زي ما هو.",
				OriginalSortText: "الترتيب اتغير؛ الـindex مبقاش يحدد الكريستالة.",
				ResetText:        "رجعنا لترتيب البداية.",
			},
		}},
	},
	{
		ID: "c62.v1.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{
			ID: "c62.v1.watch", Kind: "markdown",
			Text: "- `b := a` بيشارك العناصر؛ Sort على b يغيّر a.\n- Sort تغيّر القيم، مش أماكن الكريستالات.\n- Measure وMove وHarvest: tick لكل أمر؛ Report والحساب: صفر.",
		}},
	},
	{
		ID: "c62.v1.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c62.v1.plan", Kind: "markdown", Text: "القياسات جاهزة. اختار الثلاثة، وخطط الزيارة."},
			{
				ID: "c62.v1.q-index", Kind: "question", Text: "في `[6, 1, 9, 3]`، أنهي index فيه أكبر قيمة؟",
				Question: &LessonQuestion{Options: []string{"0", "2", "3"}, Answer: 1, Feedback: []string{
					"index 0 فيه 6، مش الأكبر.", "صح، القيمة 9 عند index 2.", "index 3 فيه 3.",
				}},
			},
			{
				ID: "c62.v1.q-clone", Kind: "question", Text: "ليه Clone قبل Sort لو المواقع مهمة؟",
				Question: &LessonQuestion{Options: []string{
					"عشان ترتيب الأحجام الأصلي يفضل محفوظ.", "عشان ننقل الكريستالات لمواقع جديدة.", "عشان نزود سعة الـslice.",
				}, Answer: 0, Feedback: []string{
					"صح؛ النسخة تسيب ترتيب الأصل زي ما هو.", "Sort تغيّر القيم، مش العالم.", "Clone مش بتضمن سعة أكبر.",
				}},
			},
			{
				ID: "c62.v1.q-scan", Kind: "question", Text: "إمتى تختار أكبر 3؟",
				Question: &LessonQuestion{Options: []string{
					"قبل ما تعرف الأحجام.", "بعد ما تعرف أحجام العشرة.", "بعد أول 3 قياسات.",
				}, Answer: 1, Feedback: []string{
					"قارن العشرة عشان تعرف الأكبر.", "صح؛ قارن العشرة قبل القرار.", "أول 3 مش بالضرورة الأكبر.",
				}},
			},
		},
	},
	{
		ID: "c62.v1.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c62.v1.turn", Kind: "markdown", Text: "شغّل الـstarter؛ اكتب خطتك بنفسك."}},
	},
}

var topThreeLessonEnglish = []LessonStep{
	{
		ID: "c62.v1.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c62.v1.mission", Kind: "markdown",
			Text: "1. The starter measured 10 crystals and saved each size by index; drone at x=9.\n2. Pick the largest 3; carry limit is 3, then harvest.\n3. Gold varies by location; see the shown target.",
		}},
	},
	{
		ID: "c62.v1.idea-step", Title: "Idea",
		Blocks: []LessonBlock{{
			ID: "c62.v1.trace", Kind: "walkthrough",
			Text: "In a sample slice, range gives index and value; Clone copies, Sort orders.",
			Code: topThreeExample, Imports: []string{"fmt", "slices"}, ExampleID: "e62-range-clone-sort",
			Walkthrough: topThreeTrace([]string{
				"Before execution.", "Illustrative sizes; fixed indexes.", "First range index and value.", "Printed pair.",
				"Next index/value.", "Last range item.", "Printed last pair.", "Loop finished.",
				"Clone copied the slice.", "Sort changed the copy.", "Original order stayed.", "Printed the copy.",
			}),
		}},
	},
	{
		ID: "c62.v1.try-step", Title: "Try it",
		Blocks: []LessonBlock{{
			ID: "c62.v1.sort-playground", Kind: "sandbox",
			Text: "Sort original or copy; watch sizes and indexes.",
			Sandbox: &LessonSandbox{
				Kind: "slice-sort-clone", Sequence: []int{5, 2, 8, 4, 1, 7, 3, 6},
				SourceLabel:    "Illustrative sizes, not mission data.",
				CloneSortLabel: "Sort copy (Clone then Sort)", OriginalSortLabel: "Sort original", ResetLabel: "Reset",
				CloneSortText:    "Sorted copy below; the original stays unchanged.",
				OriginalSortText: "Order changed; indexes no longer identify crystals.",
				ResetText:        "Original order restored.",
			},
		}},
	},
	{
		ID: "c62.v1.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{
			ID: "c62.v1.watch", Kind: "markdown",
			Text: "- `b := a` shares elements; sorting b changes a.\n- Sort changes values, not crystal positions.\n- Measure, Move, Harvest: 1 tick each; Report and computation: 0.",
		}},
	},
	{
		ID: "c62.v1.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c62.v1.plan", Kind: "markdown", Text: "Measurements are ready. Choose three and plan visits."},
			{
				ID: "c62.v1.q-index", Kind: "question", Text: "In `[6, 1, 9, 3]`, which index holds the largest value?",
				Question: &LessonQuestion{Options: []string{"0", "2", "3"}, Answer: 1, Feedback: []string{
					"Index 0 holds 6.", "Correct: 9 is at index 2.", "Index 3 holds 3.",
				}},
			},
			{
				ID: "c62.v1.q-clone", Kind: "question", Text: "Why clone before sorting when positions matter?",
				Question: &LessonQuestion{Options: []string{
					"To preserve the original size order.", "To move crystals to new locations.", "To increase the slice capacity.",
				}, Answer: 0, Feedback: []string{
					"Correct; a copy preserves the original order.", "Sort changes values, not the world.", "Clone does not guarantee more capacity.",
				}},
			},
			{
				ID: "c62.v1.q-scan", Kind: "question", Text: "When choose the largest 3?",
				Question: &LessonQuestion{Options: []string{
					"Before you know the sizes.", "After you know all ten sizes.", "After the first three measurements.",
				}, Answer: 1, Feedback: []string{
					"Compare all ten sizes to know the largest.", "Correct: compare all ten before deciding.", "The first three may not be largest.",
				}},
			},
		},
	},
	{
		ID: "c62.v1.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c62.v1.turn", Kind: "markdown", Text: "Run the starter; write your plan."}},
	},
}

const censusExample = `counts := map[string]int{}
counts["moss"]++
counts["moss"]++
counts["lichen"]++
fmt.Println(counts["crystal"])
n, ok := counts["crystal"]
fmt.Println(n, ok, len(counts))`

func censusTrace(captions []string) *LessonWalkthrough {
	return &LessonWalkthrough{Frames: []LessonFrame{
		{Line: 0, Caption: captions[0]},
		{Line: 1, Caption: captions[1], Variables: []LessonVariable{{Name: "counts", Value: "{}"}, {Name: "len(counts)", Value: "0"}}},
		{Line: 2, Caption: captions[2], Variables: []LessonVariable{{Name: "counts", Value: "{moss: 1}"}, {Name: "len(counts)", Value: "1"}}},
		{Line: 3, Caption: captions[3], Variables: []LessonVariable{{Name: "counts", Value: "{moss: 2}"}, {Name: "len(counts)", Value: "1"}}},
		{Line: 4, Caption: captions[4], Variables: []LessonVariable{{Name: "counts", Value: "{lichen: 1, moss: 2}"}, {Name: "len(counts)", Value: "2"}}},
		{Line: 5, Caption: captions[5], Variables: []LessonVariable{{Name: "counts", Value: "{lichen: 1, moss: 2}"}, {Name: "len(counts)", Value: "2"}}, Output: "0"},
		{Line: 6, Caption: captions[6], Variables: []LessonVariable{{Name: "n", Value: "0"}, {Name: "ok", Value: "false"}, {Name: "len(counts)", Value: "2"}}, Output: "0"},
		{Line: 7, Caption: captions[7], Variables: []LessonVariable{{Name: "len(counts)", Value: "2"}}, Output: "0\n0 false 2"},
	}}
}

var censusLessonArabic = []LessonStep{
	{
		ID: "c63.v1.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c63.v1.mission", Kind: "markdown",
			Text: "1. عدّ moss وlichen وcrystal في خريطة 6×4 وابعت map الأعداد.\n2. عدّ empty اختياري؛ الـstarter سجّل خانة البداية بالفعل.\n3. الذهب 47 ticks؛ والـpar هو 60.",
		}},
	},
	{
		ID: "c63.v1.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{{
			ID: "c63.v1.trace", Kind: "walkthrough",
			Text: "الـmap بتربط key بعدد؛ القراية الغايبة ترجع صفر، وcomma-ok توضّح وجود المفتاح.",
			Code: censusExample, Imports: []string{"fmt"}, ExampleID: "e63-map-count",
			Walkthrough: censusTrace([]string{
				"قبل التنفيذ.", "map فاضية؛ len=0.", "أول عدّ لـmoss بقى 1.", "نفس المفتاح؛ len زي ما هي.",
				"مفتاح جديد؛ len بقت 2.", "قراية المفتاح الغايب رجعت صفر.", "ok بتقول إن المفتاح غايب.", "الخرج بيعرض مفتاحين.",
			}),
		}},
	},
	{
		ID: "c63.v1.try-step", Title: "جرّب",
		Blocks: []LessonBlock{{
			ID: "c63.v1.map-playground", Kind: "sandbox",
			Text: "اضغط نوعًا؛ عدّاده يزيد. جرّب قراءة مفتاح مش موجود.",
			Sandbox: &LessonSandbox{
				Kind: "map-count", Kinds: []string{"moss", "lichen", "crystal"}, MissingKey: "spore",
				SourceLabel:      "أنواع للتوضيح، مش بيانات المهمة.",
				MissingReadLabel: "اقرأ مفتاح غايب: spore", ResetLabel: "من الأول",
				MapOrderLabel: "الجدول للعرض بس؛ ترتيب range على map مش مضمون.",
			},
		}},
	},
	{
		ID: "c63.v1.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{
			ID: "c63.v1.watch", Kind: "markdown",
			Text: "- الكتابة على map nil تعمل panic؛ جهّزها الأول، وrange ترتيبها مش مضمون.\n- GetEntity وMove وHarvest كل واحد 1 tick، حتى Move الفاشلة؛ Report وWorldSize والحساب 0.",
		}},
	},
	{
		ID: "c63.v1.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c63.v1.plan", Kind: "markdown", Text: "جهّز map، وزُر الخانات وعدّ الأنواع، وبعدها ابعت التقرير."},
			{
				ID: "c63.v1.q-missing", Kind: "question", Text: "لو المفتاح `crystal` مش موجود، إيه نتيجة comma-ok؟",
				Question: &LessonQuestion{Options: []string{"0, false", "0, true", "nil, false"}, Answer: 0, Feedback: []string{
					"صح؛ القيمة صفر وok=false لأن المفتاح غايب.", "ok=true معناها إن المفتاح موجود.", "قيمة int الافتراضية هي صفر، مش nil.",
				}},
			},
			{
				ID: "c63.v1.q-len", Kind: "question", Text: "لو map فيها moss=4 وlichen=2، كام `len(counts)`؟",
				Question: &LessonQuestion{Options: []string{"6", "2", "4"}, Answer: 1, Feedback: []string{
					"6 هو مجموع القيم؛ len مش بتجمعها.", "صح؛ len بتعد المفتاحين.", "4 قيمة moss، مش عدد المفاتيح.",
				}},
			},
			{
				ID: "c63.v1.q-ticks", Kind: "question", Text: "بافتراض GetEntity مرة لكل 24 خانة و23 حركة بلا تكلفة أخرى، كام tick؟",
				Question: &LessonQuestion{Options: []string{"47", "48", "60"}, Answer: 0, Feedback: []string{
					"صح؛ 24 قراءة + 23 حركة = 47، وReport والحساب مجانيان.", "في زيارة 24 خانة، فيه 23 حركة بس.", "60 هو par، مش تقدير المسار المذكور.",
				}},
			},
		},
	},
	{
		ID: "c63.v1.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c63.v1.turn", Kind: "markdown", Text: "راجع طباعة الـstarter وحدد نوع خانة البداية اللي اتعدّت."}},
	},
}

var censusLessonEnglish = []LessonStep{
	{
		ID: "c63.v1.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c63.v1.mission", Kind: "markdown",
			Text: "1. Count moss, lichen, and crystal on the 6×4 map; report the counts map.\n2. Counting empty is optional; the starter already counted its starting tile.\n3. Gold is 47 ticks; Par is 60.",
		}},
	},
	{
		ID: "c63.v1.idea-step", Title: "Idea",
		Blocks: []LessonBlock{{
			ID: "c63.v1.trace", Kind: "walkthrough",
			Text: "A map connects keys to counts; a missing read returns zero, and comma-ok shows whether the key exists.",
			Code: censusExample, Imports: []string{"fmt"}, ExampleID: "e63-map-count",
			Walkthrough: censusTrace([]string{
				"Before execution.", "Empty map; len=0.", "First moss count is 1.", "Same key; len stays 1.",
				"New key; len is 2.", "Missing key reads as zero.", "ok says the key is absent.", "Output includes two keys.",
			}),
		}},
	},
	{
		ID: "c63.v1.try-step", Title: "Try it",
		Blocks: []LessonBlock{{
			ID: "c63.v1.map-playground", Kind: "sandbox",
			Text: "Click a kind to increment it; try reading an absent key.",
			Sandbox: &LessonSandbox{
				Kind: "map-count", Kinds: []string{"moss", "lichen", "crystal"}, MissingKey: "spore",
				SourceLabel:      "Illustrative kinds, not mission data.",
				MissingReadLabel: "Read absent key: spore", ResetLabel: "Reset",
				MapOrderLabel: "Sorted for display; map range order is unspecified.",
			},
		}},
	},
	{
		ID: "c63.v1.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{
			ID: "c63.v1.watch", Kind: "markdown",
			Text: "- A nil map write panics; initialize first; map range order is unspecified.\n- GetEntity, Move, Harvest: 1 tick each, even failed Move; Report, WorldSize, computation: 0.",
		}},
	},
	{
		ID: "c63.v1.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c63.v1.plan", Kind: "markdown", Text: "Prepare a map, visit and count the tiles, then report."},
			{
				ID: "c63.v1.q-missing", Kind: "question", Text: "If `crystal` is absent, what is comma-ok?",
				Question: &LessonQuestion{Options: []string{"0, false", "0, true", "nil, false"}, Answer: 0, Feedback: []string{
					"Correct: a missing read is 0, false.", "ok=true means the key exists.", "int's zero value is 0.",
				}},
			},
			{
				ID: "c63.v1.q-len", Kind: "question", Text: "If moss=4 and lichen=2, what is `len(counts)`?",
				Question: &LessonQuestion{Options: []string{"6", "2", "4"}, Answer: 1, Feedback: []string{
					"6 is the sum, not len.", "Correct: len counts two keys.", "4 is a value, not a key count.",
				}},
			},
			{
				ID: "c63.v1.q-ticks", Kind: "question", Text: "One GetEntity per 24 tiles, 23 moves, no other costs: ticks?",
				Question: &LessonQuestion{Options: []string{"47", "48", "60"}, Answer: 0, Feedback: []string{
					"Correct: 24 + 23 = 47; Report and computation are free.", "A 24-tile visit needs 23 moves.", "60 is Par, not this estimate.",
				}},
			},
		},
	},
	{
		ID: "c63.v1.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c63.v1.turn", Kind: "markdown", Text: "Check the starter output and identify the starting tile's counted kind."}},
	},
}
