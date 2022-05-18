package api

import (
	"browserGui/internal/repository"
	json2 "encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type GameSelect struct {
	Game string `json:"game"`
}

type opener struct {

}

type OpenerAnswer struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

func NewOpener() *opener {
	return &opener{}
}

func (c *opener) Register(r *mux.Router) {
	r.HandleFunc("/opener/open", c.Open)
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

		switch game.Game {
		case "battlefield":
			p := "F:/Games/Battlefield2/bf2.exe"

			json2.NewEncoder(w).Encode(c.runApp(p))
			return
		case "minecraft":
			appdata, err := os.UserConfigDir()
			if err != nil {
				panic(err)
			}

			json2.NewEncoder(w).Encode(c.runApp(appdata + "\\.minecraft\\TL.exe"))
			return
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