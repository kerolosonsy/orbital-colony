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
