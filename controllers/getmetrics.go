package controllers

import (
	"github.com/andymeneely/git-churn/gitfuncs"
	metrics "github.com/andymeneely/git-churn/matrics"
	"github.com/andymeneely/git-churn/print"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetMetrics Swagger test
// @Summary Show an account
// @Description get string by ID
// @Tags accounts
// @Produce  json
// @repoUrl id path int true "Account ID"
// @Success 200 {object} metrics.ChurnMetrics
// @Router /accounts/{id} [get]
func GetChurnMetrics(c *gin.Context) {
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
