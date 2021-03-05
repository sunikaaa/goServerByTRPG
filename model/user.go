package model

//User ... this is commonUser
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// CreateUser ... this is creating User
func CreateUser(user *User) {
	db.Create(user)
}

// FindUser ... this is FindUser
func FindUser(u *User) User {
	var user User
	db.Where(u).First(&user)
	return user
}
