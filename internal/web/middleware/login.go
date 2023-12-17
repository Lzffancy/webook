package middleware

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {

	gob.Register(time.Now()) //go语言存入reids时候需要对 字符串序列化
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		//白名单
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		//登录态校验
		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			//ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.String(http.StatusUnauthorized, "no login")
			println("-----无效或者登录态失败----")
			return
		}
		now := time.Now()
		const updateTimeKey = "update_time"
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		// session 续期
		if now.Sub(lastUpdateTime) > time.Second*10 {
			fmt.Print("session update now")
		}
		if val == nil || (!ok) || (now.Sub(lastUpdateTime) > time.Second*10) {
			sess.Set(updateTimeKey, now)
			sess.Set("userKey", userId)
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
