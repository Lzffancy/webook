package main

import (
	"net/http"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"

	"webook/pkg/ginx/middleware/ratelimit"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// db := initDB()
	// server := initWebServer()
	// initUserHdl(db, server)
	// server.Run(":8082")

	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello ,docker is ok!")
	})
	server.Run(":8082")

}

func initDB() *gorm.DB {
	// db rest
	db, err := gorm.Open(mysql.Open("root:123@root@tcp(192.168.0.110:13316)/webook"))
	if err != nil {
		panic(err)
	}
	dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	// hdl := web.NewUserHandler()
	server := gin.Default()
	server.Use( //注册middleware，接受*Context作为参数即是HandlerFunc (中间件)
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowOrigins:     []string{"http://loaclhots:3000"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowOriginFunc:  func(origin string) bool { return true },
			MaxAge:           12 * time.Hour, //预检有效期，此期间不需要再Option请求
		}),

		func(ctx *gin.Context) {
			println("------my middleware!----")
		},
	)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "192.168.0.110:6379",
	})

	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 1).Build())
	useJWT(server)
	return server

}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us)
	hdl.RegisterRouters(server)
}

func useJWT(server *gin.Engine) {
	login := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}
	//初始化session id 使用cookies存放
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())

	//自带的redis sesion使用了aes加密
	// store, err := redis.NewStore(16, "tcp", "192.168.0.110:6379", "", []byte("REHGZQ0CA8Z528LF7ULLOX9GJ9U6XA7Y"), []byte("REHGZQ0CA8Z528LF7ULLOX9GJ9U6XA71"))
	// if err != nil {
	// 	print(err)
	// 	panic("redis seesion error")
	// }

	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
