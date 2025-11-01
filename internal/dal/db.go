package dal

import (
	"fmt"
	"go-native-webserver/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB // sql.DB, gorm.DB, depends on ORM choice...

func InitDB() error {
	c := config.GetServerConfig() // load config if needed
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	return err
}

func GetDB() *gorm.DB {
	return DB
}
