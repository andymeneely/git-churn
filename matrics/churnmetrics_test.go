package metrics

import (
	"github.com/andymeneely/git-churn/gitfuncs"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"runtime"
	"testing"
)

var churnRepo *git.Repository
var projectRepo *git.Repository

func init() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)

	churnRepo = gitfuncs.GetRepo("https://github.com/andymeneely/git-churn.git")
	//churnRepo = gitfuncs.GetRepo("/Users/raj.g/Documents/RA/git-churn")
	projectRepo = gitfuncs.GetRepo("https://github.com/ashishgalagali/SWEN610-project.git")
	//projectRepo = gitfuncs.GetRepo("/Users/raj.g/Documents/projects/SWEN610-project")
}

func TestGetChurnMetricsAll(t *testing.T) {
	churnmetricsArr, _ := GetChurnMetrics(projectRepo, "c800ce62fc8a10d5fe69adb283f06296820522c1", "", "", true, false, false)
	assert := assert.New(t)
	assert.Equal("c800ce62fc8a10d5fe69adb283f06296820522c1", churnmetricsArr.BaseCommitId)
	assert.Equal(21, len(churnmetricsArr.CommitDetails))
	commitDetail := churnmetricsArr.CommitDetails[0]
	//assert.Equal("d77e8cd611bb63178b40da8f7c4cf1900257aff0", commitDetail.CommitId)
	//assert.Equal("ag1016@rit.edu", commitDetail.CommitAuthor)
	//assert.Equal("2019-11-28 16:02:16 -0500 EST", commitDetail.DateTime)
	//assert.Equal("Merge pull request #6 from ashishgalagali/kirtana\n\nimplemented user wait and transition to game state", commitDetail.CommitMessage)
	//assert.Equal(15, len(commitDetail.ChurnMetrics))
	churnMetric := commitDetail.ChurnMetrics[0]
	assert.NotEmpty(churnMetric.FilePath)
	assert.Greater(churnMetric.DeletedLinesCount, 0)
	assert.Greater(churnMetric.SelfChurnCount, -1)
	assert.Greater(churnMetric.InteractiveChurnCount, -1)
	assert.Greater(len(churnMetric.ChurnDetails), 0)
}

func TestGetChurnMetricsAllRange(t *testing.T) {
	//https://github.com/ashishgalagali/SWEN610-project/compare/5a2bf9f4da3de056dde3d9a9c18859de124d2602...c800ce62fc8a10d5fe69adb283f06296820522c1
	churnmetrics, _ := GetChurnMetrics(projectRepo, "c800ce62fc8a10d5fe69adb283f06296820522c1", "", "4513ea9d398b549c3e00862c866adbc4566d43de", false, false, false)
	//out, err := json.Marshal(churnmetrics)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(out))
	//os.RemoveAll("churn-metrics")

	assert := assert.New(t)
	assert.Equal("c800ce62fc8a10d5fe69adb283f06296820522c1", churnmetrics.BaseCommitId)

	assert.Equal(8, len(churnmetrics.CommitDetails))
	commitDetail := churnmetrics.CommitDetails[0]
	assert.Equal("d77e8cd611bb63178b40da8f7c4cf1900257aff0", commitDetail.CommitId)
	assert.Equal("ag1016@rit.edu", commitDetail.CommitAuthor)
	//assert.Equal("2019-11-28 16:02:16 -0500 EST", commitDetail.DateTime)
	assert.Equal("Merge pull request #6 from ashishgalagali/kirtana\n\nimplemented user wait and transition to game state", commitDetail.CommitMessage)
	assert.Equal(15, len(commitDetail.ChurnMetrics))
	churnMetric := commitDetail.ChurnMetrics[0]
	assert.NotEmpty(churnMetric.FilePath)
	assert.Greater(churnMetric.DeletedLinesCount, 0)
	assert.Greater(churnMetric.SelfChurnCount, -1)
	assert.Greater(churnMetric.InteractiveChurnCount, -1)
	assert.Greater(len(churnMetric.ChurnDetails), 0)

	//assert.Equal("src/main/java/com/webcheckers/ui/WebServer.java", churnmetrics.FilePath)
	//assert.Equal(4, churnmetrics.DeletedLinesCount)
	//assert.Equal(2, churnmetrics.InteractiveChurnCount)
	//assert.Equal(2, churnmetrics.SelfChurnCount)
	//assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	//assert.Equal(3, len(churnmetrics.ChurnDetails))
	//assert.Equal("ks3057@rit.edu", churnmetrics.ChurnDetails["b742aaf3e500712668d6f76c9736637436ee695e"])
	//assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["979fe965043d49814c2fb7e7f5bae3461911b88b"])
	//assert.Equal(4, churnmetrics.FileDiffMetrics.Deletions)
	//assert.Equal(19, churnmetrics.FileDiffMetrics.Insertions)
}

