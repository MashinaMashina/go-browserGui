package api

import (
	"crypto/sha1"
	"encoding/hex"
	json2 "encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"browserGui/internal/repository"
	"github.com/gorilla/mux"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type calculator struct {
	password string
}

type Question struct {
	Text string `json:"question_text"`
}

type Answer struct {
	Answer string `json:"answer"`
}

type CalcAnswer struct {
	Message  string `json:"message"`
	Success  bool   `json:"success"`
	Code     string `json:"code"`
	Progress string `json:"progress"`
}

func CalcError(code, message, progress string) CalcAnswer {
	return CalcAnswer{
		Message:  message,
		Code:     code,
		Progress: progress,
		Success:  false,
	}
}
func CalcSuccess(code, message, progress string) CalcAnswer {
	return CalcAnswer{
		Message:  message,
		Code:     code,
		Progress: progress,
		Success:  true,
	}
}

func NewCalculator() *calculator {
	password := os.Getenv("BGUI_PASSWORD")
	if password == "" {
		password = "password"
	}

	return &calculator{
		password: password,
	}
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
	var a int
	var b int
	var answer string
	var text string

	switch rand.Intn(4) {
	case 3:
		answerInt := rand.Intn(27) + 2
		b = rand.Intn(7) + 2
		a = b * answerInt

		text = strconv.Itoa(a) + "/" + strconv.Itoa(b)
		answer = strconv.Itoa(answerInt)
	case 2:
		a = rand.Intn(97) + 2
		b = rand.Intn(7) + 2

		text = strconv.Itoa(a) + "*" + strconv.Itoa(b)
		answer = strconv.Itoa(a * b)
	case 1:
		a = rand.Intn(995) + 5
		b = rand.Intn(a-1) + 1

		text = strconv.Itoa(a) + "-" + strconv.Itoa(b)
		answer = strconv.Itoa(a - b)
	default: // case 0
		a = rand.Intn(500) + 1
		b = rand.Intn(500) + 1

		text = strconv.Itoa(a) + "+" + strconv.Itoa(b)
		answer = strconv.Itoa(a + b)
	}

	repo := repository.GetRepository(r.Context())

	if _, exists := repo.Get("failProgress"); !exists {
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

	if answer.Answer == c.password {
		repo.Set("failProgress", "0")
		w.WriteHeader(http.StatusOK)
		json2.NewEncoder(w).Encode(CalcSuccess("you_win", "", stringProgress))
		return
	}

	if Sha1String(answer.Answer) == realAnswer {
		progress := failProgress - 3

		code := "success"
		if progress <= 2 {
			progress = 2
			code = "you_win"
		}

		stringProgress = strconv.Itoa(progress)
		repo.Set("failProgress", stringProgress)

		w.WriteHeader(http.StatusOK)
		json2.NewEncoder(w).Encode(CalcSuccess(code, "", stringProgress))
	} else {
		progress := failProgress + 3

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
