package action

import (
	"os/exec"
	"xray-api/utils"

	"github.com/gin-gonic/gin"
)

var cmd exec.Cmd

func Start() (err error) {
	currentPath, _ := utils.GetCurrentPath()
	xrayBin := currentPath + "xray-core/xray"
	// xrayBin := "/usr/bin/xray"
	cmd := exec.Cmd{
		Path: xrayBin,
		Args: []string{"run", "-config", currentPath + "xray-core/config.json"},
		Dir:  currentPath + "xray-core",
	}
	err = cmd.Start()
	return
}
func Stop() (err error) {
	err = cmd.Process.Kill()
	if err != nil {
		return
	}
	cmd.Wait()
	return
}
func restart() (err error) {
	err = Stop()
	if err != nil {
		return
	}
	err = Start()
	return
}

func Restart(c *gin.Context) {
	err := restart()
	var success = false
	var msg string
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
