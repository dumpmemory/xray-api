package action

import (
	"os/exec"
	"xray-api/api"
	"xray-api/config"
	"xray-api/utils"

	"github.com/gin-gonic/gin"
)

var (
	Cmd    exec.Cmd
	Xray   config.XRAY
	Online bool
)

func Start() error {
	currentPath, _ := utils.GetCurrentPath()
	xrayBin := currentPath + "xray-core/xray"
	// xrayBin := "/usr/bin/xray"
	Cmd = exec.Cmd{
		Path: xrayBin,
		Args: []string{
			"run",
			// "-config", currentPath + "xray-core/config.json"
		},
		Dir: currentPath + "xray-core",
	}
	stdin, err := Cmd.StdinPipe()
	if err != nil {
		return err
	}
	// Cmd.Stdout = os.Stdout
	// Cmd.Stderr = os.Stderr
	err = Cmd.Start()
	if err != nil {
		return err
	}
	err = api.WriteConfig(Xray, stdin)
	if err != nil {
		return err
	}
	stdin.Close()
	go func() {
		Online = true
		Cmd.Wait()
		Online = false
		if Xray.AutoRestart {
			Start()
			api.ReAddUsers()
		}
	}()
	return nil
}
func Stop() (err error) {
	err = Cmd.Process.Kill()
	if err != nil {
		return
	}
	Cmd.Wait()
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
