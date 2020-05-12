package model

import (
	"errors"
	"time"
)

var currentID = 1 // can use iota here

type Elevator struct {
	Id              int
	CurrentPosition int
	TopFloor        int
	GoingUp         bool
	GoingDown       bool
	InUse           bool
}

func (e *Elevator) GoTo(floor int) {

	if (e.GoingDown && floor <= e.CurrentPosition) ||
		(e.GoingUp && floor >= e.CurrentPosition) {

		diff := floor - e.CurrentPosition

		if diff < 0 {
			diff = diff * -1
		}

		for i := 1; i <= diff; i++ {
			time.Sleep(99999)
			if e.GoingDown {
				e.CurrentPosition--
			} else {
				e.CurrentPosition++
			}
		}

		e.CurrentPosition = floor
		e.InUse = false
		e.GoingDown = false
		e.GoingDown = false
	}
}

func GetElevator(topFloor int) (*Elevator, error) {

	if topFloor == 0 {
		return nil, errors.New("No Value for top floor")
	}

	e := &Elevator{
		Id:              currentID,
		CurrentPosition: 1,
		TopFloor:        topFloor,
		GoingDown:       false,
		GoingUp:         false,
		InUse:           false,
	}

	currentID++

	return e, nil
}
