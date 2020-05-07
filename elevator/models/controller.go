package model

import (
	"errors"
	"fmt"
)

type Request struct {
	Id        int
	Floor     int
	GoingUp   bool
	GoingDown bool
}

type Controller struct {
	Elevators      []*Elevator
	TotalElevators int
	FloorSwitches  []*FloorSwitch
	start          bool
	requestQ       []*Request
}

// helper fucntions
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

func resetSwitch(c *Controller, atFloor int, up bool, down bool) {

	if up {
		c.FloorSwitches[atFloor-1].Up = false
	}

	if down {
		c.FloorSwitches[atFloor-1].Down = false
	}
}

func createRequest(floor int, up bool, down bool, requestId int) *Request {
	return &Request{
		Id:        requestId,
		Floor:     floor,
		GoingUp:   up,
		GoingDown: down,
	}
}

func removeRequest(c *Controller, reqId int) {
	posToRemove := -1

	for indx, val := range c.requestQ {
		if val.Id == reqId {
			posToRemove = indx
		}
	}

	if posToRemove == -1 {
		// something is wrong
		fmt.Println("##################################")
		fmt.Println("req id Not in queue")
		fmt.Println(reqId)
		fmt.Println("##################################")
	}

	c.requestQ = append(c.requestQ[:posToRemove], c.requestQ[posToRemove+1:]...)
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

func (c *Controller) Call(atFloor int, up bool, down bool) *Elevator {

	currentReq := len(c.requestQ)
	c.requestQ = append(c.requestQ, createRequest(atFloor, up, down, currentReq))

	var callFrom = -1
	var minDiff, tempMinDiff int

	minDiff = 0

	// keep running until you find a lift
	for callFrom == -1 {
		// check if there is a lift already at this floor
		for indx, val := range c.Elevators {

			if !val.InUse {
				if val.CurrentPosition == atFloor {
					val.InUse = true
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
		}
	}

	c.Elevators[callFrom].InUse = true
	setDirection(c.Elevators[callFrom], up, down)

	// set the switch to off again
	resetSwitch(c, atFloor, up, down)
	removeRequest(c, currentReq)

	return c.Elevators[callFrom]
}

func (c *Controller) StartServicing() {
	c.start = true

	for c.start {
		for _, val := range c.FloorSwitches {
			if val.Up {
				c.Call(val.FloorNumber, true, false)
			}

			if val.Down {
				c.Call(val.FloorNumber, false, true)
			}
		}
	}
}

func (c *Controller) StopServicing() {
	c.start = false
}

//will replce this
func (c *Controller) RequestFromFloor(floor int, up bool, down bool) {
	if up {
		c.FloorSwitches[floor-1].GoUp()
	}

	if down {
		c.FloorSwitches[floor-1].GoDown()
	}
}

// management

func (c *Controller) OccupiedElevators() []*Elevator {
	var occupied []*Elevator

	for _, val := range c.Elevators {
		if val.InUse {
			occupied = append(occupied, val)
		}
	}

	return occupied
}

func (c *Controller) UnOccupiedElevators() []*Elevator {
	var unOccupied []*Elevator

	for _, val := range c.Elevators {
		if !val.InUse {
			unOccupied = append(unOccupied, val)
		}
	}

	return unOccupied
}

func (c *Controller) PendingRequests() []*Request {
	return c.requestQ
}
