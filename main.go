package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nitishkumar71/gcp-cleaner/pkg/services"
)

func main() {
	fmt.Println("This is GCP Cleaner!!")
	r := setupRouter()

	port := os.Getenv("PORT")
	fmt.Println("Port: ", port)
	if port == "" {
		port = "8080"
	}
	port = ":" + port
	r.Run(port)
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	setRoutes(r)

	return r
}

func setRoutes(r *gin.Engine) {
	r.POST("/cloudrun", services.DeleteCloudRunRevisionsPost)
}