func TestGetChurnMetricsAllRangeRev(t *testing.T) {

	churnmetrics, _ := GetChurnMetrics(projectRepo, "4513ea9d398b549c3e00862c866adbc4566d43de", "", "c800ce62fc8a10d5fe69adb283f06296820522c1", false, false, false)

	assert := assert.New(t)
	assert.Equal("4513ea9d398b549c3e00862c866adbc4566d43de", churnmetrics.BaseCommitId)

	assert.Equal(8, len(churnmetrics.CommitDetails))
	commitDetail := churnmetrics.CommitDetails[0]
	assert.Equal("d77e8cd611bb63178b40da8f7c4cf1900257aff0", commitDetail.CommitId)
	assert.Equal("ag1016@rit.edu", commitDetail.CommitAuthor)
	//assert.Equal("2019-11-28 16:02:16 -0500 EST", commitDetail.DateTime)
	assert.Equal("Merge pull request #6 from ashishgalagali/kirtana\n\nimplemented user wait and transition to game state", commitDetail.CommitMessage)
	assert.Equal(15, len(commitDetail.ChurnMetrics))
	churnMetric := commitDetail.ChurnMetrics[0]
	assert.NotEmpty(churnMetric.FilePath)
	assert.Greater(churnMetric.DeletedLinesCount, 0)
	assert.Greater(churnMetric.SelfChurnCount, -1)
	assert.Greater(churnMetric.InteractiveChurnCount, -1)
	assert.Greater(len(churnMetric.ChurnDetails), 0)
	//	churnmetrics, _ := GetChurnMetrics(projectRepo, "5a2bf9f4da3de056dde3d9a9c18859de124d2602", "src/main/java/com/webcheckers/ui/WebServer.java", "c800ce62fc8a10d5fe69adb283f06296820522c1")
	//	assert := assert.New(t)
	//	assert.Equal("src/main/java/com/webcheckers/ui/WebServer.java", churnmetrics.FilePath)
	//	assert.Equal(19, churnmetrics.DeletedLinesCount)
	//	assert.Equal(4, churnmetrics.InteractiveChurnCount)
	//	assert.Equal(15, churnmetrics.SelfChurnCount)
	//	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	//	assert.Equal(5, len(churnmetrics.ChurnDetails))
	//	assert.Equal("ks3057@rit.edu", churnmetrics.ChurnDetails["9708c9a9da36928fd0b7143c74aa61694999fe5d"])
	//	assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["979fe965043d49814c2fb7e7f5bae3461911b88b"])
	//	assert.Equal(19, churnmetrics.FileDiffMetrics.Deletions)
	//	assert.Equal(4, churnmetrics.FileDiffMetrics.Insertions)
}

//
////TODO: This error is due to an error in the Blame class of go-git. Try to find a hack
////func TestGetChurnMetricsAllFailing(t *testing.T) {
////	repo := gitfuncs.Checkout("https://github.com/ashishgalagali/SWEN610-project", "c800ce62fc8a10d5fe69adb283f06296820522c1")
////	churnmetrics, _ := GetChurnMetrics(repo, "src/main/java/com/webcheckers/Application.java")
////	assert := assert.New(t)
////	assert.Equal("src/main/java/com/webcheckers/Application.java", churnmetrics.FilePath)
////	assert.Equal(2, churnmetrics.DeletedLinesCount)
////	assert.Equal(2, churnmetrics.FileDiffMetrics.Deletions)
////}
//
func TestGetChurnMetricsInteractiveChurn(t *testing.T) {
	churnmetrics, _ := GetChurnMetrics(churnRepo, "99992110e402f26ca9162f43c0e5a97b1278068a", "README.md", "", true, false, false)

	assert := assert.New(t)
	assert.Equal("99992110e402f26ca9162f43c0e5a97b1278068a", churnmetrics.BaseCommitId)

	assert.Equal(36, len(churnmetrics.CommitDetails))
	commitDetail := churnmetrics.CommitDetails[0]
	assert.Equal("180ec07da5d7a415b48fd3d9f7d5c9dd2925780e", commitDetail.CommitId)
	assert.Equal("ashishgalagali@gmail.com", commitDetail.CommitAuthor)
	//assert.Equal("2020-03-28 00:59:14 -0400 -0400", commitDetail.DateTime)
	assert.Equal("Merge pull request #19 from andymeneely/diffMetrics\n\nGetting git diff metrics for a given commit and file", commitDetail.CommitMessage)
	assert.Equal(1, len(commitDetail.ChurnMetrics))
	churnmetric := commitDetail.ChurnMetrics[0]
	assert.Equal("README.md", churnmetric.FilePath)
	assert.Equal(5, churnmetric.DeletedLinesCount)
	assert.Equal(5, churnmetric.InteractiveChurnCount)
	assert.Equal(0, churnmetric.SelfChurnCount)
	assert.Equal(1, len(churnmetric.ChurnDetails))
	assert.Equal("andy@se.rit.edu", churnmetric.ChurnDetails["79caa008ba1f9d06b34b4acc7c03d7fade185a63"])
}

