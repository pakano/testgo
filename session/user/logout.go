package user

import "net/http"

func LogOut(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieName)

	session.Values["authenticated"] = false
	session.Save(r, w)
}
