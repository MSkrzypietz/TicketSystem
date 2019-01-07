package utils

import (
	"net/http"
)

type Authenticator interface {
	Authenticate(user, password string) bool
}

type AuthenticatorFunc func(user, password string) bool

func (af AuthenticatorFunc) Authenticate(user, password string) bool {
	return af(user, password)
}

func AuthWrapper(authenticator Authenticator, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUserFromCookie(r)
		if err != nil {
			RemoveCookie(w, "requested-url-while-not-authenticated")

			http.SetCookie(w, &http.Cookie{
				Name:     "requested-url-while-not-authenticated",
				Value:    r.URL.RequestURI(),
				Path:     "/",
				HttpOnly: true,
				MaxAge:   60,
			})

			http.Redirect(w, r, "/signIn", http.StatusFound)
			return
		}

		if authenticator.Authenticate(user.Username, user.Password) {
			handler(w, r)
		} else {
			http.Redirect(w, r, ErrorUnauthorized.ErrorPageURL(), http.StatusFound)
		}
	}
}
