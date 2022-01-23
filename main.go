package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"onlineJudge/common"
	"onlineJudge/judge"
	"onlineJudge/middlewares"
	"onlineJudge/models"
	"onlineJudge/routers"
	"os"
)

const (
	DebugMode = "debug"
	ReleaseMode = "release"
	TestMode = "test"
)

func main() {
	InitConfig()
	common.InitDB()
	models.AutoMigrate()
	common.Redis()
	go judge.Judge()
	gin.SetMode(ReleaseMode)
	r := gin.Default()
	r.Use(middlewares.Cors())
	routers.MountRouter(r)
	port := viper.GetString("server.port")
	if port != "" {
		r.Run(":" + port)
	} else {
		r.Run()
	}
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir+"/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

