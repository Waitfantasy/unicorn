package restful

import (
	"errors"
	"fmt"
	"github.com/Waitfantasy/unicorn/id"
	"github.com/Waitfantasy/unicorn/service/machine"
	"github.com/Waitfantasy/unicorn/util/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type api struct {
	g *id.AtomicGenerator
	m machine.Machiner
	l *logger.Log
}

func (a *api) register() *gin.Engine {
	e := gin.Default()
	e.Use(a.l.GinMiddleware())
	// uuid api group
	g1 := e.Group("/api/v1/uuid")
	g1.GET("/make", a.uuidMake())
	g1.GET("/transfer/:uuid", a.uuidTransfer())

	// machine api group
	g2 := e.Group("/api/v1/machine")
	g2.GET("/list", a.machineList())
	g2.POST("/store", a.machineList())
	g2.POST("/delete", a.machineDelete())
	g2.POST("/replace", a.machineReplace())
	return e
}

func (a *api) uuidMake() gin.HandlerFunc {
	return func(c *gin.Context) {
		if uuid, err := a.g.Make(); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "make uuid fail.",
				"data": map[string]interface{}{
					"error": err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "make uuid success.",
				"data": map[string]interface{}{
					"uuid": uuid,
				},
			})
		}
	}
}

func (a *api) uuidTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if v := c.Param("uuid"); v == "" {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "missing uuid parameter.",
			})
		} else {
			if uuid, err := strconv.ParseUint(v, 10, 64); err != nil {
				c.JSON(http.StatusOK, map[string]interface{}{
					"success": false,
					"message": fmt.Sprintf("the uuid %s convert uint64 fail", v),
					"data": map[string]interface{}{
						"error": err.Error(),
					},
				})
			} else {
				data := a.g.Extract(uuid)
				c.JSON(http.StatusOK, map[string]interface{}{
					"success": true,
					"message": "transfer uuid success.",
					"data": map[string]interface{}{
						"machine_id": data.MachineId,
						"seq":        data.Sequence,
						"timestamp":  data.Timestamp,
						"reserved":   data.Reserved,
						"id_type":    data.IdType,
						"version":    data.Version,
					},
				})
			}
		}
	}
}

func (a *api) machineList() gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := a.m.All()
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "get all machine item fail.",
				"data": map[string]interface{}{
					"error": err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": "get all machine item success.",
				"data": map[string]interface{}{
					"machines": items,
				},
			})
		}
	}
}

func (a *api) machineStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, err := validatorIp(c)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		if item, err := a.m.Put(ip); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("put ip: %s fail.\n", ip),
				"data": map[string]interface{}{
					"error": err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("put ip: %s success.\n", ip),
				"data": map[string]interface{}{
					"machine": item,
				},
			})
		}
	}
}

func (a *api) machineDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, err := validatorIp(c)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		if item, err := a.m.Del(ip); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": fmt.Sprintf("delete ip: %s fail.\n", ip),
				"data": map[string]interface{}{
					"error": err.Error(),
				},
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("delete ip: %s success.\n", ip),
				"data": map[string]interface{}{
					"machine": item,
				},
			})
		}
	}
}

func (a *api) machineReplace() gin.HandlerFunc {
	return func(c *gin.Context) {
		inputs := make(map[string]string)
		if err := c.BindJSON(&inputs); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		oldIp, ok := inputs["oldIp"]
		if !ok {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "missing oldIp parameter",
			})
			return
		}

		newIp, ok := inputs["newIp"]
		if !ok {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "missing newIp parameter",
			})
			return
		}

		if err := a.m.Reset(oldIp, newIp); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"message": fmt.Sprintf("%s has been replaced by %s", oldIp, newIp),
			})
		}
	}
}

func validatorIp(c *gin.Context) (string, error) {
	inputs := make(map[string]string)
	if err := c.BindJSON(&inputs); err != nil {
		return "", err
	} else {
		if ip, ok := inputs["ip"]; !ok {
			return "", errors.New("missing ip parameters")
		} else {
			return ip, nil
		}
	}
}
