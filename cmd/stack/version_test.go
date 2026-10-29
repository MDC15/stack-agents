package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	stackagents "github.com/MDC15/stack-agents"
)

// binPath is the stack binary built once in TestMain.
var binPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "stack-cli-test")
	if err != nil {
		panic(err)
	}
	binPath = filepath.Join(tmp, "stack")

	build := exec.Command("go", "build", "-o", binPath, "./cmd/stack")
	// Tests run with cwd = cmd/stack; build from the repo root.
	build.Dir = "../.."
	if out, err := build.CombinedOutput(); err != nil {
		os.RemoveAll(tmp)
		panic("building stack binary: " + err.Error() + "\n" + string(out))
	}

	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}

// run executes the built binary and returns its exit code and combined
// output. Every exec gets an empty PATH and a temp HOME so agent
// detection finds nothing, whatever the test machine has installed.
func run(t *testing.T, args ...string) (int, string) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Env = []string{"PATH=", "HOME=" + t.TempDir()}
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	code := 0
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("running %v: %v", args, err)
		}
		code = exitErr.ExitCode()
	}
	return code, out.String()
}

func TestVersionOutput(t *testing.T) {
	data, err := fs.ReadFile(stackagents.Assets, "VERSION")
	if err != nil {
		t.Fatalf("reading VERSION from Assets: %v", err)
	}
	want := "stack " + strings.TrimSpace(string(data)) + "\n"

	for _, arg := range []string{"version", "--version"} {
		code, out := run(t, arg)
		if code != 0 {
			t.Errorf("stack %s: exit code = %d, want 0", arg, code)
		}
		if out != want {
			t.Errorf("stack %s: output = %q, want %q", arg, out, want)
		}
	}
}

func TestUnknownCommandExits2(t *testing.T) {
	code, _ := run(t, "bogus")
	if code != 2 {
		t.Errorf("stack bogus: exit code = %d, want 2", code)
	}
}
