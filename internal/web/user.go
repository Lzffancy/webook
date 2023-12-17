package web

import (
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
	userIdKey            = "userId"
	bizLogin             = "login"
)

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc, //外界传入
	}
}

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *service.UserService
}

func (h *UserHandler) RegisterRouters(server *gin.Engine) {
	server.POST("/users/signup", h.SignUp)
	server.POST("/users/login", h.Login)
	server.POST("/users/edit", h.Edit)
	server.GET("/users/profile", h.Profile)

	// ug := server.Group("/users")
	// ug.POST("/signup", h.SignUp)
	// ug.POST("/login", h.Login)
	// ug.GET("/edit", h.Edit)
	// ug.POST("/profile", h.Profile)

}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	// ctx.String(http.StatusOK, "hello ,your in signUp")
	//内部结构体
	type SignupReq struct {
		Email           string `json:"email"` //字段标签，这个Email在json中是email
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignupReq
	if err := ctx.Bind(&req); err != nil { //bind自动根据SignupReq结构体控制的格式进行校验和填充
		return
	}

	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "invalid email type or regex error")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "invalid isPassword type or regex error")
		return
	}

	if !isEmail {
		ctx.String(http.StatusOK, "invalid email type")
		return
	}

	if !isPassword {
		ctx.String(http.StatusOK, "invalid isPassword type")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "check Password different")
		return
	}

	err = h.svc.SignUp(ctx, domain.User{
		Password: req.Password,
		Email:    req.Email,
	})

	switch err {
	case nil:
		ctx.String(http.StatusOK, "register ok")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "email used,please change")
	default:
		ctx.String(http.StatusOK, "sys error")
	}

}
func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"` //字段标签，这个Email在json中是email
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "login error 0 ")
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)

	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			MaxAge: 60 * 60,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "login session sys error")
			return
		}
		ctx.String(http.StatusOK, "login ok")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "login InvalidUserOrPassword")
	default:
		ctx.String(http.StatusOK, "login sys error")

	}

}

func (h *UserHandler) Edit(ctx *gin.Context) {
	//用户编辑个人信息
	//每次请求都是覆盖,所以每次参数都必填
	//绑定请求结构体
	type EditReq struct {
		Nickname string `json:"nickname" binding:"required,min=3,max=20"`
		Birthday string `json:"brithday" binding:"required"`
		AboutMe  string `json:"aboutMe" binding:"required,min=3,max=200"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, " Edit param not ok")
		return
	}

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "Birthday is not ok")
		return
	}
	sess := sessions.Default(ctx)
	u := sess.Get("userId")
	if u == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		println("-----无效或者登录态失败byEdit----")
		return
	}
	uid, ok := u.(int64)
	if !ok {
		// 类型转换失败
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uid,
		Nickname: req.Nickname,
		AboutMe:  req.AboutMe,
		Birthday: birthday,
	})
	if err != nil {
		ctx.String(http.StatusOK, "sys error Edit")
		return
	}
	ctx.String(http.StatusOK, "Edit is ok")

}
func (h *UserHandler) Profile(ctx *gin.Context) {

	sess := sessions.Default(ctx)
	uid := sess.Get("userId")
	if uid == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		//一般不会走到这里，中间件校验了登录态
		println("-----无效或者登录态失败byEdit----")
		return
	}
	uid64, ok := uid.(int64)
	if !ok {
		// 类型转换失败
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	println("uin id ", uid64)
	u, err := h.svc.FindById(ctx, uid64)
	if err != nil {
		ctx.String(http.StatusOK, "sys error profile")
		return
	}
	type pUser struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	ctx.JSON(http.StatusOK, pUser{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday.Format(time.DateOnly),
	})
}
