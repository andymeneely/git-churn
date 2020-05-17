package controllers

import (
	"github.com/andymeneely/git-churn/gitfuncs"
	metrics "github.com/andymeneely/git-churn/matrics"
	"github.com/andymeneely/git-churn/print"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Get File Churn-metrics
// @Summary Provides churn details for a given commit
// @Description Provides churn details for a given commit and for the specified file
// @Produce  json
// @Param repoUrl query string true "repo Url"
// @Param commitId query string true "commit Id"
// @Param whitespace query string false "Should whitespaces be considered?"  default(true)
// @Param filepath query string true "file path"
// @Success 200 {object} metrics.FileChurnMetrics
// @Router /churn-metrics/file [get]
func GetFileChurnMetrics(c *gin.Context) {
	repoUrl := c.Query("repoUrl")
	commitId := c.Query("commitId")
	//TODO: handle error
	whitespace, _ := strconv.ParseBool(c.DefaultQuery("whitespace", "true"))
	filepath := c.DefaultQuery("filepath", "")
	var churnMetrics interface{}
	var err error
	repo := gitfuncs.Checkout(repoUrl, commitId)
	if whitespace {
		if filepath != "" {
			churnMetrics, err = metrics.GetChurnMetricsWithWhitespace(repo, filepath)
		} else {
			churnMetrics = metrics.AggrChurnMetricsWithWhitespace(repo)
		}
	} else {
		if filepath != "" {
			churnMetrics, err = metrics.GetChurnMetricsWhitespaceExcluded(repo, filepath)
		} else {
			churnMetrics = metrics.AggrChurnMetricsWhitespaceExcluded(repo)
		}
		print.CheckIfError(err)
	}
	//fmt.Println(fmt.Sprintf("%v", churnMetrics))
	//out, err := json.Marshal(churnMetrics)
	c.JSON(http.StatusOK, churnMetrics)
}

// Get Aggregate Churn-metrics
// @Summary Provides Aggregate churn details for a given commit
// @Description Provides Aggregated churn details of all the files for a given commit
// @Produce  json
// @Param repoUrl query string true "repo Url"
// @Param commitId query string true "commit Id"
// @Param whitespace query string false "Should whitespaces be considered?"  default(true)
// @Success 200 {object} metrics.AggrChurMetrics
// @Router /churn-metrics/aggr [get]
func GetAggrChurnMetrics(c *gin.Context) {
	repoUrl := c.Query("repoUrl")
	commitId := c.Query("commitId")
	//TODO: handle error
	whitespace, _ := strconv.ParseBool(c.DefaultQuery("whitespace", "true"))
	var churnMetrics interface{}
	repo := gitfuncs.Checkout(repoUrl, commitId)
	if whitespace {
		churnMetrics = metrics.AggrChurnMetricsWithWhitespace(repo)
	} else {
		churnMetrics = metrics.AggrChurnMetricsWhitespaceExcluded(repo)
	}
	//fmt.Println(fmt.Sprintf("%v", churnMetrics))
	//out, err := json.Marshal(churnMetrics)
	c.JSON(http.StatusOK, churnMetrics)
}
