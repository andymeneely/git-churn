package gitfuncs

import (
	"os"
	"testing"
)

func TestBranches(t *testing.T) {
	dir,_:= os.Getwd()
	got := Branches(dir)
	if len(got)>0{
		t.Errorf("Branches(%q) == %q, want aleast 1 branch", dir, got)
	}
}
