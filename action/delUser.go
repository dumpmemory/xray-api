package action

import (
	"xray-api/api"

	"github.com/gin-gonic/gin"

	"xray-api/utils"
)

func DelUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		var success = false
		var msg string
		uuid := c.PostForm("uuid")
		if uuid == "" {
			utils.RespondWithError(403, "uuid required", c)
			return
		}
		err := api.RemoveUser(uuid)
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
}
