package gitfuncs

import (
	"fmt"
	. "github.com/andymeneely/git-churn/print"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func LastCommit(repoPath string) string {

	r, err := git.PlainOpen(repoPath)
	CheckIfError(err)

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)

	return commit.Message
}

func Branches(repoPath string) []string {
	r, err := git.PlainOpen(repoPath)
	CheckIfError(err)

	branchIttr, _ := r.Branches()

	fmt.Println(branchIttr)
	var branches []string
	err = branchIttr.ForEach(func(ref *plumbing.Reference) error  {
		fmt.Println(ref.Name().String())
		branches = append(branches,ref.Name().String())
		return nil
	})

	return branches
}