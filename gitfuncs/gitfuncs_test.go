package gitfuncs

import (
	"os"
	"strings"
	"testing"
)

func TestBranches(t *testing.T) {
	dir,_:= os.Getwd()
	rootDit := strings.Replace(dir,"gitfuncs","",1)
	got := Branches(rootDit)
	if len(got)<=0{
		t.Errorf("Branches(%q) == %q, want aleast 1 branch", rootDit, got)
	}
}
