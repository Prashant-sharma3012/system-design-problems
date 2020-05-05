package model

type FloorSwitch struct {
	Up   bool
	Down bool
}

func GetSwitches(totalSwitches int) []*FloorSwitch {
	var f []*FloorSwitch

	for i := 0; i <= totalSwitches; i++ {
		f = append(f, &FloorSwitch{
			Up:   false,
			Down: false,
		})
	}

	return f
}

func (f *FloorSwitch) GoUp() {
	f.Up = true
}

func (f *FloorSwitch) GoDown() {
	f.Down = true
}
