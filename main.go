package main

import (
	"github.com/andymeneely/git-churn/controllers"
	_ "github.com/andymeneely/git-churn/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

var r *gin.Engine

// @title Swagger Git-Churn API
// @version 1.0
// @description APIs to get Git-churn metrices.
// @contact.name Dr. Andy Meneely
// @contact.url http://www.github.com/andymeneely
// @contact.email axmvse@rit.edu
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {

	r = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	r.LoadHTMLGlob("templates/*")
	// Initialize the routes
	initializeRoutes()
	r.GET("/churn-metrics/file", controllers.GetFileChurnMetrics)
	r.GET("/churn-metrics/aggr", controllers.GetAggrChurnMetrics)
	r.POST("/v1/churn-metrics/aggr", controllers.GetAggrChurnMetricsV1)
	r.POST("/v1/churn-metrics/file", controllers.GetFileChurnMetricsV1)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()

	//For executing the concurrent go routines in the program in parallel
	//numcpu := runtime.NumCPU()
	//runtime.GOMAXPROCS(numcpu)
	//cmd.Execute()
}

func initializeRoutes() {
	// Handle the index route
	r.GET("/", showIndexPage)
}

func showIndexPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Home Page"}, "index.html")
}

func render(c *gin.Context, data gin.H, templateName string) {
	c.HTML(http.StatusOK, templateName, data)
}
