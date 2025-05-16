package engine

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserName  string ` json:"username" db:"username"`
	Email     string `json:"email" db:"email"`
	Subject   int    `json:"sub" db:"id"`
	Audience  string `json:"aud"`
	RoleID    int    `json:"role" db:"role_id"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": claims.UserName,
		"email":    claims.Email,
		"sub":      claims.Subject,
		"aud":      "Ahmed.iq",
		"role":     claims.RoleID,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat":      time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	return tokenString, err
}

func CheckAuth(context *gin.Context) {
	// git token from cookie
	tokenString, err := context.Cookie("token")
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Unautherized user"})
		context.Abort()
		return
	}

	// decode the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRETKEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil && token == nil {
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Unautherized user "})
		context.Abort()
		return

	}
	info := token.Claims.(jwt.MapClaims)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		if info["aud"] != "Ahmed.iq" {
			context.JSON(http.StatusUnauthorized, gin.H{"msg": "not viled token"})
		}

		// chek if token is valid
		if float64(time.Now().Unix()) < info["iat"].(float64) {
			context.JSON(http.StatusUnauthorized, gin.H{"msg": "not viled token"})
		}
		// check token exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			context.JSON(http.StatusUnauthorized, gin.H{"msg": "expired token"})
			context.Abort()
			return

		}
		// chek if user exist
		// Todo check if it work

		_, err := DB.Exec("SELECT * FROM users WHERE ID=?", claims["subject"])
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"msg": "Unautherized  user not found,", "error": err.Error()})

			context.Abort()
			return

		}

		// store user id in context
		context.Set("userId", claims["sub"])
		// continue

		context.Next()

	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"msg": "Unautherized user5"})
		context.Abort()
		return
	}

}

type RefreshTokenClims struct {
	ExpireAt int64 `json:"exp"`
	IssuedAt int64 `json:"iat"`
	ID       int   `json:"id"`
	Subject  int   `json:"sub"`
}

func GenerateRefreshToken(id int) (string, error) {
	refreshToken := RefreshTokenClims{
		ExpireAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		IssuedAt: time.Now().Unix(),
		Subject:  id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": refreshToken.Subject,
		"exp": refreshToken.ExpireAt,
		"iat": refreshToken.IssuedAt,
	})
	stringToken, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		return "", err

	}
	return stringToken, nil
}
func CheckRefreshToken(refreshToken RefreshTokenClims) error {

	if time.Now().Unix() > refreshToken.ExpireAt && time.Now().Unix() <= refreshToken.IssuedAt {
		return nil
	}

	return errors.New("invalid refresh token")
}
