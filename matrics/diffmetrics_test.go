package metrics

import (
	"github.com/andymeneely/git-churn/gitfuncs"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)


func TestFileOrigin(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "6255cfe24e726c0d9222075879e7a2676ac1b5a1", nil)
	diffmetrics := CalculateDiffMetricsWithWhitespace("testdata/file.txt", changes, tree, parentTree)
	assert := assert.New(t)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(4, diffmetrics.Insertions)
	assert.Equal(0, diffmetrics.Deletions)
	assert.Equal(0, diffmetrics.LinesBefore)
	assert.Equal(4, diffmetrics.LinesAfter)
	assert.Equal(true, diffmetrics.NewFile)
	assert.Equal(false, diffmetrics.DeleteFile)
}

func TestFileAddOnly(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "f33d22b9b10a084ef494df3c9780d30c41d3f54d", nil)
	diffmetrics := CalculateDiffMetricsWithWhitespace("testdata/file.txt", changes, tree, parentTree)
	//got := fmt.Sprintf("%v", diffmetrics)
	//expected := "&{testdata/file.txt 4 0 4 8 false false}"
	//assert.Equal(t, expected, got)
	assert := assert.New(t)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(4, diffmetrics.Insertions)
	assert.Equal(0, diffmetrics.Deletions)
	assert.Equal(4, diffmetrics.LinesBefore)
	assert.Equal(8, diffmetrics.LinesAfter)
	assert.Equal(false, diffmetrics.NewFile)
	assert.Equal(false, diffmetrics.DeleteFile)
}

func TestFileDeletesOnly(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "09e4b342693bf31bfb7cead1eb9b9fd59e3eef87", nil)
	diffmetrics := CalculateDiffMetricsWithWhitespace("testdata/file.txt", changes, tree, parentTree)
	assert := assert.New(t)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(0, diffmetrics.Insertions)
	assert.Equal(1, diffmetrics.Deletions)
	assert.Equal(8, diffmetrics.LinesBefore)
	assert.Equal(7, diffmetrics.LinesAfter)
	assert.Equal(false, diffmetrics.NewFile)
	assert.Equal(false, diffmetrics.DeleteFile)
}

func TestFileChangingLines(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "00da33207bbb17a149d99301012006fbd86c80e4", nil)
	diffmetrics := CalculateDiffMetricsWithWhitespace("testdata/file.txt", changes, tree, parentTree)
	assert := assert.New(t)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(1, diffmetrics.Insertions)
	assert.Equal(1, diffmetrics.Deletions)
	assert.Equal(7, diffmetrics.LinesBefore)
	assert.Equal(7, diffmetrics.LinesAfter)
	assert.Equal(false, diffmetrics.NewFile)
	assert.Equal(false, diffmetrics.DeleteFile)
}

func TestFileDelete(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "28b27020585be592df042c61ddab562665ce84cc", nil)
	diffmetrics := CalculateDiffMetricsWithWhitespace("testdata/file.txt", changes, tree, parentTree)
	assert := assert.New(t)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(0, diffmetrics.Insertions)
	assert.Equal(9, diffmetrics.Deletions)
	assert.Equal(9, diffmetrics.LinesBefore)
	assert.Equal(0, diffmetrics.LinesAfter)
	assert.Equal(false, diffmetrics.NewFile)
	assert.Equal(true, diffmetrics.DeleteFile)
}

func TestFileOriginWhitespaceExcluded(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "6255cfe24e726c0d9222075879e7a2676ac1b5a1", nil)
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("testdata/file.txt", changes, tree, parentTree)
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(3, diffmetrics.Insertions)
	assert.Equal(0, diffmetrics.Deletions)
	assert.Equal(0, diffmetrics.LinesBefore)
	assert.Equal(3, diffmetrics.LinesAfter)
	assert.Equal(true, diffmetrics.NewFile)
	assert.Equal(false, diffmetrics.DeleteFile)
}

func TestFileAddOnlyWhitespaceExcluded(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "f33d22b9b10a084ef494df3c9780d30c41d3f54d", nil)
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("testdata/file.txt", changes, tree, parentTree)
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(3, diffmetrics.Insertions)
	assert.Equal(0, diffmetrics.Deletions)
	assert.Equal(3, diffmetrics.LinesBefore)
	assert.Equal(6, diffmetrics.LinesAfter)
	assert.Equal(false, diffmetrics.NewFile)
	assert.Equal(false, diffmetrics.DeleteFile)
}

