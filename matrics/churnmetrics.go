package metrics

import (
	"encoding/json"
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
	"sort"
	"strconv"
	"strings"
	"sync"
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

type AggChurnMetrics struct {
	FilesCount                 int
	TotalDeletedLinesCount     int
	TotalSelfChurnCount        int
	TotalInteractiveChurnCount int
}

type AggCommitDetails struct {
	CommitId        string
	CommitAuthor    string
	DateTime        string
	CommitMessage   string
	AggChurnMetrics AggChurnMetrics
}

type AggChurnMetricsOutput struct {
	BaseCommitId     string
	AggCommitDetails []AggCommitDetails
}

type AggAllCommitDetails struct {
	TotalCommits    int
	AggChurnMetrics AggChurnMetrics
}

type AggAllChurnMetricsOutput struct {
	BaseCommitId     string
	ParentCommitId   string
	AggCommitDetails AggAllCommitDetails
}

func GetChurnMetrics(repo *git.Repository, baseCommitId, filePath, parentCommitId string, whitespace bool, jsonOPToFile, printOP bool) (*ChurnMetricsOutput, error) {
	defer helper.Duration(helper.Track("GetChurnMetrics"))
	var err error
	churnMetricsOutput := new(ChurnMetricsOutput)
	commits := make([]*object.Commit, 0)
	commitDetailsArr := make([]CommitDetails, 0)

	if baseCommitId == "" {
		//	https://mirrors.edge.kernel.org/pub/software/scm/git/docs/gitrevisions.html
		baseCommitHash := gitfuncs.RevisionCommits(repo, "", "origin/master@{1}")
		baseCommitId = baseCommitHash.String()
	}
	churnMetricsOutput.BaseCommitId = baseCommitId

	if printOP {
		fmt.Println("  ______  _              ______  _                          \n / _____)(_) _          / _____)| |                         \n| /  ___  _ | |_   ___ | /      | | _   _   _   ____  ____  \n| | (___)| ||  _) (___)| |      | || \\ | | | | / ___)|  _ \\ \n| \\____/|| || |__      | \\_____ | | | || |_| || |    | | | |\n \\_____/ |_| \\___)      \\______)|_| |_| \\____||_|    |_| |_|\n                                                            ")
		PrintInGreen("Base commitID: " + baseCommitId)
		PrintInBlue("")
	}

	commits, err = gitfuncs.RevList(repo, baseCommitId, parentCommitId)
	if len(commits) == 0 {
		commits, err = gitfuncs.RevList(repo, parentCommitId, baseCommitId)
	}
	commits = commits[1:]

	commitsChannel := make(chan CommitDetails)
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(commits))

	for _, commit := range commits {
		go getCommitDetails(commit, commitsChannel, repo, baseCommitId, filePath, whitespace)
	}
	go processCommmitDetails(commitsChannel, &commitDetailsArr, printOP, &waitGroup)

	time.Sleep(5000)
	waitGroup.Wait()
	close(commitsChannel)
	//  sorts by datetime
	sort.Slice(commitDetailsArr, func(i, j int) bool { return commitDetailsArr[i].DateTime > commitDetailsArr[j].DateTime })

	churnMetricsOutput.CommitDetails = commitDetailsArr
	//Check if churn-metrics folder exists else create. The output files will be stored in this folder
	if jsonOPToFile {
		writeJsonToFile(churnMetricsOutput)
	}

	return churnMetricsOutput, err
}

