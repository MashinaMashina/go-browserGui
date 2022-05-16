package server

import (
	"browserGui/internal/api"
	"browserGui/internal/programs"
	"browserGui/internal/repository"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func NewServer() {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api/").Subrouter()

	apiRouter.Use(repository.TokenMiddleware)
	api.NewPing().Register(apiRouter)
	api.NewCalculator().Register(apiRouter)
	api.NewOpener().Register(apiRouter)

	router.HandleFunc("/keepalive/", KeepAlive)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", router)

	waitClose()

	port := "3001"

	programs.OpenBrowserDelay("http://localhost:" + port, time.Millisecond * 15)

	err := http.ListenAndServe(":" + port, nil)

	if err != nil {
		panic(err)
	}
}
