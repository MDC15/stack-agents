package stackagents

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

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
// link target exists next to the file that mentions it. The page list is
// enumerated from disk — never hardcoded, or a new page would ship
// unaccounted for.
func TestLinkIntegrity(t *testing.T) {
	if _, err := os.Stat("docs"); os.IsNotExist(err) {
		t.Skip("documentation not present (installed project), skipping")
	}
	pages := []string{"README.md", "README.vi.md", "CONTRIBUTING.md"}
	for _, dir := range []string{"docs/en", "docs/vi"} {
		for _, name := range docsDir(dir) {
			pages = append(pages, filepath.Join(dir, name))
		}
	}
	if !slicesContains(pages, filepath.Join("docs", "en", "cli.md")) {
		t.Errorf("cli.md missing from the enumerated pages; link integrity no longer covers it")
	}
	for _, path := range pages {
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
// filenames and heading counts as docs/en, by enumerating both
// directories — never a hardcoded page list, or a new page would ship
// without its translation.
func TestBilingualParity(t *testing.T) {
	if _, err := os.Stat("docs"); os.IsNotExist(err) {
		t.Skip("documentation not present (installed project), skipping")
	}
	en, vi := docsDir("docs/en"), docsDir("docs/vi")
	for _, fault := range missingMirrors(en, vi) {
		t.Errorf("%s", fault)
	}
	for _, name := range []string{"cli.md", "getting-started.md"} { // canaries: must be seen
		if !slicesContains(en, name) || !slicesContains(vi, name) {
			t.Errorf("%s not present in both docs/ trees, the enumeration misses it", name)
		}
	}
	for _, name := range intersect(en, vi) {
		enBody, _ := os.ReadFile(filepath.Join("docs/en", name))
		viBody, _ := os.ReadFile(filepath.Join("docs/vi", name))
		if got := len(headingPattern.FindAll(viBody, -1)); got != len(headingPattern.FindAll(enBody, -1)) {
			t.Errorf("docs/vi/%s has %d headings, docs/en/%s has %d", name, got, name, len(headingPattern.FindAll(enBody, -1)))
		}
	}
}

// docsDir lists the markdown files in dir, sorted.
func docsDir(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	names := []string{}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			names = append(names, e.Name())
		}
	}
	return names
}

// missingMirrors reports the fault lines: names present in en but absent
// from vi, and vice versa.
func missingMirrors(en, vi []string) []string {
	var faults []string
	for _, name := range en {
		if !slicesContains(vi, name) {
			faults = append(faults, "docs/vi lacks "+name)
		}
	}
	for _, name := range vi {
		if !slicesContains(en, name) {
			faults = append(faults, "docs/en lacks "+name)
		}
	}
	return faults
}

// TestMissingMirrors pins the fault finder the parity test rides on.
func TestMissingMirrors(t *testing.T) {
	if got := missingMirrors([]string{"a.md", "b.md"}, []string{"b.md"}); !slicesEqual(got, []string{"docs/vi lacks a.md"}) {
		t.Errorf("missingMirrors = %v, want [docs/vi lacks a.md]", got)
	}
	if got := missingMirrors([]string{"a.md", "b.md"}, []string{"b.md", "a.md"}); len(got) > 0 {
		t.Errorf("missingMirrors = %v, want none", got)
	}
	if got := missingMirrors([]string{"b.md"}, []string{"b.md", "stray.md"}); !slicesEqual(got, []string{"docs/en lacks stray.md"}) {
		t.Errorf("missingMirrors = %v, want [docs/en lacks stray.md]", got)
	}
}

// slicesEqual reports whether two string slices are equal.
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// intersect returns the names sorted in both lists.
func intersect(a, b []string) []string {
	bset := make(map[string]bool, len(b))
	for _, n := range b {
		bset[n] = true
	}
	out := []string{}
	for _, n := range a {
		if bset[n] {
			out = append(out, n)
		}
	}
	return out
}

// slicesContains reports whether s contains v.
func slicesContains(s []string, v string) bool {
	for _, n := range s {
		if n == v {
			return true
		}
	}
	return false
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
