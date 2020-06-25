package metrics

import (
	"fmt"
	"github.com/andymeneely/git-churn/gitfuncs"
	"github.com/andymeneely/git-churn/helper"
	. "github.com/andymeneely/git-churn/print"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ChurnMetrics struct {
	FilePath              string
	DeletedLinesCount     int
	SelfChurnCount        int
	InteractiveChurnCount int
	//Map of CommitId, Author
	ChurnDetails map[string]string
	//FileDiffMetrics FileDiffMetrics
}

type CommitDetails struct {
	CommitId      string
	CommitAuthor  string
	DateTime      string
	CommitMessage string
	ChurnMetrics  []ChurnMetrics
}

type ChurnMetricsOutput struct {
	BaseCommitId  string
	CommitDetails []CommitDetails
}

type AggrChurMetrics struct {
	ChurnMetricsOutput
	AggrDiffMetrics AggrDiffMetrics
}

type FileChurnMetrics struct {
	ChurnMetricsOutput
	FilePath       string
	BaseCommitId   string
	ParentCommitId string
	//Map of CommitId, Author
	ChurnDetails    map[string]string
	FileDiffMetrics FileDiffMetrics
}

func GetChurnMetricsWithWhitespace(repo *git.Repository, baseCommitId, filePath, parentCommitId string) (*ChurnMetricsOutput, error) {
	defer helper.Duration(helper.Track("GetChurnMetricsWithWhitespace"))
	var err error
	_, err = os.Stat("churn-metrics")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("churn-metrics", 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
	f, err := os.Create(filepath.Join("churn-metrics", "churn-metrics-op-"+time.Now().Format(time.RFC3339)+".txt"))
	if err != nil {
		fmt.Println(err)
		f.Close()
	}
	churnMetricsOutput := new(ChurnMetricsOutput)
	commits := make([]*object.Commit, 0)
	commitDetailsArr := make([]CommitDetails, 0)

	if baseCommitId == "" {
		baseCommitHash := gitfuncs.RevisionCommits(repo, "", "origin/master@{1}")
		baseCommitId = baseCommitHash.String()
		fmt.Println(baseCommitId)
	}
	churnMetricsOutput.BaseCommitId = baseCommitId
	fmt.Fprintln(f, "Base commitID: "+baseCommitId)
	fmt.Fprintln(f, "")

	//if parentCommitId != "" {
	commits, err = gitfuncs.RevList(repo, baseCommitId, parentCommitId)
	if len(commits) == 0 {
		commits, err = gitfuncs.RevList(repo, parentCommitId, baseCommitId)
	}
	commits = commits[1:]
	for _, commit := range commits {
		commitDetails := new(CommitDetails)
		parentCommitHash := gitfuncs.RevisionCommits(repo, baseCommitId, commit.Hash.String())
		commitDetails.CommitId = parentCommitHash.String()
		commitObj, err := repo.CommitObject(*parentCommitHash)
		CheckIfError(err)
		commitAuthor := commitObj.Author.Email
		commitDetails.CommitAuthor = commitAuthor
		commitDetails.DateTime = commitObj.Author.When.String()
		commitDetails.CommitMessage = commitObj.Message
		fmt.Fprintln(f, "\tCommitID: "+commitDetails.CommitId)
		fmt.Fprintln(f, "\tCommit Author: "+commitDetails.CommitAuthor)
		fmt.Fprintln(f, "\tDate: "+commitDetails.DateTime)
		fmt.Fprintln(f, "\tMessage: "+strings.ReplaceAll(commitDetails.CommitMessage, "\n", " "))
		fmt.Fprintln(f, "")

		churnMetricsArr, _ := calculateChurnMetrics(repo, baseCommitId, filePath, commitAuthor, parentCommitHash, f)
		commitDetails.ChurnMetrics = churnMetricsArr
		if len(churnMetricsArr) != 0 {
			commitDetailsArr = append(commitDetailsArr, *commitDetails)
		}
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "")
	}
	//} else {
	//	commitDetails := new(CommitDetails)
	//	parentCommitHash := gitfuncs.RevisionCommits(repo, baseCommitId, "")
	//	commitDetails.CommitId = parentCommitHash.String()
	//	churnMetricsArr, _ := calculateChurnMetrics(repo, baseCommitId, filePath, commitDetails, parentCommitHash)
	//	commitDetails.ChurnMetrics = churnMetricsArr
	//	commitDetailsArr = append(commitDetailsArr, *commitDetails)
	//}
	churnMetricsOutput.CommitDetails = commitDetailsArr
	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}
	return churnMetricsOutput, err
}

func GetChurnMetricsWhitespaceExcluded(repo *git.Repository, baseCommitId, filePath, parentCommitId string) (*CommitDetails, error) {
	defer helper.Duration(helper.Track("GetChurnMetricsWhitespaceExcluded"))
	parentCommitHash := gitfuncs.RevisionCommits(repo, baseCommitId, parentCommitId)
	changes, tree, parentTree := gitfuncs.CommitDiff(repo, baseCommitId, parentCommitHash)
	commitDetails := new(CommitDetails)
	_, err := calculateChurnMetrics(repo, baseCommitId, filePath, "", parentCommitHash, nil)
	diffMetrics, _ := CalculateDiffMetricsWhitespaceExcluded(filePath, changes, tree, parentTree)
	fmt.Println(diffMetrics)
	//commitDetails.FileDiffMetrics = *diffMetrics
	return commitDetails, err
}

