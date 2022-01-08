package api

import (
	"fmt"
	"net/http"
)

type ping struct {

}

func NewPing() *ping {
	return &ping{}
}

func (p *ping) Register() {
	http.HandleFunc("/api/ping", p.ping)
}

func (p *ping) ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Its work!")
}
