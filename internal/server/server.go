package server

import (
	"browserGui/internal/api"
	"browserGui/internal/programs"
	"net/http"
	"time"
)

func NewServer() {
	fs := http.FileServer(http.Dir("./static"))

	http.HandleFunc("/api/keepalive", KeepAlive)
	http.Handle("/", fs)

	api.NewPing().Register()
	api.NewCalculate().Register()

	waitClose()

	programs.OpenBrowserDelay("http://localhost:3000", time.Millisecond * 15)

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		panic(err)
	}
}
