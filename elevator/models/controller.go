package model

import (
	"errors"
	"fmt"
	"sort"
	"sync"
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

var m sync.Mutex

func (c *Controller) init() {
	// if there are total 20 floors wach floor can make 2 requests
	c.requestQ = make(chan *Request, len(c.FloorSwitches)*2)
}

// helper fucntions
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

func updateElevator(val *Elevator, up bool, down bool, atFloor int) bool {

	val.Lock()
	defer val.Unlock()

	if val.InUse {
		return false
	}

	val.PickFromFloor = append(val.PickFromFloor, atFloor)

	if up {
		val.GoingUp = true
		val.InUse = true
		val.StopAtFloor = append(val.StopAtFloor, 15)
	} else {
		val.GoingDown = true
		val.InUse = true
		val.StopAtFloor = append(val.StopAtFloor, 1)
	}

	return true
}

func log(floor, id, pos int) {
	fmt.Printf("Servicing Request for floor %d , Using Lift %d at Floor %d \n",
		floor,
		id,
		pos)
}

func (c *Controller) Call(atFloor int, up bool, down bool, wg *sync.WaitGroup) {
	var callFrom = -1
	var minDiff, tempMinDiff int

	minDiff = len(c.FloorSwitches)

	// keep running until you find a lift
	for callFrom == -1 {
		// check if there is a lift already at this floor
		for indx, val := range c.Elevators {
			if !val.InUse {
				if val.CurrentPosition == atFloor {
					if updateElevator(val, up, down, atFloor) {
						log(atFloor, val.Id, val.CurrentPosition)
						resetSwitch(c, atFloor, up, down)
						val.ServeReqs(wg, &m)
						return
					}
				}

				tempMinDiff = val.CurrentPosition - atFloor
				if tempMinDiff < 0 {
					tempMinDiff = tempMinDiff * -1
				}

				if tempMinDiff < minDiff {
					minDiff = tempMinDiff
					callFrom = indx
				}

			} else {
				if val.GoingUp && val.CurrentPosition < (atFloor-2) {
					val.Lock()
					val.PickFromFloor = append(val.PickFromFloor, atFloor)
					sort.Slice(val.PickFromFloor, func(i, j int) bool { return val.PickFromFloor[i] < val.PickFromFloor[j] })
					val.Unlock()

					resetSwitch(c, atFloor, up, down)
					fmt.Printf("Using Elevator %d already servicing \n", val.Id)
					log(atFloor, val.Id, val.CurrentPosition)
					return
				}

				if val.GoingDown && val.CurrentPosition > (atFloor+2) {
					val.Lock()
					val.PickFromFloor = append(val.PickFromFloor, atFloor)
					sort.Slice(val.PickFromFloor, func(i, j int) bool { return val.PickFromFloor[i] > val.PickFromFloor[j] })
					val.Unlock()

					resetSwitch(c, atFloor, up, down)
					fmt.Printf("Using Elevator %d already servicing \n", val.Id)
					log(atFloor, val.Id, val.CurrentPosition)
					return
				}
			}
		}
	}

	if updateElevator(c.Elevators[callFrom], up, down, atFloor) {
		log(atFloor, c.Elevators[callFrom].Id, c.Elevators[callFrom].CurrentPosition)
		resetSwitch(c, atFloor, up, down)
		c.Elevators[callFrom].ServeReqs(wg, &m)
		return
	} else {
		c.Call(atFloor, up, down, wg)
	}
}

func (c *Controller) StartServicing(wg *sync.WaitGroup) {
	c.start = true
	c.init()

	for c.start {
		for req := range c.requestQ {
			wg.Add(1)
			go c.Call(req.Floor, req.GoingUp, req.GoingDown, wg)
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
