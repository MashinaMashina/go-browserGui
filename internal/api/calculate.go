package api

import (
	"browserGui/internal/crypt"
	json2 "encoding/json"
	"fmt"
	"net/http"
)

type calculator struct {

}

type Question struct {
	QuestionText string `json:"question_text"`
	SecretAnswer string `json:"secret_answer"`
}

func NewCalculate() *calculator {
	return &calculator{}
}

func (c *calculator) Register() {
	http.HandleFunc("/api/calculator/getQuestion", c.getQuestion)
	http.HandleFunc("/api/calculator/checkAnswer", c.checkAnswer)
}

func (c *calculator) getQuestion(w http.ResponseWriter, r *http.Request) {
	text := "1+2"
	secret, err := crypt.Encrypt([]byte("3"))

	if err != nil {
		panic(err)
	}

	q := Question{
		QuestionText: text,
		SecretAnswer: secret,
	}

	json, err := json2.Marshal(q)

	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(json))
}

func (c *calculator) checkAnswer(w http.ResponseWriter, r *http.Request) {

}