package econ

import "math"

const ResearchTicks = 100

func (b *Building) researchInput() string {
	for _, good := range sortedKeys(b.In) {
		if b.In[good] >= 1 && ResearchValue(good) > 0 {
			return good
		}
	}
	return ""
}

func (s *State) StepResearch(b *Building) {
	if !b.Enabled || b.Wear >= 100 {
		return
	}
	if b.Progress == 0 {
		good := b.researchInput()
		if good == "" {
			return
		}
		b.In[good]--
		b.ActiveRecipe = Recipe{Out: good, Ticks: ResearchTicks}
	}
	b.Progress++
	if b.Progress >= b.cycleTicks(b.ActiveRecipe) {
		s.RP += ResearchValue(b.ActiveRecipe.Out) * s.ResearchGainMult()
		b.Progress, b.ActiveRecipe = 0, Recipe{}
		b.Wear = math.Min(100, b.Wear+1)
	}
}