func writeJsonToFile(output interface{}) {
	var err error
	_, err = os.Stat("churn-metrics")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("churn-metrics", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	f, err := os.Create(filepath.Join("churn-metrics", "churn-metrics-op-"+time.Now().Format(time.RFC3339)+".json"))
	if err != nil {
		fmt.Println(err)
		f.Close()
	}
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	out, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	//output :=
	_, err = fmt.Fprintln(f, string(out))
	CheckIfError(err)
}

func processCommmitDetails(commitsChannel chan CommitDetails, commitDetailsArr *[]CommitDetails, printOP bool, wg *sync.WaitGroup) {
	time.Sleep(100)
	for commitDetails := range commitsChannel {
		//for {
		//	commitDetails, ok := <-commitsChannel
		//	if !ok {
		//		break
		//	}
		//fmt.Println("Processed: " + commitDetails.CommitId)
		//if !ok {
		//	break
		//}
		if len(commitDetails.ChurnMetrics) != 0 {
			if printOP {
				PrintInYellow("\tCommitID: " + commitDetails.CommitId)
				PrintInPink("\tCommit Author: " + commitDetails.CommitAuthor)
				PrintInGrey("\tDate: " + commitDetails.DateTime)
				PrintInBlue("\tMessage: " + strings.ReplaceAll(commitDetails.CommitMessage, "\n", " "))
				PrintInBlue("")
				for _, churnMetrics := range commitDetails.ChurnMetrics {
					if printOP {
						PrintInCyan("\t\tFile Path: " + churnMetrics.FilePath)
						fmt.Println("\t\tDeleted lines count: " + strconv.Itoa(churnMetrics.DeletedLinesCount))
						fmt.Println("\t\tSelf Churn count: " + strconv.Itoa(churnMetrics.SelfChurnCount))
						fmt.Println("\t\tInteractive Churn count: " + strconv.Itoa(churnMetrics.InteractiveChurnCount))
						fmt.Println("\t\tChurn Details :")
						for k, v := range churnMetrics.ChurnDetails {
							fmt.Println("\t\t\tcommit: " + k + ", author: " + v)
						}
						fmt.Println("")
					}
				}

				PrintInBlue("")
				PrintInBlue("")
				PrintInBlue("")
			}
			*commitDetailsArr = append(*commitDetailsArr, commitDetails)
		}
		wg.Done()
		//fmt.Println("DONE")
	}
}

func getCommitDetails(commit *object.Commit, commitsChannel chan CommitDetails, repo *git.Repository, baseCommitId string, filePath string, whitespace bool) {

	commitDetails := new(CommitDetails)
	parentCommitHash := gitfuncs.RevisionCommits(repo, baseCommitId, commit.Hash.String())
	commitDetails.CommitId = parentCommitHash.String()
	commitObj, err := repo.CommitObject(*parentCommitHash)
	CheckIfError(err)
	commitAuthor := commitObj.Author.Email
	commitDetails.CommitAuthor = commitAuthor
	commitDetails.DateTime = commitObj.Author.When.String()
	commitDetails.CommitMessage = commitObj.Message
	//if printOP {
	//	PrintInYellow("\tCommitID: " + commitDetails.CommitId)
	//	PrintInPink("\tCommit Author: " + commitDetails.CommitAuthor)
	//	PrintInGrey("\tDate: " + commitDetails.DateTime)
	//	PrintInBlue("\tMessage: " + strings.ReplaceAll(commitDetails.CommitMessage, "\n", " "))
	//	PrintInBlue("")
	//}
	churnMetricsArr, _ := calculateChurnMetrics(repo, baseCommitId, filePath, commitAuthor, parentCommitHash, whitespace)
	commitDetails.ChurnMetrics = churnMetricsArr
	//if len(churnMetricsArr) != 0 {
	//	commitDetailsArr = append(commitDetailsArr, *commitDetails)
	//}
	//if printOP {
	//	PrintInBlue("")
	//	PrintInBlue("")
	//	PrintInBlue("")
	//}
	//wg.Done()
	//fmt.Println("DONE")

	commitsChannel <- *commitDetails

}

//func GetChurnMetricsWhitespaceExcluded(repo *git.Repository, baseCommitId, filePath, parentCommitId string) (*CommitDetails, error) {
//	defer helper.Duration(helper.Track("GetChurnMetricsWhitespaceExcluded"))
//	parentCommitHash := gitfuncs.RevisionCommits(repo, baseCommitId, parentCommitId)
//	changes, tree, parentTree := gitfuncs.CommitDiff(repo, baseCommitId, parentCommitHash)
//	commitDetails := new(CommitDetails)
//	_, err := calculateChurnMetrics(repo, baseCommitId, filePath, "", parentCommitHash, nil)
//	diffMetrics, _ := CalculateDiffMetricsWhitespaceExcluded(filePath, changes, tree, parentTree)
//	fmt.Println(diffMetrics)
//	//commitDetails.FileDiffMetrics = *diffMetrics
//	return commitDetails, err
//}

func calculateChurnMetrics(repo *git.Repository, baseCommitId, filePath, commitAuthor string, parentCommitHash *plumbing.Hash, whitespace bool) ([]ChurnMetrics, error) {
	//REF: https://git-scm.com/docs/gitrevisions
	changes, _, _ := gitfuncs.CommitDiff(repo, baseCommitId, parentCommitHash)
	fileDeletedLinesMap := gitfuncs.DeletedLineNumbers(changes, filePath, whitespace)

	churnMetricsArr := make([]ChurnMetrics, 0)
	//head := plumbing.NewHash(baseCommitId)
	churnMetricsChannel := make(chan ChurnMetrics)

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(fileDeletedLinesMap))
	for filePath, deletedLines := range fileDeletedLinesMap {
		go getChurnMetrics(deletedLines, filePath, churnMetricsChannel, repo, parentCommitHash, commitAuthor, &waitGroup)
	}
	go procressChurnMetrics(churnMetricsChannel, &churnMetricsArr, &waitGroup)
	time.Sleep(5000)
	waitGroup.Wait()

	close(churnMetricsChannel)
	return churnMetricsArr, nil
}

