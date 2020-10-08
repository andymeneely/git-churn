package gitfuncs

import (
	"github.com/andymeneely/git-churn/helper"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/revlist"
	"github.com/go-git/go-git/v5/storage/memory"
	"sort"
	"strings"

	. "github.com/andymeneely/git-churn/print"
)

func LastCommit(repoUrl string) string {
	// Clones the given repository in memory, creating the remote, the local
	// branches and fetching the objects, exactly as:
	//PrintInBlue("git clone " + repoUrl)

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
	//PrintInBlue("git clone " + repoUrl)

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoUrl,
	})
	//r, err := git.PlainOpen("/Users/raj.g/Documents/RA/git-churn")

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

	//PrintInBlue("git clone " + repoUrl)

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoUrl,
	})

	CheckIfError(err)
	// List all tag references, both lightweight tags and annotated tags
	//PrintInBlue("git show-ref --tag")
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
	//PrintInBlue("git clone " + repoUrl)

	r := GetRepo(repoUrl)
	w, err := r.Worktree()
	CheckIfError(err)

	// ... checking out to commit
	//PrintInBlue("git checkout %s", hash)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(hash),
	})
	CheckIfError(err)
	return r
}

func GetRepo(repoUrl string) *git.Repository {
	defer helper.Duration(helper.Track("GetRepo"))

	//PrintInBlue("git clone " + repoUrl)

	var r *git.Repository
	var err error
	//if strings.HasPrefix(repoUrl, "https://github.com") {
	r, err = git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: repoUrl,
	})
	CheckIfError(err)
	//} else {
	//	r, err = git.PlainOpen(repoUrl)
	//	CheckIfError(err)
	//}
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

//Gets the total number of lines of code in a given file in the specified commit tree
//Whitespace included
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

//Returns the total lines of code from all the files in the given commit tree and list of fine names
// Whitespace included
func LOCFilesFromTree(tree *object.Tree, c chan func() (int, []string)) {
	loc := 0
	var files []string
	tree.Files().ForEach(func(f *object.File) error {
		lines, _ := f.Lines()
		loc += len(lines)
		files = append(files, f.Name)
		return nil
	})
	c <- func() (int, []string) { return loc, files }
}

//Gets the total number of lines of code in a given file in the specified commit tree
//Whitespace excluded
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

//Returns the total lines of code from all the files in the given commit tree and list of fine names
//Whitespace excluded
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
	//PrintInBlue("git clone " + repoUrl)

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
	//PrintInBlue("git ls-tree -r HEAD")

	// ... retrieve the tree from the commit
	tree, err := commit.Tree()
	CheckIfError(err)

	return tree.Files()
}

