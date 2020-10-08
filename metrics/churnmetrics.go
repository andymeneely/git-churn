package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/andymeneely/git-churn/gitfuncs"
	"github.com/andymeneely/git-churn/helper"
	. "github.com/andymeneely/git-churn/print"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

type AggAllChurnMetricsOutput struct {
	BaseCommitId               string
	ParentCommitId             string
	TotalCommits               int
	TotalDeletedLinesCount     int
	TotalSelfChurnCount        int
	TotalInteractiveChurnCount int
}

// GetChurnMetrics Returns the Churn metrics details for the given repo
// 		baseCommitId and parentCommitId are the two commits between which the churn details are requested.
// If only one is present then root commit will be taken as default for other. If none is present then all the commits of the project are considered.
// 		whitespace: if false neglects the deleted whitespaces.
// 		jsonOPToFile: if true writes the JSON output to the a file
// 		printOP: if true prints the output into the console in a human readable form
func GetChurnMetrics(repo *git.Repository, baseCommitId, filePath, parentCommitId string, whitespace bool, jsonOPToFile, printOP bool) (*ChurnMetricsOutput, error) {
	//helper.INFO.Println("INSIDE : GetChurnMetrics")

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
		// prints GIT-CHURN banner
		fmt.Println("  ______  _              ______  _                          \n / _____)(_) _          / _____)| |                         \n| /  ___  _ | |_   ___ | /      | | _   _   _   ____  ____  \n| | (___)| ||  _) (___)| |      | || \\ | | | | / ___)|  _ \\ \n| \\____/|| || |__      | \\_____ | | | || |_| || |    | | | |\n \\_____/ |_| \\___)      \\______)|_| |_| \\____||_|    |_| |_|\n                                                            ")
		PrintInGreen("Base commit ID: " + baseCommitId)
		PrintInBlue("")
	}

	commits, err = gitfuncs.RevList(repo, baseCommitId, parentCommitId)
	if len(commits) == 0 {
		commits, err = gitfuncs.RevList(repo, parentCommitId, baseCommitId)
	}
	// neglect the 1st commit as we need not compare the commit with itself
	commits = commits[1:]

	//Ref. https://golangbot.com/buffered-channels-worker-pools/
	// Channel to hold the commit details
	commitsChannel := make(chan *object.Commit, 10)
	commitDetailsChannel := make(chan CommitDetails, 10)

	helper.INFO.Println("Commits count: " + strconv.Itoa(len(commits)))
	go allocate(commitsChannel, commits)
	done := make(chan bool)

	go processCommitDetails(commitDetailsChannel, &commitDetailsArr, printOP, done)
	noOfWorkers := 10
	//createWorkerPool(noOfWorkers)

	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go getCommitDetails(commitsChannel, commitDetailsChannel, repo, baseCommitId, filePath, whitespace, &wg)
	}
	wg.Wait()
	helper.INFO.Println("Closed commitDetailsChannel CHANNEL!!!!!")
	close(commitDetailsChannel)

	helper.INFO.Println("PROCESSED ALL " + strconv.Itoa(len(commitDetailsArr)))
	<-done
	//  sorts by datetime
	sort.Slice(commitDetailsArr, func(i, j int) bool { return commitDetailsArr[i].DateTime > commitDetailsArr[j].DateTime })

	churnMetricsOutput.CommitDetails = commitDetailsArr
	//Check if churn-metrics folder exists else create. The output files will be stored in this folder
	if jsonOPToFile {
		writeJsonToFile(churnMetricsOutput)
	}

	return churnMetricsOutput, err
}

func allocate(commitsChannel chan *object.Commit, commits []*object.Commit) {
	for _, commit := range commits {
		commitsChannel <- commit
	}
	close(commitsChannel)
}

// writeJsonToFile write the Json output into a file named churn-metrics-op-<time> into churn-metrics dir
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
	_, err = fmt.Fprintln(f, string(out))
	CheckIfError(err)
}

