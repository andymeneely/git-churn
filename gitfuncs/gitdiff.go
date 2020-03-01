package gitfuncs

import (
	"gopkg.in/src-d/go-git.v4"
	. "github.com/andymeneely/git-churn/print"


)

func Diff(repoPath string) string {

	r, err := git.PlainOpen(repoPath)
	CheckIfError(err)


}
