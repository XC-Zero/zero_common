package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.tessan.com/data-center/tessan-erp-common/component"
	"log"
)

const (
	OK            = 200
	BadRequest    = 400
	NotAuthorized = 401
	Forbidden     = 403
	PageNotFound  = 404
	InternalError = 500
	MessageLabel  = "msg"
)

func Error(ctx *gin.Context, msg any, err ...error) {
	log.Printf("[ERROR]%+v", err)
	ctx.JSON(InternalError, map[string]any{
		MessageLabel: msg,
	})
}

func Waring(ctx *gin.Context, msg any, err ...error) {
	log.Printf(`[WARN]%+v`, err)
	ctx.JSON(BadRequest, map[string]any{
		MessageLabel: msg,
	})
}

func Success(ctx *gin.Context, msg any, data any, count int64, hooks ...component.TableHookFunc) {
	if hooks != nil {
		for i := 0; i < len(hooks); i++ {
			data, _ = hooks[i](data, component.Table{})
		}
	}
	ctx.JSON(OK, map[string]any{
		MessageLabel: msg,
		"data":       data,
		"count":      count,
	})
}

func SuccessWithTable(ctx *gin.Context, msg any, data any, table component.Table, count int64, hooks ...component.TableHookFunc) {
	if hooks != nil {
		for i := 0; i < len(hooks); i++ {
			data, table = hooks[i](data, table)
		}
	}
	ctx.JSON(OK, map[string]any{
		MessageLabel: msg,
		"table":      table,
		"data":       data,
		"count":      count,
	})
}

func SuccessWithTableAndDefaultCol(ctx *gin.Context, msg, data, defaultCol any, table component.Table, count int64, hooks ...component.TableHookFunc) {
	if hooks != nil {
		for i := 0; i < len(hooks); i++ {
			data, table = hooks[i](data, table)
		}
	}

	ctx.JSON(OK, map[string]any{
		MessageLabel:  msg,
		"table":       table,
		"default_col": defaultCol,
		"data":        data,
		"count":       count,
	})
}

var AccessDeny = map[string]any{
	MessageLabel: "没有权限哦!",
	"error":      "Access denied!",
}

var AccessTimeout = map[string]any{
	MessageLabel: "登录超时!",
	"error":      "Token timeout!",
}
