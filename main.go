package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not nul;unique"`
	Password  string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()
	// defer db.Close()
	r := gin.Default()
	/*
		用户注册：
		1. 获取表单数据
		2. 验证数据的有效性
		3. 注册成功
	*/
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		// 参数获取
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")
		// 检查用户名,如果用户名为空，则返回一个10个字母的字符串
		if len(name) == 0 {
			name = RandomString(10)
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
		if isTelephoneExist(db, telephone) { // 如果用户存在就不允许注册
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已存在"})
			return
		}
		log.Println(name, telephone, password)
		// 创建用户， 如果用户不存在，就新建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)
		ctx.JSON(200, gin.H{"msg": "注册成功"})
	})
	panic(r.Run())
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 { // 用户ID不为0，说明用户存在，所以返回true
		return true
	}
	return false
}

func RandomString(n int) string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, 10)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func InitDB() *gorm.DB {
	// driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "root"
	charset := "utf-8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database, err" + err.Error())
	}
	db.AutoMigrate(&User{}) // 自动创建数据表
	return db
}
