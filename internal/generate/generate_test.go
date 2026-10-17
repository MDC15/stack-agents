package generate

import (
	"bytes"
	"testing"
)

func TestPlain(t *testing.T) {
	got := Plain([]byte("# rules\n"))
	want := []byte(Banner + "\n\n# rules\n")
	if !bytes.Equal(got, want) {
		t.Errorf("Plain([]byte(\"# rules\\n\")) = %q, want %q", got, want)
	}
}

func TestCursor(t *testing.T) {
	got := Cursor([]byte("body\n"))
	frontmatter := "---\ndescription: stack-agents operating rules\nalwaysApply: true\n---\n\n"

	if !bytes.HasPrefix(got, []byte(frontmatter)) {
		t.Fatalf("Cursor() must start with frontmatter %q, got %q", frontmatter, got)
	}
	if !bytes.Contains(got, []byte(Banner)) {
		t.Errorf("Cursor() must contain banner %q, got %q", Banner, got)
	}
	if !bytes.HasSuffix(got, []byte("body\n")) {
		t.Errorf("Cursor() must end with \"body\\n\", got %q", got)
	}
	if !bytes.Contains(got[len(frontmatter):], []byte(Banner)) {
		t.Errorf("Cursor() must contain banner %q after the frontmatter, got %q", Banner, got)
	}
}
