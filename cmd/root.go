package cmd

import (
	"fmt"
	"github.com/andymeneely/git-churn/gitfuncs"
	"github.com/andymeneely/git-churn/helper"
	metrics "github.com/andymeneely/git-churn/metrics"
	"github.com/andymeneely/git-churn/print"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&repoUrl, "repo", "r", "", "Git Repository URL on which the churn metrics has to be computed")
	//print.CheckIfError(cobra.MarkFlagRequired(pf, "repo"))
	pf.StringVarP(&commitId, "commit", "c", "", "Commit hash for which the metrics has to be computed")
	//print.CheckIfError(cobra.MarkFlagRequired(pf, "commit"))
	pf.StringVarP(&filepath, "filepath", "f", "", "File path for the file on which the commit metrics has to be computed")
	pf.StringVarP(&aggregate, "aggregate", "a", "", "Aggregate the churn metrics. \"commit\": Aggregates all files in a commit. \"all\": Aggregate all files all commits and all files")
	pf.BoolVarP(&whitespace, "whitespace", "w", true, "Excludes whitespaces while calculating the churn metrics is set to false")
	pf.BoolVarP(&jsonOPToFile, "json", "j", false, "Writes the JSON output to a file within a folder named churn-details")
	pf.BoolVarP(&printOP, "print", "p", true, "Prints the output in a human readable format")
	pf.BoolVarP(&enableLog, "logging", "l", false, "Enables logging. Defaults to false")
}

var (
	repoUrl      string
	commitId     string
	filepath     string
	whitespace   bool
	jsonOPToFile bool
	printOP      bool
	aggregate    string
	enableLog    bool

	rootCmd = &cobra.Command{
		Use:   "git-churn",
		Short: "A fast tool for collecting code churn metrics from git repositories.",
		Long: `git-churn gives the churn metrics like insertions, deletions, etc for the given commit hash in the repo specified.
                Complete documentation is available at https://github.com/andymeneely/git-churn`,
		Run: func(cmd *cobra.Command, args []string) {
			if !enableLog {
				helper.INFO.SetFlags(0)
				helper.INFO.SetOutput(ioutil.Discard)
			}
			helper.INFO.Println("\n Processing new request\n")
			//var churnMetrics interface{}
			var err error
			commitIds := strings.Split(commitId, "..")
			firstCommitId := commitIds[0]
			var secondCommitId = ""
			if len(commitIds) == 2 {
				secondCommitId = strings.TrimFunc(commitIds[1], func(r rune) bool {
					return !unicode.IsLetter(r) && !unicode.IsNumber(r)
				})
			}
			if secondCommitId == "" {
				commitId = firstCommitId
				firstCommitId = ""
			} else {
				commitId = secondCommitId
			}

			if repoUrl == "" {
				repoUrl = "."
			}
			repo := gitfuncs.GetRepo(repoUrl)
			print.PrintInBlue(repoUrl + " " + commitId + " " + filepath + " " + firstCommitId)
			helper.INFO.Println("Generating git-churn for the following: \n" + "Repo:" + repoUrl + " " + " commitId:" + commitId + " " + " filepath:" + filepath + " " + " firstCommitId:" + firstCommitId)

			if aggregate == "" {
				_, err = metrics.GetChurnMetrics(repo, commitId, filepath, firstCommitId, whitespace, jsonOPToFile, printOP)
			} else {
				_ = metrics.AggrChurnMetrics(repo, commitId, firstCommitId, aggregate, whitespace, jsonOPToFile, printOP, filepath)
			}

			print.CheckIfError(err)

			//fmt.Println(fmt.Sprintf("%v", churnMetrics))

		},
	}
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of git-churn",
	Long:  `All software has versions. This is git-churn's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("git-churn version 0.1")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
