package main

import (
	"github.com/gin-gonic/gin"
	"oceanlearn.teach/ginessential/common"
)

func main() {
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
	panic(r.Run())
}