func procressChurnMetrics(churnMetricsChannel chan ChurnMetrics, churnMetricsArr *[]ChurnMetrics, wg *sync.WaitGroup) {
	time.Sleep(500)
	for {
		churnMetrics, ok := <-churnMetricsChannel
		if !ok {
			break
		}
		*churnMetricsArr = append(*churnMetricsArr, churnMetrics)
		wg.Done()
	}
}

func getChurnMetrics(deletedLines []int, filePath string, churnMetricsChannel chan ChurnMetrics, repo *git.Repository, parentCommitHash *plumbing.Hash, commitAuthor string, wg *sync.WaitGroup) {
	//for filePath, deletedLines := range fileDeletedLinesMap {

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
			//if printOP {
			//	PrintInCyan("\t\tFile Path: " + churnMetrics.FilePath)
			//	fmt.Println("\t\tDeleted lines count: " + strconv.Itoa(churnMetrics.DeletedLinesCount))
			//	fmt.Println("\t\tSelf Churn count: " + strconv.Itoa(churnMetrics.SelfChurnCount))
			//	fmt.Println("\t\tInteractive Churn count: " + strconv.Itoa(churnMetrics.InteractiveChurnCount))
			//	fmt.Println("\t\tChurn Details :")
			//	for k, v := range churnDetails {
			//		fmt.Println("\t\t\tcommit: " + k + ", author: " + v)
			//	}
			//	fmt.Println("")
			//}

			//churnMetrics.FileDiffMetrics = *CalculateDiffMetricsWithWhitespace(filePath, changes, tree, parentTree)
			//churnMetricsArr = append(churnMetricsArr, *churnMetrics)

			churnMetricsChannel <- *churnMetrics
		} else {
			wg.Done()
			//fmt.Println(filePath + " : The specified file was a new file added in this commit. Hence, churn can't be calculated.")
		}
	} else {
		wg.Done()
	}
}

//}

func AggrChurnMetrics(repo *git.Repository, baseCommitId string, parentCommitId string, aggregate string, whitespace bool, jsonOPToFile bool, printOP bool) interface{} {
	defer helper.Duration(helper.Track("AggrChurnMetrics"))
	churnMetricsArr, _ := GetChurnMetrics(repo, baseCommitId, "", parentCommitId, whitespace, false, false)
	var aggChurnMetricsOutput interface{}
	if printOP {
		fmt.Println("  ______  _              ______  _                          \n / _____)(_) _          / _____)| |                         \n| /  ___  _ | |_   ___ | /      | | _   _   _   ____  ____  \n| | (___)| ||  _) (___)| |      | || \\ | | | | / ___)|  _ \\ \n| \\____/|| || |__      | \\_____ | | | || |_| || |    | | | |\n \\_____/ |_| \\___)      \\______)|_| |_| \\____||_|    |_| |_|\n                                                            ")
		PrintInGreen("Base commitID: " + baseCommitId)
		PrintInBlue("")
	}
	switch aggregate {
	case "commit":
		aggChurnMetricsOutput = getCommitAggChurnMetrics(churnMetricsArr, printOP)
	}
	if jsonOPToFile {
		writeJsonToFile(aggChurnMetricsOutput)
	}
	return aggChurnMetricsOutput
}

