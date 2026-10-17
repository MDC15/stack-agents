package detect

import (
	"testing"

	"github.com/MDC15/stack-agents/internal/agents"
)

// fakeFinder is a Finder with canned answers so tests of All never shell out.
type fakeFinder struct {
	path     map[string]bool
	home     map[string]bool
	exts     []string
	extsOK   bool
	versions map[string]string
}

func (f *fakeFinder) InPath(bin string) bool       { return f.path[bin] }
func (f *fakeFinder) HomeExists(rel string) bool   { return f.home[rel] }
func (f *fakeFinder) VSCodeExts() ([]string, bool) { return f.exts, f.extsOK }
func (f *fakeFinder) Version(bin string) string    { return f.versions[bin] }

func TestAll(t *testing.T) {
	tests := []struct {
		name   string
		finder *fakeFinder
		want   map[string]bool // agent ID -> Found
	}{
		{
			name:   "nothing present",
			finder: &fakeFinder{},
			want: map[string]bool{
				"claude": false, "opencode": false, "codex": false, "cursor": false,
				"copilot": false, "cline": false, "roo": false, "windsurf": false,
			},
		},
		{
			name:   "claude in PATH",
			finder: &fakeFinder{path: map[string]bool{"claude": true}},
			want:   map[string]bool{"claude": true},
		},
		{
			name:   "claude home dir alone",
			finder: &fakeFinder{home: map[string]bool{".claude": true}},
			want:   map[string]bool{"claude": true},
		},
		{
			name:   "opencode in PATH",
			finder: &fakeFinder{path: map[string]bool{"opencode": true}},
			want:   map[string]bool{"opencode": true},
		},
		{
			name:   "opencode home dir",
			finder: &fakeFinder{home: map[string]bool{".config/opencode": true}},
			want:   map[string]bool{"opencode": true},
		},
		{
			name:   "codex in PATH",
			finder: &fakeFinder{path: map[string]bool{"codex": true}},
			want:   map[string]bool{"codex": true},
		},
		{
			name:   "codex home dir",
			finder: &fakeFinder{home: map[string]bool{".codex": true}},
			want:   map[string]bool{"codex": true},
		},
		{
			name:   "cursor in PATH",
			finder: &fakeFinder{path: map[string]bool{"cursor": true}},
			want:   map[string]bool{"cursor": true},
		},
		{
			name:   "cursor home dir",
			finder: &fakeFinder{home: map[string]bool{".cursor": true}},
			want:   map[string]bool{"cursor": true},
		},
		{
			name:   "copilot extension case-insensitive",
			finder: &fakeFinder{exts: []string{"GitHub.copilot"}, extsOK: true},
			want:   map[string]bool{"copilot": true},
		},
		{
			name:   "cline extension",
			finder: &fakeFinder{exts: []string{"saoudrizwan.claude-dev"}, extsOK: true},
			want:   map[string]bool{"cline": true},
		},
		{
			name:   "roo extension",
			finder: &fakeFinder{exts: []string{"rooveterinaryinc.roo-cline"}, extsOK: true},
			want:   map[string]bool{"roo": true},
		},
		{
			name:   "windsurf in PATH",
			finder: &fakeFinder{path: map[string]bool{"windsurf": true}},
			want:   map[string]bool{"windsurf": true},
		},
		{
			name:   "windsurf home dir",
			finder: &fakeFinder{home: map[string]bool{".codeium": true}},
			want:   map[string]bool{"windsurf": true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := All(tc.finder)

			if len(got) != len(agents.All()) {
				t.Fatalf("All() returned %d results, want %d", len(got), len(agents.All()))
			}
			for i, a := range agents.All() {
				if got[i].Agent != a {
					t.Errorf("result %d agent = %+v, want %+v", i, got[i].Agent, a)
				}
				want := tc.want[a.ID]
				if got[i].Found != want {
					t.Errorf("%s: Found = %v, want %v", a.ID, got[i].Found, want)
				}
			}
		})
	}
}

func TestAllVersionWhenFound(t *testing.T) {
	f := &fakeFinder{
		path:     map[string]bool{"claude": true},
		versions: map[string]string{"claude": "2.4.1"},
	}

	for _, r := range All(f) {
		switch r.Agent.ID {
		case "claude":
			if r.Version != "2.4.1" {
				t.Errorf("claude: Version = %q, want %q", r.Version, "2.4.1")
			}
		default:
			if r.Found {
				t.Errorf("%s: Found = true, want false in this scenario", r.Agent.ID)
			}
		}
	}
}

func TestAllVersionNeverFilledWhenNotFound(t *testing.T) {
	// The finder would report a version for claude, but nothing is found.
	f := &fakeFinder{versions: map[string]string{"claude": "2.4.1"}}

	for _, r := range All(f) {
		if r.Found {
			t.Fatalf("%s: Found = true, want false", r.Agent.ID)
		}
		if r.Version != "" {
			t.Errorf("%s: Version = %q, want empty for a not-found agent", r.Agent.ID, r.Version)
		}
	}
}

func TestAllVSCodeExtsUnavailable(t *testing.T) {
	// VSCodeExts returns (nil, false) when the code binary is absent.
	f := &fakeFinder{}

	got := All(f) // must not panic

	for _, r := range got {
		switch r.Agent.ID {
		case "copilot", "cline", "roo":
			if r.Found {
				t.Errorf("%s: Found = true, want false when VSCodeExts is unavailable", r.Agent.ID)
			}
			if r.Version != "" {
				t.Errorf("%s: Version = %q, want empty", r.Agent.ID, r.Version)
			}
		}
	}
}
