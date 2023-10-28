package server

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"browserGui/internal/api"
	"browserGui/internal/programs"
	"browserGui/internal/repository"
	"github.com/gorilla/mux"
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

	port := os.Getenv("BGUI_PORT")
	if port == "" {
		port = strconv.Itoa(rand.Intn(100) + 3000)
	}

	programs.OpenBrowserDelay("http://localhost:"+port, time.Millisecond*15)

	err := http.ListenAndServe("localhost:"+port, nil)

	if err != nil {
		panic(err)
	}
}
