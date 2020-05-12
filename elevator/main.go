package main

import (
	"fmt"
	"sync"

	model "github.com/system-design-problems/elevator/models"
)

func main() {
	// Do a dry run of things justtosee if everything is working

	var wg sync.WaitGroup

	controller, _ := model.GetController(2, 20)

	fmt.Println("Starting Service")
	go controller.StartServicing(&wg)

	fmt.Println("Creating Requests")

	controller.RequestFromFloor(10, false, true)
	controller.RequestFromFloor(2, true, false)
	controller.RequestFromFloor(8, false, true)
	controller.RequestFromFloor(6, true, false)
	controller.RequestFromFloor(5, false, true)
	controller.RequestFromFloor(12, true, false)

	fmt.Println("Scanning For reqs")

	wg.Wait()

	fmt.Println("No More requests, Bye")
}
