package colony

func Power() (supply, demand float64)          { return world.Power() }
func Daylight() bool                           { return world.Daylight() }
func Storage() (used, capacity map[string]int) { return world.Storage() }
func Carry() map[Entity]int                    { return world.Carry(0) }
func CargoFree() int                           { return world.CargoFree(0) }
func BatteryLevel() int                        { return world.Battery(0) }
func Charge() bool                             { return world.Charge(0) }

func (d *Drone) Carry() map[Entity]int { return world.Carry(d.id) }
func (d *Drone) CargoFree() int        { return world.CargoFree(d.id) }
func (d *Drone) Battery() int          { return world.Battery(d.id) }
func (d *Drone) Charge() bool          { return world.Charge(d.id) }

func (b *Building) SetMode(mode string) error { return world.BuildingSetMode(b.drone, b.pos, mode) }
func (b *Building) SetPriority(n int) error   { return world.BuildingSetPriority(b.drone, b.pos, n) }
func (b *Building) Enable(on bool) error      { return world.BuildingEnable(b.drone, b.pos, on) }

func TryMelt() (int, error)            { return world.TryMelt(0) }
func TrySift() (int, error)            { return world.TrySift(0) }
func (d *Drone) TryMelt() (int, error) { return world.TryMelt(d.id) }
func (d *Drone) TrySift() (int, error) { return world.TrySift(d.id) }

func TryDig(dir Direction) error            { return world.TryDig(0, dir) }
func TryDigDown() error                     { return world.TryDigDown(0) }
func (d *Drone) TryDig(dir Direction) error { return world.TryDig(d.id, dir) }
func (d *Drone) TryDigDown() error          { return world.TryDigDown(d.id) }
