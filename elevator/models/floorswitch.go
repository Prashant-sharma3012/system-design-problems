package model

type FloorSwitch struct {
	FloorNumber int
	Up          bool
	Down        bool
}

func GetSwitches(totalSwitches int) []*FloorSwitch {
	var f []*FloorSwitch

	for i := 1; i <= totalSwitches; i++ {
		f = append(f, &FloorSwitch{
			FloorNumber: i,
			Up:          false,
			Down:        false,
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
