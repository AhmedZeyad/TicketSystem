package service

import (
	"TicketSystem/engine"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	// PhoneNumber string    `json:"phone_number" db:"phone_number"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	RoleID    int       `json:"rol_id" db:"rol_id"`
}
type RefreshToken struct {
	ID     int    `json:"id" db:"id"`
	UserID int    `json:"userId" db:"userId"`
	Token  string `json:"token" db:"token"`
}

func getUsers() ([]User, error) {
	var users []User
	println("hi form db 124")
	err := engine.DB.Select(&users, "SELECT * from users ")
	if err != nil {

		return nil, err
	}
	return users, nil

}
func (u *User) getUserById() error {

	err := engine.DB.Get(u, "SELECT * from users WHERE id = ?  ", u.ID)
	if err != nil {

		return err
	}
	return nil
}
func (u *User) getUserByEmail() error {

	err := engine.DB.Get(u, "SELECT * from users WHERE email = ?", u.Email)
	if err != nil {

		return err
	}
	return nil
}

func (u *User) InsertUser() error {
	u.CreatedAt = time.Now()
	resoult, err := engine.DB.Exec("INSERT INTO users ( name, email, rol_id, password) VALUES (?,?,?,?)", u.Name, u.Email, u.RoleID, u.Password)
	if err != nil {
		println(err.Error(), time.Now().String())
		return err
	}
	println(u.Password)

	id, err := resoult.LastInsertId()
	if err != nil {
		return err

	}
	u.ID = int(id)
	print(u.ID)
	print(u.CreatedAt.String())
	return nil
}
func (u *User) EditUser() error {
	currentUser := User{ID: u.ID}
	err := currentUser.getUserById()
	if err != nil {
		return errors.New("user not found")

	}
	kk_, err := engine.DB.Exec("UPDATE users SET name = ?, email = ?  WHERE id = ?", u.Name, u.Email, u.ID)
	if err != nil {
		println("erer", kk_, err.Error())
		return errors.New("can't update user")
	}
	*u = currentUser
	return nil
}
func dropUser(id int) error {
	result, err := engine.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // unexpected error checking rows
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil

}

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
	err := context.BindJSON(&loginInfo)

	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{
			"message": "not valid body",
		})
		return
	}

	// get real password
	var user User
	user.Email = loginInfo.Email
	err = user.getUserByEmail()
	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	// comparing passwerd with login pass

	err = comparePassword(user.Password, loginInfo.Password)
	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong phone number and password",
		})
		return
	}
	// generate token
	tokenString, err := engine.GenerateToken(engine.Claims{Subject: user.ID, UserName: loginInfo.Name, RoleID: user.RoleID, Email: user.Email})
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"msg": "failed to gen token"})
		return
	}
	// generate refresh token
	refresToken, err := engine.GenerateRefreshToken(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"msg": "failed in gen process"})
	}
	// add to db
	id, err := isRefreshTokenAvailable(user.ID)
	if err != nil {

		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if id == 0 {

		if err := InsertToken(user.ID, refresToken); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	if id > 0 {
		if err := UpdateRefresToken(tokenString, id); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

	}
	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("token", tokenString, 3600*2, "", "", true, true)
	context.SetCookie("refreshToken", refresToken, 3600*24*30, "", "", true, true)
	context.JSON(http.StatusOK, gin.H{"token": tokenString, "refreshtoken": refresToken})
}

func tokeAdding(id int) error {
	tokenID, err := isRefreshTokenAvailable(id)
	if err != nil {
		return err
	}
	if tokenID == -1 {

		return errors.New("token not found")
	}
	return nil
}
func isRefreshTokenAvailable(userId int) (int, error) {
	var refreshToken RefreshToken
	if err := engine.DB.Get(&refreshToken, "select id from refreshToken where userId=?", userId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, nil
		}
		return -1, err

	}

	return refreshToken.ID, nil
}
func InsertToken(userid int, token string) error {
	resoult, err := engine.DB.Exec(`insert into refreshToken (userId,token) values(?,?)`, userid, token)
	if err != nil {
		return err
	}
	effectedRow, err := resoult.RowsAffected()
	if err != nil {
		return err
	}
	if effectedRow == 0 {
		return errors.New("plz try again")
	}
	return nil
}
func UpdateRefresToken(token string, id int) error {
	_, err := engine.DB.Exec(`Update refreshToken set  token =? where userId=?`, token, id)
	if err != nil {
		return err
	}

	return nil
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
	user := rg.Group("/users")
	user.GET("/", engine.CheckAuth, GetUsers)
	user.GET("/:id", engine.CheckAuth, GetUserById)
	user.POST("/", engine.CheckAuth, AddUser)
	user.PUT("/", engine.CheckAuth, EditUser)
	user.DELETE("/:id", engine.CheckAuth, DeleteUser)
	user.POST("/singup", AddUser)
	user.POST("/Login", Login)
	user.POST("/Logout", engine.CheckAuth, func(context *gin.Context) {
		context.SetCookie("token", "", -1, "", "", false, true)
		context.JSON(http.StatusOK, gin.H{"msg": "logout"})
	})
}
