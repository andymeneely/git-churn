package main

import (
	"fmt"
	//"gopkg.in/src-d/go-git.v4/plumbing/object"

	//. "gopkg.in/src-d/go-git.v4/_examples"
	"github.com/andymeneely/git-churn/stringutils"
	"github.com/google/go-cmp/cmp"
	)


func main() {
	fmt.Println("Namaste, world.")
	fmt.Println(stringutils.ReverseRunes("!oG ,olleH"))
	fmt.Println(cmp.Diff("Hello World", "Hello Go"))


}

