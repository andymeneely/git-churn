package main

import (
	"github.com/andymeneely/git-churn/cmd"
	"runtime"
)

func main() {
	//For executing the concurrent go routines in the program parallelly
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
	cmd.Execute()
}
