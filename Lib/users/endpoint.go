package users

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(context *gin.Context) {

	users, err := getUsers()
	if err != nil {
		println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return

	}

	context.JSON(http.StatusOK, gin.H{
		"status": "Sucsess",
		"users":  users,
	})

}
func GetUserById(context *gin.Context) {
	struid := context.Param("id")
	uid, err := strconv.Atoi(struid)
	if err != nil {
		println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "not valid id",
		})
		return
	}
	user := User{ID: uid}
	err = user.getUserById()
	if err != nil {
		println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "user not found",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status": "Sucsess",
		"user":   user,
	})

}

func AddUser(context *gin.Context) {
	println("add user")
	var u User
	err := context.BindJSON(&u)

	if err != nil {
		println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "not valid user body",
		})
		return
	}
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		println("dASad", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "can't hash password",
		})
		return
	}
	u.Password = string(hashedPassword)

	err = u.InsertUser()
	if err != nil {
		println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  err.Error(),
			"message": "can't add user",
		})
		return

	}
	context.JSON(http.StatusOK, gin.H{
		"status": "Sucsess",
		"user":   u,
	})
}
func EditUser(context *gin.Context) {
	var newUserInfo User
	err := context.BindJSON(&newUserInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "not valid user body",
		})
		return
	}
	err = newUserInfo.EditUser()
	if err != nil {
		println(newUserInfo.ID)
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "can't update user",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status": "Sucsess",
		"user":   newUserInfo,
	})

}
func DeleteUser(context *gin.Context) {
	println("hi form drop")
	strid := context.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "not valid id",
		})
		return
	}
	err = dropUser(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"status": "Sucsess",
		})
	}
}


func Login(context *gin.Context) {
	// logi password
	var loginInfo User	
err:=context.BindJSON(&loginInfo)

if err!=nil{

	context.JSON(http.StatusBadRequest, gin.H{
		"message": "not valid body",
	})
	return
}

// get real password
var user User
user.PhoneNumber=loginInfo.PhoneNumber
err=user.getUserByPhonNumber()
if err!=nil{

	context.JSON(http.StatusBadRequest, gin.H{
		"message": "user not found",
	})
	return
}
	// comparing passwerd with login pass

err=comparePassword(user.Password,loginInfo.Password)
if err!=nil{

	context.JSON(http.StatusBadRequest, gin.H{
		"message": "wrong phone number and password",
	})
	return
}
		// generate token
		token:=jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":user.ID,
			"role":user.RoleID,
			"email":user.Email,
			"phone_number":user.PhoneNumber,
			"name":user.Name,
			"exp":time.Now().Add(time.Hour * 24*30).Unix(),
		})
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "failed to gen token"})
		return
	}
			// set toke in cookie

	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("token", tokenString, 3600*24*30, "", "", false, true)
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func hashPassword(password string) (string, error) {
	strCost := os.Getenv("COST")
	cost, err := strconv.Atoi(strCost)
	if err != nil {
		return "", err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}
func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func UserRoutes(rg *gin.RouterGroup) {
	user:=rg.Group("/users")
	user.GET("/",GetUsers)
	user.GET("/:id",GetUserById)
	user.POST("/",AddUser)
	user.PUT("/",EditUser)
	user.DELETE("/:id",DeleteUser)
	user.POST("/singup",AddUser)
	user.POST("/Login",Login)
}