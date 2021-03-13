package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"trpg.com/key"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	//ClientIDやClientSecretはGoogle API Console からコピーしてきます。
	//Endpointはgoogle、github、facebookなどがoauth2配下のパッケージに用意されています。
	config := oauth2.Config{
		ClientID:     key.GoogleID,
		ClientSecret: key.GoogleSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "https://stroy-app-by-next-js.vercel.app",      //今回はリダイレクトしないためこれ
		Scopes:       []string{"https://picasaweb.google.com/data/"}, //必要なスコープを追加
	}

	//認証のURLを取得。AuthCodeURLには文字列を渡す。CSRF攻撃の回避のため、本当はリダイレクトのコールバックで検証する。
	url := config.AuthCodeURL("test")
	fmt.Println(url) //認証のURLを表示。ブラウザにコピペする

	//リダイレクト先がないため、ブラウザで認証後に表示されるコードを入力
	var s string
	var sc = bufio.NewScanner(os.Stdin)
	if sc.Scan() {
		s = sc.Text()
	}

	//アクセストークンを取得
	token, err := config.Exchange(oauth2.NoContext, s)
	if err != nil {
		log.Fatalf("exchange error")
	}

	client := config.Client(oauth2.NoContext, token) //httpクライアントを取得

	//取得したclientでGetする
	resp, err := client.GET("https://oauth2.googleapis.com/tokeninfo")
	if err != nil {
		log.Fatalf("client get error")
	}

	//レスポンスを表示
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))

}
