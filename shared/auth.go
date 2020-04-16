package shared

import (
	"net/http"
	"os"
	"time"
	"turkdev/data"

	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
)

const (
	authCookieKey = "token"
)

/*SignedInUserClaims represents the data to generate jwt token.*/
type SignedInUserClaims struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customerid"`
	UserName   string `json:"username"`
	Email      string `json:"email"`
	Karma      int    `json:"karma"`
	jwt.StandardClaims
}

/*CustomRequest represents the http request included defaul golang http.Request*/
type CustomRequest struct {
	*http.Request
}

/*GetUser gets authenticated user.*/
func GetUser(r *http.Request) *SignedInUserClaims {
	_, user, err := IsAuthenticated(r)
	if err != nil {
		sentry.CaptureException(err)
		panic(err)
	}
	return user
}

/*IsAuthenticated checks if user is signed in or not.*/
func IsAuthenticated(r *http.Request) (bool, *SignedInUserClaims, error) {
	tokenCookie, err := r.Cookie(authCookieKey)
	if err != nil {
		if err == http.ErrNoCookie {
			return false, nil, nil
		}
		return false, nil, err
	}
	if tokenCookie.Value == "" {
		return false, nil, nil
	}
	token, err := jwt.ParseWithClaims(tokenCookie.Value, &SignedInUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err == nil {
		claims := token.Claims.(*SignedInUserClaims)
		return true, claims, nil
	}

	sentry.CaptureException(err)

	switch err.(type) {
	case nil:
		if !token.Valid {
			return false, nil, err
		}
	case *jwt.ValidationError:
		vErr := err.(*jwt.ValidationError)
		switch vErr.Errors {
		case jwt.ValidationErrorUnverifiable:
			return false, nil, nil
		case jwt.ValidationErrorSignatureInvalid:
			return false, nil, nil
		case jwt.ValidationErrorExpired:
			return false, nil, nil
		default:
			return false, nil, vErr
		}
	default:
		return false, nil, err
	}
	return false, nil, nil
}

/*SetAuthCookie sets the authentication cookie.*/
func SetAuthCookie(w http.ResponseWriter, token string, expirationTime time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:    authCookieKey,
		Value:   token,
		Expires: expirationTime,
	})
}

/*GenerateAuthToken generate jwt token for user login. */
func GenerateAuthToken(user data.User, expirationTime time.Time) (string, error) {
	claims := &SignedInUserClaims{
		ID:         user.ID,
		CustomerID: user.CustomerID,
		UserName:   user.UserName,
		Email:      user.Email,
		Karma:      user.Karma,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
