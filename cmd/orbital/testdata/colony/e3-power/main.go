package main

import (
	"orbital/colony"
	"orbital/colony/i18n"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func put(b *colony.Building, e colony.Entity, n int) {
	moved, err := b.Put(e, n)
	must(err)
	if moved != n {
		panic(i18n.T("المدخلات ناقصة", "incomplete inputs"))
	}
}

func take(b *colony.Building, e colony.Entity, n int) {
	for b.Status().Out[e] < float64(n) {
		colony.Wait()
	}
	moved, err := b.Take(e, n)
	must(err)
	if moved != n {
		panic(i18n.T("المخرجات ناقصة", "incomplete outputs"))
	}
}

func main() {
	must(colony.Build(colony.Generator))
	generator := colony.Here()
	must(generator.SetMode("solar"))
	colony.Move(colony.North)
	must(colony.Build(colony.Composter))
	composter := colony.Here()
	must(composter.SetRecipe(colony.Biomass))
	put(composter, colony.Moss, 8)
	take(composter, colony.Biomass, 2)
	put(composter, colony.Moss, 8)
	take(composter, colony.Biomass, 2)
	colony.Move(colony.South)
	colony.Move(colony.East)
	if !colony.Melt() {
		panic(i18n.T("تسييح المياه فشل", "melting failed"))
	}
	must(colony.Build(colony.Refinery))
	refinery := colony.Here()
	must(refinery.SetRecipe(colony.Fuel))
	must(refinery.SetPriority(10))
	put(refinery, colony.Biomass, 2)
	put(refinery, colony.Water, 1)
	colony.Move(colony.East)
	must(colony.Build(colony.Kiln))
	kiln := colony.Here()
	must(kiln.SetRecipe(colony.Glass))
	must(kiln.SetPriority(5))
	put(kiln, colony.Silica, 3)
	supply, demand := colony.Power()
	if supply != 10 || demand != 11 || composter.Status().Powered {
		panic(i18n.T("توزيع الطاقة بالأولوية فشل", "priority allocation failed"))
	}
	colony.Move(colony.West)
	take(refinery, colony.Fuel, 1)
	must(refinery.Enable(false))
	put(refinery, colony.Biomass, 2)
	put(refinery, colony.Water, 1)
	colony.Move(colony.West)
	put(generator, colony.Fuel, 1)
	for colony.Daylight() {
		colony.Wait()
	}
	if supply, _ := colony.Power(); supply != 7 {
		panic(i18n.T("الحد الأدنى للطاقة بالليل اتغيّر", "night baseline power changed"))
	}
	must(generator.SetMode("burner"))
	must(refinery.Enable(true))
	colony.Move(colony.East)
	take(refinery, colony.Fuel, 1)
	colony.Move(colony.West)
	if !colony.Main().Unload() {
		panic(i18n.T("تفريغ الوقود فشل", "fuel unloading failed"))
	}
	for {
		colony.Wait()
	}
}
