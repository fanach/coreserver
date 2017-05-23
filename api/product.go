package api

import (
	"github.com/zyfdegh/fanach/coreserver/entity"
	"github.com/zyfdegh/fanach/coreserver/service"
	"gopkg.in/kataras/iris.v6"
)

// GetProducts handles GET /prods
func GetProducts(ctx *iris.Context) {
	resp := entity.RespGetProducts{}

	products, err := service.GetProducts()
	if err != nil {
		resp.ErrNo = iris.StatusInternalServerError
		resp.Errmsg = err.Error()
		ctx.JSON(resp.ErrNo, resp)
		return
	}

	resp.Success = true
	if products != nil {
		resp.Products = *products
	}
	ctx.JSON(iris.StatusOK, resp)
	return
}
