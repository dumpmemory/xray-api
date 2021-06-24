package action

import (
	"strconv"
	"xray-api/api"

	"github.com/gin-gonic/gin"
)

func Traffic(c *gin.Context) {
	reset, _ := strconv.ParseBool(c.DefaultPostForm("reset", "false"))
	res := api.Traffic(reset)
	var success = true
	var msg string
	c.JSON(200, gin.H{
		"success": success,
		"msg":     msg,
		"data":    res,
	})
}
