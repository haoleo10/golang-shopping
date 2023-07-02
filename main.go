package main

import (
	"fmt"
	"log"
	"net/http"
	"shopping/api"
	"shopping/utils/graceful"

	"time"

	_ "shopping/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 电商项目
// @description 电商项目
// @version 1.0

// @host localhost:8080
// @BasePath /
func main() {
	//首先实例化gin
	r := gin.Default()
	//注册中间件
	registerMiddlewares(r)
	//注册所有的handler
	api.RegisterHandlers(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//server启动
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
	graceful.ShutdownGin(srv, time.Second*5)

}

// 注册中间件
func registerMiddlewares(r *gin.Engine) {
	r.Use(
		gin.LoggerWithFormatter(
			//日志格式
			func(param gin.LogFormatterParams) string {

				return fmt.Sprintf(
					"%s - [%s] \"%s %s %s %d %s %s\"\n",
					param.ClientIP,
					param.TimeStamp.Format(time.RFC3339),
					param.Method,
					param.Path,
					param.Request.Proto,
					param.StatusCode,
					param.Latency,
					param.ErrorMessage,
				)
			}))
	r.Use(gin.Recovery())
}
