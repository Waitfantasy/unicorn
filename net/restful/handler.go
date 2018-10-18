package restful

import (
	"fmt"
	"github.com/Waitfantasy/unicorn/id"
	"net/http"
	"strconv"
)

type handlers struct {
	generator *id.AtomicGenerator
}

func (h *handlers) register() {
	http.Handle("/uuid/make", h.uuid())
	http.Handle("/uuid/extract", h.extract())
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

func (h *handlers) extract() http.HandlerFunc {
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

		if uuid := request.URL.Query().Get("uuid"); uuid == "" {
			res.Message = "missing uuid param"
			res.Status = false
			res.success()
			return
		} else {
			if uuidUint64, err := strconv.ParseUint(uuid, 10, 64); err != nil {
				res.Message = fmt.Sprintf("the uuid %s convert uint64 fail", uuid)
				res.Status = false
				res.Data = map[string]interface{}{
					"error": err.Error(),
				}
				res.success()
				return
			} else {
				data := h.generator.Extract(uuidUint64)
				res.Message = "extract uuid success"
				res.Status = true
				res.Data = map[string]interface{}{
					"machine_id": data.MachineId,
					"seq":        data.Sequence,
					"timestamp":  data.Timestamp,
					"reserved":   data.Reserved,
					"id_type":    data.IdType,
					"version":    data.Version,
				}
				res.success()
				return
			}
		}
	}
}
