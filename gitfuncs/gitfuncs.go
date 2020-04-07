package gitfuncs

import (
	//"fmt"
	. "github.com/andymeneely/git-churn/print"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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

	//fmt.Println(commit)

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

	//fmt.Println(branchIttr)
	var branches []string
	//TODO: Check why it is only getting the master branch
	err = branchIttr.ForEach(func(ref *plumbing.Reference) error {
		//fmt.Println(ref.Name().String())
		branches = append(branches, ref.Name().String())
		return nil
	})

	return branches
}

func Tags(repoUrl string) []*plumbing.Reference {

	Info("git clone " + repoUrl)

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoUrl,
	})

	CheckIfError(err)
	// List all tag references, both lightweight tags and annotated tags
	Info("git show-ref --tag")
	var tagsArr []*plumbing.Reference

	tagrefs, err := r.Tags()
	CheckIfError(err)
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		tagsArr = append(tagsArr, t)
		return nil
	})
	CheckIfError(err)

	return tagsArr

}

func Checkout(repoUrl, hash string) *git.Repository {
	Info("git clone " + repoUrl)

	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: repoUrl,
	})

	CheckIfError(err)
	// ... retrieving the commit being pointed by HEAD

	Info("git show-ref --head HEAD")
	//ref, err := r.Head()
	//CheckIfError(err)
	//fmt.Println(ref.Hash())

	w, err := r.Worktree()
	CheckIfError(err)

	// ... checking out to commit
	Info("git checkout %s", hash)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(hash),
	})
	CheckIfError(err)

	// ... retrieving the commit being pointed by HEAD, it shows that the
	// repository is pointing to the giving commit in detached mode
	Info("git show-ref --head HEAD")
	//ref, err = r.Head()
	//CheckIfError(err)
	//fmt.Println(ref.Hash())

	return r
}

func FileLOC(repoUrl, filePath string) int {
	loc := 0
	// ... get the files iterator and print the file
	FilesIttr(repoUrl).ForEach(func(f *object.File) error {
		if f.Name == filePath {
			lines, _ := f.Lines()
			loc = len(lines)
		}
		return nil
	})
	return loc
}

func FileLOCFromTree(tree *object.Tree, filePath string) int {
	loc := 0
	tree.Files().ForEach(func(f *object.File) error {
		if f.Name == filePath {
			lines, _ := f.Lines()
			loc = len(lines)
		}
		return nil
	})
	return loc
}

func LOCFilesFromTree(tree *object.Tree) (int, []string) {
	loc := 0
	var files []string
	tree.Files().ForEach(func(f *object.File) error {
		lines, _ := f.Lines()
		loc += len(lines)
		files = append(files, f.Name)
		return nil
	})
	return loc, files
}

func FileLOCFromTreeWhitespaceExcluded(tree *object.Tree, filePath string) int {
	loc := 0
	tree.Files().ForEach(func(f *object.File) error {
		if f.Name == filePath {
			lines, _ := f.Lines()
			for _, line := range lines {
				if line != "" {
					loc += 1
				}
			}
		}
		return nil
	})
	return loc
}
func LOCFilesFromTreeWhitespaceExcluded(tree *object.Tree) (int, []string) {
	loc := 0
	var files []string
	tree.Files().ForEach(func(f *object.File) error {
		lines, _ := f.Lines()
		for _, line := range lines {
			if line != "" {
				loc += 1
			}
		}
		files = append(files, f.Name)
		return nil
	})
	return loc, files
}

func FilesIttr(repoUrl string) *object.FileIter {
	//REF: https://github.com/src-d/go-git/blob/master/_examples/showcase/main.go
	//Clones the given repository in memory, creating the remote, the local
	//branches and fetching the objects, exactly as:
	Info("git clone " + repoUrl)

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoUrl,
	})

	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)
	//fmt.Println(commit)

	// List the tree from HEAD
	Info("git ls-tree -r HEAD")

	// ... retrieve the tree from the commit
	tree, err := commit.Tree()
	CheckIfError(err)

	return tree.Files()
}

func CommitDiff(repoUrl, commitId string) (*object.Changes, *object.Tree, *object.Tree) {
	repo := Checkout(repoUrl, commitId)

	head, err := repo.Head()
	CheckIfError(err)

	commitObj, err := repo.CommitObject(head.Hash())
	CheckIfError(err)

	parentCommitObj, err := commitObj.Parent(0)
	CheckIfError(err)

	// List the tree from HEAD
	Info("git ls-tree -repo HEAD")

	// ... retrieve the tree from the commit
	tree, err := commitObj.Tree()
	CheckIfError(err)

	parentTree, err := parentCommitObj.Tree()
	CheckIfError(err)
	changes, err := parentTree.Diff(tree)
	CheckIfError(err)

	//fmt.Println(changes)
	//fmt.Println(changes.Patch())
	return &changes, tree, parentTree
}
