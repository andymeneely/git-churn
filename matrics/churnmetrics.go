package metrics

import (
	"errors"
	"github.com/andymeneely/git-churn/gitfuncs"
	"github.com/andymeneely/git-churn/helper"
	. "github.com/andymeneely/git-churn/print"
	"gopkg.in/src-d/go-git.v4"
)

type ChurnMetrics struct {
	DeletedLinesCount     int
	SelfChurnCount        int
	InteractiveChurnCount int
	CommitAuthor          string
}

type AggrChurMetrics struct {
	ChurnMetrics
	AggrDiffMetrics AggrDiffMetrics
}

type FileChurnMetrics struct {
	ChurnMetrics
	FilePath string
	//Map of CommitId, Author
	ChurnDetails    map[string]string
	FileDiffMetrics FileDiffMetrics
}

func GetChurnMetricsWithWhitespace(repo *git.Repository, filePath string) (*FileChurnMetrics, error) {
	defer helper.Duration(helper.Track("GetChurnMetricsWithWhitespace"))
	fileDeletedLinesMap, _ := gitfuncs.DeletedLineNumbers(repo)
	churnMetrics := new(FileChurnMetrics)
	err := calculateChurnMetrics(fileDeletedLinesMap, repo, filePath, churnMetrics)
	churnMetrics.FileDiffMetrics = *CalculateDiffMetricsWithWhitespace(repo, filePath)
	return churnMetrics, err
}

func GetChurnMetricsWhitespaceExcluded(repo *git.Repository, filePath string) (*FileChurnMetrics, error) {
	defer helper.Duration(helper.Track("GetChurnMetricsWhitespaceExcluded"))
	fileDeletedLinesMap, _ := gitfuncs.DeletedLineNumbersWhitespaceExcluded(repo)
	churnMetrics := new(FileChurnMetrics)
	err := calculateChurnMetrics(fileDeletedLinesMap, repo, filePath, churnMetrics)
	diffMetrics, _ := CalculateDiffMetricsWhitespaceExcluded(repo, filePath)
	churnMetrics.FileDiffMetrics = *diffMetrics
	return churnMetrics, err
}

func calculateChurnMetrics(fileDeletedLinesMap map[string][]int, repo *git.Repository, filePath string, churnMetrics *FileChurnMetrics) error {
	deletedLines := fileDeletedLinesMap[filePath]
	//REF: https://git-scm.com/docs/gitrevisions
	parentCommitHash := gitfuncs.RevisionCommits(repo, "HEAD~1")

	blame, err := gitfuncs.Blame(repo, parentCommitHash, filePath)
	if err != nil {
		return errors.New("The specified file was a new file added in this commit. Hence, churn can't be calculated.")
	}
	lines := blame.Lines

	head, _ := repo.Head()
	commitObj, err := repo.CommitObject(head.Hash())
	CheckIfError(err)
	commitAuthor := commitObj.Author.Email

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
	return nil
}

func AggrChurnMetricsWithWhitespace(repo *git.Repository) *AggrChurMetrics {
	defer helper.Duration(helper.Track("AggrChurnMetricsWithWhitespace"))
	fileDeletedLinesMap, _ := gitfuncs.DeletedLineNumbers(repo)
	churnMetrics := new(AggrChurMetrics)
	calculateAggrChurnMetrics(fileDeletedLinesMap, repo, churnMetrics)
	diffMetrics := AggrDiffMetricsWithWhitespace(repo)
	churnMetrics.AggrDiffMetrics = *diffMetrics
	return churnMetrics
}

func AggrChurnMetricsWhitespaceExcluded(repo *git.Repository) *AggrChurMetrics {
	defer helper.Duration(helper.Track("AggrChurnMetricsWithWhitespace"))
	fileDeletedLinesMap, _ := gitfuncs.DeletedLineNumbersWhitespaceExcluded(repo)
	churnMetrics := new(AggrChurMetrics)
	calculateAggrChurnMetrics(fileDeletedLinesMap, repo, churnMetrics)
	diffMetrics, _ := AggrDiffMetricsWhitespaceExcluded(repo)
	churnMetrics.AggrDiffMetrics = *diffMetrics
	return churnMetrics
}

func calculateAggrChurnMetrics(fileDeletedLinesMap map[string][]int, repo *git.Repository, churnMetrics *AggrChurMetrics) {
	parentCommitHash := gitfuncs.RevisionCommits(repo, "HEAD~1")
	head, _ := repo.Head()
	commitObj, err := repo.CommitObject(head.Hash())
	CheckIfError(err)
	commitAuthor := commitObj.Author.Email
	totalDeletedLines := 0
	totalSelfChurnCount := 0
	totalInteractiveChurnCount := 0
	for filePath, deletedLines := range fileDeletedLinesMap {
		blame, err := gitfuncs.Blame(repo, parentCommitHash, filePath)
		if err == nil && blame != nil {
			lines := blame.Lines

			for _, deletedLine := range deletedLines {
				churnAuthor := lines[deletedLine-1].Author
				if churnAuthor == commitAuthor {
					totalSelfChurnCount += 1
				} else {
					totalInteractiveChurnCount += 1
				}
				//fmt.Println(lines[deletedLine-1].Text)
			}
			totalDeletedLines += len(deletedLines)
		}
	}
	churnMetrics.DeletedLinesCount = totalDeletedLines
	churnMetrics.SelfChurnCount = totalSelfChurnCount
	churnMetrics.CommitAuthor = commitAuthor
	churnMetrics.InteractiveChurnCount = totalInteractiveChurnCount
}
