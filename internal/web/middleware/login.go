package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		//白名单
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		//登录态校验
		sess := sessions.Default(ctx)
		if sess.Get("userId") == nil {
			//ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.String(http.StatusUnauthorized, "no login")
			println("-----无效或者登录态失败----")
			return
		}
	}

}
