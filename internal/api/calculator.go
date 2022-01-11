package api

import (
	"browserGui/internal/repository"
	"crypto/sha1"
	"encoding/hex"
	json2 "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type calculator struct {

}

type Question struct {
	Text string `json:"question_text"`
}

type Answer struct {
	Answer string `json:"answer"`
}

type CalcAnswer struct {
	Message string `json:"message"`
	Success bool `json:"success"`
	Code string `json:"code"`
	Progress string `json:"progress"`
}
func CalcError(code, message, progress string) CalcAnswer {
	return CalcAnswer{
		Message: message,
		Code: code,
		Progress: progress,
		Success: false,
	}
}
func CalcSuccess(code, message, progress string) CalcAnswer {
	return CalcAnswer{
		Message: message,
		Code: code,
		Progress: progress,
		Success: true,
	}
}

func NewCalculator() *calculator {
	return &calculator{}
}
func Sha1String(s string) string {
	bytes := []byte(s)
	bytes = append(bytes, repository.GetTokenSecretKey()...)

	h := sha1.New()
	h.Write(bytes)
	return hex.EncodeToString(h.Sum(nil))
}

func (c *calculator) Register(r *mux.Router) {
	r.HandleFunc("/calculator/getQuestion", c.getQuestion)
	r.HandleFunc("/calculator/checkAnswer", c.checkAnswer)
}

func (c *calculator) getQuestion(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())

	var a int
	var b int
	var answer string
	var text string

	if rand.Intn(2) == 1 {
		a = rand.Intn(99) + 1
		b = rand.Intn(99) + 1

		text = strconv.Itoa(a) + "+" + strconv.Itoa(b)
		answer = strconv.Itoa(a + b)
	} else {
		a = rand.Intn(95) + 5
		b = rand.Intn(a - 1) + 1

		text = strconv.Itoa(a) + "-" + strconv.Itoa(b)
		answer = strconv.Itoa(a - b)
	}

	repo := repository.GetRepository(r.Context())

	if _, exists := repo.Get("failProgress"); ! exists {
		repo.Set("failProgress", "50")
	}

	repo.Set("calc_answer", Sha1String(answer))

	q := Question{
		Text: text,
	}

	json, err := json2.Marshal(q)

	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(json))
}

func (c *calculator) checkAnswer(w http.ResponseWriter, r *http.Request) {
	var answer Answer

	repo := repository.GetRepository(r.Context())

	stringProgress, _ := repo.Get("failProgress")
	failProgress, _ := strconv.Atoi(stringProgress)

	if err := json2.NewDecoder(r.Body).Decode(&answer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json2.NewEncoder(w).Encode(CalcError("cant_check_answer", "Не удается проверить ответ", stringProgress))
		return
	}

	realAnswer, _ := repository.GetRepository(r.Context()).Get("calc_answer")

	if Sha1String(answer.Answer) == realAnswer {
		progress := failProgress - 2

		code := "success"
		if progress <= 2 {
			progress = 2
			code = "you_win"

			appdata, error := os.UserConfigDir()
			if error != nil {
				panic(error)
			}

			cmd := exec.Command(appdata + "\\.minecraft\\TL.exe")
			if err := cmd.Start(); err != nil {
				panic(err)
			}
		}

		stringProgress = strconv.Itoa(progress)
		repo.Set("failProgress", stringProgress)

		w.WriteHeader(http.StatusOK)
		json2.NewEncoder(w).Encode(CalcSuccess(code, "", stringProgress))
	} else {
		progress := failProgress + 10

		code := "incorrect_answer"
		if progress >= 84 {
			progress = 84
			code = "you_lost"
		}

		stringProgress = strconv.Itoa(progress)
		repo.Set("failProgress", stringProgress)

		w.WriteHeader(http.StatusBadRequest)
		json2.NewEncoder(w).Encode(CalcError(code, "Не верный ответ", stringProgress))
	}
}