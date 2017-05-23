package api

import (
	"github.com/zyfdegh/fanach/coreserver/entity"
	"github.com/zyfdegh/fanach/coreserver/service"
	"gopkg.in/kataras/iris.v6"
)

// CreateOrder handles POST /orders
func CreateOrder(ctx *iris.Context) {
	resp := entity.RespPostOrder{}

	req := &entity.ReqPostOrder{}

	err := ctx.ReadJSON(req)
	if err != nil {
		resp.ErrNo = iris.StatusBadRequest
		resp.Errmsg = err.Error()
		ctx.JSON(resp.ErrNo, resp)
		return
	}

	order, err := service.CreateOrder(*req)
	if err != nil {
		resp.ErrNo = iris.StatusInternalServerError
		resp.Errmsg = err.Error()
		ctx.JSON(resp.ErrNo, resp)
		return
	}

	resp.Success = true
	if order != nil {
		resp.Order = *order
	}
	ctx.JSON(iris.StatusOK, resp)
	return
}
