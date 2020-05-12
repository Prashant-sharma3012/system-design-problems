package model

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var currentID = 1 // can use iota here

type Elevator struct {
	sync.Mutex
	Id              int
	CurrentPosition int
	TopFloor        int
	GoingUp         bool
	GoingDown       bool
	InUse           bool
	PickFromFloor   []int
	StopAtFloor     []int
}

func (e *Elevator) ServeReqs(wg *sync.WaitGroup, m *sync.Mutex) {
	finalPosition := 0

	for len(e.PickFromFloor) > 0 || len(e.StopAtFloor) > 0 {
		time.Sleep(999999999)

		e.Lock()
		if e.GoingDown {
			e.CurrentPosition--
		} else {
			e.CurrentPosition++
		}

		if len(e.StopAtFloor) == 1 {
			finalPosition = e.StopAtFloor[0]
		}

		if e.CurrentPosition == e.PickFromFloor[0] {
			fmt.Printf("Stopping at floor %d to pickup ", e.CurrentPosition)
			e.PickFromFloor = e.PickFromFloor[1:]
			wg.Done()
		}

		if e.CurrentPosition == e.StopAtFloor[0] {
			fmt.Printf("Stopping at floor %d to drop ", e.CurrentPosition)
			e.StopAtFloor = e.StopAtFloor[1:]
		}
		e.Unlock()
	}

	e.Lock()
	e.CurrentPosition = finalPosition
	e.InUse = false
	e.GoingDown = false
	e.GoingDown = false
	e.Unlock()
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
