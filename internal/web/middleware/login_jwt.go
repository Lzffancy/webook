package middleware

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddlewareBuilder struct{}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now()) //go语言存入reids时候需要对 字符串序列化
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		//白名单
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		//登录态校验
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			fmt.Print("---Authorization header 校验失败\n")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			fmt.Print("---Authorization header2 校验失败\n")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			fmt.Print("---jwt tokne 校验失败\n")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			fmt.Print("---jwt tokne 失效或为空\n")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			fmt.Print("---user agent 被修改--\n")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString(web.JWTKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				log.Print(err)
			}
		}
		ctx.Set("user", uc)

	}

}

type UserClaims struct {
	Id        int64
	UserAgent string
	Ssid      string
	jwt.RegisteredClaims
}
