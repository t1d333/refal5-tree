package cli

import "testing"

func Test_NewCompilerCLI(t *testing.T) {
	cli := NewCompilerCLI()

	if cli == nil {
		t.Errorf("Got nil CLI")
	}
}
