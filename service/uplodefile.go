package service

import (
	"TicketSystem/engine"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var basePath = "assets"

func UploadRoutes(rg *gin.RouterGroup) {
	api := rg.Group("assets/")
	api.POST("upload", engine.CheckAuth, UploadFile)
	api.GET("download/:filename", engine.CheckAuth, GetFile)
}
func UploadFile(context *gin.Context) {
	var isAllowed = map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"pdf":  true,
		"docx": false,
	}

	// userID := context.MustGet("userId").(int)

	// check if the post request has the file
	file, err := context.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	fileParts := strings.Split(file.Filename, ".")
	if len(fileParts) < 2 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "file name not valid",
			"success": false,
		})
		return
	}

	fileExtention := strings.ToLower(fileParts[len(fileParts)-1])
	fmt.Println(fileExtention)
	fmt.Println(file.Size)

	// the file.size i get the size in byte
	if file.Size > (5 << 20) {
		context.JSON(http.StatusBadRequest, gin.H{
			"messege": "file to big the max size in 5 MB",
			"success": false,
		})
		return
	}
	// check the file extention is allowed to uplode
	if !isAllowed[fileExtention] {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "file type not allowed",
			"success": false,
		})
		return

	}

	// Upload the file to specific dst.
	// make the file name unique

	var uniqefname = fmt.Sprintf("%d_%d.%s", time.Now().Unix(), uuid.New().ID(), fileExtention)
	// todo store name in db
	// for save paht joining
fullPath := filepath.Join(basePath, uniqefname)
	if err := context.SaveUploadedFile(file, fullPath); err != nil {
		context.String(http.StatusBadRequest, fmt.Sprintf("upload error: %s", err.Error()))
		return
	}
log.Printf("user uploded file %s , it saved as %s",file.Filename,fullPath)
	context.JSON(http.StatusOK, gin.H{
		"message": "file uploded successfly",
		"file_name":uniqefname,
		"success": true,
	})
}

/*
chlenges
File Overwriting – Uploading a file with an existing name will overwrite the original.

done Unsupported or Dangerous File Types – Users may upload .exe, .bat, or script files.

done Large File Crashes or Memory Exhaustion – Uploading big files may exceed server memory or disk limits.

done Missing Upload Folder – The upload may fail if the target folder doesn’t exist.

done No Authentication – Anyone can upload without any authorization check.

done Filename Injection (XSS) – Filenames like image<script>.jpg can cause XSS in web views.
dont user the file name of the user
done Poor Error Handling – Failing to return clear or consistent error messages.

done No Duplicate Handling – You don't detect or handle duplicate file names or content.

done Silent Failures – Errors may happen without clear logs or messages if not properly handled.
*/

func GetFile(context *gin.Context) {

	filename := context.Param("filename")
	fullPath :=filepath.Join( basePath , filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "file not found",
			"success": false,
		})
		return
	}

	context.Header("Content-Description", "File Transfer")
context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))

mimeType := mime.TypeByExtension(filepath.Ext(filename))
if mimeType != "" {
	context.Header("Content-Type", mimeType)
} else {
	context.Header("Content-Type", "application/octet-stream")
}


	context.Header("Content-Type", "application/octet-stream")

	context.File(fullPath)
	log.Printf("user downloded file %s",filename)

}


/*
chlenges
done  Path Traversal Attack – Users can request files outside the intended folder (e.g., ../../etc/passwd).

done ⚠️ File Not Found – If the file doesn’t exist, it may cause a server error or crash.

⚠️ Incorrect Content-Type – The downloaded file may not be handled correctly by the client if the wrong MIME type is set.

❌ Internal Filename Exposure – Exposing raw filenames like UUIDs may leak system info or structure.

⚠️ Large File Memory Usage – Loading large files into memory before sending can crash the server.

❌ Missing Authorization – Any logged-in user may access files they shouldn’t.

🧪 Poor Error Handling – Postman or clients may get vague or unclear errors when something fails.

❌ Lack of Logging – No tracking of who downloads what, which can be important for auditing.
*/
