package main

import (
	"github.com/andymeneely/git-churn/cmd"
	"github.com/andymeneely/git-churn/controllers"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"runtime"
)

// @title Swagger Git-Churn API
// @version 1.0
// @description APIs to get Git-churn metrices.

// @contact.name Dr. Andy Meneely
// @contact.url http://www.github.com/andymeneely
// @contact.email andy@rit.edu

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host github.com
func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	r.GET("/churn-metrics", controllers.GetChurnMetrics)

	//url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()

	//For executing the concurrent go routines in the program in parallel
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
	cmd.Execute()
}