// processCommitDetails appends each CommitDetails from the commitsChannel into commitDetailsArr. Also prints each commit detail if printOP is true.
func processCommitDetails(commitDetailsChannel chan CommitDetails, commitDetailsArr *[]CommitDetails, printOP bool, done chan bool) {
	//helper.INFO.Println("INSIDE : processCommitDetails")
	for commitDetails := range commitDetailsChannel {
		if len(commitDetails.ChurnMetrics) != 0 {
			helper.INFO.Println("Processed churn-metrics for commit: " + commitDetails.CommitId)
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
	}
	done <- true
}

// getCommitDetails adds commitDetails for the given commit into the commitsChannel
func getCommitDetails(commitsChannel chan *object.Commit, commitsDetailsChannel chan CommitDetails, repo *git.Repository, baseCommitId string, filePath string, whitespace bool, wg *sync.WaitGroup) {
	//helper.INFO.Println("INSIDE : getCommitDetails")
	for commit := range commitsChannel {
		time.Sleep(2000)
		commitDetails := new(CommitDetails)
		parentCommitHash := gitfuncs.RevisionCommits(repo, baseCommitId, commit.Hash.String())
		commitDetails.CommitId = parentCommitHash.String()
		commitObj, err := repo.CommitObject(*parentCommitHash)
		CheckIfError(err)
		commitAuthor := commitObj.Author.Email
		commitDetails.CommitAuthor = commitAuthor
		commitDetails.DateTime = commitObj.Author.When.String()
		commitDetails.CommitMessage = commitObj.Message
		churnMetricsArr, _ := calculateChurnMetrics(repo, baseCommitId, filePath, commitAuthor, parentCommitHash, whitespace)
		commitDetails.ChurnMetrics = churnMetricsArr
		commitsDetailsChannel <- *commitDetails
	}
	wg.Done()
}

type FileDeletedLines struct {
	file         string
	deletedLines []int
}

// calculateChurnMetrics calculate the churn metrics and returns the array of ChurnMetrics
func calculateChurnMetrics(repo *git.Repository, baseCommitId, filePath, commitAuthor string, parentCommitHash *plumbing.Hash, whitespace bool) ([]ChurnMetrics, error) {
	//REF: https://git-scm.com/docs/gitrevisions
	//helper.INFO.Println("INSIDE : calculateChurnMetrics")
	changes, _, _ := gitfuncs.CommitDiff(repo, baseCommitId, parentCommitHash)
	fileDeletedLinesMap := gitfuncs.DeletedLineNumbers(changes, filePath, whitespace)

	churnMetricsArr := make([]ChurnMetrics, 0)
	churnMetricsChannel := make(chan ChurnMetrics, 2)
	ipFileChannel := make(chan FileDeletedLines, 2)

	go allocateFiles(ipFileChannel, fileDeletedLinesMap)
	done := make(chan bool)
	go processChurnMetrics(churnMetricsChannel, &churnMetricsArr, done)

	noOfWorkers := 2
	//createWorkerPool(noOfWorkers)

	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go getChurnMetrics(ipFileChannel, churnMetricsChannel, repo, parentCommitHash, commitAuthor, &wg)
	}
	wg.Wait()
	//helper.INFO.Println("Closed commitDetailsChannel CHANNEL!!!!!")
	close(churnMetricsChannel)

	//helper.INFO.Println("PROCESSED ALL " + strconv.Itoa(len(churnMetricsChannel)))
	<-done

	//var waitGroup sync.WaitGroup
	//waitGroup.Add(len(fileDeletedLinesMap))
	//for filePath, deletedLines := range fileDeletedLinesMap {
	//	//go getChurnMetrics(deletedLines, filePath, churnMetricsChannel, repo, parentCommitHash, commitAuthor, &waitGroup)
	//}
	////go processChurnMetrics(churnMetricsChannel, &churnMetricsArr, &waitGroup)
	//// TODO: Research a better method to make sure all the threads are executed than using wait
	//time.Sleep(5000)
	//waitGroup.Wait()
	//
	//close(churnMetricsChannel)

	return churnMetricsArr, nil
}

func allocateFiles(ipChannel chan FileDeletedLines, fileDeletedLinesMap map[string][]int) {
	for filePath, deletedLines := range fileDeletedLinesMap {
		ipChannel <- FileDeletedLines{filePath, deletedLines}
	}
	close(ipChannel)
}

// processChurnMetrics appends each churnMetrics from churnMetricsChannel into churnMetricsArr
func processChurnMetrics(churnMetricsChannel chan ChurnMetrics, churnMetricsArr *[]ChurnMetrics, done chan bool) {
	// wait to make sure all the threads have completed execution and added the churnMetrics details into churnMetricsChannel
	// TODO: Research a better method to make sure all the threads are executed than using wait
	//helper.INFO.Println("INSIDE : processChurnMetrics")
	time.Sleep(500)
	for {
		churnMetrics, ok := <-churnMetricsChannel
		if !ok {
			break
		}
		*churnMetricsArr = append(*churnMetricsArr, churnMetrics)
		//wg.Done()
	}
	done <- true
}

// getChurnMetrics adds ChurnMetrics with churn details and count into the churnMetricsChannel for the specified deleted lines
func getChurnMetrics(ipFileChannel chan FileDeletedLines, churnMetricsChannel chan ChurnMetrics, repo *git.Repository, parentCommitHash *plumbing.Hash, commitAuthor string, wg *sync.WaitGroup) {
	//helper.INFO.Println("INSIDE : getChurnMetrics")
	for file := range ipFileChannel {
		if len(file.deletedLines) != 0 {
			churnMetrics := new(ChurnMetrics)
			blame, err := gitfuncs.Blame(repo, parentCommitHash, file.file)
			if err == nil {
				lines := blame.Lines
				churnDetails := make(map[string]string)
				selfChurnCount := 0
				interactiveChurnCount := 0
				for _, deletedLine := range file.deletedLines {
					churnAuthor := lines[deletedLine-1].Author
					if churnAuthor == commitAuthor {
						selfChurnCount += 1
					} else {
						interactiveChurnCount += 1
					}
					churnDetails[lines[deletedLine-1].Hash.String()] = churnAuthor
				}
				churnMetrics.DeletedLinesCount = len(file.deletedLines)
				churnMetrics.SelfChurnCount = selfChurnCount
				churnMetrics.InteractiveChurnCount = interactiveChurnCount
				churnMetrics.ChurnDetails = churnDetails
				churnMetrics.FilePath = file.file
				churnMetricsChannel <- *churnMetrics
			}
		}
	}
	wg.Done()
}

