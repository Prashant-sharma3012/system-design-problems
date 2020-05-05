package model

import (
	"errors"
)

type Controller struct {
	Elevators      []*Elevator
	TotalElevators int
	FloorSwitches  []*FloorSwitch
}

func GetController(numOfLifts int, topFloor int) (*Controller, error) {

	if numOfLifts == 0 || topFloor == 0 {
		return nil, errors.New("Please provide number of lifts and total number of floors")
	}

	var elevators []*Elevator
	var elevator *Elevator

	for i := 0; i <= numOfLifts; i++ {
		elevator, _ = GetElevator(topFloor)
		elevators = append(elevators, elevator)
	}

	controller := &Controller{
		Elevators:      elevators,
		TotalElevators: numOfLifts,
		FloorSwitches:  GetSwitches(topFloor),
	}

	return controller, nil
}

func setDirection(val *Elevator, up bool, down bool) {
	if up {
		val.GoingUp = true
		val.InUse = true
	}

	if down {
		val.GoingDown = true
		val.InUse = true
	}
}

func (c *Controller) Call(atFloor int, up bool, down bool) *Elevator {

	var callFrom, minDiff, tempMinDiff int

	minDiff = 0

	// check if there is a lift already at this floor
	for indx, val := range c.Elevators {
		if val.CurrentPosition == atFloor {
			setDirection(val, up, down)
			return val
		}

		tempMinDiff = val.CurrentPosition - atFloor
		if tempMinDiff < 0 {
			tempMinDiff = tempMinDiff * -1
		}

		if tempMinDiff < minDiff {
			callFrom = indx
		}
	}

	setDirection(c.Elevators[callFrom], up, down)
	return c.Elevators[callFrom]
}
