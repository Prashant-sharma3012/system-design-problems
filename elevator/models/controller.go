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
	requestQ       chan *Request
}

func (c *Controller) init() {
	// if there are total 20 floors wach floor can make 2 requests
	c.requestQ = make(chan *Request, len(c.FloorSwitches)*2)
}

// helper fucntions
func setDirection(val *Elevator, up bool, down bool) {
	if up {
		fmt.Println("Going Up")
		val.GoingUp = true
		val.InUse = true
	}

	if down {
		fmt.Println("Going Down")
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

func GetController(numOfLifts int, topFloor int) (*Controller, error) {

	if numOfLifts == 0 || topFloor == 0 {
		return nil, errors.New("Please provide number of lifts and total number of floors")
	}

	var elevators []*Elevator
	var elevator *Elevator

	for i := 1; i <= numOfLifts; i++ {
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

	var callFrom = -1
	var minDiff, tempMinDiff int

	minDiff = len(c.FloorSwitches)

	// keep running until you find a lift
	for callFrom == -1 {
		// check if there is a lift already at this floor
		for indx, val := range c.Elevators {

			if !val.InUse {
				if val.CurrentPosition == atFloor {
					fmt.Println("Servicing Request for floor")
					fmt.Println(atFloor)
					fmt.Println("Using Lift")
					fmt.Println(val.Id)
					val.InUse = true
					setDirection(val, up, down)
					return val
				}

				tempMinDiff = val.CurrentPosition - atFloor
				if tempMinDiff < 0 {
					tempMinDiff = tempMinDiff * -1
				}

				if tempMinDiff < minDiff {
					minDiff = tempMinDiff
					callFrom = indx
				}
			}
		}
	}

	fmt.Println("Servicing Request for floor")
	fmt.Println(atFloor)
	fmt.Println("from")
	fmt.Println(c.Elevators[callFrom].CurrentPosition)
	fmt.Println("Using Lift")
	fmt.Println(c.Elevators[callFrom].Id)

	c.Elevators[callFrom].InUse = true
	setDirection(c.Elevators[callFrom], up, down)

	// set the switch to off again
	resetSwitch(c, atFloor, up, down)

	return c.Elevators[callFrom]
}

func (c *Controller) StartServicing() {
	c.start = true
	c.init()

	for c.start {
		for req := range c.requestQ {
			go c.Call(req.Floor, req.GoingUp, req.GoingDown).GoTo(req.Floor)
		}
	}
}

func (c *Controller) StopServicing() {
	c.start = false
}

//will replce this
func (c *Controller) RequestFromFloor(floor int, up bool, down bool) {

	currentReq := len(c.requestQ)
	req := createRequest(floor, up, down, currentReq)

	if up {
		c.FloorSwitches[floor-1].requestUp = currentReq
		c.FloorSwitches[floor-1].GoUp()
	}

	if down {
		c.FloorSwitches[floor-1].requestDown = currentReq
		c.FloorSwitches[floor-1].GoDown()
	}

	c.requestQ <- req
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

func (c *Controller) PendingRequests() int {
	return len(c.requestQ)
}
