package api

import (
	"github.com/zyfdegh/fanach/coreserver/entity"
	"github.com/zyfdegh/fanach/coreserver/service"
	"gopkg.in/kataras/iris.v6"
)

const (
	// KeySessID is the key of session id
	KeySessID = "sessid"
	// KeySessUsername is the key of session username
	KeySessUsername = "username"
)

// PostSession create session for a user
func PostSession(ctx *iris.Context) {
	resp := entity.RespPostSess{}

	req := &entity.ReqPostSess{}
	err := ctx.ReadJSON(req)
	if err != nil {
		resp.Errmsg = err.Error()
		resp.ErrNo = iris.StatusBadRequest
		ctx.JSON(resp.ErrNo, resp)
		return
	}

	sess, err := service.CreateSess(req.Username, req.Password)
	if err != nil {
		resp.Errmsg = err.Error()
		resp.ErrNo = iris.StatusInternalServerError
		ctx.JSON(resp.ErrNo, resp)
		return
	}

	ctx.Session().Set(KeySessID, sess.SessID)
	ctx.Session().Set(KeySessUsername, req.Username)

	resp.Success = true
	resp.Sess = *sess
	ctx.JSON(iris.StatusOK, resp)
	return
}

// DeleteSession create session for a user
func DeleteSession(ctx *iris.Context) {
	resp := entity.Resp{}

	sessKey := ctx.Param("key")

	err := service.DeleteSess(sessKey)
	if err != nil {
		resp.Errmsg = err.Error()
		resp.ErrNo = iris.StatusInternalServerError
		ctx.JSON(resp.ErrNo, resp)
		return
	}

	ctx.Session().Delete(sessKey)

	resp.Success = true
	ctx.JSON(iris.StatusOK, resp)
	return
}
