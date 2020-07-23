package main

import (
	"fmt"
	"os"
	"os/signal"
)

var (
	hydrogenSignal = make(chan struct{}, 1)
	oxygenSignal   = make(chan struct{}, 2)
)

func releaser(output string, max int) func() {
	counter := 0
	return func() {
		if counter >= max {
			panic(fmt.Sprintf("No More '%v'", output))
		}
		fmt.Print(output)
		counter++
	}
}

func hydrogenGoroutine(reliser func()) func() {
	return func() {
		defer func() {
			if rcv := recover(); rcv != nil {
				fmt.Print("\n", rcv)
			}
		}()

		for {
			oxygenSignal <- struct{}{}

			<-hydrogenSignal
			reliser()
		}
	}
}

func oxygenGoroutine(reliser func()) func() {
	return func() {
		defer func() {
			if rcv := recover(); rcv != nil {
				fmt.Print("\n", rcv)
			}
		}()

		for {
			hydrogenSignal <- struct{}{}
			hydrogenSignal <- struct{}{}

			<-oxygenSignal
			<-oxygenSignal
			reliser()
		}
	}
}

func main() {
	fmt.Println("Hello, Welcome Back!")

	maxHydrogen, maxOxygen := 20, 10
	hydrogenReleaser := releaser("H", maxHydrogen)
	oxygenReleaser := releaser("O", maxOxygen)

	go hydrogenGoroutine(hydrogenReleaser)()
	go hydrogenGoroutine(hydrogenReleaser)()
	go oxygenGoroutine(oxygenReleaser)()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
