package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
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
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]

		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil

		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
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
