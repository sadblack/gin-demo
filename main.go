package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"jkdev.cn/api/common"
	"os"
)

func main() {
	//解析配置
	InitConfig()
	//初始化数据库
	db := common.InitDB()
	//最后关闭数据库
	defer db.Close()
	//使用 gin 的默认配置，
	// gin 的默认配置 = gin 的 初始化配置 + 日志filter + 异常捕获filter
	// 可以点进去看看
	r := gin.Default()
	//配置路由，即 url -> controller 的路径
	r = CollectRoute(r)
	//获取配置的 端口
	port := viper.GetString("server.port")
	if port != "" {
		//如果指定了端口，就在指定的端口上运行
		panic(r.Run(":" + port))
	}
	//默认端口 8080
	panic(r.Run()) // listen and serve on 0.0.0.0:8080
}

func InitConfig() {
	workDir, _ := os.Getwd()
	// viper 用来 解析配置文件，可以解析 json、yml、ini 等格式
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	//解析，解析完成后，默认会保存在 viper 里面，相当于内部有一个 hashmap
	err := viper.ReadInConfig()
	if err != nil {
		panic("")
	}
}
