package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"oceanlearn.teach/ginessential/common"
)

func main() {
	InitConfig() // 项目启动开始就需要读取配置
	db := common.InitDB()
	defer db.Close()
	// defer db.Close()
	r := gin.Default()
	/*
		用户注册：
		1. 获取表单数据
		2. 验证数据的有效性
		3. 注册成功
	*/
	r = CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	// 没有配置端口，使用默认的8080端口
	panic(r.Run())
}

func InitConfig() {
	workDir, _ := os.Getwd()                 // 获取当前工作目录
	viper.SetConfigName("application")       // 要读取的配置文件名称
	viper.SetConfigType("yml")               // 读取的文件的类型
	viper.AddConfigPath(workDir + "/config") // 设置配置文件的路径
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}
