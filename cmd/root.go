package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/andymeneely/git-churn/gitfuncs"
	metrics "github.com/andymeneely/git-churn/matrics"
	"github.com/andymeneely/git-churn/print"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"unicode"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	pf := rootCmd.PersistentFlags()
	pf.StringVarP(&repoUrl, "repo", "r", "", "Git Repository URL on which the churn metrics has to be computed")
	print.CheckIfError(cobra.MarkFlagRequired(pf, "repo"))
	pf.StringVarP(&commitId, "commit", "c", "", "Commit hash for which the metrics has to be computed")
	print.CheckIfError(cobra.MarkFlagRequired(pf, "commit"))
	pf.StringVarP(&filepath, "filepath", "f", "", "File path for the file on which the commit metrics has to be computed")
	pf.BoolVarP(&whitespace, "whitespace", "w", true, "Excludes whitespaces while calculating the churn metrics is set to false")
}

var (
	repoUrl    string
	commitId   string
	filepath   string
	whitespace bool

	rootCmd = &cobra.Command{
		Use:   "git-churn",
		Short: "A fast tool for collecting code churn metrics from git repositories.",
		Long: `git-churn gives the churn metrics like insertions, deletions, etc for the given commit hash in the repo specified.
                Complete documentation is available at https://github.com/andymeneely/git-churn`,
		Run: func(cmd *cobra.Command, args []string) {
			var churnMetrics interface{}
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

			repo := gitfuncs.GetRepo(repoUrl)
			if whitespace {
				if filepath != "" {
					churnMetrics, err = metrics.GetChurnMetricsWithWhitespace(repo, commitId, filepath, firstCommitId)
				} else {
					churnMetrics = metrics.AggrChurnMetricsWithWhitespace(repo, commitId)
				}
			} else {
				if filepath != "" {
					churnMetrics, err = metrics.GetChurnMetricsWhitespaceExcluded(repo, commitId, filepath, firstCommitId)
				} else {
					churnMetrics = metrics.AggrChurnMetricsWhitespaceExcluded(repo, commitId)
				}
				print.CheckIfError(err)
			}
			//fmt.Println(fmt.Sprintf("%v", churnMetrics))
			out, err := json.Marshal(churnMetrics)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(out))
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
