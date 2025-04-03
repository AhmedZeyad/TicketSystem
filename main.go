package main
import(
	"github.com/gin-gonic/gin"
)
func main() {
println("Hello World")
router := gin.Default()
router.GET("/test",func(ctx *gin.Context) {
	ctx.JSON(200,gin.H{
		"message":"Hello From API",
	})
})
router.Run(":9090")

}
