package matrics

import (
	"github.com/andymeneely/git-churn/gitfuncs"
)

type DiffMatrics struct {
	File        string
	Insertions  int
	Deletions   int
	LinesBefore int
	LinesAfter  int
	NewFile     bool
	DeleteFile  bool
}

func CalculateDiffMatrics(repoUrl, commitHash, filePath string) *DiffMatrics {
	diffMatrics := new(DiffMatrics)
	diffMatrics.File = filePath
	changes, tree, parentTree := gitfuncs.CommitDiff(repoUrl, commitHash)
	patch, _ := changes.Patch()
	//fmt.Println(changes)
	//fmt.Println(patch)
	diffStats := patch.Stats()
	//fmt.Println(diffStats)

	for _, value := range diffStats {
		if value.Name == filePath {
			diffMatrics.Insertions = value.Addition
			diffMatrics.Deletions = value.Deletion
		}
	}

	diffMatrics.LinesBefore = gitfuncs.FileLOCFromTree(parentTree, filePath)
	diffMatrics.LinesAfter = gitfuncs.FileLOCFromTree(tree, filePath)

	if diffMatrics.LinesBefore == 0 && diffMatrics.LinesAfter != 0 {
		diffMatrics.NewFile = true
	}

	if diffMatrics.LinesBefore != 0 && diffMatrics.LinesAfter == 0 {
		diffMatrics.DeleteFile = true
	}

	return diffMatrics

}
