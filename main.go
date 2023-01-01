package main

//	@title			Go-NOTE
//	@version		1.0
//	@description	Golang work.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	zhangxin
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		127.0.0.1:8880
//	@BasePath	/api/v1

import (
	Model "NoteGin/Models"

	router "NoteGin/Router"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	dataBase string = "mysql"
	userName string = "root"
	passWord string = "111111"
	IP       string = "127.0.0.1"
	Port     string = "3306"
	dbName   string = "h_list"
)

func initMysql() (err error) {
	Model.Db1, err = gorm.Open(dataBase, userName+":"+passWord+"@("+IP+":"+Port+")/"+dbName+"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return
}

func main() {
	initMysql()
	defer Model.Db1.Close()
	Model.Db1.AutoMigrate(&Model.User{}, &Model.Note{})
	engine := router.Router()
	engine.Run(":8080")
}
