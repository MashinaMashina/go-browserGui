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

			cmd := exec.Command(p)
			cmd.Dir, _ = path.Split(p)
			if err := cmd.Start(); err != nil {
				json2.NewEncoder(w).Encode(OpenerAnswer{Success: false, Message: "Не удается запустить " + p})
				return
			}

			json2.NewEncoder(w).Encode(OpenerAnswer{Success: true, Message: "Успешно"})
			return
		case "minecraft":
			appdata, err := os.UserConfigDir()
			if err != nil {
				panic(err)
			}

			p := appdata + "\\.minecraft\\TL.exe"
			cmd := exec.Command(p)
			if err = cmd.Start(); err != nil {
				json2.NewEncoder(w).Encode(OpenerAnswer{Success: false, Message: "Не удается запустить " + p})
				return
			}

			json2.NewEncoder(w).Encode(OpenerAnswer{Success: true, Message: "Успешно"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json2.NewEncoder(w).Encode(OpenerAnswer{Success: false, Message: ""})
		return
	}

	w.WriteHeader(http.StatusOK)
	json2.NewEncoder(w).Encode(OpenerAnswer{Success: false, Message: "Прогресс слишком маленький"})
}