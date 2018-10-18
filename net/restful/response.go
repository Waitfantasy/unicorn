package restful

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Message string
	Status  bool
	Data    map[string]interface{}
	w       *http.ResponseWriter
}

func (r *response) success() {
	(*r.w).WriteHeader(http.StatusOK)
	r.json()
	fmt.Fprint(*r.w)
}

func (r *response) error(code int) {
	(*r.w).WriteHeader(code)
	r.json()
	fmt.Fprint(*r.w)
}

func (r *response) json() {
	(*r.w).Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(r)
	(*r.w).Write(b)
}
