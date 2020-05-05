package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/andymeneely/git-churn/gitfuncs"
	metrics "github.com/andymeneely/git-churn/matrics"
	"github.com/andymeneely/git-churn/print"
	"github.com/spf13/cobra"
	"os"
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
	//TODO: Add whitespace exclusion bool flag
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
			repo := gitfuncs.Checkout(repoUrl, commitId)
			if whitespace {
				if filepath != "" {
					churnMetrics = metrics.GetChurnMetricsWithWhitespace(repo, filepath)
				} else {
					churnMetrics = metrics.AggrDiffMetricsWithWhitespace(repo)
				}
			} else {
				if filepath != "" {
					churnMetrics = metrics.GetChurnMetricsWhitespaceExcluded(repo, filepath)
				} else {
					churnMetrics, err = metrics.AggrDiffMetricsWhitespaceExcluded(repo)
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
