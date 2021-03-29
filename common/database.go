package common

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"oceanlearn.teach/ginessential/model"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
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
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("Failed to connect database, err" + err.Error())
	}
	db.AutoMigrate(&model.User{}) // 自动创建数据表
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
