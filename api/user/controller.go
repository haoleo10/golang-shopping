package user
//gin的内容
import (
	"net/http"
	"os"
	"shoppinggolang/config"
	"shoppinggolang/domain/user"
	"shoppinggolang/utils/api_helper"
	jwtHelper "shoppinggolang/utils/jwt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	userService *user.Service
	appConfig   *config.Configuration
}

// 实例化，service和配置文件
func NewUserController(service *user.Service, appConfig *config.Configuration) *Controller {
	return &Controller{
		userService: service,
		appConfig:   appConfig,
	}
}

// CreateUser godoc
// @Summary 根据给定的用户名和密码创建用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param CreateUserRequest body CreateUserRequest true "user information"
// @Success 201 {object} CreateUserResponse
// @Failure 400  {object} api_helper.ErrorResponse
// @Router /user [post]
func (c *Controller) CreateUser(g *gin.Context) {
	//控制器里的请求体
	var req CreateUserRequest
	//json数据绑定到结构体
	if err := g.ShouldBind(&req); err != nil {
		//如果有错误返回错误信息
		api_helper.HandleError(g, api_helper.ErrInvalidBody)
		return
	}
	//创建用户
	newUser := user.NewUser(req.Username, req.Password, req.Password2)
	err := c.userService.Create(newUser)
	if err != nil {
		api_helper.HandleError(g, err)
		return
	}

	g.JSON(
		http.StatusCreated, CreateUserResponse{
			Username: req.Username,
		})
}

// Login godoc
// @Summary 根据用户名和密码登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param LoginRequest body LoginRequest true "user information"
// @Success 200 {object} LoginResponse
// @Failure 400  {object} api_helper.ErrorResponse
// @Router /user/login [post]
func (c *Controller) Login(g *gin.Context) {
	var req LoginRequest
	//绑定数据
	if err := g.ShouldBind(&req); err != nil {
		api_helper.HandleError(g, api_helper.ErrInvalidBody)

	}
	//返回当前用户
	currentUser, err := c.userService.GetUser(req.Username, req.Password)
	if err != nil {
		api_helper.HandleError(g, err)
		return
	}
	//验证token
	decodedClaims := jwtHelper.VerifyToken(currentUser.Token, c.appConfig.SecretKey)
	if decodedClaims == nil {
		jwtClaims := jwt.NewWithClaims(
			jwt.SigningMethodHS256, jwt.MapClaims{
				"userId":   strconv.FormatInt(int64(currentUser.ID), 10),
				"username": currentUser.Username,
				"iat":      time.Now().Unix(),
				"iss":      os.Getenv("ENV"),
				"exp": time.Now().Add(
					24 *
						time.Hour).Unix(),
				"isAdmin": currentUser.IsAdmin,
			})
			//生成token
		token := jwtHelper.GenerateToken(jwtClaims, c.appConfig.SecretKey)
		//返回给用户token
		currentUser.Token = token
		//有了token更新一下用户信息
		err = c.userService.UpdateUser(&currentUser)
		if err != nil {
			api_helper.HandleError(g, err)
			return
		}
	}
//返回json格式
	g.JSON(
		http.StatusOK, LoginResponse{Username: currentUser.Username, UserId: currentUser.ID, Token: currentUser.Token})
}

// 验证token
func (c *Controller) VerifyToken(g *gin.Context) {
	token := g.GetHeader("Authorization")
	decodedClaims := jwtHelper.VerifyToken(token, c.appConfig.SecretKey)

	g.JSON(http.StatusOK, decodedClaims)

}
