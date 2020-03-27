package metrics

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFileOrigin(t *testing.T) {
	diffmetrics := CalculateDiffMetrics("https://github.com/andymeneely/git-churn", "6255cfe24e726c0d9222075879e7a2676ac1b5a1", "testdata/file.txt")
	got, err := json.Marshal(diffmetrics)
	if err != nil {
		panic(err)
	}
	expected := "{\"File\":\"testdata/file.txt\",\"Insertions\":4,\"Deletions\":0,\"LinesBefore\":0,\"LinesAfter\":4,\"NewFile\":true,\"DeleteFile\":false}"
	if string(got) != expected {
		t.Errorf("DiffMetrics expected %q but got %q", expected, got)
	}
}

func TestFileAddOnly(t *testing.T) {
	diffmetrics := CalculateDiffMetrics("https://github.com/andymeneely/git-churn", "f33d22b9b10a084ef494df3c9780d30c41d3f54d", "testdata/file.txt")
	got := fmt.Sprintf("%v", diffmetrics)
	expected := "&{testdata/file.txt 4 0 4 8 false false}"
	if got != expected {
		t.Errorf("DiffMetrics expected %q but got %q", expected, got)
	}
}

func TestFileDeletesOnly(t *testing.T) {
	diffmetrics := CalculateDiffMetrics("https://github.com/andymeneely/git-churn", "09e4b342693bf31bfb7cead1eb9b9fd59e3eef87", "testdata/file.txt")
	got := fmt.Sprintf("%v", diffmetrics)
	expected := "&{testdata/file.txt 0 1 8 7 false false}"
	if got != expected {
		t.Errorf("DiffMetrics expected %q but got %q", expected, got)
	}
}

func TestFileChangingLines(t *testing.T) {
	diffmetrics := CalculateDiffMetrics("https://github.com/andymeneely/git-churn", "00da33207bbb17a149d99301012006fbd86c80e4", "testdata/file.txt")
	got := fmt.Sprintf("%v", diffmetrics)
	expected := "&{testdata/file.txt 1 1 7 7 false false}"
	if got != expected {
		t.Errorf("DiffMetrics expected %q but got %q", expected, got)
	}
}

func TestFileDelete(t *testing.T) {
	diffmetrics := CalculateDiffMetrics("https://github.com/andymeneely/git-churn", "28b27020585be592df042c61ddab562665ce84cc", "testdata/file.txt")
	got, err := json.Marshal(diffmetrics)
	if err != nil {
		panic(err)
	}
	expected := "{\"File\":\"testdata/file.txt\",\"Insertions\":0,\"Deletions\":9,\"LinesBefore\":9,\"LinesAfter\":0,\"NewFile\":false,\"DeleteFile\":true}"
	if string(got) != expected {
		t.Errorf("DiffMetrics expected %q but got %q", expected, got)
	}
}
