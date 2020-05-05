package metrics

import (
	"github.com/andymeneely/git-churn/gitfuncs"
	"github.com/andymeneely/git-churn/helper"
	. "github.com/andymeneely/git-churn/print"
	"gopkg.in/src-d/go-git.v4"
)

type ChurnMetrics struct {
	FilePath              string
	DeletedLinesCount     int
	SelfChurnCount        int
	InteractiveChurnCount int
	CommitAuthor          string
	//Map of CommitId, Author
	ChurnDetails    map[string]string
	FileDiffMetrics FileDiffMetrics
}

func GetChurnMetricsWithWhitespace(repo *git.Repository, filePath string) *ChurnMetrics {
	defer helper.Duration(helper.Track("GetChurnMetricsWithWhitespace"))
	fileDeletedLinesMap, _ := gitfuncs.DeletedLineNumbers(repo)
	churnMetrics := new(ChurnMetrics)
	CalculateChurnMetrics(fileDeletedLinesMap, repo, filePath, churnMetrics)
	churnMetrics.FileDiffMetrics = *CalculateDiffMetricsWithWhitespace(repo, filePath)
	return churnMetrics
}

func GetChurnMetricsWhitespaceExcluded(repo *git.Repository, filePath string) *ChurnMetrics {
	defer helper.Duration(helper.Track("GetChurnMetricsWhitespaceExcluded"))
	fileDeletedLinesMap, _ := gitfuncs.DeletedLineNumbersWhitespaceExcluded(repo)
	churnMetrics := new(ChurnMetrics)
	CalculateChurnMetrics(fileDeletedLinesMap, repo, filePath, churnMetrics)
	diffMetrics, _ := CalculateDiffMetricsWhitespaceExcluded(repo, filePath)
	churnMetrics.FileDiffMetrics = *diffMetrics
	return churnMetrics
}

func CalculateChurnMetrics(fileDeletedLinesMap map[string][]int, repo *git.Repository, filePath string, churnMetrics *ChurnMetrics) {
	deletedLines := fileDeletedLinesMap[filePath]
	parentCommitHash := gitfuncs.RevisionCommits(repo, "HEAD~1")

	head, _ := repo.Head()
	commitObj, err := repo.CommitObject(head.Hash())
	CheckIfError(err)
	commitAuthor := commitObj.Author.Email

	blame, err := gitfuncs.Blame(repo, parentCommitHash, filePath)
	CheckIfError(err)
	lines := blame.Lines

	churnDetails := make(map[string]string)
	selfChurnCount := 0
	interactiveChurnCount := 0
	for _, deletedLine := range deletedLines {
		churnAuthor := lines[deletedLine-1].Author
		if churnAuthor == commitAuthor {
			selfChurnCount += 1
		} else {
			interactiveChurnCount += 1
		}
		//fmt.Println(lines[deletedLine-1].Text)
		churnDetails[lines[deletedLine-1].Hash.String()] = churnAuthor
	}
	churnMetrics.DeletedLinesCount = len(deletedLines)
	churnMetrics.FilePath = filePath
	churnMetrics.SelfChurnCount = selfChurnCount
	churnMetrics.CommitAuthor = commitAuthor
	churnMetrics.InteractiveChurnCount = interactiveChurnCount
	churnMetrics.ChurnDetails = churnDetails
}
