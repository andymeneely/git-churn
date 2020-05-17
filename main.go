package main

import (
	"github.com/andymeneely/git-churn/controllers"
	_ "github.com/andymeneely/git-churn/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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

	r := gin.Default()
	r.GET("/churn-metrics/file", controllers.GetFileChurnMetrics)
	r.GET("/churn-metrics/aggr", controllers.GetAggrChurnMetrics)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()

	//For executing the concurrent go routines in the program in parallel
	//numcpu := runtime.NumCPU()
	//runtime.GOMAXPROCS(numcpu)
	//cmd.Execute()
}
