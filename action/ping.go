package action

import (
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.Writer.WriteString("pong")
    c.Done()
}