func calculateChurnMetrics(repo *git.Repository, baseCommitId string, filePath string, commitAuthor string, parentCommitHash *plumbing.Hash, f *os.File) ([]ChurnMetrics, error) {
	//REF: https://git-scm.com/docs/gitrevisions
	changes, _, _ := gitfuncs.CommitDiff(repo, baseCommitId, parentCommitHash)
	fileDeletedLinesMap := gitfuncs.DeletedLineNumbers(changes, filePath)
	//fileDeletedLinesMap := gitfuncs.DeletedLineNumbersWhitespaceExcluded(changes)

	churnMetricsArr := make([]ChurnMetrics, 0)
	//head := plumbing.NewHash(baseCommitId)
	for filePath, deletedLines := range fileDeletedLinesMap {
		if len(deletedLines) != 0 {
			//deletedLines := fileDeletedLinesMap[filePath]
			churnMetrics := new(ChurnMetrics)
			blame, err := gitfuncs.Blame(repo, parentCommitHash, filePath)
			if err == nil {

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
				churnMetrics.SelfChurnCount = selfChurnCount
				churnMetrics.InteractiveChurnCount = interactiveChurnCount
				churnMetrics.ChurnDetails = churnDetails
				churnMetrics.FilePath = filePath
				fmt.Fprintln(f, "\t\tFile Path: "+churnMetrics.FilePath)
				fmt.Fprintln(f, "\t\tDeleted lines count: "+strconv.Itoa(churnMetrics.DeletedLinesCount))
				fmt.Fprintln(f, "\t\tSelf Churn count: "+strconv.Itoa(churnMetrics.SelfChurnCount))
				fmt.Fprintln(f, "\t\tInteractive Churn count: "+strconv.Itoa(churnMetrics.InteractiveChurnCount))
				fmt.Fprintln(f, "\t\tChurn Details :")
				for k, v := range churnDetails {
					fmt.Fprintln(f, "\t\t\tcommit: "+k+", author: "+v)
				}
				fmt.Fprintln(f, "")

				//churnMetrics.FileDiffMetrics = *CalculateDiffMetricsWithWhitespace(filePath, changes, tree, parentTree)
				churnMetricsArr = append(churnMetricsArr, *churnMetrics)

			} else {
				//fmt.Println(filePath + " : The specified file was a new file added in this commit. Hence, churn can't be calculated.")
			}
		}
	}
	return churnMetricsArr, nil
}

func AggrChurnMetricsWithWhitespace(repo *git.Repository, baseCommitId string) *AggrChurMetrics {
	defer helper.Duration(helper.Track("AggrChurnMetricsWithWhitespace"))
	changes, tree, parentTree := gitfuncs.CommitDiff(repo, baseCommitId, nil)
	fileDeletedLinesMap := gitfuncs.DeletedLineNumbers(changes, "")
	churnMetrics := new(AggrChurMetrics)
	calculateAggrChurnMetrics(fileDeletedLinesMap, repo, baseCommitId, churnMetrics)
	diffMetrics := AggrDiffMetricsWithWhitespace(changes, tree, parentTree)
	churnMetrics.AggrDiffMetrics = *diffMetrics
	return churnMetrics
}

func AggrChurnMetricsWhitespaceExcluded(repo *git.Repository, baseCommitId string) *AggrChurMetrics {
	defer helper.Duration(helper.Track("AggrChurnMetricsWithWhitespace"))
	changes, tree, parentTree := gitfuncs.CommitDiff(repo, baseCommitId, nil)
	fileDeletedLinesMap := gitfuncs.DeletedLineNumbersWhitespaceExcluded(changes)
	churnMetrics := new(AggrChurMetrics)
	calculateAggrChurnMetrics(fileDeletedLinesMap, repo, baseCommitId, churnMetrics)
	diffMetrics, _ := AggrDiffMetricsWhitespaceExcluded(changes, tree, parentTree)
	churnMetrics.AggrDiffMetrics = *diffMetrics
	return churnMetrics
}

func calculateAggrChurnMetrics(fileDeletedLinesMap map[string][]int, repo *git.Repository, baseCommitId string, churnMetrics *AggrChurMetrics) {
	parentCommitHash := gitfuncs.RevisionCommits(repo, baseCommitId, "")
	head := plumbing.NewHash(baseCommitId)
	commitObj, err := repo.CommitObject(head)
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
	//churnMetrics.DeletedLinesCount = totalDeletedLines
	//churnMetrics.SelfChurnCount = totalSelfChurnCount
	//churnMetrics.CommitAuthor = commitAuthor
	//churnMetrics.InteractiveChurnCount = totalInteractiveChurnCount
}