func TestFileDeletesOnlyWhitespaceExcluded(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "09e4b342693bf31bfb7cead1eb9b9fd59e3eef87", nil)
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("testdata/file.txt", changes, tree, parentTree)
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(0, diffmetrics.Insertions)
	assert.Equal(1, diffmetrics.Deletions)
	assert.Equal(6, diffmetrics.LinesBefore)
	assert.Equal(5, diffmetrics.LinesAfter)
	assert.Equal(false, diffmetrics.NewFile)
	assert.Equal(false, diffmetrics.DeleteFile)
}

func TestFileChangingLinesWhitespaceExcluded(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "00da33207bbb17a149d99301012006fbd86c80e4", nil)
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("testdata/file.txt", changes, tree, parentTree)
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(1, diffmetrics.Insertions)
	assert.Equal(1, diffmetrics.Deletions)
	assert.Equal(5, diffmetrics.LinesBefore)
	assert.Equal(5, diffmetrics.LinesAfter)
	assert.Equal(false, diffmetrics.NewFile)
	assert.Equal(false, diffmetrics.DeleteFile)
}

func TestFileDeleteWhitespaceExcluded(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "28b27020585be592df042c61ddab562665ce84cc", nil)
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("testdata/file.txt", changes, tree, parentTree)
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal("testdata/file.txt", diffmetrics.File)
	assert.Equal(0, diffmetrics.Insertions)
	assert.Equal(6, diffmetrics.Deletions)
	assert.Equal(6, diffmetrics.LinesBefore)
	assert.Equal(0, diffmetrics.LinesAfter)
	assert.Equal(false, diffmetrics.NewFile)
	assert.Equal(true, diffmetrics.DeleteFile)
}

func TestAggrDiffMetricsWithWhitespace(t *testing.T) {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "99992110e402f26ca9162f43c0e5a97b1278068a", nil)
	diffmetrics := AggrDiffMetricsWithWhitespace(changes, tree, parentTree)
	assert := assert.New(t)
	assert.Equal(17, diffmetrics.FilesCount)
	assert.Equal(225, diffmetrics.Insertions)
	assert.Equal(29, diffmetrics.Deletions)
	assert.Equal(2742, diffmetrics.LinesBefore)
	assert.Equal(2938, diffmetrics.LinesAfter)
	assert.Equal(1, diffmetrics.NewFiles)
	assert.Equal(0, diffmetrics.DeletedFiles)
}

func TestAggrDiffMetricsWhitespaceExcluded(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "99992110e402f26ca9162f43c0e5a97b1278068a", nil)
	diffmetrics, err := AggrDiffMetricsWhitespaceExcluded(changes, tree, parentTree)
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(17, diffmetrics.FilesCount)
	assert.Equal(217, diffmetrics.Insertions)
	assert.Equal(25, diffmetrics.Deletions)
	assert.Equal(2429, diffmetrics.LinesBefore)
	assert.Equal(2621, diffmetrics.LinesAfter)
	assert.Equal(1, diffmetrics.NewFiles)
	assert.Equal(0, diffmetrics.DeletedFiles)
}

func TestAggrFileDelete(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "28b27020585be592df042c61ddab562665ce84cc", nil)
	diffmetrics := AggrDiffMetricsWithWhitespace(changes, tree, parentTree)
	assert := assert.New(t)
	assert.Equal(5, diffmetrics.FilesCount)
	assert.Equal(0, diffmetrics.Insertions)
	assert.Equal(292, diffmetrics.Deletions)
	assert.Equal(2390, diffmetrics.LinesBefore)
	assert.Equal(2098, diffmetrics.LinesAfter)
	assert.Equal(0, diffmetrics.NewFiles)
	assert.Equal(5, diffmetrics.DeletedFiles)
	assert.Equal(diffmetrics.Deletions, diffmetrics.LinesBefore-diffmetrics.LinesAfter)
}

func TestAggrFileDeleteWhitespaceExcluded(t *testing.T) {
	changes, tree, parentTree := gitfuncs.CommitDiff(churnRepo, "28b27020585be592df042c61ddab562665ce84cc", nil)
	diffmetrics, err := AggrDiffMetricsWhitespaceExcluded(changes, tree, parentTree)
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(5, diffmetrics.FilesCount)
	assert.Equal(0, diffmetrics.Insertions)
	assert.Equal(258, diffmetrics.Deletions)
	assert.Equal(2139, diffmetrics.LinesBefore)
	assert.Equal(1881, diffmetrics.LinesAfter)
	assert.Equal(0, diffmetrics.NewFiles)
	assert.Equal(5, diffmetrics.DeletedFiles)
}
