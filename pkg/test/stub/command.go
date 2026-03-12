package stub

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func NewCommand(t *testing.T, testRunName string, envKey string) func(string, ...string) *exec.Cmd {
	t.Helper()
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{fmt.Sprintf("-test.run=^%s$", testRunName), "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{envKey + "=1"}
		return cmd
	}
}
