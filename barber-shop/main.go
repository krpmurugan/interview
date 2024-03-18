package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wakeBell chan bool
var sleepBell chan bool

func SleepBarber(name string) {
	sleepBell <- true
	fmt.Println("=========== Barber went for sleep")
}

func WakeBarber() {
	fmt.Println("=========== Barber Wakeup")
}

func BarberCutting(name string) {
	for i := 0; i < 5; i++ {
		fmt.Println("=========== Cutting in progress for ", name)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("=========== Cutting completed for ", name)
}

func BarberShop(chairs chan string, ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	sleepBell <- true
	defer wg.Done()
	fmt.Println("=========== Barber Shop Opened")
	defer fmt.Println("=========== Barber Shop Closed")
	closeShop := false
	for {
		select {
		case name := <-chairs:
			BarberCutting(name)
			if len(chairs) == 0 {
				if !closeShop {
					SleepBarber(name)
				} else {
					return
				}
			}
		case <-wakeBell:
			WakeBarber()
		case <-ctx.Done():
			if len(chairs) == 0 {
				return
			}
			closeShop = true
		}
	}
}
func CustomerEntry(name string, chairs chan string) {
	fmt.Printf("=========== Customer %s entered Barber shop\n", name)
	select {
	case <-sleepBell:
		wakeBell <- true
		chairs <- name
		fmt.Printf("=========== Got a chair for %s\n", name)
	case chairs <- name:
		fmt.Printf("=========== Got a chair for %s\n", name)
	default:
		fmt.Printf("=========== No chair for %s, hence leaving Barber shop\n", name)
	}
}

func main() {
	wakeBell = make(chan bool)
	sleepBell = make(chan bool, 1)
	chairs := make(chan string, 2)
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	go BarberShop(chairs, ctx, wg)
	<-sleepBell
	CustomerEntry("Murugan", chairs)
	CustomerEntry("Rajeev", chairs)
	CustomerEntry("Praveen", chairs)
	CustomerEntry("James", chairs)
	time.Sleep(20 * time.Second)
	CustomerEntry("David", chairs)
	cancel()
	wg.Wait()
	close(wakeBell)
	close(chairs)
	close(sleepBell)
	fmt.Println("=========== Exiting...")
}
