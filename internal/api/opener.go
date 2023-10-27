package api

import (
	json2 "encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"browserGui/internal/repository"
	"github.com/gorilla/mux"
)

type GameSelect struct {
	Game string `json:"game"`
}

type AvailableProgram struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type opener struct {
	programs []AvailableProgram
}

type OpenerAnswer struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewOpener() *opener {
	programsStr := os.Getenv("BGUI_PROGRAMS")

	if programsStr == "" {
		programsStr = `[{"name":"Minecraft","path":"%appdata%\\.tlauncher\\legacy\\Minecraft\\TL.exe"}]`
	}

	var programs []AvailableProgram
	if err := json2.Unmarshal([]byte(programsStr), &programs); err != nil {
		panic(fmt.Errorf("can not parse programs str: %w", err))
	}

	appdata, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	availablePrograms := make([]AvailableProgram, 0, len(programs))
	for _, program := range programs {
		program.Path = strings.ReplaceAll(program.Path, "%appdata%", appdata)

		if _, err = os.Stat(program.Path); err == nil {
			availablePrograms = append(availablePrograms, program)
		}
	}

	if len(availablePrograms) == 0 {
		panic("not found available programs")
	}

	return &opener{
		programs: availablePrograms,
	}
}

func (c *opener) Register(r *mux.Router) {
	r.HandleFunc("/opener/open", c.Open)
	r.HandleFunc("/opener/available", c.Available)
}

func (c *opener) Available(w http.ResponseWriter, r *http.Request) {
	json2.NewEncoder(w).Encode(c.programs)
}

func (c *opener) Open(w http.ResponseWriter, r *http.Request) {
	repo := repository.GetRepository(r.Context())

	stringProgress, _ := repo.Get("failProgress")
	progress, _ := strconv.Atoi(stringProgress)

	if progress <= 2 {
		var game GameSelect

		if err := json2.NewDecoder(r.Body).Decode(&game); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json2.NewEncoder(w).Encode(OpenerAnswer{Success: false, Message: "Не удается определить игру"})
			return
		}

		for _, program := range c.programs {
			if program.Name == game.Game {
				json2.NewEncoder(w).Encode(c.runApp(program.Path))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json2.NewEncoder(w).Encode(OpenerAnswer{Success: false, Message: "Ошибка. Не найдено приложение " + game.Game})
		return
	}

	w.WriteHeader(http.StatusOK)
	json2.NewEncoder(w).Encode(OpenerAnswer{Success: false, Message: "Прогресс слишком маленький"})
}

func (c *opener) runApp(p string) OpenerAnswer {
	p = strings.ReplaceAll(p, "\\", "/")

	cmd := exec.Command(p)
	cmd.Dir, _ = path.Split(p)
	if err := cmd.Start(); err != nil {
		return OpenerAnswer{Success: false, Message: "Не удается запустить " + p}
	}

	return OpenerAnswer{Success: true, Message: "Успешно"}
}
