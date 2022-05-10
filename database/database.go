package database

import (
	"fmt"
	"kilimanjaro-api/database/models"
	"kilimanjaro-api/database/orm"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"kilimanjaro-api/config"
)

func InitDatabase() {
	cfg := config.GetConfig()

	username := cfg.DBUser
	password := cfg.DBPass
	dbName := cfg.DBName
	dbHost := cfg.DBHost
	dbPort := cfg.DBPort
	dbType := cfg.DBType

	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbHost, dbPort, dbName)
	fmt.Println(dbUri)

	conn, err := gorm.Open(dbType, dbUri)
	if err != nil {
		fmt.Print(err)
	}

	orm.DBCon = conn

	orm.DBCon.Set("database:table_options", "ENGINE=InnoDB")
	orm.DBCon.Set("database:table_options", "collation_connection=utf8_general_ci")

	orm.DBCon.Debug().AutoMigrate(&models.User{}, &models.Product{}, &models.Vendor{})
	orm.DBCon.LogMode(false)

}
