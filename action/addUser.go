package action

import (
	api "xray-api/api"
	"xray-api/utils"

	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context) {
	var success = false
	var msg string
	uuid := c.PostForm("uuid")
	email := c.PostForm("email")
	if uuid == "" {
		utils.RespondWithError(403, "uuid required", c)
		return
	}
	if email == "" {
		utils.RespondWithError(403, "email required", c)
		return
	}
	err := api.AddUser(uuid, email)
	if err == nil {
		success = true
	} else {
		msg = err.Error()
	}
	c.JSON(200, gin.H{
		"success": success,
		"msg":     msg,
	})
}
