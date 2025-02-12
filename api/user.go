package api

import (
	"chat/service"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	// 反序列化
	var UserRegisterService service.UserRegisterService
	if err := c.ShouldBind(&UserRegisterService); err == nil {
		res := UserRegisterService.Register()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}
