package metrics

import (
	"errors"
	"fmt"
	"github.com/andymeneely/git-churn/gitfuncs"
	"strings"
)

type DiffMetrics struct {
	File        string
	Insertions  int
	Deletions   int
	LinesBefore int
	LinesAfter  int
	NewFile     bool
	DeleteFile  bool
}

func CalculateDiffMetricsWithWhitespace(repoUrl, commitHash, filePath string) *DiffMetrics {
	diffMetrics := new(DiffMetrics)
	diffMetrics.File = filePath
	changes, tree, parentTree := gitfuncs.CommitDiff(repoUrl, commitHash)
	patch, _ := changes.Patch()
	//fmt.Println(changes)
	//fmt.Println(patch)
	diffStats := patch.Stats()
	//fmt.Println(diffStats)

	//TODO: Throw error if file not exists in this commit
	for _, value := range diffStats {
		if value.Name == filePath {
			diffMetrics.Insertions = value.Addition
			diffMetrics.Deletions = value.Deletion
		}
	}

	diffMetrics.LinesBefore = gitfuncs.FileLOCFromTree(parentTree, filePath)
	diffMetrics.LinesAfter = gitfuncs.FileLOCFromTree(tree, filePath)

	if diffMetrics.LinesBefore == 0 && diffMetrics.LinesAfter != 0 {
		diffMetrics.NewFile = true
	}

	if diffMetrics.LinesBefore != 0 && diffMetrics.LinesAfter == 0 {
		diffMetrics.DeleteFile = true
	}

	return diffMetrics

}

func CalculateDiffMetricsWhitespaceExcluded(repoUrl, commitHash, filePath string) (*DiffMetrics, error) {
	diffMetrics := new(DiffMetrics)
	diffMetrics.File = filePath
	changes, tree, parentTree := gitfuncs.CommitDiff(repoUrl, commitHash)
	patch, _ := changes.Patch()

	fileDiffTexts := strings.Split(patch.String(), "diff --git a/"+filePath)
	if len(fileDiffTexts) < 2 {
		return nil, errors.New("File: " + filePath + " not found in the commitHash: " + commitHash)
	}
	fileDiff := strings.Split(fileDiffTexts[1], "+++")[1]
	fileDiff = strings.Split(fileDiff, "diff --git")[0]
	lines := strings.Split(fileDiff, "\n")

	insertions := 0
	deletions := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		fmt.Println(line)

		if strings.HasPrefix(line, "+") && line != "+" {
			insertions += 1
		}
		if strings.HasPrefix(line, "-") && line != "-" {
			deletions += 1
		}
	}

	diffMetrics.Insertions = insertions
	diffMetrics.Deletions = deletions

	diffMetrics.LinesBefore = gitfuncs.FileLOCFromTreeWhitespaceExcluded(parentTree, filePath)
	diffMetrics.LinesAfter = gitfuncs.FileLOCFromTreeWhitespaceExcluded(tree, filePath)

	if diffMetrics.LinesBefore == 0 && diffMetrics.LinesAfter != 0 {
		diffMetrics.NewFile = true
	}

	if diffMetrics.LinesBefore != 0 && diffMetrics.LinesAfter == 0 {
		diffMetrics.DeleteFile = true
	}

	return diffMetrics, nil

}
