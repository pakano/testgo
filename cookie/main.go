package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

var hashKey = securecookie.GenerateRandomKey(64)
var blocKey = securecookie.GenerateRandomKey(32)

//当blocKey为空hashKey不为空，仅为签名;
//如果blocKey不为空，则为加密。
var s = securecookie.New(hashKey, blocKey)

//对 cookie进行签名
func SetCookieHandler(w http.ResponseWriter, r *http.Request) {
	encoded, err := s.Encode("cookie-name", "cookie-value")
	if err == nil {
		cookie := http.Cookie{
			Name:     "cookie-name",
			Value:    encoded,
			Domain:   "www.baidu.com", //如果访问的是www.example.com 域名，想让其子域名访问，可设置为example.com
			Path:     "/a",
			HttpOnly: true,
			Expires:  time.Now().Add(time.Second * 5),
			//	Secure: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
		fmt.Fprintln(w, encoded)
	}
}

func GetCookieHandler(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("cookie-name"); err == nil {
		var value string
		if err = s.Decode("cookie-name", cookie.Value, &value); err == nil {
			fmt.Fprintln(w, value)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Fprintf(w, "no cookie")
	}
}

func main() {
	http.HandleFunc("/set", SetCookieHandler)
	http.HandleFunc("/b/b", GetCookieHandler)

	http.ListenAndServe(":9898", nil)
}

// curl http://:::9898/set
// curl --cookie "cookie-name=MTY2MTAwNTIxNXxEd3dBREdOdmIydHBaUzEyWVd4MVpRPT18Juh-W163OW_-CHZTd8Jq1Z7j0dHsIyvylBzk9TKpYVg=" http://:::9898/get
