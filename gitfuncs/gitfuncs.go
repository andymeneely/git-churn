package gitfuncs

import (
	"fmt"
	. "github.com/andymeneely/git-churn/print"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func LastCommit(repoUrl string) string {
	// Clones the given repository in memory, creating the remote, the local
	// branches and fetching the objects, exactly as:
	Info("git clone " + repoUrl)

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoUrl,
	})

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

func Branches(repoUrl string) []string {
	// Clones the given repository in memory, creating the remote, the local
	// branches and fetching the objects, exactly as:
	Info("git clone " + repoUrl)

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoUrl,
	})

	CheckIfError(err)

	branchIttr, _ := r.Branches()

	fmt.Println(branchIttr)
	var branches []string
	//TODO: Check why it is only getting the master branch
	err = branchIttr.ForEach(func(ref *plumbing.Reference) error {
		//fmt.Println(ref.Name().String())
		branches = append(branches, ref.Name().String())
		return nil
	})

	return branches
}
