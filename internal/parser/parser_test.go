package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResumeFile(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "r.yaml")
	data := []byte("profile:\n  name: David\n  titles: {fr: Titre, en: Title}\ncontact: {email: d@example.com}\nsummaries: {fr: Résumé, en: Summary}\n")
	if err := os.WriteFile(p, data, 0600); err != nil {
		t.Fatal(err)
	}
	r, err := ResumeFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if r.Profile.Name != "David" {
		t.Fatalf("name=%q", r.Profile.Name)
	}
}
func TestUnknownField(t *testing.T) {
	p := filepath.Join(t.TempDir(), "r.yaml")
	os.WriteFile(p, []byte("unknown: true\n"), 0600)
	if _, err := ResumeFile(p); err == nil {
		t.Fatal("expected unknown field error")
	}
}
