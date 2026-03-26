package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rdmrcv/repartnersgo/app/service"
)

// Request model to validate calculation.
type Request struct {
	Order int `json:"order" binding:"required,gt=0,lte=10000000" minimum:"1" maximum:"10000000"`

	Packages []int `json:"packages" binding:"required,dive,gt=0,lte=1000000"`
}

// Response model to have typed output.
type Response struct {
	Packs map[int]int `json:"packs"`
}

type ResponseErr struct {
	Error string `json:"error"`
}

// Solve handles the HTTP endpoint for the [service.Solve]. It validates input
// and responds properly.
//
// @Summary Receive config for packing and respond accordingly.
// @Description Receive order and packs variants and responds with optimal packing.
// @Accept json
// @Produce json
// @Param request body Request true "Request for task solution"
// @Success 200 {object} Response
// @Failure 400 {object} ResponseErr
// @Failure 500 {object} ResponseErr
// @Router /solve [post]
func Solve(c *gin.Context) {
	var req Request

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ResponseErr{Error: err.Error()})

		return
	}

	resp, err := service.Solve(req.Order, req.Packages)
	switch {
	case errors.Is(err, service.ErrParams):
		c.JSON(http.StatusBadRequest, ResponseErr{Error: err.Error()})
	case err != nil:
		c.JSON(http.StatusInternalServerError, ResponseErr{Error: err.Error()})
	default:
		c.JSON(http.StatusOK, Response{Packs: resp})
	}
}
