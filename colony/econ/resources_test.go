package econ

import (
	"math"
	"testing"
)

func TestRecipesUseKnownInputsAtOrBelowTheirTier(t *testing.T) {
	if len(RawV3) != 7 || len(Crafted) != 14 || len(Recipes) != 15 {
		t.Fatalf("resource tables: %d raw, %d crafted, %d recipes", len(RawV3), len(Crafted), len(Recipes))
	}
	for _, goods := range [][]Resource{RawV3, Crafted} {
		for _, good := range goods {
			if good.Name == "" || good.NameEn == "" || good.BasePrice <= 0 || PriceV3[good.ID] != good.BasePrice {
				t.Errorf("incomplete resource: %+v", good)
			}
		}
	}
	for _, r := range Recipes {
		if TierOf(r.Out) == 0 || r.Ticks <= 0 || FindNode(r.Building) == nil {
			t.Errorf("invalid recipe: %+v", r)
		}
		for good, qty := range r.In {
			if _, ok := PriceV3[good]; !ok || qty <= 0 {
				t.Errorf("%s consumes unknown or invalid input %s: %g", r.Out, good, qty)
			}
			if TierOf(good) > TierOf(r.Out) {
				t.Errorf("%s consumes higher-tier %s", r.Out, good)
			}
		}
	}
	for _, good := range Goods {
		if PriceV3[good] != BasePrice[good] {
			t.Errorf("old price changed for %s", good)
		}
	}
	biomass := RecipeFor("biomass")
	if len(biomass) != 2 || biomass[0].Alt || !biomass[1].Alt || len(RecipeFor("missing")) != 0 {
		t.Fatal("recipe lookup lost biomass alternatives or accepted an unknown output")
	}
}

func TestCraftedGoodsAreAcyclicAndReachableFromRaw(t *testing.T) {
	visiting, reached := map[string]bool{}, map[string]bool{}
	for _, good := range Goods {
		reached[good] = true
	}
	for _, good := range RawV3 {
		reached[good.ID] = true
	}
	var visit func(string)
	visit = func(good string) {
		if visiting[good] {
			t.Fatalf("recipe cycle at %s", good)
		}
		if reached[good] {
			return
		}
		visiting[good] = true
		recipes := RecipeFor(good)
		if len(recipes) == 0 {
			t.Fatalf("%s has no raw source or recipe", good)
		}
		for _, r := range recipes {
			for input := range r.In {
				visit(input)
			}
		}
		visiting[good], reached[good] = false, true
	}
	for _, good := range Crafted {
		visit(good.ID)
	}
}

func TestCraftedPricesRoundPrimaryRecipeValue(t *testing.T) {
	for _, good := range Crafted {
		primary := 0
		for _, r := range RecipeFor(good.ID) {
			if r.Alt {
				continue
			}
			primary++
			inputs := 0.0
			for input, qty := range r.In {
				inputs += PriceV3[input] * qty
			}
			if want := math.Round(1.6 * inputs); PriceV3[good.ID] != want {
				t.Errorf("%s price = %g, want %g", good.ID, PriceV3[good.ID], want)
			}
		}
		if primary != 1 {
			t.Errorf("%s has %d primary recipes", good.ID, primary)
		}
	}
}
