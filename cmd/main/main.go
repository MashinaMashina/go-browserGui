//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest

package main

import (
	"fmt"

	"browserGui/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("err load .env:", err)
	}

	server.NewServer()
}
