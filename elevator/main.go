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

	for i := 2; i < 10; i++ {
		controller.RequestFromFloor(10-i, false, true)
		controller.RequestFromFloor(i, true, false)
	}

	fmt.Println("Scanning For reqs")

	wg.Wait()

	fmt.Println("No More requests, Bye")
}
