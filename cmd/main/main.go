//go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest

package main

import "browserGui/internal/server"

func main() {
	server.NewServer()
}
