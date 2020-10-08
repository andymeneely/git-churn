package main

import (
	"github.com/andymeneely/git-churn/cmd"
	"runtime"
)

//var cpuprofile = flag.String("cpuprofile", "defaultProf.out", "write cpu profile to `file`")
//var memprofile = flag.String("memprofile", "defaultMem.out", "write memory profile to `file`")

func main() {
	//flag.Parse()
	//if *cpuprofile != "" {
	//	f, err := os.Create(*cpuprofile)
	//	if err != nil {
	//		log.Fatal("could not create CPU profile: ", err)
	//	}
	//	defer f.Close()
	//	if err := pprof.StartCPUProfile(f); err != nil {
	//		log.Fatal("could not start CPU profile: ", err)
	//	}
	//	defer pprof.StopCPUProfile()
	//}
	//For executing the concurrent go routines in the program parallelly
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)
	cmd.Execute()
}
