package gitfuncs

import (
	"testing"
)

func TestBranches(t *testing.T) {
	//dir,_:= os.Getwd()
	//rootDit := strings.Replace(dir,"gitfuncs","",1)
	got := Branches("https://github.com/andymeneely/git-churn")
	if len(got)<=0{
		t.Errorf("Branches(%q) == %q, want aleast 1 branch", "https://github.com/andymeneely/git-churn", got)
	}
}
