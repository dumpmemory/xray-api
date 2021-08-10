package action

import (
	"encoding/json"
	"xray-api/api"
	"xray-api/utils"

	"github.com/gin-gonic/gin"
)

func Sync(c *gin.Context) {
	var newUsers []api.User
	err2 := json.Unmarshal([]byte(c.PostForm("data")), &newUsers)
	if err2 != nil {
		utils.RespondWithError(500, "Couldn't parse JSON data", c)
		return
	}
	res := api.Sync(newUsers)
	c.JSON(200, gin.H{
		"success": true,
		"msg":     res,
	})
}
