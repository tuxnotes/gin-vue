package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"oceanlearn.teach/ginessential/common"
	"oceanlearn.teach/ginessential/model"
	"oceanlearn.teach/ginessential/util"
)

func Register(ctx *gin.Context) {
	// db := common.GetDB()
	// DB := common.InitDB()
	DB := common.GetDB()
	// 参数获取
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	// 检查用户名,如果用户名为空，则返回一个10个字母的字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	// 检查密码
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码长度不足6位"})
		return
	}

	// 检查手机号
	// if len(telephone) != 11 {
	// ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
	// return
	// }
	// 接下来采用查库的方式
	if isTelephoneExist(DB, telephone) { // 如果用户存在就不允许注册
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已存在"})
		return
	}
	log.Println(name, telephone, password)
	// 创建用户， 如果用户不存在，就新建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)
	ctx.JSON(200, gin.H{"msg": "注册成功"})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 { // 用户ID不为0，说明用户存在，所以返回true
		return true
	}
	return false
}