func AggrChurnMetrics(repo *git.Repository, baseCommitId string, parentCommitId string, aggregate string, whitespace bool, jsonOPToFile bool, printOP bool, filepath string) interface{} {
	defer helper.Duration(helper.Track("AggrChurnMetrics"))
	churnMetricsArr, _ := GetChurnMetrics(repo, baseCommitId, filepath, parentCommitId, whitespace, false, false)
	var aggChurnMetricsOutput interface{}
	if printOP {
		fmt.Println("  ______  _              ______  _                          \n / _____)(_) _          / _____)| |                         \n| /  ___  _ | |_   ___ | /      | | _   _   _   ____  ____  \n| | (___)| ||  _) (___)| |      | || \\ | | | | / ___)|  _ \\ \n| \\____/|| || |__      | \\_____ | | | || |_| || |    | | | |\n \\_____/ |_| \\___)      \\______)|_| |_| \\____||_|    |_| |_|\n                                                            ")
		PrintInBlue("")
	}
	switch aggregate {
	case "commit":
		aggChurnMetricsOutput = getCommitAggChurnMetrics(churnMetricsArr, printOP)
	case "all":
		aggChurnMetricsOutput = getAllAggChurnMetrics(churnMetricsArr, printOP)
	}
	if jsonOPToFile {
		writeJsonToFile(aggChurnMetricsOutput)
	}
	return aggChurnMetricsOutput
}

// getAllAggChurnMetrics aggregates all the churn count from all the commit aggregated churn metrics and returns AggAllChurnMetricsOutput
func getAllAggChurnMetrics(churnMetricsOutput *ChurnMetricsOutput, printOP bool) interface{} {
	var aggChurnMetricsOP = getCommitAggChurnMetrics(churnMetricsOutput, false)
	var aggAllChurnMetricsOutput AggAllChurnMetricsOutput
	aggAllChurnMetricsOutput.BaseCommitId = churnMetricsOutput.BaseCommitId
	if len(churnMetricsOutput.CommitDetails) > 0 {
		aggAllChurnMetricsOutput.ParentCommitId = churnMetricsOutput.CommitDetails[len(churnMetricsOutput.CommitDetails)-1].CommitId
		var commitsCount, totalDeletedLinesCount, totalSelfChurnCount, totalInteractiveChurnCount int
		for _, aggCommitDetails := range aggChurnMetricsOP.(AggChurnMetricsOutput).AggCommitDetails {
			commitsCount++
			totalDeletedLinesCount += aggCommitDetails.AggChurnMetrics.TotalDeletedLinesCount
			totalSelfChurnCount += aggCommitDetails.AggChurnMetrics.TotalSelfChurnCount
			totalInteractiveChurnCount += aggCommitDetails.AggChurnMetrics.TotalInteractiveChurnCount
		}
		aggAllChurnMetricsOutput.TotalCommits = commitsCount
		aggAllChurnMetricsOutput.TotalInteractiveChurnCount = totalInteractiveChurnCount
		aggAllChurnMetricsOutput.TotalSelfChurnCount = totalSelfChurnCount
		aggAllChurnMetricsOutput.TotalDeletedLinesCount = totalDeletedLinesCount

		if printOP {
			PrintInYellow("Base Commit ID: " + aggAllChurnMetricsOutput.BaseCommitId)
			PrintInPink("Parent Commit ID: " + aggAllChurnMetricsOutput.ParentCommitId)
			PrintInBlue("Total Commits: " + strconv.Itoa(aggAllChurnMetricsOutput.TotalCommits))
			fmt.Println("\tTotal Deleted lines count: " + strconv.Itoa(aggAllChurnMetricsOutput.TotalDeletedLinesCount))
			fmt.Println("\tTotal Self Churn count: " + strconv.Itoa(aggAllChurnMetricsOutput.TotalSelfChurnCount))
			fmt.Println("\tTotal Interactive Churn count: " + strconv.Itoa(aggAllChurnMetricsOutput.TotalInteractiveChurnCount))
		}
	}
	return aggAllChurnMetricsOutput
}

// getCommitAggChurnMetrics loops through each commitDetail and aggregates the churn counts and returns AggChurnMetricsOutput with aggCommitDetailsArr. It prints each aggCommitDetails if printOP is true
func getCommitAggChurnMetrics(churnMetricsOutput *ChurnMetricsOutput, printOP bool) interface{} {
	var aggChurnMetricsOP AggChurnMetricsOutput
	aggChurnMetricsOP.BaseCommitId = churnMetricsOutput.BaseCommitId
	commitDetailsArr := churnMetricsOutput.CommitDetails
	aggCommitDetailsArr := make([]AggCommitDetails, 0)

	//TODO: implement parallelism
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
