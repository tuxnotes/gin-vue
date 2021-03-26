package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	// 接下来采用查库的方式
	if isTelephoneExist(DB, telephone) { // 如果用户存在就不允许注册
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已存在"})
		return
	}
	log.Println(name, telephone, password)
	// 创建用户， 如果用户不存在，就新建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { // 如果有异常，则加密错误，这是一个系统级错误
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "加密错误"})
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword), // 保存加密后的密码
	}
	DB.Create(&newUser)
	ctx.JSON(200, gin.H{"code": 200, "msg": "注册成功"})
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	// 获取参数，登录只需要手机号和密码，不需要name
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	// 数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}

	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码长度不足6位"})
		return
	}
	// 判断手机是否存在，手机存在才能登录
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 { // 用户不存在
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	// 判断密码是否正确，用户的密码不能明文保存，用户注册的时候需要将密码加密。使用bcrypt来判断
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}
	// 密码正确发放token给前端
	// token := "11" // 先用简单的
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error: %v\n", err)
		return
	}
	// 返回结果
	ctx.JSON(200, gin.H{"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功"})
}

func Info(ctx *gin.Context) {
	// 获取用户信息的时候，用户已经通过了认证，索引从上下文中获取到用户的信息
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": user},
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 { // 用户ID不为0，说明用户存在，所以返回true
		return true
	}
	return false
}
