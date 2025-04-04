package users;

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(context *gin.Context) {

	users, err := getUsers()
	if err != nil {
		println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

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
user:=User{ID: uid}
err=user.getUserById()
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
var u User
err:=context.BindJSON(&u)

if err!=nil{
	println(err.Error())
	context.JSON(http.StatusBadRequest, gin.H{
		"message": "not valid user body",
	})
	return
}
err= u.InsertUser()
if err!=nil{
	println(err.Error())
	context.JSON(http.StatusBadRequest, gin.H{
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
		err:=context.BindJSON(&newUserInfo)
		if err!=nil	{
			context.JSON(http.StatusBadRequest, gin.H{
			"message": "not valid user body",
		})
return
		}
err=newUserInfo.EditUser()
if err!=nil{
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
		strid:=context.Param("id")
		id,err := strconv.Atoi(strid)
		if err!=nil{
context.JSON(http.StatusBadRequest, gin.H{
			"message": "not valid id",
		})
return
		}
		err=dropUser(id)
		if err!=nil{
context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
return
		}else{
context.JSON(http.StatusOK, gin.H{
	"status": "Sucsess",
})
		}
	}