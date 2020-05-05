package model

import (
	"errors"
)

var currentID = 1 // can use iota here

type Elevator struct {
	id              int
	currentPosition int
	topFloor        int
}

func (e *Elevator) GoTo(floor int) {

}

func GetElevator(topFloor int) (*Elevator, error) {

	if topFloor == 0 {
		return nil, errors.New("No Value for top floor")
	}

	e := &Elevator{
		id:              currentID,
		currentPosition: 0,
		topFloor:        topFloor,
	}

	currentID++

	return e, nil
}
