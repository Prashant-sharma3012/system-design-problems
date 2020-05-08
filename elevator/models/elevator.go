package model

import (
	"errors"
	"fmt"
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
	time.Sleep(100)
	fmt.Println("Going to")
	fmt.Println(floor)

	if (e.GoingDown && floor < e.CurrentPosition) ||
		(e.GoingUp && floor > e.CurrentPosition) {
		e.CurrentPosition = floor
		e.InUse = false
	}
	e.CurrentPosition = floor
	e.InUse = false
}

func GetElevator(topFloor int) (*Elevator, error) {

	if topFloor == 0 {
		return nil, errors.New("No Value for top floor")
	}

	e := &Elevator{
		Id:              currentID,
		CurrentPosition: 0,
		TopFloor:        topFloor,
		GoingDown:       false,
		GoingUp:         false,
		InUse:           false,
	}

	currentID++

	return e, nil
}
