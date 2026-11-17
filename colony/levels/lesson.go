package levels

type LessonStep struct {
	ID     string        `json:"id"`
	Title  string        `json:"title"`
	Blocks []LessonBlock `json:"blocks"`
}

// A question's block ID is also its persistence ID in both languages.
type LessonBlock struct {
	ID          string             `json:"id"`
	Kind        string             `json:"kind"`
	Text        string             `json:"text"`
	Note        string             `json:"note,omitempty"`
	Optional    bool               `json:"optional,omitempty"`
	Code        string             `json:"code,omitempty"`
	Imports     []string           `json:"imports,omitempty"`
	ExampleID   string             `json:"exampleId,omitempty"`
	Walkthrough *LessonWalkthrough `json:"walkthrough,omitempty"`
	Question    *LessonQuestion    `json:"question,omitempty"`
	Sandbox     *LessonSandbox     `json:"sandbox,omitempty"`
}

type LessonWalkthrough struct {
	Frames []LessonFrame `json:"frames"`
}

type LessonFrame struct {
	Line      int              `json:"line"` // zero is the initial state, before execution
	Caption   string           `json:"caption"`
	Variables []LessonVariable `json:"variables,omitempty"`
	Slice     *LessonSlice     `json:"slice,omitempty"`
	Output    string           `json:"output"`
	Panic     string           `json:"panic,omitempty"`
}

type LessonVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LessonSlice struct {
	ElementsOnly bool   `json:"elementsOnly,omitempty"` // omit storage details in introductory lessons
	Name         string `json:"name"`
	Store        string `json:"store,omitempty"`
	Values       []int  `json:"values"`
	Length       int    `json:"length"`
	Capacity     *int   `json:"capacity,omitempty"` // absent when no exact capacity is asserted
	Nil          bool   `json:"nil,omitempty"`
}

type LessonQuestion struct {
	Options  []string `json:"options"`
	Answer   int      `json:"answer"`
	Feedback []string `json:"feedback"`
}

type LessonSandbox struct {
	Kind          string              `json:"kind"`
	Presets       []LessonSlicePreset `json:"presets,omitempty"`
	AppendValue   int                 `json:"appendValue,omitempty"`
	Limit         int                 `json:"limit,omitempty"`
	GrowthLabel   string              `json:"growthLabel,omitempty"`
	ReservedLabel string              `json:"reservedLabel,omitempty"`
	LimitLabel    string              `json:"limitLabel,omitempty"`
	AppendText    string              `json:"appendText,omitempty"`
	ResetText     string              `json:"resetText,omitempty"`
	Sequence      []int               `json:"sequence,omitempty"`
	SourceLabel   string              `json:"sourceLabel,omitempty"`
	AppendLabel   string              `json:"appendLabel,omitempty"`
	ResetLabel    string              `json:"resetLabel,omitempty"`
	DoneText      string              `json:"doneText,omitempty"`
}

type LessonSlicePreset struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Values   []int  `json:"values"`
	Capacity int    `json:"capacity"`
	Text     string `json:"text"`
}
