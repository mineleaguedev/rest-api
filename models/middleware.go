package models

import (
	"net/http"
	"time"
)

type JWTMiddleware struct {
	// Realm name to display to the user
	Realm string

	// Access token secret key
	AccessTokenKey []byte

	// Refresh token secret key
	RefreshTokenKey []byte

	// Duration that an access token is valid
	AccessTokenTime time.Duration

	// Duration that a refresh token is valid
	RefreshTokenTime time.Duration

	// Set the identity key
	IdentityKey string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header
	TokenHeadName string

	// Optionally return the token as a cookie
	SendCookie bool

	// Duration that a cookie is valid
	CookieMaxAge time.Duration

	// Allow insecure cookies for development over http
	SecureCookie bool

	// Allow cookies to be accessed client side for development
	CookieHTTPOnly bool

	// Allow cookie domain change for development
	CookieDomain string

	// Allow cookie name change for development
	CookieName string

	// CookieSameSite allow use http.SameSite cookie param
	CookieSameSite http.SameSite
}
