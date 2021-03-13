package Test

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type ImplDB struct {
	DB *gorm.DB
}

type User struct {
	UserID   uint   `gorm:"primary_key;auto_increment:false"`
	UserName string `gorm:"size:255"`
}

func (i *ImplDB) initDB() {
	var err error
	// USER、PASS、DBNMEなどは各自設定してある値を記述
	DBMS := "mysql"
	USER := "root"
	PASS := "sunica"
	DBNAME := "trpg"

	CONNECT := USER + ":" + PASS + "@" + "/" + DBNAME + "?charset=utf8&parseTime=True&loc=Local"
	i.DB, err = gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic("DBへの接続に失敗しました")
	}
}

// スキーマのマイグレーション
func (i *ImplDB) initMigration() {
	i.DB.AutoMigrate(&User{})
}

func main() {
	i := &ImplDB{}
	println(i)
	// DBに接続
	i.initDB()

	// main関数が終わる際にDBの接続を切る
	defer i.DB.Close()

	// スキーマのマイグレーション
	i.initMigration()

	insertUser := User{}
	user := User{}
	user2 := User{}

	insertUser.UserID = 1
	insertUser.UserName = "hoge"

	// 作成
	// INSERT INTO users(user_id,user_name) VALUES(1,'hoge');
	i.DB.Create(&insertUser)

	// 取得
	// SELECT * FROM users WHERE user_id = 1;
	i.DB.Find(&user, "user_id = ?", 1)

	fmt.Println("取得したuserの値は", user)

	// 更新
	// UPDATE users SET user_name = 'fuga' WHERE user_id = 1 and user_name = 'hoge';
	i.DB.Model(&user).Update("user_name", "fuga")

	fmt.Println("更新後のuser:", user)

	// 削除
	// DELETE FROM users WHERE user_id = 1 and user_name = 'hoge';
	i.DB.Delete(&user)

	if err := i.DB.Find(&user2, "user_id = ?", 1).Error; err != nil {
		// エラーハンドリング
		fmt.Println("存在しませんでした")
	} else {
		fmt.Println("取得したuserの値は", user2)
	}
}
