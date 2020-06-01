package metrics

import (
	"github.com/andymeneely/git-churn/gitfuncs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetChurnMetricsAll(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/ashishgalagali/SWEN610-project", "c800ce62fc8a10d5fe69adb283f06296820522c1")
	churnmetrics, _ := GetChurnMetricsWithWhitespace(repo, "src/main/java/com/webcheckers/ui/WebServer.java")
	assert := assert.New(t)
	assert.Equal("src/main/java/com/webcheckers/ui/WebServer.java", churnmetrics.FilePath)
	assert.Equal(13, churnmetrics.DeletedLinesCount)
	assert.Equal(10, churnmetrics.InteractiveChurnCount)
	assert.Equal(3, churnmetrics.SelfChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(5, len(churnmetrics.ChurnDetails))
	assert.Equal("ks3057@rit.edu", churnmetrics.ChurnDetails["9708c9a9da36928fd0b7143c74aa61694999fe5d"])
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["16123ab124432a058ed29e7d8fb2df52c310363b"])
	assert.Equal(13, churnmetrics.FileDiffMetrics.Deletions)
}

//TODO: This error is due to an error in the Blame class of go-git. Try to find a hack
//func TestGetChurnMetricsAllFailing(t *testing.T) {
//	repo := gitfuncs.Checkout("https://github.com/ashishgalagali/SWEN610-project", "c800ce62fc8a10d5fe69adb283f06296820522c1")
//	churnmetrics, _ := GetChurnMetricsWithWhitespace(repo, "src/main/java/com/webcheckers/Application.java")
//	assert := assert.New(t)
//	assert.Equal("src/main/java/com/webcheckers/Application.java", churnmetrics.FilePath)
//	assert.Equal(2, churnmetrics.DeletedLinesCount)
//	assert.Equal(2, churnmetrics.FileDiffMetrics.Deletions)
//}

func TestGetChurnMetricsInteractiveChurn(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/andymeneely/git-churn", "99992110e402f26ca9162f43c0e5a97b1278068a")
	churnmetrics, _ := GetChurnMetricsWithWhitespace(repo, "README.md")
	assert := assert.New(t)
	assert.Equal("README.md", churnmetrics.FilePath)
	assert.Equal(5, churnmetrics.DeletedLinesCount)
	assert.Equal(5, churnmetrics.InteractiveChurnCount)
	assert.Equal(0, churnmetrics.SelfChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(1, len(churnmetrics.ChurnDetails))
	assert.Equal("andy@se.rit.edu", churnmetrics.ChurnDetails["79caa008ba1f9d06b34b4acc7c03d7fade185a63"])
	assert.Equal(5, churnmetrics.FileDiffMetrics.Deletions)
}

func TestGetChurnMetricsInteractiveChurnNewFile(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/andymeneely/git-churn", "99992110e402f26ca9162f43c0e5a97b1278068a")
	_, err := GetChurnMetricsWithWhitespace(repo, "cmd/root.go")
	assert := assert.New(t)
	assert.Equal("The specified file was a new file added in this commit. Hence, churn can't be calculated.", err.Error())
}

func TestGetChurnMetricsSelfChurn(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/andymeneely/git-churn", "c0263662b2172b3df51ae39f8075dd010573ab6b")
	churnmetrics, _ := GetChurnMetricsWithWhitespace(repo, "matrics/diffmetrics_test.go")
	assert := assert.New(t)
	assert.Equal("matrics/diffmetrics_test.go", churnmetrics.FilePath)
	assert.Equal(65, churnmetrics.DeletedLinesCount)
	assert.Equal(0, churnmetrics.InteractiveChurnCount)
	assert.Equal(65, churnmetrics.SelfChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(1, len(churnmetrics.ChurnDetails))
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["3854e533318df4f5bb9a059c76ddd8bb2464a620"])
	assert.Equal(65, churnmetrics.FileDiffMetrics.Deletions)
}

func TestGetChurnMetricsWhitespaceExcludedAll(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/ashishgalagali/SWEN610-project", "c800ce62fc8a10d5fe69adb283f06296820522c1")
	churnmetrics, _ := GetChurnMetricsWhitespaceExcluded(repo, "src/main/java/com/webcheckers/ui/WebServer.java")
	assert := assert.New(t)
	assert.Equal("src/main/java/com/webcheckers/ui/WebServer.java", churnmetrics.FilePath)
	assert.Equal(12, churnmetrics.DeletedLinesCount)
	assert.Equal(10, churnmetrics.InteractiveChurnCount)
	assert.Equal(2, churnmetrics.SelfChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(5, len(churnmetrics.ChurnDetails))
	assert.Equal("ks3057@rit.edu", churnmetrics.ChurnDetails["9708c9a9da36928fd0b7143c74aa61694999fe5d"])
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["16123ab124432a058ed29e7d8fb2df52c310363b"])
	assert.Equal(12, churnmetrics.FileDiffMetrics.Deletions)
}

func TestGetChurnMetricsWhitespaceExcludedInteractiveChurn(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/andymeneely/git-churn", "99992110e402f26ca9162f43c0e5a97b1278068a")
	churnmetrics, _ := GetChurnMetricsWhitespaceExcluded(repo, "README.md")
	assert := assert.New(t)
	assert.Equal("README.md", churnmetrics.FilePath)
	assert.Equal(3, churnmetrics.DeletedLinesCount)
	assert.Equal(3, churnmetrics.InteractiveChurnCount)
	assert.Equal(0, churnmetrics.SelfChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(1, len(churnmetrics.ChurnDetails))
	assert.Equal("andy@se.rit.edu", churnmetrics.ChurnDetails["79caa008ba1f9d06b34b4acc7c03d7fade185a63"])
	assert.Equal(3, churnmetrics.FileDiffMetrics.Deletions)
}

func TestAggrChurnMetricsWithWhitespace(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/andymeneely/git-churn", "99992110e402f26ca9162f43c0e5a97b1278068a")
	churnmetrics := AggrChurnMetricsWithWhitespace(repo)
	assert := assert.New(t)
	assert.Equal(29, churnmetrics.DeletedLinesCount)
	assert.Equal(5, churnmetrics.InteractiveChurnCount)
	assert.Equal(24, churnmetrics.SelfChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(29, churnmetrics.AggrDiffMetrics.Deletions)
}

func TestAggrChurnMetricsWhitespaceExcluded(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/andymeneely/git-churn", "99992110e402f26ca9162f43c0e5a97b1278068a")
	churnmetrics := AggrChurnMetricsWhitespaceExcluded(repo)
	assert := assert.New(t)
	assert.Equal(25, churnmetrics.DeletedLinesCount)
	assert.Equal(3, churnmetrics.InteractiveChurnCount)
	assert.Equal(22, churnmetrics.SelfChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(25, churnmetrics.AggrDiffMetrics.Deletions)
}
