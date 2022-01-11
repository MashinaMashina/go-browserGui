package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type ping struct {

}

func NewPing() *ping {
	return &ping{}
}

func (p *ping) Register(r *mux.Router) {
	r.HandleFunc("/ping", p.ping)
}

func (p *ping) ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Its work!")
}
