package user

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	//32、64位
	cookieStoreAuthKey = securecookie.GenerateRandomKey(64)
	//16位、24、32位
	cookieStoreAEncryptKey = securecookie.GenerateRandomKey(32)
)

var sessionStore *sessions.CookieStore

func init() {
	sessionStore = sessions.NewCookieStore(
		[]byte(cookieStoreAuthKey),
		[]byte(cookieStoreAEncryptKey),
	)
	sessionStore.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   60 * 15,
	}
}
