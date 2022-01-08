package server

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var continueWork chan int

func waitClose() {
	continueWork = make(chan int)

	go func() {
		quitTimer := time.NewTimer(time.Second * 11)

		for {
			select {
			case <-continueWork:
				quitTimer.Stop()
				quitTimer = time.NewTimer(time.Second * 10)
			case <-quitTimer.C:
				os.Exit(0)
			}
		}
	}()
}

func ContinueWork() {
	continueWork <- 1
}

func KeepAlive(w http.ResponseWriter, r *http.Request) {
	ContinueWork()

	fmt.Fprint(w, "OK")
}