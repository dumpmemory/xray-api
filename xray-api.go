package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"xray-api/action"
	"xray-api/api"
	"xray-api/config"
	"xray-api/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var (
	Config config.CONF
)

func main() {
	currentPath, _ := utils.GetCurrentPath()
	var configPath = flag.String("config", currentPath+"config.yaml", "config路径")
	flag.Parse()
	fmt.Println("Reading config:", *configPath)

	content, _ := ioutil.ReadFile(*configPath)
	yaml.Unmarshal(content, &Config)

	if Config.Key == "" {
		fmt.Println("No access key set. Abort.")
		return
	}
	xray := Config.Xray
	fmt.Println("======= Xray =======")
	fmt.Println("PROTOCOL:", xray.Protocol)
	fmt.Println("PORT:", xray.Port)
	fmt.Println("GRPC:", xray.Grpc)
	// api.WriteConfig(xray, currentPath+"xray-core/config.json")
	action.Xray = xray
	action.Start()

	api.Config = Config
	api.Init()
	if Config.Syncfile != "" {
		data, err := ioutil.ReadFile(Config.Syncfile)
		if err != nil {
			fmt.Println(err)
		} else {
			var users []api.User
			json.Unmarshal(data, &users)
			api.Sync(&users)
		}
	}

	if !Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Any("/ping", action.Ping)
	r.GET("/online", action.IsOnline)
	r.Use(webMiddleware)
	r.POST("/sync", action.Sync)
	r.POST("/addUser", action.AddUser)
	r.POST("/delUser", action.DelUser)
	r.POST("/traffic", action.Traffic)
	r.POST("/restart", action.Restart)
	r.POST("/stat", action.Stat)

	if Config.Listen == "" {
		Config.Listen = ":8080"
	}
	fmt.Println("======= API =======")
	fmt.Println("LISTEN:", Config.Listen)
	fmt.Println("KEY:", Config.Key)
	r.Run(Config.Listen)
}
func webMiddleware(c *gin.Context) {
	if c.Request.Header.Get("key") != Config.Key {
		utils.RespondWithError(500, "Api key Incorrect", c)
		c.Abort()
		return
	}
	c.Next()
}
