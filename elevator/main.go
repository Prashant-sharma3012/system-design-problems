package main

import (
	"fmt"
	"time"

	model "github.com/system-design-problems/elevator/models"
)

func main() {
	// Do a dry run of things justtosee if everything is working

	controller, _ := model.GetController(2, 20)

	fmt.Println("Starting Service")
	go controller.StartServicing()

	fmt.Println("Request 10 Elevators")

	for i := 2; i < 10; i++ {
		go controller.RequestFromFloor(i, true, false)
	}

	fmt.Println("Scanning For reqs")

	for controller.PendingRequests() > 0 {
		time.Sleep(9999999999)
		fmt.Println("Scanning For reqs")
	}

	fmt.Println("No More requests, Bye")
}
