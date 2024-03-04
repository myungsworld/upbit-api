package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"upbit-api/config"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	oneChan := make(chan int, 1)
	block := make(chan bool)

	var a int
	go func() {
		for {
			time.Sleep(time.Second)
			a++
			fmt.Println("a 값:", a)
			go func() {

				oneChan <- a

			}()
		}
	}()

	go func() {
		for {
			select {
			case val, ok := <-oneChan:
				block <- true
				fmt.Println("뽑아온 값:", val, ok)
				time.Sleep(time.Second * 5)

			}
		}
	}()

	<-stopChan

}

// .env 로드 , Market 상태 수집
func init() {
	config.Init()
}
