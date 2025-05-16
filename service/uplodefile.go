package service

import (
	"TicketSystem/engine"
	"fmt"
	"log"
	"net/http"
	"strings"
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
var isAllawed = map[string]bool{
    "jpg":  true,
    "jpeg": true,
    "png":  true,
    "pdf":  true,
    "docx": false,
}
var  MaxFileSize  int64 =  5000000

	userID := context.MustGet("userId")
	if userID ==0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "not valid user id",
			"value":   userID,
			"sucress": false,
		})
		return
	}
	// check if the post request has the file
	file, err := context.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"sucress": false,
		})
		return
	}
	fileExtention:= strings.Split(file.Filename, ".")
	fmt.Println(fileExtention[len(fileExtention)-1])
fmt.Println(file.Size)
// the file.size i get the size in byte
if file.Size > MaxFileSize {
	context.JSON(http.StatusBadRequest,gin.H{
		"messege":"file to big the max size in 5 MB",
		"success":false,
	})
	return
}
// check the file extention is allowed to uplode
	if !isAllawed[ fileExtention[len(fileExtention)-1] ] {
		 context.JSON(http.StatusBadRequest, gin.H{
			"message": "file type not allowed",
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
/*
chlenges 
File Overwriting – Uploading a file with an existing name will overwrite the original.

done Unsupported or Dangerous File Types – Users may upload .exe, .bat, or script files.

done Large File Crashes or Memory Exhaustion – Uploading big files may exceed server memory or disk limits.

todo Missing Upload Folder – The upload may fail if the target folder doesn’t exist.

todo No Authentication – Anyone can upload without any authorization check.

todo Filename Injection (XSS) – Filenames like image<script>.jpg can cause XSS in web views.

todo Poor Error Handling – Failing to return clear or consistent error messages.

todo No Duplicate Handling – You don't detect or handle duplicate file names or content.

todo Silent Failures – Errors may happen without clear logs or messages if not properly handled.
*/