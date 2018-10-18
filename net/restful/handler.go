package restful

import (
	"github.com/Waitfantasy/unicorn/id"
	"net/http"
)

type handlers struct {
	generator *id.AtomicGenerator
}

func (h *handlers) register() {
	http.Handle("/uuid", h.uuid())
}

func (h *handlers) uuid() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		res := response{
			w: &writer,
		}

		if request.Method != http.MethodGet {
			res.Message = "Please use GET method"
			res.Status = false
			res.success()
			return
		}

		uuid, err := h.generator.Make()
		if err != nil {
			res.Message = "Make uuid fail."
			res.Status = false
			res.Data = map[string]interface{}{
				"error": err.Error(),
			}
			res.success()
			return
		}

		res.Message = "Make uuid success."
		res.Status = true
		res.Data = map[string]interface{}{
			"uuid": uuid,
		}
		res.success()
	}
}
