package main

import (
	"fmt"
	"github.com/andymeneely/git-churn/gitfuncs"
	"os"
)


func main() {
	dir,_:= os.Getwd()
	fmt.Println(gitfuncs.LastCommit(dir))
	fmt.Println(gitfuncs.Branches("https://github.com/andymeneely/git-churn"))

}