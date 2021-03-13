package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

//User ... this is commonUser
type User struct {
	ID       int    `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email" gorm:"type:varchar(100);unique_index" `
}

// CreateUser ... this is creating User
func CreateUser(user *User) {
	fmt.Print(&db)

	db.Create(&user)
}

// FindUser ... this is FindUser
func FindUser(u *User) User {
	var user User
	db.Where(&u).First(&user)
	return user
}

// TestDB ...
func TestDB(DB *gorm.DB) {
	fmt.Println(db)

}

// ShowAll ... this is dbinnnerUsershow
func ShowAll() {
	var allUsers []User
	db.Find(&allUsers)
	fmt.Println(allUsers)
}
