package metrics

import (
	"github.com/andymeneely/git-churn/gitfuncs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetChurnMetricsAll(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/ashishgalagali/SWEN610-project", "c800ce62fc8a10d5fe69adb283f06296820522c1")
	churnmetrics, _ := GetChurnMetricsWithWhitespace(repo, "src/main/java/com/webcheckers/ui/WebServer.java", "")
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

func TestGetChurnMetricsAllRange(t *testing.T) {
	//https://github.com/ashishgalagali/SWEN610-project/compare/5a2bf9f4da3de056dde3d9a9c18859de124d2602...c800ce62fc8a10d5fe69adb283f06296820522c1
	repo := gitfuncs.Checkout("https://github.com/ashishgalagali/SWEN610-project", "c800ce62fc8a10d5fe69adb283f06296820522c1")
	churnmetrics, _ := GetChurnMetricsWithWhitespace(repo, "src/main/java/com/webcheckers/ui/WebServer.java", "5a2bf9f4da3de056dde3d9a9c18859de124d2602")
	assert := assert.New(t)
	assert.Equal("src/main/java/com/webcheckers/ui/WebServer.java", churnmetrics.FilePath)
	assert.Equal(4, churnmetrics.DeletedLinesCount)
	assert.Equal(2, churnmetrics.InteractiveChurnCount)
	assert.Equal(2, churnmetrics.SelfChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(3, len(churnmetrics.ChurnDetails))
	assert.Equal("ks3057@rit.edu", churnmetrics.ChurnDetails["b742aaf3e500712668d6f76c9736637436ee695e"])
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["979fe965043d49814c2fb7e7f5bae3461911b88b"])
	assert.Equal(4, churnmetrics.FileDiffMetrics.Deletions)
	assert.Equal(19, churnmetrics.FileDiffMetrics.Insertions)

}
func TestGetChurnMetricsAllRangeRev(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/ashishgalagali/SWEN610-project", "5a2bf9f4da3de056dde3d9a9c18859de124d2602")
	churnmetrics, _ := GetChurnMetricsWithWhitespace(repo, "src/main/java/com/webcheckers/ui/WebServer.java", "c800ce62fc8a10d5fe69adb283f06296820522c1")
	assert := assert.New(t)
	assert.Equal("src/main/java/com/webcheckers/ui/WebServer.java", churnmetrics.FilePath)
	assert.Equal(19, churnmetrics.DeletedLinesCount)
	assert.Equal(4, churnmetrics.InteractiveChurnCount)
	assert.Equal(15, churnmetrics.SelfChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(5, len(churnmetrics.ChurnDetails))
	assert.Equal("ks3057@rit.edu", churnmetrics.ChurnDetails["9708c9a9da36928fd0b7143c74aa61694999fe5d"])
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["979fe965043d49814c2fb7e7f5bae3461911b88b"])
	assert.Equal(19, churnmetrics.FileDiffMetrics.Deletions)
	assert.Equal(4, churnmetrics.FileDiffMetrics.Insertions)
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
	churnmetrics, _ := GetChurnMetricsWithWhitespace(repo, "README.md", "")
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
	_, err := GetChurnMetricsWithWhitespace(repo, "cmd/root.go", "")
	assert := assert.New(t)
	assert.Equal("The specified file was a new file added in this commit. Hence, churn can't be calculated.", err.Error())
}

func TestGetChurnMetricsSelfChurn(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/andymeneely/git-churn", "c0263662b2172b3df51ae39f8075dd010573ab6b")
	churnmetrics, _ := GetChurnMetricsWithWhitespace(repo, "matrics/diffmetrics_test.go", "")
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
	churnmetrics, _ := GetChurnMetricsWhitespaceExcluded(repo, "src/main/java/com/webcheckers/ui/WebServer.java", "")
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

func TestGetChurnMetricsWhitespaceExcludedAllRange(t *testing.T) {
	//https://github.com/ashishgalagali/SWEN610-project/compare/16c75b486a039bc34fcc5ac1ddad717d8bb49c01...7368d5fcb7eec950161ed9d13b55caf5961326b6?diff=split
	repo := gitfuncs.Checkout("https://github.com/ashishgalagali/SWEN610-project", "16c75b486a039bc34fcc5ac1ddad717d8bb49c01")
	churnmetrics, _ := GetChurnMetricsWhitespaceExcluded(repo, "README.md", "7368d5fcb7eec950161ed9d13b55caf5961326b6")
	assert := assert.New(t)
	assert.Equal("README.md", churnmetrics.FilePath)
	assert.Equal(13, churnmetrics.DeletedLinesCount)
	assert.Equal(6, churnmetrics.SelfChurnCount)
	assert.Equal(7, churnmetrics.InteractiveChurnCount)
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.CommitAuthor)
	assert.Equal(2, len(churnmetrics.ChurnDetails))
	assert.Equal("42880317+ks3057@users.noreply.github.com", churnmetrics.ChurnDetails["7368d5fcb7eec950161ed9d13b55caf5961326b6"])
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["8e6f09133b61c6eeb83d4e529c14c3754c286774"])
	assert.Equal(8, churnmetrics.FileDiffMetrics.Insertions)
	assert.Equal(13, churnmetrics.FileDiffMetrics.Deletions)
	assert.Equal(24, churnmetrics.FileDiffMetrics.LinesBefore)
	assert.Equal(19, churnmetrics.FileDiffMetrics.LinesAfter)

}
func TestGetChurnMetricsWhitespaceExcludedAllRangeRev(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/ashishgalagali/SWEN610-project", "7368d5fcb7eec950161ed9d13b55caf5961326b6")
	churnmetrics, _ := GetChurnMetricsWhitespaceExcluded(repo, "README.md", "16c75b486a039bc34fcc5ac1ddad717d8bb49c01")
	assert := assert.New(t)
	assert.Equal("README.md", churnmetrics.FilePath)
	assert.Equal(8, churnmetrics.DeletedLinesCount)
	assert.Equal(0, churnmetrics.SelfChurnCount)
	assert.Equal(8, churnmetrics.InteractiveChurnCount)
	assert.Equal("42880317+ks3057@users.noreply.github.com", churnmetrics.CommitAuthor)
	assert.Equal(2, len(churnmetrics.ChurnDetails))
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["7b56892de7fd86d1a3395a0bb10abef8ef3a033e"])
	assert.Equal("ashishgalagali@gmail.com", churnmetrics.ChurnDetails["979fe965043d49814c2fb7e7f5bae3461911b88b"])
	assert.Equal(13, churnmetrics.FileDiffMetrics.Insertions)
	assert.Equal(8, churnmetrics.FileDiffMetrics.Deletions)
	assert.Equal(19, churnmetrics.FileDiffMetrics.LinesBefore)
	assert.Equal(24, churnmetrics.FileDiffMetrics.LinesAfter)
}

func TestGetChurnMetricsWhitespaceExcludedInteractiveChurn(t *testing.T) {
	repo := gitfuncs.Checkout("https://github.com/andymeneely/git-churn", "99992110e402f26ca9162f43c0e5a97b1278068a")
	churnmetrics, _ := GetChurnMetricsWhitespaceExcluded(repo, "README.md", "")
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
