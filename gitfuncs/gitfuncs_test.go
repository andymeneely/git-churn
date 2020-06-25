package gitfuncs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBranches(t *testing.T) {
	//dir,_:= os.Getwd()
	//rootDit := strings.Replace(dir,"gitfuncs","",1)
	branches := Branches("https://github.com/andymeneely/git-churn")
	assert := assert.New(t)
	assert.NotEqual(0, len(branches))
}

func TestRevList(t *testing.T) {
	r := Checkout("https://github.com/andymeneely/git-churn", "d78e64088e11bc2fd4f36f0421be91ebac52008c")
	commits, _ := RevList(r, "d78e64088e11bc2fd4f36f0421be91ebac52008c", "a8b24a74bae39a941186e11969f70058a351327d")
	assert := assert.New(t)
	assert.Equal(6, len(commits))
	assert.Equal("2aef1cf17578358b39eca71ac910e2d6882a3fd0", commits[1].Hash.String())
	assert.Equal("c0263662b2172b3df51ae39f8075dd010573ab6b", commits[5].Hash.String())
}

func TestRevList1(t *testing.T) {
	r := GetRepo("https://github.com/ashishgalagali/SWEN610-project")
	commits, _ := RevList(r, "c800ce62fc8a10d5fe69adb283f06296820522c1", "")
	assert := assert.New(t)
	assert.Equal(22, len(commits))
}

func TestGetDistinctAuthorsEMailIds(t *testing.T) {
	r := Checkout("https://github.com/andymeneely/git-churn", "d78e64088e11bc2fd4f36f0421be91ebac52008c")
	authors, _ := GetDistinctAuthorsEMailIds(r, "d78e64088e11bc2fd4f36f0421be91ebac52008c", "cbd945aa1ddff933ffe70802bb6905e77f014bc9", "README.md")
	assert := assert.New(t)
	assert.Equal(2, len(authors))
	assert.Equal([]string{"ashishgalagali@gmail.com", "andy@se.rit.edu"}, authors)

}

func TestNoFileGetDistinctAuthorsEMailIds(t *testing.T) {
	r := Checkout("https://github.com/andymeneely/git-churn", "d78e64088e11bc2fd4f36f0421be91ebac52008c")
	authors, _ := GetDistinctAuthorsEMailIds(r, "d78e64088e11bc2fd4f36f0421be91ebac52008c", "cbd945aa1ddff933ffe70802bb6905e77f014bc9", "test.md")
	assert := assert.New(t)
	assert.Equal(0, len(authors))
}
