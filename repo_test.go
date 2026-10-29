package stackagents

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// docFiles arms both documentation checks: every file a link check
// walks, and the bilingual parity pair.
var docFiles = []string{
	"README.md",
	"README.vi.md",
	"CONTRIBUTING.md",
	"docs/en/custom-agents.md",
	"docs/en/extending.md",
	"docs/en/getting-started.md",
	"docs/en/how-it-works.md",
	"docs/en/skills.md",
	"docs/vi/custom-agents.md",
	"docs/vi/extending.md",
	"docs/vi/getting-started.md",
	"docs/vi/how-it-works.md",
	"docs/vi/skills.md",
}

// linkPattern matches markdown link targets inside (…).
var linkPattern = regexp.MustCompile(`\]\(([^)(]+)\)`)

// skipPattern matches prefixes that do not resolve on disk.
const skipPattern = "http://,https://,mailto:,#"

// skipLinkPrefixes returns the prefixes of link targets that never point
// at a file.
func skipLinkPrefixes() []string {
	return strings.Split(skipPattern, ",")
}

// bodyWithoutCode strips fenced code blocks and inline code spans so link
// syntax shown in examples is never mistaken for a real link.
func bodyWithoutCode(body string) string {
	fenced := regexp.MustCompile("(?s)```.*?```")
	body = fenced.ReplaceAllString(body, "")
	inline := regexp.MustCompile("`[^`\n]*`")
	return inline.ReplaceAllString(body, "")
}

// linksIn extracts the on-disk targets of a markdown file, keeping the
// anchor part only when the target is a same-file fragment.
func linksIn(t *testing.T, path string) []string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	body := bodyWithoutCode(string(raw))
	skips := skipLinkPrefixes()
	var targets []string
	for _, m := range linkPattern.FindAllStringSubmatch(body, -1) {
		target := m[1]
		skip := false
		for _, p := range skips {
			if strings.HasPrefix(target, p) {
				skip = true
				break
			}
		}
		if !skip {
			targets = append(targets, target)
		}
	}
	return targets
}

// TestLinkIntegrity walks every markdown page and asserts each relative
// link target exists next to the file that mentions it.
func TestLinkIntegrity(t *testing.T) {
	if _, err := os.Stat("docs"); os.IsNotExist(err) {
		t.Skip("documentation not present (installed project), skipping")
	}
	for _, path := range docFiles {
		dir := filepath.Dir(path)
		for _, target := range linksIn(t, path) {
			file := target
			if i := strings.Index(file, "#"); i >= 0 {
				file = file[:i]
			}
			if file == "" {
				continue
			}
			if _, err := os.Stat(filepath.Join(dir, file)); err != nil {
				t.Errorf("%s links to missing %s: %v", path, target, err)
			}
		}
	}
}

// headingPattern matches a markdown heading at the start of a line.
var headingPattern = regexp.MustCompile(`(?m)^#`)

// TestBilingualParity asserts the Vietnamese mirrors keep the same
// filenames and heading counts as docs/en.
func TestBilingualParity(t *testing.T) {
	if _, err := os.Stat("docs"); os.IsNotExist(err) {
		t.Skip("documentation not present (installed project), skipping")
	}
	for _, en := range []string{"custom-agents.md", "extending.md", "getting-started.md", "how-it-works.md", "skills.md"} {
		viPath := filepath.Join("docs/vi", en)
		enPath := filepath.Join("docs/en", en)
		enBody, err := os.ReadFile(enPath)
		if err != nil {
			t.Fatalf("reading %s: %v", enPath, err)
		}
		viBody, err := os.ReadFile(viPath)
		if err != nil {
			t.Errorf("docs/vi lacks %s", en)
			continue
		}
		if got := len(headingPattern.FindAll(viBody, -1)); got != len(headingPattern.FindAll(enBody, -1)) {
			t.Errorf("%s has %d headings, %s has %d", viPath, got, enPath, len(headingPattern.FindAll(enBody, -1)))
		}
	}
}

// loaderSkillNames is the shipped skill set the loader must list.
var loaderSkillNames = []string{"architect", "arena", "checkpoint", "conventions", "red-team", "tdd", "unslop", "verify"}

// TestLoaderCLI exercises the python loader the way a custom agent does.
func TestLoaderCLI(t *testing.T) {
	python, err := exec.LookPath("python3")
	if err != nil {
		t.Skip("python3 not available, skipping")
	}

	run := func(t *testing.T, args ...string) (string, int) {
		t.Helper()
		out, err := exec.Command(python, append([]string{"-I", "scripts/load_skill.py"}, args...)...).CombinedOutput()
		code := 0
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatalf("running loader: %v", err)
			}
			code = exitErr.ExitCode()
		}
		return string(out), code
	}

	out, code := run(t, "list")
	if code != 0 {
		t.Fatalf("loader list exits %d (%s)", code, out)
	}
	for _, name := range loaderSkillNames {
		if !strings.Contains(out, name) {
			t.Errorf("loader list output %q missing skill %q", out, name)
		}
	}

	out, code = run(t, "skills", "tdd")
	if code != 0 {
		t.Fatalf("loading the tdd skill exits %d (%s)", code, out)
	}
	if !strings.Contains(out, "tdd") {
		t.Errorf("loaded tdd skill output %q does not mention it", out)
	}

	if _, code = run(t, "skills", "../../etc/passwd"); code == 0 {
		t.Errorf("path traversal through the loader should fail")
	}
}