//
//func TestGetChurnMetricsInteractiveChurnNewFile(t *testing.T) {
//	_, err := GetChurnMetrics(churnRepo, "99992110e402f26ca9162f43c0e5a97b1278068a", "cmd/root.go", "")
//
//	assert := assert.New(t)
//	assert.Equal("The specified file was a new file added in this commit. Hence, churn can't be calculated.", err.Error())
//}
//
func TestGetChurnMetricsSelfChurn(t *testing.T) {
	churnmetrics, _ := GetChurnMetrics(churnRepo, "c0263662b2172b3df51ae39f8075dd010573ab6b", "matrics/diffmetrics_test.go", "", true, false, false)
	assert := assert.New(t)
	assert.Equal(6, len(churnmetrics.CommitDetails))
	commitDetail := churnmetrics.CommitDetails[0]
	assert.Equal("a8b24a74bae39a941186e11969f70058a351327d", commitDetail.CommitId)
	assert.Equal("ashishgalagali@gmail.com", commitDetail.CommitAuthor)
	//assert.Equal("2020-03-30 21:32:45 -0400 -0400", commitDetail.DateTime)
	assert.Equal("Merge pull request #22 from andymeneely/cli\n\nAdded boolean flag to exclude whitespace+ Updated README", commitDetail.CommitMessage)
	assert.Equal(1, len(commitDetail.ChurnMetrics))
	churnmetric := commitDetail.ChurnMetrics[0]
	assert.Equal(65, churnmetric.DeletedLinesCount)
	assert.Equal(0, churnmetric.InteractiveChurnCount)
	assert.Equal(65, churnmetric.SelfChurnCount)
	assert.Equal(1, len(churnmetric.ChurnDetails))
	assert.Equal("ashishgalagali@gmail.com", churnmetric.ChurnDetails["3854e533318df4f5bb9a059c76ddd8bb2464a620"])
}

//
func TestGetChurnMetricsWhitespaceExcludedAll(t *testing.T) {
	churnmetrics, _ := GetChurnMetrics(projectRepo, "c800ce62fc8a10d5fe69adb283f06296820522c1", "src/main/java/com/webcheckers/ui/WebServer.java", "", false, false, false)
	assert := assert.New(t)

	assert.Equal(20, len(churnmetrics.CommitDetails))
	commitDetail := churnmetrics.CommitDetails[1]
	assert.Equal("cef4dbea729fac483b43e130271c9e6efe93df33", commitDetail.CommitId)
	assert.Equal("ks3057@rit.edu", commitDetail.CommitAuthor)
	//assert.Equal("2019-11-28 02:33:55 -0500 EST", commitDetail.DateTime)
	assert.Equal("implemented user wait and transition to game state\n", commitDetail.CommitMessage)
	assert.Equal(1, len(commitDetail.ChurnMetrics))
	churnmetric := commitDetail.ChurnMetrics[0]
	assert.Equal("src/main/java/com/webcheckers/ui/WebServer.java", churnmetric.FilePath)
	assert.Equal(12, churnmetric.DeletedLinesCount)
	assert.Equal(10, churnmetric.SelfChurnCount)
	assert.Equal(2, churnmetric.InteractiveChurnCount)
	assert.Equal(5, len(churnmetric.ChurnDetails))
	assert.Equal("ks3057@rit.edu", churnmetric.ChurnDetails["9708c9a9da36928fd0b7143c74aa61694999fe5d"])
	assert.Equal("ashishgalagali@gmail.com", churnmetric.ChurnDetails["16123ab124432a058ed29e7d8fb2df52c310363b"])
}

