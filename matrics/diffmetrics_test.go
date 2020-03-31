package metrics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileOrigin(t *testing.T) {
	diffmetrics := CalculateDiffMetricsWithWhitespace("https://github.com/andymeneely/git-churn", "6255cfe24e726c0d9222075879e7a2676ac1b5a1", "testdata/file.txt")
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
	diffmetrics := CalculateDiffMetricsWithWhitespace("https://github.com/andymeneely/git-churn", "f33d22b9b10a084ef494df3c9780d30c41d3f54d", "testdata/file.txt")
	got := fmt.Sprintf("%v", diffmetrics)
	expected := "&{testdata/file.txt 4 0 4 8 false false}"
	assert.Equal(t, expected, got)
}

func TestFileDeletesOnly(t *testing.T) {
	diffmetrics := CalculateDiffMetricsWithWhitespace("https://github.com/andymeneely/git-churn", "09e4b342693bf31bfb7cead1eb9b9fd59e3eef87", "testdata/file.txt")
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
	diffmetrics := CalculateDiffMetricsWithWhitespace("https://github.com/andymeneely/git-churn", "00da33207bbb17a149d99301012006fbd86c80e4", "testdata/file.txt")
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
	diffmetrics := CalculateDiffMetricsWithWhitespace("https://github.com/andymeneely/git-churn", "28b27020585be592df042c61ddab562665ce84cc", "testdata/file.txt")
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
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("https://github.com/andymeneely/git-churn", "6255cfe24e726c0d9222075879e7a2676ac1b5a1", "testdata/file.txt")
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
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("https://github.com/andymeneely/git-churn", "f33d22b9b10a084ef494df3c9780d30c41d3f54d", "testdata/file.txt")
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
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("https://github.com/andymeneely/git-churn", "09e4b342693bf31bfb7cead1eb9b9fd59e3eef87", "testdata/file.txt")
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
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("https://github.com/andymeneely/git-churn", "00da33207bbb17a149d99301012006fbd86c80e4", "testdata/file.txt")
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
	diffmetrics, err := CalculateDiffMetricsWhitespaceExcluded("https://github.com/andymeneely/git-churn", "28b27020585be592df042c61ddab562665ce84cc", "testdata/file.txt")
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