func getCommitAggChurnMetrics(churnMetricsOutput *ChurnMetricsOutput, printOP bool) interface{} {
	var aggChurnMetricsOP AggChurnMetricsOutput
	aggChurnMetricsOP.BaseCommitId = churnMetricsOutput.BaseCommitId
	commitDetailsArr := churnMetricsOutput.CommitDetails
	aggCommitDetailsArr := make([]AggCommitDetails, 0)

	for _, commitDetail := range commitDetailsArr {
		aggCommitDetails := new(AggCommitDetails)
		aggCommitDetails.CommitMessage = commitDetail.CommitMessage
		aggCommitDetails.DateTime = commitDetail.DateTime
		aggCommitDetails.CommitAuthor = commitDetail.CommitAuthor
		aggCommitDetails.CommitId = commitDetail.CommitId
		if printOP {
			PrintInYellow("\tCommitID: " + aggCommitDetails.CommitId)
			PrintInPink("\tCommit Author: " + aggCommitDetails.CommitAuthor)
			PrintInGrey("\tDate: " + aggCommitDetails.DateTime)
			PrintInBlue("\tMessage: " + strings.ReplaceAll(aggCommitDetails.CommitMessage, "\n", " "))
			PrintInBlue("")
		}
		var filesCount, totalDeletedLinesCount, totalSelfChurnCount, totalInteractiveChurnCount int
		for _, churnMetics := range commitDetail.ChurnMetrics {
			filesCount++
			totalDeletedLinesCount += churnMetics.DeletedLinesCount
			totalSelfChurnCount += churnMetics.SelfChurnCount
			totalInteractiveChurnCount += churnMetics.InteractiveChurnCount
		}
		var aggChurnMetrics AggChurnMetrics
		aggChurnMetrics.FilesCount = filesCount
		aggChurnMetrics.TotalDeletedLinesCount = totalDeletedLinesCount
		aggChurnMetrics.TotalSelfChurnCount = totalSelfChurnCount
		aggChurnMetrics.TotalInteractiveChurnCount = totalInteractiveChurnCount
		aggCommitDetails.AggChurnMetrics = aggChurnMetrics
		if printOP {
			PrintInCyan("\t\tTotal File Count: " + strconv.Itoa(aggChurnMetrics.FilesCount))
			fmt.Println("\t\tTotal Deleted lines count: " + strconv.Itoa(aggChurnMetrics.TotalDeletedLinesCount))
			fmt.Println("\t\tTotal Self Churn count: " + strconv.Itoa(aggChurnMetrics.TotalSelfChurnCount))
			fmt.Println("\t\tTotal Interactive Churn count: " + strconv.Itoa(aggChurnMetrics.TotalInteractiveChurnCount))
			PrintInBlue("")
			PrintInBlue("")
			PrintInBlue("")
		}
		aggCommitDetailsArr = append(aggCommitDetailsArr, *aggCommitDetails)
	}
	aggChurnMetricsOP.AggCommitDetails = aggCommitDetailsArr
	return aggChurnMetricsOP
}

//func AggrChurnMetricsWhitespaceExcluded(repo *git.Repository, baseCommitId string) *AggrChurMetrics {
//	defer helper.Duration(helper.Track("AggrChurnMetrics"))
//	changes, tree, parentTree := gitfuncs.CommitDiff(repo, baseCommitId, nil)
//	fileDeletedLinesMap := gitfuncs.DeletedLineNumbersWhitespaceExcluded(changes)
//	churnMetrics := new(AggrChurMetrics)
//	calculateAggrChurnMetrics(fileDeletedLinesMap, repo, baseCommitId, churnMetrics)
//	diffMetrics, _ := AggrDiffMetricsWhitespaceExcluded(changes, tree, parentTree)
//	churnMetrics.AggrDiffMetrics = *diffMetrics
//	return churnMetrics
//}

//func calculateAggrChurnMetrics(fileDeletedLinesMap map[string][]int, repo *git.Repository, baseCommitId string, churnMetrics *AggrChurMetrics) {
//	parentCommitHash := gitfuncs.RevisionCommits(repo, baseCommitId, "")
//	head := plumbing.NewHash(baseCommitId)
//	commitObj, err := repo.CommitObject(head)
//	CheckIfError(err)
//	commitAuthor := commitObj.Author.Email
//	totalDeletedLines := 0
//	totalSelfChurnCount := 0
//	totalInteractiveChurnCount := 0
//	for filePath, deletedLines := range fileDeletedLinesMap {
//		blame, err := gitfuncs.Blame(repo, parentCommitHash, filePath)
//		if err == nil && blame != nil {
//			lines := blame.Lines
//
//			for _, deletedLine := range deletedLines {
//				churnAuthor := lines[deletedLine-1].Author
//				if churnAuthor == commitAuthor {
//					totalSelfChurnCount += 1
//				} else {
//					totalInteractiveChurnCount += 1
//				}
//				//fmt.Println(lines[deletedLine-1].Text)
//			}
//			totalDeletedLines += len(deletedLines)
//		}
//	}
//	//churnMetrics.DeletedLinesCount = totalDeletedLines
//	//churnMetrics.SelfChurnCount = totalSelfChurnCount
//	//churnMetrics.CommitAuthor = commitAuthor
//	//churnMetrics.InteractiveChurnCount = totalInteractiveChurnCount
//}