// Returns the changes b/n the commit and it's parent, the tree corresponding to the commit and it's parent tree
func CommitDiff(repo *git.Repository, baseCommitId string, parentCommitHash *plumbing.Hash) (*object.Changes, *object.Tree, *object.Tree) {

	head := plumbing.NewHash(baseCommitId)

	commitObj, err := repo.CommitObject(head)
	CheckIfError(err)
	//fmt.Println(commitObj.Author.Name)
	//fmt.Println(commitObj.Author.Email)
	//fmt.Println(commitObj.Author.When)
	//fmt.Println(commitObj.Author.String())

	//if parentCommitHash
	var parentCommitObj *object.Commit

	if parentCommitHash == nil {
		parentCommitObj, err = commitObj.Parent(0)
		CheckIfError(err)
	} else {
		parentCommitObj, err = repo.CommitObject(*parentCommitHash)
		CheckIfError(err)
	}
	//helper.INFO.Println("Getting the commit diff between " + baseCommitId + " & " + parentCommitObj.Hash.String())

	// List the tree from HEAD
	//PrintInBlue("git ls-tree -repo " + parentCommitObj.Hash.String())

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

// DeletedLineNumbers returns the map of file and the list of deleted lines
// 	changes: the commitDiff changes
//	filePath: if present, returns deleted lies only for that filepath
//	whitespace: if false, neglects the blank deleted lines
func DeletedLineNumbers(changes *object.Changes, filePath string, whitespace bool) map[string][]int {
	//helper.INFO.Println("Getting deleted lines for  " + strconv.Itoa(changes.Len()) + " changes")
	//changes, _, parentTree := CommitDiff(repo, parentCommitHash)
	patch, _ := changes.Patch()
	fileDeletedLinesMap := make(map[string][]int)
	for _, patch := range patch.FilePatches() {
		//fmt.Println(patch)
		lineCounter := 0
		var deletedLines []int
		fromFile, toFile := patch.Files()
		var file string
		if nil == fromFile {
			file = toFile.Path()
		} else {
			file = fromFile.Path()
		}
		if filePath != "" && file != filePath {
			continue
		}
		for _, chunk := range patch.Chunks() {
			deletedPatch := strings.Split(chunk.Content(), "\n")
			if chunk.Type() == 0 {
				if chunk.Content()[len(chunk.Content())-1] == '\n' {
					lineCounter += len(deletedPatch) - 1
				} else {
					lineCounter += len(deletedPatch)
				}
				//lineCounter += len(strings.Split(chunk.Content(), "\n")) - 1
			}
			if chunk.Type() == 2 {
				var patchLen int
				if chunk.Content()[len(chunk.Content())-1] == '\n' {
					patchLen = len(deletedPatch) - 1
				} else {
					patchLen = len(deletedPatch)
				}
				if whitespace {
					for i := 1; i <= patchLen; i++ {
						deletedLines = append(deletedLines, lineCounter+i)
					}
				} else {
					for i := 1; i <= patchLen; i++ {
						if deletedPatch[i-1] != "" {
							deletedLines = append(deletedLines, lineCounter+i)
						}
					}
				}
				lineCounter += patchLen
			}
		}
		fileDeletedLinesMap[file] = deletedLines

	}
	return fileDeletedLinesMap
}

func DeletedLineNumbersWhitespaceExcluded(changes *object.Changes) map[string][]int {
	patch, _ := changes.Patch()
	fileDeletedLinesMap := make(map[string][]int)
	for _, patch := range patch.FilePatches() {
		//fmt.Println(patch)
		lineCounter := 0
		var deletedLines []int
		for _, chunk := range patch.Chunks() {
			if chunk.Type() == 0 {
				if chunk.Content()[len(chunk.Content())-1] == '\n' {
					lineCounter += len(strings.Split(chunk.Content(), "\n")) - 1
				} else {
					lineCounter += len(strings.Split(chunk.Content(), "\n"))
				}
			}
			if chunk.Type() == 2 {
				deletedPatch := strings.Split(chunk.Content(), "\n")
				var patchLen int
				if chunk.Content()[len(chunk.Content())-1] == '\n' {
					patchLen = len(deletedPatch) - 1
				} else {
					patchLen = len(deletedPatch)
				}
				for i := 1; i <= patchLen; i++ {
					if deletedPatch[i-1] != "" {
						deletedLines = append(deletedLines, lineCounter+i)
					}
				}
				lineCounter += patchLen
			}
		}
		fromFile, toFile := patch.Files()
		if nil == fromFile {
			fileDeletedLinesMap[toFile.Path()] = deletedLines
		} else {
			fileDeletedLinesMap[fromFile.Path()] = deletedLines
		}
		//fmt.Println(deletedLines)
	}
	return fileDeletedLinesMap
}

// RevisionCommits returns the hash of the specified revision. If revision is empty then returns hash of baseCommitId~1
func RevisionCommits(r *git.Repository, baseCommitId, revision string) *plumbing.Hash {

	// Resolve revision into a sha1 commit, only some revisions are resolved
	// look at the doc to get more details
	if revision == "" {
		revision = baseCommitId + "~1"
	}
	//helper.INFO.Println("Getting revision commit hash for " + revision)
	//PrintInBlue("git rev-parse %s", revision)
	h, err := r.ResolveRevision(plumbing.Revision(revision))
	CheckIfError(err)
	return h
}

// RevList is native implementation of git rev-list command
// 	Returns list of commit objects between the given commit hash sorted by date
func RevList(r *git.Repository, beginCommit, endCommit string) ([]*object.Commit, error) {
	//TODO: should I reverse the begin and end?
	//helper.INFO.Println("Getting the RevList of commits between " + beginCommit + " and " + endCommit)
	commits := make([]*object.Commit, 0)
	var ref1hist []plumbing.Hash
	var err error
	if endCommit != "" {
		ref1hist, err = revlist.Objects(r.Storer, []plumbing.Hash{plumbing.NewHash(endCommit)}, nil)
		if err != nil {
			return nil, err
		}
	}
	ref2hist, err := revlist.Objects(r.Storer, []plumbing.Hash{plumbing.NewHash(beginCommit)}, ref1hist)
	if err != nil {
		return nil, err
	}

	for _, h := range ref2hist {
		c, err := r.CommitObject(h)
		if err != nil {
			continue
		}
		commits = append(commits, c)
	}
	//  sorts by datetime
	sort.Slice(commits, func(i, j int) bool { return commits[i].Committer.When.Unix() > commits[j].Committer.When.Unix() })
	//fmt.Println(commits)
	//helper.INFO.Println("Got " + strconv.Itoa(len(commits)) + " commits to compute the churn-metrics")

	return commits, err
}

func GetDistinctAuthorsEMailIds(r *git.Repository, beginCommit, endCommit, filePath string) ([]string, error) {

	commits, err := RevList(r, beginCommit, endCommit)
	if err != nil {
		return nil, err
	}

	var authors []string
	for _, commit := range commits {
		tree, err := commit.Tree()
		if err != nil {
			return nil, err
		}
		_, err = tree.File(filePath)
		if err != nil {
			continue
		}
		authors = append(authors, commit.Author.Email)
	}
	authors = helper.UniqueElements(authors)
	return authors, err

}

// Blame returns git.BlameResult (Last author who changed each line) for the the given has and file path
func Blame(repo *git.Repository, hash *plumbing.Hash, path string) (*git.BlameResult, error) {

	//TODO: It does not support the options mentioned in https://git-scm.com/docs/git-blame
	commitObj, err := repo.CommitObject(*hash)
	CheckIfError(err)

	// This is because the Blame throws error if the previous commit is a merge PR commit
	//if strings.Contains(commitObj.Message, "Merge pull request") {
	//	hash := RevisionCommits(repo, "HEAD^^2")
	//	commitObj, err = repo.CommitObject(*hash)
	//	fmt.Println(commitObj.Message)
	//	CheckIfError(err)
	//}

	//helper.INFO.Println("Getting Blame results for commit: " + commitObj.Hash.String() + " and file: " + path)
	//TODO: issue: https://github.com/src-d/go-git/issues/725
	blameResult, err := git.Blame(commitObj, path)

	//fmt.Println(blameResult)
	//fmt.Println(blameResult.Lines)

	return blameResult, err

}
