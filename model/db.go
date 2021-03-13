package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Db ... this is Db
var db *gorm.DB

//InitDB ... this is startDB
func InitDB() *gorm.DB {
	var err error
	db, err = gorm.Open("mysql", "root:sunica@/trpg")
	if err != nil {
		panic(err)
	}
	fmt.Println("db opened")

	// defer db.Close()
	return db
}
