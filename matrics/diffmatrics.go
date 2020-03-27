package metrics

import (
	"github.com/andymeneely/git-churn/gitfuncs"
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

func CalculateDiffMetrics(repoUrl, commitHash, filePath string) *DiffMetrics {
	diffMetrics := new(DiffMetrics)
	diffMetrics.File = filePath
	changes, tree, parentTree := gitfuncs.CommitDiff(repoUrl, commitHash)
	patch, _ := changes.Patch()
	//fmt.Println(changes)
	//fmt.Println(patch)
	diffStats := patch.Stats()
	//fmt.Println(diffStats)

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
