package main

import "orbital/colony/econ"

func experimentFixture(c *econ.Contract) *econ.State {
	parent := econ.FindContract(c.After[0])
	st := contractFixture(parent)
	st.Fail()
	st.Contracts[parent.ID] = &econ.ContractRun{Status: "done", Best: 1}
	for _, id := range parent.Reward.Grants {
		st.Blueprints[id] = true
		st.Levels[id]++
	}
	st.Normalize()
	if err := st.Accept(c.ID); err != nil {
		panic(err)
	}
	return st
}
