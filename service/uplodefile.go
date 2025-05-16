package service

import (
	"TicketSystem/engine"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var basePath = "./assets/"

func UploadRoutes(rg *gin.RouterGroup) {
	api := rg.Group("assets/")
	api.POST("upload", engine.CheckAuth, UploadFile)
	api.GET("download/:filename", engine.CheckAuth, GetFile)
}
func UploadFile(context *gin.Context) {
	userID := context.MustGet("userId")
	// the  file is come from  body
	// context.String(http.StatusBadRequest, fmt.Sprintf("upload error: %s"))
	// return
	if userID ==0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "not valid user id",
			"value":   userID,
			"sucress": false,
		})
		return
	}
	file, err := context.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"sucress": false,
		})
		return
	}
	log.Println(file.Filename)

	// Upload the file to specific dst.
	// make the file name unique
	var uniqefname = fmt.Sprintf("%d_%f_%s", time.Now().Unix(), userID, file.Filename)
	// todo store name in db

	if err := context.SaveUploadedFile(file, "./assets/"+uniqefname); err != nil {
		context.String(http.StatusBadRequest, fmt.Sprintf("upload error: %s", err.Error()))
		return
	}

	context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", uniqefname))
}
func GetFile(context *gin.Context) {
	filename := context.Param("filename")
	filePath := basePath + filename
	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Disposition", "attachment; filename="+filename)
	context.Header("Content-Type", "application/octet-stream")

	context.File(filePath)

}
