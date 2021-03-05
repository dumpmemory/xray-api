package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"xray-api/action"
	"xray-api/api"
	"xray-api/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var (
	accessKey string
)

func main() {
	currentPath, _ := utils.GetCurrentPath()
	var configPath = flag.String("config", currentPath+"config.yaml", "config路径")
	flag.Parse()
	fmt.Println("Reading config:", *configPath)

	content, _ := ioutil.ReadFile(*configPath)
	config := make(map[string]interface{})
	yaml.Unmarshal(content, &config)
	accessKeyI := config["key"]
	if accessKeyI == nil {
		fmt.Println("No access key set. Abort.")
		return
	}
	accessKey = accessKeyI.(string)

	xray := config["xray"].(map[interface{}]interface{})

	fmt.Println("=======xray=======")
	fmt.Println("PROTOCOL:", xray["protocol"])
	fmt.Println("PORT:", xray["port"])
	fmt.Println("GRPC:", xray["grpcPort"])
	fmt.Println("====================")

	api.WriteConfig(xray, currentPath+"xray-core/config.json")

	action.Start()

	api.Init("127.0.0.1", uint16(xray["grpcPort"].(int)), xray["protocol"].(string))
	if xray["protocol"].(string) == "shadowsocks" {
		api.SetSSmethod(xray["method"].(string))
	}

	if config["debug"].(bool) != true {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(webMiddleware)
	r.POST("/sync", action.Sync())
	r.POST("/addUser", action.AddUser())
	r.POST("/delUser", action.DelUser())
	r.POST("/traffic", action.Traffic())
	r.POST("/restart", action.Restart())
	r.POST("/status", action.Status())

	var address string

	addressI := config["listen"]
	if addressI == nil {
		address = "127.0.0.1:8080"
	} else {
		address = addressI.(string)
	}
	fmt.Println("LISTEN:", address)
	fmt.Println("KEY:", accessKey)
	_ = r.Run(address)
}

func webMiddleware(c *gin.Context) {
	if c.PostForm("key") == "" {
		utils.RespondWithError(401, "API token required", c)
		return
	}
	if c.PostForm("key") != accessKey {
		utils.RespondWithError(403, "API token incorrect", c)
		return
	}
	c.Next()
}
