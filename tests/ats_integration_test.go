package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/moulindavid/cv-builder/internal/ats"
	"github.com/moulindavid/cv-builder/internal/latex"
	"github.com/moulindavid/cv-builder/internal/parser"
	"github.com/moulindavid/cv-builder/internal/renderer"
	"github.com/moulindavid/cv-builder/internal/validator"
)

func root() string { _, f, _, _ := runtime.Caller(0); return filepath.Dir(filepath.Dir(f)) }
func TestPDFATSReadingOrder(t *testing.T) {
	if _, err := exec.LookPath("lualatex"); err != nil {
		t.Skip("lualatex not installed")
	}
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext not installed")
	}
	r, err := parser.ResumeFile(filepath.Join(root(), "data/resume.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := validator.Validate(r); err != nil {
		t.Fatal(err)
	}
	d := t.TempDir()
	tex := filepath.Join(d, "resume.tex")
	f, err := os.Create(tex)
	if err != nil {
		t.Fatal(err)
	}
	if err := renderer.Render(f, filepath.Join(root(), "templates/ats.tex"), r, "en"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	pdf := filepath.Join(d, "resume.pdf")
	if err := latex.Compile(tex, pdf); err != nil {
		t.Fatal(err)
	}
	text, err := ats.Extract(pdf)
	if err != nil {
		t.Fatal(err)
	}
	if err := ats.CheckOrder(text, "DAVID MOULIN", "SUMMARY", "EXPERIENCE", "TECHNICAL SKILLS", "EDUCATION"); err != nil {
		t.Fatalf("%v\nExtracted text:\n%s", err, text)
	}
	if !bytes.Contains([]byte(text), []byte("Benefiz")) {
		t.Fatal("Benefiz missing from extracted text")
	}
}
