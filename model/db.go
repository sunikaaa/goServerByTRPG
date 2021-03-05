package model

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

//InitDB ... this is startDB
func InitDB() *gorm.DB {
	var err error
	db, err := gorm.Open("mysql", "root:sunica@/trpg")
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	return db
}
