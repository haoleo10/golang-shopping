package api_helper

import (
	"shoppinggolang/utils/pagination"

	"github.com/gin-gonic/gin"
)

var userIdText = "userId"

// 从context获得用户id,把用户id转换成一个parseInt类型
func GetUserId(g *gin.Context) uint {
	return uint(pagination.ParseInt(g.GetString(userIdText), -1))
}