func TestGetChurnMetricsWhitespaceExcludedAllRange(t *testing.T) {
	// https://github.com/ashishgalagali/SWEN610-project/compare/16c75b486a039bc34fcc5ac1ddad717d8bb49c01...7368d5fcb7eec950161ed9d13b55caf5961326b6?diff=split
	churnmetrics, _ := GetChurnMetrics(projectRepo, "7368d5fcb7eec950161ed9d13b55caf5961326b6", "README.md", "16c75b486a039bc34fcc5ac1ddad717d8bb49c01", false, false, false)
	assert := assert.New(t)
	assert.Equal(2, len(churnmetrics.CommitDetails))
	commitDetail := churnmetrics.CommitDetails[1]
	assert.Equal("8e6f09133b61c6eeb83d4e529c14c3754c286774", commitDetail.CommitId)
	assert.Equal("ashishgalagali@gmail.com", commitDetail.CommitAuthor)
	//assert.Equal("2019-11-28 02:33:55 -0500 EST", commitDetail.DateTime)
	assert.Equal("Updating readme\n", commitDetail.CommitMessage)
	assert.Equal(1, len(commitDetail.ChurnMetrics))
	churnmetric := commitDetail.ChurnMetrics[0]
	assert.Equal("README.md", churnmetric.FilePath)
	assert.Equal(4, churnmetric.DeletedLinesCount)
	assert.Equal(4, churnmetric.SelfChurnCount)
	assert.Equal(0, churnmetric.InteractiveChurnCount)
	//assert.Equal("ashishgalagali@gmail.com", churnmetric.CommitAuthor)
	assert.Equal(1, len(churnmetric.ChurnDetails))
	assert.Equal("ashishgalagali@gmail.com", churnmetric.ChurnDetails["979fe965043d49814c2fb7e7f5bae3461911b88b"])
}

//assert.Equal(8, churnmetrics.FileDiffMetrics.Insertions)
//assert.Equal(13, churnmetrics.FileDiffMetrics.Deletions)
//assert.Equal(24, churnmetrics.FileDiffMetrics.LinesBefore)
//assert.Equal(19, churnmetrics.FileDiffMetrics.LinesAfter)
//}

func TestGetChurnMetricsWhitespaceExcludedAllRangeRev(t *testing.T) {
	churnmetrics, _ := GetChurnMetrics(projectRepo, "16c75b486a039bc34fcc5ac1ddad717d8bb49c01", "README.md", "7368d5fcb7eec950161ed9d13b55caf5961326b6", false, false, false)
	assert := assert.New(t)

	assert.Equal(2, len(churnmetrics.CommitDetails))
	commitDetail := churnmetrics.CommitDetails[1]
	assert.Equal("8e6f09133b61c6eeb83d4e529c14c3754c286774", commitDetail.CommitId)
	assert.Equal("ashishgalagali@gmail.com", commitDetail.CommitAuthor)
	//assert.Equal("2019-12-04 03:16:33 -0500 EST", commitDetail.DateTime)
	assert.Equal("Updating readme\n", commitDetail.CommitMessage)
	assert.Equal(1, len(commitDetail.ChurnMetrics))
	churnmetric := commitDetail.ChurnMetrics[0]
	assert.Equal("README.md", churnmetric.FilePath)
	assert.Equal(6, churnmetric.DeletedLinesCount)
	assert.Equal(6, churnmetric.SelfChurnCount)
	assert.Equal(0, churnmetric.InteractiveChurnCount)
	assert.Equal(1, len(churnmetric.ChurnDetails))
	assert.Equal("ashishgalagali@gmail.com", churnmetric.ChurnDetails["8e6f09133b61c6eeb83d4e529c14c3754c286774"])
	//assert.Equal(13, churnmetrics.FileDiffMetrics.Insertions)
	//assert.Equal(8, churnmetrics.FileDiffMetrics.Deletions)
	//assert.Equal(19, churnmetrics.FileDiffMetrics.LinesBefore)
	//assert.Equal(24, churnmetrics.FileDiffMetrics.LinesAfter)
}

