package auth

import (
	"net/http"
	"os"
	"strings"
)

func CookieName() string {
	name := os.Getenv("JWT_COOKIE_NAME")
	if name == "" {
		return "session"
	}
	return name
}

func ExtractCookieDomain(host string) string {
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}
	return host
}

func SetSessionCookie(w http.ResponseWriter, token string, host string, secure bool) {
	domain := ExtractCookieDomain(host)
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName(),
		Value:    token,
		Path:     "/",
		MaxAge:   86400, // 24h
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
	})
}

func ClearSessionCookie(w http.ResponseWriter, host string, secure bool) {
	domain := ExtractCookieDomain(host)
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName(),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Domain:   domain,
	})
}
