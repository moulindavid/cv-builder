package renderer

import (
	"bytes"
	"github.com/moulindavid/cv-builder/internal/resume"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderEscapesLatex(t *testing.T) {
	p := filepath.Join(t.TempDir(), "t.tex")
	os.WriteFile(p, []byte(`{{esc .Resume.Profile.Name}} {{text .Resume.Summaries}}`), 0600)
	r := resume.Resume{Profile: resume.Profile{Name: "D & M"}, Summaries: resume.Localized{"en": "Backend"}}
	var b bytes.Buffer
	if err := Render(&b, p, r, "en"); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(b.String(), `D \& M Backend`) {
		t.Fatalf("output=%q", b.String())
	}
}
