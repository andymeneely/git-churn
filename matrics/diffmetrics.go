package metrics

import (
	"errors"
	"github.com/andymeneely/git-churn/gitfuncs"
	"strings"
)

type DiffMetrics struct {
	Insertions  int
	Deletions   int
	LinesBefore int
	LinesAfter  int
}
type FileDiffMetrics struct {
	DiffMetrics
	File       string
	NewFile    bool
	DeleteFile bool
}
type AggrDiffMetrics struct {
	DiffMetrics
	FilesCount   int
	NewFiles     int
	DeletedFiles int
}

func CalculateDiffMetricsWithWhitespace(repoUrl, commitHash, filePath string) *FileDiffMetrics {
	diffMetrics := new(FileDiffMetrics)
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

func CalculateDiffMetricsWhitespaceExcluded(repoUrl, commitHash, filePath string) (*FileDiffMetrics, error) {
	diffMetrics := new(FileDiffMetrics)
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

func AggrDiffMetricsWithWhitespace(repoUrl, commitHash string) *AggrDiffMetrics {

	diffMetrics := new(AggrDiffMetrics)
	changes, tree, parentTree := gitfuncs.CommitDiff(repoUrl, commitHash)
	patch, _ := changes.Patch()
	//fmt.Println(changes)
	//fmt.Println(patch)
	diffStats := patch.Stats()
	//fmt.Println(diffStats)

	additions := 0
	deletions := 0
	for _, value := range diffStats {
		additions += value.Addition
		deletions += value.Deletion
	}
	diffMetrics.Insertions = additions
	diffMetrics.Deletions = deletions

	var beforeFiles []string
	var afterFiles []string
	diffMetrics.LinesBefore, beforeFiles = gitfuncs.LOCFilesFromTree(parentTree)
	diffMetrics.LinesAfter, afterFiles = gitfuncs.LOCFilesFromTree(tree)

	setFilesCounts(beforeFiles, afterFiles, diffMetrics)
	return diffMetrics
}

func setFilesCounts(beforeFiles []string, afterFiles []string, diffMetrics *AggrDiffMetrics) {
	// Putting the file names in a map to make lookup faster
	beforeSet := make(map[string]bool)
	for _, f := range beforeFiles {
		beforeSet[f] = true
	}

	// Putting the file names in a map to make lookup faster
	afterSet := make(map[string]bool)
	for _, f := range afterFiles {
		afterSet[f] = true
	}

	deletedFiles := 0
	newFiles := 0

	for _, file := range beforeFiles {
		if !afterSet[file] {
			deletedFiles += 1
		}
	}

	for _, file := range afterFiles {
		if !beforeSet[file] {
			newFiles += 1
		}
	}

	diffMetrics.NewFiles = newFiles
	diffMetrics.DeletedFiles = deletedFiles
	diffMetrics.FilesCount = len(afterFiles)

}

func AggrDiffMetricsWhitespaceExcluded(repoUrl, commitHash string) (*AggrDiffMetrics, error) {
	diffMetrics := new(AggrDiffMetrics)
	changes, tree, parentTree := gitfuncs.CommitDiff(repoUrl, commitHash)
	patch, _ := changes.Patch()

	fileDiffTexts := strings.Split(patch.String(), "diff --git a/")
	insertions := 0
	deletions := 0
	for index, _ := range fileDiffTexts {
		if index == 0 {
			continue
		}
		fileDiff := strings.Split(fileDiffTexts[index], "+++")[1]
		fileDiff = strings.Split(fileDiff, "diff --git")[0]
		lines := strings.Split(fileDiff, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "+") && line != "+" {
				insertions += 1
			}
			if strings.HasPrefix(line, "-") && line != "-" {
				deletions += 1
			}
		}
	}

	diffMetrics.Insertions = insertions
	diffMetrics.Deletions = deletions

	var beforeFiles []string
	var afterFiles []string
	diffMetrics.LinesBefore, beforeFiles = gitfuncs.LOCFilesFromTreeWhitespaceExcluded(parentTree)
	diffMetrics.LinesAfter, afterFiles = gitfuncs.LOCFilesFromTreeWhitespaceExcluded(tree)

	setFilesCounts(beforeFiles, afterFiles, diffMetrics)
	return diffMetrics, nil
}
