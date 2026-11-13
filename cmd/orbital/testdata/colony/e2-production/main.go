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

func main() {
	drone := colony.Main()
	must(drone.Build(colony.Composter))
	composter := drone.Here()
	must(composter.SetRecipe(colony.Biomass))
	for batch := 0; batch < 2; batch++ {
		n, err := composter.Put(colony.Moss, 8)
		must(err)
		if n != 8 {
			panic(i18n.T("توصيل الطحالب ناقص", "incomplete moss delivery"))
		}
		for composter.Status().Out[colony.Biomass] < 2 {
			drone.Wait()
		}
		n, err = composter.Take(colony.Biomass, 2)
		must(err)
		if n != 2 {
			panic(i18n.T("جمع الكتلة الحيوية ناقص", "incomplete biomass collection"))
		}
	}
	if !drone.Move(colony.East) || !drone.Melt() {
		panic(i18n.T("تسييح المياه فشل", "could not melt water"))
	}
	must(colony.Build(colony.Refinery))
	refinery := colony.Here()
	must(refinery.SetRecipe(colony.Fuel))
	_, err := refinery.Put(colony.Biomass, 2)
	must(err)
	_, err = refinery.Put(colony.Water, 1)
	must(err)
	for refinery.Status().Out[colony.Fuel] < 1 {
		colony.Wait()
	}
	_, err = refinery.Take(colony.Fuel, 1)
	must(err)
	if colony.Export(colony.Fuel, 1) <= 0 {
		panic(i18n.T("بيع الوقود فشل", "fuel was not sold"))
	}
	for {
		colony.Wait()
	}
}