func TestGetChurnMetricsWhitespaceExcludedInteractiveChurn(t *testing.T) {
	churnmetrics, _ := GetChurnMetrics(churnRepo, "99992110e402f26ca9162f43c0e5a97b1278068a", "README.md", "", false, false, false)
	assert := assert.New(t)
	assert.Equal(36, len(churnmetrics.CommitDetails))
	commitDetail := churnmetrics.CommitDetails[0]
	assert.Equal("180ec07da5d7a415b48fd3d9f7d5c9dd2925780e", commitDetail.CommitId)
	assert.Equal("ashishgalagali@gmail.com", commitDetail.CommitAuthor)
	//assert.Equal("2020-03-28 00:59:14 -0400 -0400", commitDetail.DateTime)
	assert.Equal("Merge pull request #19 from andymeneely/diffMetrics\n\nGetting git diff metrics for a given commit and file", commitDetail.CommitMessage)
	assert.Equal(1, len(commitDetail.ChurnMetrics))
	churnmetric := commitDetail.ChurnMetrics[0]
	assert.Equal("README.md", churnmetric.FilePath)
	assert.Equal(3, churnmetric.DeletedLinesCount)
	assert.Equal(3, churnmetric.InteractiveChurnCount)
	assert.Equal(0, churnmetric.SelfChurnCount)
	assert.Equal(1, len(churnmetric.ChurnDetails))
	assert.Equal("andy@se.rit.edu", churnmetric.ChurnDetails["79caa008ba1f9d06b34b4acc7c03d7fade185a63"])
}

func TestAggrChurnMetricsWithWhitespace(t *testing.T) {
	churnmetrics := AggrChurnMetrics(churnRepo, "99992110e402f26ca9162f43c0e5a97b1278068a", "", "commit", true, false, false)
	assert := assert.New(t)
	assert.Equal(39, len(churnmetrics.(AggChurnMetricsOutput).AggCommitDetails))
	commitDetail := churnmetrics.(AggChurnMetricsOutput).AggCommitDetails[0]
	assert.Equal("180ec07da5d7a415b48fd3d9f7d5c9dd2925780e", commitDetail.CommitId)
	assert.Equal("ashishgalagali@gmail.com", commitDetail.CommitAuthor)
	//assert.Equal("2020-03-28 00:59:14 -0400 -0400", commitDetail.DateTime)
	assert.Equal("Merge pull request #19 from andymeneely/diffMetrics\n\nGetting git diff metrics for a given commit and file", commitDetail.CommitMessage)
	churnmetric := commitDetail.AggChurnMetrics
	assert.Equal(29, churnmetric.TotalDeletedLinesCount)
	assert.Equal(5, churnmetric.TotalInteractiveChurnCount)
	assert.Equal(24, churnmetric.TotalSelfChurnCount)
	assert.Equal(4, churnmetric.FilesCount)
}

func TestAggrChurnMetricsWhitespaceExcluded(t *testing.T) {
	//churnmetrics := AggrChurnMetricsWhitespaceExcluded(churnRepo, "99992110e402f26ca9162f43c0e5a97b1278068a")
	churnmetrics := AggrChurnMetrics(churnRepo, "99992110e402f26ca9162f43c0e5a97b1278068a", "", "commit", false, false, false)
	//out, err := json.Marshal(churnmetrics)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(out))
	//assert := assert.New(t)

	assert := assert.New(t)
	assert.Equal(39, len(churnmetrics.(AggChurnMetricsOutput).AggCommitDetails))
	commitDetail := churnmetrics.(AggChurnMetricsOutput).AggCommitDetails[0]
	assert.Equal("180ec07da5d7a415b48fd3d9f7d5c9dd2925780e", commitDetail.CommitId)
	assert.Equal("ashishgalagali@gmail.com", commitDetail.CommitAuthor)
	//assert.Equal("2020-03-28 00:59:14 -0400 -0400", commitDetail.DateTime)
	assert.Equal("Merge pull request #19 from andymeneely/diffMetrics\n\nGetting git diff metrics for a given commit and file", commitDetail.CommitMessage)
	churnmetric := commitDetail.AggChurnMetrics
	assert.Equal(25, churnmetric.TotalDeletedLinesCount)
	assert.Equal(3, churnmetric.TotalInteractiveChurnCount)
	assert.Equal(22, churnmetric.TotalSelfChurnCount)
	assert.Equal(4, churnmetric.FilesCount)
}

//func TestChurnMetricsWithWhitespaceProfiling(t *testing.T) {
//	esRepo := gitfuncs.GetRepo("/Users/raj.g/Documents/Git_projects/elasticsearch")
//	churnmetrics, _ := GetChurnMetrics(esRepo, "e33a0dfe77ad530db99bdd203434d16b23999be6", "", "91cc417aba84fce2b9b8a42794ee2411708b6c71", true, false, true)
//	out, err := json.Marshal(churnmetrics)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(string(out))
//}
