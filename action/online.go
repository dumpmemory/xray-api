package action

import (
	"github.com/gin-gonic/gin"
)

func IsOnline(c *gin.Context) {
	c.JSON(200, Online)
}
