# 1 使用Go mod管理模块

初始化

```bash
$ go mod init github.com/tuxnotes/goessential
```

一般是公司域名+项目名称，比如：
```bash
go mod init oceanlearn.teach/ginessential
```
下载gin
```bash
$ go get -u github.com/gin-gonic/gin
```
# 2 实现用户注册

```go
package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		// 获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		// 数据验证
		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return // 如果手机号不符合要求就不用走下面的步骤了，直接返回
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能小于6位"})
			return // 如果密码不符合要求就不用走下面的步骤了，直接返回
		}
		// 如果name没有传，就给一个10位的随机字符串
		if len(name) == 0 {
			name = RandomString(10)
		}
		log.Println(name, telephone, password)
		// 判断手机号是否存在需要查库，所以这里这是暂时使用上面的log打印一下日志

		// 创建用户

		// 返回结果

		ctx.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	// r.Run() // 监听并在 0.0.0.0:8080 上启动服务
	panic(r.Run())
}

func RandomString(n int) string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
```
# 3 连接数据库

这里实例gorm连接数据库。首先打开gorm的官网：https://gorm.io/docs/

安装：

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
```
安装完成后，然后定义我们的model
```
type User struct {

}
```
完成代码如下：
```go
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
```
# 4 项目重构关注分离

前面引入了gorm实现了用户注册API，但是代码都写在了一个文件中，这样下去此文件将会变得越来越大，代码越来越难以维护。
因此这里重构代码，使得项目更具有结构性和维护性。

在重构之前，使用git将代码添加到版本库，并进行一次提交。

## 4.1 重构
首先将将model(即User struct)迁移到专门的包，model包。
在ginessential目录下创建model目录，在model目录下创建user.go文件

## 4.2 将handler迁移到控制器的包

handler即使r.POST参数中的func部分。

```
mkdir controller
cd controller && touch UserController.go
```

# 5 用户登录
先提交一个版本
首先定义路由，然后定义控制器

# 6 使用jwt生成token并认证路由

使用jwt实现用户认证和未知用户登录状态
首先安装jwt包到我们的项目中
```bash
go get -u github.com/dgrijalva/jwt-go
```
在common目录下创建jwt.go文件

token由三部分组成：
- 第一部分：协议头，存储token使用的加密协议
- 第二部分：负载，payload，claim部分存储的信息
- 第三部分：前两个部分加上key后的hash值

```bash
echo $第一部分 | base64 -D
```
查看第一部分的内容
```bash
echo $第二部分 | base64 -D
```
查看第二部分的内容

接下来写认证中间件，创建新的目录middleware,并在目录下创建AuthMiddleware.go文件

创建用户信息的路由

接着用中间件保护用户信息接口

# 7 处理信息返回时的敏感字段及封装统一的返回格式
## 7.1 处理敏感字段
给前端只返回用户名和手机号，其他的都不用返回，所以这里定义一个UserDto结构体
```go
type UserDto struct { // 只返回给前端用户名和手机号，其他都不用返回
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}
```
然后定义一个转换的函数
```go
func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name: user.Name,
		Telephone: user.Telephone,
	}
}
```

然后再controller中将user转换成UserDto
```go
func Info(ctx *gin.Context) {
	// 获取用户信息的时候，用户已经通过了认证，索引从上下文中获取到用户的信息
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
	})
}
```

## 7.2 封装HTTP返回
项目目录下新建response目录

然后修改controller中的注册等模块的代码
# 8 从文件中读取配置

数据库连接信息分散在各个文件，管理起来就很不方便，所以需要进行配置集中化管理。这里在项目中引入config组件：viper

安装
```bash
go get github.com/spf13/viper
```
然后在项目目录下创建一个config目录,使用yaml文件来写我们的配置项application.yml
```yaml
server:
  port: 1016

datasource:
  driveName: mysql
  host: 127.0.0.1
  port: 3306
  database: ginessential
  username: root
  password: root
  charset: utf8
```
然后在main.go中定义一个函数
```go
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
```
main函数中要先运行上面的函数所以在main函数的第一行添加：InitConfig()

然后修改数据初始化函数，即common下的database.go文件中的InitDB()函数
```go
// driverName := "mysql"
	driverName := viper.GetString("datasource.driveName")
	// host := "localhost"
	host := viper.GetString("datasource.host")
	// port := "3306"
	port := viper.GetString("datasource.port")
	// database := "ginessential"
	database := viper.GetString("datasource.database")
	// username := "root"
	username := viper.GetString("datasource.username")
	// password := "root"
	password := viper.GetString("datasource.password")
	// charset := "utf-8"
	charset := viper.GetString("datasource.charset")
```
完成上述数据库的配置后，在main函数中修改监听端口
```go
port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	// 没有配置端口，使用默认的8080端口
	panic(r.Run())
```





